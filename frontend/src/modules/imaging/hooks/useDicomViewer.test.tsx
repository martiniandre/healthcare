import { render, fireEvent, screen, waitFor } from "@testing-library/react"
import { useEffect, useState } from "react"
import { useDicomViewer } from "./useDicomViewer"
import type { ImagingStudy } from "../types"

function createFakeCanvasContext(): Record<string, unknown> {
  return {
    fillRect: vi.fn(),
    save: vi.fn(),
    translate: vi.fn(),
    scale: vi.fn(),
    restore: vi.fn(),
    beginPath: vi.fn(),
    arc: vi.fn(),
    fill: vi.fn(),
    strokeRect: vi.fn(),
    moveTo: vi.fn(),
    lineTo: vi.fn(),
    stroke: vi.fn(),
    fillText: vi.fn(),
    measureText: () => ({ width: 0 }),
  }
}

function ViewerHarness({ study }: { study?: ImagingStudy | null }) {
  const {
    activeTool,
    canvasReference,
    interfaceConfiguration,
    setActiveTool,
    applyPreset,
    handleMouseDown,
    handleMouseMove,
    handleMouseUp,
  } = useDicomViewer(study)
  const [configuration, setConfiguration] = useState({
    contrast: 1,
    brightness: 1,
    zoom: 1,
  })
  useEffect(() => {
    const syncTimer = setInterval(() => {
      setConfiguration({ ...interfaceConfiguration })
    }, 25)
    return () => clearInterval(syncTimer)
  }, [interfaceConfiguration])
  return (
    <div>
      <canvas ref={canvasReference} width={800} height={600} data-testid="viewer-canvas" />
      <div
        data-testid="drag-area"
        onMouseDown={(event) => {
          handleMouseDown(event)
        }}
        onMouseMove={(event) => {
          handleMouseMove(event)
        }}
        onMouseUp={(event) => {
          handleMouseUp(event)
        }}
      />
      <span data-testid="active-tool">{activeTool}</span>
      <span data-testid="contrast">{configuration.contrast}</span>
      <span data-testid="brightness">{configuration.brightness}</span>
      <span data-testid="zoom">{configuration.zoom}</span>
      <button
        onClick={() => {
          setActiveTool("contrast")
        }}
      >
        use-contrast
      </button>
      <button
        onClick={() => {
          applyPreset("bone")
        }}
      >
        apply-bone
      </button>
      <button
        onClick={() => {
          applyPreset("lung")
        }}
      >
        apply-lung
      </button>
      <button
        onClick={() => {
          applyPreset("soft")
        }}
      >
        apply-soft
      </button>
    </div>
  )
}

describe("useDicomViewer", () => {
  const testStudy: ImagingStudy = {
    id: "study-1",
    patient_fhir_id: "patient-1",
    title: "Tomografia de tórax",
    modality: "CT",
    status: "available",
    study_instance_uid: "1.2.3.4.5",
    created_at: "2026-01-01T10:00:00Z",
  }

  beforeEach(() => {
    vi.spyOn(HTMLCanvasElement.prototype, "getContext").mockReturnValue(
      createFakeCanvasContext() as unknown as CanvasRenderingContext2D
    )
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it("should default to the zoom tool with neutral image configuration", () => {
    render(<ViewerHarness study={testStudy} />)

    expect(screen.getByTestId("active-tool").textContent).toBe("zoom")
    expect(screen.getByTestId("contrast").textContent).toBe("1")
    expect(screen.getByTestId("brightness").textContent).toBe("1")
    expect(screen.getByTestId("zoom").textContent).toBe("1")
  })

  it("should apply the bone preset and expose updated configuration", async () => {
    render(<ViewerHarness study={testStudy} />)

    fireEvent.click(screen.getByText("apply-bone"))

    await waitFor(() => {
      expect(screen.getByTestId("contrast").textContent).toBe("2")
      expect(screen.getByTestId("brightness").textContent).toBe("0.8")
    })
  })

  it("should apply the lung preset and expose updated configuration", async () => {
    render(<ViewerHarness study={testStudy} />)

    fireEvent.click(screen.getByText("apply-lung"))

    await waitFor(() => {
      expect(screen.getByTestId("contrast").textContent).toBe("0.6")
      expect(screen.getByTestId("brightness").textContent).toBe("1.4")
    })
  })

  it("should reset to neutral values with the soft preset", async () => {
    render(<ViewerHarness study={testStudy} />)

    fireEvent.click(screen.getByText("apply-bone"))
    fireEvent.click(screen.getByText("apply-soft"))

    await waitFor(() => {
      expect(screen.getByTestId("contrast").textContent).toBe("1")
      expect(screen.getByTestId("brightness").textContent).toBe("1")
    })
  })

  it("should zoom in when dragging upwards with the zoom tool", async () => {
    render(<ViewerHarness study={testStudy} />)
    const dragArea = screen.getByTestId("drag-area")

    fireEvent.mouseDown(dragArea, { clientX: 0, clientY: 0 })
    fireEvent.mouseMove(dragArea, { clientX: 0, clientY: -100 })
    fireEvent.mouseUp(dragArea)

    await waitFor(() => {
      expect(screen.getByTestId("zoom").textContent).toBe("1.5")
    })
  })

  it("should clamp zoom to the maximum allowed value", async () => {
    render(<ViewerHarness study={testStudy} />)
    const dragArea = screen.getByTestId("drag-area")

    fireEvent.mouseDown(dragArea, { clientX: 0, clientY: 0 })
    fireEvent.mouseMove(dragArea, { clientX: 0, clientY: -5000 })
    fireEvent.mouseUp(dragArea)

    await waitFor(() => {
      expect(screen.getByTestId("zoom").textContent).toBe("5")
    })
  })

  it("should switch the active tool", () => {
    render(<ViewerHarness study={testStudy} />)

    fireEvent.click(screen.getByText("use-contrast"))

    expect(screen.getByTestId("active-tool").textContent).toBe("contrast")
  })

  it("should render without crashing when no study is provided", () => {
    render(<ViewerHarness study={null} />)

    expect(screen.getByTestId("viewer-canvas")).toBeInTheDocument()
  })
})
