package server

import (
	"context"
	"net/http"

	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
)

func (c *Server) initTracerProvider(ctx context.Context) {
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(c.config.Services.GetJaegerCollectorEndpoint()),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		c.done <- &models.InternalError{Err: err, Msg: "failed to init jaeger tracer", Path: "user.controller.getTracerProvider"}
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("user-service"),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	c.tracerProvider = tp
}

func (s *Server) initMetrics() {
	metrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.01, 0.1, 0.5, 1, 5, 10}),
		),
	)

	reg := prometheus.WrapRegistererWith(
		prometheus.Labels{"env": *s.config.Main.Env},
		prometheus.DefaultRegisterer,
	)

	reg.MustRegister(metrics)
	s.metrics = metrics

	srv := &http.Server{Addr: s.config.Services.GetUserServicePrometheusUrl()}
	m := http.NewServeMux()
	m.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{EnableOpenMetrics: true}))

	srv.Handler = m
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err := srv.Close(); err != nil {
				msg := "failed to stop http server"
				e := &models.InternalError{Err: err, Path: "user.server.initMetrics", Msg: msg}
				s.log.ErrorStruct(msg, e)
			}
		}
	}()
}
