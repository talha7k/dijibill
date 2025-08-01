<script>
	import { onMount } from 'svelte';
	import { wailsReady } from './stores/wailsStore.js';
	import { GetCustomers, CreateCustomer, UpdateCustomer, DeleteCustomer } from '../wailsjs/go/main/App.js';
	import { database } from '../wailsjs/go/models';
	import PageLayout from './components/PageLayout.svelte';
	import DataTable from './components/DataTable.svelte';
	import CustomerModal from './components/CustomerModal.svelte';
	import ConfirmationModal from './components/ConfirmationModal.svelte';
	import StatusBadge from './components/StatusBadge.svelte';
	import { showDbSuccess, showDbError } from './stores/notificationStore.js';

	let customers = [];
	let loading = false;
	let showModal = false;
	let showDeleteConfirm = false;
	let editingCustomer = null;
	let customerToDelete = null;
	let searchTerm = '';

	// Reactive statement to load customers when Wails is ready
	$: if ($wailsReady) {
		loadCustomers();
	}

	async function loadCustomers() {
		console.log('📋 Loading customers...');
		loading = true;
		try {
			customers = await GetCustomers();
			console.log('✅ Customers loaded successfully:', customers.length, 'customers');
		} catch (error) {
			console.error('❌ Error loading customers:', error);
			showDbError('load', 'customers', error);
			customers = [];
		} finally {
			loading = false;
		}
	}

	function openModal(customer = null) {
		console.log('🚪 openModal called with customer:', customer);
		editingCustomer = customer;
		showModal = true;
		console.log('🚪 Modal state - editingCustomer:', editingCustomer, 'showModal:', showModal);
	}

	function closeModal() {
		console.log('🚪 closeModal called');
		showModal = false;
		editingCustomer = null;
	}

	async function handleSaveCustomer(event) {
		console.log('💾 handleSaveCustomer called with event:', event);
		console.log('💾 Event detail:', event.detail);
		const { customer: customerData, isEditing } = event.detail;
		console.log('💾 Customer data:', customerData, 'isEditing:', isEditing);
		loading = true;
		
		try {
			const customer = new database.Customer();
			Object.assign(customer, customerData);
			console.log('💾 Created customer object:', customer);

			if (isEditing && editingCustomer) {
				customer.id = editingCustomer.id;
				console.log('💾 Updating customer with ID:', customer.id);
				console.log('💾 Customer object before update:', customer);
				const updateResult = await UpdateCustomer(customer);
				console.log('✅ Customer updated successfully, result:', updateResult);
				showDbSuccess('update', 'Customer');
			} else {
				console.log('💾 Creating new customer');
				console.log('💾 Customer object before create:', customer);
				const createResult = await CreateCustomer(customer);
				console.log('✅ Customer created successfully, result:', createResult);
				showDbSuccess('create', 'Customer');
			}

			await loadCustomers();
			closeModal();
		} catch (error) {
			console.error('❌ Error saving customer:', error);
			const operation = isEditing ? 'update' : 'create';
			showDbError(operation, 'Customer', error);
		} finally {
			loading = false;
		}
	}

	function showDeleteConfirmation(customerId) {
		console.log('🗑️ showDeleteConfirmation called with ID:', customerId);
		console.log('🗑️ Customer ID type:', typeof customerId);
		console.log('🗑️ Wails runtime ready:', $wailsReady);
		console.log('🗑️ DeleteCustomer function available:', typeof DeleteCustomer);
		
		customerToDelete = customerId;
		showDeleteConfirm = true;
	}

	function cancelDelete() {
		console.log('🗑️ Delete cancelled');
		customerToDelete = null;
		showDeleteConfirm = false;
	}

	async function confirmDelete() {
		console.log('🗑️ Delete confirmed for customer ID:', customerToDelete);
		
		if (!customerToDelete) {
			console.error('❌ No customer ID to delete');
			return;
		}

		loading = true;
		
		try {
			// Ensure we have a valid ID
			const id = parseInt(customerToDelete);
			if (isNaN(id)) {
				throw new Error(`Invalid customer ID: ${customerToDelete}`);
			}
			
			console.log('🗑️ Calling DeleteCustomer API with parsed ID:', id);
			console.log('🗑️ About to call DeleteCustomer function...');
			
			const deleteResult = await DeleteCustomer(id);
			console.log('✅ Customer deleted successfully, result:', deleteResult);
			showDbSuccess('delete', 'Customer');
			
			console.log('🔄 Reloading customer list...');
			await loadCustomers();
			console.log('✅ Customer list reloaded after deletion');
			
			// Close the confirmation modal
			showDeleteConfirm = false;
			customerToDelete = null;
			
		} catch (error) {
			console.error('❌ Error deleting customer:', error);
			console.error('❌ Error details:', {
				message: error.message,
				stack: error.stack,
				name: error.name
			});
			showDbError('delete', 'Customer', error);
		} finally {
			loading = false;
		}
	}

	$: filteredCustomers = customers.filter(customer => {
		if (!searchTerm) return true;
		const searchLower = searchTerm.toLowerCase();
		return (
			customer.name?.toLowerCase().includes(searchLower) ||
			customer.name_arabic?.toLowerCase().includes(searchLower) ||
			customer.email?.toLowerCase().includes(searchLower) ||
			customer.phone?.toLowerCase().includes(searchLower) ||
			customer.company?.toLowerCase().includes(searchLower) ||
			customer.city?.toLowerCase().includes(searchLower)
		);
	});

	// Table configuration
	const columns = [
		{ key: 'name', label: 'Name', render: (item) => `<div><div class="font-medium">${item.name || '-'}</div><div class="text-sm opacity-70">${item.name_arabic || ''}</div></div>` },
		{ key: 'email', label: 'Email', render: (item) => item.email || '-' },
		{ key: 'phone', label: 'Phone', render: (item) => item.phone || '-' },
		{ key: 'company', label: 'Company', render: (item) => item.company || '-' },
		{ key: 'city', label: 'City', render: (item) => `${item.city || '-'}, ${item.country || ''}` },
		{ key: 'status', label: 'Status' },
		{ 
			label: 'Actions',
			actions: [
				{ key: 'edit', text: 'Edit', icon: 'edit', class: 'btn-secondary', title: 'Edit Customer' },
				{ key: 'delete', text: 'Delete', icon: 'delete', class: 'btn-error', title: 'Delete Customer' }
			]
		}
	];

	const primaryAction = {
		text: 'Add Customer',
		icon: 'add'
	};

	const secondaryActions = [
		{
			text: 'Import',
			icon: 'upload'
		},
		{
			text: 'Export',
			icon: 'download'
		}
	];

	function handlePrimaryAction() {
		openModal();
	}

	function handleSecondaryAction(action) {
		if (action.text === 'Import') {
			console.log('Import clicked');
		} else if (action.text === 'Export') {
			console.log('Export clicked');
		}
	}

	function handleRowAction(event) {
    console.log('🎯 Row action triggered:', event.detail);
    const { action, item } = event.detail;
    console.log('🎯 Action type:', action);
    console.log('🎯 Row data:', item);
    console.log('🎯 Item keys:', Object.keys(item));
    console.log('🎯 Item.id:', item.id);
    console.log('🎯 Item.ID:', item.ID);
    console.log('🎯 Full item structure:', JSON.stringify(item, null, 2));
    
    if (action === 'edit') {
      console.log('✏️ Edit action triggered for customer:', item);
      openModal(item);
    } else if (action === 'delete') {
      console.log('🗑️ Delete action triggered for customer:', item);
      // Try both possible ID field names
      const customerId = item.id || item.ID;
      console.log('🗑️ Resolved customer ID:', customerId);
      showDeleteConfirmation(customerId);
    } else {
      console.log('❓ Unknown action:', action);
    }
  }

	function handleSearch(event) {
		searchTerm = event.detail.searchTerm;
	}
