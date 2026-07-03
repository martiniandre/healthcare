import { renderHook, act } from '@testing-library/react'
import { useDebounce } from './useDebounce'

beforeEach(() => {
  vi.useFakeTimers()
})

afterEach(() => {
  vi.useRealTimers()
})

it('should return initial value immediately', () => {
  const { result } = renderHook(() => useDebounce('hello', 500))

  expect(result.current).toBe('hello')
})

it('should debounce value changes', () => {
  const { result, rerender } = renderHook(
    (props: { value: string; delay: number }) => useDebounce(props.value, props.delay),
    { initialProps: { value: 'hello', delay: 500 } }
  )

  expect(result.current).toBe('hello')

  rerender({ value: 'world', delay: 500 })

  expect(result.current).toBe('hello')

  act(() => {
    vi.advanceTimersByTime(500)
  })

  expect(result.current).toBe('world')
})

it('should clear timeout on unmount', () => {
  const clearTimeoutSpy = vi.spyOn(globalThis, 'clearTimeout')

  const { unmount } = renderHook(() => useDebounce('hello', 500))

  unmount()

  expect(clearTimeoutSpy).toHaveBeenCalled()
})
