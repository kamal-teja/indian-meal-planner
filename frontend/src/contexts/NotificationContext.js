import React, { createContext, useContext, useState, useCallback } from 'react';

const NotificationContext = createContext(null);

// Notification shape: { id, message, variant: 'info'|'success'|'error', duration, action: { label, onClick } }
let nextId = 1;

export const NotificationProvider = ({ children }) => {
  const [queue, setQueue] = useState([]);

  const notify = useCallback((message, opts = {}) => {
    const id = `n_${nextId++}`;
    const notification = {
      id,
      message,
      variant: opts.variant || 'info',
      duration: opts.duration || 4000,
      action: opts.action || null,
    };
    setQueue((q) => [...q, notification]);
    // Auto-remove after duration
    setTimeout(() => setQueue((q) => q.filter(n => n.id !== id)), notification.duration + 50);
    return id;
  }, []);

  const dismiss = useCallback((id) => setQueue((q) => q.filter(n => n.id !== id)), []);

  return (
    <NotificationContext.Provider value={{ queue, notify, dismiss }}>
      {children}
    </NotificationContext.Provider>
  );
};

export const useNotification = () => {
  const ctx = useContext(NotificationContext);
  if (!ctx) throw new Error('useNotification must be used within NotificationProvider');
  return ctx;
};

export default NotificationContext;
