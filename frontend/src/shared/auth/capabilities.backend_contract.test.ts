import { readFileSync } from "node:fs"
import { resolve } from "node:path"
import { describe, it, expect } from "vitest"
import { allCapabilities } from "./capabilities"

const backendCapabilityFilePath = resolve(process.cwd(), "../backend/internal/app/policy/capability.go")

const readBackendCapabilities = (): Set<string> => {
  const backendSource = readFileSync(backendCapabilityFilePath, "utf-8")
  const capabilityMatches = backendSource.matchAll(/"([a-z_]+:[a-z_]+)"/g)
  return new Set(Array.from(capabilityMatches, (match) => match[1]))
}

describe("frontend capabilities against the backend matrix", () => {
  it("should only use capability names that the backend declares", () => {
    const backendCapabilities = readBackendCapabilities()

    const unknownToBackend = allCapabilities.filter((capability) => !backendCapabilities.has(capability))

    expect(unknownToBackend).toEqual([])
  })
})
