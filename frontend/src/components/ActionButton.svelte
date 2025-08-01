<script>
  import { createEventDispatcher } from 'svelte'
  import { getFaIcon } from '../utils/iconUtils.js'

  export let variant = 'primary' // 'primary', 'secondary', 'outline', 'ghost', 'danger'
  export let size = 'md' // 'sm', 'md', 'lg'
  export let icon = null
  export let text = ''
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
    'btn',
    `btn-${variant}`,
    `btn-${size}`,
    customClass,
    disabled || loading ? 'btn-disabled' : ''
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
  {:else if icon}
    <i class="{getFaIcon(icon)}"></i>
  {/if}
  {#if text}
    <span>{text}</span>
  {/if}
  <slot />
</button>

<style>
  .btn {
    @apply inline-flex items-center justify-center gap-2 px-4 py-2 rounded-lg font-medium transition-all duration-200 cursor-pointer;
    border: 1px solid transparent;
  }

  .btn-sm {
    @apply px-3 py-1.5 text-sm;
  }

  .btn-md {
    @apply px-4 py-2 text-sm;
  }

  .btn-lg {
    @apply px-6 py-3 text-base;
  }

  .btn-primary {
    @apply bg-purple-500 text-white border-purple-600;
  }

  .btn-primary:hover:not(.btn-disabled) {
    @apply bg-purple-700 border-purple-700;
  }

  .btn-secondary {
    @apply bg-gray-600 text-white border-gray-600;
  }

  .btn-secondary:hover:not(.btn-disabled) {
    @apply bg-gray-700 border-gray-700;
  }

  .btn-outline {
    @apply bg-transparent text-purple-400 border-purple-400;
  }

  .btn-outline:hover:not(.btn-disabled) {
    @apply bg-purple-400 text-white;
  }

  .btn-ghost {
    @apply bg-transparent text-purple-200 border-transparent;
  }

  .btn-ghost:hover:not(.btn-disabled) {
    @apply bg-purple-200 text-gray-800;
  }

  .btn-danger {
    @apply bg-red-600 text-white border-red-600;
  }

  .btn-danger:hover:not(.btn-disabled) {
    @apply bg-red-700 border-red-700;
  }

  .btn-disabled {
    @apply opacity-50 cursor-not-allowed;
  }
</style>