package collector

import (
	"context"
	"fmt"
	"testing"
	"github.com/exoscale/egoscale/v3"
	"github.com/exoscale/egoscale/v3/credentials"
	"net/http"
	// "encoding/json"
	"os"
)

// func tResponse(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	organization := v3.Organization{
// 		Address: "test addr",
// 	}
// 	if err := json.NewEncoder(w).Encode(organization); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

	// exoscaleCredentials := 
	// exoClient, err := v3.NewClient(exoscaleCredentials)
	// if err != nil {
	// 	panic(fmt.Sprintf("unable to initialize Exoscale API V3 client: %v", err))
	// }
	// exoClient = exoClient.WithEndpoint("http://localhost:9998")
	// 

var dummyExoscaleCredentials = credentials.NewStaticCredentials("EXO", "EXO")
var dummyExoscaleClient *v3.Client

func setupWebServers(m *testing.M) int {
	SetupOrganizationTestEndpoints()

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