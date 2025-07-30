import './style.css';
import './app.css';

// Import Wails API functions
import {
    CreateCustomer,
    GetCustomers,
    GetCustomerByID,
    UpdateCustomer,
    DeleteCustomer,
    CreateItemCategory,
    GetItemCategories,
    CreateSaleItem,
    GetSaleItems,
    CreateInvoice,
    GetInvoices,
    GetInvoiceByID,
    GetCompany,
    UpdateCompany,
    GenerateInvoicePDF,
    OpenPDFInViewer,
    ShowSaveDialog,
    ShowMessage,
    ShowError,
    ValidateZATCAQRCode,
    GetQRCodeInfo
} from '../wailsjs/go/main/App';

// Global state
let currentSection = 'dashboard';
let customers = [];
let saleItems = [];
let invoices = [];
let company = null;
let invoiceItemCounter = 0;
let wailsReady = false;

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    initializeNavigation();
    
    // Wait for Wails to be ready
    if (window.wails) {
        initializeApp();
    } else {
        // Fallback: wait a bit and try again
        setTimeout(() => {
            initializeApp();
        }, 1000);
    }
});

function initializeApp() {
    wailsReady = true;
    loadInitialData();
    setupEventListeners();
    
    // Set today's date as default for invoice date
    const today = new Date().toISOString().split('T')[0];
    const invoiceDateInput = document.getElementById('invoice-date');
    if (invoiceDateInput) {
        invoiceDateInput.value = today;
    }
}

// Navigation functionality
function initializeNavigation() {
    const navButtons = document.querySelectorAll('.nav-btn');
    navButtons.forEach(button => {
        button.addEventListener('click', () => {
            const section = button.getAttribute('data-section');
            showSection(section);
            
            // Update active nav button
            navButtons.forEach(btn => btn.classList.remove('active'));
            button.classList.add('active');
        });
    });
}

function showSection(sectionName) {
    // Hide all sections
    const sections = document.querySelectorAll('.content-section');
    sections.forEach(section => section.classList.remove('active'));
    
    // Show selected section
    const targetSection = document.getElementById(sectionName);
    if (targetSection) {
        targetSection.classList.add('active');
        currentSection = sectionName;
        
        // Load section-specific data
        if (wailsReady) {
            switch(sectionName) {
                case 'dashboard':
                    loadDashboardData();
                    break;
                case 'customers':
                    loadCustomersData();
                    break;
                case 'items':
                    loadItemsData();
                    break;
                case 'invoices':
                    loadInvoicesData();
                    break;
                case 'company':
                    loadCompanyData();
                    break;
            }
        }
    }
}

// Load initial data
async function loadInitialData() {
    if (!wailsReady) return;
    
    showLoading(true);
    try {
        await Promise.all([
            loadCustomers(),
            loadSaleItems(),
            loadInvoices(),
            loadCompany()
        ]);
        loadDashboardData();
    } catch (error) {
        console.error('Error loading initial data:', error);
        showErrorMessage('Error', 'Failed to load application data');
    } finally {
        showLoading(false);
    }
}

// Data loading functions
async function loadCustomers() {
    if (!wailsReady) return;
    try {
        customers = await GetCustomers();
    } catch (error) {
        console.error('Error loading customers:', error);
        customers = [];
    }
}

async function loadSaleItems() {
    if (!wailsReady) return;
    try {
        saleItems = await GetSaleItems();
    } catch (error) {
        console.error('Error loading sale items:', error);
        saleItems = [];
    }
}

async function loadInvoices() {
    if (!wailsReady) return;
    try {
        invoices = await GetInvoices();
    } catch (error) {
        console.error('Error loading invoices:', error);
        invoices = [];
    }
}

async function loadCompany() {
    if (!wailsReady) return;
    try {
        company = await GetCompany();
    } catch (error) {
        console.error('Error loading company:', error);
        company = null;
    }
}

