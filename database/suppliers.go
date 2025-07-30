package database

// Supplier operations
func (d *Database) CreateSupplier(supplier *Supplier) error {
	query := `
		INSERT INTO suppliers (company_name, company_name_arabic, contact_person, contact_person_arabic, 
		vat_number, email, phone, address, address_arabic, city, city_arabic, country, country_arabic, 
		payment_terms, active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := d.db.Exec(query, supplier.CompanyName, supplier.CompanyNameArabic, 
		supplier.ContactPerson, supplier.ContactPersonArabic, supplier.VATNumber,
		supplier.Email, supplier.Phone, supplier.Address, supplier.AddressArabic,
		supplier.City, supplier.CityArabic, supplier.Country, supplier.CountryArabic,
		supplier.PaymentTerms, supplier.Active)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	supplier.ID = int(id)
	return nil
}

func (d *Database) GetSuppliers() ([]Supplier, error) {
	query := `SELECT id, company_name, company_name_arabic, contact_person, contact_person_arabic, 
			  vat_number, email, phone, address, address_arabic, city, city_arabic, country, 
			  country_arabic, payment_terms, active, created_at, updated_at FROM suppliers ORDER BY company_name`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []Supplier
	for rows.Next() {
		var s Supplier
		err := rows.Scan(&s.ID, &s.CompanyName, &s.CompanyNameArabic, &s.ContactPerson, 
			&s.ContactPersonArabic, &s.VATNumber, &s.Email, &s.Phone, &s.Address, 
			&s.AddressArabic, &s.City, &s.CityArabic, &s.Country, &s.CountryArabic,
			&s.PaymentTerms, &s.Active, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		suppliers = append(suppliers, s)
	}
	return suppliers, nil
}

func (d *Database) GetSupplierByID(id int) (*Supplier, error) {
	query := `SELECT id, company_name, company_name_arabic, contact_person, contact_person_arabic, 
			  vat_number, email, phone, address, address_arabic, city, city_arabic, country, 
			  country_arabic, payment_terms, active, created_at, updated_at FROM suppliers WHERE id = ?`

	var s Supplier
	err := d.db.QueryRow(query, id).Scan(&s.ID, &s.CompanyName, &s.CompanyNameArabic, 
		&s.ContactPerson, &s.ContactPersonArabic, &s.VATNumber, &s.Email, &s.Phone, 
		&s.Address, &s.AddressArabic, &s.City, &s.CityArabic, &s.Country, &s.CountryArabic,
		&s.PaymentTerms, &s.Active, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (d *Database) UpdateSupplier(supplier *Supplier) error {
	query := `
		UPDATE suppliers SET company_name = ?, company_name_arabic = ?, contact_person = ?, 
		contact_person_arabic = ?, vat_number = ?, email = ?, phone = ?, address = ?, 
		address_arabic = ?, city = ?, city_arabic = ?, country = ?, country_arabic = ?, 
		payment_terms = ?, active = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err := d.db.Exec(query, supplier.CompanyName, supplier.CompanyNameArabic,
		supplier.ContactPerson, supplier.ContactPersonArabic, supplier.VATNumber,
		supplier.Email, supplier.Phone, supplier.Address, supplier.AddressArabic,
		supplier.City, supplier.CityArabic, supplier.Country, supplier.CountryArabic,
		supplier.PaymentTerms, supplier.Active, supplier.ID)
	return err
}

func (d *Database) DeleteSupplier(id int) error {
	_, err := d.db.Exec("DELETE FROM suppliers WHERE id = ?", id)
	return err
}