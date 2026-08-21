import { Gauge, HeartPulse, Ruler, Scale, Thermometer, Waves, Wind, type LucideIcon } from "lucide-react"

export interface VitalSignDisplayMetadata {
  loincCode: string
  labelKey: string | null
  IconComponent: LucideIcon
  iconClassName: string
}

export const vitalSignDisplayMetadata: VitalSignDisplayMetadata[] = [
  { loincCode: "8867-4", labelKey: "details.vitalsCard.metrics.heartRate", IconComponent: HeartPulse, iconClassName: "bg-red-50 border-red-100 text-red-600" },
  { loincCode: "8310-5", labelKey: "details.vitalsCard.metrics.bodyTemperature", IconComponent: Thermometer, iconClassName: "bg-amber-50 border-amber-100 text-amber-600" },
  { loincCode: "8480-6", labelKey: "details.vitalsCard.metrics.systolicBloodPressure", IconComponent: Gauge, iconClassName: "bg-blue-50 border-blue-100 text-blue-600" },
  { loincCode: "8462-4", labelKey: "details.vitalsCard.metrics.diastolicBloodPressure", IconComponent: Gauge, iconClassName: "bg-indigo-50 border-indigo-100 text-indigo-600" },
  { loincCode: "59408-5", labelKey: "details.vitalsCard.metrics.oxygenSaturation", IconComponent: Wind, iconClassName: "bg-sky-50 border-sky-100 text-sky-600" },
  { loincCode: "9279-1", labelKey: "details.vitalsCard.metrics.respiratoryRate", IconComponent: Waves, iconClassName: "bg-teal-50 border-teal-100 text-teal-600" },
  { loincCode: "29463-7", labelKey: "details.vitalsCard.metrics.weightKg", IconComponent: Scale, iconClassName: "bg-violet-50 border-violet-100 text-violet-600" },
  { loincCode: "8302-2", labelKey: "details.vitalsCard.metrics.heightCm", IconComponent: Ruler, iconClassName: "bg-emerald-50 border-emerald-100 text-emerald-600" },
]

export const findVitalSignDisplay = (loincCode: string): VitalSignDisplayMetadata | undefined =>
  vitalSignDisplayMetadata.find((displayMetadata) => displayMetadata.loincCode === loincCode)
