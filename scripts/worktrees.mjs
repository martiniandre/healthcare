#!/usr/bin/env node
import {
  existsSync,
  readFileSync,
  writeFileSync,
  renameSync,
  rmSync,
  mkdirSync,
  openSync,
  closeSync,
  readdirSync,
  statSync,
} from "node:fs"
import path from "node:path"
import { fileURLToPath } from "node:url"
import { execFileSync } from "node:child_process"

const REGISTRY_PARENT_DIR = process.env.WORKTREES_REGISTRY_DIR
  ? path.resolve(process.env.WORKTREES_REGISTRY_DIR)
  : path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..", "..", "healthcare-worktrees")
const REGISTRY_PATH = path.join(REGISTRY_PARENT_DIR, "registry.json")
const REGISTRY_LOCK_PATH = path.join(REGISTRY_PARENT_DIR, "registry.lock")
const LOCK_MAX_ATTEMPTS = 30
const LOCK_RETRY_DELAY_MS = 200
const STALE_LOCK_THRESHOLD_MS = 15000
const VALID_WORKTREE_TYPES = ["feat", "fix", "chore", "refactor", "docs"]

function sleep(milliseconds) {
  Atomics.wait(new Int32Array(new SharedArrayBuffer(4)), 0, 0, milliseconds)
}

function acquireRegistryLock() {
  for (let attempt = 0; attempt < LOCK_MAX_ATTEMPTS; attempt += 1) {
    try {
      const lockHandle = openSync(REGISTRY_LOCK_PATH, "wx")
      writeFileSync(REGISTRY_LOCK_PATH, String(process.pid))
      closeSync(lockHandle)
      return
    } catch (lockError) {
      if (lockError.code !== "EEXIST") throw lockError
      try {
        const lockAgeMilliseconds = Date.now() - statSync(REGISTRY_LOCK_PATH).mtimeMs
        if (lockAgeMilliseconds > STALE_LOCK_THRESHOLD_MS) {
          rmSync(REGISTRY_LOCK_PATH, { force: true })
          continue
        }
      } catch {}
      sleep(LOCK_RETRY_DELAY_MS)
    }
  }
  throw new Error(`Could not acquire registry lock after ${LOCK_MAX_ATTEMPTS} attempts: ${REGISTRY_LOCK_PATH}`)
}

function releaseRegistryLock() {
  rmSync(REGISTRY_LOCK_PATH, { force: true })
}

function readRegistry() {
  if (!existsSync(REGISTRY_PATH)) return { version: 1, worktrees: [] }
  return JSON.parse(readFileSync(REGISTRY_PATH, "utf8"))
}

function writeRegistry(registry) {
  mkdirSync(REGISTRY_PARENT_DIR, { recursive: true })
  const tempRegistryPath = `${REGISTRY_PATH}.tmp.${process.pid}`
  writeFileSync(tempRegistryPath, `${JSON.stringify(registry, null, 2)}\n`)
  rmSync(REGISTRY_PATH, { force: true })
  renameSync(tempRegistryPath, REGISTRY_PATH)
}

function parseMigrationNumbers(content) {
  const migrationNumbers = new Set()
  for (const contentLine of content.split("\n")) {
    const prefixMatch = /(\d{3})_/.exec(contentLine.trim())
    if (prefixMatch) migrationNumbers.add(Number(prefixMatch[1]))
  }
  return [...migrationNumbers]
}

function migrationNumbersInOriginMain() {
  try {
    execFileSync("git", ["fetch", "origin", "main"], { cwd: process.cwd(), stdio: "ignore", timeout: 20000 })
    const treeListing = execFileSync(
      "git",
      ["ls-tree", "-r", "--name-only", "origin/main", "--", "backend/migrations"],
      { cwd: process.cwd(), encoding: "utf8", timeout: 20000 },
    )
    return parseMigrationNumbers(treeListing)
  } catch {
    return []
  }
}

