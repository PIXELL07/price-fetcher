package main

import (
	"flag"
	"log"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3000", "listen address the service is running on")
	flag.Parse()

	svc := NewLoggingService(NewMetricService(&priceFetcher{}))
	server := NewJSONAPIServer(*listenAddr, svc)

	log.Printf("price-fetcher listening on %s", *listenAddr)
	log.Printf("endpoints: /price  /market  /info  /tickers  /batch")
	server.Run()
}
