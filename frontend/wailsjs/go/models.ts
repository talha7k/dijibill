export namespace database {
	
	export class Company {
	    id: number;
	    name: string;
	    name_arabic: string;
	    vat_number: string;
	    cr_number: string;
	    email: string;
	    phone: string;
	    address: string;
	    address_arabic: string;
	    city: string;
	    city_arabic: string;
	    country: string;
	    country_arabic: string;
	    logo: string;
	
	    static createFrom(source: any = {}) {
	        return new Company(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.name_arabic = source["name_arabic"];
	        this.vat_number = source["vat_number"];
	        this.cr_number = source["cr_number"];
	        this.email = source["email"];
	        this.phone = source["phone"];
	        this.address = source["address"];
	        this.address_arabic = source["address_arabic"];
	        this.city = source["city"];
	        this.city_arabic = source["city_arabic"];
	        this.country = source["country"];
	        this.country_arabic = source["country_arabic"];
	        this.logo = source["logo"];
	    }
	}
	export class Customer {
	    id: number;
	    name: string;
	    name_arabic: string;
	    vat_number: string;
	    email: string;
	    phone: string;
	    address: string;
	    address_arabic: string;
	    city: string;
	    city_arabic: string;
	    country: string;
	    country_arabic: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Customer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.name_arabic = source["name_arabic"];
	        this.vat_number = source["vat_number"];
	        this.email = source["email"];
	        this.phone = source["phone"];
	        this.address = source["address"];
	        this.address_arabic = source["address_arabic"];
	        this.city = source["city"];
	        this.city_arabic = source["city_arabic"];
	        this.country = source["country"];
	        this.country_arabic = source["country_arabic"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, clazz: any, asMap: boolean = false): any {
	    if (!a) {
	        return a;
	    }
	    if (a.slice && a.map) {
	        return (a as any[]).map(elem => this.convertValues(elem, clazz));
	    } else if ("object" === typeof a) {
	        if (asMap) {
	            for (const key of Object.keys(a)) {
	                a[key] = clazz && typeof clazz === 'function' ? new clazz(a[key]) : a[key];
	            }
	            return a;
	        }
	        return clazz && typeof clazz === 'function' ? new clazz(a) : a;
	    }
	    return a;
	}
	}
	export class DefaultProductSettings {
	    id: number;
	    default_stock: number;
	    default_tax_rate_id: number;
	    default_unit_id: number;
	    default_product_type: string;
	    default_product_status: boolean;
	    default_markup: number;
	    default_price_includes_tax: boolean;
	    default_price_change_allowed: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new DefaultProductSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.default_stock = source["default_stock"];
	        this.default_tax_rate_id = source["default_tax_rate_id"];
	        this.default_unit_id = source["default_unit_id"];
	        this.default_product_type = source["default_product_type"];
	        this.default_product_status = source["default_product_status"];
	        this.default_markup = source["default_markup"];
	        this.default_price_includes_tax = source["default_price_includes_tax"];
	        this.default_price_change_allowed = source["default_price_change_allowed"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, clazz: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, clazz));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = clazz && typeof clazz === 'function' ? new clazz(a[key]) : a[key];
		            }
		            return a;
		        }
		        return clazz && typeof clazz === 'function' ? new clazz(a) : a;
		    }
		    return a;
		}
	}
	export class PaymentType {
	    id: number;
	    name: string;
	    name_arabic: string;
	    code: string;
	    is_default: boolean;
	    is_active: boolean;
	    description: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new PaymentType(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.name_arabic = source["name_arabic"];
	        this.code = source["code"];
	        this.is_default = source["is_default"];
	        this.is_active = source["is_active"];
	        this.description = source["description"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, clazz: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, clazz));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = clazz && typeof clazz === 'function' ? new clazz(a[key]) : a[key];
		            }
		            return a;
		        }
		        return clazz && typeof clazz === 'function' ? new clazz(a) : a;
		    }
		    return a;
		}
	}
	export class ProductCategory {
	    id: number;
	    name: string;
	    name_arabic: string;
	    description: string;
	    description_arabic: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new ProductCategory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.name_arabic = source["name_arabic"];
	        this.description = source["description"];
	        this.description_arabic = source["description_arabic"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, clazz: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, clazz));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = clazz && typeof clazz === 'function' ? new clazz(a[key]) : a[key];
		            }
		            return a;
		        }
		        return clazz && typeof clazz === 'function' ? new clazz(a) : a;
		    }
		    return a;
		}
	}
	export class Product {
	    id: number;
	    name: string;
	    name_arabic: string;
	    description: string;
	    description_arabic: string;
	    category_id: number;
	    category_name: string;
	    category?: ProductCategory;
	    unit_price: number;
	    vat_rate: number;
	    unit: string;
	    unit_arabic: string;
	    sku: string;
	    barcode: string;
	    stock: number;
	    min_stock: number;
	    is_active: boolean;
	    color: string;
	    image_url: string;
	    service_not_using_stock: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Product(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.name_arabic = source["name_arabic"];
	        this.description = source["description"];
	        this.description_arabic = source["description_arabic"];
	        this.category_id = source["category_id"];
	        this.category_name = source["category_name"];
	        this.category = this.convertValues(source["category"], ProductCategory);
	        this.unit_price = source["unit_price"];
	        this.vat_rate = source["vat_rate"];
	        this.unit = source["unit"];
	        this.unit_arabic = source["unit_arabic"];
	        this.sku = source["sku"];
	        this.barcode = source["barcode"];
	        this.stock = source["stock"];
	        this.min_stock = source["min_stock"];
	        this.is_active = source["is_active"];
	        this.color = source["color"];
	        this.image_url = source["image_url"];
	        this.service_not_using_stock = source["service_not_using_stock"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, clazz: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, clazz));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = clazz && typeof clazz === 'function' ? new clazz(a[key]) : a[key];
		            }
		            return a;
		        }
		        return clazz && typeof clazz === 'function' ? new clazz(a) : a;
		    }
		    return a;
		}
	}
	export class SalesInvoiceItem {
	    id: number;
	    invoice_id: number;
	    product_id: number;
	    product?: Product;
	    quantity: number;
	    unit_price: number;
	    vat_rate: number;
	    vat_amount: number;
	    total_amount: number;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new SalesInvoiceItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.invoice_id = source["invoice_id"];
	        this.product_id = source["product_id"];
	        this.product = this.convertValues(source["product"], Product);
	        this.quantity = source["quantity"];
	        this.unit_price = source["unit_price"];
	        this.vat_rate = source["vat_rate"];
	        this.vat_amount = source["vat_amount"];
	        this.total_amount = source["total_amount"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, clazz: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, clazz));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = clazz && typeof clazz === 'function' ? new clazz(a[key]) : a[key];
		            }
		            return a;
		        }
		        return clazz && typeof clazz === 'function' ? new clazz(a) : a;
		    }
		    return a;
		}
	}
	export class SalesCategory {
	    id: number;
	    name: string;
	    name_arabic: string;
	    code: string;
	    description: string;
	    description_arabic: string;
	    is_default: boolean;
	    is_active: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new SalesCategory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.name_arabic = source["name_arabic"];
	        this.code = source["code"];
	        this.description = source["description"];
	        this.description_arabic = source["description_arabic"];
	        this.is_default = source["is_default"];
	        this.is_active = source["is_active"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, clazz: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, clazz));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = clazz && typeof clazz === 'function' ? new clazz(a[key]) : a[key];
		            }
		            return a;
		        }
		        return clazz && typeof clazz === 'function' ? new clazz(a) : a;
		    }
		    return a;
		}
	}
	export class SalesInvoice {
	    id: number;
	    invoice_number: string;
	    customer_id: number;
	    customer?: Customer;
	    sales_category_id: number;
	    sales_category?: SalesCategory;
	    // Go type: time
	    issue_date: any;
	    // Go type: time
	    due_date: any;
	    sub_total: number;
	    vat_amount: number;
	    total_amount: number;
	    status: string;
	    notes: string;
	    notes_arabic: string;
	    qr_code: string;
	    items?: SalesInvoiceItem[];
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new SalesInvoice(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.invoice_number = source["invoice_number"];
	        this.customer_id = source["customer_id"];
	        this.customer = this.convertValues(source["customer"], Customer);
	        this.sales_category_id = source["sales_category_id"];
	        this.sales_category = this.convertValues(source["sales_category"], SalesCategory);
	        this.issue_date = this.convertValues(source["issue_date"], null);
	        this.due_date = this.convertValues(source["due_date"], null);
	        this.sub_total = source["sub_total"];
	        this.vat_amount = source["vat_amount"];
	        this.total_amount = source["total_amount"];
	        this.status = source["status"];
	        this.notes = source["notes"];
	        this.notes_arabic = source["notes_arabic"];
	        this.qr_code = source["qr_code"];
	        this.items = this.convertValues(source["items"], SalesInvoiceItem);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, clazz: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, clazz));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = clazz && typeof clazz === 'function' ? new clazz(a[key]) : a[key];
		            }
		            return a;
		        }
		        return clazz && typeof clazz === 'function' ? new clazz(a) : a;
		    }
		    return a;
		}
	}
	export class Payment {
	    id: number;
	    invoice_id: number;
	    invoice?: SalesInvoice;
	    payment_type_id: number;
	    payment_type?: PaymentType;
	    amount: number;
	    // Go type: time
	    payment_date: any;
	    reference: string;
	    notes: string;
	    notes_arabic: string;
	    status: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Payment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.invoice_id = source["invoice_id"];
	        this.invoice = this.convertValues(source["invoice"], SalesInvoice);
	        this.payment_type_id = source["payment_type_id"];
	        this.payment_type = this.convertValues(source["payment_type"], PaymentType);
	        this.amount = source["amount"];
	        this.payment_date = this.convertValues(source["payment_date"], null);
	        this.reference = source["reference"];
	        this.notes = source["notes"];
	        this.notes_arabic = source["notes_arabic"];
	        this.status = source["status"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, clazz: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, clazz));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = clazz && typeof clazz === 'function' ? new clazz(a[key]) : a[key];
		            }
		            return a;
		        }
		        return clazz && typeof clazz === 'function' ? new clazz(a) : a;
		    }
		    return a;
		}
	}
	
	
	
	export class PurchaseInvoiceItem {
	    id: number;
	    invoice_id: number;
	    product_id: number;
	    product?: Product;
	    quantity: number;
	    unit_price: number;
	    vat_rate: number;
	    vat_amount: number;
	    total_amount: number;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new PurchaseInvoiceItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.invoice_id = source["invoice_id"];
	        this.product_id = source["product_id"];
	        this.product = this.convertValues(source["product"], Product);
	        this.quantity = source["quantity"];
	        this.unit_price = source["unit_price"];
	        this.vat_rate = source["vat_rate"];
	        this.vat_amount = source["vat_amount"];
	        this.total_amount = source["total_amount"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, clazz: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, clazz));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = clazz && typeof clazz === 'function' ? new clazz(a[key]) : a[key];
		            }
		            return a;
		        }
		        return clazz && typeof clazz === 'function' ? new clazz(a) : a;
		    }
		    return a;
		}
	}
	export class Supplier {
	    id: number;
	    company_name: string;
	    company_name_arabic: string;
	    contact_person: string;
	    contact_person_arabic: string;
	    vat_number: string;
	    email: string;
	    phone: string;
	    address: string;
	    address_arabic: string;
	    city: string;
	    city_arabic: string;
	    country: string;
	    country_arabic: string;
	    payment_terms: string;
	    active: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Supplier(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.company_name = source["company_name"];
	        this.company_name_arabic = source["company_name_arabic"];
	        this.contact_person = source["contact_person"];
	        this.contact_person_arabic = source["contact_person_arabic"];
	        this.vat_number = source["vat_number"];
	        this.email = source["email"];
	        this.phone = source["phone"];
	        this.address = source["address"];
	        this.address_arabic = source["address_arabic"];
	        this.city = source["city"];
	        this.city_arabic = source["city_arabic"];
	        this.country = source["country"];
	        this.country_arabic = source["country_arabic"];
	        this.payment_terms = source["payment_terms"];
	        this.active = source["active"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, clazz: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, clazz));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = clazz && typeof clazz === 'function' ? new clazz(a[key]) : a[key];
		            }
		            return a;
		        }
		        return clazz && typeof clazz === 'function' ? new clazz(a) : a;
		    }
		    return a;
		}
	}
	export class PurchaseInvoice {
	    id: number;
	    invoice_number: string;
	    supplier_id: number;
	    supplier?: Supplier;
	    // Go type: time
	    issue_date: any;
	    // Go type: time
	    due_date: any;
	    sub_total: number;
	    vat_amount: number;
	    total_amount: number;
	    status: string;
	    notes: string;
	    notes_arabic: string;
	    items?: PurchaseInvoiceItem[];
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new PurchaseInvoice(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.invoice_number = source["invoice_number"];
	        this.supplier_id = source["supplier_id"];
	        this.supplier = this.convertValues(source["supplier"], Supplier);
	        this.issue_date = this.convertValues(source["issue_date"], null);
	        this.due_date = this.convertValues(source["due_date"], null);
	        this.sub_total = source["sub_total"];
	        this.vat_amount = source["vat_amount"];
	        this.total_amount = source["total_amount"];
	        this.status = source["status"];
	        this.notes = source["notes"];
	        this.notes_arabic = source["notes_arabic"];
	        this.items = this.convertValues(source["items"], PurchaseInvoiceItem);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, clazz: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, clazz));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = clazz && typeof clazz === 'function' ? new clazz(a[key]) : a[key];
		            }
		            return a;
		        }
		        return clazz && typeof clazz === 'function' ? new clazz(a) : a;
		    }
		    return a;
		}
	}
	
	
	
	
	
	export class TaxRate {
	    id: number;
	    name: string;
	    name_arabic: string;
	    rate: number;
	    is_default: boolean;
	    is_active: boolean;
	    description: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new TaxRate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.name_arabic = source["name_arabic"];
	        this.rate = source["rate"];
	        this.is_default = source["is_default"];
	        this.is_active = source["is_active"];
	        this.description = source["description"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, clazz: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, clazz));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = clazz && typeof clazz === 'function' ? new clazz(a[key]) : a[key];
		            }
		            return a;
		        }
		        return clazz && typeof clazz === 'function' ? new clazz(a) : a;
		    }
		    return a;
		}
	}
	export class UnitOfMeasurement {
	    id: number;
	    value: string;
	    label: string;
	    arabic: string;
	    is_default: boolean;
	    is_active: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new UnitOfMeasurement(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.value = source["value"];
	        this.label = source["label"];
	        this.arabic = source["arabic"];
	        this.is_default = source["is_default"];
	        this.is_active = source["is_active"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, clazz: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, clazz));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = clazz && typeof clazz === 'function' ? new clazz(a[key]) : a[key];
		            }
		            return a;
		        }
		        return clazz && typeof clazz === 'function' ? new clazz(a) : a;
		    }
		    return a;
		}
	}

}

