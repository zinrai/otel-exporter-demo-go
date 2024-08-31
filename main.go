package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("sample-app"),
			semconv.ServiceVersion("0.1.0"),
		),
	)
	if err != nil {
		return nil, err
	}

	exporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithEndpoint("localhost:4318"),
	)
	if err != nil {
		return nil, err
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(1*time.Second))),
	)
	otel.SetMeterProvider(provider)

	return provider.Shutdown, nil
}

func main() {
	shutdown, err := initProvider()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	meter := otel.Meter("sample-app")

	counter, err := meter.Int64Counter(
		"sample.counter",
		metric.WithDescription("A sample counter"),
	)
	if err != nil {
		log.Fatal(err)
	}

	gauge, err := meter.Float64ObservableGauge(
		"sample.gauge",
		metric.WithDescription("A sample gauge"),
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = meter.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		o.ObserveFloat64(gauge, rand.Float64()*100)
		return nil
	}, gauge)
	if err != nil {
		log.Fatal(err)
	}

	for {
		counter.Add(context.Background(), 1, metric.WithAttributes(attribute.String("label", "value")))
		time.Sleep(1 * time.Second)
	}
}
