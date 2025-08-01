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

{#if isVisible}
  {#if asOverlay}
    <!-- Backdrop -->
    <div class="fixed inset-0 bg-black/50 z-[999998] transition-opacity duration-300"
      on:click={closeSidebar}
      on:keydown={(e) => e.key === 'Escape' && closeSidebar()}
      role="button"
      tabindex="0"
      aria-label="Close sidebar"
    ></div>
  {/if}
  
  <!-- Sidebar -->
  <div class="{asOverlay ? 'fixed left-0 top-0 h-full w-64 z-[999999]' : 'w-64 h-full'} bg-base-100/10 backdrop-blur-lg border-r border-white/20 flex flex-col transition-transform duration-300">
    <!-- Logo/Header -->
    <div class="p-4 border-b border-white/20">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
          <i class="fas fa-file-invoice w-5 h-5 text-white"></i>
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
            <i class="fas fa-tachometer-alt w-5 h-5"></i>
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
                    <i class="fas fa-{item.icon} w-4 h-4"></i>
                    <span class="text-sm">{item.label}</span>
                    {#if item.id === 'payments'}
                      <div class="ml-auto">
                        <div class="w-4 h-4 bg-accent rounded-full flex items-center justify-center">
                          <i class="fas fa-circle w-2 h-2 text-white"></i>
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