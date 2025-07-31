package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
)

// ImageCompressor handles automatic image compression
type ImageCompressor struct {
	MaxSizeKB     int64   // Maximum file size in KB (default: 150)
	MaxWidth      int     // Maximum width in pixels (default: 1920)
	MaxHeight     int     // Maximum height in pixels (default: 1080)
	JPEGQuality   int     // JPEG quality 1-100 (default: 85)
	PNGQuality    int     // PNG compression level 0-9 (default: 6)
	AllowedTypes  []string // Allowed image types
}

// NewImageCompressor creates a new image compressor with default settings
func NewImageCompressor() *ImageCompressor {
	return &ImageCompressor{
		MaxSizeKB:    150,
		MaxWidth:     1920,
		MaxHeight:    1080,
		JPEGQuality:  85,
		PNGQuality:   6,
		AllowedTypes: []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"},
	}
}

// CompressedFile represents a compressed image file
type CompressedFile struct {
	Data     []byte
	Filename string
	Size     int64
	MimeType string
}

// CompressImage compresses an uploaded image file
func (ic *ImageCompressor) CompressImage(file multipart.File, header *multipart.FileHeader) (*CompressedFile, error) {
	// Check if file is an image
	if !ic.isImageFile(header.Filename) {
		return nil, fmt.Errorf("file is not a supported image type")
	}

	// Read the original file
	originalData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Check if file is already small enough
	if int64(len(originalData)) <= ic.MaxSizeKB*1024 {
		return &CompressedFile{
			Data:     originalData,
			Filename: header.Filename,
			Size:     int64(len(originalData)),
			MimeType: header.Header.Get("Content-Type"),
		}, nil
	}

	// Decode the image
	img, format, err := image.Decode(bytes.NewReader(originalData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	// Resize if necessary
	resizedImg := ic.resizeImage(img)

	// Compress the image
	compressedData, newFormat, err := ic.compressImage(resizedImg, format, header.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to compress image: %v", err)
	}

	// Generate new filename with correct extension
	newFilename := ic.generateCompressedFilename(header.Filename, newFormat)

	// Determine MIME type
	mimeType := ic.getMimeType(newFormat)

	return &CompressedFile{
		Data:     compressedData,
		Filename: newFilename,
		Size:     int64(len(compressedData)),
		MimeType: mimeType,
	}, nil
}

// isImageFile checks if the file is a supported image type
func (ic *ImageCompressor) isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedType := range ic.AllowedTypes {
		if ext == allowedType {
			return true
		}
	}
	return false
}

// resizeImage resizes an image if it exceeds maximum dimensions
func (ic *ImageCompressor) resizeImage(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Check if resizing is needed
	if width <= ic.MaxWidth && height <= ic.MaxHeight {
		return img
	}

	// Calculate new dimensions maintaining aspect ratio
	aspectRatio := float64(width) / float64(height)
	var newWidth, newHeight int

	if width > height {
		newWidth = ic.MaxWidth
		newHeight = int(float64(ic.MaxWidth) / aspectRatio)
		if newHeight > ic.MaxHeight {
			newHeight = ic.MaxHeight
			newWidth = int(float64(ic.MaxHeight) * aspectRatio)
		}
	} else {
		newHeight = ic.MaxHeight
		newWidth = int(float64(ic.MaxHeight) * aspectRatio)
		if newWidth > ic.MaxWidth {
			newWidth = ic.MaxWidth
			newHeight = int(float64(ic.MaxWidth) / aspectRatio)
		}
	}

	// Create new image with resized dimensions
	resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	
	// Simple nearest neighbor scaling
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := int(float64(x) * float64(width) / float64(newWidth))
			srcY := int(float64(y) * float64(height) / float64(newHeight))
			resized.Set(x, y, img.At(srcX, srcY))
		}
	}

	return resized
}

