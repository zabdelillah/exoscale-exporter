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

type SecurityGroupPrometheusMetricsCollector struct {
    Context context.Context
    Client v3.Client
    SecurityGroup *prometheus.Desc
    SecurityGroupRule *prometheus.Desc
}

type NamedPrometheusMetricsCollector struct {
    Context context.Context
    Client v3.Client
    Object *prometheus.Desc
}

type PrivateNetworkPrometheusMetricsCollector NamedPrometheusMetricsCollector
type ElasticIPPrometheusMetricsCollector NamedPrometheusMetricsCollector

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

func NewSecurityGroupPrometheusMetricsCollector(ctx context.Context, cli v3.Client) *SecurityGroupPrometheusMetricsCollector {
    return &SecurityGroupPrometheusMetricsCollector{
        Context: ctx,
        Client: cli,
        SecurityGroup: prometheus.NewDesc(
            "exoscale_security_group",
            "Exoscale Security Group",
            []string{"security_group"}, nil,
        ),
        SecurityGroupRule: prometheus.NewDesc(
            "exoscale_security_group_rule",
            "Exoscale Security Group Rule",
            []string{"security_group", "start_port", "end_port", "direction", "network", "protocol"}, nil,
        ),
    }
}

func NewElasticIPPrometheusMetricsCollector(ctx context.Context, cli v3.Client) *ElasticIPPrometheusMetricsCollector {
    return &ElasticIPPrometheusMetricsCollector{
        Context: ctx,
        Client: cli,
        Object: prometheus.NewDesc(
            "exoscale_elastic_ip",
            "Exoscale Elastic IP",
            []string{"ip"}, nil,
        ),
    }
}

func NewPrivateNetworkPrometheusMetricsCollector(ctx context.Context, cli v3.Client) *PrivateNetworkPrometheusMetricsCollector {
    return &PrivateNetworkPrometheusMetricsCollector{
        Context: ctx,
        Client: cli,
        Object: prometheus.NewDesc(
            "exoscale_private_network",
            "Exoscale Private Network",
            []string{"name", "vni"}, nil,
        ),
    }
}

func (collector *LoadBalancerPrometheusMetricsCollector) Describe(channel chan<- *prometheus.Desc) {
    channel <- collector.LoadBalancer
    channel <- collector.Service
}

func (collector *SecurityGroupPrometheusMetricsCollector) Describe(channel chan<- *prometheus.Desc) {
    channel <- collector.SecurityGroup
    channel <- collector.SecurityGroupRule
}

func (collector *ElasticIPPrometheusMetricsCollector) Describe(channel chan<- *prometheus.Desc) {
    channel <- collector.Object
}

func (collector *PrivateNetworkPrometheusMetricsCollector) Describe(channel chan<- *prometheus.Desc) {
    channel <- collector.Object
}

func (collector *SecurityGroupPrometheusMetricsCollector) Collect(channel chan<- prometheus.Metric) {
    response, err := collector.Client.ListSecurityGroups(collector.Context)
    if err != nil {
        panic(fmt.Sprintf("error initializing security group counter: %v", err))
    }

    securityGroupCounter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "exoscale_security_group_count",
        Help: "Exoscale Load Balancer Counter",
    })

    securityGroupRuleCounter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "exoscale_security_group_rule_count",
        Help: "Exoscale Load Balancer Counter",
    })

    for _, securityGroup := range response.SecurityGroups {
        channel <- prometheus.MustNewConstMetric(
            collector.SecurityGroup,
            prometheus.GaugeValue,
            float64(1),
            string(securityGroup.Name),
        )
        for _, rule := range securityGroup.Rules {
            channel <- prometheus.MustNewConstMetric(
                collector.SecurityGroupRule,
                prometheus.GaugeValue,
                float64(1),
                string(securityGroup.Name),
                fmt.Sprint(rule.StartPort),
                fmt.Sprint(rule.EndPort),
                string(rule.FlowDirection),
                string(rule.Network),
                string(rule.Protocol),
            )
            securityGroupRuleCounter.Inc()
        }
        securityGroupCounter.Inc()
    }

    channel <- securityGroupCounter
    channel <- securityGroupRuleCounter
}

func (collector *PrivateNetworkPrometheusMetricsCollector) Collect(channel chan<- prometheus.Metric) {
    privateNetworkCounter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "exoscale_private_network_count",
        Help: "Exoscale Private Network Counter",
    })

    response, err := collector.Client.ListPrivateNetworks(collector.Context)
    if err != nil {
        panic(fmt.Sprintf("error initializing domain counter: %v", err))
    }

    for _, privateNetwork := range response.PrivateNetworks {
        channel <- prometheus.MustNewConstMetric(
            collector.Object,
            prometheus.GaugeValue,
            float64(1),
            string(privateNetwork.Name),
            fmt.Sprint(privateNetwork.Vni),
        )
        privateNetworkCounter.Inc()
    }

    channel <- privateNetworkCounter
}

func (collector *ElasticIPPrometheusMetricsCollector) Collect(channel chan<- prometheus.Metric) {
    elasticIpCounter := prometheus.NewCounter(prometheus.CounterOpts{
        Name: "exoscale_elastic_ip_count",
        Help: "Exoscale Load Balancer Counter",
    })

    response, err := collector.Client.ListElasticIPS(collector.Context)
    if err != nil {
        panic(fmt.Sprintf("error initializing domain counter: %v", err))
    }

    for _, elasticIp := range response.ElasticIPS {
        channel <- prometheus.MustNewConstMetric(
            collector.Object,
            prometheus.GaugeValue,
            float64(1),
            string(elasticIp.IP),
        )
        elasticIpCounter.Inc()
    }

    channel <- elasticIpCounter
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
}