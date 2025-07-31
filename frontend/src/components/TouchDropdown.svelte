<script>
  import { createEventDispatcher } from 'svelte';
  
  const dispatch = createEventDispatcher();
  
  export let options = [];
  export let value = '';
  export let placeholder = 'Select an option...';
  export let searchable = true;
  export let maxHeight = '300px';
  export let disabled = false;
  export let required = false;
  export let label = '';
  export let labelArabic = '';
  export let showIcons = false;
  export let columns = 2; // Number of columns to display
  
  let isOpen = false;
  let searchTerm = '';
  let dropdownElement;
  let searchInput;
  
  // Filter options based on search term
  $: filteredOptions = searchable && searchTerm 
    ? options.filter(option => 
        option.label?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        option.labelArabic?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        option.description?.toLowerCase().includes(searchTerm.toLowerCase())
      )
    : options;
  
  // Get selected option for display
  $: selectedOption = options.find(option => option.value === value);
  
  function toggleDropdown() {
    if (disabled) return;
    isOpen = !isOpen;
    if (isOpen && searchable) {
      setTimeout(() => searchInput?.focus(), 100);
    }
  }
  
  function selectOption(option) {
    value = option.value;
    isOpen = false;
    searchTerm = '';
    dispatch('change', { value: option.value, option });
  }
  
  function handleKeydown(event) {
    if (event.key === 'Escape') {
      isOpen = false;
      searchTerm = '';
    }
  }
  
  function handleClickOutside(event) {
    if (dropdownElement && !dropdownElement.contains(event.target)) {
      isOpen = false;
      searchTerm = '';
    }
  }
  
  // Close dropdown when clicking outside
  $: if (typeof window !== 'undefined') {
    if (isOpen) {
      document.addEventListener('click', handleClickOutside);
    } else {
      document.removeEventListener('click', handleClickOutside);
    }
  }
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="touch-dropdown-container" bind:this={dropdownElement}>
  <!-- Label -->
  {#if label || labelArabic}
    <div class="label-container mb-2">
      {#if label}
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
          {label}
          {#if required}<span class="text-red-500">*</span>{/if}
        </label>
      {/if}
      {#if labelArabic}
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 text-right" dir="rtl">
          {labelArabic}
          {#if required}<span class="text-red-500">*</span>{/if}
        </label>
      {/if}
    </div>
  {/if}
  
  <!-- Dropdown Button -->
  <div class="relative">
    <button
      type="button"
      class="touch-dropdown-button"
      class:disabled
      class:open={isOpen}
      on:click={toggleDropdown}
      {disabled}
    >
      <div class="flex items-center justify-between w-full">
        <div class="flex items-center space-x-3">
          {#if showIcons && selectedOption?.icon}
            <span class="text-lg">{selectedOption.icon}</span>
          {/if}
          <div class="text-left">
            {#if selectedOption}
              <div class="font-medium text-gray-900 dark:text-white">
                {selectedOption.label}
              </div>
              {#if selectedOption.labelArabic}
                <div class="text-sm text-gray-500 dark:text-gray-400" dir="rtl">
                  {selectedOption.labelArabic}
                </div>
              {/if}
              {#if selectedOption.description}
                <div class="text-xs text-gray-400 dark:text-gray-500">
                  {selectedOption.description}
                </div>
              {/if}
            {:else}
              <span class="text-gray-500 dark:text-gray-400">{placeholder}</span>
            {/if}
          </div>
        </div>
        <i class="fas fa-chevron-down dropdown-arrow" class:rotated={isOpen}></i>
      </div>
    </button>
    
    <!-- Dropdown Menu -->
    {#if isOpen}
      <div class="touch-dropdown-menu" style="max-height: {maxHeight}">
        <!-- Search Input -->
        {#if searchable}
          <div class="search-container">
            <input
              bind:this={searchInput}
              bind:value={searchTerm}
              type="text"
              placeholder="Search..."
              class="search-input"
            />
            <i class="fas fa-search search-icon"></i>
          </div>
        {/if}
        
        <!-- Options Grid -->
        <div class="options-grid" style="grid-template-columns: repeat({columns}, 1fr)">
          {#each filteredOptions as option (option.value)}
            <button
              type="button"
              class="option-item"
              class:selected={option.value === value}
              on:click={() => selectOption(option)}
            >
              <div class="option-content">
                {#if showIcons && option.icon}
                  <span class="option-icon">{option.icon}</span>
                {/if}
                <div class="option-text">
                  <div class="option-label">{option.label}</div>
                  {#if option.labelArabic}
                    <div class="option-label-arabic" dir="rtl">{option.labelArabic}</div>
                  {/if}
                  {#if option.description}
                    <div class="option-description">{option.description}</div>
                  {/if}
                </div>
              </div>
            </button>
          {/each}
        </div>
        
        <!-- No Results -->
        {#if filteredOptions.length === 0}
          <div class="no-results">
            <i class="fas fa-search-minus no-results-icon"></i>
            <p>No options found</p>
          </div>
        {/if}
      </div>
    {/if}
  </div>
</div>

<style>
  .touch-dropdown-container {
    position: relative;
    width: 100%;
  }
  
  .touch-dropdown-button {
    width: 100%;
    padding: 0.75rem 1rem;
    text-align: left;
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 0.5rem;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
    backdrop-filter: blur(16px);
    transition: all 200ms ease-in-out;
    color: white;
    min-height: 60px;
    cursor: pointer;
  }
  
  .touch-dropdown-button:hover {
    border-color: var(--color-primary-light);
    background: rgba(255, 255, 255, 0.15);
  }
  
  .touch-dropdown-button:focus {
    outline: none;
    border-color: var(--color-primary);
    box-shadow: 0 0 0 2px var(--color-primary-200);
  }
  
  .touch-dropdown-button.disabled {
    background: rgba(255, 255, 255, 0.5);
    color: rgba(255, 255, 255, 0.4);
    cursor: not-allowed;
  }
  
  .touch-dropdown-button.open {
    border-color: var(--color-primary);
    box-shadow: 0 0 0 2px var(--color-primary-200);
  }
  
  .dropdown-arrow {
    width: 1.25rem;
    height: 1.25rem;
    color: rgba(255, 255, 255, 0.8);
    transition: transform 200ms ease-in-out;
    flex-shrink: 0;
  }
  
  .dropdown-arrow.rotated {
    transform: rotate(180deg);
  }
  
  .touch-dropdown-menu {
    position: absolute;
    z-index: 50;
    width: 100%;
    margin-top: 0.25rem;
    background: rgba(67, 67, 67, 0.8);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 0.5rem;
    box-shadow: var(--shadow-glass);
    backdrop-filter: blur(16px);
    overflow: hidden;
    min-width: 100%;
  }
  
  .search-container {
    position: relative;
    padding: 0.75rem;
    border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  }
  
  .search-input {
    width: 100%;
    padding: 0.5rem 0.75rem 0.5rem 2.5rem;
    border: 1px solid rgba(255, 255, 255, 0.3);
    border-radius: 0.375rem;
    background: rgba(255, 255, 255, 0.1);
    color: white;
    backdrop-filter: blur(8px);
    transition: all 200ms ease-in-out;
  }
  
  .search-input::placeholder {
    color: rgba(255, 255, 255, 0.5);
  }
  
  .search-input:focus {
    outline: none;
    border-color: var(--color-primary);
    box-shadow: 0 0 0 2px var(--color-primary-200);
  }
  
  .search-icon {
    position: absolute;
    left: 1.5rem;
    top: 50%;
    transform: translateY(-50%);
    width: 1rem;
    height: 1rem;
    color: rgba(255, 255, 255, 0.6);
  }
  
  .options-grid {
    display: grid;
    gap: 0.45rem;
    padding: 0.5rem;
    max-height: 16rem;
    overflow-y: auto;
  }
  
  .option-item {
    width: 100%;
    padding: 0.75rem;
    text-align: left;
    border-radius: 0.375rem;
    transition: all 150ms ease-in-out;
    background: transparent;
    color: white;
    cursor: pointer;
    min-height: 70px;
    box-shadow: 0 0 0 1px var(--color-primary-200);

  }
  
  .option-item:hover {
    background: rgba(255, 255, 255, 0.1);
  }
  
  .option-item:focus {
    outline: none;
    background: rgba(255, 255, 255, 0.1);
  }
  
  .option-item.selected {
    background: var(--color-primary-200);
    color: var(--color-primary);
  }
  
  .option-content {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }
  
  .option-icon {
    font-size: 1.125rem;
    flex-shrink: 0;
  }
  
  .option-text {
    flex: 1;
    min-width: 0;
  }
  
  .option-label {
    font-weight: 500;
    color: white;
    font-size: 0.875rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .option-label-arabic {
    font-size: 0.75rem;
    color: rgba(255, 255, 255, 0.7);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .option-description {
    font-size: 0.75rem;
    color: rgba(255, 255, 255, 0.6);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .no-results {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 2rem;
    color: rgba(255, 255, 255, 0.6);
  }
  
  .no-results-icon {
    width: 2rem;
    height: 2rem;
    margin-bottom: 0.5rem;
  }
  
  /* Responsive adjustments */
  @media (max-width: 640px) {
    .options-grid {
      grid-template-columns: 1fr !important;
    }
  }
</style>