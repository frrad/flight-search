package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/frrad/flight-search/flight-backend/qpx"
	"github.com/frrad/flight-search/flight-backend/querydag"
	"github.com/frrad/flight-search/flight-backend/trip"
)

func queryHandler(w http.ResponseWriter, req *http.Request) {
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

	finder := qpx.NewQPXFinder(os.Getenv("QPXAPIKEY"))
	planner := trip.NewPlanner(finder)

	sols := query.AllSolutions()

	options, err := planner.ListOptions(sols)
	if err != nil {
		fmt.Fprintf(w, "ERROR")
		return
	}

	response, err := json.Marshal(options)
	if err != nil {
		fmt.Fprintf(w, "ERROR")
		return
	}

	fmt.Fprintf(w, "%s", string(response))

}

func main() {
	http.HandleFunc("/backend", queryHandler)
	http.ListenAndServe(":8080", nil)
}
