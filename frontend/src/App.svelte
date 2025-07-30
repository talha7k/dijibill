<script>
  import Dashboard from './Dashboard.svelte'
  import QRValidation from './QRValidation.svelte'
  import {Greet} from '../wailsjs/go/main/App.js'

  let currentView = 'dashboard' // Default to dashboard
  let resultText = "Welcome to DijInvoice - ZATCA Compliant Invoicing"
  let name = ''

  function greet() {
    Greet(name).then(result => resultText = result)
  }

  function switchView(view) {
    currentView = view
  }
</script>

<div class="min-h-screen bg-gradient-to-br from-primary to-secondary" data-theme="dijinvoice">
  <!-- Header -->
  <div class="navbar bg-base-100/10 backdrop-blur-lg border-b border-white/20">
    <div class="navbar-start">
      <div class="dropdown">
        <div role="button" class="btn btn-ghost lg:hidden">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h8m-8 6h16"></path>
          </svg>
        </div>
        <ul class="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 rounded-box w-52">
          <li><button on:click={() => switchView('dashboard')} class="btn btn-ghost">Dashboard</button></li>
          <li><button on:click={() => switchView('qr-validation')} class="btn btn-ghost">QR Validation</button></li>
          <li><button on:click={() => switchView('demo')} class="btn btn-ghost">Demo</button></li>
        </ul>
      </div>
      <div class="flex items-center gap-3">
        <div>
          <h1 class="text-2xl font-bold text-white">DijInvoice</h1>
          <p class="text-white/80 text-sm">ZATCA Compliant E-Invoicing</p>
        </div>
      </div>
    </div>
    
    <div class="navbar-center hidden lg:flex">
      <ul class="menu menu-horizontal px-1 gap-2">
        <li>
          <button 
            class="btn {currentView === 'dashboard' ? 'btn-primary' : 'btn-ghost text-white'}"
            on:click={() => switchView('dashboard')}
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"></path>
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5a2 2 0 012-2h4a2 2 0 012 2v6H8V5z"></path>
            </svg>
            Dashboard
          </button>
        </li>
        <li>
          <button 
            class="btn {currentView === 'qr-validation' ? 'btn-primary' : 'btn-ghost text-white'}"
            on:click={() => switchView('qr-validation')}
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z"></path>
            </svg>
            QR Validation
          </button>
        </li>
        <li>
          <button 
            class="btn {currentView === 'demo' ? 'btn-primary' : 'btn-ghost text-white'}"
            on:click={() => switchView('demo')}
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
            </svg>
            Demo
          </button>
        </li>
      </ul>
    </div>
    
    <div class="navbar-end">
      <div class="dropdown dropdown-end">
        <div role="button" class="btn btn-ghost btn-circle text-white">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"></path>
          </svg>
        </div>
        <ul class="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 rounded-box w-52">
          <li><button class="btn btn-ghost">Settings</button></li>
          <li><button class="btn btn-ghost">About</button></li>
        </ul>
      </div>
    </div>
  </div>

  <!-- Main Content -->
  <main class="{currentView === 'dashboard' ? '' : 'container mx-auto px-4 py-8'}">
    {#if currentView === 'dashboard'}
      <Dashboard />
    {:else if currentView === 'qr-validation'}
      <div class="container mx-auto px-4 py-8">
        <QRValidation />
      </div>
    {:else if currentView === 'demo'}
      <div class="container mx-auto px-4 py-8">
        <div class="hero min-h-[60vh]">
          <div class="hero-content text-center">
            <div class="max-w-md">
              <div class="card bg-base-100 shadow-2xl">
                <div class="card-body">
                  <h2 class="card-title justify-center text-2xl mb-4">Welcome Demo</h2>
                  
                  <div class="alert alert-info mb-4">
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    </svg>
                    <span>{resultText}</span>
                  </div>

                  <div class="form-control w-full">
                    <label class="label" for="name-input">
                      <span class="label-text">Enter your name</span>
                    </label>
                    <div class="join w-full">
                      <input 
                        id="name-input"
                        type="text" 
                        placeholder="Your name here..." 
                        class="input input-bordered join-item flex-1" 
                        bind:value={name}
                      />
                      <button class="btn btn-primary join-item" on:click={greet}>
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.746 0 3.332.477 4.5 1.253v13C19.832 18.477 18.246 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"></path>
                        </svg>
                        Greet
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    {/if}
  </main>

  <!-- Footer -->
  <footer class="footer footer-center p-4 bg-base-100/10 backdrop-blur-lg border-t border-white/20 text-white">
    <aside>
      <p>Copyright © 2024 - DijInvoice. ZATCA Compliant E-Invoicing Solution</p>
    </aside>
  </footer>
</div>
