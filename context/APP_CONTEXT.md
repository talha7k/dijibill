
# DIJINVOICE: Svelte & Go Invoicing App

This document provides the complete context, data models, and development rules for the project.

---

## 1. Project Overview

This project is a desktop invoicing and business management application built with **Go** for the backend, **Svelte** for the frontend, and **Wails** to package it as a cross-platform desktop application. The application is designed for small to medium-sized businesses, with a strong focus on the **Saudi Arabian market**, requiring bilingual (**English/Arabic**) support and adherence to local invoicing standards (e.g., VAT and QR codes).

The core purpose of the application is to manage company information, customers, products, and to create, track, and manage invoices efficiently.

---

## 2. Technical Stack

* **Backend:** Go
* **Frontend:** Svelte
* **Desktop Framework:** Wails (bridges Go and Svelte)
* **Database:** SQLite
* **UI Framework:** Tailwind CSS
* **UI Component Library:** DaisyUI

---

## 3. Core Functionality & Data Models

LOOK AT `models.go` file to see data models.

## 4. RULE #1: Strict DaisyUI Component Usage

When building UI, **ALWAYS** prioritize DaisyUI components over custom CSS or other UI libraries. This is the most important rule for frontend development.

### Benefits

  * **Consistency:** Ensures a single, consistent design system.
  * **Efficiency:** Built on Tailwind CSS for rapid development.
  * **Maintainability:** Semantic class names (`btn`, `card`, `modal`) improve readability.
  * **Accessibility & Theming:** Provides A11Y features and a robust theming system out of the box.

### Implementation Guidelines

1.  **Always check the [DaisyUI documentation](https://daisyui.com/components/) first** before writing any custom styles.
2.  **Use semantic class names** provided by DaisyUI (e.g., `btn-primary`, `alert-success`, `card-title`).
3.  **Use DaisyUI theme colors** (`primary`, `secondary`, `accent`, `info`, `success`, `warning`, `error`).
4.  **Combine with Tailwind utilities** for spacing (`m-4`), layout (`flex`), and responsive design (`md:`, `lg:`).

### Example

**Correct: Using DaisyUI components.**

```html
<div class="card bg-base-100 shadow-xl">
  <div class="card-body">
    <h2 class="card-title">Add Product</h2>
    <div class="form-control">
      <input type="text" class="input input-bordered" />
    </div>
    <div class="card-actions justify-end">
      <button class="btn btn-primary">Save</button>
    </div>
  </div>
</div>
```

**Incorrect: Custom CSS when a DaisyUI equivalent exists.**

```html
<style>
  .custom-card { box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
  .custom-button { background-color: #007bff; }
</style>
<div class="custom-card"><button class="custom-button">Save</button></div>
```

-----

## 5\. RULE \#2: Favor Composition Over Monolithic Components

**Break down complex UI into smaller, single-purpose, reusable components.** Avoid creating large, monolithic Svelte files that handle too many different responsibilities.

### Benefits

  * **Maintainability:** Smaller files are easier to understand, debug, and modify.
  * **Reusability:** A generic `DataInput.svelte` can be used across many forms.
  * **Testability:** It's simpler to write unit tests for small, focused components.
  * **Readability:** Parent components become more declarative, describing *what* the UI contains, not *how* every part is built.

### Example

A page for creating an invoice should not be a single, massive file.

**Avoid (Monolithic Component):** A single `InvoiceEditor.svelte` file that contains all markup and logic for the customer selector, the product line items table, date pickers, and final total calculations.

**Correct (Composed Components):** Break `InvoiceEditor.svelte` into smaller, focused components that communicate with each other.

  * `InvoiceEditor.svelte` (The main page/container)
      * Imports and uses `CustomerSelector.svelte`
      * Imports and uses `InvoiceItemTable.svelte`
      * Imports and uses `InvoiceTotals.svelte`

This makes the code in `InvoiceEditor.svelte` clean and descriptive:

```html
<script>
  import CustomerSelector from './CustomerSelector.svelte';
  import InvoiceItemTable from './InvoiceItemTable.svelte';
  import InvoiceTotals from './InvoiceTotals.svelte';
  
  let invoiceItems = [];
  let selectedCustomer = null;

  // ... logic to coordinate these components ...
</script>

<CustomerSelector on:select={(e) => selectedCustomer = e.detail} />

<InvoiceItemTable bind:items={invoiceItems} />

<InvoiceTotals {invoiceItems} />

<div class="card-actions justify-end">
    <button class="btn btn-primary">Create Invoice</button>
</div>  