function migrationNumbersOnDisk() {
  const migrationsDirectory = path.join(process.cwd(), "backend", "migrations")
  if (!existsSync(migrationsDirectory)) return []
  return parseMigrationNumbers(readdirSync(migrationsDirectory).join("\n"))
}

function formatMigrationNumber(migrationNumber) {
  return String(migrationNumber).padStart(3, "0")
}

function parseFlagArguments(flagTokens) {
  const parsedArguments = {}
  for (let tokenIndex = 0; tokenIndex < flagTokens.length; tokenIndex += 1) {
    const token = flagTokens[tokenIndex]
    if (!token.startsWith("--")) continue
    const flagName = token.slice(2)
    const nextToken = flagTokens[tokenIndex + 1]
    if (nextToken === undefined || nextToken.startsWith("--")) {
      parsedArguments[flagName] = true
    } else {
      parsedArguments[flagName] = nextToken
      tokenIndex += 1
    }
  }
  return parsedArguments
}

function assertValidSlug(slug) {
  if (!/^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(slug)) {
    throw new Error(`Invalid slug '${slug}'. Use kebab-case letters and digits only (ex: onboarding-form).`)
  }
}

function assertValidWorktreeType(worktreeType) {
  if (!VALID_WORKTREE_TYPES.includes(worktreeType)) {
    throw new Error(`Invalid worktree type '${worktreeType}'. Allowed: ${VALID_WORKTREE_TYPES.join(", ")}.`)
  }
}

function commandList(jsonOutput) {
  const registry = readRegistry()
  if (jsonOutput) {
    console.log(JSON.stringify(registry, null, 2))
    return
  }
  if (registry.worktrees.length === 0) {
    console.log("No active worktrees registered.")
    return
  }
  for (const registeredWorktree of registry.worktrees) {
    const reservedMigrations =
      registeredWorktree.migrations.length > 0
        ? registeredWorktree.migrations.map(formatMigrationNumber).join(",")
        : "-"
    console.log(
      `${registeredWorktree.status.padEnd(8)} ${registeredWorktree.slug.padEnd(34)} ${registeredWorktree.branch.padEnd(42)} migrations=[${reservedMigrations}]`,
    )
  }
}

function commandClaim(commandArgs) {
  const slug = commandArgs.slug
  const worktreeType = commandArgs.type
  const explicitBranch = commandArgs.branch
  const note = commandArgs.note
  if (!slug || !worktreeType) {
    throw new Error("Usage: claim --slug <slug> --type <type> [--branch <branch>] [--note <text>]")
  }
  assertValidSlug(slug)
  assertValidWorktreeType(worktreeType)
  const resolvedBranch = explicitBranch || `${worktreeType}/${slug}`
  const worktreeDirectory = path.join(REGISTRY_PARENT_DIR, slug)
  if (!existsSync(worktreeDirectory)) {
    throw new Error(
      `Worktree directory not found: ${worktreeDirectory}. Run 'git worktree add ../healthcare-worktrees/${slug} -b ${resolvedBranch} origin/main' before claiming.`,
    )
  }
  acquireRegistryLock()
  try {
    const registry = readRegistry()
    const collision = registry.worktrees.find(
      (registeredWorktree) => registeredWorktree.slug === slug || registeredWorktree.branch === resolvedBranch,
    )
    if (collision) {
      throw new Error(
        `Registry collision with ${collision.slug} (${collision.branch}). Choose another slug or release the existing entry first.`,
      )
    }
    const registeredWorktree = {
      slug,
      type: worktreeType,
      branch: resolvedBranch,
      path: `../healthcare-worktrees/${slug}`,
      status: "active",
      claimedAt: new Date().toISOString(),
      migrations: [],
    }
    if (note) registeredWorktree.note = note
    registry.worktrees.push(registeredWorktree)
    writeRegistry(registry)
    console.log(`Claimed ${resolvedBranch} at ../healthcare-worktrees/${slug}`)
  } finally {
    releaseRegistryLock()
  }
}

