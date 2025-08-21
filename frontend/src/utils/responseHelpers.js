// Response helper utility for handling Go backend responses
// All Go backend responses follow the format: { success: boolean, data: any, message?: string }

/**
 * Extracts data from the standardized Go backend response format
 * @param {Object} response - Axios response object
 * @returns {any} The actual data from response.data.data
 */
export const extractResponseData = (response) => {
  if (response && response.data && response.data.success) {
    return response.data.data;
  }
  
  console.warn('Unexpected response format:', response);
  return null;
};

/**
 * Safely extracts array data from backend response
 * @param {Object} response - Axios response object
 * @param {Array} defaultValue - Default value if data is not an array
 * @returns {Array} The array data or default value
 */
export const extractArrayData = (response, defaultValue = []) => {
  const data = extractResponseData(response);
  return Array.isArray(data) ? data : defaultValue;
};

/**
 * Handles API errors and extracts error messages
 * @param {Error} error - Axios error object
 * @returns {string} User-friendly error message
 */
export const extractErrorMessage = (error) => {
  if (error.response && error.response.data && error.response.data.error) {
    return error.response.data.error;
  }
  
  if (error.message) {
    return error.message;
  }
  
  return 'An unexpected error occurred';
};

// Usage examples:
// 
// const response = await mealPlannerAPI.getMealsByDate(date);
// const meals = extractArrayData(response);
// setMeals(meals);
//
// try {
//   const response = await mealPlannerAPI.getDishes();
//   const dishes = extractArrayData(response);
//   setDishes(dishes);
// } catch (error) {
//   const errorMessage = extractErrorMessage(error);
//   setError(errorMessage);
// }
