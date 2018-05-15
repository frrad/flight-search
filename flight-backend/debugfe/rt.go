package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/frrad/flight-search/flight-backend/querydag"
)

type rtData struct {
	StartDate     string
	StartAirports string
	EndDate       string
	EndAirports   string
}

func rtHandler(w http.ResponseWriter, req *http.Request) {
	rtData := collectRTData(req)
	jsonString := jsonReqFromRT(rtData)

	err := respond(w, jsonString, rtData)
	if err != nil {
		panic(err)
	}
}

func jsonReqFromRT(data rtData) string {
	dag := querydag.DAG{
		Nodes: []querydag.Node{
			querydag.Node{
				IsAirport:  false,
				Name:       "start",
				FlightsOut: []querydag.Flight{},
			},
			querydag.Node{
				IsAirport: false,
				Name:      "end",
			},
		},
	}

	starters := strings.Split(data.StartAirports, ",")
	flightOutIdx := len(starters) + 2
	for i, end := range starters {
		dag.Nodes = append(dag.Nodes, querydag.Node{
			IsAirport:  true,
			Name:       end,
			FlightsOut: []querydag.Flight{querydag.Flight{ToNode: flightOutIdx}},
		})
		dag.Nodes[0].FlightsOut = append(dag.Nodes[0].FlightsOut,
			querydag.Flight{
				ToNode: i + 2,
			},
		)
	}

	// flight out
	dag.Nodes = append(dag.Nodes, querydag.Node{
		IsAirport:  false,
		Name:       "flight out",
		FlightsOut: []querydag.Flight{},
	})

	enders := strings.Split(data.EndAirports, ",")
	flightBackIdx := len(enders) + len(starters) + 3
	for i, end := range enders {
		dag.Nodes = append(dag.Nodes, querydag.Node{
			IsAirport:  true,
			Name:       end,
			FlightsOut: []querydag.Flight{querydag.Flight{ToNode: flightBackIdx}},
		})
		dag.Nodes[flightOutIdx].FlightsOut = append(dag.Nodes[flightOutIdx].FlightsOut,
			querydag.Flight{
				ToNode: i + len(starters) + 3,
			},
		)
	}

	// flight back
	dag.Nodes = append(dag.Nodes, querydag.Node{
		IsAirport:  false,
		Name:       "flight back",
		FlightsOut: []querydag.Flight{},
	})

	for i, end := range starters {
		dag.Nodes = append(dag.Nodes, querydag.Node{
			IsAirport:  true,
			Name:       end,
			FlightsOut: []querydag.Flight{querydag.Flight{ToNode: 1}},
		})
		dag.Nodes[flightBackIdx].FlightsOut = append(dag.Nodes[flightBackIdx].FlightsOut,
			querydag.Flight{
				ToNode: i + len(starters) + len(enders) + 4,
			},
		)
	}

	ans, err := json.Marshal(dag)
	if err != nil {
		return ""
	}
	return string(ans)
}

func collectRTData(req *http.Request) rtData {
	return rtData{
		StartDate:     req.FormValue("startdate"),
		EndDate:       req.FormValue("enddate"),
		StartAirports: req.FormValue("startairports"),
		EndAirports:   req.FormValue("endairports"),
	}
}
