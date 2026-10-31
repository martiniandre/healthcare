import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { Alert, AlertDescription, AlertTitle } from './alert'

describe('Alert', () => {
  it('should render the alert with role and content', () => {
    render(
      <Alert>
        <AlertTitle>Warning</AlertTitle>
        <AlertDescription>Something happened</AlertDescription>
      </Alert>
    )
    expect(screen.getByRole('alert')).toBeInTheDocument()
    expect(screen.getByText('Warning')).toBeInTheDocument()
    expect(screen.getByText('Something happened')).toBeInTheDocument()
  })

  it('should apply the border utility', () => {
    render(
      <Alert>
        <AlertDescription>Something happened</AlertDescription>
      </Alert>
    )
    expect(screen.getByRole('alert').className).toContain('border')
  })
})
