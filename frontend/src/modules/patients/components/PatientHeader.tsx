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
}

export const PatientHeader = ({ patient, onBack }: PatientHeaderProps) => {
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
      className="sticky top-16 z-30 flex flex-col gap-2 rounded-lg border border-border bg-card px-3 py-2 text-left shadow-sm sm:px-4"
    >
      <div className="flex items-center gap-3">
        <Button variant="outline" size="sm" onClick={onBack} aria-label={t("details.back")} className="shrink-0">
          <ArrowLeft />
          <span className="hidden sm:inline">{t("details.back")}</span>
        </Button>

        <div className="min-w-0 flex-1">
          <h2 className="truncate text-base font-display font-bold leading-tight text-foreground sm:text-lg">
            {patient.full_name}
          </h2>
          <p className="hidden flex-wrap items-center gap-x-2 text-xs text-muted-foreground md:flex">
            {ageLabel && <span className="font-semibold tabular-nums text-foreground">{ageLabel}</span>}
            <span className="tabular-nums">
              {t("details.cpf")}: {patient.document_id}
            </span>
            <span className="font-mono text-[10px]">
              {t("details.fhirId")}: {patient.fhir_resource_id}
            </span>
          </p>
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
      </div>

      <div className="hidden flex-wrap items-center gap-x-6 gap-y-1.5 md:flex">
        {renderAllergyGroup(false)}
        {renderConditionGroup(false)}
      </div>
    </header>
  )
}
