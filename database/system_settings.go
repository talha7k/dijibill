package database

import "time"

// SystemSettings operations
func (d *Database) GetSystemSettings() (*SystemSettings, error) {
	query := `SELECT id, currency, language, timezone, date_format, invoice_language, zatca_enabled, auto_backup, backup_frequency, last_backup_time, 
		COALESCE(storage_type, 'local') as storage_type,
		COALESCE(storage_base_path, './app_images') as storage_base_path,
		COALESCE(file_db_path, './files.db') as file_db_path,
		COALESCE(file_db_sync_url, '') as file_db_sync_url,
		COALESCE(file_db_sync_token, '') as file_db_sync_token,
		COALESCE(file_db_auto_sync, 0) as file_db_auto_sync,
		COALESCE(file_db_sync_interval, 300) as file_db_sync_interval,
		created_at, updated_at FROM system_settings LIMIT 1`

	var s SystemSettings
	err := d.db.QueryRow(query).Scan(&s.ID, &s.Currency, &s.Language, &s.Timezone, &s.DateFormat, &s.InvoiceLanguage, &s.ZatcaEnabled, &s.AutoBackup, &s.BackupFrequency, &s.LastBackupTime, 
		&s.StorageType, &s.StorageBasePath, &s.FileDBPath, &s.FileDBSyncURL, &s.FileDBSyncToken, &s.FileDBAutoSync, &s.FileDBSyncInterval,
		&s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (d *Database) UpdateSystemSettings(settings *SystemSettings) error {
	query := `
		UPDATE system_settings SET currency = ?, language = ?, timezone = ?, date_format = ?, invoice_language = ?, zatca_enabled = ?, auto_backup = ?, backup_frequency = ?,
		storage_type = ?, storage_base_path = ?, file_db_path = ?, file_db_sync_url = ?, file_db_sync_token = ?, file_db_auto_sync = ?, file_db_sync_interval = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err := d.db.Exec(query, settings.Currency, settings.Language, settings.Timezone, settings.DateFormat, settings.InvoiceLanguage, settings.ZatcaEnabled, settings.AutoBackup, settings.BackupFrequency,
		settings.StorageType, settings.StorageBasePath, settings.FileDBPath, settings.FileDBSyncURL, settings.FileDBSyncToken, settings.FileDBAutoSync, settings.FileDBSyncInterval, settings.ID)
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