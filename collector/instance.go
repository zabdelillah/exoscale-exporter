package collector

import (
    "github.com/exoscale/egoscale/v3"
    "github.com/prometheus/client_golang/prometheus"
    "fmt"
    "context"
)

type InstancesPrometheusMetricsCollector struct {
    Context context.Context
    Client v3.Client
    Instance *prometheus.Desc
    InstancePool *prometheus.Desc
    InstancePoolSize *prometheus.Desc
    CPUs *prometheus.Desc
    GPUs *prometheus.Desc
    Memory *prometheus.Desc
}

func NewInstancesPrometheusMetricsCollector(ctx context.Context, cli v3.Client) *InstancesPrometheusMetricsCollector {
    metadata := []string{"id", "name", "family", "size", "zone"}
    return &InstancesPrometheusMetricsCollector{
        Context: ctx,
        Client: cli,
        Instance: prometheus.NewDesc(
            "exoscale_instance_up",
            "Exoscale Instance",
            metadata, nil,
        ),
        CPUs: prometheus.NewDesc(
            "exoscale_instance_cpus",
            "Exoscale Instance",
            metadata, nil,
        ),
        GPUs: prometheus.NewDesc(
            "exoscale_instance_gpus",
            "Exoscale Instance",
            metadata, nil,
        ),
        Memory: prometheus.NewDesc(
            "exoscale_instance_memory",
            "Exoscale Instance",
            metadata, nil,
        ),
        InstancePool: prometheus.NewDesc(
            "exoscale_instance_pool_up",
            "Exoscale Instance",
            []string{"id", "name"}, nil,
        ),
        InstancePoolSize: prometheus.NewDesc(
            "exoscale_instance_pool_size",
            "Exoscale Instance",
            []string{"id", "name"}, nil,
        ),
    }
}

func (collector *InstancesPrometheusMetricsCollector) Describe(channel chan<- *prometheus.Desc) {
    channel <- collector.Instance
    channel <- collector.CPUs
    channel <- collector.GPUs
    channel <- collector.Memory
    channel <- collector.InstancePool
    channel <- collector.InstancePoolSize
}

func (collector *InstancesPrometheusMetricsCollector) Collect(channel chan<- prometheus.Metric) {
    response, err := collector.Client.ListInstances(collector.Context)
    instanceCounter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "exoscale_instances_count",
        Help: "Exoscale Instance Count",
    })

    if err != nil {
        panic(fmt.Sprintf("unable to initialize Exoscale API V3 client: %v", err))
    }

    for _, instance := range response.Instances {
        instanceType, _ := collector.Client.GetInstanceType(collector.Context, instance.InstanceType.ID)
        // for _, zone := range instance.InstanceType.Zones {
            zone := "exoscale" // Unknown how to get instance zones for the moment
            metadata := []string{
                instance.ID.String(),
                instance.Name,
                string(instanceType.Family),
                string(instanceType.Size),
                string(zone),
            }
            var state int = 0;
            if (instance.State == v3.InstanceStateRunning) {
                state = 1
            }

            channel <- prometheus.MustNewConstMetric(
                collector.Instance,
                prometheus.CounterValue,
                float64(state),
                metadata...,
            )

            channel <- prometheus.MustNewConstMetric(
                collector.CPUs,
                prometheus.CounterValue,
                float64(instanceType.Cpus),
                metadata...,
            )

            channel <- prometheus.MustNewConstMetric(
                collector.GPUs,
                prometheus.CounterValue,
                float64(instanceType.Gpus),
                metadata...,
            )

            channel <- prometheus.MustNewConstMetric(
                collector.Memory,
                prometheus.CounterValue,
                float64(instanceType.Memory),
                metadata...,
            )

            instanceCounter.Inc()
        // }
    }

    // Count Instance Pools
    instancePoolsResponse, err := collector.Client.ListInstancePools(collector.Context)
    instancePoolCounter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "exoscale_instance_pool_count",
        Help: "Exoscale Instance Pool Count",
    })

    if err != nil {
        panic(fmt.Sprintf("unable to initialize Exoscale API V3 client: %v", err))
    }

    for _, pool := range instancePoolsResponse.InstancePools {
        var state int = 1;
        if (pool.State == v3.InstancePoolStateSuspended) {
            state = 0
        }
        metadata := []string{pool.ID.String(), pool.Name}
        channel <- prometheus.MustNewConstMetric(
                collector.InstancePool,
                prometheus.CounterValue,
                float64(state),
                metadata...
            )
        channel <- prometheus.MustNewConstMetric(
                collector.InstancePoolSize,
                prometheus.CounterValue,
                float64(state),
                metadata...
            )
        instancePoolCounter.Inc()
    }

    channel <- instanceCounter
    channel <- instancePoolCounter
}