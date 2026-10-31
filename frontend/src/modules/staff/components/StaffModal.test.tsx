import { describe, it, expect, vi, beforeEach } from "vitest"
import { render, screen, fireEvent, waitFor } from "@testing-library/react"
import { StaffModal } from "./StaffModal"

const mutateAsync = vi.hoisted(() => vi.fn())
const departmentsState = vi.hoisted(() => ({
  current: { data: [] as Array<{ id: string; name: string }>, isLoading: false },
}))

vi.mock("../queries", () => ({
  useCreateEmployeeMutation: () => ({ isPending: false, mutateAsync }),
  useDepartmentsQuery: () => departmentsState.current,
}))

vi.mock("../../../shared/store/auth_store", () => ({
  useAuthStore: (selector: (state: { userId: string }) => unknown) => selector({ userId: "user-1" }),
}))

describe("StaffModal", () => {
  beforeEach(() => {
    mutateAsync.mockReset()
    departmentsState.current = {
      data: [{ id: "dep-1", name: "Cardiology" }, { id: "dep-2", name: "ICU" }],
      isLoading: false,
    }
  })

  it("should mark required fields and group them in fieldsets", () => {
    render(<StaffModal isOpen onClose={vi.fn()} />)

    expect(screen.getByLabelText(/Full Name/)).toBeRequired()
    expect(screen.getAllByRole("group", { name: /Identification|Category and department|Professional registration/ })).toHaveLength(3)
    expect(screen.getByLabelText(/Corporate E-mail/)).toHaveAttribute("autocomplete", "off")
  })

  it("should show an error summary and field errors after an invalid submit", async () => {
    render(<StaffModal isOpen onClose={vi.fn()} />)

    fireEvent.click(screen.getByRole("button", { name: "Save Registration" }))

    expect(await screen.findByText("Fix 4 fields to continue")).toBeInTheDocument()
    expect(screen.getByLabelText(/Full Name/)).toHaveAttribute("aria-invalid", "true")
    expect(screen.getByLabelText(/Full Name/)).toHaveAttribute(
      "aria-describedby",
      expect.stringContaining("staff-full-name-error"),
    )
    expect(mutateAsync).not.toHaveBeenCalled()
  })

  it("should submit the unchanged payload for valid data", async () => {
    mutateAsync.mockResolvedValue({ employeeId: "emp-1" })
    const onClose = vi.fn()
    render(<StaffModal isOpen onClose={onClose} />)

    fireEvent.change(screen.getByLabelText(/Full Name/), { target: { value: "Dra. Marina Alves" } })
    fireEvent.change(screen.getByLabelText(/Corporate E-mail/), { target: { value: "marina@clinica.com" } })
    fireEvent.change(screen.getByLabelText(/Professional License/), { target: { value: "CRM-SP 12345" } })
    fireEvent.click(screen.getByRole("combobox", { name: /Department/ }))
    fireEvent.click(await screen.findByRole("option", { name: "ICU" }))
    fireEvent.click(screen.getByRole("button", { name: "Save Registration" }))

    await waitFor(() =>
      expect(mutateAsync).toHaveBeenCalledWith({
        userId: "user-1",
        fullName: "Dra. Marina Alves",
        email: "marina@clinica.com",
        role: "DOCTOR",
        crmNumber: "CRM-SP 12345",
        departmentId: "dep-2",
      }),
    )
    await waitFor(() => expect(onClose).toHaveBeenCalled())
  })

  it("should preselect the department when there is only one", () => {
    departmentsState.current = { data: [{ id: "dep-1", name: "Cardiology" }], isLoading: false }
    render(<StaffModal isOpen onClose={vi.fn()} />)

    expect(screen.getByRole("combobox", { name: /Department/ })).toHaveTextContent("Cardiology")
  })

  it("should close directly when the form is pristine", () => {
    const onClose = vi.fn()
    render(<StaffModal isOpen onClose={onClose} />)

    fireEvent.click(screen.getByRole("button", { name: "Cancel" }))

    expect(onClose).toHaveBeenCalledTimes(1)
  })

  it("should ask for confirmation before discarding a dirty form", async () => {
    const onClose = vi.fn()
    render(<StaffModal isOpen onClose={onClose} />)

    fireEvent.change(screen.getByLabelText(/Full Name/), { target: { value: "Dra. Marina" } })
    fireEvent.click(screen.getByRole("button", { name: "Cancel" }))

    expect(onClose).not.toHaveBeenCalled()
    fireEvent.click(await screen.findByRole("button", { name: "Keep editing" }))
    expect(onClose).not.toHaveBeenCalled()

    fireEvent.click(screen.getByRole("button", { name: "Cancel" }))
    fireEvent.click(await screen.findByRole("button", { name: "Discard" }))
    expect(onClose).toHaveBeenCalledTimes(1)
  })
})
