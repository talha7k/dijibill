<script>
  export let label = ''
  export let labelArabic = ''
  export let type = 'text'
  export let value = ''
  export let placeholder = ''
  export let required = false
  export let disabled = false
  export let readonly = false
  export let dir = 'ltr'
  export let id = ''
  export let name = ''
  export let min = undefined
  export let max = undefined
  export let step = undefined
  export let options = [] // For select fields
  export let rows = 3 // For textarea
  export let error = ''
  export let checked = false // For checkbox fields
  export let customClass = '' // Custom CSS class for input styling
  
  // Touch-friendly dropdown options
  export let touchFriendly = false // Enable touch-friendly dropdown mode
  export let searchable = false // Enable search in dropdown
  export let multiColumn = false // Enable two-column layout for options
  export let maxHeight = '300px' // Maximum height of dropdown
  export let showIcons = false // Show icons in options (options should have icon property)

  // Generate unique ID if not provided
  const fieldId = id || `field-${Math.random().toString(36).substr(2, 9)}`
  
  // Touch dropdown state
  let isDropdownOpen = false
  let searchTerm = ''
  let dropdownRef
  let inputRef

  // Filtered options for searchable dropdown
  $: filteredOptions = searchable && searchTerm 
    ? options.filter(option => 
        option.label.toLowerCase().includes(searchTerm.toLowerCase()) ||
        (option.labelArabic && option.labelArabic.includes(searchTerm))
      )
    : options

  // Get display text for selected value
  $: selectedOption = options.find(opt => opt.value === value)
  $: displayText = selectedOption ? selectedOption.label : (placeholder || 'Select an option')

  function handleInput(event) {
    value = event.target.value
  }

  function handleCheckboxChange(event) {
    checked = event.target.checked
  }

  function toggleDropdown() {
    if (disabled) return
    isDropdownOpen = !isDropdownOpen
    if (isDropdownOpen && searchable) {
      setTimeout(() => inputRef?.focus(), 100)
    }
  }

  function selectOption(option) {
    value = option.value
    isDropdownOpen = false
    searchTerm = ''
    
    // Dispatch change event
    const event = new CustomEvent('change', { detail: { value: option.value, option } })
    if (dropdownRef) {
      dropdownRef.dispatchEvent(event)
    }
  }

  function handleClickOutside(event) {
    if (dropdownRef && !dropdownRef.contains(event.target)) {
      isDropdownOpen = false
      searchTerm = ''
    }
  }

  function handleKeydown(event) {
    if (event.key === 'Escape') {
      isDropdownOpen = false
      searchTerm = ''
    }
  }
</script>

<svelte:window on:click={handleClickOutside} on:keydown={handleKeydown} />

