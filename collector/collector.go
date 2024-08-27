package collector
/**
 * Prometheus exporter for the Exoscale cloud provider
 */

import (
	"github.com/exoscale/egoscale/v3"
	// "github.com/exoscale/egoscale/v3/credentials"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"context"
	// "fmt"
	// "os"
	// "log"
	// "flag"
)

func PrepareCollector(ctx context.Context, cli *v3.Client) {
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
	registry.MustRegister(NewOrganizationPrometheusMetricsCollector(ctx, *cli))
	registry.MustRegister(NewAPIKeysPrometheusMetricsCollector(ctx, *cli))
	// Exoscale Compute Metrics
	registry.MustRegister(NewInstancesPrometheusMetricsCollector(ctx, *cli))
	registry.MustRegister(NewSKSClusterPrometheusMetricsCollector(ctx, *cli))
	// Exoscale Storage Metrics
	registry.MustRegister(NewSnapshotsPrometheusMetricsCollector(ctx, *cli))
	registry.MustRegister(NewSOSBucketPrometheusMetricsCollector(ctx, *cli))
	registry.MustRegister(NewTemplatesPrometheusMetricsCollector(ctx, *cli))
	// Exoscale DNS Metrics
	registry.MustRegister(NewDNSDomainPrometheusMetricsCollector(ctx, *cli))
	// Exoscale Networking Metrics
	registry.MustRegister(NewLoadBalancerPrometheusMetricsCollector(ctx, *cli))
	registry.MustRegister(NewElasticIPPrometheusMetricsCollector(ctx, *cli))
	registry.MustRegister(NewPrivateNetworkPrometheusMetricsCollector(ctx, *cli))
	registry.MustRegister(NewSecurityGroupPrometheusMetricsCollector(ctx, *cli))
	
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
}