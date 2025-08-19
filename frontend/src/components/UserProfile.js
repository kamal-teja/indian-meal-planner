import React, { useState, useEffect } from 'react';
import { User, Mail, Settings, Save, AlertCircle, CheckCircle } from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';

const UserProfile = () => {
  const { user, updateUserProfile } = useAuth();
  const [formData, setFormData] = useState({
    name: '',
    profile: {
      dietaryPreferences: [],
      spiceLevel: 'medium',
      favoriteRegions: []
    }
  });
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState({ type: '', text: '' });

  useEffect(() => {
    if (user) {
      setFormData({
        name: user.name || '',
        profile: {
          dietaryPreferences: user.profile?.dietaryPreferences || [],
          spiceLevel: user.profile?.spiceLevel || 'medium',
          favoriteRegions: user.profile?.favoriteRegions || []
        }
      });
    }
  }, [user]);

  const handleInputChange = (field, value) => {
    if (field === 'name') {
      setFormData(prev => ({ ...prev, name: value }));
    } else {
      setFormData(prev => ({
        ...prev,
        profile: { ...prev.profile, [field]: value }
      }));
    }
  };

  const handleArrayChange = (field, value, checked) => {
    setFormData(prev => ({
      ...prev,
      profile: {
        ...prev.profile,
        [field]: checked
          ? [...prev.profile[field], value]
          : prev.profile[field].filter(item => item !== value)
      }
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setMessage({ type: '', text: '' });

    try {
      const result = await updateUserProfile(formData);
      if (result.success) {
        setMessage({ type: 'success', text: 'Profile updated successfully!' });
      } else {
        setMessage({ type: 'error', text: result.error || 'Failed to update profile' });
      }
    } catch (error) {
      setMessage({ type: 'error', text: 'An unexpected error occurred' });
    } finally {
      setLoading(false);
    }
  };

  const dietaryOptions = [
    'vegetarian',
    'vegan',
    'gluten-free',
    'dairy-free',
    'low-carb',
    'high-protein'
  ];

  const spiceLevels = [
    { value: 'mild', label: 'Mild', emoji: '🟢' },
    { value: 'medium', label: 'Medium', emoji: '🟡' },
    { value: 'hot', label: 'Hot', emoji: '🔥' },
    { value: 'extra-hot', label: 'Extra Hot', emoji: '🌶️' }
  ];

  const regionOptions = [
    'North Indian',
    'South Indian',
    'Bengali',
    'Gujarati',
    'Punjabi',
    'Rajasthani',
    'Maharashtrian'
  ];

  if (!user) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-primary-50 to-secondary-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-16 w-16 border-4 border-primary-200 border-t-primary-600 mx-auto mb-4"></div>
          <h2 className="text-2xl font-display font-semibold gradient-text">
            Loading Profile...
          </h2>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center space-x-3 mb-4">
          <div className="p-2 bg-gradient-to-r from-blue-500 to-indigo-600 rounded-xl shadow-lg">
            <User className="h-8 w-8 text-white" />
          </div>
          <div>
            <h1 className="text-3xl font-display font-bold gradient-text">
              Profile Settings
            </h1>
            <p className="text-gray-600">
              Manage your account and preferences
            </p>
          </div>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Basic Information */}
        <div className="bg-white rounded-xl shadow-md p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center space-x-2">
            <Settings className="h-5 w-5" />
            <span>Basic Information</span>
          </h3>

          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Full Name
              </label>
              <input
                type="text"
                value={formData.name}
                onChange={(e) => handleInputChange('name', e.target.value)}
                className="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                placeholder="Enter your full name"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Email Address
              </label>
              <div className="flex items-center space-x-2 px-3 py-2 border border-gray-300 rounded-lg bg-gray-50">
                <Mail className="h-5 w-5 text-gray-400" />
                <span className="text-gray-600">{user.email}</span>
                <span className="text-xs text-gray-500">(cannot be changed)</span>
              </div>
            </div>
          </div>
        </div>

        {/* Dietary Preferences */}
        <div className="bg-white rounded-xl shadow-md p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">
            Dietary Preferences
          </h3>
          <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
            {dietaryOptions.map(option => (
              <label key={option} className="flex items-center space-x-2 cursor-pointer">
                <input
                  type="checkbox"
                  checked={formData.profile.dietaryPreferences.includes(option)}
                  onChange={(e) => handleArrayChange('dietaryPreferences', option, e.target.checked)}
                  className="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                />
                <span className="text-sm text-gray-700 capitalize">
                  {option.replace('-', ' ')}
                </span>
              </label>
            ))}
          </div>
        </div>

        {/* Spice Level */}
        <div className="bg-white rounded-xl shadow-md p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">
            Preferred Spice Level
          </h3>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
            {spiceLevels.map(level => (
              <label key={level.value} className="cursor-pointer">
                <input
                  type="radio"
                  name="spiceLevel"
                  value={level.value}
                  checked={formData.profile.spiceLevel === level.value}
                  onChange={(e) => handleInputChange('spiceLevel', e.target.value)}
                  className="sr-only"
                />
                <div className={`p-3 border-2 rounded-lg text-center transition-colors ${
                  formData.profile.spiceLevel === level.value
                    ? 'border-primary-500 bg-primary-50'
                    : 'border-gray-200 hover:border-gray-300'
                }`}>
                  <div className="text-2xl mb-1">{level.emoji}</div>
                  <div className="text-sm font-medium text-gray-700">{level.label}</div>
                </div>
              </label>
            ))}
          </div>
        </div>

        {/* Favorite Regions */}
        <div className="bg-white rounded-xl shadow-md p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">
            Favorite Regional Cuisines
          </h3>
          <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
            {regionOptions.map(region => (
              <label key={region} className="flex items-center space-x-2 cursor-pointer">
                <input
                  type="checkbox"
                  checked={formData.profile.favoriteRegions.includes(region)}
                  onChange={(e) => handleArrayChange('favoriteRegions', region, e.target.checked)}
                  className="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                />
                <span className="text-sm text-gray-700">{region}</span>
              </label>
            ))}
          </div>
        </div>

        {/* Message */}
        {message.text && (
          <div className={`flex items-center space-x-2 p-4 rounded-lg ${
            message.type === 'success' 
              ? 'bg-green-50 text-green-700 border border-green-200'
              : 'bg-red-50 text-red-700 border border-red-200'
          }`}>
            {message.type === 'success' ? (
              <CheckCircle className="h-5 w-5" />
            ) : (
              <AlertCircle className="h-5 w-5" />
            )}
            <span>{message.text}</span>
          </div>
        )}

        {/* Submit Button */}
        <div className="flex justify-end">
          <button
            type="submit"
            disabled={loading}
            className="flex items-center space-x-2 px-6 py-3 bg-primary-600 text-white rounded-lg hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {loading ? (
              <div className="animate-spin rounded-full h-5 w-5 border-2 border-white border-t-transparent" />
            ) : (
              <Save className="h-5 w-5" />
            )}
            <span>{loading ? 'Saving...' : 'Save Changes'}</span>
          </button>
        </div>
      </form>
    </div>
  );
};

export default UserProfile;
