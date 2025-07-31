<script>
  import { createEventDispatcher } from 'svelte';
  import { fade, fly } from 'svelte/transition';
  import introSlidesData from '../data/introSlides.json';

  const dispatch = createEventDispatcher();

  export let show = false;

  let currentSlide = 0;
  let slides = introSlidesData.slides;

  function nextSlide() {
    if (currentSlide < slides.length - 1) {
      currentSlide++;
    } else {
      completeIntro();
    }
  }

  function prevSlide() {
    if (currentSlide > 0) {
      currentSlide--;
    }
  }

  function goToSlide(index) {
    currentSlide = index;
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
      cyan: 'from-cyan-500 to-cyan-600 text-cyan-600 bg-cyan-50',
      emerald: 'from-emerald-500 to-emerald-600 text-emerald-600 bg-emerald-50',
      yellow: 'from-yellow-500 to-yellow-600 text-yellow-600 bg-yellow-50',
      red: 'from-red-500 to-red-600 text-red-600 bg-red-50',
      gray: 'from-gray-500 to-gray-600 text-gray-600 bg-gray-50'
    };
    return colorMap[color] || colorMap.blue;
  }

  $: currentSlideData = slides[currentSlide];
  $: colorClasses = getColorClasses(currentSlideData?.color);
</script>

{#if show}
  <div 
    class="fixed inset-0 bg-gradient-to-br from-indigo-900 via-purple-900 to-pink-900 z-50 flex items-center justify-center overflow-hidden"
    transition:fade={{ duration: 500 }}
  >
    <!-- Animated background elements -->
    <div class="absolute inset-0 overflow-hidden">
      <div class="absolute -top-40 -right-40 w-80 h-80 bg-purple-500 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-pulse"></div>
      <div class="absolute -bottom-40 -left-40 w-80 h-80 bg-pink-500 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-pulse" style="animation-delay: 2s;"></div>
      <div class="absolute top-40 left-40 w-80 h-80 bg-indigo-500 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-pulse" style="animation-delay: 4s;"></div>
    </div>
    
    <div 
      class="relative bg-white/10 backdrop-blur-xl border border-white/20 rounded-3xl shadow-2xl w-full max-w-4xl mx-8 overflow-hidden"
      transition:fly={{ y: 100, duration: 600, delay: 200 }}
    >
      <!-- Header -->
      <div class="flex items-center justify-between p-8 border-b border-white/20">
        <div class="flex items-center space-x-4">
          <div class="w-12 h-12 bg-gradient-to-r {colorClasses.split(' ')[0]} {colorClasses.split(' ')[1]} rounded-2xl flex items-center justify-center shadow-lg">
            <span class="text-white font-bold text-lg">{currentSlide + 1}</span>
          </div>
          <div>
            <h2 class="text-2xl font-bold text-white">Getting Started</h2>
            <p class="text-white/70 text-sm">Welcome to E-Invoice System</p>
          </div>
        </div>
        <button
          on:click={skipIntro}
          class="text-white/60 hover:text-white transition-colors p-2 rounded-full hover:bg-white/10"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>

      <!-- Content -->
      <div class="p-12 min-h-[400px] flex items-center justify-center">
        {#key currentSlide}
          <div 
            class="text-center max-w-2xl"
            in:fly={{ x: 200, duration: 400, delay: 200 }}
            out:fly={{ x: -200, duration: 300 }}
          >
            <!-- Icon -->
            <div class="mb-8 flex justify-center">
              <div class="w-32 h-32 bg-white/20 backdrop-blur-sm rounded-3xl flex items-center justify-center shadow-2xl border border-white/30">
                <svg class="w-16 h-16 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="{currentSlideData.icon}"></path>
                </svg>
              </div>
            </div>

            <!-- Title -->
            <h3 class="text-4xl font-bold text-white mb-6 leading-tight">
              {currentSlideData.title}
            </h3>

            <!-- Description -->
            <p class="text-white/80 text-xl leading-relaxed max-w-xl mx-auto">
              {currentSlideData.description}
            </p>
          </div>
        {/key}
      </div>

      <!-- Progress Indicators -->
      <div class="px-8 pb-6">
        <div class="flex justify-center space-x-3">
          {#each slides as slide, index}
            <button
              on:click={() => goToSlide(index)}
              class="w-4 h-4 rounded-full transition-all duration-500 {index === currentSlide ? 'bg-white scale-125 shadow-lg' : 'bg-white/30 hover:bg-white/50'}"
            ></button>
          {/each}
        </div>
      </div>

      <!-- Navigation -->
      <div class="flex items-center justify-between p-8 bg-white/5 backdrop-blur-sm border-t border-white/20">
        <button
          on:click={prevSlide}
          disabled={currentSlide === 0}
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
          class="flex items-center space-x-3 px-8 py-3 bg-white text-gray-900 rounded-xl hover:bg-white/90 transition-all duration-300 shadow-lg hover:shadow-xl transform hover:scale-105 font-medium"
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