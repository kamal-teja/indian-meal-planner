import { describe, it, expect, vi } from 'vitest'
import { extractResponseData, extractArrayData, extractErrorMessage } from './responseHelpers'

describe('responseHelpers', () => {
  describe('extractResponseData', () => {
    it('should extract data from successful response', () => {
      const response = {
        data: {
          success: true,
          data: { id: 1, name: 'test' }
        }
      }

      const result = extractResponseData(response)
      expect(result).toEqual({ id: 1, name: 'test' })
    })

    it('should return null for unsuccessful response', () => {
      const response = {
        data: {
          success: false,
          error: 'Something went wrong'
        }
      }

      const result = extractResponseData(response)
      expect(result).toBeNull()
    })

    it('should return null for malformed response', () => {
      const response = {
        data: null
      }

      const result = extractResponseData(response)
      expect(result).toBeNull()
    })

    it('should return null for undefined response', () => {
      const result = extractResponseData(undefined)
      expect(result).toBeNull()
    })

    it('should warn for unexpected response format', () => {
      const consoleSpy = vi.spyOn(console, 'warn').mockImplementation(() => {})
      
      const response = {
        data: {
          success: false
        }
      }

      extractResponseData(response)
      expect(consoleSpy).toHaveBeenCalledWith('Unexpected response format:', response)
      
      consoleSpy.mockRestore()
    })
  })

  describe('extractArrayData', () => {
    it('should extract array data from successful response', () => {
      const response = {
        data: {
          success: true,
          data: [{ id: 1 }, { id: 2 }]
        }
      }

      const result = extractArrayData(response)
      expect(result).toEqual([{ id: 1 }, { id: 2 }])
    })

    it('should return default value for non-array data', () => {
      const response = {
        data: {
          success: true,
          data: { id: 1, name: 'test' }
        }
      }

      const result = extractArrayData(response)
      expect(result).toEqual([])
    })

    it('should return custom default value for non-array data', () => {
      const response = {
        data: {
          success: true,
          data: null
        }
      }

      const defaultValue = [{ default: true }]
      const result = extractArrayData(response, defaultValue)
      expect(result).toEqual(defaultValue)
    })

    it('should return default value for unsuccessful response', () => {
      const response = {
        data: {
          success: false,
          error: 'Error occurred'
        }
      }

      const result = extractArrayData(response)
      expect(result).toEqual([])
    })

    it('should handle empty array data', () => {
      const response = {
        data: {
          success: true,
          data: []
        }
      }

      const result = extractArrayData(response)
      expect(result).toEqual([])
    })
  })

  describe('extractErrorMessage', () => {
    it('should extract error message from response error', () => {
      const error = {
        response: {
          data: {
            error: 'Validation failed'
          }
        }
      }

      const result = extractErrorMessage(error)
      expect(result).toBe('Validation failed')
    })

    it('should extract error message from error.message', () => {
      const error = {
        message: 'Network error'
      }

      const result = extractErrorMessage(error)
      expect(result).toBe('Network error')
    })

    it('should return default message for unknown error', () => {
      const error = {}

      const result = extractErrorMessage(error)
      expect(result).toBe('An unexpected error occurred')
    })

    it('should prioritize response error over message', () => {
      const error = {
        response: {
          data: {
            error: 'Server error'
          }
        },
        message: 'Network error'
      }

      const result = extractErrorMessage(error)
      expect(result).toBe('Server error')
    })

    it('should handle malformed error objects', () => {
      const error1 = { response: null }
      const error2 = { response: { data: null } }
      const error3 = { response: { data: { error: null } } }

      expect(extractErrorMessage(error1)).toBe('An unexpected error occurred')
      expect(extractErrorMessage(error2)).toBe('An unexpected error occurred')
      expect(extractErrorMessage(error3)).toBe('An unexpected error occurred')
    })

    it('should handle nested error structures', () => {
      const error = {
        response: {
          data: {
            error: 'User not found',
            details: 'Additional info'
          }
        }
      }

      const result = extractErrorMessage(error)
      expect(result).toBe('User not found')
    })
  })
})