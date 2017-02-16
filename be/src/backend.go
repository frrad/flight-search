package main

import (
	"encoding/json"
	// "fmt"
	"log"
	"net/http"
	"querytree"
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
	log.Println(len(t.Resolve()))
	// for _, res := range t.Resolve() {
	// 	log.Println("asdfasfd")
	// 	log.Println("\n" + res.DispFormat(0))
	// }
	log.Println("test")

}

func main() {
	http.HandleFunc("/backend", test)
	http.ListenAndServe(":8080", nil)
}
