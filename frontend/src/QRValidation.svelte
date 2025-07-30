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

<div class="qr-validation-container">
  <h2>ZATCA QR Code Validation</h2>
  
  <div class="input-section">
    <label for="qr-input">Base64 Encoded QR Code:</label>
    <textarea
      id="qr-input"
      bind:value={qrCodeInput}
      placeholder="Paste your Base64 encoded ZATCA QR code here..."
      rows="4"
      disabled={isLoading}
    ></textarea>
    
    <div class="action-buttons">
      <button 
        class="btn validate-btn" 
        on:click={validateQRCode}
        disabled={isLoading || !qrCodeInput.trim()}
      >
        {isLoading ? 'Validating...' : 'Validate QR Code'}
      </button>
      
      <button 
        class="btn sample-btn" 
        on:click={generateSample}
        disabled={isLoading}
      >
        Generate Sample
      </button>
      
      <button 
        class="btn clear-btn" 
        on:click={clearValidation}
        disabled={isLoading}
      >
        Clear
      </button>
    </div>
  </div>

  {#if error}
    <div class="error-message">
      <h3>❌ Error</h3>
      <p>{error}</p>
    </div>
  {/if}

  {#if validationResult !== null}
    <div class="validation-result {validationResult ? 'valid' : 'invalid'}">
      <h3>{validationResult ? '✅ Valid ZATCA QR Code' : '❌ Invalid QR Code'}</h3>
      
      {#if validationResult && qrInfo}
        <div class="qr-info">
          <h4>QR Code Information</h4>
          <div class="info-grid">
            <div class="info-item">
              <span class="label">Seller Name:</span>
              <span class="value">{qrInfo.seller_name || 'N/A'}</span>
            </div>
            <div class="info-item">
              <span class="label">VAT Number:</span>
              <span class="value">{qrInfo.vat_number || 'N/A'}</span>
            </div>
            <div class="info-item">
              <span class="label">Timestamp:</span>
              <span class="value">{qrInfo.timestamp || 'N/A'}</span>
            </div>
            <div class="info-item">
              <span class="label">Total Amount:</span>
              <span class="value">{qrInfo.total_amount ? `${qrInfo.total_amount} SAR` : 'N/A'}</span>
            </div>
            <div class="info-item">
              <span class="label">VAT Amount:</span>
              <span class="value">{qrInfo.vat_amount ? `${qrInfo.vat_amount} SAR` : 'N/A'}</span>
            </div>
            <div class="info-item">
              <span class="label">Cryptographic Stamp:</span>
              <span class="value">{qrInfo.hash_present ? '✅ Present' : '❌ Missing'}</span>
            </div>
            <div class="info-item">
              <span class="label">Compliance:</span>
              <span class="value">{qrInfo.compliance || 'N/A'}</span>
            </div>
          </div>
        </div>

        <div class="compliance-checks">
          <h4>ZATCA Compliance Checks</h4>
          <div class="check-list">
            <div class="check-item">
              <span class="check-icon">✅</span>
              <span>TLV Format Structure</span>
            </div>
            <div class="check-item">
              <span class="check-icon">✅</span>
              <span>Base64 Encoding</span>
            </div>
            <div class="check-item">
              <span class="check-icon">✅</span>
              <span>Required Fields Present</span>
            </div>
            <div class="check-item">
              <span class="check-icon">✅</span>
              <span>Character Limit Compliance</span>
            </div>
            <div class="check-item">
              <span class="check-icon">{qrInfo.hash_present ? '✅' : '❌'}</span>
              <span>Cryptographic Stamp</span>
            </div>
          </div>
        </div>

        <div class="technical-details">
          <h4>Technical Details</h4>
          <div class="tech-grid">
            <div class="tech-item">
              <span class="label">QR Code Length:</span>
              <span class="value">{qrCodeInput.length} characters</span>
            </div>
            <div class="tech-item">
              <span class="label">Standard:</span>
              <span class="value">ZATCA E-Invoicing</span>
            </div>
            <div class="tech-item">
              <span class="label">Encoding:</span>
              <span class="value">Base64 + TLV</span>
            </div>
          </div>
        </div>
      {/if}
    </div>
  {/if}
</div>

<style>
  .qr-validation-container {
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  }

  h2 {
    color: #2c3e50;
    text-align: center;
    margin-bottom: 30px;
    font-size: 2rem;
  }

  .input-section {
    background: #f8f9fa;
    padding: 20px;
    border-radius: 8px;
    margin-bottom: 20px;
  }

  label {
    display: block;
    margin-bottom: 8px;
    font-weight: 600;
    color: #495057;
  }

  textarea {
    width: 100%;
    padding: 12px;
    border: 2px solid #dee2e6;
    border-radius: 6px;
    font-family: 'Courier New', monospace;
    font-size: 14px;
    resize: vertical;
    transition: border-color 0.3s ease;
    box-sizing: border-box;
  }

  textarea:focus {
    outline: none;
    border-color: #007bff;
    box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.1);
  }

  textarea:disabled {
    background-color: #e9ecef;
    opacity: 0.6;
  }

  .action-buttons {
    display: flex;
    gap: 10px;
    margin-top: 15px;
    flex-wrap: wrap;
  }

  .btn {
    padding: 10px 20px;
    border: none;
    border-radius: 6px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
    font-size: 14px;
  }

  .btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .validate-btn {
    background: #28a745;
    color: white;
  }

  .validate-btn:hover:not(:disabled) {
    background: #218838;
    transform: translateY(-1px);
  }

  .sample-btn {
    background: #17a2b8;
    color: white;
  }

  .sample-btn:hover:not(:disabled) {
    background: #138496;
    transform: translateY(-1px);
  }

  .clear-btn {
    background: #6c757d;
    color: white;
  }

  .clear-btn:hover:not(:disabled) {
    background: #5a6268;
    transform: translateY(-1px);
  }

  .error-message {
    background: #f8d7da;
    border: 1px solid #f5c6cb;
    color: #721c24;
    padding: 15px;
    border-radius: 6px;
    margin-bottom: 20px;
  }

  .error-message h3 {
    margin: 0 0 10px 0;
    font-size: 1.1rem;
  }

  .error-message p {
    margin: 0;
    font-family: 'Courier New', monospace;
    font-size: 14px;
  }

  .validation-result {
    padding: 20px;
    border-radius: 8px;
    margin-bottom: 20px;
  }

  .validation-result.valid {
    background: #d4edda;
    border: 1px solid #c3e6cb;
    color: #155724;
  }

  .validation-result.invalid {
    background: #f8d7da;
    border: 1px solid #f5c6cb;
    color: #721c24;
  }

  .validation-result h3 {
    margin: 0 0 20px 0;
    font-size: 1.3rem;
  }

  .qr-info, .compliance-checks, .technical-details {
    background: white;
    padding: 15px;
    border-radius: 6px;
    margin-bottom: 15px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  }

  .qr-info h4, .compliance-checks h4, .technical-details h4 {
    margin: 0 0 15px 0;
    color: #495057;
    font-size: 1.1rem;
    border-bottom: 2px solid #e9ecef;
    padding-bottom: 8px;
  }

  .info-grid, .tech-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 10px;
  }

  .info-item, .tech-item {
    display: flex;
    justify-content: space-between;
    padding: 8px 0;
    border-bottom: 1px solid #f1f3f4;
  }

  .info-item:last-child, .tech-item:last-child {
    border-bottom: none;
  }

  .label {
    font-weight: 600;
    color: #6c757d;
  }

  .value {
    color: #495057;
    font-family: 'Courier New', monospace;
    font-size: 14px;
  }

  .check-list {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 10px;
  }

  .check-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px;
    background: #f8f9fa;
    border-radius: 4px;
  }

  .check-icon {
    font-size: 16px;
  }

  @media (max-width: 768px) {
    .qr-validation-container {
      padding: 15px;
    }

    .action-buttons {
      flex-direction: column;
    }

    .btn {
      width: 100%;
    }

    .info-grid, .tech-grid {
      grid-template-columns: 1fr;
    }

    .check-list {
      grid-template-columns: 1fr;
    }
  }
</style>