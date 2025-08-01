<script>
	import { onMount } from 'svelte';
	import { GetPayments, CreatePayment, UpdatePayment, DeletePayment, GetPaymentTypes, GetInvoices } from '../wailsjs/go/main/App.js';
	import { database } from '../wailsjs/go/models';
	import PageLayout from './components/PageLayout.svelte';
	import DataTable from './components/DataTable.svelte';
	import StatusBadge from './components/StatusBadge.svelte';
	import PaymentModal from './components/PaymentModal.svelte';
	import ConfirmationModal from './components/ConfirmationModal.svelte';
	import { showDbSuccess, showDbError } from './stores/notificationStore.js';

	/** @type {Array<any>} */
	let payments = [];
	/** @type {Array<any>} */
	let paymentTypes = [];
	/** @type {Array<any>} */
	let invoices = [];
	/** @type {string} */
	let searchTerm = '';
	/** @type {boolean} */
	let showModal = false;
	/** @type {any} */
	let editingPayment = null;
	/** @type {boolean} */
	let loading = false;
	let showDeleteConfirm = false;
	let paymentToDelete = null;
	let confirmLoading = false;

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
			showDbError('load', 'payment data', error);
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
		showModal = true;
	}

	function closeModal() {
		showModal = false;
		editingPayment = null;
	}

	async function handleSavePayment(event) {
		try {
			loading = true;
			const { payment, isEditing } = event.detail;
			
			const paymentObj = new database.Payment();
			paymentObj.invoice_id = payment.invoice_id;
			paymentObj.payment_type_id = payment.payment_type_id;
			paymentObj.amount = payment.amount;
			paymentObj.payment_date = new Date(payment.payment_date);
			paymentObj.reference = payment.reference;
			paymentObj.notes = payment.notes;
			paymentObj.notes_arabic = payment.notes_arabic;
			paymentObj.status = payment.status;

			if (isEditing) {
				paymentObj.id = editingPayment.id;
				await UpdatePayment(paymentObj);
				showDbSuccess('update', 'Payment');
			} else {
				await CreatePayment(paymentObj);
				showDbSuccess('create', 'Payment');
			}

			await loadPayments();
			closeModal();
		} catch (error) {
			console.error('Error saving payment:', error);
			const operation = isEditing ? 'update' : 'create';
			showDbError(operation, 'Payment', error);
		} finally {
			loading = false;
		}
	}

	function showDeleteConfirmation(id) {
		paymentToDelete = id;
		showDeleteConfirm = true;
	}

	function cancelDelete() {
		showDeleteConfirm = false;
		paymentToDelete = null;
	}

	async function confirmDelete() {
		if (paymentToDelete) {
			confirmLoading = true;
			try {
				await DeletePayment(paymentToDelete);
				await loadPayments();
				showDbSuccess('delete', 'Payment');
				showDeleteConfirm = false;
				paymentToDelete = null;
			} catch (error) {
				console.error('Error deleting payment:', error);
				showDbError('delete', 'Payment', error);
			} finally {
				confirmLoading = false;
			}
		}
	}

	async function handleRefundPayment(payment) {
		console.log('handleRefundPayment called with payment:', payment);
		
		// Confirm refund action
		const confirmed = confirm(`Are you sure you want to refund payment #${payment.id}? This will mark the payment as refunded.`);
		if (!confirmed) return;
		
		// Get refund reason
		const reason = prompt('Please enter a reason for the refund (optional):') || 'Payment refunded';
		
		try {
			loading = true;
			
			// Create updated payment object with refunded status
			const updatedPayment = new database.Payment();
			updatedPayment.id = payment.id;
			updatedPayment.invoice_id = payment.invoice_id;
			updatedPayment.payment_type_id = payment.payment_type_id;
			updatedPayment.amount = payment.amount;
			updatedPayment.payment_date = new Date(payment.payment_date);
			updatedPayment.reference = payment.reference;
			updatedPayment.notes = payment.notes ? `${payment.notes}\n\nRefund: ${reason}` : `Refund: ${reason}`;
			updatedPayment.notes_arabic = payment.notes_arabic;
			updatedPayment.status = 'refunded';
			
			await UpdatePayment(updatedPayment);
			showDbSuccess('update', `Payment #${payment.id} refunded`);
			await loadPayments(); // Refresh the list
		} catch (error) {
			console.error('Error refunding payment:', error);
			showDbError('update', 'Payment refund', error);
		} finally {
			loading = false;
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

	// Table configuration
	/** @type {Array<{label: string, key?: string, class?: string, render?: Function, actions?: Array<{key: string, text: string, icon?: string, class?: string, title?: string}>}>} */
	const columns = [
		{ key: 'id', label: 'Payment ID', render: (item) => `#${item.id}` },
		{ key: 'invoice_id', label: 'Invoice', render: (item) => getInvoiceNumber(item.invoice_id) },
		{ key: 'amount', label: 'Amount', render: (item) => `< class: 'btn-warning', title: 'Refund Payment' },
				{ key: 'delete', text: 'Delete', icon: 'delete', class: 'btn-error', title: 'Delete Payment' }
			]
		}
	];

	/** @type {{text: string, icon: string}} */
	const primaryAction = {
		text: 'Record Payment',
		icon: 'add'
	};

	/** @type {Array<{text: string, icon: string}>} */span class="font-semibold">${formatCurrency(item.amount)}</span>`, class: 'font-semibold' },
		{ key: 'payment_type_id', label: 'Method', render: (item) => getPaymentTypeName(item.payment_type_id) },
		{ key: 'payment_date', label: 'Date', render: (item) => formatDate(item.payment_date) },
		{ key: 'reference', label: 'Reference', render: (item) => item.reference || '-' },
		{ key: 'status', label: 'Status' },
		{
			label: 'Actions',
			actions: [
				{ key: 'edit', text: 'Edit', icon: 'edit', class: 'btn-secondary', title: 'Edit Payment' },
				{ key: 'refund', text: 'Refund', icon: 'undo',
	const secondaryActions = [
		{
			text: 'Reports',
			icon: 'report'
		}
	];

	function handlePrimaryAction() {
		openModal();
	}

	function handleSecondaryAction(action) {
		if (action.text === 'Reports') {
			console.log('Reports clicked');
		}
	}

	function handleRowAction(event) {
		const { action, item } = event.detail;
		if (action === 'edit') {
			openModal(item);
		} else if (action === 'refund') {
			handleRefundPayment(item);
		} else if (action === 'delete') {
			showDeleteConfirmation(item.id);
		}
	}

	function handleSearch(event) {
		searchTerm = event.detail.searchTerm;
	}
</script>

<PageLayout 
	title="Payments" 
	icon="payment" 
	showIndicator={true}
>
	<svelte:fragment slot="actions">
		<!-- Actions are handled by DataTable -->
	</svelte:fragment>

	<DataTable
		data={filteredPayments}
		{columns}
		{loading}
		{searchTerm}
		searchPlaceholder="Search payments..."
		emptyStateTitle="No payments recorded"
		emptyStateMessage="Start by recording your first payment"
		emptyStateIcon="payment"
		{primaryAction}
		{secondaryActions}
		on:primary-action={handlePrimaryAction}
		on:secondary-action={handleSecondaryAction}
		on:row-action={handleRowAction}
		on:search={handleSearch}
	>
		<svelte:fragment slot="cell" let:item let:column>
			{#if column.key === 'status'}
				<StatusBadge status={item.status} />
			{:else if column.render}
				{@html column.render(item)}
			{:else}
				{item[column.key] || '-'}
			{/if}
		</svelte:fragment>
	</DataTable>
</PageLayout>

<PaymentModal
	bind:show={showModal}
	{editingPayment}
	{loading}
	{invoices}
	{paymentTypes}
	on:save={handleSavePayment}
	on:close={closeModal}
/>

{#if showDeleteConfirm}
	<ConfirmationModal
		show={showDeleteConfirm}
		title="Delete Payment"
		message="Are you sure you want to delete this payment? This action cannot be undone."
		confirmText="Delete"
		cancelText="Cancel"
		confirmClass="btn-error"
		loading={confirmLoading}
		on:confirm={confirmDelete}
		on:cancel={cancelDelete}
	/>
{/if}