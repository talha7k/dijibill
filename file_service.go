package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

// FileService manages files using only the file database (no external storage)
type FileService struct {
	fileDB          *FileDatabase
	imageCompressor *ImageCompressor
}

// NewFileService creates a new file service with only file database
func NewFileService(fileDBPath string) (*FileService, error) {
	fileDB, err := NewFileDatabase(fileDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize file database: %v", err)
	}
	
	imageCompressor := NewImageCompressor()
	
	return &FileService{
		fileDB:          fileDB,
		imageCompressor: imageCompressor,
	}, nil
}

// SaveFile saves a file directly to the file database
func (fs *FileService) SaveFile(file multipart.File, header *multipart.FileHeader, category, entityType string, entityID int, userID *int) (*FileMetadata, error) {
	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %v", err)
	}
	
	// Calculate file hash for deduplication
	hash := fs.calculateFileHashFromBytes(fileContent)
	
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
			StorageType:  "database",
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
	
	// Generate unique stored name
	storedName := fs.generateUniqueFilename(header.Filename)
	relativePath := filepath.Join(category, storedName)
	
	// Create metadata record with file content
	now := time.Now()
	metadata := &FileMetadata{
		OriginalName: header.Filename,
		StoredName:   storedName,
		RelativePath: relativePath,
		FileSize:     header.Size,
		MimeType:     header.Header.Get("Content-Type"),
		Category:     category,
		EntityType:   entityType,
		EntityID:     entityID,
		StorageType:  "database",
		Hash:         hash,
		CreatedAt:    now,
		CreatedBy:    userID,
		UpdatedAt:    now,
		UpdatedBy:    userID,
		SyncStatus:   "pending",
	}
	
	// Store file content and metadata in database
	if err := fs.fileDB.CreateFileMetadataWithContent(metadata, fileContent); err != nil {
		return nil, fmt.Errorf("failed to create file metadata with content: %v", err)
	}
	
	return metadata, nil
}

// SaveFileWithCompression saves a file with automatic image compression if applicable
func (fs *FileService) SaveFileWithCompression(file multipart.File, header *multipart.FileHeader, category, entityType string, entityID int, userID *int) (*FileMetadata, error) {
	// Check if file is an image and should be compressed
	if fs.isImageFile(header.Filename) {
		// Try to compress the image
		compressedFile, err := fs.imageCompressor.CompressImage(file, header)
		if err != nil {
			// If compression fails, fall back to original file
			if _, seekErr := file.Seek(0, 0); seekErr != nil {
				return nil, fmt.Errorf("failed to reset file pointer after compression error: %v", seekErr)
			}
			return fs.SaveFile(file, header, category, entityType, entityID, userID)
		}
		
		// Create a new file reader from compressed data
		compressedReader := bytes.NewReader(compressedFile.Data)
		compressedMultipartFile := &compressedFileWrapper{
			reader: compressedReader,
			size:   compressedFile.Size,
		}
		
		// Create new header with compressed file info
		compressedHeader := &multipart.FileHeader{
			Filename: compressedFile.Filename,
			Size:     compressedFile.Size,
			Header:   make(map[string][]string),
		}
		compressedHeader.Header["Content-Type"] = []string{compressedFile.MimeType}
		
		// Save the compressed file
		return fs.SaveFile(compressedMultipartFile, compressedHeader, category, entityType, entityID, userID)
	}
	
	// For non-image files, use regular save
	return fs.SaveFile(file, header, category, entityType, entityID, userID)
}

// compressedFileWrapper wraps compressed file data to implement multipart.File interface
type compressedFileWrapper struct {
	reader *bytes.Reader
	size   int64
}

func (cfw *compressedFileWrapper) Read(p []byte) (n int, err error) {
	return cfw.reader.Read(p)
}

func (cfw *compressedFileWrapper) ReadAt(p []byte, off int64) (n int, err error) {
	// Create a new reader from the current position
	currentPos, _ := cfw.reader.Seek(0, 1) // Get current position
	cfw.reader.Seek(off, 0)                // Seek to offset
	n, err = cfw.reader.Read(p)
	cfw.reader.Seek(currentPos, 0)         // Restore original position
	return n, err
}

