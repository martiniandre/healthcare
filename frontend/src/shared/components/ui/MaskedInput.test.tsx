import { describe, it, expect, vi } from 'vitest'
import { render, fireEvent } from '@testing-library/react'
import { MaskedInput } from './MaskedInput'

describe('MaskedInput', () => {
  it('should format the slash based date mask from raw digits', () => {
    const onChange = vi.fn()
    const { container } = render(<MaskedInput mask="99/99/9999" onChange={onChange} />)
    const input = container.querySelector('input') as HTMLInputElement
    fireEvent.change(input, { target: { value: '15051990' } })
    expect(input.value).toBe('15/05/1990')
  })

  it('should keep the dash based date mask for ISO dates', () => {
    const { container } = render(<MaskedInput mask="9999-99-99" />)
    const input = container.querySelector('input') as HTMLInputElement
    fireEvent.change(input, { target: { value: '19900515' } })
    expect(input.value).toBe('1990-05-15')
  })

  it('should still forward the formatted value through the change event', () => {
    const onChange = vi.fn()
    const { container } = render(<MaskedInput mask="99/99/9999" onChange={onChange} />)
    const input = container.querySelector('input') as HTMLInputElement
    fireEvent.change(input, { target: { value: '01011990' } })
    expect(onChange).toHaveBeenCalledTimes(1)
    expect(input.value).toBe('01/01/1990')
  })

  it('should preserve non digit characters through the cpf mask', () => {
    const { container } = render(<MaskedInput mask="999.999.999-99" />)
    const input = container.querySelector('input') as HTMLInputElement
    fireEvent.change(input, { target: { value: '11122233344' } })
    expect(input.value).toBe('111.222.333-44')
  })
})