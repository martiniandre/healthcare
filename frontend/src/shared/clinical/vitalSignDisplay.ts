import { Gauge, HeartPulse, Ruler, Scale, Thermometer, Waves, Wind, type LucideIcon } from "lucide-react"

export interface VitalSignDisplayMetadata {
  loincCode: string
  labelKey: string | null
  IconComponent: LucideIcon
  iconClassName: string
}

export const vitalSignDisplayMetadata: VitalSignDisplayMetadata[] = [
  { loincCode: "8867-4", labelKey: "details.vitalsCard.metrics.heartRate", IconComponent: HeartPulse, iconClassName: "bg-destructive/10 border-destructive/20 text-destructive" },
  { loincCode: "8310-5", labelKey: "details.vitalsCard.metrics.bodyTemperature", IconComponent: Thermometer, iconClassName: "bg-warning/10 border-warning/20 text-warning" },
  { loincCode: "8480-6", labelKey: "details.vitalsCard.metrics.systolicBloodPressure", IconComponent: Gauge, iconClassName: "bg-primary/10 border-primary/20 text-primary" },
  { loincCode: "8462-4", labelKey: "details.vitalsCard.metrics.diastolicBloodPressure", IconComponent: Gauge, iconClassName: "bg-muted border-border text-muted-foreground" },
  { loincCode: "59408-5", labelKey: "details.vitalsCard.metrics.oxygenSaturation", IconComponent: Wind, iconClassName: "bg-primary/10 border-primary/20 text-primary" },
  { loincCode: "9279-1", labelKey: "details.vitalsCard.metrics.respiratoryRate", IconComponent: Waves, iconClassName: "bg-primary/10 border-primary/20 text-primary" },
  { loincCode: "29463-7", labelKey: "details.vitalsCard.metrics.weightKg", IconComponent: Scale, iconClassName: "bg-muted border-border text-muted-foreground" },
  { loincCode: "8302-2", labelKey: "details.vitalsCard.metrics.heightCm", IconComponent: Ruler, iconClassName: "bg-success/10 border-success/20 text-success" },
]

export const findVitalSignDisplay = (loincCode: string): VitalSignDisplayMetadata | undefined =>
  vitalSignDisplayMetadata.find((displayMetadata) => displayMetadata.loincCode === loincCode)
