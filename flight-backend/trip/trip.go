package trip

import (
	"log"
	"sort"
	"time"

	"github.com/frrad/flight-search/flight-backend/legfinder"
)

type TripPlanner struct {
	finder legfinder.LegFinder
}

func NewPlanner(f legfinder.LegFinder) *TripPlanner {
	return &TripPlanner{finder: f}
}

type TripOption struct {
	Id    int
	Legs  []legfinder.Leg
	Price int
}

func (tp *TripPlanner) ListOptions(tripSpecs []legfinder.TripSpec) ([]TripOption, error) {
	ans := make([][][]legfinder.Leg, len(tripSpecs))
	memo := make(map[string][]legfinder.Leg)

	for i, trip := range tripSpecs {

		ans[i] = make([][]legfinder.Leg, len(trip))

		for j, leg := range trip {
			if solns, ok := memo[leg.Hash()]; ok {
				ans[i][j] = solns
			} else {
				lookup, err := tp.finder.Find(leg)

				if err == nil {
					memo[leg.Hash()] = lookup
					ans[i][j] = memo[leg.Hash()]
				} else {
					log.Printf("error looking up leg\nleg:%+v\nerror:%+v", leg, err)
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

func consistentOptions(i int, options [][]legfinder.Leg, departAfter time.Time) []TripOption {
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
				Legs:  []legfinder.Leg{option},
				Price: option.Price,
			})
			continue
		}

		landsAt := option.Segments[len(option.Segments)-1].ArrivalTime
		subProblem := consistentOptions(i, options[1:], landsAt)
		for _, soln := range subProblem {
			ans = append(ans, TripOption{
				Id:    i,
				Legs:  append([]legfinder.Leg{option}, soln.Legs...),
				Price: option.Price + soln.Price,
			})

		}

	}

	return ans
}
