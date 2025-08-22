import React from 'react';
import { useNotification } from '../contexts/NotificationContext';

const variantClasses = {
  info: 'bg-white border-neutral-200 text-neutral-800',
  success: 'bg-green-50 border-green-200 text-green-800',
  error: 'bg-red-50 border-red-200 text-red-800',
};

const NotificationSnackbar = () => {
  const { queue, dismiss } = useNotification();
  if (!queue || queue.length === 0) return null;

  return (
    <div className="fixed top-6 right-6 z-50 flex flex-col space-y-3">
      {queue.map((n) => (
        <div key={n.id} className={`rounded-lg px-4 py-2 shadow-md border flex items-center justify-between space-x-4 ${variantClasses[n.variant] || variantClasses.info}`}>
          <div className="text-sm">{n.message}</div>
          <div className="flex items-center space-x-2">
            {n.action && (
              <button onClick={() => { n.action.onClick(); dismiss(n.id); }} className="text-sm font-medium text-accent-600 hover:underline">
                {n.action.label}
              </button>
            )}
            <button onClick={() => dismiss(n.id)} className="text-xs text-neutral-500">Dismiss</button>
          </div>
        </div>
      ))}
    </div>
  );
};

export default NotificationSnackbar;
