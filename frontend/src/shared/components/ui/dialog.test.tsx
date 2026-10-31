import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { Dialog, DialogContent, DialogTitle } from './dialog'

describe('Dialog', () => {
  it('should render the dialog content when open', () => {
    render(
      <Dialog open>
        <DialogContent aria-describedby={undefined}>
          <DialogTitle>Test Title</DialogTitle>
        </DialogContent>
      </Dialog>
    )
    expect(screen.getByRole('dialog')).toBeInTheDocument()
  })

  it('should render a backdrop overlay when open', () => {
    render(
      <Dialog open>
        <DialogContent aria-describedby={undefined}>
          <DialogTitle>Test Title</DialogTitle>
        </DialogContent>
      </Dialog>
    )
    const overlay = document.querySelector('[data-state="open"].bg-black\\/50')
    expect(overlay).not.toBeNull()
  })

  it('should not render the dialog content when closed', () => {
    render(
      <Dialog open={false}>
        <DialogContent aria-describedby={undefined}>
          <DialogTitle>Test Title</DialogTitle>
        </DialogContent>
      </Dialog>
    )
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument()
  })
})
