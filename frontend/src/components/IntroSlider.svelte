<script>
  import { createEventDispatcher } from 'svelte';
  import { fade, fly } from 'svelte/transition';
  import { quintOut } from 'svelte/easing';
  import introSlidesData from '../data/introSlides.json';

  const dispatch = createEventDispatcher();

  export let show = false;

  let currentSlide = 0;
  let slides = introSlidesData.slides;
  let isTransitioning = false;

  function nextSlide() {
    if (isTransitioning) return;
    if (currentSlide < slides.length - 1) {
      isTransitioning = true;
      currentSlide++;
      setTimeout(() => isTransitioning = false, 500);
    } else {
      completeIntro();
    }
  }

  function prevSlide() {
    if (isTransitioning) return;
    if (currentSlide > 0) {
      isTransitioning = true;
      currentSlide--;
      setTimeout(() => isTransitioning = false, 500);
    }
  }

  function goToSlide(index) {
    if (isTransitioning || index === currentSlide) return;
    isTransitioning = true;
    currentSlide = index;
    setTimeout(() => isTransitioning = false, 500);
  }

  function skipIntro() {
    completeIntro();
  }

  function completeIntro() {
    dispatch('complete');
  }

  function getColorClasses(color) {
    const colorMap = {
      blue: 'from-blue-500 to-blue-600 text-blue-600 bg-blue-50',
      green: 'from-green-500 to-green-600 text-green-600 bg-green-50',
      purple: 'from-purple-500 to-purple-600 text-purple-600 bg-purple-50',
      indigo: 'from-indigo-500 to-indigo-600 text-indigo-600 bg-indigo-50',
      pink: 'from-pink-500 to-pink-600 text-pink-600 bg-pink-50',
      orange: 'from-orange-500 to-orange-600 text-orange-600 bg-orange-50',
      teal: 'from-teal-500 to-teal-600 text-teal-600 bg-teal-50',
      cyan: 'from-cyan-500 to-cyan-600 text-cyan-600 bg-cyan-600 bg-cyan-50',
      emerald: 'from-emerald-500 to-emerald-600 text-emerald-600 bg-emerald-50',
      yellow: 'from-yellow-500 to-yellow-600 text-yellow-600 bg-yellow-50',
      red: 'from-red-500 to-red-600 text-red-600 bg-red-50',
      gray: 'from-gray-500 to-gray-600 text-gray-600 bg-gray-50'
    };
    return colorMap[color] || colorMap.blue;
  }

  function getIconSvg(iconClass) {
    // Convert FontAwesome classes to SVG paths
    const iconMap = {
      'fas fa-rocket': 'M12 2l3.09 6.26L22 9l-5.91 5.74L17.82 21 12 18.27 6.18 21l1.73-6.26L2 9l6.91-1.74L12 2z',
      'fas fa-chart-line': 'M3 17l6-6 4 4 8-8M21 7l-5 5-4-4-6 6',
      'fas fa-file-invoice-dollar': 'M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6zM14 8V3l5 5h-5zM8 12h8M8 16h8M8 20h8',
      'fas fa-cash-register': 'M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2 2zM3 7a2 2 0 012-2h14a2 2 0 012 2v2H3V7z',
      'fas fa-users': 'M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2M23 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75M13 7a4 4 0 1 1-8 0 4 4 0 0 1 8 0z',
      'fas fa-truck': 'M1 3h15v13H1zM16 8h4l3 3v5h-7V8zM6.5 18a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3zM17.5 18a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3z',
      'fas fa-boxes': 'M12 2l3.09 6.26L22 9l-5.91 5.74L17.82 21 12 18.27 6.18 21l1.73-6.26L2 9l6.91-1.74L12 2z',
      'fas fa-shopping-bag': 'M19 7h-3V6a4 4 0 0 0-8 0v1H5a1 1 0 0 0-1 1v11a3 3 0 0 0 3 3h10a3 3 0 0 0 3-3V8a1 1 0 0 0-1-1zM10 6a2 2 0 0 1 4 0v1h-4V6z',
      'fas fa-file-invoice': 'M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6zM14 8V3l5 5h-5z',
      'fas fa-credit-card': 'M21 4H3a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h18a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2zM3 10h18v8H3v-8z',
      'fas fa-user-cog': 'M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z',
      'fas fa-cog': 'M12 15.5A3.5 3.5 0 0 1 8.5 12A3.5 3.5 0 0 1 12 8.5a3.5 3.5 0 0 1 3.5 3.5 3.5 3.5 0 0 1-3.5 3.5m7.43-2.53c.04-.32.07-.64.07-.97 0-.33-.03-.66-.07-1l2.11-1.63c.19-.15.24-.42.12-.64l-2-3.46c-.12-.22-.39-.31-.61-.22l-2.49 1c-.52-.39-1.06-.73-1.69-.98l-.37-2.65A.506.506 0 0 0 14 2h-4c-.25 0-.46.18-.5.42l-.37 2.65c-.63.25-1.17.59-1.69.98l-2.49-1c-.22-.09-.49 0-.61.22l-2 3.46c-.13.22-.07.49.12.64L4.57 11c-.04.34-.07.67-.07 1 0 .33.03.65.07.97l-2.11 1.66c-.19.15-.25.42-.12.64l2 3.46c.12.22.39.3.61.22l2.49-1.01c.52.4 1.06.74 1.69.99l.37 2.65c.04.24.25.42.5.42h4c.25 0 .46-.18.5-.42l.37-2.65c.63-.26 1.17-.59 1.69-.99l2.49 1.01c.22.08.49 0 .61-.22l2-3.46c.12-.22.07-.49-.12-.64l-2.11-1.66Z',
      'fas fa-check-circle': 'M22 11.08V12a10 10 0 1 1-5.93-9.14M22 4L12 14.01l-3-3'
    };
    return iconMap[iconClass] || 'M12 2l3.09 6.26L22 9l-5.91 5.74L17.82 21 12 18.27 6.18 21l1.73-6.26L2 9l6.91-1.74L12 2z';
  }

  $: currentSlideData = slides[currentSlide];
  $: colorClasses = getColorClasses(currentSlideData?.color);
