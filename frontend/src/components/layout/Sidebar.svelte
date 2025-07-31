<script>
  import { createEventDispatcher } from 'svelte';
  
  export let isVisible = false;
  export let asOverlay = false;
  export let currentView = 'dashboard';
  export let menuItems = [];
  
  const dispatch = createEventDispatcher();
  
  function switchView(view) {
    dispatch('view-change', view);
  }
  
  function closeSidebar() {
    dispatch('close');
  }
</script>

{#if isVisible && asOverlay}
  <!-- Backdrop -->
  <div 
    class="fixed inset-0 bg-black/50 z-40 transition-opacity duration-300"
    on:click={closeSidebar}
    on:keydown={(e) => e.key === 'Escape' && closeSidebar()}
    role="button"
    tabindex="0"
    aria-label="Close sidebar"
  ></div>
  
  <!-- Sidebar as overlay -->
  <div class="fixed left-0 top-0 h-full w-64 bg-base-100/10 backdrop-blur-lg border-r border-white/20 flex flex-col z-50 transition-transform duration-300">
    <!-- Logo/Header -->
    <div class="p-4 border-b border-white/20">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
          <svg class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 24 24">
            <path d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
          </svg>
        </div>
        <div>
          <h1 class="text-lg font-bold text-white">E-Invoice</h1>
          <p class="text-xs text-white/60">الفاتورة الإلكترونية</p>
        </div>
      </div>
    </div>

    <!-- Navigation Menu -->
    <nav class="flex-1 p-4 overflow-y-auto">
      <ul class="space-y-2">
        <!-- Dashboard (standalone) -->
        <li>
          <button
            class="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-left transition-colors {currentView === 'dashboard' ? 'bg-primary text-white' : 'text-white/80 hover:bg-white/10'}"
            on:click={() => switchView('dashboard')}
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="{menuItems[0]?.icon || 'M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2 2v-2m0 0V5a2 2 0 012-2h6l2 2h6a2 2 0 012 2v2M7 13h10M7 17h4'}"></path>
            </svg>
            <span class="text-sm font-medium">{menuItems[0]?.label || 'Dashboard'}</span>
          </button>
        </li>

        <!-- Categories -->
        {#each menuItems.slice(1) as category}
          <li class="mt-6">
            <div class="text-xs font-semibold text-white/60 uppercase tracking-wider mb-2 px-3">
              {category.category}
            </div>
            <ul class="space-y-1">
              {#each category.items as item}
                <li>
                  <button
                    class="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-left transition-colors {currentView === item.id ? 'bg-primary text-white' : 'text-white/80 hover:bg-white/10'}"
                    on:click={() => switchView(item.id)}
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="{item.icon}"></path>
                    </svg>
                    <span class="text-sm">{item.label}</span>
                    {#if item.id === 'payments'}
                      <div class="ml-auto">
                        <div class="w-4 h-4 bg-accent rounded-full flex items-center justify-center">
                          <svg class="w-2 h-2 text-white" fill="currentColor" viewBox="0 0 8 8">
                            <circle cx="4" cy="4" r="3"/>
                          </svg>
                        </div>
                      </div>
                    {/if}
                  </button>
                </li>
              {/each}
            </ul>
          </li>
        {/each}
      </ul>
    </nav>
  </div>
{/if}