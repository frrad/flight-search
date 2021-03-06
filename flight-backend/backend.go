package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/frrad/flight-search/flight-backend/amadeus"
	"github.com/frrad/flight-search/flight-backend/legcacher"
	"github.com/frrad/flight-search/flight-backend/querydag"
	"github.com/frrad/flight-search/flight-backend/trip"
)

type backendHandler struct {
	planner *trip.TripPlanner
}

func (b *backendHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("received request...")

	decoder := json.NewDecoder(req.Body)

	var query querydag.DAG
	err := decoder.Decode(&query)

	if err != nil {
		log.Println("\n")
		log.Println(err)
		log.Println("invalid request")
		return
	}

	defer req.Body.Close()

	sols := query.AllSolutions()

	options, err := b.planner.ListOptions(sols)
	if err != nil {
		fmt.Fprintf(w, "ERROR listing options")
		return
	}

	response, err := json.Marshal(options)
	if err != nil {
		fmt.Fprintf(w, "ERROR marshalling response")
		return
	}

	fmt.Fprintf(w, "%s", string(response))
}

func main() {
	finder := amadeus.NewAmadeusFinder(os.Getenv("AMADEUSKEY"))

	ttl, err := time.ParseDuration("5h")
	if err != nil {
		log.Fatal("can't parse TTL")
	}

	wrapped := legcacher.NewLegCacher(finder, ttl)
	tripPlanner := trip.NewPlanner(wrapped)

	handler := &backendHandler{
		planner: tripPlanner,
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Fatal(server.ListenAndServe())

}
