package database

// DefaultProductSettings operations
func (d *Database) GetDefaultProductSettings() (*DefaultProductSettings, error) {
	query := `SELECT id, default_stock, default_tax_rate_id, default_unit_id, 
			  default_product_type, default_product_status, default_markup, 
			  default_price_includes_tax, default_price_change_allowed, created_at, updated_at 
			  FROM default_product_settings LIMIT 1`

	var dps DefaultProductSettings
	err := d.db.QueryRow(query).Scan(&dps.ID, &dps.DefaultStock, &dps.DefaultTaxRateID, &dps.DefaultUnitID,
		&dps.DefaultProductType, &dps.DefaultProductStatus,
		&dps.DefaultMarkup, &dps.DefaultPriceIncludesTax, &dps.DefaultPriceChangeAllowed, &dps.CreatedAt, &dps.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &dps, nil
}

func (d *Database) UpdateDefaultProductSettings(settings *DefaultProductSettings) error {
	query := `UPDATE default_product_settings SET 
			  default_stock = ?, default_tax_rate_id = ?, default_unit_id = ?, 
			  default_product_type = ?, default_product_status = ?, 
			  default_markup = ?, default_price_includes_tax = ?, default_price_change_allowed = ?, 
			  updated_at = CURRENT_TIMESTAMP 
			  WHERE id = ?`

	_, err := d.db.Exec(query, settings.DefaultStock, settings.DefaultTaxRateID, settings.DefaultUnitID,
		settings.DefaultProductType,
		settings.DefaultProductStatus, settings.DefaultMarkup, settings.DefaultPriceIncludesTax,
		settings.DefaultPriceChangeAllowed, settings.ID)
	return err
}