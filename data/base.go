package data

import (
	"fmt"

	mapset "github.com/deckarep/golang-set/v2"
)

func NewDAGFromSlice(nodes []Node) (*DAG, error) {
	nodesMap := make(map[string]Node)
	for _, node := range nodes {
		nodesMap[node.Name] = node
	}
	return NewDAGFromMap(nodesMap)
}

func NewDAGFromMap(nodes map[string]Node) (*DAG, error) {
	roots := make(map[string]Node)
	for key, node := range nodes {
		if len(node.Inputs) == 0 {
			roots[key] = node
		}
	}

	dag := DAG{
		Roots: roots,
		Nodes: nodes,
	}

	if !dag.IsAcyclic() {
		return nil, fmt.Errorf("dag is not acyclic")
	}
	return &dag, nil
}

type DAG struct {
	Roots map[string]Node
	Nodes map[string]Node
}

func (d *DAG) Insert(node Node) {
	d.Nodes[node.Name] = node
	if len(node.Next) == 0 {
		d.Roots[node.Name] = node
	}
}

func (d *DAG) Pop(name string) {
	delete(d.Nodes, name)
	_, ok := d.Roots[name]
	if ok {
		delete(d.Roots, name)
	}
}

// Deletes nodes that do not have corresponding node inputs/outputs.
func (d *DAG) reconcileNodesWithInputsAndOutputs() {
	// Find all nodes in inputs/outputs.
	allNames := mapset.NewSet[string]()
	for _, node := range d.Nodes {
		for _, prevNodeName := range node.Prev {
			allNames.Add(prevNodeName)
		}
		for _, nextNodeName := range node.Next {
			allNames.Add(nextNodeName)
		}
	}

	// Find nodes not in inputs/outputs.
	var nodesToDelete []string
	for nodeName := range d.Nodes {
		if !allNames.Contains(nodeName) {
			nodesToDelete = append(nodesToDelete, nodeName)
		}
	}

	// Delete nodes not in inputs/outputs.
	for _, nodeToDelete := range nodesToDelete {
		delete(d.Nodes, nodeToDelete)
	}
}

// Deletes node inputs/outputs that do not have a corresponding node.
func (d *DAG) reconcileInputsAndOutputsWithNodes() {
	// Find all unique node names.
	allNames := mapset.NewSet[string]()
	for nodeName, _ := range d.Nodes {
		allNames.Add(nodeName)
	}

	// Delete inputs/outputs not in nodes.
	for nodeName, node := range d.Nodes {
		var newPrev []string
		for _, prevName := range node.Prev {
			if allNames.Contains(prevName) {
				newPrev = append(newPrev, prevName)
			}
		}
		node.Prev = newPrev

		var newNext []string
		for _, nextName := range node.Next {
			if !allNames.Contains(nextName) {
				newNext = append(newNext, nextName)
			}
		}
		node.Next = newNext

		d.Nodes[nodeName] = node
	}
}

func (d *DAG) Reconcile() {
	d.reconcileNodesWithInputsAndOutputs()
	d.reconcileInputsAndOutputsWithNodes()

	// Delete any lingering roots.
	for rootName := range d.Roots {
		_, ok := d.Nodes[rootName]
		if !ok {
			delete(d.Roots, rootName)
		}
	}
}

func (d *DAG) Slice() ([]string, []Node) {
	var nodeNames []string
	var nodes []Node
	for nodeName, node := range d.Nodes {
		nodeNames = append(nodeNames, nodeName)
		nodes = append(nodes, node)
	}
	return nodeNames, nodes
}

func (d *DAG) Ancestors(node string) map[string]Node {
	nodes := []Node{d.Nodes[node]}
	ancestors := make(map[string]Node)
	for len(nodes) > 0 {
		node := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]
		for _, parent := range node.Prev {
			ancestors[parent] = d.Nodes[parent]
			nodes = append(nodes, d.Nodes[parent])
		}
	}
	return ancestors
}

func (d *DAG) Descendants(start string) map[string]Node {
	nodes := []Node{d.Nodes[start]}
	descendants := make(map[string]Node)
	for len(nodes) > 0 {
		node := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]
		for _, child := range node.Next {
			descendants[child] = d.Nodes[child]
			nodes = append(nodes, d.Nodes[child])
		}
	}
	return descendants
}

// A private method that uses DFS to detect cycles in the sub-graph
func (d *DAG) dfs(nodeName string, visited map[string]bool, stack map[string]bool) bool {
	// Mark the current node as visited and part of the recursion stack
	visited[nodeName] = true
	stack[nodeName] = true

	// Visit all the neighbors
	for _, v := range d.Nodes[nodeName].Outputs {
		// If the node isn't visited yet, then visit it
		if !visited[v] {
			if !d.dfs(v, visited, stack) {
				return false // Cycle detected
			}
		} else if stack[v] {
			// If the node is in the recursion stack, then there's a cycle
			return false
		}
	}

	// Remove the node from the recursion stack
	stack[nodeName] = false
	return true
}

func (d *DAG) IsAcyclic() bool {
	visited := make(map[string]bool)
	stack := make(map[string]bool)

	// Start DFS from each root
	for rootName := range d.Roots {
		if !visited[rootName] {
			if !d.dfs(rootName, visited, stack) {
				return false
			}
		}
	}

	return true
}

type Node struct {
	Name    string   `json:"name"`
	Inputs  []string `json:"inputs"`
	Outputs []string `json:"outputs"`
	Next    []string `json:"next"` // []nodeName
	Prev    []string `json:"prev"` // []nodeName
}
