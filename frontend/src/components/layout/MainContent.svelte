<script>
  import Dashboard from '../../Dashboard.svelte';
  import POS from '../../POS.svelte';
  import QRValidation from '../../QRValidation.svelte';
  import SalesInvoices from '../../SalesInvoices.svelte';
  import Customer from '../../Customer.svelte';
  import Payments from '../../Payments.svelte';
  import PurchaseInvoices from '../../PurchaseInvoices.svelte';
  import PurchaseProducts from '../../PurchaseProducts.svelte';
  import Suppliers from '../../Suppliers.svelte';
  import Products from '../../Products.svelte';
  import Users from '../../Users.svelte';
  import GeneralSettings from '../../GeneralSettings.svelte';
  import Administration from '../../Administration.svelte';
  
  export let currentView = 'dashboard';
  export let showIntro = () => {};
</script>

<main class="flex-1 overflow-y-auto">
  {#if currentView === 'dashboard'}
    <Dashboard />
  {:else if currentView === 'pos'}
    <POS />
  {:else if currentView === 'qr-validation'}
    <div class="p-6">
      <QRValidation />
    </div>
  {:else if currentView === 'sales-invoices'}
    <SalesInvoices />
  {:else if currentView === 'customer'}
    <Customer />
  {:else if currentView === 'payments'}
    <Payments />
  {:else if currentView === 'purchase-invoices'}
    <PurchaseInvoices />
  {:else if currentView === 'purchase-products'}
    <PurchaseProducts />
  {:else if currentView === 'suppliers'}
    <Suppliers />
  {:else if currentView === 'products'}
    <Products />
  {:else if currentView === 'users'}
    <Users />
  {:else if currentView === 'general-settings'}
    <GeneralSettings />
  {:else if currentView === 'administration'}
    <Administration {showIntro} />
  {/if}
</main>