#!/usr/bin/env node
import { readFileSync, readdirSync } from "node:fs"
import path from "node:path"

const STRICT = process.env.CONVENTIONS_STRICT === "1" || process.argv.includes("--strict")
const VERBOSE = process.argv.includes("--verbose")
const PROJECT_ROOT = process.cwd()
const BACKEND_ROOT = path.join(PROJECT_ROOT, "backend")
const FRONTEND_ROOT = path.join(PROJECT_ROOT, "frontend")

const IGNORED_DIRS = new Set([
  "node_modules",
  "dist",
  ".vite",
  "coverage",
  "test-results",
  "playwright-report",
  "build",
  ".git",
  "vendor",
])
const GO_EXTENSIONS = new Set([".go"])
const TS_EXTENSIONS = new Set([".ts", ".tsx", ".js", ".jsx", ".mjs", ".cjs"])
const DECLARATION_LINE_PATTERNS = [
  /^\s*(func|type|const|var|package|import)\b/,
  /^\s*(function|const|let|var|type|interface|class|export|import|default)\b/,
]

const GO_SINGLE_LETTER_PATTERNS = [
  { pattern: /func\s*\(\s*([a-z])\s+([*a-zA-Z.]+)/, kind: "receiver" },
  { pattern: /(?<![\w$])([a-z])\s*:=/, kind: "var" },
  { pattern: /\bvar\s+([a-z])\b/, kind: "var" },
]

const TS_SINGLE_LETTER_PATTERNS = [
  { pattern: /\b(?:const|let|var)\s+([a-zA-Z])\b/, kind: "decl" },
  { pattern: /\(\s*([a-zA-Z])\s*\)\s*=>/, kind: "arrow" },
  { pattern: /\(\s*([a-zA-Z])\s*,\s*([a-zA-Z])\s*\)\s*=>/, kind: "arrow" },
  { pattern: /(?<![\w$])([a-zA-Z])\s*=>/, kind: "arrow" },
]

const SWAGGER_DIRECTIVE_PATTERN = /^\/\/\s+@/

function walkFiles(directory, extensions, results) {
  let entries
  try {
    entries = readdirSync(directory, { withFileTypes: true })
  } catch {
    return
  }
  for (const entry of entries) {
    const entryPath = path.join(directory, entry.name)
    if (entry.isDirectory()) {
      if (!IGNORED_DIRS.has(entry.name) && entry.name !== "pb") walkFiles(entryPath, extensions, results)
    } else if (entry.isFile() && extensions.has(path.extname(entry.name))) {
      results.push(entryPath)
    }
  }
}

function isGeneratedFile(content) {
  const firstLines = content.slice(0, 2000)
  return /Code generated .* DO NOT EDIT/.test(firstLines) || firstLines.startsWith("// Code generated")
}

function lineNumberAt(rawContent, index) {
  let lineCount = 1
  for (let cursor = 0; cursor < index && cursor < rawContent.length; cursor += 1) {
    if (rawContent[cursor] === "\n") lineCount += 1
  }
  return lineCount
}

function tokenize(rawContent) {
  const comments = []
  const codeSegments = []
  let state = "code"
  let codeStart = 0
  let blockCommentStart = 0
  let index = 0
  while (index < rawContent.length) {
    const char = rawContent[index]
    const nextChar = rawContent[index + 1]
    if (state === "lineComment") {
      if (char === "\n") {
        comments.push({ start: codeStart - 1, end: index, type: "line", text: rawContent.slice(codeStart - 1, index) })
        codeStart = index + 1
        state = "code"
      }
      index += 1
      continue
    }
    if (state === "blockComment") {
      if (char === "*" && nextChar === "/") {
        comments.push({ start: blockCommentStart, end: index + 2, type: "block", text: rawContent.slice(blockCommentStart, index + 2) })
        codeStart = index + 2
        state = "code"
        index += 2
        continue
      }
      index += 1
      continue
    }
    if (state === "string_dq" || state === "string_sq" || state === "string_bt") {
      if (char === "\\") {
        index += 2
        continue
      }
      const closing = state === "string_dq" ? '"' : state === "string_sq" ? "'" : "`"
      if (char === closing) state = "code"
      index += 1
      continue
    }
    if (char === "/" && nextChar === "/") {
      if (codeStart < index) codeSegments.push({ start: codeStart, end: index, text: rawContent.slice(codeStart, index) })
      codeStart = index
      state = "lineComment"
      index += 2
      continue
    }
    if (char === "/" && nextChar === "*") {
      if (codeStart < index) codeSegments.push({ start: codeStart, end: index, text: rawContent.slice(codeStart, index) })
      blockCommentStart = index
      codeStart = index
      state = "blockComment"
      index += 2
      continue
    }
    if (char === '"' || char === "'" || char === "`") {
      if (codeStart < index) codeSegments.push({ start: codeStart, end: index, text: rawContent.slice(codeStart, index) })
      codeStart = index + 1
      state = char === '"' ? "string_dq" : char === "'" ? "string_sq" : "string_bt"
      index += 1
      continue
    }
    index += 1
  }
  if (state === "lineComment" && codeStart - 1 < rawContent.length) {
    comments.push({ start: codeStart - 1, end: rawContent.length, type: "line", text: rawContent.slice(codeStart - 1) })
  }
  if (state === "blockComment") {
    comments.push({ start: blockCommentStart, end: rawContent.length, type: "block", text: rawContent.slice(blockCommentStart) })
  }
  if (state === "code" && codeStart < rawContent.length) {
    codeSegments.push({ start: codeStart, end: rawContent.length, text: rawContent.slice(codeStart) })
  }
  return { comments, codeSegments }
}

function splitLines(rawContent) {
  return rawContent.split("\n")
}

function isLineCommentLine(line) {
  return /^\s*[/][/]/.test(line)
}

function isSwaggerLine(line) {
  return SWAGGER_DIRECTIVE_PATTERN.test(line.trim())
}

function findCommentBlock(rawContent, commentLineIndex) {
  const lines = splitLines(rawContent)
  let blockStart = commentLineIndex
  let blockEnd = commentLineIndex
  while (blockStart > 0 && isLineCommentLine(lines[blockStart - 1])) blockStart -= 1
  while (blockEnd < lines.length - 1 && isLineCommentLine(lines[blockEnd + 1])) blockEnd += 1
  return lines.slice(blockStart, blockEnd + 1)
}

function blockContainsSwagger(blockLines) {
  return blockLines.some(isSwaggerLine)
}

function nextNonBlankDeclares(rawContent, commentEndIndex) {
  const followingLines = rawContent.slice(commentEndIndex).split("\n")
  for (const candidateLine of followingLines) {
    if (candidateLine.trim() === "") continue
    return DECLARATION_LINE_PATTERNS.some((linePattern) => linePattern.test(candidateLine))
  }
  return false
}

function commentIsExempt(rawContent, comment) {
  if (comment.type === "block") {
    const trimmed = comment.text.trim()
    const inner = trimmed.replace(/^\/\*/, "").replace(/\*\/$/, "").trim()
    if (SWAGGER_DIRECTIVE_PATTERN.test(inner)) return true
    if (trimmed.includes("eslint-disable")) return true
    return false
  }
  const trimmed = comment.text.trim()
  if (trimmed.startsWith("///")) return true
  if (SWAGGER_DIRECTIVE_PATTERN.test(trimmed)) return true
  if (trimmed.includes("eslint-disable")) return true
  if (trimmed.includes("@ts-ignore") || trimmed.includes("@ts-expect-error")) return true
  const commentLineIndex = lineNumberAt(rawContent, comment.start) - 1
  const blockLines = findCommentBlock(rawContent, commentLineIndex)
  if (blockContainsSwagger(blockLines)) return true
  return nextNonBlankDeclares(rawContent, comment.end)
}

function checkFile(filePath) {
  const rawContent = readFileSync(filePath, "utf8")
  if (isGeneratedFile(rawContent)) return []
  const findings = []
  const { comments, codeSegments } = tokenize(rawContent)
  const isGoFile = GO_EXTENSIONS.has(path.extname(filePath))
  const singleLetterPatterns = isGoFile ? GO_SINGLE_LETTER_PATTERNS : TS_SINGLE_LETTER_PATTERNS
  const seenViolations = new Set()
  for (const segment of codeSegments) {
    for (const { pattern, kind } of singleLetterPatterns) {
      const matches = segment.text.matchAll(new RegExp(pattern.source, "g"))
      for (const match of matches) {
        const offset = match.index
        const line = lineNumberAt(rawContent, segment.start + offset)
        const variables = Array.from(match.slice(1)).filter((variable) => variable)
        const receiverType = match[2] ?? ""
        if (kind === "receiver" && /testing/.test(receiverType)) continue
        for (const variable of variables) {
          const violationKey = `${line}:${variable}:${kind}`
          if (seenViolations.has(violationKey)) continue
          seenViolations.add(violationKey)
          findings.push({ severity: "error", line, message: `variável de uma letra '${variable}' (${kind})`, file: filePath })
        }
      }
    }
  }
  for (const comment of comments) {
    const line = lineNumberAt(rawContent, comment.start)
    if (commentIsExempt(rawContent, comment)) continue
    findings.push({
      severity: "warning",
      line,
      message: "comentário em código não-gerado (regra: zero comentários; mantenha o código autoexplicativo)",
      file: filePath,
      context: comment.text.trim().slice(0, 80),
    })
  }
  return findings
}

function main() {
  const fileArgsIndex = process.argv.indexOf("--files")
  let files = []
  if (fileArgsIndex !== -1) {
    files = process.argv.slice(fileArgsIndex + 1).filter((arg) => !arg.startsWith("--"))
  } else {
    walkFiles(BACKEND_ROOT, GO_EXTENSIONS, files)
    walkFiles(FRONTEND_ROOT, TS_EXTENSIONS, files)
  }
  const allFindings = []
  for (const filePath of files) {
    const absolutePath = path.resolve(PROJECT_ROOT, filePath)
    try {
      allFindings.push(...checkFile(absolutePath))
    } catch {
      console.log(`[SKIP ] ${filePath}  não encontrado ou ilegível`)
    }
  }
  const errorFindings = allFindings.filter((finding) => finding.severity === "error")
  const warningFindings = allFindings.filter((finding) => finding.severity === "warning")
  if (VERBOSE || allFindings.length > 0) {
    for (const finding of allFindings) {
      const badge = finding.severity === "error" ? "ERROR" : "WARN "
      const relativeFile = path.relative(PROJECT_ROOT, finding.file).replaceAll("\\", "/")
      const contextSuffix = finding.context ? `  // ${finding.context}` : ""
      console.log(`[${badge}] ${relativeFile}:${finding.line}  ${finding.message}${contextSuffix}`)
    }
  }
  console.log(`Conventions check: ${files.length} arquivos, ${errorFindings.length} erro(s), ${warningFindings.length} aviso(s)`)
  if (errorFindings.length > 0) process.exitCode = 1
  else if (STRICT && warningFindings.length > 0) process.exitCode = 1
}

main()
