import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
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
    expect(screen.getByText('Sucesso')).toBeInTheDocument()
  })

  it('should render nothing when there are no toasts', () => {
    render(<ToastContainer />)
    expect(screen.queryByText('Sucesso')).not.toBeInTheDocument()
  })
})
