package collector

import (
    "github.com/exoscale/egoscale/v3"
    "github.com/prometheus/client_golang/prometheus"
    "fmt"
    "context"
)

type DNSDomainPrometheusMetricsCollector struct {
    Context context.Context
    Client v3.Client
    Domain *prometheus.Desc
    Record *prometheus.Desc
}

func NewDNSDomainPrometheusMetricsCollector(ctx context.Context, cli v3.Client) *DNSDomainPrometheusMetricsCollector {
    return &DNSDomainPrometheusMetricsCollector{
        Context: ctx,
        Client: cli,
        Domain: prometheus.NewDesc(
            "exoscale_domain",
            "Exoscale Domain",
            []string{"domain"}, nil,
        ),
        Record: prometheus.NewDesc(
            "exoscale_domain_record",
            "Exoscale Domain Record",
            []string{"domain", "name", "record", "priority", "ttl", "type"}, nil,
        ),
    }
}

func (collector *DNSDomainPrometheusMetricsCollector) Describe(channel chan<- *prometheus.Desc) {
    channel <- collector.Record
}

func (collector *DNSDomainPrometheusMetricsCollector) Collect(channel chan<- prometheus.Metric) {
    response, err := collector.Client.ListDNSDomains(collector.Context)
    counter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "exoscale_domain_count",
        Help: "Exoscale Domain Counter",
    })
    record_counter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "exoscale_domain_record_count",
        Help: "Exoscale Record Counter",
    })

    if err != nil {
        panic(fmt.Sprintf("error initializing domain counter: %v", err))
    }

    for _, domain := range response.DNSDomains {
        channel <- prometheus.MustNewConstMetric(
            collector.Domain,
            prometheus.GaugeValue,
            float64(1),
            string(domain.UnicodeName),
        )
        counter.Inc()
        record_response, err := collector.Client.ListDNSDomainRecords(collector.Context, domain.ID)
        if err != nil {
        	panic(fmt.Sprintf("error fetching dns records: %v", err))
    	}
        for _, record := range record_response.DNSDomainRecords {
        	channel <- prometheus.MustNewConstMetric(
	            collector.Record,
	            prometheus.GaugeValue,
	            float64(1),
	            string(domain.UnicodeName),
	            string(record.Name),
	            string(record.Content),
	            fmt.Sprint(record.Priority),
	            fmt.Sprint(record.Ttl),
	            string(record.Type),
        	)
        	record_counter.Inc()
        }
    }

    channel <- counter
    channel <- record_counter
}