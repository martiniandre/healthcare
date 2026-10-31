const dateOnlyPattern = /^(\d{4})-(\d{2})-(\d{2})/

export const calculateAge = (birthDate: string, referenceDate: Date = new Date()): number | null => {
  const dateOnlyMatch = dateOnlyPattern.exec(birthDate)
  const parsedBirth = dateOnlyMatch
    ? new Date(Number(dateOnlyMatch[1]), Number(dateOnlyMatch[2]) - 1, Number(dateOnlyMatch[3]))
    : new Date(birthDate)
  if (Number.isNaN(parsedBirth.getTime())) {
    return null
  }
  let age = referenceDate.getFullYear() - parsedBirth.getFullYear()
  const hadBirthdayThisYear =
    referenceDate.getMonth() > parsedBirth.getMonth() ||
    (referenceDate.getMonth() === parsedBirth.getMonth() && referenceDate.getDate() >= parsedBirth.getDate())
  if (!hadBirthdayThisYear) {
    age -= 1
  }
  return age >= 0 ? age : null
}
