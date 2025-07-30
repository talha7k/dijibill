package database

import (
	"log"
)

// insertDefaultSettings inserts default settings data
func (d *Database) insertDefaultSettings() error {
	// Insert default payment types
	var paymentCount int
	err := d.db.QueryRow("SELECT COUNT(*) FROM payment_types").Scan(&paymentCount)
	if err != nil {
		return err
	}

	if paymentCount == 0 {
		paymentTypes := []PaymentType{
			{Name: "Cash", NameArabic: "نقدي", Code: "cash", Description: "Cash payment", IsDefault: true, IsActive: true},
			{Name: "Credit Card", NameArabic: "بطاقة ائتمان", Code: "card", Description: "Credit card payment", IsDefault: false, IsActive: true},
			{Name: "Bank Transfer", NameArabic: "تحويل بنكي", Code: "bank_transfer", Description: "Bank transfer payment", IsDefault: false, IsActive: true},
		}

		for _, pt := range paymentTypes {
			if err := d.CreatePaymentType(&pt); err != nil {
				log.Printf("Warning: Could not insert payment type %s: %v", pt.Name, err)
			}
		}
	}

	// Insert default sales categories
	var salesCount int
	err = d.db.QueryRow("SELECT COUNT(*) FROM sales_categories").Scan(&salesCount)
	if err != nil {
		return err
	}

	if salesCount == 0 {
		salesCategories := []SalesCategory{
			{Name: "Retail", NameArabic: "تجزئة", Code: "retail", Description: "Retail sales", DescriptionArabic: "مبيعات التجزئة", IsDefault: true, IsActive: true},
			{Name: "Wholesale", NameArabic: "جملة", Code: "wholesale", Description: "Wholesale sales", DescriptionArabic: "مبيعات الجملة", IsDefault: false, IsActive: true},
			{Name: "Service", NameArabic: "خدمة", Code: "service", Description: "Service sales", DescriptionArabic: "مبيعات الخدمات", IsDefault: false, IsActive: true},
		}

		for _, sc := range salesCategories {
			if err := d.CreateSalesCategory(&sc); err != nil {
				log.Printf("Warning: Could not insert sales category %s: %v", sc.Name, err)
			}
		}
	}

	// Insert default tax rates
	var taxCount int
	err = d.db.QueryRow("SELECT COUNT(*) FROM tax_rates").Scan(&taxCount)
	if err != nil {
		return err
	}

	if taxCount == 0 {
		taxRates := []TaxRate{
			{Name: "Standard VAT", NameArabic: "ضريبة القيمة المضافة", Rate: 15.0, Description: "Standard VAT rate", IsDefault: true, IsActive: true},
			{Name: "Zero VAT", NameArabic: "معفى من الضريبة", Rate: 0.0, Description: "Zero VAT rate", IsDefault: false, IsActive: true},
			{Name: "Exempt", NameArabic: "معفى", Rate: 0.0, Description: "VAT exempt", IsDefault: false, IsActive: true},
		}

		for _, tr := range taxRates {
			if err := d.CreateTaxRate(&tr); err != nil {
				log.Printf("Warning: Could not insert tax rate %s: %v", tr.Name, err)
			}
		}
	}

	// Insert default units of measurement
	var unitCount int
	err = d.db.QueryRow("SELECT COUNT(*) FROM units_of_measurement").Scan(&unitCount)
	if err != nil {
		return err
	}

	if unitCount == 0 {
		units := []UnitOfMeasurement{
			{Value: "pcs", Label: "Pieces", Arabic: "قطعة", IsDefault: true, IsActive: true},
			{Value: "kg", Label: "Kilograms", Arabic: "كيلوغرام", IsDefault: false, IsActive: true},
			{Value: "m", Label: "Meters", Arabic: "متر", IsDefault: false, IsActive: true},
			{Value: "l", Label: "Liters", Arabic: "لتر", IsDefault: false, IsActive: true},
			{Value: "box", Label: "Box", Arabic: "صندوق", IsDefault: false, IsActive: true},
			{Value: "pack", Label: "Pack", Arabic: "عبوة", IsDefault: false, IsActive: true},
		}

		for _, u := range units {
			if err := d.CreateUnitOfMeasurement(&u); err != nil {
				log.Printf("Warning: Could not insert unit %s: %v", u.Label, err)
			}
		}
	}

	// Insert default product settings
	var defaultSettingsCount int
	err = d.db.QueryRow("SELECT COUNT(*) FROM default_product_settings").Scan(&defaultSettingsCount)
	if err != nil {
		return err
	}

	if defaultSettingsCount == 0 {
		// Get default IDs for foreign keys
		var defaultTaxRateID, defaultUnitID, defaultSalesCategoryID int
		
		d.db.QueryRow("SELECT id FROM tax_rates WHERE is_default = 1 LIMIT 1").Scan(&defaultTaxRateID)
		d.db.QueryRow("SELECT id FROM units_of_measurement WHERE is_default = 1 LIMIT 1").Scan(&defaultUnitID)
		d.db.QueryRow("SELECT id FROM sales_categories WHERE is_default = 1 LIMIT 1").Scan(&defaultSalesCategoryID)

		_, err = d.db.Exec(`
			INSERT INTO default_product_settings (
				default_stock, default_tax_rate_id, default_unit_id, 
				default_product_type, default_product_status, 
				default_markup, default_price_includes_tax, default_price_change_allowed
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			1, defaultTaxRateID, defaultUnitID,
			"product", true, 0.0, false, true)
		
		if err != nil {
			log.Printf("Warning: Could not insert default product settings: %v", err)
		}
	}

	return nil
}

// insertDefaultCompany inserts a default company record
func (d *Database) insertDefaultCompany() error {
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM companies").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		_, err = d.db.Exec(`
			INSERT INTO companies (name, name_arabic, vat_number, cr_number, email, phone, address, address_arabic, city, city_arabic, country, country_arabic)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			"Your Company Name", "اسم شركتك", "123456789012345", "1234567890",
			"info@company.com", "+966501234567",
			"123 Business Street", "شارع الأعمال ١٢٣",
			"Riyadh", "الرياض", "Saudi Arabia", "المملكة العربية السعودية")
		return err
	}
	return nil
}