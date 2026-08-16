import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { Badge } from './Badge'

describe('Badge', () => {
  it('should render the badge label', () => {
    render(<Badge>Active</Badge>)
    expect(screen.getByText('Active')).toBeInTheDocument()
  })

  it('should render the default variant without a dark border', () => {
    render(<Badge>Active</Badge>)
    expect(screen.getByText('Active').className).toContain('border-transparent')
  })

  it('should render the outline variant with the light border token', () => {
    render(<Badge variant="outline">Outline</Badge>)
    const badge = screen.getByText('Outline')
    expect(badge.className).toContain('border-border')
    expect(badge.className).toContain('text-foreground')
  })
})
