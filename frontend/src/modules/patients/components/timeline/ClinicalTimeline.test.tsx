import { describe, it, expect, vi, beforeEach } from "vitest"
import { render, screen, waitFor, fireEvent, within } from "@testing-library/react"
import { MemoryRouter, Route, Routes, useLocation } from "react-router-dom"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import type { TimelineItem, TimelinePage } from "../../types"
import ClinicalTimeline from "./ClinicalTimeline"

vi.mock("../../api", () => ({
  patientsApi: {
    getTimeline: vi.fn(),
  },
}))

import { patientsApi } from "../../api"

const buildItem = (overrides: Partial<TimelineItem> = {}): TimelineItem => ({
  id: "encounter:enc-1",
  type: "encounter",
  occurredAt: new Date(2026, 4, 10, 14, 30).toISOString(),
  clinicalDate: null,
  title: "Routine visit",
  summary: "Follow up",
  status: "finished",
  resourceId: "enc-1",
  encounterId: "enc-1",
  ...overrides,
})

const buildPage = (overrides: Partial<TimelinePage> = {}): TimelinePage => ({
  items: [buildItem()],
  nextCursor: null,
  partial: null,
  window: { from: "2025-11-01T00:00:00Z", to: "2026-05-01T00:00:00Z" },
  ...overrides,
})

const LocationProbe = () => {
  const location = useLocation()
  return <output data-testid="location">{`${location.pathname}${location.search}`}</output>
}

const renderTimeline = (initialUrl = "/patients/patient-1?tab=timeline") => {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[initialUrl]}>
        <LocationProbe />
        <Routes>
          <Route path="/patients/:id" element={<ClinicalTimeline patientId="patient-1" />} />
          <Route path="/patients/:id/encounters/:encounterId" element={<div>encounter page</div>} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>,
  )
}

