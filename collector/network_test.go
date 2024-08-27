package collector

import (
    "net/http"
    "github.com/exoscale/egoscale/v3"
    "testing"
)

var dummySecurityGroupsResponse = v3.ListSecurityGroupsResponse {
    SecurityGroups: []v3.SecurityGroup{
        {
            ID: "dummySecurityGroup",
            Name: "dummySecurityGroupName",
            Rules: []v3.SecurityGroupRule {
                {
                    StartPort: 80,
                    EndPort: 80,
                    FlowDirection: v3.SecurityGroupRuleFlowDirectionIngress,
                    Network: "0.0.0.0/0",
                    Protocol: v3.SecurityGroupRuleProtocolTCP,
                },
            },
        },
    },
}

var dummyPrivateNetworkResponse = v3.ListPrivateNetworksResponse {
    PrivateNetworks: []v3.PrivateNetwork {
        {
            ID: "dummyVPC",
            Name: "dummyVPCName",
            Vni: 0,
        },
    },
}

var dummyElasticIPResponse = v3.ListElasticIPSResponse {
    ElasticIPS: []v3.ElasticIP {
        {
            ID: "dummyElasticIP",
            IP: "0.0.0.0",
            Cidr: "0.0.0.0/0",
        },
    },
}

var dummyLoadBalancerResponse = v3.ListLoadBalancersResponse {
    LoadBalancers: []v3.LoadBalancer {
        {
            ID: "dummyLoadBalancer",
            Name: "dummyLoadBalancerName",
            State: v3.LoadBalancerStateCreating,
            Services: []v3.LoadBalancerService {
                {
                    ID: "dummyLoadBalancerService",
                    Name: "dummyLoadBalancerServiceName",
                    Port: 80,
                    State: v3.LoadBalancerServiceStateRunning,
                    TargetPort: 80,
                    Strategy: v3.LoadBalancerServiceStrategyRoundRobin,
                },
            },
        },
    },
}


func SetupNetworkingTestEndpoints() {
    http.HandleFunc("/load-balancer", HandleTestLoadBalancerResponse)
    http.HandleFunc("/elastic-ip", HandleTestElasticIPResponse)
    http.HandleFunc("/private-network", HandleTestPrivateNetworkResponse)
    http.HandleFunc("/security-group", HandleTestSecurityGroupResponse)
}

func HandleTestLoadBalancerResponse(w http.ResponseWriter, r *http.Request) {
    WriteObjectToResponse(w, r, dummyLoadBalancerResponse)
}
func HandleTestElasticIPResponse(w http.ResponseWriter, r *http.Request) {
    WriteObjectToResponse(w, r, dummyLoadBalancerResponse)
}
func HandleTestPrivateNetworkResponse(w http.ResponseWriter, r *http.Request) {
    WriteObjectToResponse(w, r, dummyLoadBalancerResponse)
}
func HandleTestSecurityGroupResponse(w http.ResponseWriter, r *http.Request) {
    WriteObjectToResponse(w, r, dummyLoadBalancerResponse)
}

func TestNetworkingMetricsExist(t *testing.T) {
    metrics := GetTestMetrics(t)

    metricsToCheck := []string {
        "exoscale_security_group",
        "exoscale_security_group_count",
        "exoscale_security_group_rule",
        "exoscale_private_network",
        "exoscale_private_network_count",
        "exoscale_elastic_ip",
        "exoscale_elastic_ip_count",
        "exoscale_load_balancer",
        "exoscale_load_balancer_count",
        "exoscale_load_balancer_service",
        "exoscale_load_balancer_service_count",
    }

    _, errs := CheckMetricsExist(t, metricsToCheck, metrics)
    for i := range(errs) {
        t.Errorf("Instance Metric Check Failed: %v", errs[i])
    }
}