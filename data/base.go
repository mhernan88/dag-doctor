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
	if len(nodes) == 0 {
		return nil, fmt.Errorf("dag constructor received no nodes")
	}
	roots := make(map[string]Node)
	for key, node := range nodes {
		if (len(node.Prev) == 0) || (node.Prev == nil) {
			roots[key] = node
		}
	}

	if len(roots) == 0 {
		return nil, fmt.Errorf("failed to find dag roots")
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

func SliceMap(nodes map[string]Node) ([]string, []Node) {
	var outputNodeNames []string
	var outputNodes []Node
	for nodeName, node := range nodes {
		outputNodeNames = append(outputNodeNames, nodeName)
		outputNodes = append(outputNodes, node)
	}
	return outputNodeNames, outputNodes
}

func SliceMapKeys(nodes map[string]Node) []string {
	keys, _ := SliceMap(nodes)
	return keys
}

func SliceMapValues(nodes map[string]Node) []Node {
	_, values := SliceMap(nodes)
	return values
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

func (d *DAG) CompileInputsAndOutputs() mapset.Set[string] {
	allNames := mapset.NewSet[string]()
	for _, node := range d.Nodes {
		for _, prevNodeName := range node.Prev {
			allNames.Add(prevNodeName)
		}
		for _, nextNodeName := range node.Next {
			allNames.Add(nextNodeName)
		}
	}
	return allNames
}

// Removes references to name from Prev/Next. Returns number of nodes changed.
func (d *DAG) Unlink(name string) int {
    // Delete a node from Next
    fmt.Printf("checking children for node %s\n", name)
    unlinkedNodes := mapset.NewSet[string]()
    for nodeName, node := range d.Nodes {
        var newChildren []string
        for _, childName := range node.Next {
            if childName != name {
                newChildren = append(newChildren, childName)
            } else {
                fmt.Println("killng child!")
                unlinkedNodes.Add(nodeName)
            }
        }
        node.Next = newChildren
        d.Nodes[nodeName] = node
    }

    // Delete a node from Prev
    fmt.Printf("checking parents for node %s\n", name)
    for nodeName, node := range d.Nodes {
        var newParents []string
        for _, parentName := range node.Prev {
            if parentName != name {
                newParents = append(newParents, parentName)
            } else {
                fmt.Println("killng parent!")
                unlinkedNodes.Add(nodeName)
            }
        }
        node.Prev = newParents
        d.Nodes[nodeName] = node
    }
    return len(unlinkedNodes.ToSlice())
}

// Deletes nodes that do not have corresponding node inputs/outputs.
func (d *DAG) reconcileNodesWithInputsAndOutputs() int {
	// Find all nodes in inputs/outputs.
	allNames := d.CompileInputsAndOutputs()

	// Find nodes not in inputs/outputs.
    nodesToDelete := mapset.NewSet[string]()
	for nodeName := range d.Nodes {
		if !allNames.Contains(nodeName) {
            nodesToDelete.Add(nodeName)
		}
	}

	// Delete nodes not in inputs/outputs.
	for _, nodeToDelete := range nodesToDelete.ToSlice() {
		delete(d.Nodes, nodeToDelete)
	}
    return len(nodesToDelete.ToSlice())
}

// Deletes node inputs/outputs that do not have a corresponding node.
func (d *DAG) reconcileInputsAndOutputsWithNodes() int {
	// Find all unique node names.
	allNames := mapset.NewSet[string]()
	for nodeName := range d.Nodes {
		allNames.Add(nodeName)
	}

	// Delete inputs/outputs not in nodes.
    deletedInputsAndOutputs := mapset.NewSet[string]()
	for nodeName, node := range d.Nodes {
		var newPrev []string
		for _, prevName := range node.Prev {
			if allNames.Contains(prevName) {
				newPrev = append(newPrev, prevName)
			} else {
                deletedInputsAndOutputs.Add(nodeName)
            }
		}
		node.Prev = newPrev

		var newNext []string
		for _, nextName := range node.Next {
			if allNames.Contains(nextName) {
				newNext = append(newNext, nextName)
			} else {
                deletedInputsAndOutputs.Add(nodeName)
            }
		}
		node.Next = newNext

		d.Nodes[nodeName] = node
	}
    return len(deletedInputsAndOutputs.ToSlice())
}

func (d *DAG) Reconcile() {
	// TODO: Shouldn't have to run these each twice...
	d.reconcileNodesWithInputsAndOutputs()
	d.reconcileInputsAndOutputsWithNodes()
	// d.reconcileNodesWithInputsAndOutputs()
	// d.reconcileInputsAndOutputsWithNodes()

	// Delete any lingering roots.
	for rootName := range d.Roots {
		_, ok := d.Nodes[rootName]
		if !ok {
			delete(d.Roots, rootName)
		}
	}
	isReconciled, labelsAndNotNodes, nodesAndNotLabels := d.IsReconciled()
	if !isReconciled {
		fmt.Printf("labels not nodes: %v; nodes not labels: %v", labelsAndNotNodes, nodesAndNotLabels)
		panic("recon failed!!!")
	}
}

// Returns whether dag is reconciled, labels (inputs/outputs) not in nodes, and nodes not in labels (inputs/outputs).
func (d *DAG) IsReconciled() (bool, []string, []string) {
	labels := d.CompileInputsAndOutputs()
	names := mapset.NewSetFromMapKeys[string](d.Nodes)

	labelsAndNotNodes := labels.Difference(names).ToSlice()
	NodesAndNotLabels := names.Difference(labels).ToSlice()
	return len(labels.SymmetricDifference(names).ToSlice()) == 0, labelsAndNotNodes, NodesAndNotLabels
}

func (d *DAG) Slice() ([]string, []Node) {
	return SliceMap(d.Nodes)
}

func (d *DAG) SliceKeys() []string {
	keys, _ := d.Slice()
	return keys
}

func (d *DAG) SliceValues() []Node {
	_, values := d.Slice()
	return values
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
