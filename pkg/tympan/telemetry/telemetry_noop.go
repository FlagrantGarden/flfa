//go:build !telemetry
// +build !telemetry

package telemetry

import (
	"context"
)

// In noop, Start does nothing and returns the passed context and a pair of empty strings.
func Start(ctx context.Context, honeycomb_api_key string, honeycomb_dataset string, rootSpanName string, serviceName string) (context.Context, string, string) {
	// deliberately does nothing
	return ctx, "", ""
}

// In noop, EndSpan does nothing.
func EndSpan(span string) {
	// deliberately does nothing
}

// In noop, GetSpanFromContext does nothing and returns an empty string.
func GetSpanFromContext(ctx context.Context) string {
	// deliberately does nothing
	return ""
}

// In noop, NewSpan does nothing and returns the passed context and an empty string.
func NewSpan(ctx context.Context, name string) (context.Context, string) {
	// deliberately does nothing
	return ctx, ""
}

// In noop, RecordSpanError does nothing.
func RecordSpanError(span string, err error) {
	// deliberately does nothing
}

// In noop, AddStringSpanAttribute does nothing.
func AddStringSpanAttribute(span string, key string, value string) {
	// deliberately does nothing
}

// In noop, ShutDown does nothing.
func ShutDown(ctx context.Context, provider string, span string) {
	// deliberately does nothing
}
