package querytree

import (
	"reflect"
	"testing"
)

type resolve_testcase struct {
	InputTree   Tree
	OutputTrees []Tree
	Label       string
}

var resolve_testcases = []resolve_testcase{
	resolve_testcase{
		InputTree: Tree{
			Type:        AndType,
			AirportCode: "",
			Modifier:    modifier{},
			Children:    nil,
		},
		OutputTrees: []Tree{
			Tree{
				Type:        AndType,
				AirportCode: "",
				Modifier:    modifier{},
				Children:    nil,
			},
		},
		Label: "trivial and-tree",
	},
	resolve_testcase{
		InputTree: Tree{
			Type:        AndType,
			AirportCode: "",
			Modifier:    modifier{},
			Children: []Tree{
				Tree{
					Type:        ArriveType,
					AirportCode: "ABC",
					Modifier:    modifier{},
					Children:    nil,
				},
				Tree{
					Type:        ArriveType,
					AirportCode: "XYZ",
					Modifier:    modifier{},
					Children:    nil,
				},
			},
		},
		OutputTrees: []Tree{
			Tree{
				Type:        AndType,
				AirportCode: "",
				Modifier:    modifier{},
				Children: []Tree{
					Tree{
						Type:        ArriveType,
						AirportCode: "ABC",
						Modifier:    modifier{},
						Children:    nil,
					},
					Tree{
						Type:        ArriveType,
						AirportCode: "XYZ",
						Modifier:    modifier{},
						Children:    nil,
					},
				},
			},
		},
		Label: "and-tree, two children",
	},

	resolve_testcase{
		InputTree: Tree{
			Type:        OrType,
			AirportCode: "",
			Modifier:    modifier{},
			Children: []Tree{
				Tree{
					Type:        ArriveType,
					AirportCode: "ABC",
					Modifier:    modifier{},
					Children:    nil,
				},
				Tree{
					Type:        ArriveType,
					AirportCode: "XYZ",
					Modifier:    modifier{},
					Children:    nil,
				},
			},
		},
		OutputTrees: []Tree{
			Tree{
				Type:        OrType,
				AirportCode: "",
				Modifier:    modifier{},
				Children: []Tree{
					Tree{
						Type:        ArriveType,
						AirportCode: "ABC",
						Modifier:    modifier{},
						Children:    nil,
					},
				},
			},
			Tree{
				Type:        OrType,
				AirportCode: "",
				Modifier:    modifier{},
				Children: []Tree{
					Tree{
						Type:        ArriveType,
						AirportCode: "XYZ",
						Modifier:    modifier{},
						Children:    nil,
					},
				},
			},
		},
		Label: "or-tree, two children",
	},
	resolve_testcase{
		InputTree: Tree{
			Type:        AndType,
			AirportCode: "",
			Modifier:    modifier{},
			Children: []Tree{
				Tree{
					Type:     OrType,
					Modifier: modifier{},
					Children: []Tree{
						Tree{
							Type:        DepartType,
							AirportCode: "ABC",
							Modifier:    modifier{},
							Children:    nil,
						},
						Tree{
							Type:        DepartType,
							AirportCode: "DEF",
							Modifier:    modifier{},
							Children:    nil,
						},
					},
				},
				Tree{
					Type:     OrType,
					Modifier: modifier{},
					Children: []Tree{
						Tree{
							Type:        ArriveType,
							AirportCode: "UVW",
							Modifier:    modifier{},
							Children:    nil,
						},
						Tree{
							Type:        ArriveType,
							AirportCode: "XYZ",
							Modifier:    modifier{},
							Children:    nil,
						},
					},
				},
			},
		},
		OutputTrees: []Tree{
			Tree{
				Type:        AndType,
				AirportCode: "",
				Children: []Tree{
					Tree{
						Type:        OrType,
						AirportCode: "",
						Children: []Tree{
							Tree{
								Type:        DepartType,
								AirportCode: "ABC",
								Children:    nil,
								Modifier:    modifier{},
							},
						},
						Modifier: modifier{}},
					Tree{
						Type:        OrType,
						AirportCode: "",
						Children: []Tree{
							Tree{
								Type:        ArriveType,
								AirportCode: "UVW",
								Children:    nil,
								Modifier:    modifier{},
							},
						},
						Modifier: modifier{},
					},
				},
				Modifier: modifier{},
			},
			Tree{
				Type:        AndType,
				AirportCode: "",
				Children: []Tree{
					Tree{
						Type:        OrType,
						AirportCode: "",
						Children: []Tree{
							Tree{
								Type:        DepartType,
								AirportCode: "DEF",
								Children:    nil,
								Modifier:    modifier{},
							},
						},
						Modifier: modifier{}},
					Tree{
						Type:        OrType,
						AirportCode: "",
						Children: []Tree{
							Tree{
								Type:        ArriveType,
								AirportCode: "UVW",
								Children:    nil,
								Modifier:    modifier{},
							},
						},
						Modifier: modifier{},
					},
				},
				Modifier: modifier{},
			},
			Tree{
				Type:        AndType,
				AirportCode: "",
				Children: []Tree{
					Tree{
						Type:        OrType,
						AirportCode: "",
						Children: []Tree{
							Tree{
								Type:        DepartType,
								AirportCode: "ABC",
								Children:    nil,
								Modifier:    modifier{},
							},
						},
						Modifier: modifier{}},
					Tree{
						Type:        OrType,
						AirportCode: "",
						Children: []Tree{
							Tree{
								Type:        ArriveType,
								AirportCode: "XYZ",
								Children:    nil,
								Modifier:    modifier{},
							},
						},
						Modifier: modifier{},
					},
				},
				Modifier: modifier{},
			},
			Tree{
				Type:        AndType,
				AirportCode: "",
				Children: []Tree{
					Tree{
						Type:        OrType,
						AirportCode: "",
						Children: []Tree{
							Tree{
								Type:        DepartType,
								AirportCode: "DEF",
								Children:    nil,
								Modifier:    modifier{},
							},
						},
						Modifier: modifier{}},
					Tree{
						Type:        OrType,
						AirportCode: "",
						Children: []Tree{
							Tree{
								Type:        ArriveType,
								AirportCode: "XYZ",
								Children:    nil,
								Modifier:    modifier{},
							},
						},
						Modifier: modifier{},
					},
				},
				Modifier: modifier{},
			},
		},
		Label: "several levels",
	},
}

func TestResolve(test *testing.T) {
	for _, testcase := range resolve_testcases {
		if !reflect.DeepEqual(testcase.InputTree.Resolve(),
			testcase.OutputTrees) {
			test.Errorf("label: %s\nexpected %v\nFound    %v",
				testcase.Label,
				testcase.OutputTrees,
				testcase.InputTree.Resolve(),
			)
		}
	}
}
