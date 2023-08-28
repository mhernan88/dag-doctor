package data

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func backwardLinkLeaf(leaf *Node, nodes map[string]*Node) (*Node, error) {
	for _, parentName := range leaf.PrevLabels {
		_, parentExists := leaf.Prev[parentName]
		if parentExists {
			continue
		}

		node, nodeExists := nodes[parentName]
		if !nodeExists {
			return nil, fmt.Errorf("parent '%s' not found in nodes", parentName)
		}
		node, err := backwardLinkLeaf(node, nodes)
		if err != nil {
			return nil, err
		}
		leaf.Prev[parentName] = node
	}
	return leaf, nil
}

func backwardLinkDAG(nodes map[string]*Node) (map[string]*Node, error) {
	leaves := make(map[string]*Node)
	for name, node := range nodes {
		if len(node.NextLabels) == 0 {
			leaves[name] = node
		}
	}

	linkedLeaves := make(map[string]*Node)
	for _, leaf := range leaves {
		linkedLeaf, err := backwardLinkLeaf(leaf, nodes)
		if err != nil {
			return nil, err
		}
		linkedLeaves[linkedLeaf.Name] = linkedLeaf
	}
	return leaves, nil
}

func forwardLinkRoot(root *Node, nodes map[string]*Node) (*Node, error) {
	for _, childName := range root.NextLabels {
		_, childExists := root.Next[childName]
		if childExists {
			continue
		}

		node, nodeExists := nodes[childName]
		if !nodeExists {
			return nil, fmt.Errorf("child '%s' not found in nodes", childName)
		}
		node, err := forwardLinkRoot(node, nodes)
		if err != nil {
			return nil, err
		}
		root.Next[childName] = node
	}
	return root, nil
}

func isAcyclicNode(node *Node) bool {
	if node == nil {
		return true
	}

	if node.State == "visiting" {
		return false
	}

	if node.State == "visited" {
		return true
	}

	node.State = "visiting"

	for _, nextNode := range node.Next {
		if !isAcyclicNode(nextNode) {
			return false
		}
	}

	node.State = "visited"

	return true
}

func isAcyclicGraph(roots map[string]*Node) bool {
	for _, root := range roots {
		if !isAcyclicNode(root) {
			return false
		}
	}
	return true
}

func forwardLinkDAG(nodes map[string]*Node) (map[string]*Node, error) {
	roots := make(map[string]*Node)
	for name, node := range nodes {
		if len(node.PrevLabels) == 0 {
			roots[name] = node
		}
	}

	linkedRoots := make(map[string]*Node)
	for _, root := range roots {
		linkedRoot, err := forwardLinkRoot(root, nodes)
		if err != nil {
			return nil, err
		}
		linkedRoots[linkedRoot.Name] = linkedRoot
	}
	return roots, nil
}

func LoadDAG(filename string) (map[string]*Node, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	dagContainer := make(map[string]map[string]*Node)
	if err := json.Unmarshal(bytes, &dagContainer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	nodes, ok := dagContainer["nodes"]
	if !ok {
		return nil, fmt.Errorf("'nodes' key not found in dat")
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("0 nodes found in catalog")
	}

	for _, node := range nodes {
		node.Next = make(map[string]*Node)
		node.Prev = make(map[string]*Node)
		node.State = "unvisited"
	}

	var dag map[string]*Node

	dag, err = backwardLinkDAG(nodes)
	if err != nil {
		return nil, err
	}

	dag, err = forwardLinkDAG(nodes)
	if err != nil {
		return nil, err
	}

	for name, node := range nodes {
		if len(node.PrevLabels) != len(node.Prev) {
			return nil, fmt.Errorf(
				"node '%s' had %d PrevLabels and %d PrevNodes",
				name,
				len(node.PrevLabels),
				len(node.Prev),
			)
		}

		if len(node.NextLabels) != len(node.Next) {
			return nil, fmt.Errorf(
				"node '%s' had %d NextLabels and %d NextNodes",
				name,
				len(node.NextLabels),
				len(node.Next),
			)
		}
	}

	if !isAcyclicGraph(dag) {
		return nil, fmt.Errorf("nodes are not acyclic")
	}

	return dag, nil
}
