package querytree

import (
	"fmt"
	"strconv"
	"strings"
)

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

type Tree struct {
	Type node_type

	// Only filled if ARRIVE or DEPART type
	AirportCode string

	Children []Tree
	Modifier modifier
}

func (t *Tree) copy() Tree {
	new_one := Tree{
		Type:        t.Type,
		AirportCode: t.AirportCode,
		Modifier:    t.Modifier,
	}
	for _, subtree := range t.Children {
		new_one.Children = append(new_one.Children, subtree.copy())
	}
	return new_one
}

func (tree *Tree) DispFormat(depth int) string {
	front_pad := strings.Repeat(" ", depth)

	answer := front_pad + fmt.Sprintf("type: %s\n", tree.Type)

	if mod := disp_mod(tree.Modifier); mod != "" {
		answer += front_pad + fmt.Sprintf("modifier: %s\n", mod)
	}
	if tree.AirportCode != "" {
		answer += front_pad + fmt.Sprintf("code: %s\n", tree.AirportCode)
	}

	answer += front_pad + "children:\n"
	for _, child := range tree.Children {
		answer += child.DispFormat(depth + 2)
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
