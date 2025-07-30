<script>
  // Payments functionality will be implemented here
  let payments = []
  let isLoading = false
</script>

<div class="p-6">
  <div class="card bg-base-100/20 backdrop-blur-lg border border-white/20 shadow-xl">
    <div class="card-body">
      <h2 class="card-title text-white text-2xl mb-4">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"></path>
        </svg>
        Payments
        <div class="ml-auto">
          <div class="w-4 h-4 bg-accent rounded-full flex items-center justify-center">
            <svg class="w-2 h-2 text-white" fill="currentColor" viewBox="0 0 8 8">
              <circle cx="4" cy="4" r="3"/>
            </svg>
          </div>
        </div>
      </h2>
      
      <div class="flex justify-between items-center mb-6">
        <div class="flex gap-2">
          <button class="btn btn-primary">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
            </svg>
            Record Payment
          </button>
          <button class="btn btn-outline btn-secondary">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
            </svg>
            Reports
          </button>
        </div>
        <div class="form-control">
          <input type="text" placeholder="Search payments..." class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/60" />
        </div>
      </div>

      <div class="overflow-x-auto">
        <table class="table table-zebra">
          <thead>
            <tr class="text-white/80">
              <th>Payment ID</th>
              <th>Customer</th>
              <th>Invoice</th>
              <th>Amount</th>
              <th>Method</th>
              <th>Date</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody class="text-white/70">
            <tr>
              <td colspan="7" class="text-center py-8">
                <div class="flex flex-col items-center gap-2">
                  <svg class="w-12 h-12 text-white/40" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"></path>
                  </svg>
                  <p class="text-white/60">No payments recorded</p>
                  <button class="btn btn-primary btn-sm">Record your first payment</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</div>