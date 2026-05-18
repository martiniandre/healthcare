import { create } from "zustand"

interface LayoutState {
  isMobileSidebarOpen: boolean
  toggleMobileSidebar: () => void
  closeMobileSidebar: () => void
}

export const useLayoutStore = create<LayoutState>((set) => ({
  isMobileSidebarOpen: false,
  toggleMobileSidebar: () => set((state) => ({ isMobileSidebarOpen: !state.isMobileSidebarOpen })),
  closeMobileSidebar: () => set({ isMobileSidebarOpen: false }),
}))
