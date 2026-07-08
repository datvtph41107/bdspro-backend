package seo_domain

import "strings"

type SeoQualityInput struct {
	HasH1             bool
	HasInternalLinks  bool
	HasStructuredData bool
	CanonicalUnique   bool
}

type SeoQualityResult struct {
	Score    uint32
	Warnings []string
}

func CalculateSeoQuality(page *SeoDomain, input SeoQualityInput) SeoQualityResult {
	if page == nil {
		return SeoQualityResult{Score: 0, Warnings: []string{"page_nil"}}
	}

	var score uint32
	warnings := make([]string, 0, 10)

	if strings.TrimSpace(page.Title) != "" {
		score += 10
	} else {
		warnings = append(warnings, "missing_title")
	}

	if strings.TrimSpace(page.Description) != "" {
		score += 10
	} else {
		warnings = append(warnings, "missing_description")
	}

	if strings.TrimSpace(page.CanonicalURL) != "" {
		score += 10
	} else {
		warnings = append(warnings, "missing_canonical")
	}

	if input.CanonicalUnique {
		score += 10
	} else {
		warnings = append(warnings, "canonical_not_verified_unique")
	}

	contentLen := len(strings.TrimSpace(page.Content))
	switch {
	case contentLen >= 1200:
		score += 20
	case contentLen >= 500:
		score += 12
		warnings = append(warnings, "content_medium")
	case contentLen > 0:
		score += 5
		warnings = append(warnings, "content_thin")
	default:
		warnings = append(warnings, "missing_content")
	}

	if input.HasH1 || strings.TrimSpace(page.Title) != "" {
		score += 10
	} else {
		warnings = append(warnings, "missing_h1")
	}

	if input.HasInternalLinks || len(page.InternalLinks) > 0 {
		score += 10
	} else {
		warnings = append(warnings, "missing_internal_links")
	}

	if page.RefSource == SeoRefSourceManual || strings.TrimSpace(string(page.RefSnapshotJSON)) != "" {
		score += 10
	} else {
		warnings = append(warnings, "missing_source_snapshot")
	}

	if input.HasStructuredData || strings.Contains(strings.ToLower(string(page.Metadata)), "schema") {
		score += 10
	} else {
		warnings = append(warnings, "missing_structured_data")
	}

	if score > 100 {
		score = 100
	}

	return SeoQualityResult{Score: score, Warnings: warnings}
}
