import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { Dialog, DialogContent, DialogTitle } from './Dialog'

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

  it('should render an opaque backdrop overlay with blur when open', () => {
    render(
      <Dialog open>
        <DialogContent aria-describedby={undefined}>
          <DialogTitle>Test Title</DialogTitle>
        </DialogContent>
      </Dialog>
    )
    const overlay = document.querySelector('[data-state="open"].bg-black\\/60')
    expect(overlay).not.toBeNull()
    expect(overlay?.className).toContain('backdrop-blur-[2px]')
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
