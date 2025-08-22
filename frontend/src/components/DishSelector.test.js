import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import DishSelector from './DishSelector'

// Mock the auth context
const mockUser = {
  id: 'user1',
  name: 'Test User',
  favorites: ['dish1', 'dish3']
}

const mockToggleFavorite = vi.fn()

vi.mock('../contexts/AuthContext', () => ({
  useAuth: () => ({
    user: mockUser,
    toggleFavoriteWithState: mockToggleFavorite
  })
}))

// Mock child components
vi.mock('./AddDishForm', () => ({
  default: ({ isOpen, onClose, onSuccess }) => 
    isOpen ? (
      <div data-testid="add-dish-form">
        <button onClick={onClose}>Cancel</button>
        <button onClick={() => onSuccess({ id: 'new-dish', name: 'New Dish' })}>
          Save
        </button>
      </div>
    ) : null
}))

vi.mock('./ui/CustomDropdown', () => ({
  default: ({ value, onChange, options, placeholder }) => (
    <select 
      data-testid="custom-dropdown"
      value={value} 
      onChange={(e) => onChange(e.target.value)}
    >
      <option value="">{placeholder}</option>
      {options.map(option => (
        <option key={option.value} value={option.value}>
          {option.label}
        </option>
      ))}
    </select>
  )
}))

