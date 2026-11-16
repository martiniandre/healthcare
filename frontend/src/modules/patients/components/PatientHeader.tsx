import type { ReactNode } from "react"
import { ArrowLeft } from "lucide-react"
import { useTranslation } from "react-i18next"
import { Button } from "../../../shared/components/ui/button"
import { Popover, PopoverContent, PopoverTrigger } from "../../../shared/components/ui/popover"
import { Skeleton } from "../../../shared/components/ui/skeleton"
import { StatusBadge } from "../../../shared/components/ui/status-badge"
import { Action, Feature } from "../../../shared/auth/AbilityContext"
import { useAbility } from "../../../shared/auth/useAbility"
import { usePatientAllergiesQuery, usePatientConditionsQuery } from "../queries"
import { ACTIVE_CONDITION_STATUSES } from "./clinicalStatus"
import { calculateAge } from "./patientAge"
import { PatientAlertGroup } from "./PatientAlertGroup"

interface PatientRepresentation {
  patient_id: string
  fhir_resource_id: string
  full_name: string
  birth_date: string
  document_id: string
  phone_number: string
}

interface PatientHeaderProps {
  patient: PatientRepresentation
  onBack: () => void
  actions?: ReactNode
}

const buildPatientInitials = (fullName: string): string => {
  const nameParts = fullName.trim().split(/\s+/).filter(Boolean)
  if (nameParts.length === 0) return "?"
  const lastInitial = nameParts.length > 1 ? nameParts[nameParts.length - 1].charAt(0) : ""
  return `${nameParts[0].charAt(0)}${lastInitial}`.toUpperCase()
}

