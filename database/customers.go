package database

// Customer operations
func (d *Database) CreateCustomer(customer *Customer) error {
	query := `
		INSERT INTO customers (name, name_arabic, vat_number, email, phone, address, address_arabic, city, city_arabic, country, country_arabic)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := d.db.Exec(query, customer.Name, customer.NameArabic, customer.VATNumber,
		customer.Email, customer.Phone, customer.Address, customer.AddressArabic,
		customer.City, customer.CityArabic, customer.Country, customer.CountryArabic)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	customer.ID = int(id)
	return nil
}

func (d *Database) GetCustomers() ([]Customer, error) {
	query := `SELECT id, name, name_arabic, vat_number, email, phone, address, address_arabic, 
			  city, city_arabic, country, country_arabic, created_at, updated_at FROM customers ORDER BY name`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var c Customer
		err := rows.Scan(&c.ID, &c.Name, &c.NameArabic, &c.VATNumber, &c.Email, &c.Phone,
			&c.Address, &c.AddressArabic, &c.City, &c.CityArabic, &c.Country, &c.CountryArabic,
			&c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}

func (d *Database) GetCustomerByID(id int) (*Customer, error) {
	query := `SELECT id, name, name_arabic, vat_number, email, phone, address, address_arabic, 
			  city, city_arabic, country, country_arabic, created_at, updated_at FROM customers WHERE id = ?`

	var c Customer
	err := d.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.NameArabic, &c.VATNumber, &c.Email, &c.Phone,
		&c.Address, &c.AddressArabic, &c.City, &c.CityArabic, &c.Country, &c.CountryArabic,
		&c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (d *Database) UpdateCustomer(customer *Customer) error {
	query := `
		UPDATE customers SET name = ?, name_arabic = ?, vat_number = ?, email = ?, phone = ?, 
		address = ?, address_arabic = ?, city = ?, city_arabic = ?, country = ?, country_arabic = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err := d.db.Exec(query, customer.Name, customer.NameArabic, customer.VATNumber,
		customer.Email, customer.Phone, customer.Address, customer.AddressArabic,
		customer.City, customer.CityArabic, customer.Country, customer.CountryArabic, customer.ID)
	return err
}

func (d *Database) DeleteCustomer(id int) error {
	_, err := d.db.Exec("DELETE FROM customers WHERE id = ?", id)
	return err
}