import React, { useState, useEffect, useRef } from 'react';
import { format, addDays, subDays } from 'date-fns';
import { ChevronLeft, ChevronRight, Plus, X, Trash, ShoppingCart, Coffee, Sun, Moon, Sunset, Sparkles, MoreHorizontal } from 'lucide-react';
import { mealPlannerAPI } from '../services/api';
import MealCard from './MealCard';
import DishSelector from './DishSelector';
import IngredientsList from './IngredientsList';
import MealRecommendations from './MealRecommendations';
import ConfirmDialog from './ConfirmDialog';
import UndoSnackbar from './UndoSnackbar';
import { useNotification } from '../contexts/NotificationContext';

const MEAL_TYPES = [
  { id: 'breakfast', name: 'Breakfast', icon: <Coffee className="h-5 w-5" />, color: 'bg-warm-100 border-warm-200', bgColor: 'bg-warm-50', textColor: 'text-warm-700' },
  { id: 'lunch', name: 'Lunch', icon: <Sun className="h-5 w-5" />, color: 'bg-sage-100 border-sage-200', bgColor: 'bg-sage-50', textColor: 'text-sage-700' },
  { id: 'dinner', name: 'Dinner', icon: <Moon className="h-5 w-5" />, color: 'bg-lavender-100 border-lavender-200', bgColor: 'bg-lavender-50', textColor: 'text-lavender-700' },
  { id: 'snack', name: 'Snacks', icon: <Sunset className="h-5 w-5" />, color: 'bg-secondary-100 border-secondary-200', bgColor: 'bg-secondary-50', textColor: 'text-secondary-700' }
];

