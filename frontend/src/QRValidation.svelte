<script>
  import { ValidateZATCAQRCode, GetQRCodeInfo, GenerateSampleQRCode } from '../wailsjs/go/main/App.js'

  let qrCodeInput = ''
  let validationResult = null
  let qrInfo = null
  let isLoading = false
  let error = null

  async function validateQRCode() {
    if (!qrCodeInput.trim()) {
      error = 'Please enter a QR code'
      return
    }

    isLoading = true
    error = null
    validationResult = null
    qrInfo = null

    try {
      // Validate the QR code
      const isValid = await ValidateZATCAQRCode(qrCodeInput.trim())
      validationResult = isValid

      if (isValid) {
        // Get QR code information
        const info = await GetQRCodeInfo(qrCodeInput.trim())
        qrInfo = info
      }
    } catch (err) {
      error = err.toString()
      validationResult = false
    } finally {
      isLoading = false
    }
  }

  async function generateSample() {
    isLoading = true
    error = null

    try {
      const sampleQR = await GenerateSampleQRCode()
      qrCodeInput = sampleQR
      await validateQRCode()
    } catch (err) {
      error = err.toString()
    } finally {
      isLoading = false
    }
  }

  function clearValidation() {
    qrCodeInput = ''
    validationResult = null
    qrInfo = null
    error = null
  }
</script>

