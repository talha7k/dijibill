/**
 * Centralized icon utility for Font Awesome icon mapping
 * This ensures consistent icon usage across the entire application
 */

// Comprehensive icon mapping for all common use cases
const iconMap = {
  // Navigation & UI
  'menu': 'fas fa-bars',
  'close': 'fas fa-times',
  'search': 'fas fa-search',
  'filter': 'fas fa-filter',
  'sort': 'fas fa-sort',
  'refresh': 'fas fa-sync-alt',
  'settings': 'fas fa-cog',
  'home': 'fas fa-home',
  'dashboard': 'fas fa-tachometer-alt',
  
  // Actions
  'add': 'fas fa-plus',
  'edit': 'fas fa-edit',
  'delete': 'fas fa-trash',
  'save': 'fas fa-save',
  'cancel': 'fas fa-times',
  'confirm': 'fas fa-check',
  'copy': 'fas fa-copy',
  'download': 'fas fa-download',
  'upload': 'fas fa-upload',
  'print': 'fas fa-print',
  
  // Status & Feedback
  'loading': 'fas fa-spinner',
  'success': 'fas fa-check-circle',
  'error': 'fas fa-exclamation-circle',
  'warning': 'fas fa-exclamation-triangle',
  'info': 'fas fa-info-circle',
  'check': 'fas fa-check',
  'times': 'fas fa-times',
  
  // Business Objects
  'product': 'fas fa-box',
  'products': 'fas fa-boxes',
  'invoice': 'fas fa-file-invoice',
  'invoices': 'fas fa-file-invoice-dollar',
  'customer': 'fas fa-user',
  'customers': 'fas fa-users',
  'supplier': 'fas fa-truck',
  'suppliers': 'fas fa-shipping-fast',
  'payment': 'fas fa-credit-card',
  'payments': 'fas fa-money-bill-wave',
  'category': 'fas fa-tags',
  'categories': 'fas fa-tag',
  'report': 'fas fa-chart-bar',
  'reports': 'fas fa-chart-line',
  
  // User & Auth
  'user': 'fas fa-user',
  'users': 'fas fa-users',
  'login': 'fas fa-sign-in-alt',
  'logout': 'fas fa-sign-out-alt',
  'profile': 'fas fa-user-circle',
  'lock': 'fas fa-lock',
  'unlock': 'fas fa-unlock',
  
  // Data & Files
  'file': 'fas fa-file',
  'folder': 'fas fa-folder',
  'image': 'fas fa-image',
  'document': 'fas fa-file-alt',
  'pdf': 'fas fa-file-pdf',
  'excel': 'fas fa-file-excel',
  'csv': 'fas fa-file-csv',
  
  // Communication
  'email': 'fas fa-envelope',
  'phone': 'fas fa-phone',
  'message': 'fas fa-comment',
  'notification': 'fas fa-bell',
  
  // Commerce
  'cart': 'fas fa-shopping-cart',
  'pos': 'fas fa-cash-register',
  'barcode': 'fas fa-barcode',
  'qr': 'fas fa-qrcode',
  'price': 'fas fa-dollar-sign',
  'discount': 'fas fa-percent',
  'tax': 'fas fa-receipt',
  
  // Navigation arrows
  'arrow-left': 'fas fa-arrow-left',
  'arrow-right': 'fas fa-arrow-right',
  'arrow-up': 'fas fa-arrow-up',
  'arrow-down': 'fas fa-arrow-down',
  'chevron-left': 'fas fa-chevron-left',
  'chevron-right': 'fas fa-chevron-right',
  'chevron-up': 'fas fa-chevron-up',
  'chevron-down': 'fas fa-chevron-down',
  
  // Common UI elements
  'calendar': 'fas fa-calendar',
  'clock': 'fas fa-clock',
  'location': 'fas fa-map-marker-alt',
  'link': 'fas fa-link',
  'external-link': 'fas fa-external-link-alt',
  'eye': 'fas fa-eye',
  'eye-slash': 'fas fa-eye-slash',
  'heart': 'fas fa-heart',
  'star': 'fas fa-star',
  'bookmark': 'fas fa-bookmark',
  
  // System
  'database': 'fas fa-database',
  'server': 'fas fa-server',
  'cloud': 'fas fa-cloud',
  'wifi': 'fas fa-wifi',
  'signal': 'fas fa-signal',
  'battery': 'fas fa-battery-full',
  
  // Fallback icons
  'table': 'fas fa-table',
  'list': 'fas fa-list',
  'grid': 'fas fa-th',
  'rocket': 'fas fa-rocket',
  'magic': 'fas fa-magic'
}

/**
 * Get Font Awesome class for an icon name
 * @param {string} iconName - The icon name or Font Awesome class
 * @param {string} fallback - Fallback icon if not found
 * @returns {string} Font Awesome class string
 */
export function getFaIcon(iconName, fallback = 'fas fa-circle') {
  if (!iconName) return fallback
  
  // If it's already a Font Awesome class, return as is
  if (iconName.includes('fa-')) {
    return iconName.startsWith('fa') ? iconName : `fas ${iconName}`
  }
  
  // Look up in our icon map
  return iconMap[iconName] || fallback
}

/**
 * Get icon HTML element string for use in templates
 * @param {string} iconName - The icon name
 * @param {string} className - Additional CSS classes
 * @param {string} fallback - Fallback icon if not found
 * @returns {string} HTML string for the icon
 */
export function getIconHtml(iconName, className = '', fallback = 'fas fa-circle') {
  const iconClass = getFaIcon(iconName, fallback)
  const combinedClass = className ? `${iconClass} ${className}` : iconClass
  return `<i class="${combinedClass}"></i>`
}

/**
 * Create an icon element programmatically
 * @param {string} iconName - The icon name
 * @param {string} className - Additional CSS classes
 * @param {string} fallback - Fallback icon if not found
 * @returns {HTMLElement} Icon element
 */
export function createIconElement(iconName, className = '', fallback = 'fas fa-circle') {
  const i = document.createElement('i')
  const iconClass = getFaIcon(iconName, fallback)
  i.className = className ? `${iconClass} ${className}` : iconClass
  return i
}

/**
 * Get all available icon names
 * @returns {string[]} Array of available icon names
 */
export function getAvailableIcons() {
  return Object.keys(iconMap)
}

export default {
  getFaIcon,
  getIconHtml,
  createIconElement,
  getAvailableIcons
}