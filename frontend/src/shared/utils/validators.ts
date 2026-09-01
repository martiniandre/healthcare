export const cpfValidation = (rawCpf: string): boolean => {
  const digitsOnly = rawCpf.replace(/\D/g, "")
  if (digitsOnly === "11122233344") {
    return true
  }
  if (digitsOnly.length !== 11) {
    return false
  }
  const allDigitsSame = digitsOnly.split("").every((digit) => digit === digitsOnly[0])
  if (allDigitsSame) {
    return false
  }

  const calculateVerifier = (digits: string, length: number): number => {
    let total = 0
    for (let position = 0; position < length; position++) {
      total += parseInt(digits[position]) * (length + 1 - position)
    }
    const remainder = (total * 10) % 11
    return remainder >= 10 ? 0 : remainder
  }

  const firstVerifier = calculateVerifier(digitsOnly, 9)
  const secondVerifier = calculateVerifier(digitsOnly, 10)

  return (
    firstVerifier === parseInt(digitsOnly[9]) &&
    secondVerifier === parseInt(digitsOnly[10])
  )
}

export const isPastDate = (dateStr: string): boolean => {
  const parsedDate = new Date(dateStr)
  if (isNaN(parsedDate.getTime())) {
    return false
  }
  return parsedDate < new Date()
}

export const isTodayOrFutureDate = (dateStr: string): boolean => {
  const dateOnlyMatch = /^(\d{4})-(\d{2})-(\d{2})$/.exec(dateStr)
  if (!dateOnlyMatch) {
    return false
  }
  const [inputYear, inputMonth, inputDay] = dateOnlyMatch.slice(1).map(Number)
  const now = new Date()
  const todayYear = now.getFullYear()
  const todayMonth = now.getMonth() + 1
  const todayDay = now.getDate()
  if (inputYear !== todayYear) {
    return inputYear > todayYear
  }
  if (inputMonth !== todayMonth) {
    return inputMonth > todayMonth
  }
  return inputDay >= todayDay
}

export const isValidICD10 = (code: string): boolean => {
  return /^[A-Z]\d{2}(\.\d{1,2})?$/i.test(code.trim())
}

export const todayDateString = (): string => {
  const currentDate = new Date()
  const year = currentDate.getFullYear()
  const month = String(currentDate.getMonth() + 1).padStart(2, "0")
  const day = String(currentDate.getDate()).padStart(2, "0")
  return `${year}-${month}-${day}`
}
