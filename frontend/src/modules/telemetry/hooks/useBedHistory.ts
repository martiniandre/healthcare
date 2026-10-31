import { useState } from "react"
import { appendBedSamples, type BedHistoryMap } from "../utils/bed_helpers"
import type { BedPatient } from "../types"

interface BedHistoryState {
  history: BedHistoryMap
  updatedAt: number
}

export const useBedHistory = (beds: BedPatient[], updatedAt: number): BedHistoryMap => {
  const [state, setState] = useState<BedHistoryState>({ history: {}, updatedAt: 0 })

  if (updatedAt !== 0 && updatedAt !== state.updatedAt && beds.length > 0) {
    setState({ history: appendBedSamples(state.history, beds), updatedAt })
  }

  return state.history
}
