<script>
  import { PopulateSampleData } from '../wailsjs/go/main/App.js'
  
  // Props
  export let showIntro = () => {}
  
  // Administration functionality will be implemented here
  let systemInfo = {
    version: '1.0.0',
    uptime: '2 days, 14 hours',
    database: 'Connected',
    zatcaStatus: 'Active'
  }
  let isLoading = false
  let sampleDataLoading = false
  let sampleDataMessage = ''

  async function populateSampleData() {
    sampleDataLoading = true
    sampleDataMessage = ''
    
    try {
      await PopulateSampleData()
      sampleDataMessage = 'Sample data populated successfully! Created 5 customers, 5 products, 5 payment types, and 5 invoices.'
    } catch (error) {
      sampleDataMessage = `Error: ${error}`
    } finally {
      sampleDataLoading = false
    }
  }
</script>

<div class="p-6">
  <div class="card bg-base-100/20 backdrop-blur-lg border border-white/20 shadow-xl">
    <div class="card-body">
      <h2 class="card-title text-white text-2xl mb-4">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path>
        </svg>
        Administration
      </h2>
      
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- System Status -->
        <div class="card bg-base-100/10 border border-white/10">
          <div class="card-body">
            <h3 class="card-title text-white text-lg mb-4">System Status</h3>
            
            <div class="stats stats-vertical shadow bg-base-100/10">
              <div class="stat">
                <div class="stat-title text-white/60">Application Version</div>
                <div class="stat-value text-primary text-lg">{systemInfo.version}</div>
              </div>
              
              <div class="stat">
                <div class="stat-title text-white/60">System Uptime</div>
                <div class="stat-value text-success text-lg">{systemInfo.uptime}</div>
              </div>
              
              <div class="stat">
                <div class="stat-title text-white/60">Database Status</div>
                <div class="stat-value text-success text-lg">{systemInfo.database}</div>
              </div>
              
              <div class="stat">
                <div class="stat-title text-white/60">ZATCA Integration</div>
                <div class="stat-value text-success text-lg">{systemInfo.zatcaStatus}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- System Actions -->
        <div class="card bg-base-100/10 border border-white/10">
          <div class="card-body">
            <h3 class="card-title text-white text-lg mb-4">System Actions</h3>
            
            <div class="space-y-3">
              <button 
                class="btn btn-outline btn-success w-full" 
                class:loading={sampleDataLoading}
                disabled={sampleDataLoading}
                on:click={populateSampleData}
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
                </svg>
                {sampleDataLoading ? 'Populating...' : 'Populate Sample Data (5 each)'}
              </button>
              
              <button class="btn btn-outline btn-primary w-full" on:click={showIntro}>
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
                Show App Introduction
              </button>
              
              <button class="btn btn-outline btn-info w-full">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4-4m0 0L8 8m4-4v12"></path>
                </svg>
                Backup Database
              </button>
              
              <button class="btn btn-outline btn-warning w-full">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
                </svg>
                Clear Cache
              </button>
              
              <button class="btn btn-outline btn-secondary w-full">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
                </svg>
                Generate Reports
              </button>
              
              <button class="btn btn-outline btn-error w-full">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
                </svg>
                System Maintenance
              </button>
            </div>

            <!-- Sample Data Message -->
            {#if sampleDataMessage}
              <div class="alert mt-4" class:alert-success={!sampleDataMessage.includes('Error')} class:alert-error={sampleDataMessage.includes('Error')}>
                <div>
                  <svg class="w-6 h-6 shrink-0 stroke-current" fill="none" viewBox="0 0 24 24">
                    {#if sampleDataMessage.includes('Error')}
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    {:else}
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    {/if}
                  </svg>
                  <span>{sampleDataMessage}</span>
                </div>
              </div>
            {/if}
          </div>
        </div>
      </div>

      <!-- Logs Section -->
      <div class="card bg-base-100/10 border border-white/10 mt-6">
        <div class="card-body">
          <h3 class="card-title text-white text-lg mb-4">Recent System Logs</h3>
          
          <div class="overflow-x-auto">
            <table class="table table-zebra">
              <thead>
                <tr class="text-white/80">
                  <th>Timestamp</th>
                  <th>Level</th>
                  <th>Message</th>
                  <th>User</th>
                </tr>
              </thead>
              <tbody class="text-white/70">
                <tr>
                  <td colspan="4" class="text-center py-8">
                    <div class="flex flex-col items-center gap-2">
                      <svg class="w-12 h-12 text-white/40" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                      </svg>
                      <p class="text-white/60">No recent logs</p>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>