package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type query_tree struct {
	Type node_type

	// Only filled if ARRIVE or DEPART type
	AirportCode string

	Children []query_tree
	Modifier modifier
}

type node_type string

const (
	OrType     node_type = "OR"
	AndType    node_type = "AND"
	ArriveType node_type = "ARRIVE"
	DepartType node_type = "DEPART"
)

type modifier struct {
	PriceAdjustment int32
}

func disp_format(tree query_tree, depth int) string {
	front_pad := strings.Repeat(" ", depth)

	answer := front_pad + fmt.Sprintf("type: %s\n", tree.Type)

	if mod := disp_mod(tree.Modifier); mod != "" {
		answer += front_pad + fmt.Sprintf("modifier: %s\n", mod)
	}

	answer += front_pad + "children:\n"
	for _, child := range tree.Children {
		answer += front_pad + disp_format(child, depth+2)
	}

	return answer
}

// Display modifier
func disp_mod(mod modifier) string {
	if mod.PriceAdjustment != 0 {
		return strconv.Itoa(int(mod.PriceAdjustment))
	}

	return ""
}


// Resolve one query tree to several constrained ones.
func resolve(qt query_tree) []query_tree{
	return []query_tree{}
}



func test(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)

	var t query_tree

	err := decoder.Decode(&t)

	if err != nil {
		log.Println("\n")
		log.Println(err)
		log.Println("invalid request")
		return
	}

	defer req.Body.Close()

	log.Println("\n" + disp_format(t, 0))
}

func main() {
	http.HandleFunc("/backend", test)
	http.ListenAndServe(":8080", nil)
}


