package trip

import (
	"fmt"
	"log"
	"sort"
	"time"

	"golang.org/x/sync/errgroup"

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

func (tp *TripPlanner) FetchOptions(tripSpecs []legfinder.TripSpec) (map[string][]legfinder.Leg, error) {
	distinctLegs := []legfinder.LegSpec{}
	seenSpecs := map[string]bool{}
	for _, trip := range tripSpecs {
		for _, leg := range trip {
			if _, ok := seenSpecs[leg.Hash()]; ok {
				continue
			}
			distinctLegs = append(distinctLegs, leg)
			seenSpecs[leg.Hash()] = true
		}
	}

	type kv struct {
		k string
		v []legfinder.Leg
	}

	queries := make(chan legfinder.LegSpec, 3)
	queryResults := make(chan kv, 3)
	aggResults := make(chan map[string][]legfinder.Leg)

	go func(legs []legfinder.LegSpec) {
		for _, leg := range legs {
			queries <- leg
		}
		close(queries)
	}(distinctLegs)

	go func() {
		ans := map[string][]legfinder.Leg{}
		for x := range queryResults {
			ans[x.k] = x.v
		}
		aggResults <- ans
		close(aggResults)
	}()

	eg := &errgroup.Group{}

	for i := 0; i < 10; i++ {
		i := i
		eg.Go(func() error {
			for leg := range queries {
				log.Printf("Worker %d dispatching query\n", i)
				results, err := tp.finder.Find(leg)

				if err != nil {
					return err
				}

				queryResults <- kv{
					k: leg.Hash(),
					v: results,
				}
				log.Printf("Worker %d received response\n", i)
			}
			return nil
		})

	}

	err := eg.Wait()
	close(queryResults)

	if err != nil {
		return nil, err
	}

	return <-aggResults, nil
}

func (tp *TripPlanner) ListOptions(tripSpecs []legfinder.TripSpec) ([]TripOption, error) {
	ans := make([][][]legfinder.Leg, len(tripSpecs))
	memo, err := tp.FetchOptions(tripSpecs)
	if err != nil {
		log.Printf("error fetching options from legfinder. \nerror:%+v", err)
		return nil, err
	}

	for i, trip := range tripSpecs {
		ans[i] = make([][]legfinder.Leg, len(trip))

		for j, leg := range trip {
			if solns, ok := memo[leg.Hash()]; ok {
				ans[i][j] = solns
			} else {
				log.Printf("error finding leg results\nleg:%+v", leg)
				return nil, fmt.Errorf("Didn't find leg in fetched results")
			}
		}
	}

	log.Println("Finding consistent options...")

	tripOptions := []TripOption{}
	for i, options := range ans {
		theseOptions := consistentOptions(i, options, time.Now())

		if len(theseOptions) > 20 {
			theseOptions = theseOptions[:20]
		}

		tripOptions = append(tripOptions, theseOptions...)
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
