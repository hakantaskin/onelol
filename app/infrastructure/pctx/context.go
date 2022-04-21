package pctx

import (
	"context"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ContextKey ...
type ContextKey string

const (
	ctxKeyMYSQLTx     ContextKey = "mysql-tx"
	ctxKeyLoggerEntry ContextKey = "logger-entry"
)

// WithDBTransaction returns a context with MYSQL transaction
func WithDBTransaction(ctx context.Context, value *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxKeyMYSQLTx, value)
}

// DBTransaction returns MYSQL transaction
func DBTransaction(ctx context.Context) *gorm.DB {
	value, _ := ctx.Value(ctxKeyMYSQLTx).(*gorm.DB)
	return value
}

// WithLoggerEntry returns a context with logrus entry
func WithLoggerEntry(ctx context.Context, entry *logrus.Entry) context.Context {
	return context.WithValue(ctx, ctxKeyLoggerEntry, entry)
}

// LoggerEntry gets a logrus entry from ctx, fallback to a default one if the entry is not found
func LoggerEntry(ctx context.Context) *logrus.Entry {
	entry, ok := ctx.Value(ctxKeyLoggerEntry).(*logrus.Entry)

	if !ok {
		entry = logrus.StandardLogger().WithContext(ctx)
	}

	return entry
}
