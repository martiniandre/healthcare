import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { Table, TableHeader, TableBody, TableFooter, TableRow, TableHead, TableCell } from './Table'

describe('Table', () => {
  it('should render table structure', () => {
    render(
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow>
            <TableCell>John</TableCell>
          </TableRow>
        </TableBody>
        <TableFooter>
          <TableRow>
            <TableCell>Total</TableCell>
          </TableRow>
        </TableFooter>
      </Table>
    )
    expect(screen.getByText('Name')).toBeInTheDocument()
    expect(screen.getByText('John')).toBeInTheDocument()
    expect(screen.getByText('Total')).toBeInTheDocument()
  })

  it('should apply light border tokens to rows and footer', () => {
    render(
      <Table>
        <TableBody>
          <TableRow>
            <TableCell>John</TableCell>
          </TableRow>
        </TableBody>
        <TableFooter>
          <TableRow>
            <TableCell>Total</TableCell>
          </TableRow>
        </TableFooter>
      </Table>
    )
    const bodyRow = screen.getByText('John').closest('tr')
    expect(bodyRow?.className).toContain('border-b')
    expect(bodyRow?.className).toContain('border-border')

    const footer = screen.getByText('Total').closest('tfoot')
    expect(footer?.className).toContain('border-t')
    expect(footer?.className).toContain('border-border')
  })
})
