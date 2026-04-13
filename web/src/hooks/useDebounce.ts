/**
 * useDebounce Hook
 *
 * Custom React hook for debouncing a value.
 * Useful for delaying API calls while user is typing.
 */

import { useState, useEffect } from 'react';

/**
 * Debounces a value by the specified delay.
 *
 * @param value - The value to debounce
 * @param delay - The delay in milliseconds
 * @returns The debounced value
 */
export function useDebounce<T>(value: T, delay: number): T {
  // ==================== State ====================
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  // ==================== Effects ====================

  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(timer);
    };
  }, [value, delay]);

  return debouncedValue;
}

export default useDebounce;
