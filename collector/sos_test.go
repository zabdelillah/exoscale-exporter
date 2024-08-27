package collector

import (
    "net/http"
    "github.com/exoscale/egoscale/v3"
    "testing"
)

var dummySOSBucketUsage = v3.ListSOSBucketsUsageResponse {
    SOSBucketsUsage: []v3.SOSBucketUsage{
        {
            Name: "dummySOSBucket",
            Size: 64000,
        },
    },
}

func SetupSOSTestEndpoints() {
    http.HandleFunc("/sos-buckets-usage", HandleTestSOSResponse)
}

func HandleTestSOSResponse(w http.ResponseWriter, r *http.Request) {
    WriteObjectToResponse(w, r, dummySOSBucketUsage)
}

func TestSOSMetricsExist(t *testing.T) {
    metrics := GetTestMetrics(t)

    metricsToCheck := []string {
        "exoscale_sos_bucket_count",
    }

    _, errs := CheckMetricsExist(t, metricsToCheck, metrics)
    for i := range(errs) {
        t.Errorf("Instance Metric Check Failed: %v", errs[i])
    }
}