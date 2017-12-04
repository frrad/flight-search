package qpx

import (
	"bytes"
	"encoding/json"
	//	"fmt"
	"net/http"
	"os"
)

//                                 _
//  _ __ ___  __ _ _   _  ___  ___| |_
// | '__/ _ \/ _` | | | |/ _ \/ __| __|
// | | |  __/ (_| | |_| |  __/\__ \ |_
// |_|  \___|\__, |\__,_|\___||___/\__|
//              |_|

// Just a wrapper for qpxRequest
type request struct {
	Request qpxRequest `json:"request"`
}

type qpxRequest struct {
	Passengers passengerCounts `json:"passengers"`
	Slice      []sliceInput    `json:"slice"`
	Solutions  int             `json:"solutions"`
}

type sliceInput struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Date        string `json:"date"`
}

type passengerCounts struct {
	AdultCount int `json:"adultCount"`
}

//  _ __ ___  ___ _ __   ___  _ __  ___  ___
// | '__/ _ \/ __| '_ \ / _ \| '_ \/ __|/ _ \
// | | |  __/\__ \ |_) | (_) | | | \__ \  __/
// |_|  \___||___/ .__/ \___/|_| |_|___/\___|
//               |_|

type qpxResponse struct {
	Kind  string `json:"kind"`
	Trips trips  `json:"trips"`
}

type trips struct {
	Kind       string       `json:"kind"`
	RequestId  string       `json:"requestId"`
	Data       data         `json:"data"`
	TripOption []tripOption `json:"tripOption"`
}

type data struct {
	Kind string `json:"kind"`
}

type tripOption struct {
	Kind      string  `json:"kind"`
	SaleTotal string  `json:"saleTotal"`
	Id        string  `json:"id"`
	Slice     []slice `json:"slice"`
}

type slice struct {
	Kind     string    `json:"kind"`
	Duration int       `json:"duration"`
	Segment  []segment `json:"segment"`
}

type segment struct {
	Kind     string `json:"kind"`
	Duration int    `json:"duration"`
	Flight   flight `json:"flight"`
	Leg      []leg  `json:"leg"`
}

type leg struct {
	ArrivalTime   string `json:"arrivalTime"`
	DepartureTime string `json:"departureTime"`
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
}

type flight struct {
	Carrier string `json:"carrier"`
	Number  string `json:"number"`
}

func testRequest() request {
	var answer request
	answer.Request.Solutions = 10
	answer.Request.Passengers.AdultCount = 1
	answer.Request.Slice = []sliceInput{sliceInput{}}
	answer.Request.Slice[0].Origin = "MCO"
	answer.Request.Slice[0].Destination = "DAY"
	answer.Request.Slice[0].Date = "2017-09-01"
	return answer
}

func getAPIKey() string {
	return os.Getenv("QPXAPIKEY")
}

func CallQPX() string {
	requestJson, err := json.Marshal(testRequest())
	if err != nil {
		return "error"
	}

	response, err := http.Post("https://www.googleapis.com/qpxExpress/v1/trips/search"+"?key="+getAPIKey(), "application/json", bytes.NewBuffer(requestJson))

	if err != nil {
		return "error posting request"
	}

	if response.StatusCode != 200 {
		return "error"
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	// newStr := buf.String()

	var response_coded qpxResponse
	err = json.Unmarshal(buf.Bytes(), &response_coded)

	answer, err := json.Marshal(response_coded)
	return string(answer)
}