func (cfw *compressedFileWrapper) Seek(offset int64, whence int) (int64, error) {
	return cfw.reader.Seek(offset, whence)
}

func (cfw *compressedFileWrapper) Close() error {
	return nil // bytes.Reader doesn't need closing
}

// isImageFile checks if the file is an image based on extension
func (fs *FileService) isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
	for _, imgExt := range imageExtensions {
		if ext == imgExt {
			return true
		}
	}
	return false
}

// GetFileContent returns the file content from the database
func (fs *FileService) GetFileContent(fileID int) ([]byte, *FileMetadata, error) {
	metadata, err := fs.fileDB.GetFileMetadataByID(fileID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get file metadata: %v", err)
	}
	
	if metadata == nil {
		return nil, nil, fmt.Errorf("file not found")
	}
	
	content, err := fs.fileDB.GetFileContent(fileID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get file content: %v", err)
	}
	
	return content, metadata, nil
}

// GetFilesByEntity returns all files for a specific entity
func (fs *FileService) GetFilesByEntity(entityType string, entityID int) ([]FileMetadata, error) {
	return fs.fileDB.GetFileMetadataByEntity(entityType, entityID)
}

// DeleteFile removes a file from the database
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
	
	// Only delete the actual file content if no other records reference it
	if otherFiles == nil || otherFiles.ID == fileID {
		if err := fs.fileDB.DeleteFileContent(fileID); err != nil {
			// Log error but don't fail the operation since metadata is already deleted
			fmt.Printf("Warning: failed to delete file content from database: %v\n", err)
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
	if err := fs.validateFileType(header.Filename, allowedTypes); err != nil {
		return err
	}
	
	// Validate file size
	if err := fs.validateFileSize(header.Size, maxSizeMB); err != nil {
		return err
	}
	
	return nil
}

// validateFileType checks if the uploaded file type is allowed
func (fs *FileService) validateFileType(filename string, allowedTypes []string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	
	for _, allowedType := range allowedTypes {
		if ext == allowedType {
			return nil
		}
	}
	
	return fmt.Errorf("file type %s not allowed. Allowed types: %v", ext, allowedTypes)
}

// validateFileSize checks if the uploaded file size is within limits
func (fs *FileService) validateFileSize(size int64, maxSizeMB int64) error {
	maxSizeBytes := maxSizeMB * 1024 * 1024
	if size > maxSizeBytes {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size of %d MB", size, maxSizeMB)
	}
	return nil
}

// GetCompressionSettings returns current image compression settings
func (fs *FileService) GetCompressionSettings() map[string]interface{} {
	return fs.imageCompressor.GetCompressionInfo()
}

// UpdateCompressionSettings updates image compression settings
func (fs *FileService) UpdateCompressionSettings(maxSizeKB int64, maxWidth, maxHeight, jpegQuality, pngQuality int) {
	fs.imageCompressor.UpdateSettings(maxSizeKB, maxWidth, maxHeight, jpegQuality, pngQuality)
}

// calculateFileHashFromBytes calculates SHA256 hash of the file content
func (fs *FileService) calculateFileHashFromBytes(content []byte) string {
	hasher := sha256.New()
	hasher.Write(content)
	return hex.EncodeToString(hasher.Sum(nil))
}

// generateUniqueFilename creates a unique filename based on timestamp and hash
func (fs *FileService) generateUniqueFilename(originalFilename string) string {
	// Get file extension
	ext := filepath.Ext(originalFilename)
	
	// Create a hash of the original filename and current time
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s_%d", originalFilename, time.Now().UnixNano())))
	hash := hex.EncodeToString(hasher.Sum(nil))
	
	// Use timestamp and hash for uniqueness
	timestamp := time.Now().Format("20060102_150405")
	
	return fmt.Sprintf("%s_%s%s", timestamp, hash[:8], ext)
}

// Close closes the file database connection
func (fs *FileService) Close() error {
	return fs.fileDB.Close()
}