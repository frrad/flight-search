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
					Type:        AirportType,
					AirportCode: "ABC",
					Modifier:    modifier{},
					Children:    nil,
				},
				Tree{
					Type:        AirportType,
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
						Type:        AirportType,
						AirportCode: "ABC",
						Modifier:    modifier{},
						Children:    nil,
					},
					Tree{
						Type:        AirportType,
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
					Type:        AirportType,
					AirportCode: "ABC",
					Modifier:    modifier{},
					Children:    nil,
				},
				Tree{
					Type:        AirportType,
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
						Type:        AirportType,
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
						Type:        AirportType,
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
							Type:        AirportType,
							AirportCode: "ABC",
							Modifier:    modifier{},
							Children:    nil,
						},
						Tree{
							Type:        AirportType,
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
							Type:        AirportType,
							AirportCode: "UVW",
							Modifier:    modifier{},
							Children:    nil,
						},
						Tree{
							Type:        AirportType,
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
								Type:        AirportType,
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
								Type:        AirportType,
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
								Type:        AirportType,
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
								Type:        AirportType,
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
								Type:        AirportType,
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
								Type:        AirportType,
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
								Type:        AirportType,
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
								Type:        AirportType,
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

type simplify_testcase struct {
	InputTree  Tree
	OutputTree Tree
	Label      string
}

var simplify_testcases = []simplify_testcase{
	simplify_testcase{
		InputTree:  Tree{},
		OutputTree: Tree{},
		Label:      "empty tree",
	},
	simplify_testcase{
		InputTree: Tree{
			Type:        AndType,
			AirportCode: "",
			Children: []Tree{
				Tree{
					Type:        OrType,
					AirportCode: "",
					Children: []Tree{
						Tree{
							Type:        AirportType,
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
							Type:        AirportType,
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
		OutputTree: Tree{
			Type: AndType,
			Children: []Tree{
				Tree{
					Type:        AirportType,
					AirportCode: "ABC",
				},
				Tree{
					Type:        AirportType,
					AirportCode: "XYZ",
				},
			},
		},
		Label: "non-empty tree",
	},
}

func TestSimplify(test *testing.T) {
	for _, testcase := range simplify_testcases {
		if !reflect.DeepEqual(*testcase.InputTree.Simplify(),
			testcase.OutputTree) {
			test.Errorf("label: %s\nexpected %v\nFound    %v",
				testcase.Label,
				testcase.OutputTree,
				testcase.InputTree.Simplify(),
			)
		}
	}
}
