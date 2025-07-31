<script>
  import { createEventDispatcher } from 'svelte'
  import FormField from './FormField.svelte'
  
  /** @type {Array<any>} */
  export let data = []
  /** @type {Array<{label: string, key?: string, class?: string, render?: Function, actions?: Array<{key: string, text: string, icon?: string, class?: string, title?: string}>}>} */
  export let columns = []
  /** @type {boolean} */
  export let loading = false
  /** @type {string} */
  export let searchTerm = ''
  /** @type {string} */
  export let searchPlaceholder = 'Search...'
  /** @type {string} */
  export let emptyStateTitle = 'No data found'
  /** @type {string} */
  export let emptyStateMessage = 'No items to display'
  /** @type {string} */
  export let emptyStateIcon = 'fa-table'
  /** @type {boolean} */
  export let showSearch = true
  /** @type {boolean} */
  export let showActions = true
  /** @type {{icon?: string, text?: string} | null} */
  export let primaryAction = null
  /** @type {Array<{icon?: string, text?: string}>} */
  export let secondaryActions = []

  const dispatch = createEventDispatcher()

  function handleSearch(event) {
    dispatch('search', { searchTerm: event.target.value })
  }

  function handlePrimaryAction() {
    if (primaryAction) {
      dispatch('primary-action')
    }
  }

  function handleSecondaryAction(action) {
    dispatch('secondary-action', { action })
  }

  // Helper function to convert FontAwesome icon names to SVG icons
  function getSvgIcon(iconName) {
    const icons = {
      'fa-table': 'M3 3h18v18H3V3zm2 2v14h14V5H5zm2 2h10v2H7V7zm0 4h10v2H7v-2zm0 4h10v2H7v-2z',
      'fa-file-invoice': 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z',
      'fa-users': 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z',
      'fa-credit-card': 'M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z',
      'fa-shopping-cart': 'M3 3h2l.4 2M7 13h10l4-8H5.4m0 0L7 13m0 0l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17M17 13v4a2 2 0 01-2 2H9a2 2 0 01-2-2v-4m8 0V9a2 2 0 00-2-2H9a2 2 0 00-2 2v4.01',
      'fa-truck': 'M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z',
      'fa-ruler': 'M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zM7 3H5v12a2 2 0 104 0V3z',
      'fa-percentage': 'M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z',
      'fa-box': 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4',
      'fa-tags': 'M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z',
      'fa-plus': 'M12 4v16m8-8H4',
      'fa-edit': 'M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z',
      'fa-trash': 'M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16',
      'fa-eye': 'M15 12a3 3 0 11-6 0 3 3 0 016 0z M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z'
    }
    return icons[iconName] || icons['fa-table'] // fallback to table icon
  }

  function handleRowAction(action, item) {
    dispatch('row-action', { action, item })
  }
</script>