// compressImage compresses an image with progressive quality reduction
func (ic *ImageCompressor) compressImage(img image.Image, originalFormat, filename string) ([]byte, string, error) {
	// Try JPEG compression first (usually smaller)
	if data, err := ic.compressAsJPEG(img, ic.JPEGQuality); err == nil {
		if int64(len(data)) <= ic.MaxSizeKB*1024 {
			return data, "jpeg", nil
		}
	}

	// Try progressive quality reduction for JPEG
	for quality := ic.JPEGQuality; quality >= 20; quality -= 10 {
		if data, err := ic.compressAsJPEG(img, quality); err == nil {
			if int64(len(data)) <= ic.MaxSizeKB*1024 {
				return data, "jpeg", nil
			}
		}
	}

	// Try PNG compression if original was PNG
	if originalFormat == "png" {
		if data, err := ic.compressAsPNG(img); err == nil {
			if int64(len(data)) <= ic.MaxSizeKB*1024 {
				return data, "png", nil
			}
		}
	}

	// If still too large, try more aggressive JPEG compression
	for quality := 15; quality >= 5; quality -= 5 {
		if data, err := ic.compressAsJPEG(img, quality); err == nil {
			if int64(len(data)) <= ic.MaxSizeKB*1024 {
				return data, "jpeg", nil
			}
		}
	}

	// Last resort: very low quality JPEG
	data, err := ic.compressAsJPEG(img, 1)
	if err != nil {
		return nil, "", fmt.Errorf("failed to compress image to target size: %v", err)
	}

	return data, "jpeg", nil
}

// compressAsJPEG compresses an image as JPEG with specified quality
func (ic *ImageCompressor) compressAsJPEG(img image.Image, quality int) ([]byte, error) {
	var buf bytes.Buffer
	
	options := &jpeg.Options{
		Quality: quality,
	}
	
	if err := jpeg.Encode(&buf, img, options); err != nil {
		return nil, err
	}
	
	return buf.Bytes(), nil
}

// compressAsPNG compresses an image as PNG
func (ic *ImageCompressor) compressAsPNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	
	encoder := &png.Encoder{
		CompressionLevel: png.BestCompression,
	}
	
	if err := encoder.Encode(&buf, img); err != nil {
		return nil, err
	}
	
	return buf.Bytes(), nil
}

// generateCompressedFilename generates a new filename with correct extension
func (ic *ImageCompressor) generateCompressedFilename(originalFilename, format string) string {
	ext := filepath.Ext(originalFilename)
	nameWithoutExt := strings.TrimSuffix(originalFilename, ext)
	
	switch format {
	case "jpeg":
		return nameWithoutExt + ".jpg"
	case "png":
		return nameWithoutExt + ".png"
	default:
		return nameWithoutExt + ".jpg"
	}
}

// getMimeType returns the MIME type for a given format
func (ic *ImageCompressor) getMimeType(format string) string {
	switch format {
	case "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	default:
		return "image/jpeg"
	}
}

// GetCompressionInfo returns information about the compression settings
func (ic *ImageCompressor) GetCompressionInfo() map[string]interface{} {
	return map[string]interface{}{
		"max_size_kb":    ic.MaxSizeKB,
		"max_width":      ic.MaxWidth,
		"max_height":     ic.MaxHeight,
		"jpeg_quality":   ic.JPEGQuality,
		"png_quality":    ic.PNGQuality,
		"allowed_types":  ic.AllowedTypes,
	}
}

// UpdateSettings allows updating compression settings
func (ic *ImageCompressor) UpdateSettings(maxSizeKB int64, maxWidth, maxHeight, jpegQuality, pngQuality int) {
	if maxSizeKB > 0 {
		ic.MaxSizeKB = maxSizeKB
	}
	if maxWidth > 0 {
		ic.MaxWidth = maxWidth
	}
	if maxHeight > 0 {
		ic.MaxHeight = maxHeight
	}
	if jpegQuality >= 1 && jpegQuality <= 100 {
		ic.JPEGQuality = jpegQuality
	}
	if pngQuality >= 0 && pngQuality <= 9 {
		ic.PNGQuality = pngQuality
	}
}