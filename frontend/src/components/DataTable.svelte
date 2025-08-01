<script>
  import { createEventDispatcher } from 'svelte'
  import FormField from './FormField.svelte'
  import ActionButton from './ActionButton.svelte'
  import IconButton from './IconButton.svelte'
  import { getFaIcon } from '../utils/iconUtils.js'
  
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
  export let emptyStateIcon = 'table'
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
            <ActionButton
              variant="primary"
              icon={primaryAction.icon}
              text={primaryAction.text || 'Action'}
              {loading}
              on:click={handlePrimaryAction}
            />
          {/if}
          
          {#each secondaryActions as action}
            <ActionButton
              variant="outline"
              icon={action.icon}
              text={action.text || 'Action'}
              disabled={loading}
              on:click={() => handleSecondaryAction(action)}
            />
          {/each}
        </div>
      {/if}

      {#if showSearch}
        <div class="search-container">
          <div class="search-icon">
            <i class="fas fa-search"></i>
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
        <i class="fas fa-spinner fa-spin"></i>
        <span>Loading...</span>
      </div>
    {:else if data.length === 0}
      <div class="empty-state">
        <i class="{getFaIcon(emptyStateIcon)} text-4xl mb-4 text-gray-400"></i>
        <h3>{emptyStateTitle}</h3>
        <p>
          {searchTerm ? 'No items match your search criteria.' : emptyStateMessage}
        </p>
        {#if !searchTerm && primaryAction}
          <ActionButton
            variant="primary"
            icon={primaryAction.icon}
            text={primaryAction.text}
            on:click={handlePrimaryAction}
          />
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
                          {#if action.icon}
                            <IconButton
                              icon={action.icon}
                              variant={action.key === 'delete' ? 'danger' : 'secondary'}
                              title={action.title || action.text}
                              disabled={loading}
                              on:click={() => handleRowAction(action.key, item)}
                            />
                          {:else}
                            <ActionButton
                              variant="ghost"
                              size="sm"
                              text={action.text}
                              title={action.title || action.text}
                              disabled={loading}
                              customClass={action.class || ''}
                              on:click={() => handleRowAction(action.key, item)}
                            />
                          {/if}
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

  .loading-state i {
    font-size: 1.25rem;
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 3rem;
    text-align: center;
  }

  .empty-state i {
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