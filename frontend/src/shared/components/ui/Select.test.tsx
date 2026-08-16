import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from './Select'

describe('Select', () => {
  it('should render the trigger with the input border token', () => {
    render(
      <Select>
        <SelectTrigger aria-label="Pick an option">
          <SelectValue placeholder="Select..." />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="one">One</SelectItem>
        </SelectContent>
      </Select>
    )
    const trigger = screen.getByRole('combobox')
    expect(trigger.className).toContain('border-input')
  })

  it('should render the trigger placeholder', () => {
    render(
      <Select>
        <SelectTrigger aria-label="Pick an option">
          <SelectValue placeholder="Select..." />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="one">One</SelectItem>
        </SelectContent>
      </Select>
    )
    expect(screen.getByText('Select...')).toBeInTheDocument()
  })
})
