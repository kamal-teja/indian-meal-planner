import React, { useState, useRef, useEffect } from 'react';
import { ChevronDown, Check } from 'lucide-react';

const CustomDropdown = ({ 
  value, 
  onChange, 
  options, 
  placeholder = "Select an option",
  className = "",
  disabled = false,
  error = false 
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef(null);

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  const handleSelect = (option) => {
    onChange(option.value);
    setIsOpen(false);
  };

  const selectedOption = options.find(option => option.value === value);

  return (
    <div className={`relative ${className}`} ref={dropdownRef}>
      <button
        type="button"
        onClick={() => !disabled && setIsOpen(!isOpen)}
        disabled={disabled}
        className={`
          w-full appearance-none bg-white rounded-xl px-4 py-2.5 
          text-neutral-700 font-medium shadow-sm transition-all duration-200 hover:shadow-md 
          focus:ring-2 focus:ring-accent-500 focus:outline-none cursor-pointer text-left relative
          ${error 
            ? 'border border-red-300 hover:border-red-400 focus:border-red-500 focus:ring-red-500' 
            : 'border border-neutral-300 hover:border-accent-400 focus:border-accent-500'
          }
          ${disabled ? 'opacity-50 cursor-not-allowed' : ''}
          ${isOpen ? (error ? 'border-red-500 ring-2 ring-red-500' : 'border-accent-500 ring-2 ring-accent-500') : ''}
        `}
      >
        <span className={`block pr-8 ${selectedOption ? 'text-neutral-700' : 'text-neutral-400'}`}>
          {selectedOption ? selectedOption.label : placeholder}
        </span>
        <ChevronDown 
          className={`h-4 w-4 text-neutral-500 transition-transform duration-200 absolute right-3 top-1/2 transform -translate-y-1/2 ${
            isOpen ? 'rotate-180' : ''
          }`} 
        />
      </button>

      {isOpen && !disabled && (
        <div className="absolute z-50 w-full mt-1 bg-white border border-neutral-200 rounded-xl shadow-lg max-h-60 overflow-auto">
          <div className="py-1">
            {options.map((option) => (
              <button
                key={option.value}
                type="button"
                onClick={() => handleSelect(option)}
                className={`
                  w-full px-4 py-2.5 text-left hover:bg-accent-50 focus:bg-accent-50 
                  focus:outline-none transition-colors duration-150 flex items-center justify-between
                  ${value === option.value ? 'bg-accent-50 text-accent-700' : 'text-neutral-700'}
                `}
              >
                <span>{option.label}</span>
                {value === option.value && (
                  <Check className="h-4 w-4 text-accent-600" />
                )}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default CustomDropdown;