<div class="data-table-container">
  <div class="table-header">
    <div class="table-controls">
      {#if showActions}
        <div class="action-buttons">
          {#if primaryAction}
            <button 
              class="btn btn-primary" 
              on:click={handlePrimaryAction}
              disabled={loading}
            >
              {#if primaryAction && 'icon' in primaryAction && primaryAction.icon}
                <i class="fas {primaryAction.icon}"></i>
              {/if}
              {primaryAction && 'text' in primaryAction ? primaryAction.text : 'Action'}
            </button>
          {/if}
          
          {#each secondaryActions as action}
            <button 
              class="btn btn-secondary-outline" 
              on:click={() => handleSecondaryAction(action)}
              disabled={loading}
            >
              {#if action && 'icon' in action && action.icon}
                <i class="fas {action.icon}"></i>
              {/if}
              {action && 'text' in action ? action.text : 'Action'}
            </button>
          {/each}
        </div>
      {/if}

      {#if showSearch}
        <div class="search-container">
          <div class="search-icon">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
            </svg>
          </div>
          <FormField
            type="text"
            placeholder={searchPlaceholder}
            bind:value={searchTerm}
            disabled={loading}
            label=""
            customClass="pl-10"
            on:input={handleSearch}
          />
        </div>
      {/if}
    </div>
  </div>

  <div class="table-content">
    {#if loading}
      <div class="loading-state">
        <svg class="w-5 h-5 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
        </svg>
        <span>Loading...</span>
      </div>
    {:else if data.length === 0}
      <div class="empty-state">
        <svg class="w-12 h-12 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={getSvgIcon(emptyStateIcon)}></path>
        </svg>
        <h3>{emptyStateTitle}</h3>
        <p>
          {searchTerm ? 'No items match your search criteria.' : emptyStateMessage}
        </p>
        {#if !searchTerm && primaryAction}
          <button class="btn btn-primary" on:click={handlePrimaryAction}>
            {#if primaryAction.icon}
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={getSvgIcon(primaryAction.icon)}></path>
              </svg>
            {/if}
            {primaryAction.text}
          </button>
        {/if}
      </div>
    {:else}
      <div class="table-wrapper">
        <table class="data-table">
          <thead>
            <tr>
              {#each columns as column}
                <th class={column.class || ''}>{column.label}</th>
              {/each}
              {#if columns.some(col => col.actions)}
                <th class="actions-column">Actions</th>
              {/if}
            </tr>
          </thead>
          <tbody>
            {#each data as item, index}
              <tr>
                {#each columns as column}
                  <td class={column.class || ''}>
                    <slot name="cell" {item} {column} {index}>
                      {#if column.render}
                        {@html column.render(item)}
                      {:else if column.key}
                        {item[column.key] || '-'}
                      {:else}
                        -
                      {/if}
                    </slot>
                  </td>
                {/each}
                {#if columns.some(col => col.actions)}
                  <td class="actions-column">
                    <div class="action-buttons">
                      <slot name="actions" {item} {index}>
                        {#each (columns.find(col => col.actions)?.actions || []) as action}
                          <button
                            class="btn btn-sm {action.class || 'btn-secondary'}"
                            on:click={() => handleRowAction(action.key, item)}
                            title={action.title || action.text}
                            disabled={loading}
                          >
                            {#if action.icon}
                              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={getSvgIcon(action.icon)}></path>
                              </svg>
                            {:else}
                              {action.text}
                            {/if}
                          </button>
                        {/each}
                      </slot>
                    </div>
                  </td>
                {/if}
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </div>
</div>

<style>
  .data-table-container {
    display: flex;
    flex-direction: column;
    height: 100%;
  }

  .table-header {
    padding: 1.5rem;
    border-bottom: 1px solid var(--color-border-light);
  }

  .table-controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
  }

  .action-buttons {
    display: flex;
    gap: 0.75rem;
    align-items: center;
  }

  .search-container {
    position: relative;
    min-width: 300px;
  }

  .search-icon {
    position: absolute;
    left: 1rem;
    top: 50%;
    transform: translateY(-50%);
    color: var(--color-text-muted);
    font-size: 0.875rem;
  }

  .search-input {
    width: 100%;
    padding: 0.75rem 1rem 0.75rem 2.5rem;
    background: var(--color-glass-white);
    border: 1px solid var(--color-border-light);
    border-radius: var(--radius-md);
    color: var(--color-text-primary);
    font-size: 0.875rem;
    transition: var(--transition-normal);
  }

  .search-input::placeholder {
    color: var(--color-text-muted);
  }

  .search-input:focus {
    outline: none;
    border-color: var(--color-primary);
    box-shadow: 0 0 0 3px var(--color-primary-100);
  }

  .table-content {
    flex: 1;
    overflow: hidden;
  }

  .loading-state {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 3rem;
    color: var(--color-text-secondary);
    gap: 0.75rem;
  }

  .loading-state svg {
    width: 1.25rem;
    height: 1.25rem;
  }

  .animate-spin {
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 3rem;
    text-align: center;
  }

  .empty-state svg {
    color: var(--color-text-muted);
    margin-bottom: 1rem;
  }

  .empty-state h3 {
    margin: 0 0 0.5rem 0;
    color: var(--color-text-primary);
    font-size: 1.25rem;
    font-weight: 600;
  }

  .empty-state p {
    margin: 0 0 1.5rem 0;
    color: var(--color-text-secondary);
    font-size: 0.875rem;
  }

  .table-wrapper {
    overflow: auto;
    height: 100%;
    /* Enhanced scrollbar visibility */
    scrollbar-width: thin;
    scrollbar-color: rgba(255, 255, 255, 0.3) rgba(255, 255, 255, 0.1);
  }

  /* Webkit scrollbar styling for better visibility */
  .table-wrapper::-webkit-scrollbar {
    width: 8px;
    height: 8px;
  }

  .table-wrapper::-webkit-scrollbar-track {
    background: rgba(255, 255, 255, 0.1);
    border-radius: 4px;
  }

  .table-wrapper::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.3);
    border-radius: 4px;
    transition: background 0.2s ease;
  }

  .table-wrapper::-webkit-scrollbar-thumb:hover {
    background: rgba(255, 255, 255, 0.5);
  }

  .table-wrapper::-webkit-scrollbar-corner {
    background: rgba(255, 255, 255, 0.1);
  }

  .data-table {
    width: 100%;
    min-width: 800px; /* Ensure minimum width to trigger horizontal scroll */
    border-collapse: collapse;
  }

  .data-table th {
    padding: 1rem;
    text-align: left;
    font-weight: 600;
    color: var(--color-text-secondary);
    font-size: 0.875rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    border-bottom: 1px solid var(--color-border-light);
    background: var(--color-glass-white-light);
    white-space: nowrap; /* Prevent header text wrapping */
    position: sticky;
    top: 0;
    z-index: 1;
  }

  .data-table td {
    padding: 1rem;
    color: var(--color-text-primary);
    border-bottom: 1px solid var(--color-border-lighter);
    font-size: 0.875rem;
    white-space: nowrap; /* Prevent cell content wrapping */
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 200px; /* Set max width for better control */
  }

  .data-table tbody tr {
    transition: var(--transition-fast);
  }

  .data-table tbody tr:hover {
    background: var(--color-glass-white-light);
  }

  .actions-column {
    width: 120px;
    text-align: center;
  }

  .actions-column .action-buttons {
    justify-content: center;
    gap: 0.5rem;
  }

  /* Responsive */
  @media (max-width: 768px) {
    .table-controls {
      flex-direction: column;
      align-items: stretch;
      gap: 1rem;
    }

    .search-container {
      min-width: auto;
    }

    .action-buttons {
      justify-content: center;
    }
  }
</style>