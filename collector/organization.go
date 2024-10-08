package collector

import (
    "github.com/exoscale/egoscale/v3"
    "github.com/prometheus/client_golang/prometheus"
    "fmt"
    "context"
)

type OrganizationPrometheusMetricsCollector struct {
    Context context.Context
    Client v3.Client
    Balance *prometheus.Desc
    Usage *prometheus.Desc
    SSHKey *prometheus.Desc
}

func NewOrganizationPrometheusMetricsCollector(ctx context.Context, cli v3.Client) *OrganizationPrometheusMetricsCollector {
    return &OrganizationPrometheusMetricsCollector{
        Context: ctx,
        Client: cli,
        Balance: prometheus.NewDesc(
            "exoscale_organization_balance",
            "Current balance on exoscale organization",
            []string{"organization_id", "organization_name"}, nil,
        ),
        Usage: prometheus.NewDesc(
            "exoscale_organization_usage",
            "Current balance on exoscale organization",
            []string{"organization_id", "organization_name"}, nil,
        ),
        SSHKey: prometheus.NewDesc(
            "exoscale_ssh_key",
            "SSH Keys stored on Exoscale",
            []string{"fingerprint", "name"}, nil,
        ),
    }
}

func (collector *OrganizationPrometheusMetricsCollector) Describe(channel chan<- *prometheus.Desc) {
    channel <- collector.Balance
    channel <- collector.Usage
    channel <- collector.SSHKey
}

func (collector *OrganizationPrometheusMetricsCollector) Collect(channel chan<- prometheus.Metric) {
    organization, err := collector.Client.GetOrganization(collector.Context)

    if err != nil {
        panic(fmt.Sprintf("unable to initialize Exoscale API V3 client: %v", err))
    }

    channel <- prometheus.MustNewConstMetric(
        collector.Balance,
        prometheus.GaugeValue,
        organization.Balance,
        organization.ID.String(), 
        organization.Name,
    )

    channel <- prometheus.MustNewConstMetric(
        collector.Usage,
        prometheus.GaugeValue,
        (0 - organization.Balance),
        organization.ID.String(), 
        organization.Name,
    )

    sshKeys, err := collector.Client.ListSSHKeys(collector.Context)
    for _, sshKey := range(sshKeys.SSHKeys) {
        channel <- prometheus.MustNewConstMetric(
            collector.SSHKey,
            prometheus.GaugeValue,
            float64(1),
            sshKey.Fingerprint, 
            sshKey.Name,
        )
    }
}