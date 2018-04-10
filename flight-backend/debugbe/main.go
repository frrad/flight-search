package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"text/template"

	"github.com/frrad/flight-search/flight-backend/qpx"
	"github.com/frrad/flight-search/flight-backend/querydag"
	"github.com/frrad/flight-search/flight-backend/trip"
)

type pageData struct {
	DAGJSON        string
	UnmarshalError string
	GraphImage     string
	BEResponseHTML string
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	page := `
<html>

<head><title>Debug Frontend</title><head>

{{if .UnmarshalError}}
Error unmarshalling json:
<br>
{{.UnmarshalError}}
<br>
{{end}}

{{if .GraphImage}}
<img src="data:image/png;base64,{{.GraphImage}}">
<br>
{{end}}

{{if .BEResponseHTML}}
{{.BEResponseHTML}}
<br>
{{end}}


<br>
<textarea name="dagjson" cols="80" rows="20" form="03670">{{.DAGJSON}}</textarea>
<form method="post" id="03670">
<input type="submit" value="Submit">
</form>
</html>
`

	defaultDAG := `
{"Nodes":[{"IsAirport":false,"Name":"start","FlightsOut":[{"ToNode":1,"Dates":null},{"ToNode":2,"Dates":null}]},{"IsAirport":true,"Name":"SFO","FlightsOut":[{"ToNode":3,"Dates":["2018-04-15"]}]},{"IsAirport":true,"Name":"OAK","FlightsOut":[{"ToNode":3,"Dates":["2018-04-16"]}]},{"IsAirport":false,"Name":"<flight>","FlightsOut":[{"ToNode":4,"Dates":null},{"ToNode":5,"Dates":null}]},{"IsAirport":true,"Name":"MCO","FlightsOut":[{"ToNode":6,"Dates":null}]},{"IsAirport":true,"Name":"MIA","FlightsOut":[{"ToNode":6,"Dates":null}]},{"IsAirport":false,"Name":"<flight>","FlightsOut":[{"ToNode":7,"Dates":["2018-04-20"]}]},{"IsAirport":true,"Name":"JFK","FlightsOut":[{"ToNode":8,"Dates":null}]},{"IsAirport":false,"Name":"end","FlightsOut":[]}]}
`

	tmpl := template.New("debugfe.html")
	tmpl, err := tmpl.Parse(page)
	if err != nil {
		panic(err)
	}

	jsonString := req.FormValue("dagjson")
	if jsonString == "" {
		jsonString = defaultDAG
	}

	data := pageData{}

	dag := querydag.DAG{}
	err = json.Unmarshal([]byte(jsonString), &dag)

	if err != nil {
		data.DAGJSON = jsonString
		data.UnmarshalError = fmt.Sprintf("%v", err)

		err = tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
		return
	}

	b, _ := json.MarshalIndent(dag, "", "  ")
	data.DAGJSON = string(b)

	data.GraphImage = drawGraph(dag)
	beResp, err := queryBackend(dag)

	if err == nil {
		data.BEResponseHTML = formatResponse(beResp)
	} else {
		data.BEResponseHTML = fmt.Sprintf("Error querying backend:\n<br>\n%+v<br>\n%s", err, beResp)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func queryBackend(dag querydag.DAG) ([]trip.TripOption, error) {
	url := "http://localhost:8080/backend"

	jsonData, err := json.Marshal(dag)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	var beresp []trip.TripOption

	log.Printf("received response. length %d\n", len(buf.String()))

	err = json.Unmarshal(buf.Bytes(), &beresp)
	if err != nil {
		log.Printf("Error: Can't unmarshal response:\n%s", buf.String())
		return nil, err
	}

	return beresp, nil
}

func formatResponse(options []trip.TripOption) string {
	log.Printf("Formatting Response...\n")

	ans := "<table>\n"

	for _, trip := range options {
		ans += fmt.Sprintf("%s\n", formatTrip(trip))
	}

	return ans + "</table>\n"
}

func formatMoney(x int) string {
	dollars, cents := x/100, x%100
	return fmt.Sprintf("$%d.%2d", dollars, cents)
}

func formatTrip(trip trip.TripOption) string {
	ans := fmt.Sprintf("<tr>\n  <td>%d</td>\n  <td>%s</td>\n", trip.Id, formatMoney(trip.Price))
	for _, leg := range trip.Legs {
		ans += fmt.Sprintf("  <td>%s</td>\n", formatLeg(leg))
	}

	return ans + "</tr>"
}

func formatLeg(leg qpx.Leg) string {
	ans := "<table>\n"
	ans += "<tr><th>price</th>"
	for i := 0; i < len(leg.Segments); i++ {
		ans += fmt.Sprintf("<th>seg%d</th>", i)
	}
	ans += "</tr>\n"
	ans += fmt.Sprintf("<tr><td>%d</td>", leg.Price)
	for _, seg := range leg.Segments {
		ans += formatSeg(seg)
	}
	ans += "</tr>\n"
	return ans + "</table>\n"
}

func formatSeg(seg qpx.Segment) string {
	return fmt.Sprintf("<td>%s %s<br>\n%s -> %s<br>\n%v<br>\n%v</td>",
		seg.Airlines,
		seg.FlightNumber,
		seg.Origin,
		seg.Destination,
		seg.DepartureTime,
		seg.ArrivalTime,
	)
}

func main() {
	http.HandleFunc("/debugfe.html", HelloServer)
	log.Fatal(http.ListenAndServe(":9109", nil))
}

func drawGraph(dag querydag.DAG) string {
	filestring := dag.Viz()
	filename := "/tmp/23590.dot"
	ioutil.WriteFile(filename, []byte(filestring), 0644)

	imagePath := "/home/frederick/Downloads/image.png"
	renderCmd := exec.Command("dot", "-Tpng", filename,
		"-o", imagePath)
	renderCmd.Run()

	baseCmd := exec.Command("base64", "-w", "0", imagePath)
	out, _ := baseCmd.Output()

	return string(out)
}
