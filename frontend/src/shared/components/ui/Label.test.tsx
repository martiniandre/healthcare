import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { Label } from './Label'

describe('Label', () => {
  it('should render the children', () => {
    render(<Label>First name</Label>)
    expect(screen.getByText('First name')).toBeInTheDocument()
  })

  it('should render with the label element semantic', () => {
    const { container } = render(<Label htmlFor="name">Name</Label>)
    expect(container.querySelector('label')?.getAttribute('for')).toBe('name')
  })

  it('should apply the base label typography token', () => {
    render(<Label>Email</Label>)
    expect(screen.getByText('Email').className).toContain('font-semibold')
  })
})
