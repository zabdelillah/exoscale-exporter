package collector

import (
    "github.com/exoscale/egoscale/v3"
    "github.com/prometheus/client_golang/prometheus"
    "fmt"
    "context"
)

type TemplatesPrometheusMetricsCollector struct {
    Context context.Context
    Client v3.Client
    Volume *prometheus.Desc
    Template *prometheus.Desc
}

func NewTemplatesPrometheusMetricsCollector(ctx context.Context, cli v3.Client) *TemplatesPrometheusMetricsCollector {
    metadata := []string{"id", "name"}
    return &TemplatesPrometheusMetricsCollector{
        Context: ctx,
        Client: cli,
        Template: prometheus.NewDesc(
            "exoscale_snapshot_size",
            "Exoscale Storage Template",
            metadata, nil,
        ),
        Volume: prometheus.NewDesc(
            "exoscale_volume_size",
            "Exoscale Storage Template",
            metadata, nil,
        ),
    }
}

func (collector *TemplatesPrometheusMetricsCollector) Describe(channel chan<- *prometheus.Desc) {
    channel <- collector.Template
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