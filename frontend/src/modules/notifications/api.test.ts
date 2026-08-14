import { notificationsApi } from "./api"

vi.mock("../../shared/utils/http", () => ({
  http: {
    get: vi.fn(),
    post: vi.fn(),
  },
}))

import { http } from "../../shared/utils/http"

describe("notificationsApi", () => {
  afterEach(() => {
    vi.clearAllMocks()
  })

  it("should list notifications with the default pagination", async () => {
    vi.mocked(http.get).mockResolvedValue({ notifications: [], total: 0 })

    await notificationsApi.list()

    expect(http.get).toHaveBeenCalledWith("/notifications?limit=50&offset=0")
  })

  it("should list notifications with custom pagination", async () => {
    vi.mocked(http.get).mockResolvedValue({ notifications: [], total: 0 })

    await notificationsApi.list(20, 40)

    expect(http.get).toHaveBeenCalledWith("/notifications?limit=20&offset=40")
  })

  it("should mark a notification as read", async () => {
    vi.mocked(http.post).mockResolvedValue(undefined)

    await notificationsApi.markRead("notification-1")

    expect(http.post).toHaveBeenCalledWith("/notifications/notification-1/read")
  })

  it("should fetch the unread count", async () => {
    vi.mocked(http.get).mockResolvedValue({ unread_count: 3 })

    await notificationsApi.getUnreadCount()

    expect(http.get).toHaveBeenCalledWith("/notifications/unread-count")
  })
})
