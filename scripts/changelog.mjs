#!/usr/bin/env node
import { execFileSync } from "node:child_process"
import { writeFileSync } from "node:fs"
import path from "node:path"

const PROJECT_ROOT = process.cwd()
const CHANGELOG_PATH = path.join(PROJECT_ROOT, "CHANGELOG.md")

const TYPE_ORDER = ["feat", "fix", "docs", "refactor", "perf", "test", "style", "chore", "ci", "build"]
const TYPE_LABELS = {
  feat: "Features",
  fix: "Bug fixes",
  docs: "Documentation",
  refactor: "Refactors",
  perf: "Performance",
  test: "Tests",
  style: "Style",
  chore: "Chores",
  ci: "CI",
  build: "Build",
}
const COMMIT_PATTERN = /^(\w+)(\(([^)]+)\))?:\s+(.+)$/

function runGit(args) {
  return execFileSync("git", args, { encoding: "utf8" }).trim()
}

function loadCommits() {
  const rawLog = runGit(["log", "--pretty=format:%H|%s", "--date=short"])
  if (rawLog === "") return []
  return rawLog.split("\n").map((line) => {
    const [hash, ...subjectParts] = line.split("|")
    return { hash, subject: subjectParts.join("|") }
  })
}

function classify(commit) {
  const match = COMMIT_PATTERN.exec(commit.subject)
  if (!match) return { hash: commit.hash, type: "outros", scope: null, message: commit.subject }
  return { hash: commit.hash, type: match[1], scope: match[3] ?? null, message: match[4] }
}

function main() {
  const commits = loadCommits()
  if (commits.length === 0) {
    console.log("Nenhum commit encontrado; CHANGELOG.md não foi gerado.")
    return
  }
  const groups = new Map()
  for (const commit of commits) {
    const classified = classify(commit)
    if (!groups.has(classified.type)) groups.set(classified.type, [])
    groups.get(classified.type).push(classified)
  }
  const orderedTypes = TYPE_ORDER.filter((type) => groups.has(type))
  const remainingTypes = Array.from(groups.keys()).filter((type) => !TYPE_ORDER.includes(type)).sort()
  const sections = orderedTypes.concat(remainingTypes)

  const lines = ["# Changelog", "", "Todas as mudanças relevantes por tipo, geradas automaticamente a partir dos commits Conventional Commits.", ""]
  for (const type of sections) {
    const header = TYPE_LABELS[type] ?? type.charAt(0).toUpperCase() + type.slice(1)
    lines.push(`## ${header}`, "")
    for (const item of groups.get(type)) {
      const scopePrefix = item.scope ? `**${item.scope}:** ` : ""
      lines.push(`- ${scopePrefix}${item.message} (\`${item.hash.slice(0, 7)}\`)`)
    }
    lines.push("")
  }
  lines.push(`Últimos ${commits.length} commits incluídos.`, "")
  const changelogContent = lines.join("\n")
  writeFileSync(CHANGELOG_PATH, changelogContent, "utf8")
  console.log(`CHANGELOG.md gerado com ${commits.length} commits em ${sections.length} grupo(s).`)
}

main()
