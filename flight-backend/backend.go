package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/frrad/flight-search/flight-backend/qpx"
	"github.com/frrad/flight-search/flight-backend/querydag"
	"github.com/frrad/flight-search/flight-backend/querytree"
	"github.com/frrad/flight-search/flight-backend/trip"
)

func test(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)

	var t querytree.Tree

	err := decoder.Decode(&t)

	if err != nil {
		log.Println("\n")
		log.Println(err)
		log.Println("invalid request")
		return
	}

	defer req.Body.Close()

	log.Println("\n" + t.DispFormat(0))
	log.Println(len(t.Reduce()))
	// for _, res := range t.Resolve() {
	// 	log.Println("asdfasfd")
	// 	log.Println("\n" + res.DispFormat(0))
	// }
	log.Println("test")

}

func main() {

	// http.HandleFunc("/backend", test)
	// http.ListenAndServe(":8080", nil)

	test := querydag.DAG{
		Nodes: []querydag.Node{
			{
				Name: "start",
				FlightsOut: []querydag.Flight{
					{ToNode: 1},
					{ToNode: 2},
				},
			},
			{
				Name:      "SFO",
				IsAirport: true,
				FlightsOut: []querydag.Flight{
					{
						ToNode: 3,
						Dates:  []string{"2018-04-01"},
					},
				},
			},
			{
				Name:      "OAK",
				IsAirport: true,
				FlightsOut: []querydag.Flight{
					{
						ToNode: 3,
						Dates:  []string{"2018-04-01", "2018-04-02"},
					},
				},
			},
			{
				Name:      "<>",
				IsAirport: false,
				FlightsOut: []querydag.Flight{
					{ToNode: 4},
					{ToNode: 5},
				},
			},
			{
				Name:      "MCO",
				IsAirport: true,
				FlightsOut: []querydag.Flight{
					{
						ToNode: 6,
					},
				},
			},

			{
				Name:      "MIA",
				IsAirport: true,
				FlightsOut: []querydag.Flight{
					{
						ToNode: 6,
					},
				},
			},

			{Name: "end"},
		},
	}

	fmt.Println(test.Viz())

	finder := qpx.NewQPXFinder(os.Getenv("QPXAPIKEY"))
	planner := trip.NewPlanner(finder)

	sols := test.AllSolutions()
	fmt.Println(planner.ListOptions(sols))

}
