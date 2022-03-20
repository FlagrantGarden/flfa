//go:build telemetry
// +build telemetry

package telemetry

import (
	"context"
	"os"
	"runtime"
	"strings"

	"github.com/denisbrodbeck/machineid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/sdk/resource"
)

// Start is used to configure and initialize telemetry for an application. In the current implementation, it only
// supports sending telemetry to honeycomb.io. It:
//
// 1. Creates a new OpenTelemetry exporter, configured for Honeycomb with the passed api key and dataset.
//
// 2. Creates a new OpenTelemetry resource, setting the specified service name.
//
// 3. Creates and configures a new OpenTelemetry Tracer provider, passing the exporter and resource. This is used
// to record span information and export it.
//
// 4. It creates the root span and adds a few attributes to it, including a unique (but non-identifiable) uuid for the
// machine running the application, the operating system of the machine, the architecture of the machine, and whether or
// not the application is running in a CI context.
//
// 5. It returns the modified context, configured trace provider, and created root span.
func Start(ctx context.Context, honeycomb_api_key string, honeycomb_dataset string, rootSpanName string, serviceName string) (context.Context, *sdktrace.TracerProvider, trace.Span) {
	var tp *sdktrace.TracerProvider
	// if telemetry is turned on and honeycomb is configured, hook it up
	api_key_set := honeycomb_api_key != "not_set" && honeycomb_api_key != ""
	dataset_set := honeycomb_dataset != "not_set" && honeycomb_dataset != ""
	if api_key_set && dataset_set {
		exp, err := otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint("api.honeycomb.io:443"),
			otlptracegrpc.WithHeaders(map[string]string{
				"x-honeycomb-team":    honeycomb_api_key,
				"x-honeycomb-dataset": honeycomb_dataset,
			}),
			otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
		)
		if err != nil {
			log.Fatal().Msgf("failed to initialize exporter: %v", err)
		}

		res, err := resource.New(ctx,
			resource.WithAttributes(
				// the service name used to display traces in backends
				semconv.ServiceNameKey.String(serviceName),
			),
		)
		if err != nil {
			log.Fatal().Msgf("failed to initialize respource: %v", err)
		}

		// Create a new tracer provider with a batch span processor and the otlp exporter.
		// Add a resource attribute service.name that identifies the service in the Honeycomb UI.
		tp = sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exp),
			sdktrace.WithResource(res),
		)

		// Set the Tracer Provider and the W3C Trace Context propagator as globals
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(
			propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
		)
	} else {
		var unset_values []string
		if !api_key_set {
			unset_values = append(unset_values, "API Key")
		}
		if !dataset_set {
			unset_values = append(unset_values, "Dataset")
		}
		log.Fatal().Msgf("Unable to load honeycomb: %s must be set and not empty", strings.Join(unset_values, " and "))
	}

	tracer := otel.GetTracerProvider().Tracer("")

	var span trace.Span
	ctx, span = tracer.Start(ctx, rootSpanName)

	// The Protected ID is hashed base on application name to prevent any
	// accidental leakage of a reversable ID.
	machineUUID, _ := machineid.ProtectedID(serviceName)

	AddStringSpanAttribute(span, "uuid", machineUUID)
	AddStringSpanAttribute(span, "osinfo/os", runtime.GOOS)
	AddStringSpanAttribute(span, "osinfo/arch", runtime.GOARCH)
	runningInCi := os.Getenv("CI")
	if runningInCi == "" {
		runningInCi = "false"
	}
	AddStringSpanAttribute(span, "ci", strings.ToLower(runningInCi))

	return ctx, tp, span
}

// Close a span; this makes it immutable
func EndSpan(span trace.Span) {
	span.End()
}

// Returns the current span from context
func GetSpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// Create a new span; the span will need to be closed somewhere up the call stack
func NewSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	tracer := otel.GetTracerProvider().Tracer("")
	return tracer.Start(ctx, name)
}

// Records an error to the span; if err is nil, this function does nothing
func RecordSpanError(span trace.Span, err error) {
	span.RecordError(err)
}

// Create a new attribute and attach it to the specified span
func AddStringSpanAttribute(span trace.Span, key string, value string) {
	attr := attribute.Key(key)
	span.SetAttributes(attr.String(value))
}

// Close the root span and then the provider; this will send all events.
func ShutDown(ctx context.Context, provider *sdktrace.TracerProvider, span trace.Span) {
	// The parent span must be closed before the provider can be.
	span.End()
	err := provider.Shutdown(ctx)
	if err != nil {
		log.Fatal().Msgf("failed to shut down telemetry provider: %v", err)
	}
}
