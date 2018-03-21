package qpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type QPXFinder struct {
	apiKey string
}

func NewQPXFinder(key string) *QPXFinder {
	return &QPXFinder{
		apiKey: key,
	}
}

type LegSpec struct {
	Origin      string
	Destination string
	Dates       []string // "YYYY-MM-DD"
}

type qpxSpec struct {
	Origin      string
	Destination string
	Date        string // "YYYY-MM-DD"
}

type Leg struct {
	Price    int // In pennies USD
	Segments []Segment
	//	fake     string
}

type Segment struct {
	Airlines      string
	FlightNumber  string
	ArrivalTime   string
	DepartureTime string
	Origin        string
	Destination   string
}

func (l LegSpec) Hash() string {
	return strings.Join(append(l.Dates, l.Origin, l.Destination), "$")
}

func (f *QPXFinder) Find(spec LegSpec) ([]Leg, error) {
	ans := []Leg{}
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

func (f *QPXFinder) findOneDate(spec qpxSpec) ([]Leg, error) {

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

func interpretResp(resp qpxResponse) ([]Leg, error) {
	ans := make([]Leg, len(resp.Trips.TripOption))
	for i, option := range resp.Trips.TripOption {
		segmentsFromTrip := make([]Segment, 0)

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

					segmentsFromTrip = append(segmentsFromTrip, Segment{
						Airlines:      flight.Carrier,
						FlightNumber:  flight.Number,
						ArrivalTime:   leg.ArrivalTime,
						DepartureTime: leg.DepartureTime,
						Origin:        leg.Origin,
						Destination:   leg.Destination,
					})

				}
			}
		}

		//		asdf, _ := json.Marshal(option)
		ans[i] = Leg{
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

	response, err := http.Post("https://www.googleapis.com/qpxExpress/v1/trips/search"+"?key="+f.apiKey, "application/json", bytes.NewBuffer(requestJson))

	if err != nil {
		return resp, err
	}

	if response.StatusCode != 200 {
		return resp, fmt.Errorf("Error! QPX says: %d", response.StatusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)

	err = json.Unmarshal(buf.Bytes(), &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
