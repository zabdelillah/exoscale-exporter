package collector

import (
	"context"
	"fmt"
	"testing"
	"github.com/exoscale/egoscale/v3"
	// "github.com/exoscale/egoscale/v3/credentials"
	"net/http"
	"encoding/json"
	"strings"
	 "io/ioutil"
)

var dummyOrganization = v3.Organization{
	Address: "test",
	Balance: 5.00,
	City: "Geneva",
	Country: "Switzerland",
	Currency: "CHF",
	Name: "Go Tests",
	ID: "go-tests",
	Postcode: "G0T 3ST5",
}

func SetupOrganizationTestEndpoints() {
	http.HandleFunc("/organization", HandleTestOrganizationResponse)
	// Others
	http.HandleFunc("/instance", HandleTestOrganizationResponse)
	http.HandleFunc("/instance-pool", HandleTestOrganizationResponse)
	http.HandleFunc("/sos-buckets-usage", HandleTestOrganizationResponse)
	http.HandleFunc("/snapshot", HandleTestOrganizationResponse)
	http.HandleFunc("/sks-cluster", HandleTestOrganizationResponse)
	http.HandleFunc("/block-storage", HandleTestOrganizationResponse)
	http.HandleFunc("/api-key", HandleTestOrganizationResponse)
}

func HandleTestOrganizationResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dummyOrganization); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Tests that the local http server works and we're injecting
// the dummyOrganization properly into the egoscale/v3 api
func TestOrganizationDummyResponse(t *testing.T) {
	organization, err := dummyExoscaleClient.GetOrganization(context.Background())
	if err != nil {
		fmt.Printf("unable to initialize Exoscale API V3 client: %v", err)
	}
	if ((organization.Address != dummyOrganization.Address) ||
		(organization.Balance != dummyOrganization.Balance) ||
		(organization.City != dummyOrganization.City) ||
		(organization.Country != dummyOrganization.Country) ||
		(organization.Currency != dummyOrganization.Currency) ||
		(organization.Name != dummyOrganization.Name) ||
		(organization.Postcode != dummyOrganization.Postcode)) {
		t.Fail()
	}
}

func GetTestMetrics(t *testing.T) string {
	resp, err := http.Get("http://localhost:9998/metrics")
    if err != nil {
        t.Errorf("http.Get() error: %v", err)
    }
    defer resp.Body.Close()

    // Read the response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        t.Errorf("ioutil.ReadAll() error: %v", err)
    }

    return string(body)
}

func TestOrganizationMetricsExist(t *testing.T) {
	metrics := GetTestMetrics(t)

    balance_metric := fmt.Sprintf(
    	"exoscale_organization_balance{organization_id=\"%s\",organization_name=\"%s\"}",
    	dummyOrganization.ID, dummyOrganization.Name)

    usage_metric := fmt.Sprintf(
    	"exoscale_organization_usage{organization_id=\"%s\",organization_name=\"%s\"}",
    	dummyOrganization.ID, dummyOrganization.Name)

    // Check if the expected content exists in the response
    if !strings.Contains(metrics, balance_metric) {
        t.Errorf("Metric %s not found", "exoscale_organization_balance")
    }

    if !strings.Contains(metrics, usage_metric) {
        t.Errorf("Metric %s not found", "exoscale_organization_balance")
    }
}
