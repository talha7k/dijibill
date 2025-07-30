package database

// Company operations
func (d *Database) GetCompany() (*Company, error) {
	query := `SELECT id, name, name_arabic, vat_number, cr_number, email, phone, address, address_arabic, city, city_arabic, country, country_arabic, COALESCE(logo, '') as logo FROM companies LIMIT 1`

	var c Company
	err := d.db.QueryRow(query).Scan(&c.ID, &c.Name, &c.NameArabic, &c.VATNumber, &c.CRNumber, &c.Email, &c.Phone,
		&c.Address, &c.AddressArabic, &c.City, &c.CityArabic, &c.Country, &c.CountryArabic, &c.Logo)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (d *Database) UpdateCompany(company *Company) error {
	query := `
		UPDATE companies SET name = ?, name_arabic = ?, vat_number = ?, cr_number = ?, email = ?, phone = ?, 
		address = ?, address_arabic = ?, city = ?, city_arabic = ?, country = ?, country_arabic = ?, logo = NULLIF(?, '')
		WHERE id = ?`

	_, err := d.db.Exec(query, company.Name, company.NameArabic, company.VATNumber, company.CRNumber,
		company.Email, company.Phone, company.Address, company.AddressArabic, company.City, company.CityArabic,
		company.Country, company.CountryArabic, company.Logo, company.ID)
	return err
}