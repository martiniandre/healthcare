export const calculateTotalPages = (totalPatientsCount: number, itemsPerPage: number) =>
  Math.max(1, Math.ceil(totalPatientsCount / itemsPerPage))