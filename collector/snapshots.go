package collector

import (
	"github.com/exoscale/egoscale/v3"
	"github.com/prometheus/client_golang/prometheus"
	"fmt"
	"context"
)

type SnapshotsPrometheusMetricsCollector struct {
	Context context.Context
	Client v3.Client
	Size *prometheus.Desc
	Count *prometheus.Desc
}

func NewSnapshotsPrometheusMetricsCollector(ctx context.Context, cli v3.Client) *SnapshotsPrometheusMetricsCollector {
	return &SnapshotsPrometheusMetricsCollector{
		Context: ctx,
		Client: cli,
		Count: prometheus.NewDesc(
			"exoscale_snapshot_count",
			"Amount of snapshots",
			nil, nil,
		),
		Size: prometheus.NewDesc(
			"exoscale_snapshot_size",
			"Current balance on exoscale organization",
			[]string{"id", "name"}, nil,
		),
	}
}

func (collector *SnapshotsPrometheusMetricsCollector) Describe(channel chan<- *prometheus.Desc) {
	channel <- collector.Size
}

func (collector *SnapshotsPrometheusMetricsCollector) Collect(channel chan<- prometheus.Metric) {
	response, err := collector.Client.ListSnapshots(collector.Context)

	if err != nil {
		panic(fmt.Sprintf("unable to initialize Exoscale API V3 client: %v", err))
	}

	channel <- prometheus.MustNewConstMetric(
			collector.Count,
			prometheus.CounterValue,
			0,
		)

	for _, snapshot := range response.Snapshots {
		channel <- prometheus.MustNewConstMetric(
			collector.Size,
			prometheus.GaugeValue,
			float64(snapshot.Size),
			snapshot.ID.String(), 
			snapshot.Name,
		)
	}
}