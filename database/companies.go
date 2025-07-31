package database

import (
	"fmt"
	"time"
)

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

func (d *Database) GetCompanies() ([]Company, error) {
	query := `SELECT id, name, name_arabic, vat_number, cr_number, email, phone, address, address_arabic, city, city_arabic, country, country_arabic, COALESCE(logo, '') as logo FROM companies ORDER BY name`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying companies: %v", err)
	}
	defer rows.Close()

	var companies []Company
	for rows.Next() {
		var c Company
		err := rows.Scan(&c.ID, &c.Name, &c.NameArabic, &c.VATNumber, &c.CRNumber, &c.Email, &c.Phone,
			&c.Address, &c.AddressArabic, &c.City, &c.CityArabic, &c.Country, &c.CountryArabic, &c.Logo)
		if err != nil {
			return nil, fmt.Errorf("error scanning company: %v", err)
		}
		companies = append(companies, c)
	}

	return companies, nil
}

func (d *Database) GetCompanyByID(id int) (*Company, error) {
	query := `SELECT id, name, name_arabic, vat_number, cr_number, email, phone, address, address_arabic, city, city_arabic, country, country_arabic, COALESCE(logo, '') as logo FROM companies WHERE id = ?`

	var c Company
	err := d.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.NameArabic, &c.VATNumber, &c.CRNumber, &c.Email, &c.Phone,
		&c.Address, &c.AddressArabic, &c.City, &c.CityArabic, &c.Country, &c.CountryArabic, &c.Logo)
	if err != nil {
		return nil, fmt.Errorf("error getting company: %v", err)
	}
	return &c, nil
}

func (d *Database) CreateCompany(company *Company) error {
	query := `INSERT INTO companies (name, name_arabic, vat_number, cr_number, email, phone, address, address_arabic, city, city_arabic, country, country_arabic, logo) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NULLIF(?, ''))`

	result, err := d.db.Exec(query, company.Name, company.NameArabic, company.VATNumber, company.CRNumber,
		company.Email, company.Phone, company.Address, company.AddressArabic, company.City, company.CityArabic,
		company.Country, company.CountryArabic, company.Logo)
	if err != nil {
		return fmt.Errorf("error creating company: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting company ID: %v", err)
	}

	company.ID = int(id)
	return nil
}

func (d *Database) UpdateCompany(company *Company) error {
	query := `
		UPDATE companies SET name = ?, name_arabic = ?, vat_number = ?, cr_number = ?, email = ?, phone = ?, 
		address = ?, address_arabic = ?, city = ?, city_arabic = ?, country = ?, country_arabic = ?, logo = NULLIF(?, ''), updated_at = ?
		WHERE id = ?`

	_, err := d.db.Exec(query, company.Name, company.NameArabic, company.VATNumber, company.CRNumber,
		company.Email, company.Phone, company.Address, company.AddressArabic, company.City, company.CityArabic,
		company.Country, company.CountryArabic, company.Logo, time.Now(), company.ID)
	return err
}