// Dashboard functionality
function loadDashboardData() {
    const totalCustomers = (customers || []).length;
    const totalInvoices = (invoices || []).length;
    const totalItems = (saleItems || []).length;
    const totalRevenue = (invoices || []).reduce((sum, invoice) => sum + (invoice.totalAmount || 0), 0);
    
    document.getElementById('total-customers').textContent = totalCustomers;
    document.getElementById('total-invoices').textContent = totalInvoices;
    document.getElementById('total-items').textContent = totalItems;
    document.getElementById('total-revenue').textContent = `$${totalRevenue.toFixed(2)}`;
}

// Customer management
function loadCustomersData() {
    const tbody = document.getElementById('customers-table');
    tbody.innerHTML = '';
    
    (customers || []).forEach(customer => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${customer.nameEn || customer.nameAr || 'N/A'}</td>
            <td>${customer.email || 'N/A'}</td>
            <td>${customer.phone || 'N/A'}</td>
            <td>${customer.vatNumber || 'N/A'}</td>
            <td class="action-buttons">
                <button class="action-btn btn-secondary" onclick="editCustomer(${customer.id})">
                    <i class="fas fa-edit"></i>
                </button>
                <button class="action-btn btn-danger" onclick="deleteCustomer(${customer.id})">
                    <i class="fas fa-trash"></i>
                </button>
            </td>
        `;
        tbody.appendChild(row);
    });
}

// Sale items management
function loadItemsData() {
    const tbody = document.getElementById('items-table');
    tbody.innerHTML = '';
    
    (saleItems || []).forEach(item => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${item.nameEn || item.nameAr || 'N/A'}</td>
            <td>General</td>
            <td>$${item.price ? item.price.toFixed(2) : '0.00'}</td>
            <td>${item.vatRate || 0}%</td>
            <td class="action-buttons">
                <button class="action-btn btn-secondary" onclick="editSaleItem(${item.id})">
                    <i class="fas fa-edit"></i>
                </button>
                <button class="action-btn btn-danger" onclick="deleteSaleItem(${item.id})">
                    <i class="fas fa-trash"></i>
                </button>
            </td>
        `;
        tbody.appendChild(row);
    });
}

