import { useState } from "react"
import { useTranslation } from "react-i18next"
import { Checkbox } from "@/shared/components/ui/checkbox"
import { Field, FieldDescription, FieldError, FieldLabel } from "@/shared/components/ui/field"
import { Input } from "@/shared/components/ui/input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/shared/components/ui/select"
import { Switch } from "@/shared/components/ui/switch"
import { Textarea } from "@/shared/components/ui/textarea"
import { Showcase } from "./SectionShell"

export const FormsShowcase = () => {
  const { t } = useTranslation("designSystem")
  const [isNotificationEnabled, setIsNotificationEnabled] = useState(true)

  return (
    <>
      <Showcase title={t("components.forms.input")}>
        <div className="grid gap-4 md:grid-cols-2">
          <Field>
            <FieldLabel htmlFor="ds-input-default">{t("components.forms.name")}</FieldLabel>
            <Input id="ds-input-default" placeholder={t("components.forms.namePlaceholder")} />
            <FieldDescription>{t("components.forms.hint")}</FieldDescription>
          </Field>
          <Field data-invalid="true">
            <FieldLabel htmlFor="ds-input-invalid">{t("components.forms.name")}</FieldLabel>
            <Input
              id="ds-input-invalid"
              aria-invalid="true"
              aria-describedby="ds-input-invalid-error"
              defaultValue={t("components.forms.invalidValue")}
            />
            <FieldError id="ds-input-invalid-error">{t("components.forms.error")}</FieldError>
          </Field>
        </div>
      </Showcase>
      <Showcase title={t("components.forms.textarea")}>
        <Field>
          <FieldLabel htmlFor="ds-textarea">{t("components.forms.notes")}</FieldLabel>
          <Textarea id="ds-textarea" placeholder={t("components.forms.notesPlaceholder")} />
        </Field>
      </Showcase>
      <Showcase title={t("components.forms.select")}>
        <Field>
          <FieldLabel htmlFor="ds-select">{t("components.forms.specialty")}</FieldLabel>
          <Select defaultValue="cardiology">
            <SelectTrigger id="ds-select" className="w-full">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="cardiology">{t("components.forms.cardiology")}</SelectItem>
              <SelectItem value="neurology">{t("components.forms.neurology")}</SelectItem>
            </SelectContent>
          </Select>
        </Field>
      </Showcase>
      <Showcase title={t("components.forms.toggles")}>
        <div className="flex flex-col gap-3">
          <Field orientation="horizontal">
            <Checkbox id="ds-checkbox" defaultChecked />
            <FieldLabel htmlFor="ds-checkbox">{t("components.forms.consent")}</FieldLabel>
          </Field>
          <Field orientation="horizontal">
            <Switch
              id="ds-switch"
              checked={isNotificationEnabled}
              onCheckedChange={setIsNotificationEnabled}
            />
            <FieldLabel htmlFor="ds-switch">{t("components.forms.notifications")}</FieldLabel>
          </Field>
        </div>
      </Showcase>
    </>
  )
}