describe("ClinicalTimeline", () => {
  beforeEach(() => {
    vi.mocked(patientsApi.getTimeline).mockReset()
  })

  it("should render a skeleton while the first page loads", () => {
    vi.mocked(patientsApi.getTimeline).mockReturnValue(new Promise(() => {}))
    renderTimeline()

    expect(screen.getByRole("status", { name: "Loading timeline..." })).toHaveAttribute("aria-busy", "true")
  })

  it("should group entries by day with the newest day first", async () => {
    vi.mocked(patientsApi.getTimeline).mockResolvedValue(
      buildPage({
        items: [
          buildItem({ id: "report:r-1", type: "report", title: "Chest report", resourceId: "r-1", status: "final" }),
          buildItem({ id: "encounter:enc-1" }),
          buildItem({
            id: "allergy:a-1",
            type: "allergy",
            title: "Penicillin",
            resourceId: "a-1",
            status: "active",
            encounterId: null,
            occurredAt: new Date(2026, 4, 8, 9, 0).toISOString(),
          }),
        ],
      }),
    )
    renderTimeline()

    const dayGroups = await screen.findAllByTestId("timeline-day")
    expect(dayGroups).toHaveLength(2)
    expect(within(dayGroups[0]).getAllByTestId("timeline-entry")).toHaveLength(2)
    expect(within(dayGroups[1]).getByText("Penicillin")).toBeInTheDocument()
    expect(within(dayGroups[0]).getByRole("heading", { level: 3 }).textContent).toContain("2026")
  })

  it("should show the localized type label, status badge and clinical date metadata", async () => {
    vi.mocked(patientsApi.getTimeline).mockResolvedValue(
      buildPage({
        items: [
          buildItem({
            id: "condition:c-1",
            type: "condition",
            title: "Hypertension",
            resourceId: "c-1",
            status: "active",
            encounterId: null,
            clinicalDate: new Date(2026, 3, 1, 8, 0).toISOString(),
          }),
        ],
      }),
    )
    renderTimeline()

    const entry = await screen.findByTestId("timeline-entry")
    expect(within(entry).getByText("Condition")).toBeInTheDocument()
    expect(within(entry).getByText("Active").closest("[data-status]")).toHaveAttribute("data-status", "attention")
    expect(within(entry).getByText(/Clinical date:/)).toBeInTheDocument()
    expect(within(entry).queryByRole("button", { name: "Open encounter" })).toBeNull()
  })

  it("should hide the clinical date when it matches the recorded time", async () => {
    const sameMoment = new Date(2026, 4, 10, 14, 30).toISOString()
    vi.mocked(patientsApi.getTimeline).mockResolvedValue(
      buildPage({ items: [buildItem({ clinicalDate: sameMoment, occurredAt: sameMoment })] }),
    )
    renderTimeline()

    await screen.findByTestId("timeline-entry")
    expect(screen.queryByText(/Clinical date:/)).toBeNull()
  })

  it("should navigate to the related encounter from an entry", async () => {
    vi.mocked(patientsApi.getTimeline).mockResolvedValue(
      buildPage({
        items: [buildItem({ id: "report:r-1", type: "report", resourceId: "r-1", encounterId: "enc-9", status: "final" })],
      }),
    )
    renderTimeline()

    fireEvent.click(await screen.findByRole("button", { name: "Open encounter" }))

    expect(await screen.findByText("encounter page")).toBeInTheDocument()
    expect(screen.getByTestId("location").textContent).toBe("/patients/patient-1/encounters/enc-9")
  })

  it("should request the default window without type filters first", async () => {
    vi.mocked(patientsApi.getTimeline).mockResolvedValue(buildPage())
    renderTimeline()

    await screen.findByTestId("timeline-entry")
    expect(patientsApi.getTimeline).toHaveBeenCalledWith("patient-1", {
      types: [],
      cursor: undefined,
      limit: 30,
    })
  })

  it("should refetch with the selected types and persist them in the url", async () => {
    vi.mocked(patientsApi.getTimeline).mockResolvedValue(buildPage())
    renderTimeline()
    await screen.findByTestId("timeline-entry")

    fireEvent.click(screen.getByRole("button", { name: "Report" }))
    fireEvent.click(screen.getByRole("button", { name: "Allergy" }))

    await waitFor(() =>
      expect(patientsApi.getTimeline).toHaveBeenLastCalledWith(
        "patient-1",
        expect.objectContaining({ types: ["allergy", "report"] }),
      ),
    )
    expect(screen.getByTestId("location").textContent).toContain("types=allergy%2Creport")
    expect(screen.getByRole("button", { name: "Report" })).toHaveAttribute("aria-pressed", "true")
    expect(screen.getByRole("button", { name: "All types" })).toHaveAttribute("aria-pressed", "false")
  })

  it("should clear the type selection with the all types chip", async () => {
    vi.mocked(patientsApi.getTimeline).mockResolvedValue(buildPage())
    renderTimeline("/patients/patient-1?tab=timeline&types=report")
    await screen.findByTestId("timeline-entry")

    fireEvent.click(screen.getByRole("button", { name: "All types" }))

    await waitFor(() =>
      expect(patientsApi.getTimeline).toHaveBeenLastCalledWith("patient-1", expect.objectContaining({ types: [] })),
    )
    expect(screen.getByTestId("location").textContent).not.toContain("types=")
  })

  it("should send a computed window for the thirty day period read from the url", async () => {
    vi.mocked(patientsApi.getTimeline).mockResolvedValue(buildPage())
    renderTimeline("/patients/patient-1?tab=timeline&period=last30Days")

    await screen.findByTestId("timeline-entry")
    const request = vi.mocked(patientsApi.getTimeline).mock.calls[0][1]
    expect(request?.from).toBeTruthy()
    expect(request?.to).toBeTruthy()
  })

  it("should send the custom range typed in the date inputs", async () => {
    vi.mocked(patientsApi.getTimeline).mockResolvedValue(buildPage())
    renderTimeline("/patients/patient-1?tab=timeline&period=custom")
    await screen.findByTestId("timeline-entry")

    fireEvent.change(screen.getByLabelText("From"), { target: { value: "2026-01-01" } })
    fireEvent.change(screen.getByLabelText("To"), { target: { value: "2026-01-31" } })

    await waitFor(() =>
      expect(patientsApi.getTimeline).toHaveBeenLastCalledWith(
        "patient-1",
        expect.objectContaining({
          from: new Date(2026, 0, 1, 0, 0, 0, 0).toISOString(),
          to: new Date(2026, 0, 31, 23, 59, 59, 999).toISOString(),
        }),
      ),
    )
  })

  it("should preselect the six month period in the select", async () => {
    vi.mocked(patientsApi.getTimeline).mockResolvedValue(buildPage())
    renderTimeline()
    await screen.findByTestId("timeline-entry")

    expect(screen.getByRole("combobox")).toHaveTextContent("Last 6 months")
  })

  it("should show the empty state with a hint to widen the period", async () => {
    vi.mocked(patientsApi.getTimeline).mockResolvedValue(buildPage({ items: [] }))
    renderTimeline()

    expect(await screen.findByText("No records in this period")).toBeInTheDocument()
    expect(screen.getByText(/Widen the period/)).toBeInTheDocument()
  })

  it("should show an error alert and retry the request", async () => {
    vi.mocked(patientsApi.getTimeline).mockRejectedValueOnce(new Error("boom"))
    vi.mocked(patientsApi.getTimeline).mockResolvedValue(buildPage())
    renderTimeline()

    expect(await screen.findByText("Could not load the data")).toBeInTheDocument()
    fireEvent.click(screen.getByRole("button", { name: "Try again" }))

    expect(await screen.findByTestId("timeline-entry")).toBeInTheDocument()
    expect(patientsApi.getTimeline).toHaveBeenCalledTimes(2)
  })

  it("should warn about failed types while still listing the other entries and retry", async () => {
    vi.mocked(patientsApi.getTimeline).mockResolvedValue(buildPage({ partial: { failedTypes: ["imaging", "report"] } }))
    renderTimeline()

    const warning = await screen.findByText("Incomplete timeline")
    expect(warning.closest("[data-slot='alert']")).toHaveTextContent("Imaging and Report")
    expect(screen.getByTestId("timeline-entry")).toBeInTheDocument()

    fireEvent.click(screen.getByRole("button", { name: "Try again" }))
    await waitFor(() => expect(patientsApi.getTimeline).toHaveBeenCalledTimes(2))
  })

  it("should show the warning instead of the empty state when only failures happened", async () => {
    vi.mocked(patientsApi.getTimeline).mockResolvedValue(buildPage({ items: [], partial: { failedTypes: ["allergy"] } }))
    renderTimeline()

    expect(await screen.findByText("Incomplete timeline")).toBeInTheDocument()
    expect(screen.queryByText("No records in this period")).toBeNull()
  })

  it("should load the next page with the cursor and append its entries", async () => {
    vi.mocked(patientsApi.getTimeline)
      .mockResolvedValueOnce(buildPage({ nextCursor: "cursor-2" }))
      .mockResolvedValueOnce(
        buildPage({
          items: [
            buildItem({
              id: "allergy:a-1",
              type: "allergy",
              title: "Penicillin",
              resourceId: "a-1",
              encounterId: null,
              status: "active",
              occurredAt: new Date(2026, 3, 2, 9, 0).toISOString(),
            }),
          ],
        }),
      )
    renderTimeline()

    fireEvent.click(await screen.findByRole("button", { name: "Load more" }))

    expect(await screen.findByText("Penicillin")).toBeInTheDocument()
    expect(patientsApi.getTimeline).toHaveBeenLastCalledWith("patient-1", expect.objectContaining({ cursor: "cursor-2" }))
    expect(screen.getAllByTestId("timeline-entry")).toHaveLength(2)
    expect(screen.queryByRole("button", { name: "Load more" })).toBeNull()
  })

  it("should show a spinner while the next page loads", async () => {
    vi.mocked(patientsApi.getTimeline)
      .mockResolvedValueOnce(buildPage({ nextCursor: "cursor-2" }))
      .mockReturnValueOnce(new Promise(() => {}))
    renderTimeline()

    fireEvent.click(await screen.findByRole("button", { name: "Load more" }))

    const loadingButton = await screen.findByRole("button", { name: /Loading more records/ })
    expect(loadingButton).toBeDisabled()
    expect(within(loadingButton).getByRole("status")).toBeInTheDocument()
  })
})
