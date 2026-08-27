import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { Textarea } from './Textarea'

describe('Textarea', () => {
  it('should render a textarea element', () => {
    render(<Textarea placeholder="Write here" />)
    expect(screen.getByPlaceholderText('Write here').tagName).toBe('TEXTAREA')
  })

  it('should handle change events', () => {
    const onChange = vi.fn()
    render(<Textarea onChange={onChange} />)
    fireEvent.change(screen.getByRole('textbox'), { target: { value: 'note' } })
    expect(onChange).toHaveBeenCalledTimes(1)
  })

  it('should apply the surface border token', () => {
    render(<Textarea />)
    expect(screen.getByRole('textbox').className).toContain('border-input')
  })
})
