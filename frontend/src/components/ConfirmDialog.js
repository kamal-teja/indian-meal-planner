import React from 'react';

const ConfirmDialog = ({ open, title, description, onConfirm, onCancel, confirmLabel = 'Confirm', cancelLabel = 'Cancel' }) => {
  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="absolute inset-0 bg-black/40" onClick={onCancel}></div>
      <div className="bg-white rounded-lg shadow-lg max-w-sm w-full p-6 z-10 border border-neutral-100">
        <h3 className="text-lg font-semibold text-neutral-800 mb-2">{title}</h3>
        <p className="text-sm text-accent-600 mb-4">{description}</p>
        <div className="flex justify-end space-x-3">
          <button onClick={onCancel} className="px-3 py-2 rounded-md bg-white border border-neutral-200 text-neutral-700 hover:bg-accent-50">{cancelLabel}</button>
          <button onClick={onConfirm} className="px-3 py-2 rounded-md bg-warm-600 text-white hover:bg-warm-700">{confirmLabel}</button>
        </div>
      </div>
    </div>
  );
};

export default ConfirmDialog;
