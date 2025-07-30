
# DIJINVOICE: Svelte & Go Invoicing App

A cross-platform desktop invoicing and business management app for SMEs, built with **Go**, **Svelte**, and **Wails**, tailored for the **Saudi Arabian market** with bilingual (English/Arabic) support and local invoice standards (VAT, QR).

---

## 1. Tech Stack

- **Backend:** Go  
- **Frontend:** Svelte  
- **Desktop Framework:** Wails  
- **Database:** SQLite  
- **Styling:** Tailwind CSS + DaisyUI  

---

## 2. Data Models

See `models.go` for all data structures.

---

## 3. RULE #1: Use DaisyUI First

**Always use DaisyUI components before writing custom CSS.**

### ✅ Benefits
- Consistent UI  
- Faster dev with Tailwind  
- Readable, semantic class names  
- Built-in A11Y + theming  

### ✅ Guidelines
- Check [DaisyUI docs](https://daisyui.com/components/) first  
- Use semantic classes like `btn`, `input`, `card`  
- Stick to DaisyUI theme colors: `primary`, `success`, etc.  
- Combine with Tailwind for layout/spacing (`flex`, `mb-4`, `md:grid`, etc.)

### ✅ Example

```html
<div class="card bg-base-100 shadow-xl">
  <div class="card-body">
    <h2 class="card-title">Add Product</h2>
    <div class="form-control">
      <label class="label"><span class="label-text">Product Name</span></label>
      <input type="text" class="input input-bordered" />
    </div>
    <div class="card-actions justify-end">
      <button class="btn btn-primary">Save</button>
    </div>
  </div>
</div>

## 4. RULE #2: Compose, Don’t Cram

**Break complex UIs into smaller reusable Svelte components.**

### ✅ Structure

<InvoiceEditor>
  <CustomerSelector />
  <InvoiceItemTable />
  <InvoiceTotals />
  <div class="card-actions justify-end">
    <button class="btn btn-primary">Create Invoice</button>
  </div>
</InvoiceEditor>

Each component should have one responsibility. Avoid monolithic `.svelte` files.


## 5. RULE #3: Consistent Form Layout

**Labels above fields. Everything left-aligned.**

### ✅ Guidelines

* Use `form-control`
* Use `label` with `label-text` above each input
* Add vertical spacing with `mb-4`
* Avoid inline labels or placeholders as labels

### ✅ Example

<div class="form-control mb-4">
  <label class="label"><span class="label-text">Customer Name</span></label>
  <input type="text" class="input input-bordered w-full" />
</div>