// Invoice management
function loadInvoicesData() {
    const tbody = document.getElementById('invoices-table');
    tbody.innerHTML = '';
    
    (invoices || []).forEach(invoice => {
        const customer = (customers || []).find(c => c.id === invoice.customerID);
        const customerName = customer ? (customer.nameEn || customer.nameAr) : 'Unknown';
        const invoiceDate = new Date(invoice.invoiceDate).toLocaleDateString();
        
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${invoice.invoiceNumber || 'N/A'}</td>
            <td>${customerName}</td>
            <td>${invoiceDate}</td>
            <td>$${invoice.totalAmount ? invoice.totalAmount.toFixed(2) : '0.00'}</td>
            <td><span class="status-badge status-${invoice.status || 'draft'}">${invoice.status || 'draft'}</span></td>
            <td class="action-buttons">
                <button class="action-btn btn-success" onclick="generatePDF(${invoice.id})" title="Generate PDF">
                    <i class="fas fa-file-pdf"></i>
                </button>
                <button class="action-btn btn-secondary" onclick="viewInvoice(${invoice.id})" title="View">
                    <i class="fas fa-eye"></i>
                </button>
                <button class="action-btn btn-danger" onclick="deleteInvoice(${invoice.id})" title="Delete">
                    <i class="fas fa-trash"></i>
                </button>
            </td>
        `;
        tbody.appendChild(row);
    });
}

// Company management
function loadCompanyData() {
    if (company) {
        document.getElementById('company-name-en').value = company.nameEn || '';
        document.getElementById('company-name-ar').value = company.nameAr || '';
        document.getElementById('company-vat').value = company.vatNumber || '';
        document.getElementById('company-cr').value = company.commercialRegistration || '';
        document.getElementById('company-address-en').value = company.addressEn || '';
        document.getElementById('company-address-ar').value = company.addressAr || '';
        document.getElementById('company-phone').value = company.phone || '';
        document.getElementById('company-email').value = company.email || '';
    }
}

// Modal functionality
function showModal(modalId) {
    const overlay = document.getElementById('modal-overlay');
    const modal = document.getElementById(modalId);
    
    // Hide all modals first
    const allModals = overlay.querySelectorAll('.modal');
    allModals.forEach(m => m.style.display = 'none');
    
    // Show selected modal
    modal.style.display = 'block';
    overlay.classList.add('active');
}

function closeModal() {
    const overlay = document.getElementById('modal-overlay');
    overlay.classList.remove('active');
    
    // Reset forms
    const forms = overlay.querySelectorAll('form');
    forms.forEach(form => form.reset());
    
    // Clear invoice items
    const container = document.getElementById('invoice-items-container');
    if (container) {
        container.innerHTML = '';
        invoiceItemCounter = 0;
        updateInvoiceTotals();
    }
}

// Customer modal functions
function showCustomerModal(customerId = null) {
    showModal('customer-modal');
    
    if (customerId) {
        // Edit mode - populate form with customer data
        const customer = customers.find(c => c.id === customerId);
        if (customer) {
            document.getElementById('customer-name-en').value = customer.nameEn || '';
            document.getElementById('customer-name-ar').value = customer.nameAr || '';
            document.getElementById('customer-email').value = customer.email || '';
            document.getElementById('customer-phone').value = customer.phone || '';
            document.getElementById('customer-vat').value = customer.vatNumber || '';
            document.getElementById('customer-address-en').value = customer.addressEn || '';
            document.getElementById('customer-address-ar').value = customer.addressAr || '';
            
            // Store customer ID for update
            document.getElementById('customer-form').setAttribute('data-customer-id', customerId);
        }
    }
}

// Sale item modal functions
function showItemModal(itemId = null) {
    showModal('item-modal');
    
    if (itemId) {
        // Edit mode - populate form with item data
        const item = saleItems.find(i => i.id === itemId);
        if (item) {
            document.getElementById('item-name-en').value = item.nameEn || '';
            document.getElementById('item-name-ar').value = item.nameAr || '';
            document.getElementById('item-price').value = item.price || '';
            document.getElementById('item-vat-rate').value = item.vatRate || '';
            document.getElementById('item-description-en').value = item.descriptionEn || '';
            document.getElementById('item-description-ar').value = item.descriptionAr || '';
            
            // Store item ID for update
            document.getElementById('item-form').setAttribute('data-item-id', itemId);
        }
    }
}

// Invoice modal functions
function showInvoiceModal() {
    showModal('invoice-modal');
    populateCustomerDropdown();
    addInvoiceItem(); // Add first item by default
}

function populateCustomerDropdown() {
    const select = document.getElementById('invoice-customer');
    select.innerHTML = '<option value="">Select Customer</option>';
    
    (customers || []).forEach(customer => {
        const option = document.createElement('option');
        option.value = customer.id;
        option.textContent = customer.nameEn || customer.nameAr || `Customer ${customer.id}`;
        select.appendChild(option);
    });
}

function addInvoiceItem() {
    const container = document.getElementById('invoice-items-container');
    const itemDiv = document.createElement('div');
    itemDiv.className = 'invoice-item';
    itemDiv.setAttribute('data-item-id', invoiceItemCounter);
    
    itemDiv.innerHTML = `
        <div class="form-group">
            <label>Sale Item</label>
            <select name="saleItemID" onchange="updateItemPrice(this)" required>
                <option value="">Select Item</option>
                ${(saleItems || []).map(item => 
                    `<option value="${item.id}" data-price="${item.price}" data-vat="${item.vatRate}">
                        ${item.nameEn || item.nameAr || `Item ${item.id}`}
                    </option>`
                ).join('')}
            </select>
        </div>
        <div class="form-group">
            <label>Quantity</label>
            <input type="number" name="quantity" min="1" value="1" onchange="updateInvoiceTotals()" required>
        </div>
        <div class="form-group">
            <label>Unit Price</label>
            <input type="number" name="unitPrice" step="0.01" min="0" onchange="updateInvoiceTotals()" required>
        </div>
        <div class="form-group">
            <label>VAT Rate (%)</label>
            <input type="number" name="vatRate" step="0.01" min="0" max="100" onchange="updateInvoiceTotals()" required>
        </div>
        <button type="button" class="remove-item" onclick="removeInvoiceItem(this)">
            <i class="fas fa-times"></i>
        </button>
    `;
    
    container.appendChild(itemDiv);
    invoiceItemCounter++;
}

function removeInvoiceItem(button) {
    button.parentElement.remove();
    updateInvoiceTotals();
}

function updateItemPrice(select) {
    const option = select.selectedOptions[0];
    if (option && option.dataset.price) {
        const itemDiv = select.closest('.invoice-item');
        const priceInput = itemDiv.querySelector('input[name="unitPrice"]');
        const vatInput = itemDiv.querySelector('input[name="vatRate"]');
        
        priceInput.value = option.dataset.price;
        vatInput.value = option.dataset.vat || 15;
        
        updateInvoiceTotals();
    }
}

function updateInvoiceTotals() {
    const items = document.querySelectorAll('.invoice-item');
    let subtotal = 0;
    let vatAmount = 0;
    
    items.forEach(item => {
        const quantity = parseFloat(item.querySelector('input[name="quantity"]').value) || 0;
        const unitPrice = parseFloat(item.querySelector('input[name="unitPrice"]').value) || 0;
        const vatRate = parseFloat(item.querySelector('input[name="vatRate"]').value) || 0;
        
        const lineTotal = quantity * unitPrice;
        const lineVat = (lineTotal * vatRate) / 100;
        
        subtotal += lineTotal;
        vatAmount += lineVat;
    });
    
    const total = subtotal + vatAmount;
    
    document.getElementById('invoice-subtotal').textContent = `$${subtotal.toFixed(2)}`;
    document.getElementById('invoice-vat').textContent = `$${vatAmount.toFixed(2)}`;
    document.getElementById('invoice-total').textContent = `$${total.toFixed(2)}`;
}

// Event listeners setup
function setupEventListeners() {
    // Customer form submission
    document.getElementById('customer-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        await handleCustomerSubmit(e.target);
    });
    
    // Item form submission
    document.getElementById('item-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        await handleItemSubmit(e.target);
    });
    
    // Invoice form submission
    document.getElementById('invoice-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        await handleInvoiceSubmit(e.target);
    });
    
    // Company form submission
    document.getElementById('company-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        await handleCompanySubmit(e.target);
    });
    
    // Modal overlay click to close
    document.getElementById('modal-overlay').addEventListener('click', (e) => {
        if (e.target.id === 'modal-overlay') {
            closeModal();
        }
    });
}

// Form submission handlers
async function handleCustomerSubmit(form) {
    if (!wailsReady) return;
    showLoading(true);
    try {
        const formData = new FormData(form);
        const customer = {
            nameEn: formData.get('nameEn'),
            nameAr: formData.get('nameAr'),
            email: formData.get('email'),
            phone: formData.get('phone'),
            vatNumber: formData.get('vatNumber'),
            addressEn: formData.get('addressEn'),
            addressAr: formData.get('addressAr')
        };
        
        const customerId = form.getAttribute('data-customer-id');
        if (customerId) {
            // Update existing customer
            customer.id = parseInt(customerId);
            await UpdateCustomer(customer);
            showSuccessMessage('Success', 'Customer updated successfully');
        } else {
            // Create new customer
            await CreateCustomer(customer);
            showSuccessMessage('Success', 'Customer created successfully');
        }
        
        await loadCustomers();
        if (currentSection === 'customers') {
            loadCustomersData();
        }
        closeModal();
    } catch (error) {
        console.error('Error saving customer:', error);
        showErrorMessage('Error', 'Failed to save customer');
    } finally {
        showLoading(false);
    }
}

async function handleItemSubmit(form) {
    if (!wailsReady) return;
    showLoading(true);
    try {
        const formData = new FormData(form);
        const item = {
            nameEn: formData.get('nameEn'),
            nameAr: formData.get('nameAr'),
            price: parseFloat(formData.get('price')),
            vatRate: parseFloat(formData.get('vatRate')),
            descriptionEn: formData.get('descriptionEn'),
            descriptionAr: formData.get('descriptionAr'),
            categoryID: 1 // Default category
        };
        
        const itemId = form.getAttribute('data-item-id');
        if (itemId) {
            // Update existing item
            item.id = parseInt(itemId);
            // Note: UpdateSaleItem method needs to be implemented in Go
            showErrorMessage('Error', 'Update functionality not yet implemented');
        } else {
            // Create new item
            await CreateSaleItem(item);
            showSuccessMessage('Success', 'Sale item created successfully');
        }
        
        await loadSaleItems();
        if (currentSection === 'items') {
            loadItemsData();
        }
        closeModal();
    } catch (error) {
        console.error('Error saving sale item:', error);
        showErrorMessage('Error', 'Failed to save sale item');
    } finally {
        showLoading(false);
    }
}

async function handleInvoiceSubmit(form) {
    if (!wailsReady) return;
    showLoading(true);
    try {
        const formData = new FormData(form);
        const items = [];
        
        // Collect invoice items
        const itemDivs = document.querySelectorAll('.invoice-item');
        itemDivs.forEach(div => {
            const saleItemID = parseInt(div.querySelector('select[name="saleItemID"]').value);
            const quantity = parseFloat(div.querySelector('input[name="quantity"]').value);
            const unitPrice = parseFloat(div.querySelector('input[name="unitPrice"]').value);
            const vatRate = parseFloat(div.querySelector('input[name="vatRate"]').value);
            
            if (saleItemID && quantity && unitPrice >= 0) {
                items.push({
                    saleItemID,
                    quantity,
                    unitPrice,
                    vatRate
                });
            }
        });
        
        if (items.length === 0) {
            showErrorMessage('Error', 'Please add at least one item to the invoice');
            return;
        }
        
        const invoice = {
            customerID: parseInt(formData.get('customerID')),
            invoiceDate: formData.get('invoiceDate'),
            items: items
        };
        
        await CreateInvoice(invoice);
        showSuccessMessage('Success', 'Invoice created successfully');
        
        await loadInvoices();
        if (currentSection === 'invoices') {
            loadInvoicesData();
        }
        closeModal();
    } catch (error) {
        console.error('Error creating invoice:', error);
        showErrorMessage('Error', 'Failed to create invoice');
    } finally {
        showLoading(false);
    }
}

async function handleCompanySubmit(form) {
    if (!wailsReady) return;
    showLoading(true);
    try {
        const formData = new FormData(form);
        const companyData = {
            nameEn: formData.get('nameEn'),
            nameAr: formData.get('nameAr'),
            vatNumber: formData.get('vatNumber'),
            commercialRegistration: formData.get('commercialRegistration'),
            addressEn: formData.get('addressEn'),
            addressAr: formData.get('addressAr'),
            phone: formData.get('phone'),
            email: formData.get('email')
        };
        
        await UpdateCompany(companyData);
        showSuccessMessage('Success', 'Company information updated successfully');
        
        await loadCompany();
    } catch (error) {
        console.error('Error updating company:', error);
        showErrorMessage('Error', 'Failed to update company information');
    } finally {
        showLoading(false);
    }
}

// Action functions
async function generatePDF(invoiceId) {
    if (!wailsReady) return;
    showLoading(true);
    try {
        const filePath = await GenerateInvoicePDF(invoiceId);
        showSuccessMessage('Success', `PDF generated successfully: ${filePath}`);
        
        // Optionally open the PDF
        if (confirm('Would you like to open the PDF?')) {
            await OpenPDFInViewer(filePath);
        }
    } catch (error) {
        console.error('Error generating PDF:', error);
        showErrorMessage('Error', 'Failed to generate PDF');
    } finally {
        showLoading(false);
    }
}

function editCustomer(customerId) {
    showCustomerModal(customerId);
}

function editSaleItem(itemId) {
    showItemModal(itemId);
}

async function deleteCustomer(customerId) {
    if (!wailsReady) return;
    if (confirm('Are you sure you want to delete this customer?')) {
        showLoading(true);
        try {
            await DeleteCustomer(customerId);
            showSuccessMessage('Success', 'Customer deleted successfully');
            await loadCustomers();
            if (currentSection === 'customers') {
                loadCustomersData();
            }
        } catch (error) {
            console.error('Error deleting customer:', error);
            showErrorMessage('Error', 'Failed to delete customer');
        } finally {
            showLoading(false);
        }
    }
}

function deleteSaleItem(itemId) {
    showErrorMessage('Error', 'Delete functionality not yet implemented');
}

function deleteInvoice(invoiceId) {
    showErrorMessage('Error', 'Delete functionality not yet implemented');
}

function viewInvoice(invoiceId) {
    showErrorMessage('Info', 'View functionality not yet implemented');
}

// Utility functions
function showLoading(show) {
    const loading = document.getElementById('loading');
    if (loading) {
        if (show) {
            loading.classList.add('active');
        } else {
            loading.classList.remove('active');
        }
    }
}

function showSuccessMessage(title, message) {
    if (wailsReady) {
        ShowMessage(title, message).catch(console.error);
    } else {
        alert(`${title}: ${message}`);
    }
}

function showErrorMessage(title, message) {
    if (wailsReady) {
        ShowError(title, message).catch(console.error);
    } else {
        alert(`${title}: ${message}`);
    }
}

// QR Code Validation Functions
async function validateQRCode() {
    if (!wailsReady) return;
    
    const qrInput = document.getElementById('qr-code-input');
    const qrCode = qrInput.value.trim();
    
    if (!qrCode) {
        showErrorMessage('Error', 'Please enter a QR code to validate');
        return;
    }
    
    // Hide previous results
    document.getElementById('qr-validation-results').style.display = 'none';
    document.getElementById('qr-validation-error').style.display = 'none';
    
    showLoading(true);
    
    try {
        // Validate the QR code
        const isValid = await ValidateZATCAQRCode(qrCode);
        
        if (isValid) {
            // Get QR code information
            const qrInfo = await GetQRCodeInfo(qrCode);
            displayQRValidationResults(qrInfo, qrCode);
        } else {
            displayQRValidationError('Invalid ZATCA QR code format or structure');
        }
    } catch (error) {
        console.error('Error validating QR code:', error);
        displayQRValidationError(error.message || 'Failed to validate QR code');
    } finally {
        showLoading(false);
    }
}

function displayQRValidationResults(qrInfo, qrCode) {
    // Update QR code information
    document.getElementById('qr-seller-name').textContent = qrInfo.sellerName || 'N/A';
    document.getElementById('qr-vat-number').textContent = qrInfo.vatNumber || 'N/A';
    document.getElementById('qr-timestamp').textContent = qrInfo.timestamp || 'N/A';
    document.getElementById('qr-total-amount').textContent = qrInfo.totalAmount || 'N/A';
    document.getElementById('qr-vat-amount').textContent = qrInfo.vatAmount || 'N/A';
    document.getElementById('qr-hash-present').textContent = qrInfo.hashPresent ? 'Yes' : 'No';
    
    // Update technical details
    document.getElementById('qr-length').textContent = `${qrCode.length} characters`;
    
    // Update status badge
    const statusBadge = document.getElementById('validation-status');
    statusBadge.textContent = 'VALID';
    statusBadge.className = 'status-badge';
    
    // Show results
    document.getElementById('qr-validation-results').style.display = 'block';
    document.getElementById('qr-validation-error').style.display = 'none';
}

function displayQRValidationError(errorMessage) {
    document.getElementById('error-message').textContent = errorMessage;
    document.getElementById('qr-validation-error').style.display = 'block';
    document.getElementById('qr-validation-results').style.display = 'none';
}

function clearQRValidation() {
    document.getElementById('qr-code-input').value = '';
    document.getElementById('qr-validation-results').style.display = 'none';
    document.getElementById('qr-validation-error').style.display = 'none';
}

// Make functions globally available
window.showCustomerModal = showCustomerModal;
window.showItemModal = showItemModal;
window.showInvoiceModal = showInvoiceModal;
window.closeModal = closeModal;
window.addInvoiceItem = addInvoiceItem;
window.removeInvoiceItem = removeInvoiceItem;
window.updateItemPrice = updateItemPrice;
window.updateInvoiceTotals = updateInvoiceTotals;
window.generatePDF = generatePDF;
window.editCustomer = editCustomer;
window.editSaleItem = editSaleItem;
window.deleteCustomer = deleteCustomer;
window.deleteSaleItem = deleteSaleItem;
window.deleteInvoice = deleteInvoice;
window.viewInvoice = viewInvoice;
window.validateQRCode = validateQRCode;
window.clearQRValidation = clearQRValidation;