</script>

{#if show}
  <div 
    class="fixed inset-0 bg-gradient-to-br from-indigo-900 via-purple-900 to-pink-900 z-50 flex items-center justify-center overflow-hidden"
    transition:fade={{ duration: 600, easing: quintOut }}
  >
    <!-- Animated background elements -->
    <div class="absolute inset-0 overflow-hidden">
      <div class="absolute -top-40 -right-40 w-80 h-80 bg-purple-500 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-pulse"></div>
      <div class="absolute -bottom-40 -left-40 w-80 h-80 bg-pink-500 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-pulse" style="animation-delay: 2s;"></div>
      <div class="absolute top-40 left-40 w-80 h-80 bg-indigo-500 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-pulse" style="animation-delay: 4s;"></div>
    </div>
    
    <div 
      class="relative bg-white/10 backdrop-blur-xl border border-white/20 rounded-3xl shadow-2xl w-full max-w-6xl mx-8 overflow-hidden"
      transition:fly={{ y: 50, duration: 800, delay: 300, easing: quintOut }}
    >
      <!-- Header -->
      <div class="flex items-center justify-between p-8 border-b border-white/20">
        <div class="flex items-center space-x-4">
          <div class="w-12 h-12 bg-gradient-to-r {colorClasses.split(' ')[0]} {colorClasses.split(' ')[1]} rounded-2xl flex items-center justify-center shadow-lg">
            <span class="text-white font-bold text-lg">{currentSlide + 1}</span>
          </div>
          <div>
            <h2 class="text-2xl font-bold text-white">Getting Started</h2>
            <p class="text-white/70 text-sm">Welcome to DijiBill</p>
          </div>
        </div>
        <button
          on:click={skipIntro}
          class="text-white/60 hover:text-white transition-colors duration-300 p-2 rounded-full hover:bg-white/10"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>

      <!-- Content -->
      <div class="p-8">
        {#key currentSlide}
          <div 
            class="grid grid-cols-1 lg:grid-cols-2 gap-8 items-center"
            in:fly={{ x: 100, duration: 500, delay: 100, easing: quintOut }}
            out:fly={{ x: -100, duration: 300, easing: quintOut }}
          >
            <!-- Left side - Main content -->
            <div class="text-center lg:text-left">
              <!-- Icon -->
              <div class="mb-6 flex justify-center lg:justify-start">
                <div class="w-24 h-24 bg-white/20 backdrop-blur-sm rounded-2xl flex items-center justify-center shadow-xl border border-white/30">
                  <svg class="w-12 h-12 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="{getIconSvg(currentSlideData.icon)}"></path>
                  </svg>
                </div>
              </div>

              <!-- Title -->
              <h3 class="text-3xl lg:text-4xl font-bold text-white mb-4 leading-tight">
                {currentSlideData.title}
              </h3>

              <!-- Description -->
              <p class="text-white/80 text-lg leading-relaxed">
                {currentSlideData.description}
              </p>
            </div>

            <!-- Right side - Features -->
            <div class="space-y-4">
              <h4 class="text-xl font-semibold text-white mb-6 text-center lg:text-left">Key Features</h4>
              <div class="grid gap-4">
                {#each currentSlideData.features as feature, index}
                  <div 
                    class="bg-white/10 backdrop-blur-sm rounded-xl p-4 border border-white/20 hover:bg-white/15 transition-all duration-300"
                    style="animation-delay: {index * 100}ms"
                    in:fly={{ x: 50, duration: 400, delay: 200 + (index * 100), easing: quintOut }}
                  >
                    <div class="flex items-start space-x-3">
                      <div class="w-8 h-8 bg-white/20 rounded-lg flex items-center justify-center flex-shrink-0 mt-1">
                        <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                        </svg>
                      </div>
                      <div>
                        <h5 class="font-semibold text-white text-sm mb-1">{feature.title}</h5>
                        <p class="text-white/70 text-xs leading-relaxed">{feature.description}</p>
                      </div>
                    </div>
                  </div>
                {/each}
              </div>
            </div>
          </div>
        {/key}
      </div>

      <!-- Progress Indicators -->
      <div class="px-8 pb-6">
        <div class="flex justify-center space-x-3">
          {#each slides as slide, index}
            <button
              on:click={() => goToSlide(index)}
              class="w-3 h-3 rounded-full transition-all duration-300 {index === currentSlide ? 'bg-white scale-125 shadow-lg' : 'bg-white/30 hover:bg-white/50'}"
              disabled={isTransitioning}
            ></button>
          {/each}
        </div>
      </div>

      <!-- Navigation -->
      <div class="flex items-center justify-between p-8 bg-white/5 backdrop-blur-sm border-t border-white/20">
        <button
          on:click={prevSlide}
          disabled={currentSlide === 0 || isTransitioning}
          class="flex items-center space-x-3 px-6 py-3 text-white/70 hover:text-white disabled:opacity-30 disabled:cursor-not-allowed transition-all duration-300 rounded-xl hover:bg-white/10"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
          </svg>
          <span class="font-medium">Previous</span>
        </button>

        <div class="text-white/60 font-medium">
          {currentSlide + 1} of {slides.length}
        </div>

        <button
          on:click={nextSlide}
          disabled={isTransitioning}
          class="flex items-center space-x-3 px-8 py-3 bg-white text-gray-900 rounded-xl hover:bg-white/90 transition-all duration-300 shadow-lg hover:shadow-xl transform hover:scale-105 font-medium disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
        >
          <span>{currentSlide === slides.length - 1 ? 'Get Started' : 'Next'}</span>
          {#if currentSlide === slides.length - 1}
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
          {:else}
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
            </svg>
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  /* Custom scrollbar for content if needed */
  :global(.intro-slider-content::-webkit-scrollbar) {
    width: 6px;
  }
  
  :global(.intro-slider-content::-webkit-scrollbar-track) {
    background: #f1f1f1;
    border-radius: 3px;
  }
  
  :global(.intro-slider-content::-webkit-scrollbar-thumb) {
    background: #c1c1c1;
    border-radius: 3px;
  }
  
  :global(.intro-slider-content::-webkit-scrollbar-thumb:hover) {
    background: #a8a8a8;
  }
</style>