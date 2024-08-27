package collector

import (
    "github.com/exoscale/egoscale/v3"
    "github.com/prometheus/client_golang/prometheus"
    "fmt"
    "context"
)

type LoadBalancerPrometheusMetricsCollector struct {
    Context context.Context
    Client v3.Client
    LoadBalancer *prometheus.Desc
    Service *prometheus.Desc
}

func NewLoadBalancerPrometheusMetricsCollector(ctx context.Context, cli v3.Client) *LoadBalancerPrometheusMetricsCollector {
    return &LoadBalancerPrometheusMetricsCollector{
        Context: ctx,
        Client: cli,
        LoadBalancer: prometheus.NewDesc(
            "exoscale_load_balancer",
            "Exoscale Load Balancer",
            []string{"name"}, nil,
        ),
        Service: prometheus.NewDesc(
            "exoscale_load_balancer_service",
            "Exoscale Load Balancing Service",
            []string{"name", "service", "port", "target_port", "strategy"}, nil,
        ),
    }
}

func (collector *LoadBalancerPrometheusMetricsCollector) Describe(channel chan<- *prometheus.Desc) {
    channel <- collector.LoadBalancer
    channel <- collector.Service
}

func (collector *LoadBalancerPrometheusMetricsCollector) Collect(channel chan<- prometheus.Metric) {
    loadBalancerCounter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "exoscale_load_balancer_count",
        Help: "Exoscale Load Balancer Counter",
    })
    loadBalancerServiceCounter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "exoscale_load_balancer_service_count",
        Help: "Exoscale Load Balancer Counter",
    })
    securityGroupCounter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "exoscale_security_group_count",
        Help: "Exoscale Load Balancer Counter",
    })
    privateNetworkCounter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "exoscale_private_network_count",
        Help: "Exoscale Load Balancer Counter",
    })
    elasticIpCounter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "exoscale_elastic_ip_count",
        Help: "Exoscale Load Balancer Counter",
    })

    response, err := collector.Client.ListLoadBalancers(collector.Context)
    if err != nil {
        panic(fmt.Sprintf("error initializing domain counter: %v", err))
    }

    for _, loadBalancer := range response.LoadBalancers {
        status := 1
        if loadBalancer.State != v3.LoadBalancerStateRunning {
            status = 0
        }
        channel <- prometheus.MustNewConstMetric(
            collector.LoadBalancer,
            prometheus.GaugeValue,
            float64(status),
            string(loadBalancer.Name),
        )
        for _, service := range loadBalancer.Services {
            serviceStatus := 1
            if service.State != v3.LoadBalancerServiceStateRunning {
                serviceStatus = 0
            }
            channel <- prometheus.MustNewConstMetric(
                collector.Service,
                prometheus.GaugeValue,
                float64(serviceStatus),
                string(loadBalancer.Name),
                string(service.Name),
                fmt.Sprint(service.Port),
                fmt.Sprint(service.TargetPort),
                string(service.Strategy),
            )
            loadBalancerServiceCounter.Inc()
        }
        loadBalancerCounter.Inc()
    }

    channel <- loadBalancerCounter
    channel <- loadBalancerServiceCounter
    channel <- securityGroupCounter
    channel <- privateNetworkCounter
    channel <- elasticIpCounter
}