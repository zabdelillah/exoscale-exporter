package collector

import (
    "github.com/exoscale/egoscale/v3"
    "github.com/prometheus/client_golang/prometheus"
    "fmt"
    "context"
)

type SKSClusterPrometheusMetricsCollector struct {
    Context context.Context
    Client v3.Client
    Cluster *prometheus.Desc
    NodePool *prometheus.Desc
    NodePoolSize *prometheus.Desc
    NodePoolDiskSize *prometheus.Desc
}

func NewSKSClusterPrometheusMetricsCollector(ctx context.Context, cli v3.Client) *SKSClusterPrometheusMetricsCollector {
    metadata := []string{"id", "name", "level", "version"}
    nodePoolMetadata := []string{"id", "name", "version"}
    return &SKSClusterPrometheusMetricsCollector{
        Context: ctx,
        Client: cli,
        Cluster: prometheus.NewDesc(
            "exoscale_sks_cluster_up",
            "Exoscale SKS Cluster",
            metadata, nil,
        ),
        NodePool: prometheus.NewDesc(
            "exoscale_sks_nodepool_up",
            "Exoscale SKS NodePool",
            nodePoolMetadata, nil,
        ),
        NodePoolSize: prometheus.NewDesc(
            "exoscale_sks_nodepool_size",
            "Exoscale SKS NodePool",
            nodePoolMetadata, nil,
        ),
        NodePoolDiskSize: prometheus.NewDesc(
            "exoscale_sks_nodepool_disk_size",
            "Exoscale SKS NodePool",
            nodePoolMetadata, nil,
        ),
    }
}

func (collector *SKSClusterPrometheusMetricsCollector) Describe(channel chan<- *prometheus.Desc) {
    channel <- collector.Cluster
    channel <- collector.NodePool
    channel <- collector.NodePoolSize
    channel <- collector.NodePoolDiskSize
}

func (collector *SKSClusterPrometheusMetricsCollector) Collect(channel chan<- prometheus.Metric) {
    response, err := collector.Client.ListSKSClusters(collector.Context)
    counter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "exoscale_sks_cluster_count",
        Help: "Exoscale Buckets Count",
    })

    if err != nil {
        panic(fmt.Sprintf("unable to initialize Exoscale API V3 client: %v", err))
    }

    for _, cluster := range response.SKSClusters {
        clusterNodeSize := prometheus.NewCounter(prometheus.CounterOpts{
            Name: "exoscale_sks_cluster_size",
            Help: "Exoscale Buckets Count",
        })

        metadata := []string{
            cluster.ID.String(),
            cluster.Name,
            string(cluster.Level),
            cluster.Version}
        // Push cluster up status
        state := 0
        if (cluster.State == v3.SKSClusterStateRunning) {
            state = 1
        }
        channel <- prometheus.MustNewConstMetric(
            collector.Cluster,
            prometheus.CounterValue,
            float64(state),
            metadata...,
        )

        for _, pool := range cluster.Nodepools {
            nodePoolMetadata := []string{
                pool.ID.String(),
                pool.Name,
                pool.Version}
            state = 0
            if (pool.State == v3.SKSNodepoolStateRunning) {
                state = 1
            }
            channel <- prometheus.MustNewConstMetric(
                collector.NodePool,
                prometheus.CounterValue,
                float64(state),
                nodePoolMetadata...,
            )

            channel <- prometheus.MustNewConstMetric(
                collector.NodePoolSize,
                prometheus.CounterValue,
                float64(pool.Size),
                nodePoolMetadata...,
            )

            channel <- prometheus.MustNewConstMetric(
                collector.NodePoolDiskSize,
                prometheus.CounterValue,
                float64(pool.DiskSize),
                nodePoolMetadata...,
            )

            clusterNodeSize.Add(float64(pool.Size))
        }
        counter.Inc()
        channel <- clusterNodeSize
    }

    channel <- counter
}