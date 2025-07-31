package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// ensureDir creates a directory if it doesn't exist
func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// FileDatabase handles all file-related database operations separately
type FileDatabase struct {
	db *sql.DB
}

// FileMetadata represents file information in the separate file database
type FileMetadata struct {
	ID           int       `json:"id" db:"id"`
	OriginalName string    `json:"original_name" db:"original_name"`
	StoredName   string    `json:"stored_name" db:"stored_name"`
	RelativePath string    `json:"relative_path" db:"relative_path"`
	FileSize     int64     `json:"file_size" db:"file_size"`
	MimeType     string    `json:"mime_type" db:"mime_type"`
	Category     string    `json:"category" db:"category"` // e.g., "company_logos", "product_images"
	EntityType   string    `json:"entity_type" db:"entity_type"` // e.g., "company", "product"
	EntityID     int       `json:"entity_id" db:"entity_id"`
	StorageType  string    `json:"storage_type" db:"storage_type"` // "local" or "network"
	Hash         string    `json:"hash" db:"hash"` // File hash for deduplication
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	CreatedBy    *int      `json:"created_by" db:"created_by"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy    *int      `json:"updated_by" db:"updated_by"`
	SyncStatus   string    `json:"sync_status" db:"sync_status"` // "pending", "synced", "failed"
	LastSyncAt   *time.Time `json:"last_sync_at" db:"last_sync_at"`
}

// NewFileDatabase creates a new file database connection
func NewFileDatabase(dbPath string) (*FileDatabase, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := ensureDir(dir); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, err
	}

	fileDB := &FileDatabase{db: db}
	
	// Create tables
	if err := fileDB.createTables(); err != nil {
		return nil, err
	}

	return fileDB, nil
}

// createTables creates the file metadata tables
func (fd *FileDatabase) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS file_metadata (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			original_name TEXT NOT NULL,
			stored_name TEXT NOT NULL,
			relative_path TEXT NOT NULL UNIQUE,
			file_size INTEGER NOT NULL,
			mime_type TEXT NOT NULL,
			category TEXT NOT NULL,
			entity_type TEXT NOT NULL,
			entity_id INTEGER NOT NULL,
			storage_type TEXT NOT NULL DEFAULT 'local',
			hash TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_by INTEGER,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_by INTEGER,
			sync_status TEXT NOT NULL DEFAULT 'pending',
			last_sync_at DATETIME
		)`,
		`CREATE INDEX IF NOT EXISTS idx_file_metadata_entity ON file_metadata(entity_type, entity_id)`,
		`CREATE INDEX IF NOT EXISTS idx_file_metadata_category ON file_metadata(category)`,
		`CREATE INDEX IF NOT EXISTS idx_file_metadata_storage_type ON file_metadata(storage_type)`,
		`CREATE INDEX IF NOT EXISTS idx_file_metadata_sync_status ON file_metadata(sync_status)`,
		`CREATE INDEX IF NOT EXISTS idx_file_metadata_hash ON file_metadata(hash)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_file_metadata_path ON file_metadata(relative_path)`,
	}

	for _, query := range queries {
		if _, err := fd.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %s, error: %v", query, err)
		}
	}

	return nil
}

