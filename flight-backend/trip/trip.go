package trip

import (
	"github.com/frrad/flight-search/flight-backend/qpx"
)

type TripPlanner struct {
	finder *qpx.QPXFinder
}

func NewPlanner(f *qpx.QPXFinder) *TripPlanner {
	return &TripPlanner{finder: f}
}

// Choose an option for each leg
type TripOptions [][]qpx.Leg

func (tp *TripPlanner) ListOptions(tripSpecs [][]qpx.LegSpec) []TripOptions {
	ans := make([]TripOptions, len(tripSpecs))

	memo := make(map[string][]qpx.Leg)
	for i, trip := range tripSpecs {

		ans[i] = make([][]qpx.Leg, len(trip))

		for j, leg := range trip {
			if solns, ok := memo[leg.Hash()]; ok {
				ans[i][j] = solns
			} else {
				lookup, err := tp.finder.Find(leg)

				if err == nil {
					memo[leg.Hash()] = lookup
					ans[i][j] = memo[leg.Hash()]
				} else {
					panic(err)
				}
			}
		}
	}

	return ans
}
