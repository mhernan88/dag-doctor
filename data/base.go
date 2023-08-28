package data

import (
	mapset "github.com/deckarep/golang-set/v2"
)

type Dataset struct {
	Filename string `json:"filename"`
}

type Node struct {
	Name       string           `json:"name"`
	Inputs     []string         `json:"inputs"`
	Outputs    []string         `json:"outputs"`
	NextLabels []string         `json:"next"` // []nodeName
	PrevLabels []string         `json:"prev"` // []nodeName
	Next       map[string]*Node `json:"-"`    // map[nodeName]nodePointer
	Prev       map[string]*Node `json:"-"`    // map[nodeName]nodePointer
	State      string           `json:"-"`
}

func NewNodeMap(nodeslice []*Node) map[string]*Node {
	dag := make(map[string]*Node)
	for _, node := range nodeslice {
		dag[node.Name] = node
	}
	return dag
}

func SliceNodeMap(nodemap map[string]*Node) ([]string, []*Node) {
	var names []string
	var nodes []*Node
	for name, node := range nodemap {
		names = append(names, name)
		nodes = append(nodes, node)
	}
	return names, nodes
}

// Returns a) a map of all descendant nodes and b) a list of downstream leaf nodes.
func GetNodeDescendants(sources []*Node) (
	map[string]*Node,
	map[string]*Node,
) {
	var nodes []*Node
	for _, source := range sources {
		nodes = append(nodes, source)
	}

	descendants := make(map[string]*Node)
	leaves := make(map[string]*Node)

	for len(nodes) > 0 {
		node := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]

		for childName, child := range node.Next {
			descendants[childName] = child
			if len(child.Next) == 0 {
				leaves[childName] = child
			}
			nodes = append(nodes, child)
		}
	}
	return descendants, leaves
}

// Returns a) a map of all ancestor nodes and b) a list of upstream root nodes.
func GetNodeAncestors(sources []*Node) (
	map[string]*Node,
	map[string]*Node,
) {
	var nodes []*Node
	for _, source := range sources {
		nodes = append(nodes, source)
	}

	ancestors := make(map[string]*Node)
	roots := make(map[string]*Node)

	for len(nodes) > 0 {
		node := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]

		for parentName, parent := range node.Prev {
			ancestors[parentName] = parent
			if len(parent.Prev) == 0 {
				roots[parentName] = parent
			}
			nodes = append(nodes, parent)
		}
	}
	return ancestors, roots
}

func UniqueNodes(roots map[string]*Node) mapset.Set[string] {
	var nodes []*Node
	for _, root := range roots {
		nodes = append(nodes, root)
	}

	nodeNames := mapset.NewSet[string]()

	for len(nodes) > 0 {
		node := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]

		nodeNames.Add(node.Name)
		for _, parent := range node.Prev {
			nodes = append(nodes, parent)
		}
	}
	return nodeNames
}

func LoadNodesIntoSet(nodes []string) mapset.Set[string] {
	uniqueNodes := mapset.NewSet[string]()
	for _, node := range nodes {
		uniqueNodes.Add(node)
	}
	return uniqueNodes
}
