package main
/**
 * Prometheus exporter for the Exoscale cloud provider
 */

import (
	"zai.dev/m/v2/collector"
	"github.com/exoscale/egoscale/v3"
	"github.com/exoscale/egoscale/v3/credentials"
	"net/http"
	"context"
	"fmt"
	"os"
	"log"
	"flag"
)

func main() {
	var port string

	flag.StringVar(&port, "listen", ":9999", "address and port to listen on")
	flag.Parse()

	ctx := context.Background()
	api_key, api_key_exists := os.LookupEnv("EXOSCALE_API_KEY")
	api_secret, api_secret_exists := os.LookupEnv("EXOSCALE_API_SECRET")

	if !api_key_exists {
		log.Printf("EXOSCALE_API_KEY not provided")
	}
	if !api_secret_exists {
		log.Printf("EXOSCALE_API_KEY not provided")
	}

	exoscaleCredentials := credentials.NewStaticCredentials(api_key, api_secret)
	exoClient, err := v3.NewClient(exoscaleCredentials)

	if err != nil {
		panic(fmt.Sprintf("unable to initialize Exoscale API V3 client: %v", err))
	}

	collector.PrepareCollector(ctx, exoClient)

	http.ListenAndServe(port, nil)
}