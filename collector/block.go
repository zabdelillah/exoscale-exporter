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

type TemplatesPrometheusMetricsCollector struct {
    Context context.Context
    Client v3.Client
    Volume *prometheus.Desc
    Template *prometheus.Desc
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

func NewTemplatesPrometheusMetricsCollector(ctx context.Context, cli v3.Client) *TemplatesPrometheusMetricsCollector {
    metadata := []string{"id", "name"}
    return &TemplatesPrometheusMetricsCollector{
        Context: ctx,
        Client: cli,
        Volume: prometheus.NewDesc(
            "exoscale_volume_size",
            "Exoscale Storage Template",
            metadata, nil,
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


func (collector *TemplatesPrometheusMetricsCollector) Describe(channel chan<- *prometheus.Desc) {
    channel <- collector.Volume
}

func (collector *TemplatesPrometheusMetricsCollector) Collect(channel chan<- prometheus.Metric) {
    response, err := collector.Client.ListBlockStorageVolumes(collector.Context)
    counter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "exoscale_volume_count",
        Help: "Exoscale Templates Counter",
    })

    if err != nil {
        panic(fmt.Sprintf("unable to initialize Exoscale API V3 client: %v", err))
    }

    for _, volume := range response.BlockStorageVolumes {
        metadata := []string{volume.ID.String(), volume.Name}
        // for _, zone := range template.Zones {
            channel <- prometheus.MustNewConstMetric(
                collector.Volume,
                prometheus.GaugeValue,
                float64(volume.Size),
                metadata...
            )
            counter.Inc()
        // }
    }

    channel <- counter
}