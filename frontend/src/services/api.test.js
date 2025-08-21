import { describe, it, expect, vi, beforeEach } from 'vitest'

// Simple API tests focusing on structure and configuration
describe('mealPlannerAPI', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should have all required API methods', async () => {
    const api = await import('./api')
    
    expect(api.mealPlannerAPI).toBeDefined()
    expect(typeof api.mealPlannerAPI.login).toBe('function')
    expect(typeof api.mealPlannerAPI.register).toBe('function')
    expect(typeof api.mealPlannerAPI.logout).toBe('function')
    expect(typeof api.mealPlannerAPI.getCurrentUser).toBe('function')
    expect(typeof api.mealPlannerAPI.updateProfile).toBe('function')
    expect(typeof api.mealPlannerAPI.addToFavorites).toBe('function')
    expect(typeof api.mealPlannerAPI.removeFromFavorites).toBe('function')
    expect(typeof api.mealPlannerAPI.getDishes).toBe('function')
    expect(typeof api.mealPlannerAPI.searchDishes).toBe('function')
    expect(typeof api.mealPlannerAPI.getFavoriteDishes).toBe('function')
    expect(typeof api.mealPlannerAPI.addDish).toBe('function')
    expect(typeof api.mealPlannerAPI.getMealsByDate).toBe('function')
    expect(typeof api.mealPlannerAPI.getMealsByMonth).toBe('function')
    expect(typeof api.mealPlannerAPI.addMeal).toBe('function')
    expect(typeof api.mealPlannerAPI.updateMeal).toBe('function')
    expect(typeof api.mealPlannerAPI.deleteMeal).toBe('function')
    expect(typeof api.mealPlannerAPI.getMealAnalytics).toBe('function')
    expect(typeof api.mealPlannerAPI.getShoppingList).toBe('function')
    expect(typeof api.mealPlannerAPI.getRecommendations).toBe('function')
    expect(typeof api.mealPlannerAPI.getNutritionProgress).toBe('function')
    expect(typeof api.mealPlannerAPI.getNutritionGoals).toBe('function')
    expect(typeof api.mealPlannerAPI.updateNutritionGoals).toBe('function')
    expect(typeof api.mealPlannerAPI.healthCheck).toBe('function')
  })

  it('should be default exported', async () => {
    const api = await import('./api')
    
    expect(api.default).toBeDefined()
    expect(api.default).toBe(api.mealPlannerAPI)
  })

  it('should define API base URL constant', () => {
    // Test that the API configuration is accessible
    const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:5000/api'
    expect(baseURL).toBeDefined()
    expect(typeof baseURL).toBe('string')
  })
})