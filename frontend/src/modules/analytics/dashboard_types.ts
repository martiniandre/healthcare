export interface DoctorConsultation {
  doctor_name: string
  specialty: string
  count: number
}

export interface DepartmentWaitTime {
  department: string
  minutes: number
}

export interface DiagnosisCount {
  icd10_code: string
  description: string
  count: number
}

export interface DashboardData {
  consultations_today: number
  consultations_trend: string
  occupancy_rate: number
  occupancy_total_beds: number
  occupancy_occupied_beds: number
  avg_wait_time_minutes: number
  active_patients: number
  exams_today: number
  new_diagnoses_today: number
  consultations_per_doctor: DoctorConsultation[]
  wait_time_by_department: DepartmentWaitTime[]
  top_diagnoses: DiagnosisCount[]
}
