import React, { useState, useEffect } from 'react';
import { ShoppingCart, Calendar, Download, Check, Plus, Minus } from 'lucide-react';
import { mealPlannerAPI } from '../services/api';
import { format, addDays, startOfWeek, endOfWeek } from 'date-fns';

const ShoppingList = () => {
  const [shoppingList, setShoppingList] = useState(null);
  const [loading, setLoading] = useState(false);
  const [dateRange, setDateRange] = useState({
    startDate: format(new Date(), 'yyyy-MM-dd'),
    endDate: format(addDays(new Date(), 6), 'yyyy-MM-dd')
  });
  const [checkedItems, setCheckedItems] = useState({});

  useEffect(() => {
    if (dateRange.startDate && dateRange.endDate) {
      generateShoppingList();
    }
  }, [dateRange]);

  const generateShoppingList = async () => {
    try {
      setLoading(true);
      const response = await mealPlannerAPI.getShoppingList(
        dateRange.startDate,
        dateRange.endDate
      );
      setShoppingList(response.data);
      setCheckedItems({});
    } catch (error) {
      console.error('Error generating shopping list:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleDateRangeChange = (field, value) => {
    setDateRange(prev => ({
      ...prev,
      [field]: value
    }));
  };

  const setQuickRange = (days) => {
    const start = new Date();
    const end = addDays(start, days - 1);
    setDateRange({
      startDate: format(start, 'yyyy-MM-dd'),
      endDate: format(end, 'yyyy-MM-dd')
    });
  };

  const setWeekRange = () => {
    const start = startOfWeek(new Date(), { weekStartsOn: 1 }); // Monday
    const end = endOfWeek(new Date(), { weekStartsOn: 1 }); // Sunday
    setDateRange({
      startDate: format(start, 'yyyy-MM-dd'),
      endDate: format(end, 'yyyy-MM-dd')
    });
  };

  const toggleItem = (itemName) => {
    setCheckedItems(prev => ({
      ...prev,
      [itemName]: !prev[itemName]
    }));
  };

  const exportList = () => {
    if (!shoppingList) return;

    const listText = shoppingList.items
      .map(item => `${checkedItems[item.name] ? '✓' : '☐'} ${item.name} (needed for ${item.count} dishes)`)
      .join('\n');

    const fullText = `Shopping List (${shoppingList.period.startDate} to ${shoppingList.period.endDate})\n\n${listText}`;

    // Create and download file
    const blob = new Blob([fullText], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `shopping-list-${shoppingList.period.startDate}-to-${shoppingList.period.endDate}.txt`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  const checkedCount = Object.values(checkedItems).filter(Boolean).length;
  const totalItems = shoppingList?.totalItems || 0;

  return (
    <div className="max-w-4xl mx-auto">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center space-x-3 mb-6">
          <div className="p-2 bg-gradient-to-r from-green-500 to-emerald-600 rounded-xl shadow-lg">
            <ShoppingCart className="h-8 w-8 text-white" />
          </div>
          <div>
            <h1 className="text-3xl font-display font-bold gradient-text">
              Shopping List
            </h1>
            <p className="text-gray-600">
              Generate ingredient lists from your meal plans
            </p>
          </div>
        </div>

        {/* Date Range Selector */}
        <div className="bg-white rounded-xl shadow-md p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center space-x-2">
            <Calendar className="h-5 w-5" />
            <span>Select Date Range</span>
          </h3>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Start Date
              </label>
              <input
                type="date"
                value={dateRange.startDate}
                onChange={(e) => handleDateRangeChange('startDate', e.target.value)}
                className="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                End Date
              </label>
              <input
                type="date"
                value={dateRange.endDate}
                onChange={(e) => handleDateRangeChange('endDate', e.target.value)}
                className="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              />
            </div>
          </div>

          <div className="flex flex-wrap gap-2">
            <button
              onClick={() => setQuickRange(3)}
              className="px-3 py-2 text-sm bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors"
            >
              Next 3 days
            </button>
            <button
              onClick={() => setQuickRange(7)}
              className="px-3 py-2 text-sm bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors"
            >
              Next 7 days
            </button>
            <button
              onClick={setWeekRange}
              className="px-3 py-2 text-sm bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors"
            >
              This week
            </button>
          </div>
        </div>
      </div>

      {/* Shopping List */}
      {loading ? (
        <div className="text-center py-16">
          <div className="animate-spin rounded-full h-16 w-16 border-4 border-primary-200 border-t-primary-600 mx-auto mb-4"></div>
          <h2 className="text-2xl font-display font-semibold gradient-text">
            Generating Shopping List...
          </h2>
        </div>
      ) : shoppingList ? (
        <div className="space-y-6">
          {/* Summary */}
          <div className="bg-white rounded-xl shadow-md p-6">
            <div className="flex items-center justify-between mb-4">
              <div>
                <h3 className="text-xl font-semibold text-gray-900">
                  Shopping List Summary
                </h3>
                <p className="text-gray-600">
                  {format(new Date(shoppingList.period.startDate), 'MMM dd')} -{' '}
                  {format(new Date(shoppingList.period.endDate), 'MMM dd, yyyy')}
                </p>
              </div>
              <button
                onClick={exportList}
                className="flex items-center space-x-2 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
              >
                <Download className="h-4 w-4" />
                <span>Export</span>
              </button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="bg-blue-50 rounded-lg p-4">
                <div className="text-2xl font-bold text-blue-700">{totalItems}</div>
                <div className="text-sm text-blue-600">Total Items</div>
              </div>
              <div className="bg-green-50 rounded-lg p-4">
                <div className="text-2xl font-bold text-green-700">{checkedCount}</div>
                <div className="text-sm text-green-600">Items Checked</div>
              </div>
              <div className="bg-orange-50 rounded-lg p-4">
                <div className="text-2xl font-bold text-orange-700">{totalItems - checkedCount}</div>
                <div className="text-sm text-orange-600">Items Remaining</div>
              </div>
            </div>

            {/* Progress Bar */}
            <div className="mt-4">
              <div className="flex items-center justify-between text-sm text-gray-600 mb-1">
                <span>Progress</span>
                <span>{totalItems > 0 ? Math.round((checkedCount / totalItems) * 100) : 0}%</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div 
                  className="bg-green-500 h-2 rounded-full transition-all duration-300"
                  style={{ width: `${totalItems > 0 ? (checkedCount / totalItems) * 100 : 0}%` }}
                ></div>
              </div>
            </div>
          </div>

          {/* Shopping Items */}
          <div className="bg-white rounded-xl shadow-md p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              Ingredients List
            </h3>
            
            {shoppingList.items.length > 0 ? (
              <div className="space-y-3">
                {shoppingList.items.map((item, index) => (
                  <div
                    key={index}
                    className={`flex items-center space-x-4 p-4 border rounded-lg transition-all duration-200 ${
                      checkedItems[item.name]
                        ? 'bg-green-50 border-green-200'
                        : 'hover:bg-gray-50'
                    }`}
                  >
                    <button
                      onClick={() => toggleItem(item.name)}
                      className={`flex-shrink-0 w-6 h-6 rounded-full border-2 flex items-center justify-center transition-colors ${
                        checkedItems[item.name]
                          ? 'bg-green-500 border-green-500 text-white'
                          : 'border-gray-300 hover:border-green-400'
                      }`}
                    >
                      {checkedItems[item.name] && <Check className="h-3 w-3" />}
                    </button>

                    <div className="flex-1 min-w-0">
                      <p className={`text-lg font-medium ${
                        checkedItems[item.name] 
                          ? 'text-green-800 line-through' 
                          : 'text-gray-900'
                      }`}>
                        {item.name}
                      </p>
                      <p className="text-sm text-gray-500">
                        Needed for {item.count} dish{item.count !== 1 ? 'es' : ''}: {item.dishes.join(', ')}
                      </p>
                    </div>

                    <div className="flex-shrink-0">
                      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                        checkedItems[item.name]
                          ? 'bg-green-100 text-green-800'
                          : 'bg-gray-100 text-gray-800'
                      }`}>
                        {item.count}x
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-8">
                <ShoppingCart className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <p className="text-gray-500">
                  No ingredients found for the selected date range.
                  Try planning some meals first!
                </p>
              </div>
            )}
          </div>
        </div>
      ) : (
        <div className="text-center py-16">
          <ShoppingCart className="h-16 w-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-2xl font-display font-semibold text-gray-700 mb-2">
            Ready to Generate Your Shopping List
          </h3>
          <p className="text-gray-500">
            Select a date range above to generate a shopping list from your planned meals.
          </p>
        </div>
      )}
    </div>
  );
};

export default ShoppingList;
