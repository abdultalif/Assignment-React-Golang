import React, { useEffect } from 'react';
import type { ToastProps } from '../types/user';

// Toast Notification Component
const Toast: React.FC<ToastProps> = ({ message, type, isVisible, onClose }) => {
  useEffect(() => {
    if (isVisible) {
      const timer = setTimeout(onClose, 4000);
      return () => clearTimeout(timer);
    }
  }, [isVisible, onClose]);

  if (!isVisible) return null;

  return (
    <div className="fixed top-4 right-4 z-50 animate-in slide-in-from-top-8 duration-300">
      <div className={`p-4 rounded-xl shadow-lg backdrop-blur-sm ${
        type === 'success' 
          ? 'bg-green-100/90 dark:bg-green-900/90 text-green-800 dark:text-green-200 border border-green-200 dark:border-green-800' 
          : 'bg-red-100/90 dark:bg-red-900/90 text-red-800 dark:text-red-200 border border-red-200 dark:border-red-800'
      }`}>
        <div className="flex items-center space-x-3">
          <div className="flex-shrink-0">
            {type === 'success' ? (
              <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
            ) : (
              <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            )}
          </div>
          <p className="font-medium">{message}</p>
          <button
            onClick={onClose}
            className="flex-shrink-0 ml-4 p-1 hover:bg-black/10 rounded transition-colors duration-200"
          >
            <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
};

export default Toast;