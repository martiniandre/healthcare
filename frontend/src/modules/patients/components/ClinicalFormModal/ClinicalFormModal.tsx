import { Controller, useForm, type DefaultValues, type FieldPath, type Resolver, type SubmitHandler } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslation } from "react-i18next"
import type { z } from "zod"
import { Button } from "../../../../shared/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "../../../../shared/components/ui/dialog"
import { Input } from "../../../../shared/components/ui/input"
import { Textarea } from "../../../../shared/components/ui/textarea"
import { Field, FieldError, FieldLabel } from "../../../../shared/components/ui/field"
import { Spinner } from "../../../../shared/components/ui/spinner"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../../../../shared/components/ui/select"

export interface ClinicalFormOption {
  value: string
  labelKey?: string
  label?: string
}

export interface ClinicalFormField<FormData extends Record<string, unknown>> {
  name: FieldPath<FormData>
  labelKey: string
  placeholderKey?: string
  kind?: "text" | "number" | "textarea" | "select"
  options?: ClinicalFormOption[]
}

export interface ClinicalFormConfig<FormData extends Record<string, unknown>> {
  titleKey: string
  confirmKey: string
  cancelKey?: string
  fields: ClinicalFormField<FormData>[]
  schema: z.ZodType<FormData>
  defaultValues?: Partial<FormData>
  resetsAfterSubmit?: boolean
  transformOnSubmit?: (formData: FormData) => FormData
}

interface ClinicalFormModalProps<FormData extends Record<string, unknown>> {
  isOpen: boolean
  onClose: () => void
  onSubmit: (formData: FormData) => void
  isPending: boolean
  config: ClinicalFormConfig<FormData>
}

export function ClinicalFormModal<FormData extends Record<string, unknown>>({
  isOpen,
  onClose,
  onSubmit,
  isPending,
  config,
}: ClinicalFormModalProps<FormData>) {
  const { t } = useTranslation("patients")

  const {
    register,
    handleSubmit,
    control,
    reset,
    formState: { errors },
  } = useForm<FormData>({
    resolver: zodResolver(config.schema as never) as Resolver<FormData, unknown, FormData>,
    defaultValues: config.defaultValues as DefaultValues<FormData>,
  })

  if (!isOpen) {
    return null
  }

  const handleFormSubmit = (formData: FormData) => {
    const preparedFormData = config.transformOnSubmit ? config.transformOnSubmit(formData) : formData
    onSubmit(preparedFormData)
    if (config.resetsAfterSubmit) {
      reset()
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle className="text-left">{t(config.titleKey)}</DialogTitle>
        </DialogHeader>
        <form
          onSubmit={handleSubmit(handleFormSubmit as SubmitHandler<FormData>)}
          className="flex flex-col gap-4 text-left mt-4"
          noValidate
        >
          {config.fields.map((field) => {
            const errorMessage = errors[field.name as keyof FormData]?.message
            let fieldControl
            if (field.kind === "select") {
              fieldControl = (
                <Controller
                  control={control}
                  name={field.name}
                  render={({ field: fieldApi }) => (
                    <Select onValueChange={fieldApi.onChange} value={fieldApi.value as string}>
                      <SelectTrigger className="w-full" aria-invalid={Boolean(errorMessage)}>
                        <SelectValue
                          placeholder={field.placeholderKey ? t(field.placeholderKey) : undefined}
                        />
                      </SelectTrigger>
                      <SelectContent>
                        {field.options?.map((option) => (
                          <SelectItem key={option.value} value={option.value}>
                            {option.labelKey ? t(option.labelKey) : option.label ?? option.value}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  )}
                />
              )
            } else if (field.kind === "textarea") {
              fieldControl = (
                <Textarea
                  className="h-24 resize-none"
                  aria-invalid={Boolean(errorMessage)}
                  placeholder={field.placeholderKey ? t(field.placeholderKey) : undefined}
                  {...register(field.name)}
                />
              )
            } else {
              fieldControl = (
                <Input
                  type={field.kind === "number" ? "number" : "text"}
                  step={field.kind === "number" ? "any" : undefined}
                  placeholder={field.placeholderKey ? t(field.placeholderKey) : undefined}
                  aria-invalid={Boolean(errorMessage)}
                  {...register(
                    field.name,
                    field.kind === "number" ? { valueAsNumber: true } : undefined
                  )}
                />
              )
            }
            return (
              <Field className="gap-1" data-invalid={Boolean(errorMessage)} key={field.name}>
                <FieldLabel className="text-xs font-semibold text-muted-foreground">{t(field.labelKey)}</FieldLabel>
                {fieldControl}
                {errorMessage && <FieldError>{errorMessage as string}</FieldError>}
              </Field>
            )
          })}
          <div className="flex gap-3 justify-end mt-4">
            <Button variant="outline" type="button" onClick={onClose}>
              {t(config.cancelKey ?? "modal.cancel")}
            </Button>
            <Button type="submit" disabled={isPending}>
              {isPending && <Spinner />}
              {t(config.confirmKey)}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
