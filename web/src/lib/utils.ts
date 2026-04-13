/**
 * Utility Functions
 *
 * Shared utility functions used across the application.
 */

import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

/**
 * Combines multiple class values into a single className string.
 * Uses clsx for conditional classes and tailwind-merge to resolve conflicts.
 *
 * @param inputs - Class values to combine
 * @returns Merged className string
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
