<script>
  export let label = ''
  export let labelArabic = ''
  export let type = 'text'
  export let value = ''
  export let placeholder = ''
  export let required = false
  export let disabled = false
  export let dir = 'ltr'
  export let id = ''
  export let name = ''
  export let min = undefined
  export let max = undefined
  export let step = undefined
  export let options = [] // For select fields
  export let rows = 3 // For textarea
  export let error = ''

  // Generate unique ID if not provided
  const fieldId = id || `field-${Math.random().toString(36).substr(2, 9)}`

  function handleInput(event) {
    value = event.target.value
  }
</script>

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
    <select 
      {id}
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
  {:else if type === 'textarea'}
    <textarea
      {id}
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
        {id}
        {name}
        type="checkbox"
        class="checkbox-primary"
        bind:checked={value}
        {required}
        {disabled}
      />
      <span>{placeholder}</span>
    </label>
  {:else}
    <input 
      {id}
      {name}
      {type}
      class="input-glass {error ? 'border-error' : ''}"
      bind:value
      {placeholder}
      {required}
      {disabled}
      {dir}
      {min}
      {max}
      {step}
      on:input={handleInput}
    />
  {/if}
  
  {#if error}
    <div class="text-error text-sm mt-1">{error}</div>
  {/if}
</div>