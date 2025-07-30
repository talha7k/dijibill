package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"dijibill/database"
	"github.com/skip2/go-qrcode"
)

// ZATCAQRService handles ZATCA-compliant QR code generation
type ZATCAQRService struct{}

// NewZATCAQRService creates a new ZATCA QR service
func NewZATCAQRService() *ZATCAQRService {
	return &ZATCAQRService{}
}

// ZATCAQRData represents the required fields for ZATCA QR code
type ZATCAQRData struct {
	SellerName  string    // Tag 1: Seller's name
	VATNumber   string    // Tag 2: VAT registration number of the seller
	Timestamp   time.Time // Tag 3: Time stamp of the invoice (ISO 8601)
	TotalAmount float64   // Tag 4: Invoice total (with VAT)
	VATAmount   float64   // Tag 5: VAT total
	Hash        []byte    // Tag 6: Cryptographic stamp (SHA256 hash)
}

// GenerateZATCAQRCodeOnDemand generates a ZATCA-compliant QR code in Base64 PNG format on-demand
// This method doesn't store the QR code in the database, making it more efficient for large datasets
func (q *ZATCAQRService) GenerateZATCAQRCodeOnDemand(invoice *database.Invoice, company *database.Company) (string, error) {
	// Prepare QR data
	qrData := ZATCAQRData{
		SellerName:  company.Name,
		VATNumber:   company.VATNumber,
		Timestamp:   invoice.IssueDate,
		TotalAmount: invoice.TotalAmount,
		VATAmount:   invoice.VATAmount,
	}

	// Generate cryptographic stamp (simplified - in production, use proper digital signature)
	qrData.Hash = q.generateCryptographicStamp(qrData)

	// Encode to TLV format
	tlvBytes, err := q.encodeTLV(qrData)
	if err != nil {
		return "", fmt.Errorf("failed to encode TLV: %v", err)
	}

	// Encode TLV to Base64 (this is the actual QR code content)
	tlvBase64 := base64.StdEncoding.EncodeToString(tlvBytes)

	// Validate length (max 700 characters as per ZATCA specs)
	if len(tlvBase64) > 700 {
		return "", fmt.Errorf("QR code exceeds maximum length of 700 characters: %d", len(tlvBase64))
	}

	// Generate QR code image from the TLV Base64 data
	qrCodePNG, err := qrcode.Encode(tlvBase64, qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code image: %v", err)
	}

	// Convert PNG to Base64 for embedding in HTML
	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCodePNG)

	return qrCodeBase64, nil
}

// GenerateZATCAQRCode generates a ZATCA-compliant QR code in Base64 PNG format
func (q *ZATCAQRService) GenerateZATCAQRCode(invoice *database.Invoice, company *database.Company) (string, error) {
	fmt.Printf("DEBUG: Starting QR generation for invoice %d\n", invoice.ID)
	
	// Prepare QR data
	qrData := ZATCAQRData{
		SellerName:  company.Name,
		VATNumber:   company.VATNumber,
		Timestamp:   invoice.IssueDate,
		TotalAmount: invoice.TotalAmount,
		VATAmount:   invoice.VATAmount,
	}
	
	fmt.Printf("DEBUG: QR data prepared - Seller: %s, VAT: %s\n", qrData.SellerName, qrData.VATNumber)

	// Generate cryptographic stamp (simplified - in production, use proper digital signature)
	qrData.Hash = q.generateCryptographicStamp(qrData)
	fmt.Printf("DEBUG: Hash generated, length: %d\n", len(qrData.Hash))

	// Encode to TLV format
	tlvBytes, err := q.encodeTLV(qrData)
	if err != nil {
		fmt.Printf("DEBUG: TLV encoding failed: %v\n", err)
		return "", fmt.Errorf("failed to encode TLV: %v", err)
	}
	fmt.Printf("DEBUG: TLV encoded, length: %d\n", len(tlvBytes))

	// Encode TLV to Base64 (this is the actual QR code content)
	tlvBase64 := base64.StdEncoding.EncodeToString(tlvBytes)
	fmt.Printf("DEBUG: TLV Base64 encoded, length: %d\n", len(tlvBase64))

	// Validate length (max 700 characters as per ZATCA specs)
	if len(tlvBase64) > 700 {
		return "", fmt.Errorf("QR code exceeds maximum length of 700 characters: %d", len(tlvBase64))
	}

	// Generate QR code image from the TLV Base64 data
	qrCodePNG, err := qrcode.Encode(tlvBase64, qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code image: %v", err)
	}

	// Convert PNG to Base64 for embedding in HTML
	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCodePNG)
	fmt.Printf("DEBUG: QR code image generated, base64 length: %d\n", len(qrCodeBase64))

	return qrCodeBase64, nil
}

