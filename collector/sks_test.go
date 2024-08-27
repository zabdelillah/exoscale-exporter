package collector

import (
	"net/http"
	"github.com/exoscale/egoscale/v3"
	"testing"
)

var dummySKSCluster = v3.ListSKSClustersResponse {
	SKSClusters: []v3.SKSCluster{
		{
			Name: "dummyName",
			State: v3.SKSClusterStateRunning,
			Nodepools: []v3.SKSNodepool{
				{
					ID: "dummySKSNodePoolID",
					Name: "dummySKSNodePoolName",
					Version: "dummySKSNodePoolVersion",
					State: v3.SKSNodepoolStateRunning,
					Size: 3,
					DiskSize: 16,
				},
			},
		},
	},
}

func SetupSKSTestEndpoints() {
	http.HandleFunc("/sks-cluster", HandleTestSKSResponse)
}

func HandleTestSKSResponse(w http.ResponseWriter, r *http.Request) {
	WriteObjectToResponse(w, r, dummySKSCluster)
}

func TestSKSMetricsExist(t *testing.T) {
	metrics := GetTestMetrics(t)

	metricsToCheck := []string {
		"exoscale_sks_cluster_count",
		"exoscale_sks_cluster_size",
        "exoscale_sks_cluster_up",
        "exoscale_sks_nodepool_up",
        "exoscale_sks_nodepool_size",
        "exoscale_sks_nodepool_disk_size",
    }

    _, errs := CheckMetricsExist(t, metricsToCheck, metrics)
    for i := range(errs) {
        t.Errorf("Instance Metric Check Failed: %v", errs[i])
    }
}