<div class="container mx-auto max-w-4xl p-4">
  <!-- Header -->
  <div class="text-center mb-8">
    <h1 class="text-4xl font-bold text-white mb-2">ZATCA QR Code Validation</h1>
    <p class="text-white/80">Validate and analyze ZATCA compliant QR codes</p>
  </div>

  <!-- Input Section -->
  <div class="card bg-base-100 shadow-2xl mb-6">
    <div class="card-body">
      <h2 class="card-title text-2xl mb-4">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z"></path>
        </svg>
        QR Code Input
      </h2>
      
      <div class="form-control w-full">
        <label class="label" for="qr-input">
          <span class="label-text font-semibold">Base64 Encoded QR Code</span>
          <span class="label-text-alt">Paste your ZATCA QR code here</span>
        </label>
        <textarea
          id="qr-input"
          class="textarea textarea-bordered h-32 font-mono text-sm"
          placeholder="Paste your Base64 encoded ZATCA QR code here..."
          bind:value={qrCodeInput}
          disabled={isLoading}
        ></textarea>
      </div>
      
      <div class="card-actions justify-end mt-4 gap-2">
        <button 
          class="btn btn-outline btn-secondary"
          on:click={clearValidation}
          disabled={isLoading}
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
          </svg>
          Clear
        </button>
        
        <button 
          class="btn btn-info"
          on:click={generateSample}
          disabled={isLoading}
        >
          {#if isLoading}
            <span class="loading loading-spinner loading-sm"></span>
          {:else}
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4"></path>
            </svg>
          {/if}
          Generate Sample
        </button>
        
        <button 
          class="btn btn-primary"
          on:click={validateQRCode}
          disabled={isLoading || !qrCodeInput.trim()}
        >
          {#if isLoading}
            <span class="loading loading-spinner loading-sm"></span>
            Validating...
          {:else}
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
            Validate QR Code
          {/if}
        </button>
      </div>
    </div>
  </div>

  <!-- Error Alert -->
  {#if error}
    <div class="alert alert-error mb-6">
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
      </svg>
      <div>
        <h3 class="font-bold">Validation Error</h3>
        <div class="text-xs font-mono">{error}</div>
      </div>
    </div>
  {/if}

  <!-- Validation Results -->
  {#if validationResult !== null}
    <div class="card bg-base-100 shadow-2xl mb-6">
      <div class="card-body">
        <div class="flex items-center gap-3 mb-4">
          {#if validationResult}
            <div class="badge badge-success badge-lg">
              <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
              </svg>
              Valid
            </div>
            <h2 class="card-title text-success">✅ Valid ZATCA QR Code</h2>
          {:else}
            <div class="badge badge-error badge-lg">
              <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
              Invalid
            </div>
            <h2 class="card-title text-error">❌ Invalid QR Code</h2>
          {/if}
        </div>

        {#if validationResult && qrInfo}
          <!-- QR Code Information -->
          <div class="collapse collapse-open bg-base-200 mb-4">
            <div class="collapse-title text-xl font-medium">
              <svg class="w-5 h-5 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
              QR Code Information
            </div>
            <div class="collapse-content">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div class="stat bg-base-100 rounded-lg">
                  <div class="stat-title">Seller Name</div>
                  <div class="stat-value text-lg">{qrInfo.seller_name || 'N/A'}</div>
                </div>
                <div class="stat bg-base-100 rounded-lg">
                  <div class="stat-title">VAT Number</div>
                  <div class="stat-value text-lg font-mono">{qrInfo.vat_number || 'N/A'}</div>
                </div>
                <div class="stat bg-base-100 rounded-lg">
                  <div class="stat-title">Timestamp</div>
                  <div class="stat-value text-lg">{qrInfo.timestamp || 'N/A'}</div>
                </div>
                <div class="stat bg-base-100 rounded-lg">
                  <div class="stat-title">Total Amount</div>
                  <div class="stat-value text-lg text-success">{qrInfo.total_amount ? `${qrInfo.total_amount} SAR` : 'N/A'}</div>
                </div>
                <div class="stat bg-base-100 rounded-lg">
                  <div class="stat-title">VAT Amount</div>
                  <div class="stat-value text-lg text-warning">{qrInfo.vat_amount ? `${qrInfo.vat_amount} SAR` : 'N/A'}</div>
                </div>
                <div class="stat bg-base-100 rounded-lg">
                  <div class="stat-title">Compliance Status</div>
                  <div class="stat-value text-lg">{qrInfo.compliance || 'N/A'}</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Compliance Checks -->
          <div class="collapse collapse-open bg-base-200 mb-4">
            <div class="collapse-title text-xl font-medium">
              <svg class="w-5 h-5 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
              ZATCA Compliance Checks
            </div>
            <div class="collapse-content">
              <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                <div class="flex items-center gap-2 p-3 bg-base-100 rounded-lg">
                  <div class="badge badge-success">✓</div>
                  <span>TLV Format Structure</span>
                </div>
                <div class="flex items-center gap-2 p-3 bg-base-100 rounded-lg">
                  <div class="badge badge-success">✓</div>
                  <span>Base64 Encoding</span>
                </div>
                <div class="flex items-center gap-2 p-3 bg-base-100 rounded-lg">
                  <div class="badge badge-success">✓</div>
                  <span>Required Fields Present</span>
                </div>
                <div class="flex items-center gap-2 p-3 bg-base-100 rounded-lg">
                  <div class="badge badge-success">✓</div>
                  <span>Character Limit Compliance</span>
                </div>
                <div class="flex items-center gap-2 p-3 bg-base-100 rounded-lg">
                  <div class="badge {qrInfo.hash_present ? 'badge-success' : 'badge-error'}">{qrInfo.hash_present ? '✓' : '✗'}</div>
                  <span>Cryptographic Stamp</span>
                </div>
                <div class="flex items-center gap-2 p-3 bg-base-100 rounded-lg">
                  <div class="badge badge-success">✓</div>
                  <span>ZATCA Standard</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Technical Details -->
          <div class="collapse collapse-open bg-base-200">
            <div class="collapse-title text-xl font-medium">
              <svg class="w-5 h-5 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
              </svg>
              Technical Details
            </div>
            <div class="collapse-content">
              <div class="stats stats-vertical lg:stats-horizontal shadow">
                <div class="stat">
                  <div class="stat-title">QR Code Length</div>
                  <div class="stat-value text-primary">{qrCodeInput.length}</div>
                  <div class="stat-desc">characters</div>
                </div>
                <div class="stat">
                  <div class="stat-title">Standard</div>
                  <div class="stat-value text-secondary">ZATCA</div>
                  <div class="stat-desc">E-Invoicing</div>
                </div>
                <div class="stat">
                  <div class="stat-title">Encoding</div>
                  <div class="stat-value text-accent">Base64</div>
                  <div class="stat-desc">+ TLV Format</div>
                </div>
              </div>
            </div>
          </div>
        {/if}
      </div>
    </div>
  {/if}
</div>