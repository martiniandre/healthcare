export const touchTargetClassName = "h-11 md:h-9"

export const selectTouchTargetClassName = "w-full data-[size=default]:h-11 md:data-[size=default]:h-9"

export const joinDescribedBy = (...identifiers: Array<string | false | undefined>): string | undefined => {
  const joinedIdentifiers = identifiers.filter(Boolean).join(" ")
  return joinedIdentifiers || undefined
}
