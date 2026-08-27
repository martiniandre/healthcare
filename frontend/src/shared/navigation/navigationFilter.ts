import { navigationGroups } from "./navigationConfig"
import type { SidebarNavigationGroup, SidebarNavigationItem } from "./navigationConfig"

export const isNavigationItemVisible = (
  navigationItem: SidebarNavigationItem,
  currentRole: string | null
): boolean => {
  const matchesRoleScope = navigationItem.patientOnly
    ? currentRole === "PATIENT"
    : navigationItem.staffOnly
      ? currentRole !== "PATIENT"
      : true
  const matchesAdminScope = navigationItem.adminOnly ? currentRole === "ADMIN" : true
  return matchesRoleScope && matchesAdminScope
}

export const isNavigationItemActive = (
  navigationItem: SidebarNavigationItem,
  currentPathname: string
): boolean => {
  const exactPathMatch = currentPathname === navigationItem.path
  const nestedPathMatch = navigationItem.path !== "/" && currentPathname.startsWith(navigationItem.path)
  const homePathMatch =
    navigationItem.path === "/" &&
    (currentPathname === "/" || currentPathname.startsWith("/patients"))
  return exactPathMatch || nestedPathMatch || homePathMatch
}

export const getVisibleNavigationGroups = (
  currentRole: string | null,
  sourceGroups: SidebarNavigationGroup[] = navigationGroups
): SidebarNavigationGroup[] =>
  sourceGroups
    .map((navigationGroup) => ({
      ...navigationGroup,
      items: navigationGroup.items.filter((navigationItem) =>
        isNavigationItemVisible(navigationItem, currentRole)
      ),
    }))
    .filter((navigationGroup) => navigationGroup.items.length > 0)