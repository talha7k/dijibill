import { writable } from 'svelte/store';

// Store to track Wails runtime initialization status
export const wailsReady = writable(false);
export const wailsError = writable(null);

// Function to initialize Wails runtime
export async function initializeWailsRuntime() {
  console.log('🚀 Initializing Wails runtime...');
  
  // Reset stores
  wailsReady.set(false);
  wailsError.set(null);
  
  let attempts = 0;
  const maxAttempts = 100; // 10 seconds max wait (100ms intervals)
  
  while (attempts < maxAttempts) {
    try {
      // Check if window object exists (for SSR compatibility)
      if (typeof window === 'undefined') {
        console.log('⏳ Window object not available yet...');
        await new Promise(resolve => setTimeout(resolve, 100));
        attempts++;
        continue;
      }

      // Check for Wails runtime availability
      if (window['go'] && window['go']['main'] && window['go']['main']['App']) {
        console.log('✅ Wails runtime initialized successfully');
        wailsReady.set(true);
        wailsError.set(null);
        return true;
      }
      
      // Check if we're in development mode and should wait longer
      const isDev = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';
      if (isDev && attempts < 50) {
        console.log(`⏳ Attempt ${attempts + 1}/${maxAttempts} - Waiting for Wails runtime (dev mode)...`);
      } else if (attempts % 10 === 0) {
        console.log(`⏳ Attempt ${attempts + 1}/${maxAttempts} - Waiting for Wails runtime...`);
      }
      
      await new Promise(resolve => setTimeout(resolve, 100));
      attempts++;
    } catch (error) {
      console.error('❌ Error during Wails runtime check:', error);
      await new Promise(resolve => setTimeout(resolve, 100));
      attempts++;
    }
  }
  
  const errorMessage = `Wails runtime failed to initialize after ${maxAttempts * 100}ms. Please ensure the Wails development server is running.`;
  console.error('❌', errorMessage);
  wailsError.set(errorMessage);
  wailsReady.set(false);
  return false;
}

// Helper function to check if Wails is ready
export function isWailsReady() {
  if (typeof window === 'undefined') return false;
  return window['go'] && window['go']['main'] && window['go']['main']['App'];
}

// Helper function to wait for Wails with a promise
export function waitForWails(timeoutMs = 10000) {
  return new Promise((resolve, reject) => {
    if (isWailsReady()) {
      resolve(true);
      return;
    }

    const startTime = Date.now();
    const checkInterval = setInterval(() => {
      if (isWailsReady()) {
        clearInterval(checkInterval);
        resolve(true);
      } else if (Date.now() - startTime > timeoutMs) {
        clearInterval(checkInterval);
        reject(new Error(`Wails runtime not available after ${timeoutMs}ms`));
      }
    }, 100);
  });
}