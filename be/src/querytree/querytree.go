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

func (t *Tree) Resolve() []Tree {

	// True if t is not AND or OR type
	if len(t.Children) == 0 {

		return []Tree{t.copy()}
	}

	down := [][]Tree{}
	combos := 1
	for _, child := range t.Children {
		resolved_child := child.Resolve()
		down = append(down, resolved_child)
		combos *= len(resolved_child)
	}

	answer := []Tree{}
	if t.Type == OrType {

		for _, inside := range down {
			for _, tree := range inside {
				new_one := Tree{
					Type:        t.Type,
					AirportCode: t.AirportCode,
					Modifier:    t.Modifier,
					Children:    []Tree{tree.copy()},
				}
				answer = append(answer, new_one)
			}
		}

		return answer
	}

	// T must be AND type

	for i := 0; i < combos; i++ {

		children := []Tree{}
		index := i
		for _, layer := range down {
			children = append(children, layer[index%len(layer)].copy())
			index /= len(layer)
		}

		new_one := Tree{
			Type:        t.Type,
			AirportCode: t.AirportCode,
			Modifier:    t.Modifier,
			Children:    children,
		}
		answer = append(answer, new_one)
	}

	return answer
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