<div class="form-control">
  <label for={fieldId} class="label-glass">
    {label}
    {#if labelArabic}
      <span class="text-white/60 text-sm">({labelArabic})</span>
    {/if}
    {#if required}
      <span class="text-error ml-1">*</span>
    {/if}
  </label>
  
  {#if type === 'select'}
    {#if touchFriendly}
      <!-- Touch-friendly dropdown -->
      <div class="relative" bind:this={dropdownRef}>
        <button
          type="button"
          class="select-glass w-full text-left flex items-center justify-between {error ? 'border-error' : ''} {disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}"
          on:click={toggleDropdown}
          {disabled}
        >
          <span class="truncate {!selectedOption ? 'text-white/60' : ''}">{displayText}</span>
          <i class="fas fa-chevron-down transition-transform duration-200 {isDropdownOpen ? 'rotate-180' : ''}"></i>
        </button>
        
        {#if isDropdownOpen}
          <div class="absolute z-50 w-full mt-1 bg-gray-800/95 backdrop-blur-sm border border-white/20 rounded-lg shadow-xl" style="max-height: {maxHeight}">
            {#if searchable}
              <div class="p-3 border-b border-white/10">
                <input
                  bind:this={inputRef}
                  type="text"
                  bind:value={searchTerm}
                  placeholder="Search options..."
                  class="w-full px-3 py-2 bg-gray-700/50 border border-white/20 rounded text-white placeholder-white/60 focus:outline-none focus:border-blue-400"
                />
              </div>
            {/if}
            
            <div class="overflow-y-auto" style="max-height: calc({maxHeight} - {searchable ? '60px' : '0px'})">
              {#if multiColumn && filteredOptions.length > 6}
                <div class="grid grid-cols-2 gap-1 p-2">
                  {#each filteredOptions as option}
                    <button
                      type="button"
                      class="flex items-center gap-2 px-3 py-2 text-left hover:bg-white/10 rounded transition-colors duration-150 {value === option.value ? 'bg-blue-600/30 text-blue-300' : 'text-white'}"
                      on:click={() => selectOption(option)}
                    >
                      {#if showIcons && option.icon}
                        <i class="fas {option.icon} text-sm"></i>
                      {/if}
                      <div class="flex-1 min-w-0">
                        <div class="truncate text-sm">{option.label}</div>
                        {#if option.labelArabic}
                          <div class="truncate text-xs text-white/60" dir="rtl">{option.labelArabic}</div>
                        {/if}
                      </div>
                    </button>
                  {/each}
                </div>
              {:else}
                <div class="py-1">
                  {#each filteredOptions as option}
                    <button
                      type="button"
                      class="flex items-center gap-3 w-full px-4 py-3 text-left hover:bg-white/10 transition-colors duration-150 {value === option.value ? 'bg-blue-600/30 text-blue-300' : 'text-white'}"
                      on:click={() => selectOption(option)}
                    >
                      {#if showIcons && option.icon}
                        <i class="fas {option.icon} text-lg"></i>
                      {/if}
                      <div class="flex-1">
                        <div class="font-medium">{option.label}</div>
                        {#if option.labelArabic}
                          <div class="text-sm text-white/60" dir="rtl">{option.labelArabic}</div>
                        {/if}
                        {#if option.description}
                          <div class="text-xs text-white/50 mt-1">{option.description}</div>
                        {/if}
                      </div>
                      {#if value === option.value}
                        <i class="fas fa-check text-blue-400"></i>
                      {/if}
                    </button>
                  {/each}
                </div>
              {/if}
              
              {#if filteredOptions.length === 0}
                <div class="px-4 py-6 text-center text-white/60">
                  <i class="fas fa-search text-2xl mb-2"></i>
                  <div>No options found</div>
                </div>
              {/if}
            </div>
          </div>
        {/if}
      </div>
    {:else}
      <!-- Standard select dropdown -->
      <select 
        id={fieldId}
        {name}
        class="select-glass {error ? 'border-error' : ''}"
        bind:value
        {required}
        {disabled}
      >
        <option value="">{placeholder || 'Select an option'}</option>
        {#each options as option}
          <option value={option.value}>{option.label}</option>
        {/each}
      </select>
    {/if}
  {:else if type === 'textarea'}
    <textarea
      id={fieldId}
      {name}
      class="textarea-glass {error ? 'border-error' : ''}"
      bind:value
      {placeholder}
      {required}
      {disabled}
      {dir}
      {rows}
      on:input={handleInput}
    ></textarea>
  {:else if type === 'checkbox'}
    <label class="checkbox-glass">
      <input 
        id={fieldId}
        {name}
        type="checkbox"
        class="checkbox-primary"
        bind:checked
        {required}
        {disabled}
        on:change={handleCheckboxChange}
      />
      <span>{placeholder}</span>
    </label>
  {:else if type === 'number'}
    <input 
      id={fieldId}
      {name}
      type="number"
      class="input-glass {error ? 'border-error' : ''}"
      bind:value
      {placeholder}
      {required}
      {disabled}
      {readonly}
      {dir}
      {min}
      {max}
      {step}
      on:input={handleInput}
    />
  {:else if type === 'email'}
    <input 
      id={fieldId}
      {name}
      type="email"
      class="input-glass {error ? 'border-error' : ''}"
      bind:value
      {placeholder}
      {required}
      {disabled}
      {dir}
      on:input={handleInput}
    />
  {:else if type === 'tel'}
    <input 
      id={fieldId}
      {name}
      type="tel"
      class="input-glass {error ? 'border-error' : ''}"
      bind:value
      {placeholder}
      {required}
      {disabled}
      {dir}
      on:input={handleInput}
    />
  {:else if type === 'date'}
    <input 
      id={fieldId}
      {name}
      type="date"
      class="input-glass {error ? 'border-error' : ''}"
      bind:value
      {required}
      {disabled}
      {min}
      {max}
      on:input={handleInput}
    />
  {:else if type === 'datetime-local'}
    <input 
      id={fieldId}
      {name}
      type="datetime-local"
      class="input-glass {error ? 'border-error' : ''}"
      bind:value
      {required}
      {disabled}
      {min}
      {max}
      on:input={handleInput}
    />
  {:else if type === 'time'}
    <input 
      id={fieldId}
      {name}
      type="time"
      class="input-glass {error ? 'border-error' : ''}"
      bind:value
      {required}
      {disabled}
      {min}
      {max}
      on:input={handleInput}
    />
  {:else}
    <input 
      id={fieldId}
      {name}
      type="text"
      class="input-glass {customClass} {error ? 'border-error' : ''}"
      bind:value
      {placeholder}
      {required}
      {disabled}
      {dir}
      on:input={handleInput}
    />
  {/if}
  
  {#if error}
    <div class="text-error text-sm mt-1">{error}</div>
  {/if}
</div>