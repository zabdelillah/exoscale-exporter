package main

import (
	"zai.dev/m/v2/collector"
	"github.com/exoscale/egoscale/v3"
	"github.com/exoscale/egoscale/v3/credentials"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"context"
	"fmt"
	"os"
)

func main() {
	// Initialize Credentials
	ctx := context.Background()
	// Login to Exoscale
	exoscaleCredentials := credentials.NewStaticCredentials(
		os.Getenv("EXOSCALE_API_KEY"), os.Getenv("EXOSCALE_API_SECRET"),
	)
	exoClient, err := v3.NewClient(exoscaleCredentials)
	if err != nil {
		panic(fmt.Sprintf("unable to initialize Exoscale API V3 client: %v", err))
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	registry.MustRegister(prometheus.NewGoCollector())
	// Organization Metrics
	registry.MustRegister(collector.NewOrganizationPrometheusMetricsCollector(ctx, *exoClient))
	registry.MustRegister(collector.NewAPIKeysPrometheusMetricsCollector(ctx, *exoClient))
	// Instance Metrics
	registry.MustRegister(collector.NewInstancesPrometheusMetricsCollector(ctx, *exoClient))
	registry.MustRegister(collector.NewSnapshotsPrometheusMetricsCollector(ctx, *exoClient))
	registry.MustRegister(collector.NewSKSClusterPrometheusMetricsCollector(ctx, *exoClient))
	// Storage Metrics
	registry.MustRegister(collector.NewSOSBucketPrometheusMetricsCollector(ctx, *exoClient))
	registry.MustRegister(collector.NewTemplatesPrometheusMetricsCollector(ctx, *exoClient))
	
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	http.ListenAndServe(":9999", nil)
}