// CreateFileMetadata inserts a new file metadata record
func (fd *FileDatabase) CreateFileMetadata(metadata *FileMetadata) error {
	query := `
		INSERT INTO file_metadata (
			original_name, stored_name, relative_path, file_size, mime_type,
			category, entity_type, entity_id, storage_type, hash,
			created_at, created_by, updated_at, updated_by, sync_status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	result, err := fd.db.Exec(query,
		metadata.OriginalName, metadata.StoredName, metadata.RelativePath,
		metadata.FileSize, metadata.MimeType, metadata.Category,
		metadata.EntityType, metadata.EntityID, metadata.StorageType, metadata.Hash,
		metadata.CreatedAt, metadata.CreatedBy, metadata.UpdatedAt, metadata.UpdatedBy,
		metadata.SyncStatus,
	)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	metadata.ID = int(id)
	return nil
}

// GetFileMetadataByEntity retrieves file metadata for a specific entity
func (fd *FileDatabase) GetFileMetadataByEntity(entityType string, entityID int) ([]FileMetadata, error) {
	query := `
		SELECT id, original_name, stored_name, relative_path, file_size, mime_type,
			   category, entity_type, entity_id, storage_type, hash,
			   created_at, created_by, updated_at, updated_by, sync_status, last_sync_at
		FROM file_metadata
		WHERE entity_type = ? AND entity_id = ?
		ORDER BY created_at DESC
	`
	
	rows, err := fd.db.Query(query, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var files []FileMetadata
	for rows.Next() {
		var file FileMetadata
		err := rows.Scan(
			&file.ID, &file.OriginalName, &file.StoredName, &file.RelativePath,
			&file.FileSize, &file.MimeType, &file.Category, &file.EntityType,
			&file.EntityID, &file.StorageType, &file.Hash, &file.CreatedAt,
			&file.CreatedBy, &file.UpdatedAt, &file.UpdatedBy, &file.SyncStatus,
			&file.LastSyncAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	
	return files, nil
}

// GetFileMetadataByID retrieves a specific file metadata record
func (fd *FileDatabase) GetFileMetadataByID(id int) (*FileMetadata, error) {
	query := `
		SELECT id, original_name, stored_name, relative_path, file_size, mime_type,
			   category, entity_type, entity_id, storage_type, hash,
			   created_at, created_by, updated_at, updated_by, sync_status, last_sync_at
		FROM file_metadata
		WHERE id = ?
	`
	
	var file FileMetadata
	err := fd.db.QueryRow(query, id).Scan(
		&file.ID, &file.OriginalName, &file.StoredName, &file.RelativePath,
		&file.FileSize, &file.MimeType, &file.Category, &file.EntityType,
		&file.EntityID, &file.StorageType, &file.Hash, &file.CreatedAt,
		&file.CreatedBy, &file.UpdatedAt, &file.UpdatedBy, &file.SyncStatus,
		&file.LastSyncAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return &file, nil
}

// GetPendingSyncFiles retrieves files that need to be synced
func (fd *FileDatabase) GetPendingSyncFiles() ([]FileMetadata, error) {
	query := `
		SELECT id, original_name, stored_name, relative_path, file_size, mime_type,
			   category, entity_type, entity_id, storage_type, hash,
			   created_at, created_by, updated_at, updated_by, sync_status, last_sync_at
		FROM file_metadata
		WHERE sync_status = 'pending' OR sync_status = 'failed'
		ORDER BY created_at ASC
	`
	
	rows, err := fd.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var files []FileMetadata
	for rows.Next() {
		var file FileMetadata
		err := rows.Scan(
			&file.ID, &file.OriginalName, &file.StoredName, &file.RelativePath,
			&file.FileSize, &file.MimeType, &file.Category, &file.EntityType,
			&file.EntityID, &file.StorageType, &file.Hash, &file.CreatedAt,
			&file.CreatedBy, &file.UpdatedAt, &file.UpdatedBy, &file.SyncStatus,
			&file.LastSyncAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	
	return files, nil
}

// UpdateSyncStatus updates the sync status of a file
func (fd *FileDatabase) UpdateSyncStatus(id int, status string) error {
	now := time.Now()
	query := `UPDATE file_metadata SET sync_status = ?, last_sync_at = ? WHERE id = ?`
	_, err := fd.db.Exec(query, status, now, id)
	return err
}

// DeleteFileMetadata removes a file metadata record
func (fd *FileDatabase) DeleteFileMetadata(id int) error {
	query := `DELETE FROM file_metadata WHERE id = ?`
	_, err := fd.db.Exec(query, id)
	return err
}

// GetFileByHash checks if a file with the same hash already exists (deduplication)
func (fd *FileDatabase) GetFileByHash(hash string) (*FileMetadata, error) {
	query := `
		SELECT id, original_name, stored_name, relative_path, file_size, mime_type,
			   category, entity_type, entity_id, storage_type, hash,
			   created_at, created_by, updated_at, updated_by, sync_status, last_sync_at
		FROM file_metadata
		WHERE hash = ?
		LIMIT 1
	`
	
	var file FileMetadata
	err := fd.db.QueryRow(query, hash).Scan(
		&file.ID, &file.OriginalName, &file.StoredName, &file.RelativePath,
		&file.FileSize, &file.MimeType, &file.Category, &file.EntityType,
		&file.EntityID, &file.StorageType, &file.Hash, &file.CreatedAt,
		&file.CreatedBy, &file.UpdatedAt, &file.UpdatedBy, &file.SyncStatus,
		&file.LastSyncAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return &file, nil
}

// Close closes the file database connection
func (fd *FileDatabase) Close() error {
	return fd.db.Close()
}

// GetDatabaseStats returns statistics about the file database
func (fd *FileDatabase) GetDatabaseStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Total files
	var totalFiles int
	err := fd.db.QueryRow("SELECT COUNT(*) FROM file_metadata").Scan(&totalFiles)
	if err != nil {
		return nil, err
	}
	stats["total_files"] = totalFiles
	
	// Total size
	var totalSize int64
	err = fd.db.QueryRow("SELECT COALESCE(SUM(file_size), 0) FROM file_metadata").Scan(&totalSize)
	if err != nil {
		return nil, err
	}
	stats["total_size_bytes"] = totalSize
	
	// Pending sync count
	var pendingSync int
	err = fd.db.QueryRow("SELECT COUNT(*) FROM file_metadata WHERE sync_status = 'pending'").Scan(&pendingSync)
	if err != nil {
		return nil, err
	}
	stats["pending_sync"] = pendingSync
	
	// Files by storage type
	rows, err := fd.db.Query("SELECT storage_type, COUNT(*) FROM file_metadata GROUP BY storage_type")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	storageStats := make(map[string]int)
	for rows.Next() {
		var storageType string
		var count int
		if err := rows.Scan(&storageType, &count); err != nil {
			return nil, err
		}
		storageStats[storageType] = count
	}
	stats["by_storage_type"] = storageStats
	
	return stats, nil
}