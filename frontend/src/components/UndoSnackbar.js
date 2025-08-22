import React, { useEffect } from 'react';

const UndoSnackbar = ({ open, message = 'Deleted', onUndo, onClose, duration = 6000 }) => {
  useEffect(() => {
    if (!open) return;
    const t = setTimeout(() => {
      onClose();
    }, duration);
    return () => clearTimeout(t);
  }, [open, duration, onClose]);

  if (!open) return null;

  return (
    <div className="fixed top-6 right-6 z-50">
      <div className="bg-white border border-neutral-200 shadow-md rounded-lg px-4 py-2 flex items-center space-x-4">
        <div className="text-sm text-neutral-800">{message}</div>
        <button onClick={onUndo} className="text-sm text-accent-600 font-medium hover:underline">Undo</button>
        <button onClick={onClose} className="ml-2 text-xs text-neutral-500">Dismiss</button>
      </div>
    </div>
  );
};

export default UndoSnackbar;
