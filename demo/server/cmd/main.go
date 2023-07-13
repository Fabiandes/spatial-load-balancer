package main

import (
	"context"
	"fmt"
	"log"

	"github.com/fabiandes/slb/demo/server"
	"github.com/fabiandes/slb/demo/simulation"
	"github.com/fabiandes/slb/demo/util"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/zap"
	"google.golang.org/grpc/credentials"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// SimulationSize represents the number of entities within the simulation.
const SimulationSize = 100

// The size of the world in meters.
const (
	MapWidth  = 1000
	MapHeight = 1000
)

// Configuration for OpenTelemetry.
var (
	serviceName  = "slb-demo"
	collectorURL = "127.0.0.1:4317"
	insecure     = true
)

// initTracer creates a global OpenTelemetry Tracer and returns a cleanup function.
func initTracer(ctx context.Context) func(context.Context) error {
	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if insecure {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)

	if err != nil {
		log.Fatal(err)
	}
	resources, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Printf("Could not set resources: %v", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()), // TODO: Explore different sampling methods.
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}

// initLogger returns a zap logger configured to work with OpenTelemetry.
func initLogger() *otelzap.SugaredLogger {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	l, _ := cfg.Build()
	log := otelzap.New(l)
	return log.Sugar()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	l := initLogger()

	// Set up OpenTelemetry.
	cleanup := initTracer(ctx)
	defer cleanup(context.Background())

	// Generate a world
	w := util.GenerateWorld(SimulationSize, MapWidth, MapHeight)

	// Create a server and listen for requests.
	go func() {
		s := server.New(l, w)
		if err := s.Listen(ctx); err != nil {
			fmt.Printf("Server failed to listen: %v\n", err)
			cancel()
		}
	}()

	// Run a simulation.
	s := simulation.New(l, w)
	if err := s.Simulate(ctx); err != nil {
		fmt.Printf("An error occurred while running the simulation: %v\n", err)
		cancel()
	}
}
