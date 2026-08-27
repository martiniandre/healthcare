import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { Checkbox } from './Checkbox'

describe('Checkbox', () => {
  it('should render a checkbox input', () => {
    render(<Checkbox aria-label="Accept" />)
    expect(screen.getByRole('checkbox')).toBeInTheDocument()
  })

  it('should toggle when clicked', () => {
    const onCheckedChange = vi.fn()
    render(<Checkbox onCheckedChange={onCheckedChange} />)
    fireEvent.click(screen.getByRole('checkbox'))
    expect(onCheckedChange).toHaveBeenCalledTimes(1)
  })

  it('should be checked when the checked prop is set', () => {
    render(<Checkbox checked aria-label="Done" />)
    expect(screen.getByRole('checkbox')).toHaveAttribute('data-state', 'checked')
  })
})
