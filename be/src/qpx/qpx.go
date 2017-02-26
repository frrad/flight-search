package qpx

import (
	"bytes"
	"encoding/json"
	//	"fmt"
	"net/http"
	"os"
)

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
}

func testRequest() request {
	var answer request
	answer.Request.Solutions = 10
	answer.Request.Passengers.AdultCount = 1
	answer.Request.Slice = []sliceInput{sliceInput{}}
	answer.Request.Slice[0].Origin = "LAX"
	answer.Request.Slice[0].Destination = "SFO"
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
	newStr := buf.String()

	return newStr
}
