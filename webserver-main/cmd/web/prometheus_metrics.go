package main

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Structure that store Metric server dependencies and functions
// including HTTP serving
type Metrics struct {
	Registry *prometheus.Registry
	Address  string
	//DB           *sql.DB
	HTTPRequests *prometheus.CounterVec
	HTTPDuration *prometheus.HistogramVec
}

// NewCustomMetrics used to initialize Metrics structure as
// a model for exposing application metrics
// It will initialize new custom registry and fill it
// with metrics
// The only entry point for using with main server server functionality
// It is be a part of application dependencies, to provide access to sql.DB
func NewCustomMetrics() *Metrics {
	// Initialize custom registry
	reg := prometheus.NewRegistry()

	m := &Metrics{
		Registry: reg,
		HTTPRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_request_total",
				Help: "Total numbers of HTTP request processed.",
			},
			[]string{"path", "method", "status"},
		),
		HTTPDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Histogram of response latencies for HTTP requests.",
				Buckets: prometheus.DefBuckets,
			}, []string{"path", "method"},
		),
	}

	// Register custom metrics
	reg.MustRegister(m.HTTPRequests)
	reg.MustRegister(m.HTTPDuration)

	// Register DB metrics
	//reg.MustRegister(collectors.NewDBStatsCollector(m.DB, "debut"))

	// Register default metrics about Go runtime and compilation
	reg.MustRegister(collectors.NewBuildInfoCollector())
	reg.MustRegister(collectors.NewGoCollector())
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	return m
}

// Helper structure that will be used to trap
// response and it's status code
type responseWriterTrap struct {
	http.ResponseWriter
	statusCode int
}

// Helper function used to capture the response response code
// to use it as a metric label in the future
func (rwt *responseWriterTrap) WriterHeader(code int) {
	rwt.statusCode = code
	rwt.ResponseWriter.WriteHeader(code)
}

// Middleware used as wrapper for HTTP handlers.
// Captures HTTP requests and their duration
func (m *Metrics) MetricsHTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		writerTrap := &responseWriterTrap{w, http.StatusOK}
		next.ServeHTTP(writerTrap, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(writerTrap.statusCode)

		m.HTTPRequests.WithLabelValues(r.URL.Path, r.Method, status).Inc()
		m.HTTPDuration.WithLabelValues(r.URL.Path, r.Method).Observe(duration)
	})
}

// Used to serve HTTP server with single /metrics path
// with prometheus metric.
// By default server will be opened on address :2112
// with path /metrics
// Server will be launched inside dedicated goroutine
// Require initialized Metrics structure
// Accepts address/port on which server will be available
func (m *Metrics) MetricsListenAndServe() {
	// If in Metrics structure address is not
	// explicitelly specified, we will set it to :2112
	if m.Address == "" {
		m.Address = ":2112"
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	metricsMux := http.NewServeMux()

	// Be careful with dereferenced values
	metricsMux.Handle("/metrics", promhttp.HandlerFor(
		m.Registry,
		promhttp.HandlerOpts{EnableOpenMetrics: true},
	))

	metricsServer := &http.Server{
		Addr:    m.Address,
		Handler: m.MetricsHTTPMiddleware(metricsMux),
	}

	go func() {
		logger.Info("Metrics handling is up!", "address", metricsServer.Addr)
		err := metricsServer.ListenAndServe()
		logger.Error(err.Error())
	}()
}
