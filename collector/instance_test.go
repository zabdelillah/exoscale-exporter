package collector

import (
    "net/http"
    "github.com/exoscale/egoscale/v3"
    "testing"
)

var dummyInstanceType = v3.InstanceType{
    ID: "dummy",
    Family: v3.InstanceTypeFamilyStandard,
    Size: v3.InstanceTypeSizeMedium,
}

var dummyInstanceResponse = v3.ListInstancesResponse {
    Instances: []v3.ListInstancesResponseInstances{
        {
            ID: "instanceUUID",
            Name: "dummyInstance",
            InstanceType: &dummyInstanceType,
        },
    },
}

var dummyInstanceTypes = v3.ListInstanceTypesResponse {
    InstanceTypes: []v3.InstanceType{
        dummyInstanceType,
    },
}

var dummyInstancePoolResponse = v3.ListInstancePoolsResponse {
    InstancePools: []v3.InstancePool{
        {
            State: v3.InstancePoolStateRunning,
            Name: "dummyInstancePool",
            ID: "dummyInstancePool",
        },
    },
}

func SetupInstanceTestEndpoints() {
    http.HandleFunc("/instance", HandleTestInstanceResponse)
    http.HandleFunc("/instance-type/dummy", HandleTestDummyInstanceTypeResponse)
    http.HandleFunc("/instance-pool", HandleTestInstancePoolResponse)
}

func HandleTestInstanceResponse(w http.ResponseWriter, r *http.Request) {
    WriteObjectToResponse(w, r, dummyInstanceResponse)
}

func HandleTestDummyInstanceTypeResponse(w http.ResponseWriter, r *http.Request) {
    WriteObjectToResponse(w, r, dummyInstanceTypes)
}

func HandleTestInstancePoolResponse(w http.ResponseWriter, r *http.Request) {
    WriteObjectToResponse(w, r, dummyInstancePoolResponse)
}

func TestInstanceMetricsExist(t *testing.T) {
    metrics := GetTestMetrics(t)

    metricsToCheck := []string {
        "exoscale_instances_count",
        "exoscale_instance_up",
        "exoscale_instance_cpus",
        "exoscale_instance_gpus",
        "exoscale_instance_memory",
    }

    _, errs := CheckMetricsExist(t, metricsToCheck, metrics)
    for i := range(errs) {
        t.Errorf("Instance Metric Check Failed: %v", errs[i])
    }
}

func TestInstancePoolMetricsExist(t *testing.T) {
    metrics := GetTestMetrics(t)

    metricsToCheck := []string {
        "exoscale_instance_pool_up",
        "exoscale_instance_pool_size",
    }

    _, errs := CheckMetricsExist(t, metricsToCheck, metrics)
    for i := range(errs) {
        t.Errorf("Instance Metric Check Failed: %v", errs[i])
    }
}