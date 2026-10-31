import { useEffect, useId, useState } from "react"
import {
  Controller,
  useForm,
  useWatch,
  type Control,
  type DefaultValues,
  type FieldPath,
  type PathValue,
  type Resolver,
  type SubmitHandler,
} from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslation } from "react-i18next"
import type { z } from "zod"
import { Dialog } from "../../../../shared/components/ui/dialog"
import { DiscardChangesDialog } from "../../../../shared/components/ui/discard-changes-dialog"
import { FieldLegend, FieldSet } from "../../../../shared/components/ui/field"
import { Input } from "../../../../shared/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../../../../shared/components/ui/select"
import { Textarea } from "../../../../shared/components/ui/textarea"
import { useDiscardGuard } from "../../../../shared/hooks/useDiscardGuard"
import { cn } from "../../../../shared/utils/cn"
import { FormDialogLayout } from "../form/FormDialogLayout"
import { FormErrorSummary, type FormErrorSummaryItem } from "../form/FormErrorSummary"
import { FormField, type FormFieldControlProperties } from "../form/FormField"
import { FormSubmitError } from "../form/FormSubmitError"
import { useSubmissionLock } from "../form/useSubmissionLock"

const ERROR_SUMMARY_MIN_FIELDS = 5

export interface ClinicalFormOption {
  value: string
  labelKey?: string
  label?: string
}

export interface ClinicalFormField<FormData extends Record<string, unknown>> {
  name: FieldPath<FormData>
  labelKey: string
  placeholderKey?: string
  descriptionKey?: string
  kind?: "text" | "number" | "textarea" | "select"
  options?: ClinicalFormOption[]
  required?: boolean
  inputMode?: "text" | "numeric" | "decimal"
  autoCapitalize?: "none" | "sentences" | "words" | "characters"
  spellCheck?: boolean
  maxLength?: number
  range?: { minimum: number; maximum: number }
  span?: "full" | "half"
  sectionKey?: string
  autofillTarget?: FieldPath<FormData>
}

export interface ClinicalFormConfig<FormData extends Record<string, unknown>> {
  titleKey: string
  descriptionKey?: string
  confirmKey: string
  cancelKey?: string
  fields: ClinicalFormField<FormData>[]
  schema: z.ZodType<FormData>
  defaultValues?: Partial<FormData>
  transformOnSubmit?: (formData: FormData) => FormData
}

interface ClinicalFormModalProps<FormData extends Record<string, unknown>> {
  isOpen: boolean
  onClose: () => void
  onSubmit: (formData: FormData) => void | Promise<unknown>
  isPending: boolean
  hasSubmitError?: boolean
  config: ClinicalFormConfig<FormData>
}

interface CharacterCounterProps<FormData extends Record<string, unknown>> {
  control: Control<FormData>
  name: FieldPath<FormData>
  maximum: number
}

const CharacterCounter = <FormData extends Record<string, unknown>>({
  control,
  name,
  maximum,
}: CharacterCounterProps<FormData>) => {
  const { t } = useTranslation("patients")
  const currentValue = useWatch({ control, name }) as unknown
  const current = typeof currentValue === "string" ? currentValue.length : 0

  return (
    <span className="shrink-0 text-xs tabular-nums text-muted-foreground">
      {t("form.characterCount", { current, max: maximum })}
    </span>
  )
}

interface FieldSection<FormData extends Record<string, unknown>> {
  sectionKey?: string
  fields: ClinicalFormField<FormData>[]
}

const groupFieldsBySection = <FormData extends Record<string, unknown>>(
  fields: ClinicalFormField<FormData>[]
): FieldSection<FormData>[] => {
  const sections: FieldSection<FormData>[] = []
  for (const field of fields) {
    const lastSection = sections.at(-1)
    if (lastSection && lastSection.sectionKey === field.sectionKey) {
      lastSection.fields.push(field)
    } else {
      sections.push({ sectionKey: field.sectionKey, fields: [field] })
    }
  }
  return sections
}

const stripTrailingCode = (optionText: string) => optionText.replace(/\s*\([^)]*\)\s*$/, "").trim()

