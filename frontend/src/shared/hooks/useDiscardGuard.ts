import { useCallback, useState } from "react"

interface DiscardGuardOptions {
  isDirty: boolean
  onClose: () => void
}

export const useDiscardGuard = ({ isDirty, onClose }: DiscardGuardOptions) => {
  const [isConfirmationOpen, setIsConfirmationOpen] = useState(false)

  const requestClose = useCallback(() => {
    if (isDirty) {
      setIsConfirmationOpen(true)
      return
    }
    onClose()
  }, [isDirty, onClose])

  const keepEditing = useCallback(() => setIsConfirmationOpen(false), [])

  const discardChanges = useCallback(() => {
    setIsConfirmationOpen(false)
    onClose()
  }, [onClose])

  return { isConfirmationOpen, requestClose, keepEditing, discardChanges }
}
