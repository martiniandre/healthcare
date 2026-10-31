import { useId, useMemo, useState } from "react"
import { useForm, useWatch, type Resolver } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslation } from "react-i18next"
import { useAbility } from "../../../../shared/auth/useAbility"
import { Action, Feature } from "../../../../shared/auth/types"
import { Dialog } from "../../../../shared/components/ui/dialog"
import { DiscardChangesDialog } from "../../../../shared/components/ui/discard-changes-dialog"
import { Input } from "../../../../shared/components/ui/input"
import { Textarea } from "../../../../shared/components/ui/textarea"
import { useDiscardGuard } from "../../../../shared/hooks/useDiscardGuard"
import { hasBlockingAlert, mergeAlerts } from "../../medicationSafety"
import {
  getMedicationPrescriptionSchema,
  type MedicationPrescriptionFormData,
} from "../../patient_schemas"
import { useMedicationSafetyCheck } from "../../queries"
import type { MedicationSafetyAlert } from "../../types"
import { FormDialogLayout } from "../form/FormDialogLayout"
import { FormField } from "../form/FormField"
import { FormSubmitError } from "../form/FormSubmitError"
import { useSubmissionLock } from "../form/useSubmissionLock"
import { MedicationOverrideSection } from "../medication/MedicationOverrideSection"
import { SafetyAlerts } from "../medication/SafetyAlerts"

export interface MedicationSubmission {
  medicationDisplay: string
  dosageInstruction: string
  overrideJustification?: string
}

export type MedicationSubmitOutcome =
  | { status: "success" }
  | { status: "blocked"; alerts: MedicationSafetyAlert[] }
  | { status: "failed" }

interface MedicationModalProps {
  patientId: string
  isOpen: boolean
  onClose: () => void
  onSubmit: (submission: MedicationSubmission) => Promise<MedicationSubmitOutcome>
  isPending: boolean
}

interface AlertsForMedication {
  medicationName: string
  alerts: MedicationSafetyAlert[]
}

const EMPTY_FORM_VALUES: MedicationPrescriptionFormData = {
  medicationDisplay: "",
  dosageInstruction: "",
  overrideJustification: "",
  riskAcknowledged: false,
  unverifiedAcknowledged: false,
}

const MEDICATION_NAME_MAX_LENGTH = 255
const MEDICATION_DOSAGE_MAX_LENGTH = 1000

const NO_REQUIREMENTS = { requiresOverride: false, requiresUnverifiedAcknowledgement: false }

