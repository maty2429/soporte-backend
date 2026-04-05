package services

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"gorm.io/gorm"

	"soporte/internal/core/domain"
)

// isDuplicateKeyError reports whether err is a unique-constraint violation.
func isDuplicateKeyError(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	lower := strings.ToLower(err.Error())
	return strings.Contains(lower, "duplicate key") || strings.Contains(lower, "unique constraint")
}

// splitRutDV parses a raw Chilean RUT string (e.g. "12.345.678-9" or "123456789") into its
// numeric body and check digit (DV). Returns a ValidationError on bad input.
func splitRutDV(raw string) (string, string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(raw))
	normalized = strings.ReplaceAll(normalized, ".", "")
	normalized = strings.ReplaceAll(normalized, "-", "")
	normalized = strings.ReplaceAll(normalized, " ", "")

	if len(normalized) < 2 {
		return "", "", domain.ValidationError("rut must include number and dv", nil)
	}

	runes := []rune(normalized)
	rut := string(runes[:len(runes)-1])
	dv := string(runes[len(runes)-1])

	for _, r := range rut {
		if !unicode.IsDigit(r) {
			return "", "", domain.ValidationError("rut must contain only digits before dv", nil)
		}
	}

	if !unicode.IsDigit(runes[len(runes)-1]) && runes[len(runes)-1] != 'K' {
		return "", "", domain.ValidationError(fmt.Sprintf("dv '%s' is invalid", dv), nil)
	}

	return rut, dv, nil
}

// wrapServiceError passes through domain errors unchanged and wraps anything else
// as an InternalError with the given operation name.
func wrapServiceError(op string, err error) error {
	var appErr *domain.Error
	if errors.As(err, &appErr) {
		return err
	}
	return domain.InternalError(op, err)
}

// ptrOrDefaultBool returns the dereferenced value of p if non-nil, otherwise def.
func ptrOrDefaultBool(p *bool, def bool) bool {
	if p == nil {
		return def
	}
	return *p
}

// normalizePagination clamps limit and offset to sane values.
// If limit is 0 or negative it falls back to DefaultListLimit.
// If limit exceeds MaxListLimit it is clamped to MaxListLimit.
// Negative offsets are clamped to 0.
func normalizePagination(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = DefaultListLimit
	}
	if limit > MaxListLimit {
		limit = MaxListLimit
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
