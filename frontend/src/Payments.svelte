<script>
	import { onMount } from 'svelte';
	import { GetPayments, CreatePayment, UpdatePayment, DeletePayment, GetPaymentTypes, GetInvoices } from '../wailsjs/go/main/App.js';
	import { database } from '../wailsjs/go/models';
	import PageLayout from './components/PageLayout.svelte';
	import DataTable from './components/DataTable.svelte';
	import StatusBadge from './components/StatusBadge.svelte';
	import PaymentModal from './components/PaymentModal.svelte';

	let payments = [];
	let paymentTypes = [];
	let invoices = [];
	let searchTerm = '';
	let showModal = false;
	let editingPayment = null;
	let loading = false;

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
			} else {
				await CreatePayment(paymentObj);
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
	const columns = [
		{ key: 'id', label: 'Payment ID', render: (item) => `#${item.id}` },
		{ key: 'invoice_id', label: 'Invoice', render: (item) => getInvoiceNumber(item.invoice_id) },
		{ key: 'amount', label: 'Amount', render: (item) => `<span class="font-semibold">${formatCurrency(item.amount)}</span>`, class: 'font-semibold' },
		{ key: 'payment_type_id', label: 'Method', render: (item) => getPaymentTypeName(item.payment_type_id) },
		{ key: 'payment_date', label: 'Date', render: (item) => formatDate(item.payment_date) },
		{ key: 'reference', label: 'Reference', render: (item) => item.reference || '-' },
		{ 
			key: 'status', 
			label: 'Status',
			actions: [
				{ key: 'edit', icon: 'fa-edit', class: 'btn-secondary', title: 'Edit Payment' },
				{ key: 'delete', icon: 'fa-trash', class: 'btn-error', title: 'Delete Payment' }
			]
		}
	];

	const primaryAction = {
		text: 'Record Payment',
		icon: 'fa-plus'
	};

	const secondaryActions = [
		{
			text: 'Reports',
			icon: 'fa-chart-bar'
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
		} else if (action === 'delete') {
			deletePayment(item.id);
		}
	}

	function handleSearch(event) {
		searchTerm = event.detail.searchTerm;
	}
</script>

<PageLayout 
	title="Payments" 
	icon="fa-credit-card" 
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
		emptyStateIcon="fa-credit-card"
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