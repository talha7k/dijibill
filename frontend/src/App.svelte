<script>
  import logo from './assets/images/logo-universal.png'
  import QRValidation from './QRValidation.svelte'
  import {Greet} from '../wailsjs/go/main/App.js'

  let currentView = 'qr-validation' // Default to QR validation
  let resultText = "Welcome to DijInvoice - ZATCA Compliant Invoicing"
  let name

  function greet() {
    Greet(name).then(result => resultText = result)
  }

  function switchView(view) {
    currentView = view
  }
</script>

<main>
  <header>
    <img alt="DijInvoice logo" id="logo" src="{logo}">
    <h1>DijInvoice</h1>
    <p class="subtitle">ZATCA Compliant E-Invoicing Solution</p>
  </header>

  <nav class="navigation">
    <button 
      class="nav-btn {currentView === 'qr-validation' ? 'active' : ''}"
      on:click={() => switchView('qr-validation')}
    >
      QR Validation
    </button>
    <button 
      class="nav-btn {currentView === 'demo' ? 'active' : ''}"
      on:click={() => switchView('demo')}
    >
      Demo
    </button>
  </nav>

  <div class="content">
    {#if currentView === 'qr-validation'}
      <QRValidation />
    {:else if currentView === 'demo'}
      <div class="demo-section">
        <div class="result" id="result">{resultText}</div>
        <div class="input-box" id="input">
          <input autocomplete="off" bind:value={name} class="input" id="name" type="text" placeholder="Enter your name"/>
          <button class="btn" on:click={greet}>Greet</button>
        </div>
      </div>
    {/if}
  </div>
</main>

<style>
  main {
    min-height: 100vh;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  }

  header {
    text-align: center;
    padding: 2rem 0;
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  }

  #logo {
    display: block;
    width: 80px;
    height: 80px;
    margin: 0 auto 1rem;
    border-radius: 50%;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  }

  h1 {
    color: white;
    margin: 0;
    font-size: 2.5rem;
    font-weight: 700;
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  }

  .subtitle {
    color: rgba(255, 255, 255, 0.9);
    margin: 0.5rem 0 0;
    font-size: 1.1rem;
    font-weight: 300;
  }

  .navigation {
    display: flex;
    justify-content: center;
    gap: 1rem;
    padding: 1rem;
    background: rgba(255, 255, 255, 0.05);
    backdrop-filter: blur(10px);
  }

  .nav-btn {
    padding: 0.75rem 1.5rem;
    border: 2px solid rgba(255, 255, 255, 0.3);
    background: rgba(255, 255, 255, 0.1);
    color: white;
    border-radius: 25px;
    cursor: pointer;
    transition: all 0.3s ease;
    font-weight: 600;
    font-size: 1rem;
  }

  .nav-btn:hover {
    background: rgba(255, 255, 255, 0.2);
    border-color: rgba(255, 255, 255, 0.5);
    transform: translateY(-2px);
  }

  .nav-btn.active {
    background: rgba(255, 255, 255, 0.9);
    color: #667eea;
    border-color: white;
  }

  .content {
    padding: 2rem;
    min-height: calc(100vh - 200px);
  }

  .demo-section {
    max-width: 500px;
    margin: 0 auto;
    background: rgba(255, 255, 255, 0.95);
    padding: 2rem;
    border-radius: 15px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
    backdrop-filter: blur(10px);
  }

  .result {
    height: auto;
    line-height: 1.5;
    margin: 0 0 1.5rem 0;
    padding: 1rem;
    background: #f8f9fa;
    border-radius: 8px;
    color: #495057;
    text-align: center;
    font-size: 1.1rem;
  }

  .input-box {
    display: flex;
    gap: 1rem;
    align-items: center;
  }

  .input-box .btn {
    padding: 0.75rem 1.5rem;
    border-radius: 8px;
    border: none;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    cursor: pointer;
    font-weight: 600;
    transition: all 0.3s ease;
    white-space: nowrap;
  }

  .input-box .btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
  }

  .input-box .input {
    flex: 1;
    border: 2px solid #dee2e6;
    border-radius: 8px;
    outline: none;
    padding: 0.75rem 1rem;
    background-color: white;
    font-size: 1rem;
    transition: border-color 0.3s ease;
  }

  .input-box .input:hover {
    border-color: #adb5bd;
  }

  .input-box .input:focus {
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }

  @media (max-width: 768px) {
    .content {
      padding: 1rem;
    }

    .demo-section {
      padding: 1.5rem;
    }

    .input-box {
      flex-direction: column;
      align-items: stretch;
    }

    .nav-btn {
      padding: 0.5rem 1rem;
      font-size: 0.9rem;
    }

    h1 {
      font-size: 2rem;
    }
  }
</style>
