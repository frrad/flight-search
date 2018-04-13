package amadeus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/frrad/flight-search/flight-backend/legfinder"
)

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
		price := result.Fare.TotalPrice // dollars.cents string
	}
	return ans, nil
}

func (a *AmadeusLegFinder) callAPI(origin, destination, date string) (*amadeusResponse, error) {
	log.Println(origin, destination, date)

	urlTemplate := "https://api.sandbox.amadeus.com/v1.2/flights/low-fare-search?apikey=%s&origin=%s&destination=%s&departure_date=%s"
	url := fmt.Sprintf(urlTemplate, a.apiKey, origin, destination, date)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("error calling Amadeus", err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Println("Amadeus returned code", resp.StatusCode)
		return nil, fmt.Errorf("Amadeus returned code %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
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

type amadeusResult struct {
	Fare        fareInfo `json:"fare"`
	Itineraries []struct {
		Outbound struct {
			Flights []struct {
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
			} `json:"flights"`
		} `json:"outbound"`
	} `json:"itineraries"`
}
