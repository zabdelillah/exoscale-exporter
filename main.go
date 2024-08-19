package main
/**
 * Prometheus exporter for the Exoscale cloud provider
 */

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
	"log"
	"flag"
)

func main() {
	var port string

	flag.StringVar(&port, "listen", ":9999", "address and port to listen on")
	flag.Parse()

	ctx := context.Background()
	api_key, api_key_exists := os.LookupEnv("EXOSCALE_API_KEY")
	api_secret, api_secret_exists := os.LookupEnv("EXOSCALE_API_SECRET")

	if !api_key_exists {
		log.Printf("EXOSCALE_API_KEY not provided")
	}
	if !api_secret_exists {
		log.Printf("EXOSCALE_API_KEY not provided")
	}

	exoscaleCredentials := credentials.NewStaticCredentials(api_key, api_secret)
	exoClient, err := v3.NewClient(exoscaleCredentials)

	if err != nil {
		panic(fmt.Sprintf("unable to initialize Exoscale API V3 client: %v", err))
	}

	/**
	 * Each collector follows the same structure of containing the required functions;
	 * 	- Describe()
	 * 	- Collect()
	 * 	
	 * 	Each resource will hold a counter metric, incrementing on each one found and post-fixed
	 * 		by *_count. Every metric will be prefixed with exoscale_*, resulting in the
	 * 		naming convention of exoscale_<resource>_<metric>.
	 */
	registry := prometheus.NewRegistry()
	// Go Application Metrics
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	registry.MustRegister(prometheus.NewGoCollector())
	// Exoscale Organization-Level Metrics
	registry.MustRegister(collector.NewOrganizationPrometheusMetricsCollector(ctx, *exoClient))
	registry.MustRegister(collector.NewAPIKeysPrometheusMetricsCollector(ctx, *exoClient))
	// Exoscale Compute Metrics
	registry.MustRegister(collector.NewInstancesPrometheusMetricsCollector(ctx, *exoClient))
	registry.MustRegister(collector.NewSKSClusterPrometheusMetricsCollector(ctx, *exoClient))
	// Exoscale Storage Metrics
	registry.MustRegister(collector.NewSnapshotsPrometheusMetricsCollector(ctx, *exoClient))
	registry.MustRegister(collector.NewSOSBucketPrometheusMetricsCollector(ctx, *exoClient))
	registry.MustRegister(collector.NewTemplatesPrometheusMetricsCollector(ctx, *exoClient))
	
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	http.ListenAndServe(port, nil)
}