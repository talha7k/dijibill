<script>
  import { createEventDispatcher, onMount } from 'svelte';
  import { fade, fly } from 'svelte/transition';
  import { quintOut } from 'svelte/easing';
  import { tweened } from 'svelte/motion';
  import introSlidesData from '../data/introSlides.json';

  const dispatch = createEventDispatcher();

  export let show = false;

  let currentSlide = 0;
  let slides = introSlidesData.slides;
  let isTransitioning = false;
  let slideContainer;
  let preloadedData = new Map();
  
  // Performance optimizations
  const slideProgress = tweened(0, { duration: 300, easing: quintOut });
  
  // Preload and cache slide data
  onMount(() => {
    // Preload all slide data and color classes
    slides.forEach((slide, index) => {
      preloadedData.set(index, {
        ...slide,
        colorClasses: getColorClasses(slide.color),
        iconSvg: getIconSvg(slide.icon)
      });
    });
    
    // Preload next slide in background
    if (slides.length > 1) {
      preloadSlide(1);
    }
  });

  function preloadSlide(index) {
    if (index >= 0 && index < slides.length && !preloadedData.has(index)) {
      const slide = slides[index];
      preloadedData.set(index, {
        ...slide,
        colorClasses: getColorClasses(slide.color),
        iconSvg: getIconSvg(slide.icon)
      });
    }
  }

  function nextSlide() {
    if (isTransitioning) return;
    if (currentSlide < slides.length - 1) {
      isTransitioning = true;
      const nextIndex = currentSlide + 1;
      
      // Preload next slide in background
      if (nextIndex + 1 < slides.length) {
        preloadSlide(nextIndex + 1);
      }
      
      currentSlide = nextIndex;
      slideProgress.set(nextIndex);
      
      // Use requestAnimationFrame for smoother transitions
      requestAnimationFrame(() => {
        setTimeout(() => isTransitioning = false, 300);
      });
    } else {
      completeIntro();
    }
  }

  function prevSlide() {
    if (isTransitioning) return;
    if (currentSlide > 0) {
      isTransitioning = true;
      const prevIndex = currentSlide - 1;
      
      currentSlide = prevIndex;
      slideProgress.set(prevIndex);
      
      requestAnimationFrame(() => {
        setTimeout(() => isTransitioning = false, 300);
      });
    }
  }

  function goToSlide(index) {
    if (isTransitioning || index === currentSlide) return;
    isTransitioning = true;
    
    // Preload target slide if not already loaded
    preloadSlide(index);
    
    currentSlide = index;
    slideProgress.set(index);
    
    requestAnimationFrame(() => {
      setTimeout(() => isTransitioning = false, 300);
    });
  }

  function skipIntro() {
    completeIntro();
  }

  function completeIntro() {
    dispatch('complete');
  }

  // Optimized color mapping with memoization
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

  function getColorClasses(color) {
    return colorMap[color] || colorMap.blue;
  }

  // Optimized icon mapping with memoization
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

  function getIconSvg(iconClass) {
    return iconMap[iconClass] || iconMap['fas fa-rocket'];
  }

  // Use preloaded data for better performance
  $: currentSlideData = preloadedData.get(currentSlide) || slides[currentSlide];
  $: colorClasses = currentSlideData?.colorClasses || getColorClasses(currentSlideData?.color);
</script>

