package querytree

import (
	"fmt"
	//	"log"
	"strconv"
	"strings"
)

type node_type string

const (
	OrType      node_type = "OR"
	AndType     node_type = "AND"
	AirportType node_type = "AIRPORT"
	TripType    node_type = "TRIP"
)

type modifier struct {
	PriceAdjustment int32
}

type Tree struct {
	Type     node_type
	Modifier modifier

	// Only filled if AIRPORT type
	AirportCode string

	// Filled if and only if TRIP type
	Depart *Tree
	Arrive *Tree

	// Filled if and only if AND or OR type
	Children []Tree
}

func (t *Tree) Simplify() *Tree {

	// The leaf case
	if t.Type == AirportType {
		return t
	}

	if len(t.Children) == 1 {
		// must be AND or OR type
		child := t.Children[0].Simplify()
		child.CombineModifier(t.Modifier)

		return child
	}

	if len(t.Children) > 1 {
		// must be AND or OR type
		for i, child := range t.Children {
			t.Children[i] = *child.Simplify()
		}
		return t
	}

	// Must be TRIP type
	if t.Arrive != nil {
		t.Arrive = t.Arrive.Simplify()
	}
	if t.Depart != nil {
		t.Depart = t.Depart.Simplify()
	}
	return t

}

// Adds given modifier to t
func (t *Tree) CombineModifier(mod modifier) {
	t.Modifier.PriceAdjustment += mod.PriceAdjustment
}

func (t *Tree) Resolve() []Tree {
	// Invalid type.
	if t.Type == "" || t.Type == AndType && len(t.Children) == 0 {
		return []Tree{t.copy()}
	}

	// AIRPORT case
	if t.Type == AirportType {
		return []Tree{t.copy()}
	}

	//  TRIP case here. In some ways this is a less general AND...
	if t.Type == TripType {
		// log.Println("TRIP TYPE")

		arrivals := t.Arrive.Resolve()
		departures := t.Depart.Resolve()

		answer := []Tree{}
		for _, arrive := range arrivals {
			// log.Println("I = ", i)
			for _, depart := range departures {
				// log.Println("J = ", j)
				entry := t.copyTop()

				arriver := arrive.copy()
				entry.Arrive = &arriver

				departer := depart.copy()
				entry.Depart = &departer

				answer = append(answer, entry)
			}
		}
		return answer
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

func (t *Tree) Reduce() []Tree {
	resolved := t.Simplify().Resolve()
	for i, tree := range resolved {
		resolved[i] = *tree.Simplify()
	}
	return resolved
}

// Copy top level values
func (t *Tree) copyTop() Tree {
	new_one := Tree{
		Type:        t.Type,
		AirportCode: t.AirportCode,
		Modifier:    t.Modifier,
	}
	return new_one
}

// Deep copy
func (t *Tree) copy() Tree {
	new_one := t.copyTop()
	for _, subtree := range t.Children {
		new_one.Children = append(new_one.Children, subtree.copy())
	}

	if t.Depart != nil {
		depart := t.Depart.copy()
		new_one.Depart = &depart
	}
	if t.Arrive != nil {
		arrive := t.Arrive.copy()
		new_one.Arrive = &arrive
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

	if tree.Arrive != nil {
		answer += front_pad + "Arrive:\n"
		answer += tree.Arrive.DispFormat(depth + 2)
	}
	if tree.Depart != nil {
		answer += front_pad + "Depart:\n"
		answer += tree.Depart.DispFormat(depth + 2)
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
