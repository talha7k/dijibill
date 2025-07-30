package main

import (
	"embed"
	"fmt"
	"html/template"
)

//go:embed templates/*.html
var templateFS embed.FS

// TemplateService handles loading and managing invoice templates
type TemplateService struct {
	englishTemplate  *template.Template
	arabicTemplate   *template.Template
	bilingualTemplate *template.Template
}

// NewTemplateService creates a new template service
func NewTemplateService() (*TemplateService, error) {
	service := &TemplateService{}
	
	// Load English template
	englishContent, err := templateFS.ReadFile("templates/invoice_english.html")
	if err != nil {
		return nil, fmt.Errorf("failed to read English template: %w", err)
	}
	
	service.englishTemplate, err = template.New("invoice_english").Parse(string(englishContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse English template: %w", err)
	}
	
	// Load Arabic template
	arabicContent, err := templateFS.ReadFile("templates/invoice_arabic.html")
	if err != nil {
		return nil, fmt.Errorf("failed to read Arabic template: %w", err)
	}
	
	service.arabicTemplate, err = template.New("invoice_arabic").Parse(string(arabicContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse Arabic template: %w", err)
	}
	
	// Load Bilingual template
	bilingualContent, err := templateFS.ReadFile("templates/invoice_bilingual.html")
	if err != nil {
		return nil, fmt.Errorf("failed to read Bilingual template: %w", err)
	}
	
	service.bilingualTemplate, err = template.New("invoice_bilingual").Parse(string(bilingualContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse Bilingual template: %w", err)
	}
	
	return service, nil
}

// GetTemplate returns the appropriate template based on language
func (ts *TemplateService) GetTemplate(language string) *template.Template {
	switch language {
	case "arabic", "ar":
		return ts.arabicTemplate
	case "bilingual":
		return ts.bilingualTemplate
	default:
		return ts.englishTemplate
	}
}