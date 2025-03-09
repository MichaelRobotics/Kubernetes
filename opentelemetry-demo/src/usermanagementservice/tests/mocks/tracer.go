package mocks

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// MockTracer is a mock implementation of trace.Tracer
type MockTracer struct{}

// MockSpan is a mock implementation of trace.Span
type MockSpan struct {
	trace.Span
	name       string
	attributes []attribute.KeyValue
	recordedErrors []error
}

// Start creates a new span and context
func (t *MockTracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	span := &MockSpan{
		Span: noop.Tracer.Start(ctx, spanName),
		name: spanName,
	}
	return trace.ContextWithSpan(ctx, span), span
}

// RecordError records an error
func (s *MockSpan) RecordError(err error, opts ...trace.EventOption) {
	s.recordedErrors = append(s.recordedErrors, err)
}

// SetStatus sets the status of the span
func (s *MockSpan) SetStatus(code codes.Code, description string) {
	// This is a mock implementation, so we don't need to do anything
}

// SetAttributes sets attributes on the span
func (s *MockSpan) SetAttributes(attributes ...attribute.KeyValue) {
	s.attributes = append(s.attributes, attributes...)
}

// End completes the span
func (s *MockSpan) End(options ...trace.SpanEndOption) {
	// This is a mock implementation, so we don't need to do anything
}

// NewMockTracer creates a new mock tracer
func NewMockTracer() trace.Tracer {
	return &MockTracer{}
} 