package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
	"crypto/md5"
	"encoding/hex"
)

// StorageType represents the type of storage backend
type StorageType string

const (
	StorageTypeLocal   StorageType = "local"
	StorageTypeNetwork StorageType = "network"
)

// StorageConfig holds configuration for different storage backends
type StorageConfig struct {
	Type        StorageType `json:"type"`
	BasePath    string      `json:"base_path"`     // Local path or network base URL
	NetworkAuth *NetworkAuth `json:"network_auth,omitempty"` // For network storage authentication
}

// NetworkAuth holds authentication for network storage
type NetworkAuth struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
	Headers  map[string]string `json:"headers,omitempty"`
}

// StorageService provides a unified interface for file storage
type StorageService struct {
	config *StorageConfig
}

// NewStorageService creates a new storage service with the given configuration
func NewStorageService(config *StorageConfig) *StorageService {
	return &StorageService{
		config: config,
	}
}

// SaveFile saves a file and returns the relative path that should be stored in the database
func (s *StorageService) SaveFile(file multipart.File, header *multipart.FileHeader, category string) (string, error) {
	// Generate a unique filename
	filename := s.generateUniqueFilename(header.Filename)
	
	// Create category subdirectory path
	relativePath := filepath.Join(category, filename)
	
	switch s.config.Type {
	case StorageTypeLocal:
		return s.saveFileLocal(file, relativePath)
	case StorageTypeNetwork:
		return s.saveFileNetwork(file, relativePath)
	default:
		return "", fmt.Errorf("unsupported storage type: %s", s.config.Type)
	}
}

// GetFileURL returns the full URL/path to access a file
func (s *StorageService) GetFileURL(relativePath string) string {
	if relativePath == "" {
		return ""
	}
	
	switch s.config.Type {
	case StorageTypeLocal:
		return filepath.Join(s.config.BasePath, relativePath)
	case StorageTypeNetwork:
		// For network storage, construct the full URL
		baseURL := strings.TrimSuffix(s.config.BasePath, "/")
		return fmt.Sprintf("%s/%s", baseURL, strings.ReplaceAll(relativePath, "\\", "/"))
	default:
		return relativePath
	}
}

// DeleteFile removes a file from storage
func (s *StorageService) DeleteFile(relativePath string) error {
	if relativePath == "" {
		return nil
	}
	
	switch s.config.Type {
	case StorageTypeLocal:
		return s.deleteFileLocal(relativePath)
	case StorageTypeNetwork:
		return s.deleteFileNetwork(relativePath)
	default:
		return fmt.Errorf("unsupported storage type: %s", s.config.Type)
	}
}

// saveFileLocal saves a file to the local filesystem
func (s *StorageService) saveFileLocal(file multipart.File, relativePath string) (string, error) {
	// Create the full path
	fullPath := filepath.Join(s.config.BasePath, relativePath)
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}
	
	// Create the destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()
	
	// Copy the uploaded file to the destination
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}
	
	return relativePath, nil
}

// deleteFileLocal removes a file from the local filesystem
func (s *StorageService) deleteFileLocal(relativePath string) error {
	fullPath := filepath.Join(s.config.BasePath, relativePath)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %v", err)
	}
	return nil
}

// saveFileNetwork saves a file to network storage (placeholder for future implementation)
func (s *StorageService) saveFileNetwork(file multipart.File, relativePath string) (string, error) {
	// TODO: Implement network storage upload (HTTP PUT, WebDAV, etc.)
	// This would use the NetworkAuth configuration for authentication
	return "", fmt.Errorf("network storage not yet implemented")
}

// deleteFileNetwork removes a file from network storage (placeholder for future implementation)
func (s *StorageService) deleteFileNetwork(relativePath string) error {
	// TODO: Implement network storage deletion (HTTP DELETE, WebDAV, etc.)
	return fmt.Errorf("network storage not yet implemented")
}

// generateUniqueFilename creates a unique filename based on timestamp and hash
func (s *StorageService) generateUniqueFilename(originalFilename string) string {
	// Get file extension
	ext := filepath.Ext(originalFilename)
	
	// Create a hash of the original filename and current time
	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("%s_%d", originalFilename, time.Now().UnixNano())))
	hash := hex.EncodeToString(hasher.Sum(nil))
	
	// Use timestamp and hash for uniqueness
	timestamp := time.Now().Format("20060102_150405")
	
	return fmt.Sprintf("%s_%s%s", timestamp, hash[:8], ext)
}

// ValidateFileType checks if the uploaded file type is allowed
func (s *StorageService) ValidateFileType(filename string, allowedTypes []string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	
	for _, allowedType := range allowedTypes {
		if ext == allowedType {
			return nil
		}
	}
	
	return fmt.Errorf("file type %s not allowed. Allowed types: %v", ext, allowedTypes)
}

// ValidateFileSize checks if the uploaded file size is within limits
func (s *StorageService) ValidateFileSize(size int64, maxSizeMB int64) error {
	maxSizeBytes := maxSizeMB * 1024 * 1024
	if size > maxSizeBytes {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size of %d MB", size, maxSizeMB)
	}
	return nil
}