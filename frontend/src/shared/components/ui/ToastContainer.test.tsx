import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { act, fireEvent, render, screen } from '@testing-library/react'
import { ToastContainer } from './ToastContainer'
import { useToastStore } from '../../store/toast_store'

describe('ToastContainer', () => {
  beforeEach(() => {
    useToastStore.setState({ toasts: [] })
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('should render the toast message with a light border token', () => {
    useToastStore.setState({
      toasts: [{ id: 'toast-1', message: 'Operação concluída', type: 'success' }],
    })
    render(<ToastContainer />)
    const toastElement = screen.getByText('Operação concluída').closest('.pointer-events-auto')
    expect(toastElement?.className).toContain('border-border')
  })

  it('should render the toast title for success type', () => {
    useToastStore.setState({
      toasts: [{ id: 'toast-1', message: 'Operação concluída', type: 'success' }],
    })
    render(<ToastContainer />)
    expect(screen.getByText('Success')).toBeInTheDocument()
  })

  it('should render nothing when there are no toasts', () => {
    render(<ToastContainer />)
    expect(screen.queryByText('Success')).not.toBeInTheDocument()
  })
})

describe('ToastContainer accessibility and timing', () => {
  beforeEach(() => {
    useToastStore.setState({ toasts: [] })
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('should announce errors as alerts and keep them until dismissed', () => {
    useToastStore.getState().addToast('Falha ao salvar', 'error')
    render(<ToastContainer />)
    expect(screen.getByRole('alert')).toHaveTextContent('Falha ao salvar')
    act(() => {
      vi.advanceTimersByTime(60_000)
    })
    expect(screen.getByRole('alert')).toBeInTheDocument()
  })

  it('should announce success politely and dismiss it after the default duration', () => {
    useToastStore.getState().addToast('Salvo', 'success')
    render(<ToastContainer />)
    expect(screen.getByRole('status')).toHaveTextContent('Salvo')
    act(() => {
      vi.advanceTimersByTime(6_100)
    })
    expect(screen.queryByRole('status')).not.toBeInTheDocument()
  })

  it('should pause the dismissal while hovered', () => {
    useToastStore.getState().addToast('Salvo', 'success')
    render(<ToastContainer />)
    fireEvent.mouseEnter(screen.getByRole('status'))
    act(() => {
      vi.advanceTimersByTime(10_000)
    })
    expect(screen.getByRole('status')).toBeInTheDocument()
  })
})
