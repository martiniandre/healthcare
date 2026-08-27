import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { Separator } from './Separator'

describe('Separator', () => {
  it('should render a horizontal separator by default', () => {
    const { container } = render(<Separator decorative={false} />)
    const separator = screen.getByRole('separator')
    expect(separator).toHaveTextContent('')
    expect(container.querySelector('[data-orientation="horizontal"]')).not.toBeNull()
  })

  it('should render a vertical separator when orientation is vertical', () => {
    const { container } = render(<Separator decorative={false} orientation="vertical" />)
    expect(container.querySelector('[data-orientation="vertical"]')).not.toBeNull()
  })
})