function commandReserveMigration(commandArgs) {
  const slug = commandArgs.slug
  if (!slug) throw new Error("Usage: reserve-migration --slug <slug> [--name <snake_case_name>]")
  const originMainNumbers = migrationNumbersInOriginMain()
  const diskNumbers = migrationNumbersOnDisk()
  acquireRegistryLock()
  try {
    const registry = readRegistry()
    const registeredWorktree = registry.worktrees.find((candidate) => candidate.slug === slug)
    if (!registeredWorktree) {
      throw new Error(`No registered worktree for slug '${slug}'. Run claim first.`)
    }
    const reservedMigrationNumbers = new Set([...originMainNumbers, ...diskNumbers])
    for (const otherWorktree of registry.worktrees) {
      for (const migrationNumber of otherWorktree.migrations) reservedMigrationNumbers.add(migrationNumber)
    }
    for (const migrationNumber of registeredWorktree.migrations) reservedMigrationNumbers.add(migrationNumber)
    const nextMigrationNumber = Math.max(0, ...reservedMigrationNumbers) + 1
    registeredWorktree.migrations.push(nextMigrationNumber)
    writeRegistry(registry)
    const formattedMigrationName = commandArgs.name
      ? `_${commandArgs.name.replace(/[^a-z0-9_]+/g, "_")}`
      : ""
    console.log(
      `Reserved migration ${formatMigrationNumber(nextMigrationNumber)} for ${slug}: backend/migrations/${formatMigrationNumber(nextMigrationNumber)}${formattedMigrationName}_{up|down}.sql`,
    )
  } finally {
    releaseRegistryLock()
  }
}

function commandRelease(commandArgs) {
  const slug = commandArgs.slug
  if (!slug) throw new Error("Usage: release --slug <slug>")
  acquireRegistryLock()
  try {
    const registry = readRegistry()
    const releasedWorktree = registry.worktrees.find((registeredWorktree) => registeredWorktree.slug === slug)
    if (!releasedWorktree) throw new Error(`No registered worktree for slug '${slug}'.`)
    registry.worktrees = registry.worktrees.filter((registeredWorktree) => registeredWorktree.slug !== slug)
    writeRegistry(registry)
    const freedMigrations =
      releasedWorktree.migrations.length > 0
        ? ` freed migration(s) ${releasedWorktree.migrations.map(formatMigrationNumber).join(", ")}`
        : ""
    console.log(`Released ${releasedWorktree.branch}${freedMigrations}`)
  } finally {
    releaseRegistryLock()
  }
}

function printUsage() {
  console.log(
    [
      "Worktree coordination registry helpers for parallel features",
      "  list [--json]                                   list active worktrees and reserved identifiers",
      "  claim --slug <kebab-slug> --type <tipo> [--branch <branch>] [--note <text>]",
      "  reserve-migration --slug <kebab-slug> [--name <snake_case_name>]",
      "  release --slug <kebab-slug>",
      "",
      `Registry: ${REGISTRY_PATH}`,
    ].join("\n"),
  )
}

function main() {
  const [commandName, ...flagTokens] = process.argv.slice(2)
  const jsonOutput = flagTokens.includes("--json")
  try {
    switch (commandName) {
      case "list":
        commandList(jsonOutput)
        break
      case "claim":
        commandClaim(parseFlagArguments(flagTokens))
        break
      case "reserve-migration":
        commandReserveMigration(parseFlagArguments(flagTokens))
        break
      case "release":
        commandRelease(parseFlagArguments(flagTokens))
        break
      default:
        printUsage()
        if (commandName) process.exitCode = 1
    }
  } catch (error) {
    console.error(`[worktrees] ${error.message}`)
    process.exitCode = 1
  }
}

main()