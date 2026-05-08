package observe

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
)

var _initMeter = sync.OnceValues(func() (*metric.MeterProvider, error) {
	exporter, err := otlpmetrichttp.New(context.Background())

	if err != nil {
		return nil, err
	}
	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
	)
	otel.SetMeterProvider(mp) // 必须注册全局 MeterProvider
	return mp, nil
})

func InitMeter() (*metric.MeterProvider, error) {
	return _initMeter()
}
