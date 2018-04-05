package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"text/template"

	"github.com/frrad/flight-search/flight-backend/querydag"
)

type pageData struct {
	DAGJSON        string
	UnmarshalError string
	GraphImage     string
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

<br>
<textarea name="dagjson" cols="80" rows="20" form="03670">{{.DAGJSON}}</textarea>
<form method="post" id="03670">
<input type="submit" value="Submit">
</form>
</html>
`

	defaultDAG := `
{"Nodes":[{"IsAirport":false,"Name":"start","FlightsOut":[{"ToNode":1,"Dates":null},{"ToNode":2,"Dates":null}]},{"IsAirport":true,"Name":"SFO","FlightsOut":[{"ToNode":3,"Dates":["2018-04-01"]}]},{"IsAirport":true,"Name":"OAK","FlightsOut":[{"ToNode":3,"Dates":["2018-04-01","2018-04-02"]}]},{"IsAirport":false,"Name":"fake","FlightsOut":[{"ToNode":4,"Dates":null},{"ToNode":5,"Dates":null}]},{"IsAirport":true,"Name":"MCO","FlightsOut":[{"ToNode":6,"Dates":null}]},{"IsAirport":true,"Name":"MIA","FlightsOut":[{"ToNode":6,"Dates":null}]},{"IsAirport":false,"Name":"end","FlightsOut":null}]}`

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

	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
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
