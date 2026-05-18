import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...classInputs: ClassValue[]) {
  return twMerge(clsx(classInputs))
}
