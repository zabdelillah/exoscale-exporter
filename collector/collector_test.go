package collector

import (
	"context"
	"fmt"
	"testing"
	"github.com/exoscale/egoscale/v3"
	"github.com/exoscale/egoscale/v3/credentials"
	"net/http"
	"encoding/json"
	"os"
	"errors"
	"io/ioutil"
	"strings"
)

var dummyExoscaleCredentials = credentials.NewStaticCredentials("EXO", "EXO")
var dummyExoscaleClient *v3.Client

func WriteObjectToResponse(w http.ResponseWriter, r *http.Request, instance interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(instance); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func CheckMetricExists(metric string, metrics string) (string, error) {
	if !strings.Contains(metrics, metric) {
		return "", errors.New(fmt.Sprintf("Metric '%s' not found", metric))
	}

	return metric, nil
}

func CheckMetricsExist(t *testing.T, metricsToCheck []string, metrics string) ([]string, []error) {
	var notFoundErrors []error
	var foundMetrics []string

	for i := range(metricsToCheck) {
        _, err := CheckMetricExists(metricsToCheck[i], metrics)
        if err != nil {
            notFoundErrors = append(notFoundErrors, err)
        } else {
        	foundMetrics = append(foundMetrics, metricsToCheck[i])
        }
    }

    return foundMetrics, notFoundErrors
}

func setupWebServers(m *testing.M) int {
	SetupBlockStorageTestEndpoints()
	SetupDNSTestEndpoints()
	SetupIAMTestEndpoints()
	SetupInstanceTestEndpoints()
	SetupOrganizationTestEndpoints()
	SetupSOSTestEndpoints()
	SetupSKSTestEndpoints()

	PrepareCollector(context.Background(), dummyExoscaleClient)

	go http.ListenAndServe(":9998", nil)
	return m.Run()
}

func TestMain(m *testing.M) {
	var err error
	dummyExoscaleClient, err = v3.NewClient(dummyExoscaleCredentials)
	if err != nil {
		panic(fmt.Sprintf("unable to initialize Exoscale API V3 client: %v", err))
	}
	dummyExoscaleClient = dummyExoscaleClient.WithEndpoint("http://localhost:9998")
	os.Exit(setupWebServers(m))
}