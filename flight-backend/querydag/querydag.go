package querydag

import (
	"fmt"
	"strings"

	"github.com/frrad/flight-search/flight-backend/legfinder"
)

type DAG struct {
	// List of nodes, first is "start" last is "end" but no other constraints on order
	Nodes []Node
}

type Node struct {
	IsAirport  bool
	Name       string
	FlightsOut []Flight
}

type Flight struct {
	ToNode int      // index of destination node in Nodes list
	Dates  []string // omit == any date
}

func (dag DAG) Viz() string {
	ans := "digraph {\n"
	ans += "  rankdir=LR\n"
	for i, n := range dag.Nodes {
		shape := "oval"
		if n.IsAirport {
			shape = "box"
		}

		ans += fmt.Sprintf("  %d [label=\"%s\" shape=\"%s\"];\n", i, n.Name, shape)
		for _, flight := range n.FlightsOut {
			ans += fmt.Sprintf("  %d -> %d [ label=\"%s\" ];\n", i, flight.ToNode, strings.Join(flight.Dates, ";"))
		}
	}
	ans += "}\n"

	return ans
}

func (dag DAG) AllSolutions() []legfinder.TripSpec {
	paths := dag.pathsFromI(0)

	specs := []legfinder.TripSpec{}
	for _, path := range paths {
		specs = append(specs, dag.pathToSpecs(path))
	}
	return specs
}

// Extract leg specifications from a path through DAG
func (dag DAG) pathToSpecs(path []int) []legfinder.LegSpec {
	specs := []legfinder.LegSpec{}

	state := "start"
	thisSpec := legfinder.LegSpec{}
	for i := 0; i < len(path); i += 2 {
		nodeIndex, flightIndex := path[i], path[i+1]
		node := dag.Nodes[nodeIndex]
		flight := node.FlightsOut[flightIndex]

		switch state {
		case "start":
			if !node.IsAirport {
				break
			}
			thisSpec.Origin = node.Name
			if len(flight.Dates) > 0 {
				thisSpec.Dates = append([]string{}, flight.Dates...)
				state = "accepting-dates"
				break
			}
			state = "accepting-anydate"

		case "accepting-anydate":
			// We're looking for the destination airport, and no dates
			// are yet specified.
			if node.IsAirport {
				panic("underspecified")
			}
			if len(flight.Dates) > 0 {
				thisSpec.Dates = append([]string{}, flight.Dates...)
				state = "accepting-dates"
			}

		case "accepting-dates":
			if node.IsAirport {
				thisSpec.Destination = node.Name
				specs = append(specs, thisSpec)
				thisSpec = legfinder.LegSpec{Origin: node.Name}
				if len(flight.Dates) > 0 {
					thisSpec.Dates = append([]string{}, flight.Dates...)
					state = "accepting-dates"
					break
				}
				state = "accepting-anydate"
				break
			}

			if len(flight.Dates) > 0 {
				newDates := []string{}
				for _, d1 := range flight.Dates {
					for _, d2 := range thisSpec.Dates {
						if d1 == d2 {
							newDates = append(newDates, d1)
						}
					}
				}

				thisSpec.Dates = newDates
			}

		}

	}

	return specs
}

// List of all paths from Ith node to the end
func (dag DAG) pathsFromI(i int) [][]int {
	if isEnd(i, dag.Nodes) {
		return [][]int{{}}
	}
	ans := [][]int{}

	node := dag.Nodes[i]
	for j, flight := range node.FlightsOut {
		for _, path := range dag.pathsFromI(flight.ToNode) {
			ans = append(ans, append([]int{i, j}, path...))
		}
	}

	if len(ans) == 0 {
		return [][]int{{}}
	}

	return ans
}

func isEnd(i int, nodes []Node) bool {
	return i == len(nodes)-1
}
