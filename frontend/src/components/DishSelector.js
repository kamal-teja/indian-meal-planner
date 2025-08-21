import React, { useState, useEffect, useCallback, useRef } from 'react';
import { X, Search, Filter, Flame, Clock, Plus, Heart, Loader, Sun, Sunset, Moon, Coffee, Check } from 'lucide-react';
import AddDishForm from './AddDishForm';
import CustomDropdown from './ui/CustomDropdown';
import { useAuth } from '../contexts/AuthContext';

const DishSelector = ({ loadDishes, onSelect, onClose, mealType, isEditing, onAddDish }) => {
  const { user, toggleFavoriteWithState } = useAuth();
  const [dishes, setDishes] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedType, setSelectedType] = useState('all');
  const [selectedCuisine, setSelectedCuisine] = useState('all');
  const [showAddDishForm, setShowAddDishForm] = useState(false);
  const [dishFavorites, setDishFavorites] = useState({});
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const [cuisines, setCuisines] = useState([]);
  const [selectedDishId, setSelectedDishId] = useState(null);
  const [isSelecting, setIsSelecting] = useState(false);
  const [isSuccess, setIsSuccess] = useState(false);
  const isInitialLoadRef = useRef(true);
  
  const observer = useRef();
  
  const loadMoreDishes = useCallback(async () => {
    if (loading || !hasMore) return;
    
    try {
      setLoading(true);
      const nextPage = currentPage + 1;
      
      const params = {
        page: nextPage,
        limit: 20,
        ...(searchTerm && { search: searchTerm }),
        ...(selectedType !== 'all' && { type: selectedType }),
        ...(selectedCuisine !== 'all' && { cuisine: selectedCuisine })
      };
      
      const response = await loadDishes(params);
      setDishes(prev => [...prev, ...(response.dishes || [])]);
      setCurrentPage(nextPage);
      setHasMore(response.pagination?.hasMore || false);
    } catch (error) {
      console.error('Error loading more dishes:', error);
    } finally {
      setLoading(false);
    }
  }, [loadDishes, loading, hasMore, currentPage, searchTerm, selectedType, selectedCuisine]);

  const lastDishElementRef = useCallback(node => {
    if (loading) return;
    if (observer.current) observer.current.disconnect();
    observer.current = new IntersectionObserver(entries => {
      if (entries[0].isIntersecting && hasMore) {
        loadMoreDishes();
      }
    });
    if (node) observer.current.observe(node);
  }, [loading, hasMore, loadMoreDishes]);

  const loadInitialDishes = useCallback(async () => {
    try {
      setLoading(true);
      const params = {
        page: 1,
        limit: 20
      };
      
      const response = await loadDishes(params);
      setDishes(response.dishes || []);
      setCurrentPage(1);
      setHasMore(response.pagination?.hasMore || false);
      setTotalCount(response.pagination?.totalCount || 0);
      
      // Extract unique cuisines for filter dropdown
      const uniqueCuisines = [...new Set((response.dishes || []).map(dish => dish.cuisine))];
      setCuisines(uniqueCuisines.sort());
      
      // Mark initial load as complete
      isInitialLoadRef.current = false;
    } catch (error) {
      console.error('Error loading dishes:', error);
    } finally {
      setLoading(false);
    }
  }, [loadDishes]);

  const resetAndLoadDishes = useCallback(async () => {
    try {
      setLoading(true);
      setDishes([]);
      setCurrentPage(1);
      
      const params = {
        page: 1,
        limit: 20,
        ...(searchTerm && { search: searchTerm }),
        ...(selectedType !== 'all' && { type: selectedType }),
        ...(selectedCuisine !== 'all' && { cuisine: selectedCuisine })
      };
      
      const response = await loadDishes(params);
      setDishes(response.dishes || []);
      setHasMore(response.pagination?.hasMore || false);
      setTotalCount(response.pagination?.totalCount || 0);
    } catch (error) {
      console.error('Error loading dishes:', error);
    } finally {
      setLoading(false);
    }
  }, [loadDishes, searchTerm, selectedType, selectedCuisine]); // Restore dependencies

  useEffect(() => {
    // Load initial dishes when component mounts
    loadInitialDishes();
  }, []); // Only run once on mount

  useEffect(() => {
    // Skip on initial load - loadInitialDishes will handle it
    if (isInitialLoadRef.current) return;
    
    // Reset and reload dishes when filters change (debounce if needed)
    const timer = setTimeout(() => {
      resetAndLoadDishes();
    }, 300);
    
    return () => clearTimeout(timer);
  }, [searchTerm, selectedType, selectedCuisine, resetAndLoadDishes]);

  const getMealTypeIcon = (mealType) => {
    const iconMap = {
      breakfast: <Coffee className="h-6 w-6 text-white/80" />,
      lunch: <Sun className="h-6 w-6 text-white/80" />,
      dinner: <Moon className="h-6 w-6 text-white/80" />,
      snack: <Sunset className="h-6 w-6 text-white/80" />
    };
    return iconMap[mealType] || <Coffee className="h-6 w-6 text-white/80" />;
  };

  const getTypeColor = (type) => {
    return type === 'Veg' 
      ? 'badge-veg' 
      : 'badge-non-veg';
  };

  const getCuisineColor = (cuisine) => {
    const colors = {
      'North Indian': 'bg-warm-100 text-warm-800 border border-warm-200',
      'South Indian': 'bg-sage-100 text-sage-800 border border-sage-200',
      'Punjabi': 'bg-secondary-100 text-secondary-800 border border-secondary-200',
      'Hyderabadi': 'bg-lavender-100 text-lavender-800 border border-lavender-200',
      'Bengali': 'bg-accent-100 text-accent-800 border border-accent-200',
      'Gujarati': 'bg-neutral-100 text-neutral-800 border border-neutral-200',
    };
    return colors[cuisine] || 'bg-neutral-100 text-neutral-800 border border-neutral-200';
  };

  const handleAddDish = async (dishData) => {
    try {
      const newDish = await onAddDish(dishData);
      setShowAddDishForm(false);
      // Add the new dish to the beginning of the list
      if (newDish) {
        setDishes(prev => [newDish, ...prev]);
        onSelect(newDish);
      }
    } catch (error) {
      console.error('Error adding dish:', error);
      throw error;
    }
  };

  const handleDishSelect = async (dish) => {
    try {
      setSelectedDishId(dish.id);
      setIsSelecting(true);
      
      // Add a small delay for visual feedback
      await new Promise(resolve => setTimeout(resolve, 400));
      
      // Call the parent onSelect
      await onSelect(dish);
      
      // Show success state
      setIsSelecting(false);
      setIsSuccess(true);
      
      // Show success state briefly before closing
      await new Promise(resolve => setTimeout(resolve, 500));
      
    } catch (error) {
      console.error('Error selecting dish:', error);
    } finally {
      setIsSelecting(false);
      setSelectedDishId(null);
      setIsSuccess(false);
    }
  };

  const handleToggleFavorite = async (e, dishId) => {
    e.stopPropagation();
    if (!user) return;

    const currentIsFavorite = isDishFavorite({ id: dishId });
    const result = await toggleFavoriteWithState(dishId, currentIsFavorite);
    if (result.success) {
      setDishFavorites(prev => ({
        ...prev,
        [dishId]: result.isFavorite
      }));
    } else {
      console.error('Failed to toggle favorite:', result.error);
    }
  };

  const isDishFavorite = (dish) => {
    if (dishFavorites[dish.id] !== undefined) {
      return dishFavorites[dish.id];
    }
    return dish.isFavorite || false;
  };

  return (
    <div className="fixed inset-0 bg-black/40 backdrop-blur-sm flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-lg shadow-xl max-w-4xl w-full overflow-hidden border border-neutral-200 flex flex-col" style={{ height: '80vh', minHeight: '600px' }}>
        {/* Header */}
        <div className="bg-secondary-600 p-6 text-white flex-shrink-0">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <span className="flex items-center justify-center">{getMealTypeIcon(mealType)}</span>
              <div>
                <h2 className="text-2xl font-display font-bold">
                  {isEditing ? 'Change' : 'Select'} {mealType.charAt(0).toUpperCase() + mealType.slice(1)}
                </h2>
                <p className="text-white/80">Choose from our delicious Indian dishes</p>
              </div>
            </div>
            <div className="flex items-center space-x-2">
              <button
                onClick={() => setShowAddDishForm(true)}
                className="flex items-center space-x-2 px-4 py-2 bg-white/20 hover:bg-white/30 text-white rounded-lg transition-colors"
              >
                <Plus className="w-4 h-4" />
                <span>Add New Dish</span>
              </button>
              <button
                onClick={onClose}
                className="text-white/70 hover:text-white hover:bg-white/10 p-2 rounded-lg transition-all duration-200"
              >
                <X className="h-6 w-6" />
              </button>
            </div>
          </div>
        </div>

        {/* Filters */}
        <div className="p-6 border-b border-neutral-200 bg-neutral-50 flex-shrink-0">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {/* Search */}
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-neutral-400" />
              <input
                type="text"
                placeholder="Search dishes..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full pl-10 pr-4 py-3 border border-neutral-300 rounded-lg focus:ring-2 focus:ring-accent-500 focus:border-accent-500"
              />
            </div>

            {/* Type Filter */}
            <CustomDropdown
              value={selectedType}
              onChange={setSelectedType}
              options={[
                { value: 'all', label: 'All Types' },
                { value: 'Veg', label: 'Vegetarian' },
                { value: 'Non-Veg', label: 'Non-Vegetarian' }
              ]}
              placeholder="Select Type"
              className="w-full"
            />

            {/* Cuisine Filter */}
            <CustomDropdown
              value={selectedCuisine}
              onChange={setSelectedCuisine}
              options={[
                { value: 'all', label: 'All Cuisines' },
                ...cuisines.map(cuisine => ({ value: cuisine, label: cuisine }))
              ]}
              placeholder="Select Cuisine"
              className="w-full"
            />
          </div>

          <div className="flex items-center justify-between mt-4">
            <span className="text-sm text-neutral-600">
              {totalCount > 0 ? `${dishes.length} of ${totalCount} dish${totalCount !== 1 ? 'es' : ''}` : 'No dishes found'}
            </span>
            <button
              onClick={() => {
                setSearchTerm('');
                setSelectedType('all');
                setSelectedCuisine('all');
              }}
              className="text-sm text-accent-600 hover:text-accent-700 font-medium"
            >
              Clear filters
            </button>
          </div>
        </div>

        {/* Dishes Grid */}
        <div className="flex-1 p-6 overflow-y-auto">
          {loading && dishes.length === 0 ? (
            /* Initial Loading Skeleton */
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {[...Array(6)].map((_, index) => (
                <div key={index} className="p-4 border border-neutral-200 rounded-lg bg-white">
                  <div className="flex items-start space-x-4">
                    <div className="relative flex-shrink-0">
                      <div className="w-20 h-20 bg-neutral-200 rounded-lg animate-pulse"></div>
                      <div className="absolute -top-2 -right-2 w-12 h-6 bg-neutral-200 rounded-full animate-pulse"></div>
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="h-5 bg-neutral-200 rounded animate-pulse mb-2" style={{ width: `${75 + (index * 5)}%` }}></div>
                      <div className="h-4 bg-neutral-200 rounded animate-pulse mb-3 w-1/2"></div>
                      <div className="flex items-center space-x-3 mb-2">
                        <div className="h-4 bg-neutral-200 rounded animate-pulse w-16"></div>
                        <div className="h-4 bg-neutral-200 rounded animate-pulse w-16"></div>
                      </div>
                      <div className="flex space-x-1">
                        <div className="h-6 bg-neutral-200 rounded animate-pulse w-12"></div>
                        <div className="h-6 bg-neutral-200 rounded animate-pulse w-12"></div>
                        <div className="h-6 bg-neutral-200 rounded animate-pulse w-8"></div>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          ) : dishes.length === 0 && !loading ? (
            <div className="text-center py-12">
              <Filter className="h-16 w-16 text-neutral-300 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-neutral-900 mb-2">No dishes found</h3>
              <p className="text-neutral-600">Try adjusting your search or filters</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {dishes.map((dish, index) => {
                const isSelected = selectedDishId === dish.id;
                const isCurrentlySelecting = isSelecting && isSelected;
                const isSuccessState = isSuccess && isSelected;
                
                return (
                <div
                  key={dish.id}
                  ref={index === dishes.length - 1 ? lastDishElementRef : null}
                  onClick={() => !isSelecting && !isSuccess && handleDishSelect(dish)}
                  className={`
                    p-4 border rounded-lg cursor-pointer transition-all duration-300 group bg-white relative overflow-hidden
                    ${(isSelected && isSelecting) || isSuccessState
                      ? 'border-accent-500 shadow-lg scale-105 bg-accent-50' 
                      : 'border-neutral-200 hover:border-accent-300 hover:shadow-md'
                    }
                    ${(isSelecting && !isSelected) || isSuccess ? 'opacity-50 pointer-events-none' : ''}
                  `}
                >
                  <div className="flex items-start space-x-4">
                    <div className="relative flex-shrink-0">
                      <img
                        src={dish.image}
                        alt={dish.name}
                        className="w-20 h-20 object-cover rounded-lg group-hover:scale-105 transition-transform"
                        onError={(e) => {
                          e.target.src = `https://images.unsplash.com/photo-1546833999-b9f581a1996d?w=300&h=200&fit=crop`;
                        }}
                      />
                      <div className="absolute -top-2 -right-2">
                        <span className={`px-2 py-1 text-xs font-medium rounded-full ${getTypeColor(dish.type)}`}>
                          {dish.type}
                        </span>
                      </div>
                      {user && (
                        <button
                          onClick={(e) => handleToggleFavorite(e, dish.id)}
                          className={`absolute -top-2 -left-2 p-1.5 bg-white rounded-full shadow-md transition-colors ${
                            isDishFavorite(dish)
                              ? 'text-red-500 hover:text-red-600'
                              : 'text-neutral-400 hover:text-red-500'
                          }`}
                        >
                          <Heart className={`h-4 w-4 ${isDishFavorite(dish) ? 'fill-current' : ''}`} />
                        </button>
                      )}
                    </div>

                    <div className="flex-1 min-w-0">
                      <h3 className="font-semibold text-neutral-900 group-hover:text-accent-600 transition-colors">
                        {dish.name}
                      </h3>
                      <span className={`inline-block px-2 py-1 text-xs font-medium rounded-full mt-1 ${getCuisineColor(dish.cuisine)}`}>
                        {dish.cuisine}
                      </span>

                      <div className="flex items-center space-x-3 mt-2 text-sm text-neutral-600">
                        <div className="flex items-center space-x-1">
                          <Flame className="h-4 w-4 text-orange-500" />
                          <span className="font-medium">{dish.calories} cal</span>
                        </div>
                        <div className="flex items-center space-x-1">
                          <Clock className="h-4 w-4 text-blue-500" />
                          <span>30 min</span>
                        </div>
                      </div>

                      {dish.ingredients && dish.ingredients.length > 0 && (
                        <div className="mt-2">
                          <div className="flex flex-wrap gap-1">
                            {dish.ingredients.slice(0, 2).map((ingredient, index) => (
                              <span
                                key={index}
                                className="px-1.5 py-0.5 bg-neutral-100 text-neutral-600 text-xs rounded"
                              >
                                {ingredient}
                              </span>
                            ))}
                            {dish.ingredients.length > 2 && (
                              <span className="px-1.5 py-0.5 bg-neutral-100 text-neutral-500 text-xs rounded">
                                +{dish.ingredients.length - 2}
                              </span>
                            )}
                          </div>
                        </div>
                      )}
                    </div>
                  </div>
                  
                  {/* Selection Loading/Success Overlay */}
                  {(isCurrentlySelecting || isSuccessState) && (
                    <div className="absolute inset-0 bg-accent-500/10 flex items-center justify-center rounded-lg">
                      <div className="bg-white rounded-full p-3 shadow-lg">
                        {isCurrentlySelecting ? (
                          <div className="flex items-center space-x-2">
                            <Loader className="h-5 w-5 animate-spin text-accent-600" />
                            <span className="text-sm font-medium text-accent-700">
                              Adding to {mealType}...
                            </span>
                          </div>
                        ) : (
                          <div className="flex items-center space-x-2">
                            <div className="h-5 w-5 bg-green-500 rounded-full flex items-center justify-center">
                              <Check className="h-3 w-3 text-white" />
                            </div>
                            <span className="text-sm font-medium text-green-700">
                              Added successfully!
                            </span>
                          </div>
                        )}
                      </div>
                    </div>
                  )}
                </div>
                );
              })}
            </div>
          )}
          
          {/* Loading indicator */}
          {loading && (
            <div className="flex justify-center items-center py-8">
              <Loader className="h-6 w-6 animate-spin text-accent-600" />
              <span className="ml-2 text-neutral-600">Loading more dishes...</span>
            </div>
          )}
          
          {/* End of results indicator */}
          {!hasMore && dishes.length > 0 && (
            <div className="text-center py-4 text-gray-500 text-sm">
              You've reached the end of the list
            </div>
          )}
        </div>
      </div>
      
      {/* Add Dish Form Modal */}
      {showAddDishForm && (
        <AddDishForm
          onSubmit={handleAddDish}
          onClose={() => setShowAddDishForm(false)}
        />
      )}
    </div>
  );
};

export default DishSelector;
