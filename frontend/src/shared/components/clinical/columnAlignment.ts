const ACTIONS_COLUMN_ID = "actions"

export const isActionsColumn = (columnId: string): boolean => columnId === ACTIONS_COLUMN_ID

export const actionsHeaderClassName = "text-right"

export const actionsHeaderContentClassName = "justify-end"

export const actionsCellClassName = "text-right [&>div]:justify-end"
