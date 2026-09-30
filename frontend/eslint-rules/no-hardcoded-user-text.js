const PROSE_PATTERN = /[A-Za-zÀ-ÿ]{2,}/

const ALLOWED_SHORT_TEXT = new Set([
  "°C",
  "°F",
  "BPM",
  "SpO₂",
  "SpO2",
  "UID",
  "DICOM",
  "DCM",
  "CRM",
  "COREN",
  "N/I",
  "N/A",
  "N/D",
  "FC",
  "VS",
  "IA",
  "AP",
  "DP",
  "MB",
  "KB",
  "GB",
  "HealthCare",
  "FHIR R4 · gRPC-Web",
  "ECG",
  "PACS",
  "FHIR",
  "CT",
  "MR",
  "US",
  "DX",
  "CR",
  "v",
])

const normalizeText = (rawText) => rawText.replace(/\s+/g, " ").trim()

const isProse = (rawText) => {
  const normalizedText = normalizeText(rawText)
  if (normalizedText.length === 0) {
    return false
  }
  if (ALLOWED_SHORT_TEXT.has(normalizedText)) {
    return false
  }
  return PROSE_PATTERN.test(normalizedText)
}

const hasEscapeHatch = (node) => {
  let currentNode = node
  while (currentNode && currentNode.type !== "Program") {
    if (currentNode.type === "JSXOpeningElement" || currentNode.type === "JSXAttribute") {
      const attributeName = currentNode.type === "JSXAttribute"
        ? currentNode.name?.name
        : currentNode.attributes?.some?.(
            (attribute) => attribute.name?.name === "data-i18n-allow",
          )
      if (attributeName === "data-i18n-allow") {
        return true
      }
    }
    currentNode = currentNode.parent
  }
  return false
}

const userFacingAttributes = new Set(["aria-label", "placeholder", "title", "alt"])

export const noHardcodedUserText = {
  meta: {
    type: "problem",
    docs: {
      description:
        "Disallow hardcoded user-facing prose in JSX text nodes and user-facing attributes. Use i18next keys via useTranslation instead.",
    },
    schema: [],
    messages: {
      hardcodedJsxText:
        'Hardcoded user-facing text "{{text}}". Move it to an i18next key and render it with t().',
      hardcodedAttribute:
        'Hardcoded user-facing text in "{{attribute}}" attribute: "{{text}}". Move it to an i18next key and render it with t().',
    },
  },
  create(context) {
    const reportAttribute = (attributeNode) => {
      const attributeName = attributeNode.name?.name
      if (!userFacingAttributes.has(attributeName)) {
        return
      }
      const attributeValue = attributeNode.value
      if (!attributeValue || attributeValue.type !== "Literal" || typeof attributeValue.value !== "string") {
        return
      }
      if (!isProse(attributeValue.value)) {
        return
      }
      context.report({
        node: attributeNode,
        messageId: "hardcodedAttribute",
        data: { attribute: attributeName, text: normalizeText(attributeValue.value) },
      })
    }

    return {
      JSXText(node) {
        if (!isProse(node.value)) {
          return
        }
        if (hasEscapeHatch(node)) {
          return
        }
        context.report({
          node,
          messageId: "hardcodedJsxText",
          data: { text: normalizeText(node.value) },
        })
      },
      JSXAttribute: reportAttribute,
    }
  },
}
