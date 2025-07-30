<script>
	import { onMount } from 'svelte';
	import { GetCustomers, CreateCustomer, UpdateCustomer, DeleteCustomer } from '../wailsjs/go/main/App.js';
	import { database } from '../wailsjs/go/models';
	import PageLayout from './components/PageLayout.svelte';
	import DataTable from './components/DataTable.svelte';
	import StatusBadge from './components/StatusBadge.svelte';
	import CustomerModal from './components/CustomerModal.svelte';

	let customers = [];
	let searchTerm = '';
	let showModal = false;
	let editingCustomer = null;
	let loading = false;

	onMount(async () => {
		await loadCustomers();
	});

	async function loadCustomers() {
		loading = true;
		try {
			customers = await GetCustomers();
		} catch (error) {
			console.error('Error loading customers:', error);
			customers = [];
		} finally {
			loading = false;
		}
	}

	function openModal(customer = null) {
		editingCustomer = customer;
		showModal = true;
	}

	function closeModal() {
		showModal = false;
		editingCustomer = null;
	}

	async function handleSaveCustomer(event) {
		const customerData = event.detail;
		loading = true;
		
		try {
			const customer = new database.Customer();
			Object.assign(customer, customerData);

			if (editingCustomer) {
				customer.id = editingCustomer.id;
				await UpdateCustomer(customer);
			} else {
				await CreateCustomer(customer);
			}

			await loadCustomers();
			closeModal();
		} catch (error) {
			console.error('Error saving customer:', error);
			alert('Error saving customer: ' + error.message);
		} finally {
			loading = false;
		}
	}

	async function deleteCustomer(id) {
		if (confirm('Are you sure you want to delete this customer?')) {
			loading = true;
			try {
				await DeleteCustomer(id);
				await loadCustomers();
			} catch (error) {
				console.error('Error deleting customer:', error);
				alert('Error deleting customer: ' + error.message);
			} finally {
				loading = false;
			}
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
		{ 
			key: 'status', 
			label: 'Status',
			actions: [
				{ key: 'edit', icon: 'fa-edit', class: 'btn-secondary', title: 'Edit Customer' },
				{ key: 'delete', icon: 'fa-trash', class: 'btn-error', title: 'Delete Customer' }
			]
		}
	];

	const primaryAction = {
		text: 'Add Customer',
		icon: 'fa-plus'
	};

	const secondaryActions = [
		{
			text: 'Import',
			icon: 'fa-upload'
		},
		{
			text: 'Export',
			icon: 'fa-download'
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
		const { action, item } = event.detail;
		if (action === 'edit') {
			openModal(item);
		} else if (action === 'delete') {
			deleteCustomer(item.id);
		}
	}

	function handleSearch(event) {
		searchTerm = event.detail.searchTerm;
	}
</script>

<PageLayout 
	title="Customers" 
	icon="fa-users" 
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
		emptyStateIcon="fa-users"
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
		{editingCustomer}
		{loading}
		on:save={handleSaveCustomer}
		on:close={closeModal}
	/>
{/if}

<style>
  .customer-name .name-primary {
    font-weight: 500;
    color: var(--text-primary);
  }

  .customer-name .name-secondary {
    font-size: 0.875rem;
    color: var(--text-secondary);
    margin-top: 2px;
  }

  .city-info, .country-info {
    line-height: 1.4;
  }

  .city-arabic, .country-arabic {
    font-size: 0.875rem;
    color: var(--text-secondary);
    margin-top: 2px;
  }

  .action-buttons {
    display: flex;
    gap: 0.5rem;
  }

  .loading-state {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 3rem;
    color: var(--text-secondary);
    gap: 0.75rem;
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 3rem;
    text-align: center;
  }

  .empty-state i {
    font-size: 3rem;
    color: var(--text-muted);
    margin-bottom: 1rem;
  }

  .empty-state h3 {
    margin: 0 0 0.5rem 0;
    color: var(--text-primary);
  }

  .empty-state p {
    margin: 0 0 1.5rem 0;
    color: var(--text-secondary);
  }
</style>