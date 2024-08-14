package collector

import (
    "github.com/exoscale/egoscale/v3"
    "github.com/prometheus/client_golang/prometheus"
    "fmt"
    "context"
)

type SOSBucketPrometheusMetricsCollector struct {
    Context context.Context
    Client v3.Client
    Bucket *prometheus.Desc
}

func NewSOSBucketPrometheusMetricsCollector(ctx context.Context, cli v3.Client) *SOSBucketPrometheusMetricsCollector {
    return &SOSBucketPrometheusMetricsCollector{
        Context: ctx,
        Client: cli,
        Bucket: prometheus.NewDesc(
            "exoscale_sos_bucket",
            "Exoscale SOS Bucket",
            []string{"name", "zone"}, nil,
        ),
    }
}

func (collector *SOSBucketPrometheusMetricsCollector) Describe(channel chan<- *prometheus.Desc) {
    channel <- collector.Bucket
}

func (collector *SOSBucketPrometheusMetricsCollector) Collect(channel chan<- prometheus.Metric) {
    response, err := collector.Client.ListSOSBucketsUsage(collector.Context)
    counter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "exoscale_sos_bucket_count",
        Help: "Exoscale Buckets Count",
    })

    if err != nil {
        panic(fmt.Sprintf("unable to initialize Exoscale API V3 client: %v", err))
    }

    for _, bucket := range response.SOSBucketsUsage {
        channel <- prometheus.MustNewConstMetric(
            collector.Bucket,
            prometheus.GaugeValue,
            float64(bucket.Size),
            bucket.Name,
            string(bucket.ZoneName),
        )
        counter.Inc()
    }

    channel <- counter
}