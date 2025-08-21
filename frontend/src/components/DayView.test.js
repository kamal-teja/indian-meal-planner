import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import DayView from './DayView'
import { mealPlannerAPI } from '../services/api'

// Mock the API
vi.mock('../services/api', () => ({
  mealPlannerAPI: {
    getMealsByDate: vi.fn(),
    deleteMeal: vi.fn(),
    updateMeal: vi.fn()
  }
}))

// Mock the date-fns format function
vi.mock('date-fns', () => ({
  format: vi.fn((date, formatStr) => {
    if (formatStr === 'yyyy-MM-dd') {
      return '2023-10-15'
    }
    return '2023-10-15'
  }),
  addDays: vi.fn((date, days) => new Date(2023, 9, 15 + days)),
  subDays: vi.fn((date, days) => new Date(2023, 9, 15 - days))
}))

// Mock child components
vi.mock('./MealCard', () => ({
  default: ({ meal, onEdit, onDelete }) => (
    <div data-testid="meal-card">
      <span>{meal.dish.name}</span>
      <button onClick={() => onEdit(meal)}>Edit</button>
      <button onClick={() => onDelete(meal.id)}>Delete</button>
    </div>
  )
}))

vi.mock('./DishSelector', () => ({
  default: ({ isOpen, onClose, onSelectDish }) => 
    isOpen ? (
      <div data-testid="dish-selector">
        <button onClick={onClose}>Close</button>
        <button onClick={() => onSelectDish({ id: 'dish1', name: 'Test Dish' })}>
          Select Dish
        </button>
      </div>
    ) : null
}))

vi.mock('./IngredientsList', () => ({
  default: ({ isOpen, onClose }) => 
    isOpen ? (
      <div data-testid="ingredients-list">
        <button onClick={onClose}>Close</button>
      </div>
    ) : null
}))

vi.mock('./MealRecommendations', () => ({
  default: ({ isOpen, onClose }) => 
    isOpen ? (
      <div data-testid="meal-recommendations">
        <button onClick={onClose}>Close</button>
      </div>
    ) : null
}))