export const MedicationModal = ({ patientId, isOpen, onClose, onSubmit, isPending }: MedicationModalProps) => {
  const { t } = useTranslation("patients")
  const ability = useAbility()
  const canOverride = ability.can(Action.Create, Feature.MedicationRequest)
  const [serverAlerts, setServerAlerts] = useState<AlertsForMedication | null>(null)
  const [lateAlerts, setLateAlerts] = useState<AlertsForMedication | null>(null)
  const [failedVerificationName, setFailedVerificationName] = useState<string | null>(null)
  const [isVerifying, setIsVerifying] = useState(false)
  const [hasSubmitFailed, setHasSubmitFailed] = useState(false)
  const formId = useId()
  const runWithSubmissionLock = useSubmissionLock()

  const baseResolver = useMemo(
    () =>
      zodResolver(getMedicationPrescriptionSchema(t, NO_REQUIREMENTS) as never) as Resolver<
        MedicationPrescriptionFormData,
        unknown,
        MedicationPrescriptionFormData
      >,
    [t],
  )

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    setError,
    control,
    formState: { errors, isDirty },
  } = useForm<MedicationPrescriptionFormData>({
    resolver: baseResolver,
    defaultValues: EMPTY_FORM_VALUES,
    mode: "onTouched",
  })

  const medicationName = useWatch({ control, name: "medicationDisplay" })
  const riskAcknowledged = useWatch({ control, name: "riskAcknowledged" })
  const unverifiedAcknowledged = useWatch({ control, name: "unverifiedAcknowledged" })
  const dosageText = useWatch({ control, name: "dosageInstruction" })
  const trimmedName = medicationName.trim()

  const safety = useMedicationSafetyCheck(patientId, medicationName, isOpen)

  const alertsForCurrentName = (alertsForMedication: AlertsForMedication | null): MedicationSafetyAlert[] =>
    alertsForMedication && alertsForMedication.medicationName === trimmedName ? alertsForMedication.alerts : []

  const combinedAlerts = mergeAlerts(
    mergeAlerts(safety.alerts, alertsForCurrentName(lateAlerts)),
    alertsForCurrentName(serverAlerts),
  )
  const isBlocking = hasBlockingAlert(combinedAlerts)
  const hasServerBlock = hasBlockingAlert(alertsForCurrentName(serverAlerts))
  const effectiveStatus = failedVerificationName === trimmedName ? "error" : safety.status
  const requiresOverride = isBlocking && canOverride
  const isSubmitBlocked = isPending || isVerifying || effectiveStatus === "loading" || (isBlocking && !canOverride)

  const closeAndReset = () => {
    reset(EMPTY_FORM_VALUES)
    setServerAlerts(null)
    setLateAlerts(null)
    setFailedVerificationName(null)
    setHasSubmitFailed(false)
    onClose()
  }

  const { isConfirmationOpen, requestClose, keepEditing, discardChanges } = useDiscardGuard({
    isDirty,
    onClose: closeAndReset,
  })

  const validateRequirements = (
    formData: MedicationPrescriptionFormData,
    requirements: { requiresOverride: boolean; requiresUnverifiedAcknowledgement: boolean },
  ): boolean => {
    const validation = getMedicationPrescriptionSchema(t, requirements).safeParse(formData)
    if (validation.success) return true
    for (const issue of validation.error.issues) {
      const fieldName = issue.path[0]
      if (fieldName === "overrideJustification" || fieldName === "riskAcknowledged" || fieldName === "unverifiedAcknowledged") {
        setError(fieldName, { message: issue.message })
      }
    }
    return false
  }

  const submitPrescription = async (formData: MedicationPrescriptionFormData) => {
    if (isPending) return
    await runWithSubmissionLock(() => runSubmission(formData))
  }

  const runSubmission = async (formData: MedicationPrescriptionFormData) => {
    setHasSubmitFailed(false)
    setIsVerifying(true)
    const verification = await safety.verifyNow()
    setIsVerifying(false)

    let blockingNow = isBlocking
    let unverifiedNow: boolean
    if (verification.verified) {
      setFailedVerificationName(null)
      setLateAlerts({ medicationName: trimmedName, alerts: verification.alerts })
      blockingNow = hasBlockingAlert(mergeAlerts(verification.alerts, alertsForCurrentName(serverAlerts)))
      unverifiedNow = false
    } else {
      setFailedVerificationName(trimmedName)
      unverifiedNow = true
    }

    const requirements = {
      requiresOverride: blockingNow && canOverride,
      requiresUnverifiedAcknowledgement: unverifiedNow,
    }
    if (blockingNow && !canOverride) return
    if (!validateRequirements(formData, requirements)) return

    const outcome = await onSubmit({
      medicationDisplay: formData.medicationDisplay,
      dosageInstruction: formData.dosageInstruction,
      overrideJustification: blockingNow ? formData.overrideJustification.trim() : undefined,
    })

    if (outcome.status === "success") {
      reset(EMPTY_FORM_VALUES)
      setServerAlerts(null)
      setLateAlerts(null)
      setFailedVerificationName(null)
    } else if (outcome.status === "blocked") {
      setServerAlerts({ medicationName: trimmedName, alerts: outcome.alerts })
    } else {
      setHasSubmitFailed(true)
    }
  }

  if (!isOpen) {
    return null
  }

  const medicationNameError = errors.medicationDisplay?.message
  const dosageError = errors.dosageInstruction?.message

  return (
    <>
      <Dialog open={isOpen} onOpenChange={(open) => !open && requestClose()}>
        <FormDialogLayout
          title={t("modals.medication.title")}
          onSubmit={handleSubmit(submitPrescription)}
          cancelLabel={t("modal.cancel")}
          onCancel={requestClose}
          submitLabel={requiresOverride ? t("medicationSafety.override.confirm") : t("modals.medication.confirm")}
          isPending={isPending || isVerifying}
          isSubmitDisabled={isSubmitBlocked}
        >
          {hasSubmitFailed ? <FormSubmitError /> : null}
          <p className="text-xs text-muted-foreground">{t("form.requiredLegend")}</p>
          <FormField
            id={`${formId}-medication-name`}
            label={t("modals.medication.name")}
            description={t("modals.medication.nameHint")}
            error={medicationNameError}
          >
            {(controlProperties) => (
              <Input
                type="text"
                placeholder={t("modals.medication.namePlaceholder")}
                autoComplete="off"
                autoCapitalize="sentences"
                spellCheck={false}
                maxLength={MEDICATION_NAME_MAX_LENGTH}
                className="min-h-11 md:min-h-9"
                {...controlProperties}
                {...register("medicationDisplay")}
              />
            )}
          </FormField>
          <FormField
            id={`${formId}-medication-dosage`}
            label={t("modals.medication.dosage")}
            description={t("modals.medication.dosageHint")}
            error={dosageError}
            counter={
              <span className="shrink-0 text-xs tabular-nums text-muted-foreground">
                {t("form.characterCount", { current: dosageText.length, max: MEDICATION_DOSAGE_MAX_LENGTH })}
              </span>
            }
          >
            {(controlProperties) => (
              <Textarea
                className="min-h-24 resize-y"
                placeholder={t("modals.medication.dosagePlaceholder")}
                autoCapitalize="sentences"
                maxLength={MEDICATION_DOSAGE_MAX_LENGTH}
                {...controlProperties}
                {...register("dosageInstruction")}
              />
            )}
          </FormField>

          <SafetyAlerts
            status={effectiveStatus}
            alerts={combinedAlerts}
            onRetry={() => {
              setFailedVerificationName(null)
              safety.retry()
            }}
            unverifiedAcknowledged={unverifiedAcknowledged}
            onUnverifiedAcknowledgedChange={(acknowledged) =>
              setValue("unverifiedAcknowledged", acknowledged, { shouldValidate: false })
            }
            unverifiedAcknowledgementError={errors.unverifiedAcknowledged?.message}
          />

          {isBlocking && (
            <MedicationOverrideSection
              canOverride={canOverride}
              justificationProps={register("overrideJustification")}
              justificationError={errors.overrideJustification?.message}
              riskAcknowledged={riskAcknowledged}
              onRiskAcknowledgedChange={(acknowledged) =>
                setValue("riskAcknowledged", acknowledged, { shouldValidate: false })
              }
              riskAcknowledgedError={errors.riskAcknowledged?.message}
              blockedByServer={hasServerBlock}
            />
          )}
        </FormDialogLayout>
      </Dialog>
      <DiscardChangesDialog open={isConfirmationOpen} onKeepEditing={keepEditing} onDiscard={discardChanges} />
    </>
  )
}
