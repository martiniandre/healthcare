import { auditLogsApi } from "./api"

vi.mock("../../shared/utils/http", () => ({
  http: {
    get: vi.fn(),
    post: vi.fn(),
  },
}))

import { http } from "../../shared/utils/http"

describe("auditLogsApi", () => {
  afterEach(() => {
    vi.clearAllMocks()
  })

  it("should list audit logs without filters", async () => {
    vi.mocked(http.get).mockResolvedValue({ audit_logs: [], total: 0 })

    await auditLogsApi.listAuditLogs({})

    expect(http.get).toHaveBeenCalledWith("/audit-logs")
  })

  it("should append every provided filter to the query string", async () => {
    vi.mocked(http.get).mockResolvedValue({ audit_logs: [], total: 0 })

    await auditLogsApi.listAuditLogs({
      action: "PAGE_VIEW",
      email: "medico@clinica.com",
      status: "GRANTED",
      startDate: "2026-01-01",
      endDate: "2026-01-31",
    })

    expect(http.get).toHaveBeenCalledWith(
      "/audit-logs?action=PAGE_VIEW&email=medico%40clinica.com&status=GRANTED&startDate=2026-01-01&endDate=2026-01-31"
    )
  })

  it("should skip the All sentinel values", async () => {
    vi.mocked(http.get).mockResolvedValue({ audit_logs: [], total: 0 })

    await auditLogsApi.listAuditLogs({ action: "All", status: "All" })

    expect(http.get).toHaveBeenCalledWith("/audit-logs")
  })

  it("should create an audit log entry", async () => {
    vi.mocked(http.post).mockResolvedValue(undefined)

    await auditLogsApi.createAuditLog({
      method: "PAGE_VIEW",
      correlation_id: "Viewed page: /dashboard",
      access_granted: true,
    })

    expect(http.post).toHaveBeenCalledWith("/audit-logs", {
      method: "PAGE_VIEW",
      correlation_id: "Viewed page: /dashboard",
      access_granted: true,
    })
  })
})
