<script>
	import { onMount } from 'svelte';
	import { GetPayments, CreatePayment, UpdatePayment, DeletePayment, GetPaymentTypes, GetInvoices } from '../wailsjs/go/main/App.js';
	import { database } from '../wailsjs/go/models';

	let payments = [];
	let paymentTypes = [];
	let invoices = [];
	let searchTerm = '';
	let showModal = false;
	let editingPayment = null;
	let loading = false;

	// Form data
	let paymentData = {
		invoice_id: 0,
		payment_type_id: 0,
		amount: 0,
		payment_date: new Date().toISOString().split('T')[0],
		reference: '',
		notes: '',
		notes_arabic: '',
		status: 'completed'
	};

	onMount(async () => {
		loading = true;
		try {
			await Promise.all([
				loadPayments(),
				loadPaymentTypes(),
				loadInvoices()
			]);
		} catch (error) {
			console.error('Error loading payment data:', error);
		} finally {
			loading = false;
		}
	});

	async function loadPayments() {
		try {
			console.log('Loading payments...');
			payments = await GetPayments();
			console.log('Payments loaded:', payments.length);
		} catch (error) {
			console.error('Error loading payments:', error);
			payments = []; // Ensure payments is always an array
		}
	}

	async function loadPaymentTypes() {
		try {
			console.log('Loading payment types...');
			paymentTypes = await GetPaymentTypes();
			console.log('Payment types loaded:', paymentTypes.length);
		} catch (error) {
			console.error('Error loading payment types:', error);
			paymentTypes = []; // Ensure paymentTypes is always an array
		}
	}

	async function loadInvoices() {
		try {
			console.log('Loading invoices...');
			invoices = await GetInvoices();
			console.log('Invoices loaded:', invoices.length);
		} catch (error) {
			console.error('Error loading invoices:', error);
			invoices = []; // Ensure invoices is always an array
		}
	}

	function openModal(payment = null) {
		editingPayment = payment;
		if (payment) {
			paymentData = {
				invoice_id: payment.invoice_id,
				payment_type_id: payment.payment_type_id,
				amount: payment.amount,
				payment_date: new Date(payment.payment_date).toISOString().split('T')[0],
				reference: payment.reference || '',
				notes: payment.notes || '',
				notes_arabic: payment.notes_arabic || '',
				status: payment.status || 'completed'
			};
		} else {
			resetForm();
		}
		showModal = true;
	}

	function closeModal() {
		showModal = false;
		editingPayment = null;
		resetForm();
	}

	function resetForm() {
		paymentData = {
			invoice_id: 0,
			payment_type_id: 0,
			amount: 0,
			payment_date: new Date().toISOString().split('T')[0],
			reference: '',
			notes: '',
			notes_arabic: '',
			status: 'completed'
		};
	}

	async function savePayment() {
		try {
			loading = true;
			
			const payment = new database.Payment();
			payment.invoice_id = Number(paymentData.invoice_id);
			payment.payment_type_id = Number(paymentData.payment_type_id);
			payment.amount = Number(paymentData.amount);
			payment.payment_date = new Date(paymentData.payment_date);
			payment.reference = paymentData.reference;
			payment.notes = paymentData.notes;
			payment.notes_arabic = paymentData.notes_arabic;
			payment.status = paymentData.status;

			if (editingPayment) {
				payment.id = editingPayment.id;
				await UpdatePayment(payment);
			} else {
				await CreatePayment(payment);
			}

			await loadPayments();
			closeModal();
		} catch (error) {
			console.error('Error saving payment:', error);
			alert('Error saving payment: ' + error.message);
		} finally {
			loading = false;
		}
	}

	async function deletePayment(id) {
		if (confirm('Are you sure you want to delete this payment?')) {
			try {
				loading = true;
				await DeletePayment(id);
				await loadPayments();
			} catch (error) {
				console.error('Error deleting payment:', error);
				alert('Error deleting payment: ' + error.message);
			} finally {
				loading = false;
			}
		}
	}

	function formatDate(dateString) {
		return new Date(dateString).toLocaleDateString();
	}

	function formatCurrency(amount) {
		return new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'SAR'
		}).format(amount);
	}

	function getInvoiceNumber(invoiceId) {
		if (!invoices || invoices.length === 0) return `Invoice #${invoiceId}`;
		const invoice = invoices.find(inv => inv.id === invoiceId);
		return invoice ? invoice.invoice_number : `Invoice #${invoiceId}`;
	}

	function getPaymentTypeName(paymentTypeId) {
		if (!paymentTypes || paymentTypes.length === 0) return 'Unknown';
		const paymentType = paymentTypes.find(pt => pt.id === paymentTypeId);
		return paymentType ? paymentType.name : 'Unknown';
	}

	function getStatusBadgeClass(status) {
		switch (status) {
			case 'completed': return 'badge-success';
			case 'pending': return 'badge-warning';
			case 'failed': return 'badge-error';
			case 'cancelled': return 'badge-neutral';
			default: return 'badge-neutral';
		}
	}

	$: filteredPayments = (() => {
		console.log('Filtering payments:', { 
			paymentsLength: payments.length, 
			invoicesLength: invoices.length, 
			paymentTypesLength: paymentTypes.length,
			searchTerm,
			loading 
		});
		
		if (!payments || payments.length === 0) {
			return [];
		}
		
		if (!searchTerm) {
			return payments;
		}
		
		return payments.filter(payment => {
			const searchLower = searchTerm.toLowerCase();
			const invoiceNumber = getInvoiceNumber(payment.invoice_id).toLowerCase();
			const paymentTypeName = getPaymentTypeName(payment.payment_type_id).toLowerCase();
			const reference = (payment.reference || '').toLowerCase();
			const notes = (payment.notes || '').toLowerCase();
			
			return invoiceNumber.includes(searchLower) ||
				   paymentTypeName.includes(searchLower) ||
				   reference.includes(searchLower) ||
				   notes.includes(searchLower) ||
				   payment.status.toLowerCase().includes(searchLower);
		});
	})();
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
					<button class="btn btn-primary" on:click={() => openModal()}>
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
					<input 
						type="text" 
						placeholder="Search payments..." 
						class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/60"
						bind:value={searchTerm}
					/>
				</div>
			</div>

			<div class="overflow-x-auto">
				<table class="table table-zebra">
					<thead>
						<tr class="text-white/80">
							<th>Payment ID</th>
							<th>Invoice</th>
							<th>Amount</th>
							<th>Method</th>
							<th>Date</th>
							<th>Reference</th>
							<th>Status</th>
							<th>Actions</th>
						</tr>
					</thead>
					<tbody class="text-white/70">
						{#if loading}
							<tr>
								<td colspan="8" class="text-center py-8">
									<span class="loading loading-spinner loading-md"></span>
									Loading payments...
								</td>
							</tr>
						{:else if filteredPayments.length === 0}
							<tr>
								<td colspan="8" class="text-center py-8">
									<div class="flex flex-col items-center gap-2">
										<svg class="w-12 h-12 text-white/40" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"></path>
										</svg>
										<p class="text-white/60">
											{searchTerm ? 'No payments found matching your search.' : 'No payments recorded'}
										</p>
										{#if !searchTerm}
											<button class="btn btn-primary btn-sm" on:click={() => openModal()}>
												Record your first payment
											</button>
										{/if}
									</div>
								</td>
							</tr>
						{:else}
							{#each filteredPayments as payment (payment.id)}
								<tr>
									<td>#{payment.id}</td>
									<td>{getInvoiceNumber(payment.invoice_id)}</td>
									<td class="font-semibold">{formatCurrency(payment.amount)}</td>
									<td>{getPaymentTypeName(payment.payment_type_id)}</td>
									<td>{formatDate(payment.payment_date)}</td>
									<td>{payment.reference || '-'}</td>
									<td>
										<span class="badge {getStatusBadgeClass(payment.status)}">
											{payment.status}
										</span>
									</td>
									<td>
										<div class="flex gap-2">
											<button 
												class="btn btn-sm btn-ghost text-white/70 hover:text-white"
												on:click={() => openModal(payment)}
												title="Edit Payment"
											>
												<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
												</svg>
											</button>
											<button 
												class="btn btn-sm btn-ghost text-red-400 hover:text-red-300"
												on:click={() => deletePayment(payment.id)}
												title="Delete Payment"
											>
												<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
													<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
												</svg>
											</button>
										</div>
									</td>
								</tr>
							{/each}
						{/if}
					</tbody>
				</table>
			</div>
		</div>
	</div>
</div>

{#if showModal}
	<div class="modal modal-open">
		<div class="modal-box max-w-2xl bg-base-100/90 backdrop-blur-lg border border-white/20">
			<h3 class="font-bold text-lg mb-4 text-white">
				{editingPayment ? 'Edit Payment' : 'Record New Payment'}
			</h3>

			<form on:submit|preventDefault={savePayment}>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="form-control">
						<label class="label" for="invoice_id">
							<span class="label-text text-white/80">Invoice *</span>
						</label>
						<select 
							id="invoice_id"
							class="select select-bordered bg-white/10 border-white/20 text-white"
							bind:value={paymentData.invoice_id}
							required
						>
							<option value={0}>Select Invoice</option>
							{#each invoices as invoice}
								<option value={invoice.id}>
									{invoice.invoice_number} - {formatCurrency(invoice.total_amount)}
								</option>
							{/each}
						</select>
					</div>

					<div class="form-control">
						<label class="label" for="payment_type_id">
							<span class="label-text text-white/80">Payment Method *</span>
						</label>
						<select 
							id="payment_type_id"
							class="select select-bordered bg-white/10 border-white/20 text-white"
							bind:value={paymentData.payment_type_id}
							required
						>
							<option value={0}>Select Payment Method</option>
							{#each paymentTypes as paymentType}
								<option value={paymentType.id}>{paymentType.name}</option>
							{/each}
						</select>
					</div>

					<div class="form-control">
						<label class="label" for="amount">
							<span class="label-text text-white/80">Amount *</span>
						</label>
						<input 
							id="amount"
							type="number" 
							step="0.01"
							min="0"
							class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/60"
							bind:value={paymentData.amount}
							required
						/>
					</div>

					<div class="form-control">
						<label class="label" for="payment_date">
							<span class="label-text text-white/80">Payment Date *</span>
						</label>
						<input 
							id="payment_date"
							type="date"
							class="input input-bordered bg-white/10 border-white/20 text-white"
							bind:value={paymentData.payment_date}
							required
						/>
					</div>

					<div class="form-control">
						<label class="label" for="reference">
							<span class="label-text text-white/80">Reference</span>
						</label>
						<input 
							id="reference"
							type="text"
							class="input input-bordered bg-white/10 border-white/20 text-white placeholder-white/60"
							placeholder="Check number, transaction ID, etc."
							bind:value={paymentData.reference}
						/>
					</div>

					<div class="form-control">
						<label class="label" for="status">
							<span class="label-text text-white/80">Status *</span>
						</label>
						<select 
							id="status"
							class="select select-bordered bg-white/10 border-white/20 text-white"
							bind:value={paymentData.status}
							required
						>
							<option value="completed">Completed</option>
							<option value="pending">Pending</option>
							<option value="failed">Failed</option>
							<option value="cancelled">Cancelled</option>
						</select>
					</div>

					<div class="form-control md:col-span-2">
						<label class="label" for="notes">
							<span class="label-text text-white/80">Notes</span>
						</label>
						<textarea 
							id="notes"
							class="textarea textarea-bordered bg-white/10 border-white/20 text-white placeholder-white/60"
							placeholder="Additional notes about this payment"
							bind:value={paymentData.notes}
						></textarea>
					</div>

					<div class="form-control md:col-span-2">
						<label class="label" for="notes_arabic">
							<span class="label-text text-white/80">Notes (Arabic)</span>
						</label>
						<textarea 
							id="notes_arabic"
							class="textarea textarea-bordered bg-white/10 border-white/20 text-white placeholder-white/60"
							placeholder="ملاحظات إضافية حول هذه الدفعة"
							bind:value={paymentData.notes_arabic}
							dir="rtl"
						></textarea>
					</div>
				</div>

				<div class="modal-action">
					<button type="button" class="btn btn-ghost text-white/70" on:click={closeModal}>Cancel</button>
					<button 
						type="submit" 
						class="btn btn-primary"
						disabled={loading}
					>
						{#if loading}
							<span class="loading loading-spinner loading-sm"></span>
						{/if}
						{editingPayment ? 'Update Payment' : 'Record Payment'}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<style>
	.badge-success {
		@apply bg-green-100 text-green-800 border-green-200;
	}

	.badge-warning {
		@apply bg-yellow-100 text-yellow-800 border-yellow-200;
	}

	.badge-error {
		@apply bg-red-100 text-red-800 border-red-200;
	}

	.badge-neutral {
		@apply bg-gray-100 text-gray-800 border-gray-200;
	}
</style>