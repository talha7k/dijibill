<script>
  import { createEventDispatcher } from 'svelte'
  
  export let data = []
  export let columns = []
  export let loading = false
  export let searchTerm = ''
  export let searchPlaceholder = 'Search...'
  export let emptyStateTitle = 'No data found'
  export let emptyStateMessage = 'No items to display'
  export let emptyStateIcon = 'fa-table'
  export let showSearch = true
  export let showActions = true
  export let primaryAction = null
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
          <i class="fas fa-search search-icon"></i>
          <input
            type="text"
            placeholder={searchPlaceholder}
            class="search-input"
            value={searchTerm}
            on:input={handleSearch}
            disabled={loading}
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
        <i class="fas {emptyStateIcon}"></i>
        <h3>{emptyStateTitle}</h3>
        <p>
          {searchTerm ? 'No items match your search criteria.' : emptyStateMessage}
        </p>
        {#if !searchTerm && primaryAction}
          <button class="btn btn-primary" on:click={handlePrimaryAction}>
            {#if primaryAction.icon}
              <i class="fas {primaryAction.icon}"></i>
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
                              <i class="fas {action.icon}"></i>
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
    font-size: 3rem;
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
  }

  .data-table {
    width: 100%;
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
  }

  .data-table td {
    padding: 1rem;
    color: var(--color-text-primary);
    border-bottom: 1px solid var(--color-border-lighter);
    font-size: 0.875rem;
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