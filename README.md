# Dijibill: Cross-Platform Invoicing & Business Management for SMEs

Dijibill is a desktop application designed to streamline invoicing and business operations for Small and Medium-sized Enterprises (SMEs). Built with a focus on local compliance and user experience, it offers bilingual support (English/Arabic) and adheres to Saudi Arabian invoice standards, including VAT calculation and ZATCA-compliant QR codes.

---

## Features

* **Invoice Generation**: Create and manage sales invoices with line items, tax calculations, and QR codes.
* **Customer Management**: Maintain a comprehensive database of customer details.
* **Product Catalog**: Manage products and services with pricing, categories, and units of measurement.
* **Company Settings**: Configure business information, tax rates, payment types, and default settings.
* **Local Compliance**: Full support for Saudi Arabian VAT standards, ZATCA QR code generation and validation, and bilingual output (English/Arabic).
* **Cross-Platform**: Available on Windows, macOS, and Linux.
* **Intuitive UI**: A modern, responsive interface leveraging a custom design system with glass morphism effects.

---

## Tech Stack

* **Backend**: Go 1.21+
* **Frontend**: Svelte 4 + Vite
* **Desktop Framework**: Wails v2
* **Database**: SQLite (built-in)
* **Styling**: Tailwind CSS + DaisyUI + Custom Design System

---

## Installation & Setup

To get Dijibill up and running on your local machine:

### Prerequisites

* Go 1.21+
* Node.js (LTS version) & npm
* Wails v2 (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Steps

1. **Clone the repository**:
   ```bash
   git clone https://github.com/talha7k/dijibill.git
   cd dijibill
   ```

2. **Install frontend dependencies**:
   ```bash
   cd frontend
   npm install
   cd ..
   ```

3. **Install Go dependencies**:
   ```bash
   go mod tidy
   ```

4. **Run in development mode**:
   ```bash
   wails dev
   ```
   This command will start the Go backend and the Svelte frontend development server, opening the application in a new window.

---

## Development Workflow

### Useful Commands

* **Run in development mode**:
  ```bash
  wails dev
  ```
* **Build for production**:
  ```bash
  wails build
  ```
  To build for a specific platform (e.g., Windows):
  ```bash
  wails build -platform windows/amd64
  ```
* **Generate Wails bindings (if Go method signatures change)**:
  ```bash
  wails generate
  ```
* **Clean Wails build cache**:
  ```bash
  wails clean
  ```

### Code Style Guidelines

* **Go**: Adhere to `gofmt` and standard Go conventions. Ensure clear error handling.
* **Svelte**: Use 2-space indentation. Prioritize Svelte's reactivity and component composition.
* **Styling**: **Always use the defined Design System classes first.** Only use custom CSS or Tailwind utilities for layout and spacing when design system components aren't sufficient.

---

## Project Structure Overview

```
dijibill/
├── app.go                 # Wails app context and API methods
├── database.go           # Database operations and CRUD logic
├── main.go              # Application entry point
├── models.go            # Data structures and types
├── pdf_service.go       # PDF generation service for invoices
├── qr_service.go        # QR code generation service (ZATCA compliant)
├── frontend/
│   ├── src/
│   │   ├── App.svelte           # Main application shell and routing
│   │   ├── Dashboard.svelte     # Dashboard view
│   │   ├── GeneralSettings.svelte # Settings management
│   │   ├── Products.svelte      # Product management view
│   │   ├── Customer.svelte      # Customer management view
│   │   ├── SalesInvoices.svelte # Invoice management view
│   │   ├── components/          # Reusable Svelte UI components
│   │   ├── styles/             # Design system and global styles
│   │   └── wailsjs/            # Generated Wails bindings (Go -> JS)
│   ├── package.json
│   ├── tailwind.config.js
│   └── vite.config.js
```

---

## Saudi Arabian Compliance Highlights

Dijibill is specifically tailored for the Saudi Arabian market, focusing on:

* **ZATCA Compliant QR Codes**: Invoices include QR codes meeting the ZATCA e-invoicing specifications.
* **VAT Handling**: Accurate calculation and display of Value Added Tax (VAT) as per local regulations.
* **Bilingual Support**: Full English and Arabic language support throughout the application and for generated invoices.
* **Local Standards**: Default currency (SAR), date format (DD/MM/YYYY), and RTL support for Arabic.

---

## Contributing

We welcome contributions to Dijibill! Please ensure you:

1. **Follow the Design System**: Utilize existing components and styling guidelines.
2. **Write Tests**: Include unit tests for new features and bug fixes.
3. **Document Changes**: Update relevant sections of the documentation for significant additions or alterations.
4. **Maintain Code Style**: Adhere to the established Go and Svelte code styles.
5. **Test Thoroughly**: Verify your changes work correctly across different platforms if applicable.

---

## License

[Add license information here, e.g., MIT License, Apache 2.0, etc.]

---