describe('DishSelector Component', () => {
  const mockLoadDishes = vi.fn()
  const mockOnSelect = vi.fn()
  const mockOnClose = vi.fn()
  const mockOnAddDish = vi.fn()

  const mockDishes = [
    {
      id: 'dish1',
      name: 'Paneer Butter Masala',
      type: 'Veg',
      cuisine: 'North Indian',
      spiceLevel: 'Medium',
      cookingTime: 30,
      calories: 350,
      isFavorite: true
    },
    {
      id: 'dish2',
      name: 'Chicken Curry',
      type: 'Non-Veg',
      cuisine: 'South Indian',
      spiceLevel: 'Hot',
      cookingTime: 45,
      calories: 400,
      isFavorite: false
    },
    {
      id: 'dish3',
      name: 'Biryani',
      type: 'Non-Veg',
      cuisine: 'Hyderabadi',
      spiceLevel: 'Medium',
      cookingTime: 60,
      calories: 500,
      isFavorite: true
    }
  ]

  const mockResponse = {
    dishes: mockDishes,
    pagination: {
      hasMore: false,
      totalCount: 3
    }
  }

  beforeEach(() => {
    vi.clearAllMocks()
    mockLoadDishes.mockResolvedValue(mockResponse)
  })

  it('should render DishSelector component', async () => {
    render(
      <DishSelector
        loadDishes={mockLoadDishes}
        onSelect={mockOnSelect}
        onClose={mockOnClose}
        mealType="breakfast"
        onAddDish={mockOnAddDish}
      />
    )
    
    // Should show the modal overlay
    expect(screen.getByText('Select a Dish for Breakfast')).toBeInTheDocument()
    
    // Should show search input
    expect(screen.getByPlaceholderText('Search dishes...')).toBeInTheDocument()
    
    // Should show filter dropdowns
    expect(screen.getByText('All Types')).toBeInTheDocument()
    expect(screen.getByText('All Cuisines')).toBeInTheDocument()
  })

  it('should load and display dishes', async () => {
    render(
      <DishSelector
        loadDishes={mockLoadDishes}
        onSelect={mockOnSelect}
        onClose={mockOnClose}
        mealType="lunch"
        onAddDish={mockOnAddDish}
      />
    )
    
    // Wait for dishes to load
    await waitFor(() => {
      expect(mockLoadDishes).toHaveBeenCalledWith({
        page: 1,
        limit: 20
      })
    })

    // Should display dishes
    await waitFor(() => {
      expect(screen.getByText('Paneer Butter Masala')).toBeInTheDocument()
      expect(screen.getByText('Chicken Curry')).toBeInTheDocument()
      expect(screen.getByText('Biryani')).toBeInTheDocument()
    })
  })

  it('should handle search functionality', async () => {
    const user = userEvent.setup()
    render(
      <DishSelector
        loadDishes={mockLoadDishes}
        onSelect={mockOnSelect}
        onClose={mockOnClose}
        mealType="dinner"
        onAddDish={mockOnAddDish}
      />
    )
    
    const searchInput = screen.getByPlaceholderText('Search dishes...')
    
    // Type in search input
    await user.type(searchInput, 'paneer')
    
    // Should trigger search with debounce
    await waitFor(() => {
      expect(mockLoadDishes).toHaveBeenCalledWith({
        page: 1,
        limit: 20,
        search: 'paneer'
      })
    }, { timeout: 2000 })
  })

  it('should handle type filter', async () => {
    const user = userEvent.setup()
    render(
      <DishSelector
        loadDishes={mockLoadDishes}
        onSelect={mockOnSelect}
        onClose={mockOnClose}
        mealType="snack"
        onAddDish={mockOnAddDish}
      />
    )
    
    // Find and change type filter
    const typeDropdown = screen.getAllByTestId('custom-dropdown')[0]
    await user.selectOptions(typeDropdown, 'Veg')
    
    // Should filter by vegetarian dishes
    await waitFor(() => {
      expect(mockLoadDishes).toHaveBeenCalledWith({
        page: 1,
        limit: 20,
        type: 'Veg'
      })
    })
  })

  it('should handle cuisine filter', async () => {
    const user = userEvent.setup()
    render(
      <DishSelector
        loadDishes={mockLoadDishes}
        onSelect={mockOnSelect}
        onClose={mockOnClose}
        mealType="breakfast"
        onAddDish={mockOnAddDish}
      />
    )
    
    // Find and change cuisine filter
    const cuisineDropdown = screen.getAllByTestId('custom-dropdown')[1]
    await user.selectOptions(cuisineDropdown, 'North Indian')
    
    // Should filter by cuisine
    await waitFor(() => {
      expect(mockLoadDishes).toHaveBeenCalledWith({
        page: 1,
        limit: 20,
        cuisine: 'North Indian'
      })
    })
  })

  it('should handle dish selection', async () => {
    const user = userEvent.setup()
    render(
      <DishSelector
        loadDishes={mockLoadDishes}
        onSelect={mockOnSelect}
        onClose={mockOnClose}
        mealType="lunch"
        onAddDish={mockOnAddDish}
      />
    )
    
    // Wait for dishes to load
    await waitFor(() => {
      expect(screen.getByText('Paneer Butter Masala')).toBeInTheDocument()
    })
    
    // Click on a dish
    await user.click(screen.getByText('Paneer Butter Masala'))
    
    // Should call onSelect with the dish
    expect(mockOnSelect).toHaveBeenCalledWith(mockDishes[0])
  })

  it('should handle favorite toggle', async () => {
    const user = userEvent.setup()
    render(
      <DishSelector
        loadDishes={mockLoadDishes}
        onSelect={mockOnSelect}
        onClose={mockOnClose}
        mealType="dinner"
        onAddDish={mockOnAddDish}
      />
    )
    
    // Wait for dishes to load
    await waitFor(() => {
      expect(screen.getByText('Chicken Curry')).toBeInTheDocument()
    })
    
    // Find and click favorite button for a non-favorite dish
    const favoriteButtons = screen.getAllByRole('button', { name: /favorite/i })
    await user.click(favoriteButtons[1]) // Chicken Curry favorite button
    
    // Should call toggle favorite
    expect(mockToggleFavorite).toHaveBeenCalledWith('dish2')
  })

  it('should close modal when close button is clicked', async () => {
    const user = userEvent.setup()
    render(
      <DishSelector
        loadDishes={mockLoadDishes}
        onSelect={mockOnSelect}
        onClose={mockOnClose}
        mealType="breakfast"
        onAddDish={mockOnAddDish}
      />
    )
    
    // Find and click close button
    const closeButton = screen.getByRole('button', { name: /close/i })
    await user.click(closeButton)
    
    // Should call onClose
    expect(mockOnClose).toHaveBeenCalled()
  })

  it('should open add dish form', async () => {
    const user = userEvent.setup()
    render(
      <DishSelector
        loadDishes={mockLoadDishes}
        onSelect={mockOnSelect}
        onClose={mockOnClose}
        mealType="lunch"
        onAddDish={mockOnAddDish}
      />
    )
    
    // Find and click "Add New Dish" button
    const addButton = screen.getByText('Add New Dish')
    await user.click(addButton)
    
    // Should open add dish form
    expect(screen.getByTestId('add-dish-form')).toBeInTheDocument()
  })

  it('should handle adding new dish', async () => {
    const user = userEvent.setup()
    render(
      <DishSelector
        loadDishes={mockLoadDishes}
        onSelect={mockOnSelect}
        onClose={mockOnClose}
        mealType="dinner"
        onAddDish={mockOnAddDish}
      />
    )
    
    // Open add dish form
    await user.click(screen.getByText('Add New Dish'))
    
    // Save new dish
    await user.click(screen.getByText('Save'))
    
    // Should call onAddDish and close form
    expect(mockOnAddDish).toHaveBeenCalledWith({ id: 'new-dish', name: 'New Dish' })
    expect(screen.queryByTestId('add-dish-form')).not.toBeInTheDocument()
  })

  it('should display dish information correctly', async () => {
    render(
      <DishSelector
        loadDishes={mockLoadDishes}
        onSelect={mockOnSelect}
        onClose={mockOnClose}
        mealType="breakfast"
        onAddDish={mockOnAddDish}
      />
    )
    
    // Wait for dishes to load
    await waitFor(() => {
      expect(screen.getByText('Paneer Butter Masala')).toBeInTheDocument()
    })
    
    // Should display dish details
    expect(screen.getAllByText('North Indian')).toHaveLength(2) // One in filter, one in dish card
    expect(screen.getByText('350 cal')).toBeInTheDocument()
    expect(screen.getAllByText('30 min')).toHaveLength(3) // Multiple dishes may have same time
  })

  it('should handle loading state', async () => {
    mockLoadDishes.mockImplementation(() => new Promise(() => {})) // Never resolves
    
    render(
      <DishSelector
        loadDishes={mockLoadDishes}
        onSelect={mockOnSelect}
        onClose={mockOnClose}
        mealType="lunch"
        onAddDish={mockOnAddDish}
      />
    )
    
    // Should show loading indicator
    await waitFor(() => {
      expect(mockLoadDishes).toHaveBeenCalled()
    })
  })

  it('should handle empty dishes list', async () => {
    mockLoadDishes.mockResolvedValue({
      dishes: [],
      pagination: { hasMore: false, totalCount: 0 }
    })
    
    render(
      <DishSelector
        loadDishes={mockLoadDishes}
        onSelect={mockOnSelect}
        onClose={mockOnClose}
        mealType="snack"
        onAddDish={mockOnAddDish}
      />
    )
    
    // Should show empty state message
    await waitFor(() => {
      expect(screen.getByText(/no dishes found/i)).toBeInTheDocument()
    })
  })
})