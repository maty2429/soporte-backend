package database

import (
	"context"
	"errors"
	"log/slog"
	"time"

	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm"
)

// slogGormLogger adapts *slog.Logger to the gorm.io/gorm/logger.Interface so
// that all GORM output (SQL traces, slow-query warnings, errors) flows through
// the application's structured logger instead of writing to stderr in its own
// text format.
type slogGormLogger struct {
	log                *slog.Logger
	level              gormlogger.LogLevel
	slowQueryThreshold time.Duration
}

// NewGormLogger returns a gormlogger.Interface backed by the given slog.Logger.
// slowQueryThreshold controls when a query is logged as WARN instead of DEBUG.
// Set it to 0 to disable slow-query logging.
func NewGormLogger(log *slog.Logger, slowQueryThreshold time.Duration) gormlogger.Interface {
	return &slogGormLogger{
		log:                log.WithGroup("gorm"),
		level:              gormlogger.Warn,
		slowQueryThreshold: slowQueryThreshold,
	}
}

// LogMode implements gormlogger.Interface.
func (l *slogGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	clone := *l
	clone.level = level
	return &clone
}

// Info implements gormlogger.Interface (GORM uses this for migration messages).
func (l *slogGormLogger) Info(ctx context.Context, msg string, args ...any) {
	if l.level >= gormlogger.Info {
		l.log.InfoContext(ctx, msg, toSlogArgs(args)...)
	}
}

// Warn implements gormlogger.Interface.
func (l *slogGormLogger) Warn(ctx context.Context, msg string, args ...any) {
	if l.level >= gormlogger.Warn {
		l.log.WarnContext(ctx, msg, toSlogArgs(args)...)
	}
}

// Error implements gormlogger.Interface.
func (l *slogGormLogger) Error(ctx context.Context, msg string, args ...any) {
	if l.level >= gormlogger.Error {
		l.log.ErrorContext(ctx, msg, toSlogArgs(args)...)
	}
}

// Trace implements gormlogger.Interface and is called for every SQL statement.
//
// Routing logic:
//   - err != nil (excluding ErrRecordNotFound) → ERROR
//   - latency > slowQueryThreshold              → WARN  (slow query)
//   - otherwise                                 → DEBUG (normal query)
func (l *slogGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.level <= gormlogger.Silent {
		return
	}

	latency := time.Since(begin)
	sql, rows := fc()

	attrs := []any{
		"latency", latency,
		"rows", rows,
		"sql", sql,
	}

	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound):
		if l.level >= gormlogger.Error {
			l.log.ErrorContext(ctx, "query error", append(attrs, "error", err)...)
		}

	case l.slowQueryThreshold > 0 && latency > l.slowQueryThreshold:
		if l.level >= gormlogger.Warn {
			l.log.WarnContext(ctx, "slow query", attrs...)
		}

	default:
		if l.level >= gormlogger.Info {
			l.log.DebugContext(ctx, "query", attrs...)
		}
	}
}

// toSlogArgs converts GORM's variadic printf-style args to slog key-value pairs.
// GORM passes args as alternating key/value pairs when using %v formatting,
// but the format string itself is the message, so we wrap unknowns safely.
func toSlogArgs(args []any) []any {
	// GORM passes formatted values; wrap them under a generic key so slog
	// doesn't panic on odd-length slices.
	if len(args)%2 != 0 {
		return append([]any{"details"}, args...)
	}
	return args
}
