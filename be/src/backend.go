package main

import (
	"encoding/json"
	"log"
	"net/http"
	"querytree"
)

// Resolve one query tree to several constrained ones.
func resolve(qt querytree.Tree) []querytree.Tree {
	if len(qt.Children) == 0 {
		return []querytree.Tree{qt}
	}
	// if qt.Type != OrType {

	// }

	return []querytree.Tree{}
}

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
}

func main() {
	http.HandleFunc("/backend", test)
	http.ListenAndServe(":8080", nil)
}
