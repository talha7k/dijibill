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
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
    transition:fade={{ duration: 300 }}
  >
    <div 
      class="bg-white rounded-2xl shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-hidden"
      transition:fly={{ y: 50, duration: 300 }}
    >
      <!-- Header -->
      <div class="flex items-center justify-between p-6 border-b border-gray-200">
        <div class="flex items-center space-x-3">
          <div class="w-8 h-8 bg-gradient-to-r {colorClasses.split(' ')[0]} {colorClasses.split(' ')[1]} rounded-lg flex items-center justify-center">
            <span class="text-white font-bold text-sm">{currentSlide + 1}</span>
          </div>
          <h2 class="text-xl font-bold text-gray-900">Getting Started</h2>
        </div>
        <button
          on:click={skipIntro}
          class="text-gray-400 hover:text-gray-600 transition-colors"
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
            class="text-center"
            in:fly={{ x: 100, duration: 300, delay: 150 }}
            out:fly={{ x: -100, duration: 300 }}
          >
            <!-- Icon -->
            <div class="mb-6 flex justify-center">
              <div class="w-20 h-20 {colorClasses.split(' ')[3]} rounded-full flex items-center justify-center">
                <svg class="w-10 h-10 {colorClasses.split(' ')[2]}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="{currentSlideData.icon}"></path>
                </svg>
              </div>
            </div>

            <!-- Title -->
            <h3 class="text-2xl font-bold text-gray-900 mb-4">
              {currentSlideData.title}
            </h3>

            <!-- Description -->
            <p class="text-gray-600 text-lg leading-relaxed max-w-lg mx-auto">
              {currentSlideData.description}
            </p>
          </div>
        {/key}
      </div>

      <!-- Progress Indicators -->
      <div class="px-8 pb-4">
        <div class="flex justify-center space-x-2">
          {#each slides as slide, index}
            <button
              on:click={() => goToSlide(index)}
              class="w-3 h-3 rounded-full transition-all duration-300 {index === currentSlide ? 'bg-gradient-to-r ' + colorClasses.split(' ')[0] + ' ' + colorClasses.split(' ')[1] : 'bg-gray-300 hover:bg-gray-400'}"
            ></button>
          {/each}
        </div>
      </div>

      <!-- Navigation -->
      <div class="flex items-center justify-between p-6 bg-gray-50 border-t border-gray-200">
        <button
          on:click={prevSlide}
          disabled={currentSlide === 0}
          class="flex items-center space-x-2 px-4 py-2 text-gray-600 hover:text-gray-800 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
          </svg>
          <span>Previous</span>
        </button>

        <div class="text-sm text-gray-500">
          {currentSlide + 1} of {slides.length}
        </div>

        <button
          on:click={nextSlide}
          class="flex items-center space-x-2 px-6 py-2 bg-gradient-to-r {colorClasses.split(' ')[0]} {colorClasses.split(' ')[1]} text-white rounded-lg hover:shadow-lg transition-all duration-300"
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