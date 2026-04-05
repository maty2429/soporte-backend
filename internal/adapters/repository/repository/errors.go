package repository

import (
	"context"
	"database/sql/driver"
	"errors"
	"strings"

	"gorm.io/gorm"

	"soporte/internal/core/domain"
)

// wrapDBError converts database errors to domain errors without exposing
// internal database structure. This prevents leaking table names, column
// names, and other schema details to API responses.
func wrapDBError(resource string, err error) error {
	if err == nil {
		return nil
	}

	if isContextError(err) {
		return domain.ServiceUnavailableError("database query timed out", err)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.NotFoundError(resource, err)
	}

	if isUniqueViolation(err) {
		return domain.ConflictError(resource+" already exists", err)
	}

	if isForeignKeyViolation(err) {
		return domain.ValidationError("invalid reference", err)
	}

	if isConnectionError(err) {
		return domain.ServiceUnavailableError("database is not available", err)
	}

	return domain.InternalError("database operation failed", err)
}

// wrapListError wraps errors from list operations.
func wrapListError(resource string, err error) error {
	if err == nil {
		return nil
	}

	if isContextError(err) {
		return domain.ServiceUnavailableError("database query timed out", err)
	}

	if isConnectionError(err) {
		return domain.ServiceUnavailableError("database is not available", err)
	}

	return domain.InternalError("failed to list "+resource, err)
}

// wrapCreateError wraps errors from create operations.
func wrapCreateError(resource string, err error) error {
	if err == nil {
		return nil
	}

	if isContextError(err) {
		return domain.ServiceUnavailableError("database query timed out", err)
	}

	if isUniqueViolation(err) {
		return domain.ConflictError(resource+" already exists", err)
	}

	if isForeignKeyViolation(err) {
		return domain.ValidationError("invalid reference", err)
	}

	if isConnectionError(err) {
		return domain.ServiceUnavailableError("database is not available", err)
	}

	return domain.InternalError("failed to create "+resource, err)
}

// wrapUpdateError wraps errors from update operations.
func wrapUpdateError(resource string, err error) error {
	if err == nil {
		return nil
	}

	if isContextError(err) {
		return domain.ServiceUnavailableError("database query timed out", err)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.NotFoundError(resource, err)
	}

	if isUniqueViolation(err) {
		return domain.ConflictError(resource+" already exists", err)
	}

	if isForeignKeyViolation(err) {
		return domain.ValidationError("invalid reference", err)
	}

	if isConnectionError(err) {
		return domain.ServiceUnavailableError("database is not available", err)
	}

	return domain.InternalError("failed to update "+resource, err)
}

// isUniqueViolation checks if the error is a unique constraint violation.
// Supports PostgreSQL and SQLite error messages.
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())

	// PostgreSQL: "duplicate key value violates unique constraint"
	if strings.Contains(msg, "duplicate key") {
		return true
	}

	// SQLite: "UNIQUE constraint failed"
	if strings.Contains(msg, "unique constraint") {
		return true
	}

	return false
}

// isForeignKeyViolation checks if the error is a foreign key violation.
// Supports PostgreSQL and SQLite error messages.
func isForeignKeyViolation(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())

	// PostgreSQL: "violates foreign key constraint"
	if strings.Contains(msg, "foreign key") {
		return true
	}

	// SQLite: "FOREIGN KEY constraint failed"
	if strings.Contains(msg, "foreign key constraint") {
		return true
	}

	return false
}

// wrapDeleteError wraps errors from delete operations.
func wrapDeleteError(resource string, err error) error {
	if err == nil {
		return nil
	}

	if isContextError(err) {
		return domain.ServiceUnavailableError("database query timed out", err)
	}

	if isForeignKeyViolation(err) {
		return domain.ConflictError(resource+" has dependent records", err)
	}

	if isConnectionError(err) {
		return domain.ServiceUnavailableError("database is not available", err)
	}

	return domain.InternalError("failed to delete "+resource, err)
}

// isContextError reports whether the error is a context cancellation or deadline.
// Both mean the query was abandoned — return 503 so the client knows to retry.
func isContextError(err error) bool {
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)
}

// derefBool safely dereferences a *bool, returning false if nil.
func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// derefStr safely dereferences a *string, returning "" if nil.
func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// isConnectionError reports whether the error is a database connectivity failure.
// Covers: connection refused, dropped connections, bad connections from the pool,
// and broken pipes — all of which indicate the DB is unreachable at runtime.
func isConnectionError(err error) bool {
	if errors.Is(err, driver.ErrBadConn) {
		return true
	}

	msg := strings.ToLower(err.Error())

	return strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "unexpected eof") ||
		strings.Contains(msg, "bad connection") ||
		strings.Contains(msg, "no such host") ||
		strings.Contains(msg, "dial tcp")
}
