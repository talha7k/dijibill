# File Management System

This document describes the new file management system that uses a separate database for storing file metadata, keeping the main database lightweight for synchronization.

## Architecture Overview

The file management system consists of three main components:

1. **Storage Service** (`storage_service.go`) - Handles physical file storage (local/network)
2. **File Database** (`file_database.go`) - Manages file metadata in a separate SQLite database
3. **File Service** (`file_service.go`) - Integrates storage and database operations

## Key Features

- **Separate Database**: File metadata is stored in a separate SQLite database to keep the main database lightweight
- **Flexible Storage**: Supports both local filesystem and network storage
- **File Deduplication**: Uses SHA256 hashing to avoid storing duplicate files
- **Sync Management**: Tracks sync status for remote synchronization
- **Category Organization**: Files can be categorized (logo, invoice, receipt, etc.)
- **Entity Association**: Files can be linked to specific entities (company, invoice, etc.)

## Configuration

The system is configured through the `SystemSettings` model with these new fields:

```go
type SystemSettings struct {
    // Storage configuration
    StorageType     string `json:"storage_type"`     // "local" or "network"
    StorageBasePath string `json:"storage_base_path"` // Base path for file storage
    
    // File database configuration
    FileDBPath         string `json:"file_db_path"`          // Path to separate file database
    FileDBSyncURL      string `json:"file_db_sync_url"`      // URL for remote sync
    FileDBSyncToken    string `json:"file_db_sync_token"`    // Authentication token
    FileDBAutoSync     bool   `json:"file_db_auto_sync"`     // Enable automatic sync
    FileDBSyncInterval int    `json:"file_db_sync_interval"` // Sync interval in seconds
}
```

## Usage Examples

### 1. Initialize File Service

```go
// Get system settings
settings, err := mainDB.GetSystemSettings()
if err != nil {
    return err
}

// Initialize file service
fileService, err := database.NewFileService(
    settings.FileDBPath,      // Path to separate file database
    settings.StorageType,     // Storage type (local/network)
    settings.StorageBasePath, // Base path for file storage
)
if err != nil {
    return err
}
defer fileService.Close()
```

### 2. Save a File

```go
// Save a company logo
metadata, err := fileService.SaveFile(
    file,           // io.Reader (file content)
    fileHeader,     // *multipart.FileHeader (file info)
    "logo",         // category
    "company",      // entity type
    1,              // entity ID
    userID,         // user ID
)
if err != nil {
    return err
}

fmt.Printf("File saved: %s\n", metadata.StoredName)
```

### 3. Get File URL

```go
url, err := fileService.GetFileURL(metadata.ID)
if err != nil {
    return err
}

fmt.Printf("File URL: %s\n", url)
```

### 4. Get Files by Entity

```go
// Get all files for a company
files, err := fileService.GetFilesByEntity("company", 1)
if err != nil {
    return err
}

for _, file := range files {
    fmt.Printf("File: %s (%s)\n", file.OriginalName, file.Category)
}
```

### 5. Delete a File

```go
err := fileService.DeleteFile(metadata.ID, userID)
if err != nil {
    return err
}
```

### 6. Sync Operations

```go
// Get files pending sync
pendingFiles, err := fileService.GetPendingSyncFiles()
if err != nil {
    return err
}

// Mark file as synced
err = fileService.MarkFileSynced(fileID)
if err != nil {
    return err
}

// Mark file sync as failed
err = fileService.MarkFileSyncFailed(fileID, "Connection timeout")
if err != nil {
    return err
}
```

## File Categories

The system supports various file categories:

- `logo` - Company logos
- `invoice` - Invoice attachments
- `receipt` - Receipt images
- `product` - Product images
- `document` - General documents
- `backup` - Backup files

## Storage Types

### Local Storage
- Files are stored in the local filesystem
- Path configured via `StorageBasePath`
- Files organized in subdirectories by category and date

### Network Storage
- Files are stored on a remote server
- Base URL configured via `StorageBasePath`
- Supports HTTP/HTTPS protocols

## Database Schema

The file metadata database contains a single table:

```sql
CREATE TABLE file_metadata (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    original_name TEXT NOT NULL,
    stored_name TEXT NOT NULL UNIQUE,
    relative_path TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    mime_type TEXT NOT NULL,
    category TEXT NOT NULL,
    entity_type TEXT,
    entity_id INTEGER,
    storage_type TEXT NOT NULL,
    hash TEXT NOT NULL,
    sync_status TEXT DEFAULT 'pending',
    sync_error TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER
);
```

## Migration

The system includes automatic migration support:

1. **Migration 006**: Adds initial storage configuration to `system_settings`
2. **Migration 008**: Updates storage configuration for separate file database

## Benefits

1. **Lightweight Main Database**: File metadata is separate, keeping main database fast
2. **Independent Sync**: Files can be synchronized separately from main data
3. **Deduplication**: Prevents storage of duplicate files
4. **Flexible Storage**: Easy to switch between local and network storage
5. **Audit Trail**: Tracks who created/modified files and when
6. **Sync Management**: Built-in support for remote synchronization

## Best Practices

1. Always use the `FileService` for file operations
2. Set appropriate file size limits using `ValidateFile`
3. Use meaningful categories for better organization
4. Regularly monitor sync status for failed uploads
5. Backup the file database separately from the main database
6. Use HTTPS for network storage to ensure security