package collector

import (
	"net/http"
	"github.com/exoscale/egoscale/v3"
	"testing"
)

var dummyAPIKeys = v3.ListAPIKeysResponse {
	APIKeys: []v3.IAMAPIKey{
		{
			Key: "dummyKey",
			Name: "dummyName",
			RoleID: "dummyUUID",
		},
	},
}

func SetupIAMTestEndpoints() {
	http.HandleFunc("/api-key", HandleTestAPIKeysResponse)
}

func HandleTestAPIKeysResponse(w http.ResponseWriter, r *http.Request) {
	WriteObjectToResponse(w, r, dummyAPIKeys)
}

func TestIAMMetricsExist(t *testing.T) {
	metrics := GetTestMetrics(t)

	_, err := CheckMetricExists("exoscale_iam_key", metrics)
	if err != nil {
		t.Fail()
	}
}