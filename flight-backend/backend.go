package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/frrad/flight-search/flight-backend/qpx"
	"github.com/frrad/flight-search/flight-backend/querytree"
)

func test(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)

	var t querytree.Tree

	err := decoder.Decode(&t)

	if err != nil {
		log.Println("\n")
		log.Println(err)
		log.Println("invalid request")
		return
	}

	defer req.Body.Close()

	log.Println("\n" + t.DispFormat(0))
	log.Println(len(t.Reduce()))
	// for _, res := range t.Resolve() {
	// 	log.Println("asdfasfd")
	// 	log.Println("\n" + res.DispFormat(0))
	// }
	log.Println("test")

}

func main() {
	finder := qpx.NewQPXFinder(os.Getenv("QPXAPIKEY"))
	fmt.Println(finder.Find(
		qpx.LegSpec{
			Origin:      "SFO",
			Destination: "LAX",
			Date:        "2018-04-01",
		}))

	// http.HandleFunc("/backend", test)
	// http.ListenAndServe(":8080", nil)
}
