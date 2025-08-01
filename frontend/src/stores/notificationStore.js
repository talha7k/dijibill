import { writable } from 'svelte/store';

// Create a writable store for notifications
export const notifications = writable([]);

// Notification types
export const NOTIFICATION_TYPES = {
  SUCCESS: 'success',
  ERROR: 'error',
  WARNING: 'warning',
  INFO: 'info'
};

// Auto-dismiss timeouts for different types
const AUTO_DISMISS_TIMES = {
  [NOTIFICATION_TYPES.SUCCESS]: 4000,
  [NOTIFICATION_TYPES.INFO]: 5000,
  [NOTIFICATION_TYPES.WARNING]: 6000,
  [NOTIFICATION_TYPES.ERROR]: 8000
};

let notificationId = 0;

/**
 * Add a notification to the store
 * @param {Object} notification - The notification object
 * @param {string} notification.type - Type of notification (success, error, warning, info)
 * @param {string} notification.title - Title of the notification
 * @param {string} notification.message - Message content
 * @param {boolean} notification.persistent - Whether the notification should auto-dismiss
 * @param {number} notification.duration - Custom duration in ms (overrides default)
 */
export function addNotification({ type = NOTIFICATION_TYPES.INFO, title, message, persistent = false, duration }) {
  const id = ++notificationId;
  const notification = {
    id,
    type,
    title,
    message,
    persistent,
    duration: duration || AUTO_DISMISS_TIMES[type],
    timestamp: Date.now()
  };

  // Add to store
  notifications.update(items => [...items, notification]);

  // Auto-dismiss if not persistent
  if (!persistent) {
    setTimeout(() => {
      removeNotification(id);
    }, notification.duration);
  }

  return id;
}

/**
 * Remove a notification by ID
 * @param {number} id - The notification ID to remove
 */
export function removeNotification(id) {
  notifications.update(items => items.filter(item => item.id !== id));
}

/**
 * Clear all notifications
 */
export function clearAllNotifications() {
  notifications.set([]);
}

// Convenience functions for different notification types
export function showSuccess(title, message, options = {}) {
  return addNotification({
    ...options,
    type: NOTIFICATION_TYPES.SUCCESS,
    title,
    message
  });
}

export function showError(title, message, options = {}) {
  return addNotification({
    ...options,
    type: NOTIFICATION_TYPES.ERROR,
    title,
    message
  });
}

export function showWarning(title, message, options = {}) {
  return addNotification({
    ...options,
    type: NOTIFICATION_TYPES.WARNING,
    title,
    message
  });
}

export function showInfo(title, message, options = {}) {
  return addNotification({
    ...options,
    type: NOTIFICATION_TYPES.INFO,
    title,
    message
  });
}

// Database operation helpers
export function showDbSuccess(operation, entity) {
  const messages = {
    create: `${entity} created successfully!`,
    update: `${entity} updated successfully!`,
    delete: `${entity} deleted successfully!`,
    save: `${entity} saved successfully!`
  };
  
  return showSuccess('Success', messages[operation] || `${operation} completed successfully!`);
}

export function showDbError(operation, entity, error) {
  const messages = {
    create: `Failed to create ${entity}`,
    update: `Failed to update ${entity}`,
    delete: `Failed to delete ${entity}`,
    save: `Failed to save ${entity}`,
    load: `Failed to load ${entity}`
  };
  
  const title = messages[operation] || `${operation} failed`;
  const message = error?.message || error || 'An unexpected error occurred. Please try again.';
  
  return showError(title, message);
}