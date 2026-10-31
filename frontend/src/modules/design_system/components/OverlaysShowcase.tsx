import { useTranslation } from "react-i18next"
import { Button } from "@/shared/components/ui/button"
import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/shared/components/ui/dialog"
import { Popover, PopoverContent, PopoverDescription, PopoverHeader, PopoverTitle, PopoverTrigger } from "@/shared/components/ui/popover"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/shared/components/ui/tooltip"
import { Showcase } from "./SectionShell"

export const OverlaysShowcase = () => {
  const { t } = useTranslation("designSystem")

  return (
    <>
      <Showcase title={t("components.tooltip.title")}>
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger asChild>
              <Button variant="outline">{t("components.tooltip.trigger")}</Button>
            </TooltipTrigger>
            <TooltipContent>{t("components.tooltip.content")}</TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </Showcase>
      <Showcase title={t("components.popover.title")}>
        <Popover>
          <PopoverTrigger asChild>
            <Button variant="outline">{t("components.popover.trigger")}</Button>
          </PopoverTrigger>
          <PopoverContent>
            <PopoverHeader>
              <PopoverTitle>{t("components.popover.heading")}</PopoverTitle>
              <PopoverDescription>{t("components.popover.description")}</PopoverDescription>
            </PopoverHeader>
          </PopoverContent>
        </Popover>
      </Showcase>
      <Showcase title={t("components.dialog.title")}>
        <Dialog>
          <DialogTrigger asChild>
            <Button variant="outline">{t("components.dialog.trigger")}</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{t("components.dialog.heading")}</DialogTitle>
              <DialogDescription>{t("components.dialog.description")}</DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <DialogClose asChild>
                <Button variant="outline">{t("components.dialog.cancel")}</Button>
              </DialogClose>
              <Button>{t("components.dialog.confirm")}</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </Showcase>
    </>
  )
}
