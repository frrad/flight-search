package qpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/frrad/flight-search/flight-backend/legfinder"
)

type QPXFinder struct {
	apiKey string
}

func NewQPXFinder(key string) *QPXFinder {
	return &QPXFinder{
		apiKey: key,
	}
}

type qpxSpec struct {
	Origin      string
	Destination string
	Date        string // "YYYY-MM-DD"
}

func (f *QPXFinder) Find(spec legfinder.LegSpec) ([]legfinder.Leg, error) {
	ans := []legfinder.Leg{}
	for _, date := range spec.Dates {
		oneDay, err := f.findOneDate(
			qpxSpec{
				Origin:      spec.Origin,
				Destination: spec.Destination,
				Date:        date,
			})
		if err != nil {
			return ans, err
		}
		ans = append(ans, oneDay...)
	}
	return ans, nil
}

func (f *QPXFinder) findOneDate(spec qpxSpec) ([]legfinder.Leg, error) {

	req := qpxRequest{
		Solutions: 10,
		Passengers: passengerCounts{
			AdultCount: 1,
		},
		Slice: []sliceInput{
			{
				Origin:      spec.Origin,
				Destination: spec.Destination,
				Date:        spec.Date,
			},
		},
	}

	wrapped := requestWrapper{
		Request: req,
	}

	resp, err := f.callQPX(wrapped)
	if err != nil {
		return nil, err
	}

	return interpretResp(resp)
}

func interpretResp(resp qpxResponse) ([]legfinder.Leg, error) {
	ans := make([]legfinder.Leg, len(resp.Trips.TripOption))
	for i, option := range resp.Trips.TripOption {
		segmentsFromTrip := make([]legfinder.Segment, 0)

		costStr := option.SaleTotal
		if costStr[:3] != "USD" {
			return nil, fmt.Errorf("Not in USD: %s", costStr[:3])
		}

		d, err := strconv.ParseFloat(costStr[3:], 64)
		if err != nil {
			return nil, fmt.Errorf("Trouble parsing cost: %s", costStr[3:])
		}

		for _, s := range option.Slice {
			for _, seg := range s.Segment {
				flight := seg.Flight
				for _, leg := range seg.Leg {

					layout := "2006-01-02T15:04-07:00"

					parsedArrival, _ := time.Parse(layout, leg.ArrivalTime)
					parsedDeparture, _ := time.Parse(layout, leg.DepartureTime)

					segmentsFromTrip = append(segmentsFromTrip, legfinder.Segment{
						Airlines:      flight.Carrier,
						FlightNumber:  flight.Number,
						ArrivalTime:   parsedArrival,
						DepartureTime: parsedDeparture,
						Origin:        leg.Origin,
						Destination:   leg.Destination,
					})

				}
			}
		}

		//		asdf, _ := json.Marshal(option)
		ans[i] = legfinder.Leg{
			Price:    int(d * 100),
			Segments: segmentsFromTrip,
			//			fake:     string(asdf),
		}
	}

	return ans, nil
}

func (f *QPXFinder) callQPX(req requestWrapper) (qpxResponse, error) {
	var resp qpxResponse

	requestJson, err := json.Marshal(req)
	if err != nil {
		return resp, err
	}

	log.Println("making QPX request")
	log.Println(requestJson)

	response, err := http.Post("https://www.googleapis.com/qpxExpress/v1/trips/search"+"?key="+f.apiKey, "application/json", bytes.NewBuffer(requestJson))

	if err != nil {
		return resp, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)

	if response.StatusCode != 200 {
		log.Println(buf.String())
		return resp, fmt.Errorf("Error! QPX says: %d", response.StatusCode)
	}

	err = json.Unmarshal(buf.Bytes(), &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
