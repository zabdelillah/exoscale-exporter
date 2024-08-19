package collector

import (
	"github.com/exoscale/egoscale/v3"
	"github.com/prometheus/client_golang/prometheus"
	"fmt"
	"context"
)

type APIKeysPrometheusMetricsCollector struct {
	Context context.Context
	Client v3.Client
	Count *prometheus.Desc
	Key *prometheus.Desc
}

func NewAPIKeysPrometheusMetricsCollector(ctx context.Context, cli v3.Client) *APIKeysPrometheusMetricsCollector {
	return &APIKeysPrometheusMetricsCollector{
		Context: ctx,
		Client: cli,
		// Count: prometheus.NewDesc(
		// 	"exoscale_iam_key_count",
		// 	"Amount of snapshots",
		// 	nil, nil,
		// ),
		Key: prometheus.NewDesc(
			"exoscale_iam_key",
			"Exoscale IAM Key",
			[]string{"key", "name", "role"}, nil,
		),
	}
}

func (collector *APIKeysPrometheusMetricsCollector) Describe(channel chan<- *prometheus.Desc) {
	channel <- collector.Key
}

func (collector *APIKeysPrometheusMetricsCollector) Collect(channel chan<- prometheus.Metric) {
	response, err := collector.Client.ListAPIKeys(collector.Context)
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "exoscale_iam_key_count",
		Help: "Exoscale IAM Key Count",
	})

	if err != nil {
		panic(fmt.Sprintf("unable to initialize Exoscale API V3 client: %v", err))
	}

	for _, key := range response.APIKeys {
		channel <- prometheus.MustNewConstMetric(
			collector.Key,
			prometheus.GaugeValue,
			0,
			key.Key,
			key.Name, 
			key.RoleID.String(),
		)
		counter.Inc()
	}

	channel <- counter
}