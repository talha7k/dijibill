<script>
  import { onMount, createEventDispatcher } from 'svelte'

  // Props
  export let tabs = []
  export let activeTab = ''
  export let variant = 'standard' // 'standard' or 'glass'
  export let showScrollButtons = true
  export let scrollAmount = 200

  // Event dispatcher
  const dispatch = createEventDispatcher()

  // Internal state
  let tabsContainer
  let canScrollLeft = false
  let canScrollRight = false

  // Tab scrolling functions
  function scrollLeft() {
    if (tabsContainer) {
      tabsContainer.scrollBy({ left: -scrollAmount, behavior: 'smooth' })
    }
  }

  function scrollRight() {
    if (tabsContainer) {
      tabsContainer.scrollBy({ left: scrollAmount, behavior: 'smooth' })
    }
  }

  // Check if scrolling is needed
  function updateScrollButtons() {
    if (tabsContainer) {
      canScrollLeft = tabsContainer.scrollLeft > 0
      canScrollRight = tabsContainer.scrollLeft < (tabsContainer.scrollWidth - tabsContainer.clientWidth)
    }
  }

  // Handle tab click
  function handleTabClick(tabId) {
    dispatch('tabChange', { tabId })
  }

  // Setup scroll listener on mount
  onMount(() => {
    if (tabsContainer) {
      tabsContainer.addEventListener('scroll', updateScrollButtons)
      // Initial check
      setTimeout(updateScrollButtons, 100)
    }

    return () => {
      if (tabsContainer) {
        tabsContainer.removeEventListener('scroll', updateScrollButtons)
      }
    }
  })

  // Reactive statement to update scroll buttons when tabs change
  $: if (tabsContainer && tabs.length > 0) {
    setTimeout(updateScrollButtons, 100)
  }
</script>

{#if variant === 'glass'}
  <!-- Glass variant (for Products page) -->
  <div class="tabs-glass">
    {#each tabs as tab}
      <button 
        class="tab-glass {activeTab === tab.id ? 'tab-glass-active' : ''}"
        on:click={() => handleTabClick(tab.id)}
      >
        {#if tab.icon}
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={tab.icon}></path>
          </svg>
        {/if}
        {tab.label || tab.name}
      </button>
    {/each}
  </div>
{:else}
  <!-- Standard variant (for General Settings page) -->
  <div class="bg-white/10 backdrop-blur-lg border-b border-white/20">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="relative flex items-center">
        <!-- Left scroll button -->
        {#if showScrollButtons}
          <button
            on:click={scrollLeft}
            class="flex-shrink-0 p-2 text-white/60 hover:text-white disabled:opacity-50 disabled:cursor-not-allowed"
            disabled={!canScrollLeft}
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
          </button>
        {/if}

        <!-- Scrollable tabs container -->
        <div 
          bind:this={tabsContainer}
          class="flex-1 overflow-x-auto scrollbar-hide"
          style="scrollbar-width: none; -ms-overflow-style: none;"
          on:scroll={updateScrollButtons}
        >
          <nav class="flex space-x-2 py-4" style="min-width: max-content;">
            {#each tabs as tab}
              <button
                on:click={() => handleTabClick(tab.id)}
                class="tab-standard {activeTab === tab.id ? 'tab-standard-active' : ''}"
              >
                {#if tab.icon}
                  <span class="mr-2 text-lg">{tab.icon}</span>
                {/if}
                {tab.label || tab.name}
              </button>
            {/each}
          </nav>
        </div>

        <!-- Right scroll button -->
        {#if showScrollButtons}
          <button
            on:click={scrollRight}
            class="flex-shrink-0 p-2 text-white/60 hover:text-white disabled:opacity-50 disabled:cursor-not-allowed"
            disabled={!canScrollRight}
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>
        {/if}
      </div>
    </div>
  </div>
{/if}

<style>
  /* Hide scrollbar for Chrome, Safari and Opera */
  .scrollbar-hide::-webkit-scrollbar {
    display: none;
  }

  /* Hide scrollbar for IE, Edge and Firefox */
  .scrollbar-hide {
    -ms-overflow-style: none;  /* IE and Edge */
    scrollbar-width: none;  /* Firefox */
  }
</style>