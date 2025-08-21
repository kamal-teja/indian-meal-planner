import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { createContext } from 'react'
import AuthContext from './AuthContext'

// Mock localStorage for testing
const mockLocalStorage = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn()
}
global.localStorage = mockLocalStorage

describe('AuthContext', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockLocalStorage.getItem.mockReturnValue(null)
  })

  it('should provide AuthContext', () => {
    expect(AuthContext).toBeDefined()
    expect(typeof AuthContext).toBe('object')
  })

  it('should be a React context', () => {
    // Test that it's a valid React context
    expect(AuthContext._currentValue).toBeDefined()
    expect(AuthContext.Provider).toBeDefined()
    expect(AuthContext.Consumer).toBeDefined()
  })

  it('should have displayName for debugging', () => {
    // Check if context has a display name for better debugging
    expect(AuthContext.displayName || AuthContext._context?.displayName).toBeTruthy()
  })
})