const DayView = ({ loadDishes, onAddDish }) => {
  const [selectedDate, setSelectedDate] = useState(new Date());
  const [meals, setMeals] = useState([]);
  const [loading, setLoading] = useState(false);
  const [showDishSelector, setShowDishSelector] = useState(false);
  const [selectedMealType, setSelectedMealType] = useState(null);
  const [editingMeal, setEditingMeal] = useState(null);
  const [openMenuKey, setOpenMenuKey] = useState(null);
  const [hoverMenuKey, setHoverMenuKey] = useState(null);
  const [confirmState, setConfirmState] = useState({ open: false, title: '', description: '', onConfirm: null });
  const hideTimerRef = useRef(null);
  const menuRef = useRef(null);
  const firstMenuItemRef = useRef(null);
  const [animatingGroupKey, setAnimatingGroupKey] = useState(null);
  const mountedRef = useRef(false);
  // hover-delete and quick-delete removed: deletions are explicit via the menu
  const [snackbar, setSnackbar] = useState({ open: false, message: '', date: null, dishId: null });
  const { notify } = useNotification();
  const [showIngredientsPanel, setShowIngredientsPanel] = useState(false);
  const [showRecommendations, setShowRecommendations] = useState(false);
  const newAddedDishIdRef = useRef(null);
  const recommendationsRef = useRef(null);
  const [highlightRecommendations, setHighlightRecommendations] = useState(false);

  useEffect(() => {
    loadMealsForDate();
  }, [selectedDate]);

  // Click-away listener + cleanup for hide timer
  useEffect(() => {
    const handleDocMouseDown = (e) => {
      if (openMenuKey && menuRef.current && !menuRef.current.contains(e.target)) {
        setOpenMenuKey(null);
        setHoverMenuKey(null);
      }
    };

    document.addEventListener('mousedown', handleDocMouseDown);
    return () => {
      document.removeEventListener('mousedown', handleDocMouseDown);
      if (hideTimerRef.current) {
        clearTimeout(hideTimerRef.current);
        hideTimerRef.current = null;
      }
    };
  }, [openMenuKey]);

  // Focus the first menu item when a menu opens
  useEffect(() => {
    if (openMenuKey) {
      // small timeout to ensure element is mounted
      setTimeout(() => {
        try { firstMenuItemRef.current?.focus(); } catch (e) { /* ignore */ }
      }, 0);
    }
  }, [openMenuKey]);

  // When recommendations are shown, scroll them into view so user sees them
  useEffect(() => {
    if (showRecommendations && recommendationsRef.current) {
      // small timeout to ensure layout is settled
      setTimeout(() => {
        try {
          recommendationsRef.current.scrollIntoView({ behavior: 'smooth', block: 'start' });
          // briefly highlight the recommendations block
          setHighlightRecommendations(true);
          setTimeout(() => setHighlightRecommendations(false), 900);
        } catch (e) {
          // ignore
        }
      }, 80);
    }
  }, [showRecommendations]);

  const loadMealsForDate = async () => {
    try {
      setLoading(true);
      const dateStr = format(selectedDate, 'yyyy-MM-dd');
      const response = await mealPlannerAPI.getMealsByDate(dateStr);
      // Backend returns { success: true, data: [meals] }, so we need response.data.data
      setMeals(response.data.data || []);
    } catch (error) {
      console.error('Error loading meals:', error);
      setMeals([]); // Set empty array on error to prevent reduce error
    } finally {
      setLoading(false);
    }
  };

  const handleAddMeal = (mealType) => {
    setSelectedMealType(mealType);
    setEditingMeal(null);
    setShowDishSelector(true);
  };

  const handleEditMeal = (meal) => {
    setSelectedMealType(meal.mealType);
    setEditingMeal(meal);
    setShowDishSelector(true);
  };

  const handleDeleteOneMeal = async (meal) => {
    if (!meal) return;
    try {
      await mealPlannerAPI.deleteMeal(meal.id);
      // trigger small count-decrease animation for this group
      const groupKey = meal?.dish?.id || meal?.dishId || `meal-${meal.id}`;
      setAnimatingGroupKey(groupKey);
      setTimeout(() => setAnimatingGroupKey(null), 450);
      await loadMealsForDate();
      setOpenMenuKey(null);
    } catch (error) {
      console.error('Error deleting meal:', error);
    }
  };

  const handleDeleteAllMeals = async (groupKey) => {
    if (!groupKey) return;
    try {
      // If groupKey is an actual dishId (not a fallback like `meal-<id>`), use date+dish API for bulk delete
      const isDishId = !groupKey.startsWith('meal-');
      if (isDishId) {
        const dateStr = format(selectedDate, 'yyyy-MM-dd');
        await mealPlannerAPI.deleteMealsByDateAndDish(dateStr, groupKey);
      } else {
        // fallback to deleting by individual IDs
        const toDelete = meals.filter(m => {
          const dishId = m?.dish?.id || m?.dishId || `meal-${m.id}`;
          return dishId === groupKey;
        });
        const ids = toDelete.map(m => m.id);
        const resp = await mealPlannerAPI.deleteMealsBulk(ids);
        // backend returns { success: true, undoToken }
        if (resp?.data?.undoToken) {
          const token = resp.data.undoToken;
          notify(`${toDelete.length} meal(s) deleted`, { variant: 'success', action: { label: 'Undo', onClick: async () => { await mealPlannerAPI.undoByToken(token); await loadMealsForDate(); } } });
        }
      }
      setAnimatingGroupKey(groupKey);
      setTimeout(() => setAnimatingGroupKey(null), 450);
  await loadMealsForDate();
      setOpenMenuKey(null);
    } catch (error) {
      console.error('Error deleting meals:', error);
    }
  };

  const addMeal = async (mealData) => {
    try {
      const resp = await mealPlannerAPI.addMeal(mealData);
      // store dish id so we can scroll to it after reload
      if (mealData?.dishId) newAddedDishIdRef.current = mealData.dishId;
      await loadMealsForDate(); // Refresh meals after adding

      // scroll to the newly added dish group if present
      if (newAddedDishIdRef.current) {
        setTimeout(() => {
          try {
            const el = document.querySelector(`[data-dish-id="${newAddedDishIdRef.current}"]`);
            if (el) el.scrollIntoView({ behavior: 'smooth', block: 'center' });
          } catch (e) {
            // ignore
          }
          newAddedDishIdRef.current = null;
        }, 150);
      }
    } catch (error) {
      console.error('Error adding meal:', error);
    }
  };

  const handleDishSelect = async (dish) => {
    try {
      const dateStr = format(selectedDate, 'yyyy-MM-dd');
      
      if (editingMeal) {
        // Update existing meal
        await mealPlannerAPI.updateMeal(editingMeal.id, dish.id);
      } else {
        // Add new meal
        await mealPlannerAPI.addMeal({
          date: dateStr,
          mealType: selectedMealType,
          dishId: dish.id
        });
      }
  await loadMealsForDate();
      setShowDishSelector(false);
      setEditingMeal(null);
    } catch (error) {
      console.error('Error saving meal:', error);
    }
  };

  const handleDeleteMeal = async (mealId) => {
    // Confirmation handled by ConfirmDialog in menu flow; preserve API call for other uses
    try {
      await mealPlannerAPI.deleteMeal(mealId);
      await loadMealsForDate();
    } catch (error) {
      console.error('Error deleting meal:', error);
    }
  };

  const navigateDate = (direction) => {
    if (direction === 'prev') {
      setSelectedDate(subDays(selectedDate, 1));
    } else {
      setSelectedDate(addDays(selectedDate, 1));
    }
  };

  const getMealsForType = (mealType) => {
    // Safety check: ensure meals is an array before calling filter
    if (!Array.isArray(meals)) {
      console.warn('meals is not an array:', meals);
      return [];
    }
    return meals.filter(meal => meal.mealType === mealType);
  };

  // Quick-delete keyboard shortcut removed for explicit UX.

  const getTotalCalories = () => {
    // Safety check: ensure meals is an array before calling reduce
    if (!Array.isArray(meals)) {
      console.warn('meals is not an array:', meals);
      return 0;
    }
    return meals.reduce((total, meal) => total + (meal.dish?.calories || 0), 0);
  };

  return (
    <div className="max-w-7xl mx-auto p-6">
      <div className={`grid transition-all duration-300 gap-6 ${showIngredientsPanel ? 'grid-cols-1 lg:grid-cols-3' : 'grid-cols-1'}`}>
        {/* Main Content */}
        <div className={`space-y-6 ${showIngredientsPanel ? 'lg:col-span-2' : ''}`}>
          {/* Date Navigation */}
          <div className="card-elevated p-6">
            <div className="flex items-center justify-between">
              <button
                onClick={() => navigateDate('prev')}
                className="p-3 rounded-lg hover:bg-accent-50 transition-all duration-200 hover:shadow-md"
              >
                <ChevronLeft className="h-6 w-6 text-accent-600" />
              </button>
              
              <div className="text-center">
                <h2 className="text-3xl font-display font-bold text-neutral-800">
                  {format(selectedDate, 'EEEE')}
                </h2>
                <p className="text-lg text-accent-600">
                  {format(selectedDate, 'MMMM d, yyyy')}
                </p>
                <div className="mt-3 px-6 py-2 bg-accent-100 rounded-full inline-block border border-accent-200">
                  <span className="text-sm font-medium text-accent-700">
                    {getTotalCalories()} calories today
                  </span>
                </div>
              </div>
              
              <div className="flex items-center space-x-2">
                <button
                  onClick={() => setShowRecommendations(!showRecommendations)}
                  className={`px-4 py-3 rounded-lg transition-all duration-200 flex items-center space-x-2 font-medium ${
                    showRecommendations 
                      ? 'bg-lavender-100 text-lavender-700 hover:bg-lavender-200 shadow-md border border-lavender-200' 
                      : 'bg-white hover:bg-accent-50 text-accent-600 hover:shadow-md border border-accent-200 shadow-sm'
                  }`}
                  title="Toggle Meal Recommendations"
                >
                  <Sparkles className="h-4 w-4" />
                  <span className="hidden sm:inline text-sm">Recommendations</span>
                </button>
                <button
                  onClick={() => setShowIngredientsPanel(!showIngredientsPanel)}
                  className={`px-4 py-3 rounded-lg transition-all duration-200 flex items-center space-x-2 font-medium relative ${
                    showIngredientsPanel 
                      ? 'bg-sage-100 text-sage-700 hover:bg-sage-200 shadow-md border border-sage-200' 
                      : 'bg-white hover:bg-accent-50 text-accent-600 hover:shadow-md border border-accent-200 shadow-sm'
                  }`}
                  title="Toggle Groceries"
                >
                  <ShoppingCart className="h-4 w-4" />
                  <span className="hidden sm:inline text-sm">Groceries</span>
                  {meals.length > 0 && (
                    <span className="absolute -top-1 -right-1 bg-accent-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center shadow-sm">
                      {meals.reduce((count, meal) => count + (meal.dish?.ingredients?.length || 0), 0)}
                    </span>
                  )}
                </button>
                <button
                  onClick={() => navigateDate('next')}
                  className="p-3 rounded-lg hover:bg-accent-50 transition-all duration-200 hover:shadow-md"
                >
                  <ChevronRight className="h-6 w-6 text-accent-600" />
                </button>
              </div>
            </div>
          </div>

          {/* Meal Sections */}
          {loading ? (
            <div className="flex justify-center py-12">
              <div className="animate-spin rounded-full h-12 w-12 border-4 border-primary-200 border-t-primary-600"></div>
            </div>
          ) : (
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {MEAL_TYPES.map((mealType) => {
                const mealTypemeals = getMealsForType(mealType.id);

                // Group meals by dish id so duplicate dishes show as a single card with a count
                const groupedMealsMap = mealTypemeals.reduce((acc, meal) => {
                  const dishId = meal?.dish?.id || meal?.dishId || `meal-${meal.id}`;
                  if (!acc[dishId]) {
                    acc[dishId] = {
                      key: dishId,
                      count: 0,
                      representative: meal
                    };
                  }
                  acc[dishId].count += 1;
                  return acc;
                }, {});

                const groupedMeals = Object.values(groupedMealsMap);

                return (
                  <div key={mealType.id} className="card-elevated p-6">
                    <div className="flex items-center justify-between mb-6">
                      <div className="flex items-center space-x-3">
                        <div className={`p-3 ${mealType.color} rounded-lg border`}>
                          <span className="flex items-center justify-center">{mealType.icon}</span>
                        </div>
                        <div>
                          <h3 className={`text-xl font-display font-semibold ${mealType.textColor}`}>
                            {mealType.name}
                          </h3>
                          <p className="text-sm text-accent-600">
                            {mealTypemeals.length} meal{mealTypemeals.length !== 1 ? 's' : ''}
                          </p>
                        </div>
                      </div>
                      
                      <button
                        onClick={() => handleAddMeal(mealType.id)}
                        className={`p-3 ${mealType.bgColor} hover:shadow-md text-accent-600 rounded-lg transition-all duration-200 hover:scale-105`}
                      >
                        <Plus className="h-5 w-5" />
                      </button>
                    </div>

                      <div className="space-y-4">
                      {groupedMeals.length === 0 ? (
                        <div className="text-center py-8 text-accent-500">
                          <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-accent-100 flex items-center justify-center">
                            <span className="flex items-center justify-center">{mealType.icon}</span>
                          </div>
                          <p className="text-sm mb-3">No meals planned</p>
                          <button
                            onClick={() => handleAddMeal(mealType.id)}
                            className="btn-secondary text-sm"
                          >
                            Add your first meal
                          </button>
                        </div>
                      ) : (
                        groupedMeals.map((group) => {
                          const meal = group.representative;
                          return (
                            <div key={group.key} className="relative group" tabIndex={0} data-dish-id={group.representative?.dish?.id || group.key}>
                              <MealCard meal={meal} count={group.count} animateDecrease={animatingGroupKey === group.key} />
                              <div className="absolute bottom-2 right-2">
                                <div
                                  className="relative"
                                  onMouseEnter={() => {
                                    if (hideTimerRef.current) { clearTimeout(hideTimerRef.current); hideTimerRef.current = null; }
                                    setHoverMenuKey(group.key);
                                  }}
                                  onMouseLeave={() => {
                                    // delay hiding so users can move to the dropdown without it disappearing
                                    hideTimerRef.current = setTimeout(() => { setHoverMenuKey(null); }, 180);
                                  }}
                                >
                                  {/* Icon visible only when hovering the meal card (parent has 'group') or when menu is open */}
                                  <button
                                    onClick={() => {
                                      if (group.count === 1) {
                                        setConfirmState({
                                          open: true,
                                          title: 'Delete',
                                          description: 'Delete this meal?',
                                          onConfirm: async () => { await handleDeleteOneMeal(meal); }
                                        });
                                      } else {
                                        setOpenMenuKey(openMenuKey === group.key ? null : group.key);
                                      }
                                    }}
                                    onKeyDown={(e) => {
                                      if (e.key === 'Enter' || e.key === ' ') {
                                        e.preventDefault();
                                        if (group.count === 1) {
                                          setConfirmState({ open: true, title: 'Delete', description: 'Delete this meal?', onConfirm: async () => { await handleDeleteOneMeal(meal); } });
                                        } else {
                                          setOpenMenuKey(openMenuKey === group.key ? null : group.key);
                                        }
                                      } else if (e.key === 'Escape') {
                                        setOpenMenuKey(null);
                                        setHoverMenuKey(null);
                                      }
                                    }}
                                    aria-haspopup={group.count > 1}
                                    aria-expanded={openMenuKey === group.key}
                                    aria-label="Meal actions"
                                    className={`p-1.5 bg-white/80 text-neutral-600 rounded-full border border-neutral-200 hover:text-accent-700 hover:bg-accent-50 transition-colors opacity-0 group-hover:opacity-100 ${openMenuKey === group.key || hoverMenuKey === group.key ? 'opacity-100 pointer-events-auto' : 'pointer-events-none'}`}
                                  >
                                    <Trash className="h-4 w-4" />
                                  </button>

                                  {/* Dropdown: appears only when hovering the icon/dropdown or when explicitly opened */}
                                  {group.count > 1 && (openMenuKey === group.key || hoverMenuKey === group.key) && (
                                    <div
                                      ref={menuRef}
                                      onMouseEnter={() => {
                                        if (hideTimerRef.current) { clearTimeout(hideTimerRef.current); hideTimerRef.current = null; }
                                        setHoverMenuKey(group.key);
                                      }}
                                      onMouseLeave={() => {
                                        hideTimerRef.current = setTimeout(() => { setHoverMenuKey(null); }, 180);
                                      }}
                                      className={`absolute right-0 bottom-12 w-44 bg-white rounded-md shadow-lg border border-neutral-100 z-20`}
                                      role="menu"
                                      aria-label="Meal delete menu"
                                      onKeyDown={(e) => {
                                        const items = Array.from(menuRef.current?.querySelectorAll('button')) || [];
                                        const idx = items.indexOf(document.activeElement);
                                        if (e.key === 'ArrowDown') {
                                          e.preventDefault();
                                          const next = items[(idx + 1) % items.length];
                                          next?.focus();
                                        } else if (e.key === 'ArrowUp') {
                                          e.preventDefault();
                                          const prev = items[(idx - 1 + items.length) % items.length];
                                          prev?.focus();
                                        } else if (e.key === 'Escape') {
                                          setOpenMenuKey(null);
                                          setHoverMenuKey(null);
                                        }
                                      }}
                                    >
                                      <button
                                        ref={firstMenuItemRef}
                                        onClick={() => setConfirmState({
                                          open: true,
                                          title: 'Delete one',
                                          description: 'Delete one instance of this meal?',
                                          onConfirm: async () => { await handleDeleteOneMeal(meal); }
                                        })}
                                        className="w-full text-left px-3 py-2 hover:bg-accent-50 flex items-center space-x-2"
                                      >
                                        <Trash className="h-4 w-4 text-neutral-700" />
                                        <span className="text-sm">Delete one</span>
                                      </button>
                                      {group.count > 1 && (
                                        <button
                                          onClick={() => setConfirmState({
                                            open: true,
                                            title: 'Delete all',
                                            description: 'Delete all instances of this dish for this day?',
                                            onConfirm: async () => { await handleDeleteAllMeals(group.key); }
                                          })}
                                          className="w-full text-left px-3 py-2 hover:bg-warm-50 flex items-center space-x-2"
                                        >
                                          <Trash className="h-4 w-4 text-warm-700" />
                                          <span className="text-sm">Delete all</span>
                                        </button>
                                      )}
                                    </div>
                                  )}
                                </div>
                              </div>
                            </div>
                          );
                        })
                      )}
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>

  {/* Recommendations used to be a right column; now render full-width below so Groceries stays on the right */}

        {/* Ingredients Panel */}
        {showIngredientsPanel && (
          <div className="lg:col-span-1">
            <div className="sticky top-6">
              <IngredientsList
                meals={meals}
                selectedDate={selectedDate}
                onClose={() => setShowIngredientsPanel(false)}
              />
            </div>
          </div>
        )}
      </div>

      {/* Full-width Recommendations section shown below meals */}
      {showRecommendations && (
        <div
          className={`mt-6 transition-transform duration-300 ${highlightRecommendations ? 'animate-pulse scale-101 ring-4 ring-accent-200/60 rounded-lg' : ''}`}
          ref={recommendationsRef}
        >
          <MealRecommendations
            onAddMeal={addMeal}
            selectedDate={format(selectedDate, 'yyyy-MM-dd')}
            mealType="lunch"
          />
        </div>
      )}

      {/* Dish Selector Modal */}
      {showDishSelector && (
        <DishSelector
          loadDishes={loadDishes}
          onSelect={handleDishSelect}
          onClose={() => setShowDishSelector(false)}
          mealType={selectedMealType}
          isEditing={!!editingMeal}
          onAddDish={onAddDish}
        />
      )}
      <ConfirmDialog
        open={confirmState.open}
        title={confirmState.title}
        description={confirmState.description}
        onCancel={() => setConfirmState({ open: false, title: '', description: '', onConfirm: null })}
        onConfirm={async () => {
          if (confirmState.onConfirm) await confirmState.onConfirm();
          setConfirmState({ open: false, title: '', description: '', onConfirm: null });
        }}
      />
      <UndoSnackbar
        open={snackbar.open}
        message={snackbar.message}
        onUndo={async () => {
          if (!snackbar.date || !snackbar.dishId) return;
          try {
            await mealPlannerAPI.undoDeleteByDateAndDish(snackbar.date, snackbar.dishId);
            await loadMealsForDate();
          } catch (err) {
            console.error('Undo failed', err);
          } finally {
            setSnackbar({ open: false, message: '', date: null, dishId: null });
          }
        }}
        onClose={() => setSnackbar({ open: false, message: '', date: null, dishId: null })}
      />
    </div>
  );
};

export default DayView;

