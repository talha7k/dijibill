# File Management System

This document describes the file management system that uses a separate database for storing both file metadata and file content, keeping the main database lightweight for synchronization.

## Architecture Overview

The file management system consists of two main components:

1. **File Database** (`file_database.go`) - Manages file metadata and content in a separate SQLite database
2. **File Service** (`file_service.go`) - Provides high-level file operations with compression and deduplication

## Key Features

- **Database Storage**: File content and metadata are stored directly in a separate SQLite database
- **File Deduplication**: Uses SHA256 hashing to avoid storing duplicate files
- **Image Compression**: Automatic compression for image files to save space
- **Sync Management**: Tracks sync status for remote synchronization
- **Category Organization**: Files can be categorized (logo, invoice, receipt, etc.)
- **Entity Association**: Files can be linked to specific entities (company, invoice, etc.)

## Configuration

The system is configured through the `SystemSettings` model with these fields:

```go
type SystemSettings struct {
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
fileService, err := NewFileService(settings.FileDBPath)
if err != nil {
    return err
}
defer fileService.Close()
```

### 2. Save a File

```go
// Save a company logo with automatic compression
metadata, err := fileService.SaveFileWithCompression(
    file,           // multipart.File (file content)
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

### 3. Get File Content

```go
content, metadata, err := fileService.GetFileContent(fileID)
if err != nil {
    return err
}

fmt.Printf("File: %s, Size: %d bytes\n", metadata.OriginalName, len(content))
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

## Database Schema

The file metadata database contains two tables:

### file_metadata table
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
    storage_type TEXT NOT NULL DEFAULT 'database',
    hash TEXT NOT NULL,
    sync_status TEXT DEFAULT 'pending',
    sync_error TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER
);
```

### file_content table
```sql
CREATE TABLE file_content (
    file_id INTEGER PRIMARY KEY,
    content BLOB NOT NULL,
    FOREIGN KEY (file_id) REFERENCES file_metadata(id) ON DELETE CASCADE
);
```

## Migration

The system includes automatic migration support:

1. **Migration 006**: Adds initial storage configuration to `system_settings`
2. **Migration 008**: Updates storage configuration for separate file database

## Benefits

1. **Lightweight Main Database**: File content and metadata are separate, keeping main database fast
2. **Independent Sync**: Files can be synchronized separately from main data
3. **Deduplication**: Prevents storage of duplicate files using SHA256 hashing
4. **Automatic Compression**: Images are automatically compressed to save space
5. **Database Integrity**: File content is stored with ACID guarantees
6. **Simplified Architecture**: No need for filesystem management or network storage configuration
5. **Audit Trail**: Tracks who created/modified files and when
6. **Sync Management**: Built-in support for remote synchronization

## Best Practices

1. Always use the `FileService` for file operations
2. Set appropriate file size limits using `ValidateFile`
3. Use meaningful categories for better organization
4. Regularly monitor sync status for failed uploads
5. Backup the file database separately from the main database
6. Use HTTPS for network storage to ensure security