</script>

<PageLayout 
	title="Customers" 
	icon="users" 
	showIndicator={true}
>
	<svelte:fragment slot="actions">
		<!-- Actions are handled by DataTable -->
	</svelte:fragment>

	<DataTable
		data={filteredCustomers}
		{columns}
		{loading}
		{searchTerm}
		searchPlaceholder="Search customers..."
		emptyStateTitle="No customers found"
		emptyStateMessage="Start by adding your first customer"
		emptyStateIcon="users"
		{primaryAction}
		{secondaryActions}
		on:primary-action={handlePrimaryAction}
		on:secondary-action={handleSecondaryAction}
		on:row-action={handleRowAction}
		on:search={handleSearch}
	>
		<svelte:fragment slot="cell" let:item let:column>
			{#if column.key === 'status'}
				<StatusBadge status={item.active ? 'active' : 'inactive'} />
			{:else if column.render}
				{@html column.render(item)}
			{:else}
				{item[column.key] || '-'}
			{/if}
		</svelte:fragment>
	</DataTable>
</PageLayout>

{#if showModal}
	<CustomerModal
		show={showModal}
		{editingCustomer}
		{loading}
		on:save={handleSaveCustomer}
		on:close={closeModal}
	/>
{/if}

{#if showDeleteConfirm}
	<ConfirmationModal
		show={showDeleteConfirm}
		title="Delete Customer"
		message="Are you sure you want to delete this customer? This action cannot be undone."
		confirmText="Delete"
		cancelText="Cancel"
		confirmClass="btn-error"
		{loading}
		on:confirm={confirmDelete}
		on:cancel={cancelDelete}
	/>
{/if}