export function ClinicalFormModal<FormData extends Record<string, unknown>>({
  isOpen,
  onClose,
  onSubmit,
  isPending,
  hasSubmitError = false,
  config,
}: ClinicalFormModalProps<FormData>) {
  const { t } = useTranslation("patients")
  const formId = useId()
  const runWithSubmissionLock = useSubmissionLock()
  const [wasOpen, setWasOpen] = useState(isOpen)
  const [hasAttemptedSubmit, setHasAttemptedSubmit] = useState(false)

  const {
    register,
    handleSubmit,
    control,
    reset,
    getValues,
    setValue,
    formState: { errors, isDirty, submitCount },
  } = useForm<FormData>({
    resolver: zodResolver(config.schema as never) as Resolver<FormData, unknown, FormData>,
    defaultValues: config.defaultValues as DefaultValues<FormData>,
    mode: "onTouched",
  })

  const { isConfirmationOpen, requestClose, keepEditing, discardChanges } = useDiscardGuard({
    isDirty,
    onClose,
  })

  if (wasOpen !== isOpen) {
    setWasOpen(isOpen)
    if (!isOpen) {
      setHasAttemptedSubmit(false)
    }
  }

  const defaultValuesKey = JSON.stringify(config.defaultValues ?? {})

  useEffect(() => {
    const defaultValues = JSON.parse(defaultValuesKey) as DefaultValues<FormData>
    if (isOpen) {
      reset(defaultValues, { keepDirtyValues: true })
    } else {
      reset(defaultValues)
    }
  }, [isOpen, defaultValuesKey, reset])

  if (!isOpen) {
    return null
  }

  const submitForm = async (formData: FormData) => {
    if (isPending) {
      return
    }
    await runWithSubmissionLock(async () => {
      setHasAttemptedSubmit(true)
      const preparedFormData = config.transformOnSubmit ? config.transformOnSubmit(formData) : formData
      await onSubmit(preparedFormData)
    })
  }

  const buildFieldId = (field: ClinicalFormField<FormData>) => `${formId}-${field.name}`

  const resolveOptionText = (option: ClinicalFormOption | undefined) => {
    if (!option) {
      return ""
    }
    return option.labelKey ? t(option.labelKey) : (option.label ?? option.value)
  }

  const autofillTargetFromOption = (
    field: ClinicalFormField<FormData>,
    previousValue: string,
    selectedValue: string
  ) => {
    if (!field.autofillTarget) {
      return
    }
    const findOption = (value: string) => field.options?.find((option) => option.value === value)
    const previousText = stripTrailingCode(resolveOptionText(findOption(previousValue)))
    const currentTargetValue = String(getValues(field.autofillTarget) ?? "")
    if (currentTargetValue === "" || currentTargetValue === previousText) {
      setValue(
        field.autofillTarget,
        stripTrailingCode(resolveOptionText(findOption(selectedValue))) as PathValue<FormData, FieldPath<FormData>>,
        { shouldDirty: true, shouldValidate: submitCount > 0 }
      )
    }
  }

  const getErrorMessage = (field: ClinicalFormField<FormData>) => {
    const message = errors[field.name as keyof FormData]?.message
    return typeof message === "string" ? message : undefined
  }

  const renderControl = (field: ClinicalFormField<FormData>, controlProperties: FormFieldControlProperties) => {
    const placeholder = field.placeholderKey ? t(field.placeholderKey) : undefined
    if (field.kind === "select") {
      return (
        <Controller
          control={control}
          name={field.name}
          render={({ field: fieldApi }) => (
            <Select
              value={(fieldApi.value as string) ?? ""}
              onValueChange={(selectedValue) => {
                const previousValue = String(fieldApi.value ?? "")
                fieldApi.onChange(selectedValue)
                autofillTargetFromOption(field, previousValue, selectedValue)
              }}
              onOpenChange={(open) => {
                if (!open) {
                  fieldApi.onBlur()
                }
              }}
            >
              <SelectTrigger
                ref={fieldApi.ref}
                className="min-h-11 w-full md:min-h-9"
                {...controlProperties}
                aria-required={undefined}
              >
                <SelectValue placeholder={placeholder} />
              </SelectTrigger>
              <SelectContent>
                {field.options?.map((option) => (
                  <SelectItem key={option.value} value={option.value}>
                    {resolveOptionText(option)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
        />
      )
    }
    if (field.kind === "textarea") {
      return (
        <Textarea
          className="min-h-24 resize-y"
          placeholder={placeholder}
          maxLength={field.maxLength}
          autoCapitalize={field.autoCapitalize ?? "sentences"}
          spellCheck={field.spellCheck}
          {...controlProperties}
          {...register(field.name)}
        />
      )
    }
    const isNumber = field.kind === "number"
    return (
      <Input
        type={isNumber ? "number" : "text"}
        step={isNumber ? "any" : undefined}
        min={field.range?.minimum}
        max={field.range?.maximum}
        inputMode={field.inputMode ?? (isNumber ? "decimal" : undefined)}
        autoComplete="off"
        autoCapitalize={field.autoCapitalize ?? (isNumber ? undefined : "sentences")}
        spellCheck={field.spellCheck ?? (isNumber ? undefined : true)}
        maxLength={field.maxLength}
        placeholder={placeholder}
        className={cn("min-h-11 md:min-h-9", field.autoCapitalize === "characters" && "uppercase")}
        onWheel={isNumber ? (event) => event.currentTarget.blur() : undefined}
        {...controlProperties}
        {...register(field.name, isNumber ? { valueAsNumber: true } : undefined)}
      />
    )
  }

  const renderField = (field: ClinicalFormField<FormData>) => {
    const isRequired = field.required !== false
    const description = field.range
      ? t("form.range", { min: field.range.minimum, max: field.range.maximum })
      : field.descriptionKey
        ? t(field.descriptionKey)
        : undefined
    const showsCounter = field.kind === "textarea" && field.maxLength !== undefined
    return (
      <FormField
        key={field.name}
        id={buildFieldId(field)}
        label={t(field.labelKey)}
        required={isRequired}
        description={description}
        error={getErrorMessage(field)}
        className={field.span === "half" ? undefined : "md:col-span-2"}
        counter={
          showsCounter ? (
            <CharacterCounter control={control} name={field.name} maximum={field.maxLength as number} />
          ) : undefined
        }
      >
        {(controlProperties) => renderControl(field, controlProperties)}
      </FormField>
    )
  }

  const hasRequiredFields = config.fields.some((field) => field.required !== false)
  const errorSummaryItems: FormErrorSummaryItem[] =
    config.fields.length >= ERROR_SUMMARY_MIN_FIELDS && submitCount > 0
      ? config.fields.flatMap((field) => {
          const message = getErrorMessage(field)
          return message ? [{ fieldId: buildFieldId(field), label: t(field.labelKey), message }] : []
        })
      : []

  const sections = groupFieldsBySection(config.fields)

  return (
    <>
      <Dialog open={isOpen} onOpenChange={(open) => !open && requestClose()}>
        <FormDialogLayout
          title={t(config.titleKey)}
          description={config.descriptionKey ? t(config.descriptionKey) : undefined}
          onSubmit={handleSubmit(submitForm as SubmitHandler<FormData>)}
          cancelLabel={t(config.cancelKey ?? "modal.cancel")}
          onCancel={requestClose}
          submitLabel={t(config.confirmKey)}
          isPending={isPending}
        >
          {hasSubmitError && hasAttemptedSubmit ? <FormSubmitError /> : null}
          <FormErrorSummary items={errorSummaryItems} />
          {hasRequiredFields && config.fields.length > 1 ? (
            <p className="text-xs text-muted-foreground">{t("form.requiredLegend")}</p>
          ) : null}
          {sections.map((section) =>
            section.sectionKey ? (
              <FieldSet key={section.sectionKey} className="gap-3">
                <FieldLegend variant="label" className="mb-0 text-sm font-semibold">
                  {t(section.sectionKey)}
                </FieldLegend>
                <div className="grid gap-4 md:grid-cols-2">{section.fields.map(renderField)}</div>
              </FieldSet>
            ) : (
              <div key={section.fields[0].name} className="grid gap-4 md:grid-cols-2">
                {section.fields.map(renderField)}
              </div>
            )
          )}
        </FormDialogLayout>
      </Dialog>
      <DiscardChangesDialog open={isConfirmationOpen} onKeepEditing={keepEditing} onDiscard={discardChanges} />
    </>
  )
}
