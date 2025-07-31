package database

import "time"

// SystemSettings operations
func (d *Database) GetSystemSettings() (*SystemSettings, error) {
	query := `SELECT id, currency, language, timezone, date_format, invoice_language, zatca_enabled, auto_backup, backup_frequency, last_backup_time, 
		COALESCE(storage_type, 'local') as storage_type,
		COALESCE(storage_local_path, './app_images') as storage_local_path,
		COALESCE(storage_s3_bucket, '') as storage_s3_bucket,
		COALESCE(storage_s3_region, '') as storage_s3_region,
		COALESCE(storage_s3_prefix, '') as storage_s3_prefix,
		created_at, updated_at FROM system_settings LIMIT 1`

	var s SystemSettings
	err := d.db.QueryRow(query).Scan(&s.ID, &s.Currency, &s.Language, &s.Timezone, &s.DateFormat, &s.InvoiceLanguage, &s.ZatcaEnabled, &s.AutoBackup, &s.BackupFrequency, &s.LastBackupTime, 
		&s.StorageType, &s.StorageLocalPath, &s.StorageS3Bucket, &s.StorageS3Region, &s.StorageS3Prefix,
		&s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (d *Database) UpdateSystemSettings(settings *SystemSettings) error {
	query := `
		UPDATE system_settings SET currency = ?, language = ?, timezone = ?, date_format = ?, invoice_language = ?, zatca_enabled = ?, auto_backup = ?, backup_frequency = ?,
		storage_type = ?, storage_local_path = ?, storage_s3_bucket = ?, storage_s3_region = ?, storage_s3_prefix = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err := d.db.Exec(query, settings.Currency, settings.Language, settings.Timezone, settings.DateFormat, settings.InvoiceLanguage, settings.ZatcaEnabled, settings.AutoBackup, settings.BackupFrequency,
		settings.StorageType, settings.StorageLocalPath, settings.StorageS3Bucket, settings.StorageS3Region, settings.StorageS3Prefix, settings.ID)
	return err
}

func (d *Database) UpdateLastBackupTime(backupTime time.Time) error {
	query := `UPDATE system_settings SET last_backup_time = ?, updated_at = CURRENT_TIMESTAMP WHERE id = 1`
	_, err := d.db.Exec(query, backupTime)
	return err
}

func (d *Database) CreateSystemSettings(settings *SystemSettings) error {
	query := `INSERT INTO system_settings (currency, language, timezone, date_format, invoice_language, zatca_enabled, auto_backup, backup_frequency, last_backup_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, err := d.db.Exec(query, settings.Currency, settings.Language, settings.Timezone, settings.DateFormat, settings.InvoiceLanguage, settings.ZatcaEnabled, settings.AutoBackup, settings.BackupFrequency, settings.LastBackupTime)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	settings.ID = int(id)
	return nil
}