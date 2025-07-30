
# dijibill: Complete Application Documentation

A cross-platform desktop invoicing and business management app for SMEs, built with **Go**, **Svelte**, and **Wails**, tailored for the **Saudi Arabian market** with bilingual (English/Arabic) support and local invoice standards (VAT, QR).

---

## Table of Contents

1. [Tech Stack](#1-tech-stack)
2. [Application Architecture](#2-application-architecture)
3. [Data Models](#3-data-models)
4. [Design System](#4-design-system)
5. [Component Guidelines](#5-component-guidelines)
6. [Backend API](#6-backend-api)
7. [Frontend Structure](#7-frontend-structure)
8. [Development Workflow](#8-development-workflow)
9. [Deployment](#9-deployment)
10. [Troubleshooting](#10-troubleshooting)

---

## 1. Tech Stack

- **Backend:** Go 1.21+
- **Frontend:** Svelte 4 + Vite
- **Desktop Framework:** Wails v2
- **Database:** SQLite with Go database/sql
- **Styling:** Tailwind CSS + DaisyUI + Custom Design System
- **Build Tool:** Vite (frontend), Go build (backend)
- **Package Manager:** npm (frontend), Go modules (backend)

---

## 2. Application Architecture

### 2.1 Project Structure

```
dijibill/
├── app.go                 # Wails app context and API methods
├── database.go           # Database operations and CRUD
├── main.go              # Application entry point
├── models.go            # Data structures and types
├── pdf_service.go       # PDF generation service
├── qr_service.go        # QR code generation service
├── context/             # Documentation
│   ├── APP_CONTEXT.md   # This file
│   └── ZATCA_QR_COMPLIANCE.md
└── frontend/
    ├── src/
    │   ├── App.svelte           # Main application component
    │   ├── Dashboard.svelte     # Dashboard view
    │   ├── GeneralSettings.svelte # Settings management
    │   ├── Products.svelte      # Product management
    │   ├── Customer.svelte      # Customer management
    │   ├── SalesInvoices.svelte # Invoice management
    │   ├── components/          # Reusable components
    │   ├── styles/             # Design system and styles
    │   └── wailsjs/            # Generated Wails bindings
    ├── package.json
    ├── tailwind.config.js
    └── vite.config.js
```

### 2.2 Application Flow

1. **Initialization**: Wails starts Go backend and Svelte frontend
2. **Database**: SQLite database is initialized with required tables
3. **Navigation**: Single-page application with component-based routing
4. **Data Flow**: Frontend calls Go methods via Wails bindings
5. **State Management**: Svelte stores and reactive variables
6. **UI Updates**: Real-time updates through Svelte reactivity

---

## 3. Data Models

### 3.1 Core Models (see `models.go`)

- **Company**: Business information and settings
- **Customer**: Customer details and contact information
- **Product**: Product/service catalog with pricing
- **Invoice**: Sales invoices with line items
- **PaymentType**: Payment method definitions
- **SalesCategory**: Product categorization
- **TaxRate**: VAT and tax configurations
- **UnitOfMeasurement**: Measurement units for products
- **DefaultProductSettings**: Default values for new products

### 3.2 Model Relationships

```
Company (1) ←→ (∞) Invoice
Customer (1) ←→ (∞) Invoice
Invoice (1) ←→ (∞) InvoiceItem
Product (1) ←→ (∞) InvoiceItem
TaxRate (1) ←→ (∞) Product
PaymentType (1) ←→ (∞) Invoice
SalesCategory (1) ←→ (∞) Product
UnitOfMeasurement (1) ←→ (∞) Product
```

---

## 4. Design System

### 4.1 Overview

The application uses a centralized design system located in `/frontend/src/styles/design-system.css` that provides:

- **Consistent color palette**
- **Standardized component styles**
- **Responsive design patterns**
- **Accessibility features**
- **Glass morphism effects**

### 4.2 Color System

#### Primary Colors
```css
--color-primary: #667eea        /* Main brand color */
--color-primary-light: #8b9df0  /* Lighter variant */
--color-primary-dark: #4c63d2   /* Darker variant */
```

#### Status Colors
```css
--color-success: #36d399        /* Success states */
--color-warning: #fbbd23        /* Warning states */
--color-error: #f87272          /* Error states */
--color-info: #3abff8           /* Info states */
```

#### Glass/Backdrop Colors
```css
--color-glass-white: rgba(255, 255, 255, 0.1)
--color-glass-white-strong: rgba(255, 255, 255, 0.2)
```

### 4.3 Component Classes

#### Glass Cards
```html
<!-- Basic glass card -->
<div class="glass-card">
  <h3 class="heading-3">Card Title</h3>
  <p class="text-secondary">Card content</p>
</div>

<!-- Strong glass effect -->
<div class="glass-card-strong">
  Content with stronger glass effect
</div>

<!-- Subtle glass effect -->
<div class="glass-card-subtle">
  Content with subtle glass effect
</div>
```

#### Buttons
```html
<!-- Primary buttons -->
<button class="btn-primary">Primary Action</button>
<button class="btn-primary-outline">Primary Outline</button>
<button class="btn-primary-ghost">Primary Ghost</button>

<!-- Glass buttons -->
<button class="btn-glass">Glass Button</button>
<button class="btn-glass-outline">Glass Outline</button>

<!-- Status buttons -->
<button class="btn-success">Success</button>
<button class="btn-warning">Warning</button>
<button class="btn-error">Error</button>
<button class="btn-info">Info</button>

<!-- Button sizes -->
<button class="btn-primary btn-sm">Small</button>
<button class="btn-primary">Default</button>
<button class="btn-primary btn-lg">Large</button>
<button class="btn-primary btn-xl">Extra Large</button>

<!-- Icon buttons -->
<button class="btn-icon btn-primary">
  <svg class="w-5 h-5">...</svg>
</button>
```

#### Form Inputs
```html
<!-- Glass inputs (for dark backgrounds) -->
<div class="form-control">
  <label class="label-glass">Field Label</label>
  <input type="text" class="input-glass" placeholder="Enter text">
</div>

<!-- Standard inputs (for light backgrounds) -->
<div class="form-control">
  <label class="label-standard">Field Label</label>
  <input type="text" class="input-standard" placeholder="Enter text">
</div>

<!-- Select inputs -->
<select class="select-glass">
  <option>Option 1</option>
  <option>Option 2</option>
</select>

<!-- Textarea -->
<textarea class="textarea-glass" placeholder="Enter description"></textarea>
```

#### Typography
```html
<!-- Headings -->
<h1 class="heading-1">Main Title</h1>
<h2 class="heading-2">Section Title</h2>
<h3 class="heading-3">Subsection Title</h3>
<h4 class="heading-4">Component Title</h4>

<!-- Text variants -->
<p class="text-primary">Primary text (white)</p>
<p class="text-secondary">Secondary text (white/80)</p>
<p class="text-muted">Muted text (white/60)</p>
<p class="text-disabled">Disabled text (white/40)</p>

<!-- Dark theme text -->
<p class="text-dark-primary">Dark primary text</p>
<p class="text-dark-secondary">Dark secondary text</p>
<p class="text-dark-muted">Dark muted text</p>
```

#### Status Badges
```html
<span class="badge-success">Active</span>
<span class="badge-warning">Pending</span>
<span class="badge-error">Inactive</span>
<span class="badge-info">Info</span>
<span class="badge-primary">Featured</span>
<span class="badge-neutral">Default</span>
```

#### Tables
```html
<!-- Glass table (for dark backgrounds) -->
<table class="table-glass">
  <thead>
    <tr>
      <th>Column 1</th>
      <th>Column 2</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>Data 1</td>
      <td>Data 2</td>
    </tr>
  </tbody>
</table>

<!-- Standard table (for light backgrounds) -->
<table class="table-standard">
  <!-- Same structure -->
</table>
```

#### Modals
```html
<div class="modal-glass">
  <div class="modal-backdrop"></div>
  <div class="modal-content">
    <div class="modal-header">
      <h3 class="heading-3">Modal Title</h3>
    </div>
    <div class="modal-body">
      Modal content goes here
    </div>
    <div class="modal-footer">
      <button class="btn-glass">Cancel</button>
      <button class="btn-primary">Save</button>
    </div>
  </div>
</div>
```

#### Navigation
```html
<!-- Glass navigation -->
<nav class="nav-glass">
  <div class="nav-section-title">Main</div>
  <button class="nav-item nav-item-active">
    <span>Dashboard</span>
  </button>
  <button class="nav-item">
    <span>Products</span>
  </button>
</nav>
```

#### Tabs
```html
<!-- Glass tabs -->
<div class="tabs-glass">
  <button class="tab-glass tab-glass-active">Tab 1</button>
  <button class="tab-glass">Tab 2</button>
  <button class="tab-glass">Tab 3</button>
</div>

<!-- Standard tabs -->
<div class="tabs-standard">
  <button class="tab-standard tab-standard-active">Tab 1</button>
  <button class="tab-standard">Tab 2</button>
</div>
```

### 4.4 Utility Classes

```html
<!-- Background gradients -->
<div class="bg-gradient-primary">Primary gradient background</div>
<div class="bg-gradient-secondary">Secondary gradient background</div>
<div class="bg-gradient-dark">Dark gradient background</div>

<!-- Text gradient -->
<h1 class="text-gradient">Gradient text effect</h1>

<!-- Glass effects -->
<div class="backdrop-blur-glass">Enhanced blur effect</div>
<div class="shadow-glass">Glass shadow effect</div>

<!-- Loading states -->
<div class="loading-spinner"></div>
<div class="loading-spinner-sm"></div>
<div class="loading-spinner-lg"></div>

<!-- Loading overlay -->
<div class="loading-overlay">
  <div class="loading-spinner-lg"></div>
</div>
```

### 4.5 Responsive Design

The design system includes responsive utilities:

```html
<!-- Responsive buttons -->
<button class="btn-primary btn-responsive">Full width on mobile</button>

<!-- Responsive cards -->
<div class="glass-card">
  <!-- Automatically adjusts margins on different screen sizes -->
</div>
```

### 4.6 Accessibility Features

- **Focus states**: All interactive elements have visible focus indicators
- **High contrast support**: Enhanced borders and colors in high contrast mode
- **Reduced motion**: Respects user's motion preferences
- **Screen reader support**: Semantic HTML and ARIA labels

---

## 5. Component Guidelines

### 5.1 RULE #1: Use Design System First

**Always use the design system classes before writing custom CSS.**

#### ✅ Benefits
- Consistent UI across the application
- Faster development with pre-built components
- Automatic accessibility features
- Responsive design out of the box
- Easy maintenance and updates

#### ✅ Guidelines
- Check the design system documentation first
- Use semantic classes like `btn-primary`, `input-glass`, `glass-card`
- Stick to the defined color palette
- Combine with Tailwind for layout and spacing

#### ✅ Example
```html
<div class="glass-card">
  <div class="p-6">
    <h2 class="heading-2 mb-4">Add Product</h2>
    <div class="space-y-4">
      <div class="form-control">
        <label class="label-glass">Product Name</label>
        <input type="text" class="input-glass" />
      </div>
      <div class="flex justify-end space-x-3">
        <button class="btn-glass">Cancel</button>
        <button class="btn-primary">Save</button>
      </div>
    </div>
  </div>
</div>
```

### 5.2 RULE #2: Component Composition

**Break complex UIs into smaller, reusable Svelte components.**

#### ✅ Structure Example
```html
<InvoiceEditor>
  <CustomerSelector />
  <InvoiceItemTable />
  <InvoiceTotals />
  <InvoiceActions />
</InvoiceEditor>
```

#### ✅ Component Responsibilities
- Each component should have a single responsibility
- Keep components under 200 lines when possible
- Use props for data input and events for data output
- Avoid deeply nested component hierarchies

#### ✅ HorizontalTabs Component

**Universal horizontal tab navigation component for consistent tab interfaces.**

##### Props
- `tabs` (Array): Array of tab objects with `id`, `label`, and `icon` properties
- `activeTab` (String): Currently active tab ID
- `variant` (String): 'standard' or 'glass' styling variant
- `showScrollButtons` (Boolean): Enable/disable scroll buttons for overflow
- `scrollAmount` (Number): Scroll distance in pixels (default: 200)

##### Events
- `tabChange`: Dispatched when a tab is clicked, contains the new tab ID

##### Usage Examples
```html
<!-- Glass variant for Products page -->
<HorizontalTabs 
  {tabs}
  {activeTab}
  variant="glass"
  showScrollButtons={false}
  on:tabChange={handleTabChange}
/>

<!-- Standard variant with scrolling for Settings page -->
<HorizontalTabs 
  {tabs}
  {activeTab}
  variant="standard"
  showScrollButtons={true}
  scrollAmount={150}
  on:tabChange={handleTabChange}
/>
```

##### Tab Object Structure
```javascript
const tabs = [
  {
    id: 'products',
    label: 'Products',
    icon: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4'
  },
  {
    id: 'categories',
    label: 'Categories',
    icon: 'M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10'
  }
]
```

##### Implementation
```javascript
function handleTabChange(event) {
  activeTab = event.detail
  // Additional logic for tab switching
}
```

### 5.3 RULE #3: Consistent Form Layout

**Use the design system form components for consistency.**

#### ✅ Form Structure
```html
<form class="space-y-6">
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <div class="form-control">
      <label class="label-glass">Customer Name</label>
      <input type="text" class="input-glass" bind:value={customer.name} />
    </div>
    <div class="form-control">
      <label class="label-glass">Email</label>
      <input type="email" class="input-glass" bind:value={customer.email} />
    </div>
  </div>
  
  <div class="flex justify-end space-x-3">
    <button type="button" class="btn-glass">Cancel</button>
    <button type="submit" class="btn-primary">Save Customer</button>
  </div>
</form>
```

### 5.4 RULE #4: State Management

**Use Svelte's reactive features effectively.**

#### ✅ Reactive Variables
```javascript
// Reactive declarations
$: fullName = `${firstName} ${lastName}`;
$: isValid = email.includes('@') && password.length >= 8;

// Reactive statements
$: if (selectedCustomer) {
  loadCustomerData(selectedCustomer.id);
}
```

#### ✅ Stores for Global State
```javascript
// stores.js
import { writable } from 'svelte/store';

export const currentUser = writable(null);
export const notifications = writable([]);
export const theme = writable('dijibill');
```

---

## 6. Backend API

### 6.1 API Structure

All backend methods are exposed through Wails context in `app.go`:

#### Company Management
- `GetCompany()` - Get company information
- `UpdateCompany(company)` - Update company details

#### Customer Management
- `GetCustomers()` - Get all customers
- `GetCustomer(id)` - Get specific customer
- `CreateCustomer(customer)` - Create new customer
- `UpdateCustomer(customer)` - Update customer
- `DeleteCustomer(id)` - Delete customer

#### Product Management
- `GetProducts()` - Get all products
- `GetProduct(id)` - Get specific product
- `CreateProduct(product)` - Create new product
- `UpdateProduct(product)` - Update product
- `DeleteProduct(id)` - Delete product
- `GetProductCategories()` - Get product categories

#### Invoice Management
- `GetInvoices()` - Get all invoices
- `GetInvoice(id)` - Get specific invoice
- `CreateInvoice(invoice)` - Create new invoice
- `UpdateInvoice(invoice)` - Update invoice
- `DeleteInvoice(id)` - Delete invoice

#### Settings Management
- `GetTaxRates()` - Get tax rates
- `GetPaymentTypes()` - Get payment types
- `GetSalesCategories()` - Get sales categories
- `GetUnitsOfMeasurement()` - Get units
- `GetDefaultProductSettings()` - Get default product settings
- `UpdateDefaultProductSettings(settings)` - Update defaults

#### Utility Functions
- `GeneratePDF(invoiceId)` - Generate PDF for invoice
- `GenerateQRCode(data)` - Generate QR code
- `ValidateZATCAQR(qrData)` - Validate ZATCA QR compliance

### 6.2 Database Operations

Database operations are handled in `database.go`:

```go
// Example CRUD operations
func (a *App) GetProducts() ([]Product, error) {
    // Implementation
}

func (a *App) CreateProduct(product Product) error {
    // Implementation with validation
}
```

### 6.3 Error Handling

```go
// Consistent error handling pattern
func (a *App) GetCustomer(id int) (*Customer, error) {
    if id <= 0 {
        return nil, fmt.Errorf("invalid customer ID")
    }
    
    // Database operation
    customer, err := a.db.GetCustomer(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get customer: %w", err)
    }
    
    return customer, nil
}
```

---

## 7. Frontend Structure

### 7.1 Main Components

#### App.svelte
- Main application shell
- Navigation sidebar
- Route handling
- Global state management

#### Dashboard.svelte
- Business overview
- Key metrics and charts
- Recent activity
- Quick actions

#### GeneralSettings.svelte
- Horizontal scrollable tabs
- Company settings
- System preferences
- Default product settings
- Tax rates, payment types, etc.

#### Products.svelte
- Product catalog management
- Category management
- Product modal with tabs
- Search and filtering

### 7.2 Component Structure

```
components/
├── CompanySettings.svelte      # Company information form
├── SystemSettings.svelte       # System preferences
├── DefaultProductSettings.svelte # Default product values
├── TaxRateSettings.svelte      # Tax rate management
├── PaymentTypeSettings.svelte  # Payment type management
├── SalesCategorySettings.svelte # Sales category management
├── UnitSettings.svelte         # Unit of measurement settings
├── HorizontalTabs.svelte       # Universal horizontal tab navigation
├── ProductModal.svelte         # Product creation/editing
├── ProductBasicTab.svelte      # Basic product information
├── ProductPriceTab.svelte      # Pricing and tax settings
├── ProductDescriptionTab.svelte # Product descriptions
├── ProductImageTab.svelte      # Product images
└── ProductCategoryModal.svelte        # Product category management
```

### 7.3 Wails Integration

```javascript
// Import Wails functions
import { GetProducts, CreateProduct } from '../wailsjs/go/main/App.js';
import { main } from '../wailsjs/go/models';

// Use in components
async function loadProducts() {
  try {
    const products = await GetProducts();
    // Handle response
  } catch (error) {
    console.error('Failed to load products:', error);
  }
}
```

---

## 8. Development Workflow

### 8.1 Setup

```bash
# Clone repository
git clone <repository-url>
cd dijibill

# Install frontend dependencies
cd frontend
npm install

# Install Go dependencies
cd ..
go mod tidy

# Run in development mode
wails dev
```

### 8.2 Development Commands

```bash
# Frontend development server
cd frontend && npm run dev

# Wails development mode
wails dev

# Build for production
wails build

# Generate Wails bindings
wails generate

# Run tests
go test ./...
cd frontend && npm test
```

### 8.3 Code Style

#### Go Code Style
- Follow Go conventions and `gofmt`
- Use meaningful variable names
- Add comments for exported functions
- Handle errors appropriately

#### Svelte Code Style
- Use 2-space indentation
- Follow Svelte conventions
- Use design system classes
- Keep components focused and small

#### CSS/Styling
- Use design system classes first
- Avoid custom CSS when possible
- Follow BEM naming for custom classes
- Use Tailwind utilities for layout

---

## 9. Deployment

### 9.1 Build Process

```bash
# Build for current platform
wails build

# Build for specific platform
wails build -platform windows/amd64
wails build -platform darwin/amd64
wails build -platform linux/amd64
```

### 9.2 Distribution

- **Windows**: `.exe` installer
- **macOS**: `.app` bundle or `.dmg`
- **Linux**: `.AppImage` or distribution packages

### 9.3 Configuration

- Database: SQLite file in user data directory
- Logs: Application logs in system log directory
- Settings: Stored in database and local storage

---

## 10. Troubleshooting

### 10.1 Common Issues

#### Build Issues
```bash
# Clear Wails cache
wails clean

# Regenerate bindings
wails generate

# Clean and rebuild
rm -rf build/
wails build
```

#### Frontend Issues
```bash
# Clear node modules
rm -rf frontend/node_modules
cd frontend && npm install

# Clear Vite cache
rm -rf frontend/.vite
```

#### Database Issues
- Check SQLite file permissions
- Verify database schema matches models
- Check for database locks

### 10.2 Debugging

#### Frontend Debugging
- Use browser developer tools
- Check Wails console for Go errors
- Use Svelte DevTools extension

#### Backend Debugging
- Add logging to Go functions
- Use Go debugger (delve)
- Check database queries

### 10.3 Performance

#### Frontend Performance
- Use Svelte's reactive features efficiently
- Avoid unnecessary re-renders
- Optimize large lists with virtual scrolling

#### Backend Performance
- Use database indexes appropriately
- Cache frequently accessed data
- Optimize database queries

---

## 11. Saudi Arabian Compliance

### 11.1 ZATCA Requirements

- **QR Code**: ZATCA-compliant QR codes on invoices
- **VAT**: Proper VAT calculation and display
- **Arabic Support**: Bilingual invoice support
- **Digital Signatures**: E-invoice signing capability

### 11.2 Local Standards

- **Currency**: Saudi Riyal (SAR) as default
- **Date Format**: DD/MM/YYYY format
- **Language**: Arabic RTL support
- **Tax Rates**: Configurable VAT rates

---

## 12. Future Enhancements

### 12.1 Planned Features

- **Multi-company support**
- **Advanced reporting and analytics**
- **Integration with payment gateways**
- **Mobile companion app**
- **Cloud synchronization**
- **Advanced inventory management**

### 12.2 Technical Improvements

- **Database migrations system**
- **Plugin architecture**
- **Advanced caching**
- **Real-time collaboration**
- **API for third-party integrations**

---

## Contributing

When contributing to this project:

1. **Follow the design system** - Use established patterns and components
2. **Write tests** - Add unit tests for new functionality
3. **Document changes** - Update this documentation for significant changes
4. **Follow conventions** - Maintain consistent code style
5. **Test thoroughly** - Verify changes work across platforms

---

## License

[Add license information here]

---

*Last updated: [Current Date]*
*Version: 1.0.0*
