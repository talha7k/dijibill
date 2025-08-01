<script>
  import { createEventDispatcher } from 'svelte'
  import { getFaIcon } from '../utils/iconUtils.js'

  export let icon = ''
  export let variant = 'secondary' // 'primary', 'secondary', 'danger', 'ghost'
  export let size = 'sm' // 'xs', 'sm', 'md'
  export let disabled = false
  export let loading = false
  export let title = ''
  export let customClass = ''

  const dispatch = createEventDispatcher()

  function handleClick() {
    if (!disabled && !loading) {
      dispatch('click')
    }
  }

  $: buttonClass = [
    'icon-btn',
    `icon-btn-${variant}`,
    `icon-btn-${size}`,
    customClass,
    disabled || loading ? 'icon-btn-disabled' : ''
  ].filter(Boolean).join(' ')
</script>

<button 
  class={buttonClass}
  on:click={handleClick}
  {disabled}
  {title}
>
  {#if loading}
    <i class="fas fa-spinner fa-spin"></i>
  {:else}
    <i class="{getFaIcon(icon)}"></i>
  {/if}
</button>

<style>
  .icon-btn {
    @apply inline-flex items-center justify-center rounded-lg font-medium transition-all duration-200 cursor-pointer;
    border: 1px solid transparent;
  }

  .icon-btn-xs {
    @apply w-6 h-6 text-xs;
  }

  .icon-btn-sm {
    @apply w-8 h-8 text-sm;
  }

  .icon-btn-md {
    @apply w-10 h-10 text-base;
  }

  .icon-btn-primary {
    @apply bg-blue-600 text-white border-blue-600;
  }

  .icon-btn-primary:hover:not(.icon-btn-disabled) {
    @apply bg-blue-700 border-blue-700;
  }

  .icon-btn-secondary {
    @apply bg-gray-100 text-gray-600 border-gray-200;
  }

  .icon-btn-secondary:hover:not(.icon-btn-disabled) {
    @apply bg-gray-200 text-gray-800 border-gray-300;
  }

  .icon-btn-danger {
    @apply bg-red-100 text-red-600 border-red-200;
  }

  .icon-btn-danger:hover:not(.icon-btn-disabled) {
    @apply bg-red-200 text-red-800 border-red-300;
  }

  .icon-btn-ghost {
    @apply bg-transparent text-gray-500 border-transparent;
  }

  .icon-btn-ghost:hover:not(.icon-btn-disabled) {
    @apply bg-gray-100 text-gray-700;
  }

  .icon-btn-disabled {
    @apply opacity-50 cursor-not-allowed;
  }
</style>