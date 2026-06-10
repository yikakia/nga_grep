package http

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var httpMeter = sync.OnceValue(func() metric.Meter {
	return otel.Meter("api")
})

var httpRateLimitMeter = sync.OnceValues(func() (metric.Int64Counter, error) {
	return httpMeter().Int64Counter("api_rate_limit")
})

func recordRateLimit(ctx context.Context, isAllow bool) {
	if counter, err := httpRateLimitMeter(); err == nil {
		counter.Add(ctx, 1, metric.WithAttributes(attribute.Bool("allow", isAllow)))
	}
}
