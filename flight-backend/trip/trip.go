package trip

import (
	"log"
	"sort"
	"time"

	"github.com/frrad/flight-search/flight-backend/qpx"
)

type TripPlanner struct {
	finder *qpx.QPXFinder
}

func NewPlanner(f *qpx.QPXFinder) *TripPlanner {
	return &TripPlanner{finder: f}
}

type TripOption struct {
	Id    int
	Legs  []qpx.Leg
	Price int
}

func (tp *TripPlanner) ListOptions(tripSpecs []qpx.TripSpec) ([]TripOption, error) {
	ans := make([][][]qpx.Leg, len(tripSpecs))
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
					return nil, err
				}
			}
		}
	}

	log.Println("Finding consistent options...")

	tripOptions := []TripOption{}
	for i, options := range ans {
		tripOptions = append(tripOptions, consistentOptions(i, options, time.Now())...)
	}

	log.Printf("Found %d options", len(tripOptions))

	sort.Slice(tripOptions, func(i, j int) bool {
		return tripOptions[i].Price < tripOptions[j].Price
	})

	return tripOptions, nil
}

func consistentOptions(i int, options [][]qpx.Leg, departAfter time.Time) []TripOption {
	ans := []TripOption{}
	if len(options) == 0 {
		return ans
	}

	for _, option := range options[0] {
		if option.Segments[0].DepartureTime.Before(departAfter) {
			continue
		}

		if len(options) == 1 {
			ans = append(ans, TripOption{
				Id:    i,
				Legs:  []qpx.Leg{option},
				Price: option.Price,
			})
			continue
		}

		landsAt := option.Segments[len(option.Segments)-1].ArrivalTime
		subProblem := consistentOptions(i, options[1:], landsAt)
		for _, soln := range subProblem {
			ans = append(ans, TripOption{
				Id:    i,
				Legs:  append([]qpx.Leg{option}, soln.Legs...),
				Price: option.Price + soln.Price,
			})

		}

	}

	return ans
}
