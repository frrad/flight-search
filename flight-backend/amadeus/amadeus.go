package amadeus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/frrad/flight-search/flight-backend/legfinder"
)

var priceRE = regexp.MustCompile("^([0-9]*)\\.([0-9]*)$")

const timeLayout = "2006-01-02T15:04"

type AmadeusLegFinder struct {
	apiKey string
}

func NewAmadeusFinder(key string) *AmadeusLegFinder {
	return &AmadeusLegFinder{
		apiKey: key,
	}
}

func (a *AmadeusLegFinder) Find(spec legfinder.LegSpec) ([]legfinder.Leg, error) {
	ans := []legfinder.Leg{}

	for _, date := range spec.Dates {
		response, err := a.callAPI(spec.Origin, spec.Destination, date)
		if err != nil {
			if _, ok := err.(noResultsError); ok {
				return []legfinder.Leg{}, nil
			}

			return nil, err
		}
		if response.Currency != "USD" {
			log.Println("non-us currency detected", response.Currency)
			return nil, fmt.Errorf("non-USD == wat do")
		}

		legs, err := legsFromAmadeusResults(response.Results)
		if err != nil {
			return nil, err
		}
		ans = append(ans, legs...)

	}

	return ans, nil
}

func legsFromAmadeusResults(results []amadeusResult) ([]legfinder.Leg, error) {
	// todo: finish funciton

	ans := []legfinder.Leg{}
	for _, result := range results {
		priceStr := result.Fare.TotalPrice
		price, err := extractPrice(priceStr)
		if err != nil {
			return nil, err
		}

		newLegs, err := legsFromItins(price, result.Itineraries)
		if err != nil {
			return nil, err
		}

		ans = append(ans, newLegs...)

	}
	return ans, nil
}

func legsFromItins(price int, itins []Itinerary) ([]legfinder.Leg, error) {
	legs := []legfinder.Leg{}

	for _, it := range itins {

		segs := []legfinder.Segment{}

		for _, flight := range it.Outbound.Flights {
			arrive, err := time.Parse(timeLayout, flight.ArrivesAt)
			if err != nil {
				return nil, err
			}
			depart, err := time.Parse(timeLayout, flight.DepartsAt)
			if err != nil {
				return nil, err
			}

			segs = append(segs, legfinder.Segment{
				Airlines:      flight.OperatingAirline,
				FlightNumber:  flight.FlightNumber,
				ArrivalTime:   arrive,
				DepartureTime: depart,
				Origin:        flight.Origin.Airport,
				Destination:   flight.Destination.Airport,
			})
		}

		legs = append(legs, legfinder.Leg{
			Price:    price,
			Segments: segs,
		})

	}

	return legs, nil
}

func extractPrice(priceStr string) (int, error) {
	extracted := priceRE.FindStringSubmatch(priceStr)

	if len(extracted) != 3 {
		return 0, fmt.Errorf("Trouble parsing price, %s", priceStr)
	}

	dollars, cents := extracted[1], extracted[2]

	d, err := strconv.Atoi(dollars)
	if err != nil {
		return 0, err
	}
	c, err := strconv.Atoi(cents)
	if err != nil {
		return 0, err
	}

	return d*100 + c, nil
}

func (a *AmadeusLegFinder) callAPI(origin, destination, date string) (*amadeusResponse, error) {

	urlTemplate := "https://api.sandbox.amadeus.com/v1.2/flights/low-fare-search?apikey=%s&origin=%s&destination=%s&departure_date=%s"
	url := fmt.Sprintf(urlTemplate, a.apiKey, origin, destination, date)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("error calling Amadeus", err)
		return nil, err
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Println("Amadeus returned code", resp.StatusCode)
		log.Println(string(contents))

		errInfo := amadeusError{}
		err := json.Unmarshal(contents, &errInfo)
		if err != nil {
			log.Println("error unmarshalling error", string(contents))
			return nil, err
		}

		if errInfo.Message == "No result found." {
			return nil, noResultsError{}
		}

		return nil, fmt.Errorf("Amadeus returned code %d.\n%+v", resp.StatusCode, errInfo)
	}

	if err != nil {
		log.Println("error reading response body", err)
		return nil, err
	}

	var result amadeusResponse
	err = json.Unmarshal(contents, &result)
	if err != nil {
		log.Println("error unmarshalling", string(contents))
		return nil, err
	}

	return &result, nil
}

type noResultsError struct{}

func (x noResultsError) Error() string {
	return ""
}

type amadeusError struct {
	Status   int    `json:"status"`
	Message  string `json:"message"`
	MoreInfo string `json:"more_info"`
}

type amadeusResponse struct {
	Currency string          `json:"currency"`
	Results  []amadeusResult `json:"results"`
}

type fareInfo struct {
	PricePerAdult struct {
		Tax       string `json:"tax"`
		TotalFare string `json:"total_fare"`
	} `json:"price_per_adult"`
	Restrictions struct {
		ChangePenalties bool `json:"change_penalties"`
		Refundable      bool `json:"refundable"`
	} `json:"restrictions"`
	TotalPrice string `json:"total_price"`
}

type Itinerary struct {
	Outbound struct {
		Flights []Flight `json:"flights"`
	} `json:"outbound"`
}

type Flight struct {
	Aircraft    string `json:"aircraft"`
	ArrivesAt   string `json:"arrives_at"`
	BookingInfo struct {
		BookingCode    string `json:"booking_code"`
		SeatsRemaining int64  `json:"seats_remaining"`
		TravelClass    string `json:"travel_class"`
	} `json:"booking_info"`
	DepartsAt   string `json:"departs_at"`
	Destination struct {
		Airport  string `json:"airport"`
		Terminal string `json:"terminal"`
	} `json:"destination"`
	FlightNumber     string `json:"flight_number"`
	MarketingAirline string `json:"marketing_airline"`
	OperatingAirline string `json:"operating_airline"`
	Origin           struct {
		Airport  string `json:"airport"`
		Terminal string `json:"terminal"`
	} `json:"origin"`
}

type amadeusResult struct {
	Fare        fareInfo    `json:"fare"`
	Itineraries []Itinerary `json:"itineraries"`
}