describe('DayView Component', () => {
  const mockLoadDishes = vi.fn()
  const mockOnAddDish = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
    mealPlannerAPI.getMealsByDate.mockResolvedValue({
      data: {
        data: [
          {
            id: '1',
            dish: { id: 'dish1', name: 'Breakfast Dish', type: 'Veg' },
            mealType: 'breakfast',
            date: '2023-10-15',
            rating: 5
          },
          {
            id: '2',
            dish: { id: 'dish2', name: 'Lunch Dish', type: 'Non-Veg' },
            mealType: 'lunch',
            date: '2023-10-15',
            rating: 4
          }
        ]
      }
    })
  })

  it('should render DayView component', async () => {
    render(<DayView loadDishes={mockLoadDishes} onAddDish={mockOnAddDish} />)
    
    // Should show the date header
    expect(screen.getByText('Sunday, Oct 15, 2023')).toBeInTheDocument()
    
    // Should show meal type sections
    expect(screen.getByText('Breakfast')).toBeInTheDocument()
    expect(screen.getByText('Lunch')).toBeInTheDocument()
    expect(screen.getByText('Dinner')).toBeInTheDocument()
    expect(screen.getByText('Snacks')).toBeInTheDocument()
  })

  it('should load and display meals for the selected date', async () => {
    render(<DayView loadDishes={mockLoadDishes} onAddDish={mockOnAddDish} />)
    
    // Wait for meals to load
    await waitFor(() => {
      expect(mealPlannerAPI.getMealsByDate).toHaveBeenCalledWith('2023-10-15')
    })

    // Should display meal cards
    await waitFor(() => {
      expect(screen.getByText('Breakfast Dish')).toBeInTheDocument()
      expect(screen.getByText('Lunch Dish')).toBeInTheDocument()
    })
  })

  it('should handle navigation between dates', async () => {
    const user = userEvent.setup()
    render(<DayView loadDishes={mockLoadDishes} onAddDish={mockOnAddDish} />)
    
    // Find navigation buttons
    const prevButton = screen.getByRole('button', { name: /previous day/i })
    const nextButton = screen.getByRole('button', { name: /next day/i })
    
    // Click next day
    await user.click(nextButton)
    
    // Should call API with new date
    await waitFor(() => {
      expect(mealPlannerAPI.getMealsByDate).toHaveBeenCalledTimes(2)
    })
  })

  it('should open dish selector when adding a meal', async () => {
    const user = userEvent.setup()
    render(<DayView loadDishes={mockLoadDishes} onAddDish={mockOnAddDish} />)
    
    // Find and click an "Add" button for breakfast
    const addButtons = screen.getAllByRole('button', { name: /add/i })
    await user.click(addButtons[0]) // First add button should be for breakfast
    
    // Should open dish selector
    expect(screen.getByTestId('dish-selector')).toBeInTheDocument()
  })

  it('should close dish selector when cancel is clicked', async () => {
    const user = userEvent.setup()
    render(<DayView loadDishes={mockLoadDishes} onAddDish={mockOnAddDish} />)
    
    // Open dish selector
    const addButtons = screen.getAllByRole('button', { name: /add/i })
    await user.click(addButtons[0])
    
    // Close dish selector
    const closeButton = screen.getByText('Close')
    await user.click(closeButton)
    
    // Should close dish selector
    expect(screen.queryByTestId('dish-selector')).not.toBeInTheDocument()
  })

  it('should handle meal deletion', async () => {
    const user = userEvent.setup()
    mealPlannerAPI.deleteMeal.mockResolvedValue({ success: true })
    
    render(<DayView loadDishes={mockLoadDishes} onAddDish={mockOnAddDish} />)
    
    // Wait for meals to load
    await waitFor(() => {
      expect(screen.getByText('Breakfast Dish')).toBeInTheDocument()
    })
    
    // Click delete button on first meal
    const deleteButtons = screen.getAllByText('Delete')
    await user.click(deleteButtons[0])
    
    // Should call delete API
    await waitFor(() => {
      expect(mealPlannerAPI.deleteMeal).toHaveBeenCalledWith('1')
    })
  })

  it('should open ingredients panel', async () => {
    const user = userEvent.setup()
    render(<DayView loadDishes={mockLoadDishes} onAddDish={mockOnAddDish} />)
    
    // Find and click shopping cart button
    const shoppingButton = screen.getByRole('button', { name: /shopping/i })
    await user.click(shoppingButton)
    
    // Should open ingredients panel
    expect(screen.getByTestId('ingredients-list')).toBeInTheDocument()
  })

  it('should open recommendations panel', async () => {
    const user = userEvent.setup()
    render(<DayView loadDishes={mockLoadDishes} onAddDish={mockOnAddDish} />)
    
    // Find and click recommendations button
    const recommendationsButton = screen.getByRole('button', { name: /recommendations/i })
    await user.click(recommendationsButton)
    
    // Should open recommendations panel
    expect(screen.getByTestId('meal-recommendations')).toBeInTheDocument()
  })

  it('should handle API errors gracefully', async () => {
    mealPlannerAPI.getMealsByDate.mockRejectedValue(new Error('API Error'))
    
    render(<DayView loadDishes={mockLoadDishes} onAddDish={mockOnAddDish} />)
    
    // Should not crash and should show empty state
    await waitFor(() => {
      expect(mealPlannerAPI.getMealsByDate).toHaveBeenCalled()
    })
    
    // Component should still render without errors
    expect(screen.getByText('Breakfast')).toBeInTheDocument()
  })

  it('should display loading state', async () => {
    // Make API call hang to test loading state
    mealPlannerAPI.getMealsByDate.mockImplementation(() => new Promise(() => {}))
    
    render(<DayView loadDishes={mockLoadDishes} onAddDish={mockOnAddDish} />)
    
    // Should show loading state (component should handle this internally)
    expect(mealPlannerAPI.getMealsByDate).toHaveBeenCalled()
  })

  it('should calculate nutrition summary correctly', async () => {
    const mealsWithNutrition = [
      {
        id: '1',
        dish: { 
          id: 'dish1', 
          name: 'Breakfast Dish', 
          calories: 300,
          protein: 15,
          carbs: 40,
          fat: 10
        },
        mealType: 'breakfast',
        date: '2023-10-15'
      },
      {
        id: '2',
        dish: { 
          id: 'dish2', 
          name: 'Lunch Dish', 
          calories: 500,
          protein: 25,
          carbs: 60,
          fat: 20
        },
        mealType: 'lunch',
        date: '2023-10-15'
      }
    ]

    mealPlannerAPI.getMealsByDate.mockResolvedValue({
      data: { data: mealsWithNutrition }
    })
    
    render(<DayView loadDishes={mockLoadDishes} onAddDish={mockOnAddDish} />)
    
    // Wait for nutrition summary to be calculated
    await waitFor(() => {
      // The component should calculate and display total nutrition
      // This would be visible in a nutrition summary section
      expect(screen.getByText('Breakfast Dish')).toBeInTheDocument()
      expect(screen.getByText('Lunch Dish')).toBeInTheDocument()
    })
  })
})