package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"
)

// FileService integrates storage service with file database for complete file management
type FileService struct {
	storageService *StorageService
	fileDB         *FileDatabase
}

// NewFileService creates a new file service with storage and database components
func NewFileService(storageConfig *StorageConfig, fileDBPath string) (*FileService, error) {
	storageService := NewStorageService(storageConfig)
	
	fileDB, err := NewFileDatabase(fileDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize file database: %v", err)
	}
	
	return &FileService{
		storageService: storageService,
		fileDB:         fileDB,
	}, nil
}

// SaveFile saves a file to storage and records metadata in the file database
func (fs *FileService) SaveFile(file multipart.File, header *multipart.FileHeader, category, entityType string, entityID int, userID *int) (*FileMetadata, error) {
	// Calculate file hash for deduplication
	hash, err := fs.calculateFileHash(file)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate file hash: %v", err)
	}
	
	// Reset file pointer after hash calculation
	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to reset file pointer: %v", err)
	}
	
	// Check if file already exists (deduplication)
	existingFile, err := fs.fileDB.GetFileByHash(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing file: %v", err)
	}
	
	if existingFile != nil {
		// File already exists, create a new metadata record pointing to the same file
		now := time.Now()
		metadata := &FileMetadata{
			OriginalName: header.Filename,
			StoredName:   existingFile.StoredName,
			RelativePath: existingFile.RelativePath,
			FileSize:     header.Size,
			MimeType:     header.Header.Get("Content-Type"),
			Category:     category,
			EntityType:   entityType,
			EntityID:     entityID,
			StorageType:  fs.storageService.config.Type.String(),
			Hash:         hash,
			CreatedAt:    now,
			CreatedBy:    userID,
			UpdatedAt:    now,
			UpdatedBy:    userID,
			SyncStatus:   "pending",
		}
		
		if err := fs.fileDB.CreateFileMetadata(metadata); err != nil {
			return nil, fmt.Errorf("failed to create file metadata: %v", err)
		}
		
		return metadata, nil
	}
	
	// Save new file to storage
	relativePath, err := fs.storageService.SaveFile(file, header, category)
	if err != nil {
		return nil, fmt.Errorf("failed to save file to storage: %v", err)
	}
	
	// Create metadata record
	now := time.Now()
	metadata := &FileMetadata{
		OriginalName: header.Filename,
		StoredName:   filepath.Base(relativePath),
		RelativePath: relativePath,
		FileSize:     header.Size,
		MimeType:     header.Header.Get("Content-Type"),
		Category:     category,
		EntityType:   entityType,
		EntityID:     entityID,
		StorageType:  fs.storageService.config.Type.String(),
		Hash:         hash,
		CreatedAt:    now,
		CreatedBy:    userID,
		UpdatedAt:    now,
		UpdatedBy:    userID,
		SyncStatus:   "pending",
	}
	
	if err := fs.fileDB.CreateFileMetadata(metadata); err != nil {
		// If database insert fails, try to clean up the stored file
		fs.storageService.DeleteFile(relativePath)
		return nil, fmt.Errorf("failed to create file metadata: %v", err)
	}
	
	return metadata, nil
}

// GetFileURL returns the URL/path to access a file
func (fs *FileService) GetFileURL(fileID int) (string, error) {
	metadata, err := fs.fileDB.GetFileMetadataByID(fileID)
	if err != nil {
		return "", fmt.Errorf("failed to get file metadata: %v", err)
	}
	
	if metadata == nil {
		return "", fmt.Errorf("file not found")
	}
	
	return fs.storageService.GetFileURL(metadata.RelativePath), nil
}

// GetFilesByEntity returns all files for a specific entity
func (fs *FileService) GetFilesByEntity(entityType string, entityID int) ([]FileMetadata, error) {
	return fs.fileDB.GetFileMetadataByEntity(entityType, entityID)
}

// DeleteFile removes a file from both storage and database
func (fs *FileService) DeleteFile(fileID int) error {
	metadata, err := fs.fileDB.GetFileMetadataByID(fileID)
	if err != nil {
		return fmt.Errorf("failed to get file metadata: %v", err)
	}
	
	if metadata == nil {
		return fmt.Errorf("file not found")
	}
	
	// Check if other records reference the same file (deduplication)
	otherFiles, err := fs.fileDB.GetFileByHash(metadata.Hash)
	if err != nil {
		return fmt.Errorf("failed to check for other file references: %v", err)
	}
	
	// Delete metadata record
	if err := fs.fileDB.DeleteFileMetadata(fileID); err != nil {
		return fmt.Errorf("failed to delete file metadata: %v", err)
	}
	
	// Only delete the actual file if no other records reference it
	if otherFiles == nil || otherFiles.ID == fileID {
		if err := fs.storageService.DeleteFile(metadata.RelativePath); err != nil {
			// Log error but don't fail the operation since metadata is already deleted
			fmt.Printf("Warning: failed to delete file from storage: %v\n", err)
		}
	}
	
	return nil
}

// GetPendingSyncFiles returns files that need to be synced
func (fs *FileService) GetPendingSyncFiles() ([]FileMetadata, error) {
	return fs.fileDB.GetPendingSyncFiles()
}

// MarkFileSynced updates the sync status of a file
func (fs *FileService) MarkFileSynced(fileID int) error {
	return fs.fileDB.UpdateSyncStatus(fileID, "synced")
}

// MarkFileSyncFailed updates the sync status of a file to failed
func (fs *FileService) MarkFileSyncFailed(fileID int) error {
	return fs.fileDB.UpdateSyncStatus(fileID, "failed")
}

// GetDatabaseStats returns statistics about the file database
func (fs *FileService) GetDatabaseStats() (map[string]interface{}, error) {
	return fs.fileDB.GetDatabaseStats()
}

// ValidateFile validates file type and size
func (fs *FileService) ValidateFile(header *multipart.FileHeader, allowedTypes []string, maxSizeMB int64) error {
	// Validate file type
	if err := fs.storageService.ValidateFileType(header.Filename, allowedTypes); err != nil {
		return err
	}
	
	// Validate file size
	if err := fs.storageService.ValidateFileSize(header.Size, maxSizeMB); err != nil {
		return err
	}
	
	return nil
}

// calculateFileHash calculates SHA256 hash of the file content
func (fs *FileService) calculateFileHash(file multipart.File) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// Close closes the file database connection
func (fs *FileService) Close() error {
	return fs.fileDB.Close()
}

// String method for StorageType to convert to string
func (st StorageType) String() string {
	return string(st)
}