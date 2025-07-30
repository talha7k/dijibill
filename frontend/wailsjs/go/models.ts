export namespace main {
	
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
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SaleItem {
	    id: number;
	    name: string;
	    name_arabic: string;
	    description: string;
	    description_arabic: string;
	    category_id: number;
	    unit_price: number;
	    vat_rate: number;
	    unit: string;
	    unit_arabic: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new SaleItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.name_arabic = source["name_arabic"];
	        this.description = source["description"];
	        this.description_arabic = source["description_arabic"];
	        this.category_id = source["category_id"];
	        this.unit_price = source["unit_price"];
	        this.vat_rate = source["vat_rate"];
	        this.unit = source["unit"];
	        this.unit_arabic = source["unit_arabic"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class InvoiceItem {
	    id: number;
	    invoice_id: number;
	    sale_item_id: number;
	    sale_item?: SaleItem;
	    quantity: number;
	    unit_price: number;
	    vat_rate: number;
	    vat_amount: number;
	    total_amount: number;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new InvoiceItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.invoice_id = source["invoice_id"];
	        this.sale_item_id = source["sale_item_id"];
	        this.sale_item = this.convertValues(source["sale_item"], SaleItem);
	        this.quantity = source["quantity"];
	        this.unit_price = source["unit_price"];
	        this.vat_rate = source["vat_rate"];
	        this.vat_amount = source["vat_amount"];
	        this.total_amount = source["total_amount"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Invoice {
	    id: number;
	    invoice_number: string;
	    customer_id: number;
	    customer?: Customer;
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
	    items?: InvoiceItem[];
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Invoice(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.invoice_number = source["invoice_number"];
	        this.customer_id = source["customer_id"];
	        this.customer = this.convertValues(source["customer"], Customer);
	        this.issue_date = this.convertValues(source["issue_date"], null);
	        this.due_date = this.convertValues(source["due_date"], null);
	        this.sub_total = source["sub_total"];
	        this.vat_amount = source["vat_amount"];
	        this.total_amount = source["total_amount"];
	        this.status = source["status"];
	        this.notes = source["notes"];
	        this.notes_arabic = source["notes_arabic"];
	        this.qr_code = source["qr_code"];
	        this.items = this.convertValues(source["items"], InvoiceItem);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class ItemCategory {
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
	        return new ItemCategory(source);
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
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

