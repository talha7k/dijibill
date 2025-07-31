<script>
  import { createEventDispatcher, onMount, onDestroy } from 'svelte';
  import { gsap } from 'gsap';
  import introSlidesData from '../data/introSlides.json';

  const dispatch = createEventDispatcher();

  export let show = false;

  let currentSlide = 0;
  let slides = introSlidesData.slides;
  let isTransitioning = false;
  let slideContainer;
  let mainContainer;
  let preloadedData = new Map();
  let tl; // GSAP timeline
  
  // Preload and cache slide data
  onMount(() => {
    // Preload all slide data and color classes
    slides.forEach((slide, index) => {
      preloadedData.set(index, {
        ...slide,
        colorClasses: getColorClasses(slide.color)
      });
    });
    
    // Initialize GSAP animations
    if (show) {
      initializeAnimations();
    }
  });

  onDestroy(() => {
    if (tl) {
      tl.kill();
    }
  });

  function initializeAnimations() {
    if (!mainContainer) return;
    
    // Create master timeline
    tl = gsap.timeline();
    
    // Initial entrance animation
    tl.fromTo(mainContainer, 
      { opacity: 0, scale: 0.9 },
      { opacity: 1, scale: 1, duration: 0.8, ease: "power3.out" }
    );
    
    // Animate background orbs
    gsap.to(".bg-orb-1", {
      rotation: 360,
      duration: 20,
      repeat: -1,
      ease: "none"
    });
    
    gsap.to(".bg-orb-2", {
      rotation: -360,
      duration: 25,
      repeat: -1,
      ease: "none"
    });
    
    gsap.to(".bg-orb-3", {
      rotation: 360,
      duration: 30,
      repeat: -1,
      ease: "none"
    });
    
    // Animate slide content
    animateSlideContent();
  }

  function animateSlideContent() {
    const slideContent = document.querySelector('.slide-content');
    const features = document.querySelectorAll('.feature-card');
    const icon = document.querySelector('.slide-icon-wrapper');
    const title = document.querySelector('.slide-title');
    const description = document.querySelector('.slide-description');
    
    if (!slideContent) return;
    
    // Create slide animation timeline
    const slideTl = gsap.timeline();
    
    // Animate icon
    if (icon) {
      slideTl.fromTo(icon,
        { scale: 0, rotation: -180 },
        { scale: 1, rotation: 0, duration: 0.6, ease: "back.out(1.7)" }
      );
    }
    
    // Animate title
    if (title) {
      slideTl.fromTo(title,
        { y: 30, opacity: 0 },
        { y: 0, opacity: 1, duration: 0.5, ease: "power2.out" },
        "-=0.3"
      );
    }
    
    // Animate description
    if (description) {
      slideTl.fromTo(description,
        { y: 20, opacity: 0 },
        { y: 0, opacity: 1, duration: 0.5, ease: "power2.out" },
        "-=0.3"
      );
    }
    
    // Animate features with stagger
    if (features.length > 0) {
      slideTl.fromTo(features,
        { x: 50, opacity: 0 },
        { 
          x: 0, 
          opacity: 1, 
          duration: 0.4, 
          ease: "power2.out",
          stagger: 0.1
        },
        "-=0.2"
      );
    }
  }

  function preloadSlide(index) {
    if (index >= 0 && index < slides.length && !preloadedData.has(index)) {
      const slide = slides[index];
      preloadedData.set(index, {
        ...slide,
        colorClasses: getColorClasses(slide.color)
      });
    }
  }

  function nextSlide() {
    if (isTransitioning) return;
    if (currentSlide < slides.length - 1) {
      isTransitioning = true;
      const nextIndex = currentSlide + 1;
      
      // Animate slide transition with GSAP
      const currentContent = document.querySelector('.slide-content');
      if (currentContent) {
        gsap.to(currentContent, {
          x: -100,
          opacity: 0,
          duration: 0.3,
          ease: "power2.in",
          onComplete: () => {
            currentSlide = nextIndex;
            // Trigger reactive update and then animate in
            setTimeout(() => {
              animateSlideContent();
              isTransitioning = false;
            }, 50);
          }
        });
      } else {
        currentSlide = nextIndex;
        setTimeout(() => {
          animateSlideContent();
          isTransitioning = false;
        }, 50);
      }
    } else {
      completeIntro();
    }
  }

  function prevSlide() {
    if (isTransitioning) return;
    if (currentSlide > 0) {
      isTransitioning = true;
      const prevIndex = currentSlide - 1;
      
      // Animate slide transition with GSAP
      const currentContent = document.querySelector('.slide-content');
      if (currentContent) {
        gsap.to(currentContent, {
          x: 100,
          opacity: 0,
          duration: 0.3,
          ease: "power2.in",
          onComplete: () => {
            currentSlide = prevIndex;
            // Trigger reactive update and then animate in
            setTimeout(() => {
              animateSlideContent();
              isTransitioning = false;
            }, 50);
          }
        });
      } else {
        currentSlide = prevIndex;
        setTimeout(() => {
          animateSlideContent();
          isTransitioning = false;
        }, 50);
      }
    }
  }

  function goToSlide(index) {
    if (isTransitioning || index === currentSlide) return;
    isTransitioning = true;
    
    // Preload target slide if not already loaded
    preloadSlide(index);
    
    // Animate slide transition with GSAP
    const currentContent = document.querySelector('.slide-content');
    if (currentContent) {
      gsap.to(currentContent, {
        scale: 0.8,
        opacity: 0,
        duration: 0.3,
        ease: "power2.in",
        onComplete: () => {
          currentSlide = index;
          // Trigger reactive update and then animate in
          setTimeout(() => {
            animateSlideContent();
            isTransitioning = false;
          }, 50);
        }
      });
    } else {
      currentSlide = index;
      setTimeout(() => {
        animateSlideContent();
        isTransitioning = false;
      }, 50);
    }
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

  // Use preloaded data for better performance
  $: currentSlideData = preloadedData.get(currentSlide) || slides[currentSlide];
  $: colorClasses = currentSlideData?.colorClasses || getColorClasses(currentSlideData?.color);
  
  // Reactive statement to trigger animations when show changes
  $: if (show && mainContainer) {
    initializeAnimations();
  }
  
  // Reactive statement to trigger slide content animations when currentSlide changes
  $: if (currentSlide !== undefined && !isTransitioning) {
    setTimeout(() => animateSlideContent(), 100);
  }
</script>

{#if show}
  <div 
    class="intro-slider-overlay"
    bind:this={mainContainer}
  >
    <!-- Animated background elements with GSAP -->
    <div class="intro-background">
      <div class="bg-orb bg-orb-1"></div>
      <div class="bg-orb bg-orb-2"></div>
      <div class="bg-orb bg-orb-3"></div>
    </div>
    
    <!-- Use PageLayout for proper structure -->
    <div class="intro-page-layout">
      <!-- Header -->
      <div class="intro-header">
        <div class="intro-header-content">
          <div class="slide-number-badge {colorClasses.split(' ')[0]} {colorClasses.split(' ')[1]}">
            <span class="text-white font-bold text-lg">{currentSlide + 1}</span>
          </div>
          <div class="intro-title-section">
            <h2 class="intro-main-title">Getting Started</h2>
            <p class="intro-subtitle">Welcome to DijiBill</p>
          </div>
        </div>
        <button
          on:click={skipIntro}
          class="skip-button"
          aria-label="Skip intro"
        >
          <i class="fas fa-times text-xl"></i>
        </button>
      </div>

      <!-- Main content area with proper constraints -->
      <div class="intro-content-container">
        <div class="intro-content-wrapper">
          {#key currentSlide}
            <div class="slide-content">
              <!-- Left side - Main content -->
              <div class="slide-main-content">
                <!-- Icon -->
                <div class="slide-icon-container">
                  <div class="slide-icon-wrapper">
                    <i class="{currentSlideData?.icon || 'fas fa-rocket'} text-white text-5xl"></i>
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

              <!-- Right side - Features -->
              <div class="slide-features">
                <h4 class="slide-features-title">Key Features</h4>
                <div class="features-grid">
                  {#each (currentSlideData?.features || []) as feature, index}
                    <div class="feature-card" style="--delay: {index * 50}ms">
                      <div class="feature-content">
                        <div class="feature-icon">
                          <i class="fas fa-check text-white text-sm"></i>
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

      <!-- Bottom section with progress and navigation -->
      <div class="intro-footer">
        <!-- Progress Indicators -->
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

        <!-- Navigation -->
        <div class="navigation-container">
          <button
            on:click={prevSlide}
            disabled={currentSlide === 0 || isTransitioning}
            class="nav-button nav-button-prev"
            aria-label="Previous slide"
          >
            <i class="fas fa-chevron-left"></i>
            <span class="font-medium hidden sm:inline">Previous</span>
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
            <span class="hidden sm:inline">{currentSlide === slides.length - 1 ? 'Get Started' : 'Next'}</span>
            <span class="sm:hidden">{currentSlide === slides.length - 1 ? 'Start' : 'Next'}</span>
            {#if currentSlide === slides.length - 1}
              <i class="fas fa-check-circle"></i>
            {:else}
              <i class="fas fa-chevron-right"></i>
            {/if}
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}

<style>
  /* Main overlay container */
  .intro-slider-overlay {
    position: fixed;
    inset: 0;
    background: linear-gradient(to bottom right, #312e81, #7c3aed, #ec4899);
    z-index: 50;
    overflow: hidden;
  }

  /* Background elements */
  .intro-background {
    position: absolute;
    inset: 0;
    overflow: hidden;
  }

  /* GSAP-powered background orbs */
  .bg-orb {
    position: absolute;
    width: 20rem;
    height: 20rem;
    border-radius: 50%;
    mix-blend-mode: multiply;
    filter: blur(3rem);
    opacity: 0.7;
    will-change: transform;
  }

  .bg-orb-1 {
    top: -10rem;
    right: -10rem;
    background: #a855f7;
  }

  .bg-orb-2 {
    bottom: -10rem;
    left: -10rem;
    background: #ec4899;
  }

  .bg-orb-3 {
    top: 10rem;
    left: 10rem;
    background: #6366f1;
  }

  /* Page layout structure */
  .intro-page-layout {
    position: relative;
    height: 100vh;
    display: flex;
    flex-direction: column;
    max-width: 100vw;
    overflow: hidden;
  }

  /* Header section */
  .intro-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1.5rem 2rem;
    border-bottom: 1px solid rgba(255, 255, 255, 0.2);
    background: rgba(255, 255, 255, 0.05);
    backdrop-filter: blur(12px);
    flex-shrink: 0;
  }

  .intro-header-content {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .intro-title-section {
    display: flex;
    flex-direction: column;
  }

  .intro-main-title {
    font-size: 1.5rem;
    font-weight: 700;
    color: white;
    margin: 0;
    line-height: 1.2;
  }

  .intro-subtitle {
    color: rgba(255, 255, 255, 0.7);
    font-size: 0.875rem;
    margin: 0;
  }

  /* Content container with proper constraints */
  .intro-content-container {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem;
    overflow-y: auto;
    min-height: 0; /* Important for flex child */
  }

  .intro-content-wrapper {
    width: 100%;
    max-width: 1200px;
    margin: 0 auto;
  }

  /* Footer section */
  .intro-footer {
    background: rgba(255, 255, 255, 0.05);
    backdrop-filter: blur(12px);
    border-top: 1px solid rgba(255, 255, 255, 0.2);
    flex-shrink: 0;
  }

  /* Slide components */
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
    flex-shrink: 0;
  }

  .skip-button {
    color: rgba(255, 255, 255, 0.6);
    transition: color 0.3s ease;
    padding: 0.75rem;
    border-radius: 50%;
    will-change: transform;
    border: none;
    background: none;
    cursor: pointer;
  }

  .skip-button:hover {
    color: white;
    background-color: rgba(255, 255, 255, 0.1);
  }

  .slide-content {
    display: grid;
    grid-template-columns: 1fr;
    gap: 2rem;
    align-items: start;
    will-change: transform;
    min-height: 400px;
  }

  @media (min-width: 1024px) {
    .slide-content {
      grid-template-columns: 1fr 1fr;
      gap: 4rem;
      align-items: center;
    }
  }

  .slide-main-content {
    text-align: center;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  @media (min-width: 1024px) {
    .slide-main-content {
      text-align: left;
    }
  }

  .slide-icon-container {
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
    border-radius: 1.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
    border: 1px solid rgba(255, 255, 255, 0.3);
    will-change: transform;
  }

  @media (min-width: 1024px) {
    .slide-icon-wrapper {
      width: 8rem;
      height: 8rem;
    }
  }

  .slide-title {
    font-size: 2rem;
    font-weight: 700;
    color: white;
    margin: 0;
    line-height: 1.2;
  }

  @media (min-width: 1024px) {
    .slide-title {
      font-size: 2.5rem;
    }
  }

  .slide-description {
    color: rgba(255, 255, 255, 0.8);
    font-size: 1.125rem;
    line-height: 1.75;
    max-width: 500px;
    margin: 0 auto;
  }

  @media (min-width: 1024px) {
    .slide-description {
      margin: 0;
      font-size: 1.25rem;
    }
  }

  .slide-features {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .slide-features-title {
    font-size: 1.25rem;
    font-weight: 600;
    color: white;
    margin: 0;
    text-align: center;
  }

  @media (min-width: 1024px) {
    .slide-features-title {
      text-align: left;
      font-size: 1.5rem;
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
    transition: background-color 0.3s ease, transform 0.3s ease;
    will-change: transform;
  }

  .feature-card:hover {
    background-color: rgba(255, 255, 255, 0.15);
    transform: translateY(-2px);
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
    margin-top: 0.125rem;
  }

  .feature-text {
    flex: 1;
  }

  .feature-title {
    font-weight: 600;
    color: white;
    font-size: 0.875rem;
    margin: 0 0 0.25rem 0;
  }

  .feature-description {
    color: rgba(255, 255, 255, 0.7);
    font-size: 0.8125rem;
    line-height: 1.5;
    margin: 0;
  }

  .progress-container {
    padding: 1rem 2rem;
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
    padding: 1.5rem 2rem;
  }

  .nav-button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem 1.5rem;
    border-radius: 0.75rem;
    transition: all 0.3s ease;
    border: none;
    cursor: pointer;
    font-weight: 500;
    will-change: transform;
    font-size: 0.875rem;
  }

  .nav-button-prev {
    color: rgba(255, 255, 255, 0.7);
    background: none;
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
    font-size: 0.875rem;
  }

  /* Responsive adjustments */
  @media (max-width: 768px) {
    .intro-header {
      padding: 1rem;
    }

    .intro-content-container {
      padding: 1rem;
    }

    .navigation-container {
      padding: 1rem;
    }

    .slide-content {
      gap: 1.5rem;
      min-height: 300px;
    }

    .slide-title {
      font-size: 1.75rem;
    }

    .slide-description {
      font-size: 1rem;
    }

    .nav-button {
      padding: 0.5rem 1rem;
      font-size: 0.8125rem;
    }

    .nav-button-next {
      padding: 0.5rem 1.5rem;
    }
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
</style>