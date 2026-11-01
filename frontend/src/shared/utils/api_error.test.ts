import { extractApiErrorDetails } from "./api_error"

describe("extractApiErrorDetails", () => {
  it("should read code, message and reason from the response body", () => {
    const requestError = {
      response: { data: { code: "INVALID_DICOM", error: "invalid file", reason: "missing preamble" } },
    }

    expect(extractApiErrorDetails(requestError)).toEqual({
      code: "INVALID_DICOM",
      message: "invalid file",
      reason: "missing preamble",
    })
  })

  it("should ignore blank and non-string values", () => {
    const requestError = { response: { data: { code: 422, error: "  ", reason: null } } }

    expect(extractApiErrorDetails(requestError)).toEqual({
      code: undefined,
      message: undefined,
      reason: undefined,
    })
  })

  it("should return an empty object for errors without a response body", () => {
    expect(extractApiErrorDetails(new Error("network"))).toEqual({})
    expect(extractApiErrorDetails(null)).toEqual({})
    expect(extractApiErrorDetails({ response: { data: "text" } })).toEqual({})
  })
})