export const PatientHeader = ({ patient, onBack, actions }: PatientHeaderProps) => {
  const { t } = useTranslation("patients")
  const ability = useAbility()
  const canReadAllergies = ability.can(Action.Read, Feature.Allergy)
  const canReadConditions = ability.can(Action.Read, Feature.Condition)

  const allergiesQuery = usePatientAllergiesQuery(canReadAllergies ? patient.fhir_resource_id : "")
  const conditionsQuery = usePatientConditionsQuery(canReadConditions ? patient.fhir_resource_id : "")

  const allergyNames = (allergiesQuery.data ?? [])
    .filter((allergy) => allergy.clinical_status === "active")
    .map((allergy) => allergy.allergen_display)
  const conditionNames = (conditionsQuery.data ?? [])
    .filter((condition) => ACTIVE_CONDITION_STATUSES.includes(condition.clinical_status))
    .map((condition) => condition.code_display)

  const age = calculateAge(patient.birth_date)
  const isAlertsLoading =
    (canReadAllergies && allergiesQuery.isLoading) || (canReadConditions && conditionsQuery.isLoading)
  const hasAlertsError =
    (canReadAllergies && allergiesQuery.isError) || (canReadConditions && conditionsQuery.isError)
  const alertsTotal = allergyNames.length + conditionNames.length

  const renderAllergyGroup = (showAllItems: boolean) =>
    canReadAllergies ? (
      <PatientAlertGroup
        testId="patient-header-allergies"
        label={
          allergyNames.length > 0
            ? t("header.allergiesCount", { count: allergyNames.length })
            : t("header.allergies")
        }
        items={allergyNames}
        itemStatus="critical"
        emptyLabel={t("header.noKnownAllergies")}
        isLoading={allergiesQuery.isLoading}
        isError={allergiesQuery.isError}
        onRetry={() => void allergiesQuery.refetch()}
        showAllItems={showAllItems}
      />
    ) : null

  const renderConditionGroup = (showAllItems: boolean) =>
    canReadConditions ? (
      <PatientAlertGroup
        testId="patient-header-conditions"
        label={
          conditionNames.length > 0
            ? t("header.conditionsCount", { count: conditionNames.length })
            : t("header.conditions")
        }
        items={conditionNames}
        itemStatus="critical"
        emptyLabel={t("header.noActiveConditions")}
        isLoading={conditionsQuery.isLoading}
        isError={conditionsQuery.isError}
        onRetry={() => void conditionsQuery.refetch()}
        showAllItems={showAllItems}
      />
    ) : null

  const alertsChip = isAlertsLoading ? (
    <Skeleton data-testid="patient-header-chip-loading" className="h-6 w-20" />
  ) : hasAlertsError ? (
    <StatusBadge status="attention">{t("header.alertsError")}</StatusBadge>
  ) : alertsTotal > 0 ? (
    <StatusBadge status="critical">{t("header.alertsCount", { count: alertsTotal })}</StatusBadge>
  ) : (
    <StatusBadge status="neutral">{t("header.noAlerts")}</StatusBadge>
  )

  const ageLabel = age === null ? null : t("header.age", { count: age })

  return (
    <header
      data-testid="patient-sticky-header"
      className="sticky top-14 z-30 overflow-hidden rounded-lg border border-border bg-card text-left shadow-xs"
    >
      <div className="flex flex-wrap items-center gap-3 p-4 sm:gap-4 sm:p-5">
        <Button variant="outline" size="icon" onClick={onBack} aria-label={t("details.back")} className="shrink-0">
          <ArrowLeft />
        </Button>

        <span
          aria-hidden="true"
          className="hidden size-12 shrink-0 items-center justify-center rounded-full bg-sage-100 font-display text-lg font-semibold text-sage-800 sm:flex dark:bg-primary/15 dark:text-primary"
        >
          {buildPatientInitials(patient.full_name)}
        </span>

        <div className="min-w-0 flex-1">
          <h2 className="break-words font-display text-xl font-semibold leading-tight text-foreground sm:text-2xl">
            {patient.full_name}
          </h2>
          <p className="mt-1 hidden flex-wrap items-center gap-x-3 gap-y-1 text-sm text-muted-foreground md:flex">
            {ageLabel && <span className="font-semibold tabular-nums text-foreground">{ageLabel}</span>}
            <span aria-hidden="true" className="size-1 rounded-full bg-border" />
            <span className="tabular-nums">
              {t("details.cpf")}: {patient.document_id}
            </span>
            <span aria-hidden="true" className="size-1 rounded-full bg-border" />
            <span className="font-mono text-xs">
              {t("details.fhirId")}: {patient.fhir_resource_id}
            </span>
          </p>
          {ageLabel && <p className="mt-0.5 text-sm font-semibold tabular-nums text-foreground md:hidden">{ageLabel}</p>}
        </div>

        <div className="md:hidden">
          <Popover>
            <PopoverTrigger asChild>
              <button
                type="button"
                aria-label={t("header.alertsDetails")}
                className="rounded-md focus-visible:outline-2 focus-visible:outline-ring"
              >
                {alertsChip}
              </button>
            </PopoverTrigger>
            <PopoverContent align="end" className="flex w-72 flex-col gap-3">
              {ageLabel && <span className="text-sm font-semibold tabular-nums">{ageLabel}</span>}
              <span className="text-xs tabular-nums text-muted-foreground">
                {t("details.cpf")}: {patient.document_id}
              </span>
              <span className="font-mono text-[10px] text-muted-foreground">
                {t("details.fhirId")}: {patient.fhir_resource_id}
              </span>
              {renderAllergyGroup(true)}
              {renderConditionGroup(true)}
            </PopoverContent>
          </Popover>
        </div>

        {actions && <div className="flex w-full flex-wrap items-center gap-2 sm:w-auto">{actions}</div>}
      </div>

      <div className="hidden flex-wrap items-center gap-x-8 gap-y-2 border-t border-border bg-muted/40 px-5 py-2.5 md:flex">
        {renderAllergyGroup(false)}
        {renderConditionGroup(false)}
      </div>
    </header>
  )
}