// encodeTLV encodes the QR data in Tag-Length-Value format according to ZATCA specifications
func (q *ZATCAQRService) encodeTLV(data ZATCAQRData) ([]byte, error) {
	var buffer bytes.Buffer

	// Tag 1: Seller's name
	if err := q.addTLVField(&buffer, 1, []byte(data.SellerName)); err != nil {
		return nil, err
	}

	// Tag 2: VAT registration number
	if err := q.addTLVField(&buffer, 2, []byte(data.VATNumber)); err != nil {
		return nil, err
	}

	// Tag 3: Timestamp (ISO 8601 format)
	timestamp := data.Timestamp.Format("2006-01-02T15:04:05Z")
	if err := q.addTLVField(&buffer, 3, []byte(timestamp)); err != nil {
		return nil, err
	}

	// Tag 4: Invoice total (with VAT)
	totalStr := fmt.Sprintf("%.2f", data.TotalAmount)
	if err := q.addTLVField(&buffer, 4, []byte(totalStr)); err != nil {
		return nil, err
	}

	// Tag 5: VAT total
	vatStr := fmt.Sprintf("%.2f", data.VATAmount)
	if err := q.addTLVField(&buffer, 5, []byte(vatStr)); err != nil {
		return nil, err
	}

	// Tag 6: Cryptographic stamp (hash)
	if err := q.addTLVField(&buffer, 6, data.Hash); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// addTLVField adds a single TLV field to the buffer
func (q *ZATCAQRService) addTLVField(buffer *bytes.Buffer, tag byte, value []byte) error {
	// Validate length (must fit in one byte for tags 1-5)
	if tag <= 5 && len(value) > 255 {
		return fmt.Errorf("value too long for tag %d: %d bytes (max 255)", tag, len(value))
	}

	// Tag 6 (hash) has fixed length of 32 bytes
	if tag == 6 && len(value) != 32 {
		return fmt.Errorf("hash must be exactly 32 bytes, got %d", len(value))
	}

	// Write tag (1 byte)
	buffer.WriteByte(tag)

	// Write length (1 byte)
	buffer.WriteByte(byte(len(value)))

	// Write value
	buffer.Write(value)

	return nil
}

// generateCryptographicStamp generates a cryptographic stamp (simplified implementation)
// In production, this should use proper digital signatures with ZATCA-approved certificates
func (q *ZATCAQRService) generateCryptographicStamp(data ZATCAQRData) []byte {
	// Create a string representation of the invoice data
	dataString := fmt.Sprintf("%s|%s|%s|%.2f|%.2f",
		data.SellerName,
		data.VATNumber,
		data.Timestamp.Format("2006-01-02T15:04:05Z"),
		data.TotalAmount,
		data.VATAmount,
	)

	// Generate SHA256 hash
	hash := sha256.Sum256([]byte(dataString))
	return hash[:]
}

// ValidateZATCAQRCode validates a ZATCA QR code format
func (q *ZATCAQRService) ValidateZATCAQRCode(qrCodeBase64 string) error {
	// Check length
	if len(qrCodeBase64) > 700 {
		return fmt.Errorf("QR code exceeds maximum length of 700 characters")
	}

	// Decode Base64
	tlvBytes, err := base64.StdEncoding.DecodeString(qrCodeBase64)
	if err != nil {
		return fmt.Errorf("invalid Base64 encoding: %v", err)
	}

	// Parse TLV structure
	_, err = q.parseTLV(tlvBytes)
	if err != nil {
		return fmt.Errorf("invalid TLV structure: %v", err)
	}

	return nil
}

// parseTLV parses TLV-encoded bytes back to QR data (for validation)
func (q *ZATCAQRService) parseTLV(tlvBytes []byte) (*ZATCAQRData, error) {
	data := &ZATCAQRData{}
	offset := 0

	for offset < len(tlvBytes) {
		if offset+2 > len(tlvBytes) {
			return nil, fmt.Errorf("incomplete TLV field at offset %d", offset)
		}

		tag := tlvBytes[offset]
		length := int(tlvBytes[offset+1])
		offset += 2

		if offset+length > len(tlvBytes) {
			return nil, fmt.Errorf("incomplete TLV value at offset %d", offset)
		}

		value := tlvBytes[offset : offset+length]
		offset += length

		switch tag {
		case 1:
			data.SellerName = string(value)
		case 2:
			data.VATNumber = string(value)
		case 3:
			timestamp, err := time.Parse("2006-01-02T15:04:05Z", string(value))
			if err != nil {
				return nil, fmt.Errorf("invalid timestamp format: %v", err)
			}
			data.Timestamp = timestamp
		case 4:
			if _, err := fmt.Sscanf(string(value), "%f", &data.TotalAmount); err != nil {
				return nil, fmt.Errorf("invalid total amount format: %v", err)
			}
		case 5:
			if _, err := fmt.Sscanf(string(value), "%f", &data.VATAmount); err != nil {
				return nil, fmt.Errorf("invalid VAT amount format: %v", err)
			}
		case 6:
			if length != 32 {
				return nil, fmt.Errorf("hash must be 32 bytes, got %d", length)
			}
			data.Hash = make([]byte, 32)
			copy(data.Hash, value)
		default:
			return nil, fmt.Errorf("unknown tag: %d", tag)
		}
	}

	// Validate required fields
	if data.SellerName == "" {
		return nil, fmt.Errorf("missing seller name (tag 1)")
	}
	if data.VATNumber == "" {
		return nil, fmt.Errorf("missing VAT number (tag 2)")
	}
	if data.Timestamp.IsZero() {
		return nil, fmt.Errorf("missing timestamp (tag 3)")
	}
	if data.TotalAmount == 0 {
		return nil, fmt.Errorf("missing total amount (tag 4)")
	}
	if len(data.Hash) == 0 {
		return nil, fmt.Errorf("missing cryptographic stamp (tag 6)")
	}

	return data, nil
}

// GetQRCodeInfo returns human-readable information about a QR code
func (q *ZATCAQRService) GetQRCodeInfo(qrCodeBase64 string) (map[string]interface{}, error) {
	data, err := q.parseTLV([]byte(qrCodeBase64))
	if err != nil {
		// Try to decode from Base64 first
		tlvBytes, decodeErr := base64.StdEncoding.DecodeString(qrCodeBase64)
		if decodeErr != nil {
			return nil, fmt.Errorf("failed to decode Base64: %v", decodeErr)
		}

		data, err = q.parseTLV(tlvBytes)
		if err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"seller_name":  data.SellerName,
		"vat_number":   data.VATNumber,
		"timestamp":    data.Timestamp.Format("2006-01-02T15:04:05Z"),
		"total_amount": data.TotalAmount,
		"vat_amount":   data.VATAmount,
		"hash_present": len(data.Hash) == 32,
		"compliance":   "ZATCA E-Invoicing Regulation",
	}, nil
}