{#if show}
  <div 
    class="fixed inset-0 bg-gradient-to-br from-indigo-900 via-purple-900 to-pink-900 z-50 flex items-center justify-center overflow-hidden"
    transition:fade={{ duration: 600, easing: quintOut }}
  >
    <!-- Optimized animated background elements with CSS transforms -->
    <div class="absolute inset-0 overflow-hidden">
      <div class="bg-orb bg-orb-1"></div>
      <div class="bg-orb bg-orb-2"></div>
      <div class="bg-orb bg-orb-3"></div>
    </div>
    
    <div 
      class="relative bg-white/10 backdrop-blur-xl border border-white/20 rounded-3xl shadow-2xl w-full max-w-6xl mx-8 overflow-hidden"
      transition:fly={{ y: 50, duration: 800, delay: 300, easing: quintOut }}
      bind:this={slideContainer}
    >
      <!-- Header -->
      <div class="flex items-center justify-between p-8 border-b border-white/20">
        <div class="flex items-center space-x-4">
          <div class="slide-number-badge {colorClasses.split(' ')[0]} {colorClasses.split(' ')[1]}">
            <span class="text-white font-bold text-lg">{currentSlide + 1}</span>
          </div>
          <div>
            <h2 class="text-2xl font-bold text-white">Getting Started</h2>
            <p class="text-white/70 text-sm">Welcome to DijiBill</p>
          </div>
        </div>
        <button
          on:click={skipIntro}
          class="skip-button"
          aria-label="Skip intro"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>

      <!-- Content with optimized transitions -->
      <div class="p-8">
        <div class="slide-content-wrapper">
          {#key currentSlide}
            <div 
              class="slide-content"
              in:fly={{ x: 100, duration: 400, delay: 50, easing: quintOut }}
              out:fly={{ x: -100, duration: 200, easing: quintOut }}
            >
              <!-- Left side - Main content -->
              <div class="slide-main-content">
                <!-- Icon with optimized rendering -->
                <div class="slide-icon-container">
                  <div class="slide-icon-wrapper">
                    <svg class="w-12 h-12 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="{currentSlideData?.iconSvg || getIconSvg(currentSlideData?.icon)}"></path>
                    </svg>
                  </div>
                </div>

                <!-- Title -->
                <h3 class="slide-title">
                  {currentSlideData?.title || ''}
                </h3>

                <!-- Description -->
                <p class="slide-description">
                  {currentSlideData?.description || ''}
                </p>
              </div>

              <!-- Right side - Features with staggered animations -->
              <div class="slide-features">
                <h4 class="slide-features-title">Key Features</h4>
                <div class="features-grid">
                  {#each (currentSlideData?.features || []) as feature, index}
                    <div 
                      class="feature-card"
                      style="--delay: {index * 50}ms"
                      in:fly={{ x: 30, duration: 300, delay: 100 + (index * 50), easing: quintOut }}
                    >
                      <div class="feature-content">
                        <div class="feature-icon">
                          <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                          </svg>
                        </div>
                        <div class="feature-text">
                          <h5 class="feature-title">{feature.title}</h5>
                          <p class="feature-description">{feature.description}</p>
                        </div>
                      </div>
                    </div>
                  {/each}
                </div>
              </div>
            </div>
          {/key}
        </div>
      </div>

      <!-- Optimized Progress Indicators -->
      <div class="progress-container">
        <div class="progress-indicators">
          {#each slides as slide, index}
            <button
              on:click={() => goToSlide(index)}
              class="progress-dot {index === currentSlide ? 'active' : ''}"
              disabled={isTransitioning}
              aria-label="Go to slide {index + 1}"
            ></button>
          {/each}
        </div>
      </div>

      <!-- Navigation with improved accessibility -->
      <div class="navigation-container">
        <button
          on:click={prevSlide}
          disabled={currentSlide === 0 || isTransitioning}
          class="nav-button nav-button-prev"
          aria-label="Previous slide"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
          </svg>
          <span class="font-medium">Previous</span>
        </button>

        <div class="slide-counter">
          {currentSlide + 1} of {slides.length}
        </div>

        <button
          on:click={nextSlide}
          disabled={isTransitioning}
          class="nav-button nav-button-next"
          aria-label={currentSlide === slides.length - 1 ? 'Get started' : 'Next slide'}
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
  /* Hardware-accelerated animations for better performance */
  .bg-orb {
    position: absolute;
    width: 20rem;
    height: 20rem;
    border-radius: 50%;
    mix-blend-mode: multiply;
    filter: blur(3rem);
    opacity: 0.7;
    will-change: transform;
    animation: pulse-optimized 6s ease-in-out infinite;
  }

  .bg-orb-1 {
    top: -10rem;
    right: -10rem;
    background: #a855f7;
    animation-delay: 0s;
  }

  .bg-orb-2 {
    bottom: -10rem;
    left: -10rem;
    background: #ec4899;
    animation-delay: 2s;
  }

  .bg-orb-3 {
    top: 10rem;
    left: 10rem;
    background: #6366f1;
    animation-delay: 4s;
  }

  @keyframes pulse-optimized {
    0%, 100% {
      transform: scale(1) translateZ(0);
      opacity: 0.7;
    }
    50% {
      transform: scale(1.1) translateZ(0);
      opacity: 0.5;
    }
  }

  /* Optimized slide components */
  .slide-number-badge {
    width: 3rem;
    height: 3rem;
    background: linear-gradient(to right, var(--tw-gradient-stops));
    border-radius: 1rem;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
    will-change: transform;
  }

  .skip-button {
    color: rgba(255, 255, 255, 0.6);
    transition: color 0.3s ease;
    padding: 0.5rem;
    border-radius: 9999px;
    will-change: transform;
  }

  .skip-button:hover {
    color: white;
    background-color: rgba(255, 255, 255, 0.1);
  }

  .slide-content-wrapper {
    position: relative;
    min-height: 400px;
  }

  .slide-content {
    display: grid;
    grid-template-columns: 1fr;
    gap: 2rem;
    align-items: center;
    will-change: transform;
  }

  @media (min-width: 1024px) {
    .slide-content {
      grid-template-columns: 1fr 1fr;
    }
  }

  .slide-main-content {
    text-align: center;
  }

  @media (min-width: 1024px) {
    .slide-main-content {
      text-align: left;
    }
  }

  .slide-icon-container {
    margin-bottom: 1.5rem;
    display: flex;
    justify-content: center;
  }

  @media (min-width: 1024px) {
    .slide-icon-container {
      justify-content: flex-start;
    }
  }

  .slide-icon-wrapper {
    width: 6rem;
    height: 6rem;
    background-color: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(4px);
    border-radius: 1rem;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
    border: 1px solid rgba(255, 255, 255, 0.3);
    will-change: transform;
  }

  .slide-title {
    font-size: 1.875rem;
    font-weight: 700;
    color: white;
    margin-bottom: 1rem;
    line-height: 1.25;
  }

  @media (min-width: 1024px) {
    .slide-title {
      font-size: 2.25rem;
    }
  }

  .slide-description {
    color: rgba(255, 255, 255, 0.8);
    font-size: 1.125rem;
    line-height: 1.75;
  }

  .slide-features {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .slide-features-title {
    font-size: 1.25rem;
    font-weight: 600;
    color: white;
    margin-bottom: 1.5rem;
    text-align: center;
  }

  @media (min-width: 1024px) {
    .slide-features-title {
      text-align: left;
    }
  }

  .features-grid {
    display: grid;
    gap: 1rem;
  }

  .feature-card {
    background-color: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(4px);
    border-radius: 0.75rem;
    padding: 1rem;
    border: 1px solid rgba(255, 255, 255, 0.2);
    transition: background-color 0.3s ease;
    will-change: transform;
    animation-delay: var(--delay);
  }

  .feature-card:hover {
    background-color: rgba(255, 255, 255, 0.15);
  }

  .feature-content {
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
  }

  .feature-icon {
    width: 2rem;
    height: 2rem;
    background-color: rgba(255, 255, 255, 0.2);
    border-radius: 0.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    margin-top: 0.25rem;
  }

  .feature-text {
    flex: 1;
  }

  .feature-title {
    font-weight: 600;
    color: white;
    font-size: 0.875rem;
    margin-bottom: 0.25rem;
  }

  .feature-description {
    color: rgba(255, 255, 255, 0.7);
    font-size: 0.75rem;
    line-height: 1.75;
  }

  .progress-container {
    padding: 0 2rem 1.5rem;
  }

  .progress-indicators {
    display: flex;
    justify-content: center;
    gap: 0.75rem;
  }

  .progress-dot {
    width: 0.75rem;
    height: 0.75rem;
    border-radius: 50%;
    background-color: rgba(255, 255, 255, 0.3);
    transition: all 0.3s ease;
    border: none;
    cursor: pointer;
    will-change: transform;
  }

  .progress-dot:hover {
    background-color: rgba(255, 255, 255, 0.5);
  }

  .progress-dot.active {
    background-color: white;
    transform: scale(1.25);
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
  }

  .progress-dot:disabled {
    cursor: not-allowed;
  }

  .navigation-container {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 2rem;
    background-color: rgba(255, 255, 255, 0.05);
    backdrop-filter: blur(4px);
    border-top: 1px solid rgba(255, 255, 255, 0.2);
  }

  .nav-button {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1.5rem;
    border-radius: 0.75rem;
    transition: all 0.3s ease;
    border: none;
    cursor: pointer;
    font-weight: 500;
    will-change: transform;
  }

  .nav-button-prev {
    color: rgba(255, 255, 255, 0.7);
  }

  .nav-button-prev:hover:not(:disabled) {
    color: white;
    background-color: rgba(255, 255, 255, 0.1);
  }

  .nav-button-prev:disabled {
    opacity: 0.3;
    cursor: not-allowed;
  }

  .nav-button-next {
    background-color: white;
    color: #374151;
    padding: 0.75rem 2rem;
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
  }

  .nav-button-next:hover:not(:disabled) {
    background-color: rgba(255, 255, 255, 0.9);
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
    transform: scale(1.05);
  }

  .nav-button-next:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    transform: none;
  }

  .slide-counter {
    color: rgba(255, 255, 255, 0.6);
    font-weight: 500;
  }

  /* Performance optimizations */
  * {
    backface-visibility: hidden;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }

  /* Reduce layout thrashing */
  .slide-content,
  .feature-card,
  .progress-dot,
  .nav-button {
    contain: layout style paint;
  }

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