package collector

import (
    "net/http"
    "github.com/exoscale/egoscale/v3"
    "testing"
)

var dummyDNSResponse = v3.ListDNSDomainsResponse {
    DNSDomains: []v3.DNSDomain{
        {
            ID: "dummyDNSRecord",
            UnicodeName: "test.dev",
        },
    },
}

var dummyDNSRecordResponse = v3.ListDNSDomainRecordsResponse {
    DNSDomainRecords: []v3.DNSDomainRecord {
        {
            Content: "www.test.dev",
            ID: "dummyRecordID",
            Name: "test.dev",
            Priority: 10,
            Ttl: 300,
            Type: v3.DNSDomainRecordTypeA,
        },
    },
}

func SetupDNSTestEndpoints() {
    http.HandleFunc("/dns-domain", HandleTestDomainResponse)
    http.HandleFunc("/dns-domain/dummyDNSRecord/record", HandleTestDomainRecordResponse)
}

func HandleTestDomainResponse(w http.ResponseWriter, r *http.Request) {
    WriteObjectToResponse(w, r, dummyDNSResponse)
}

func HandleTestDomainRecordResponse(w http.ResponseWriter, r *http.Request) {
    WriteObjectToResponse(w, r, dummyDNSRecordResponse)
}

func TestDNSMetricsExist(t *testing.T) {
    metrics := GetTestMetrics(t)

    metricsToCheck := []string {
        "exoscale_domain",
        "exoscale_domain_count",
        "exoscale_domain_record",
        "exoscale_domain_record_count",
    }

    _, errs := CheckMetricsExist(t, metricsToCheck, metrics)
    for i := range(errs) {
        t.Errorf("Instance Metric Check Failed: %v", errs[i])
    }
}