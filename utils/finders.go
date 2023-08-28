package utils

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/mhernan88/dag-bisect/data"
)

func FlattenRootNamesMap(nodes map[string]*data.Node) []string {
	var nodesSlice []*data.Node
	for _, node := range nodes {
		nodesSlice = append(nodesSlice, node)
	}
	return FlattenRootNames(nodesSlice)
}
func FlattenRootNames(nodes []*data.Node) []string {
	output := mapset.NewSet[string]()

	var node *data.Node
	for len(nodes) > 0 {
		node = nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]

		if len(node.Prev) == 0 {
			output.Add(node.Name)
		}
	}
	return output.ToSlice()
}

func FlattenLeafNamesMap(nodes map[string]*data.Node) []string {
	var nodesSlice []*data.Node
	for _, node := range nodes {
		nodesSlice = append(nodesSlice, node)
	}
	return FlattenLeafNames(nodesSlice)
}

func FlattenLeafNames(nodes []*data.Node) []string {
	output := mapset.NewSet[string]()

	var node *data.Node
	for len(nodes) > 0 {
		node = nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]

		if len(node.Next) == 0 {
			output.Add(node.Name)
		}

		for _, child := range node.Next {
			nodes = append(nodes, child)
		}
	}
	return output.ToSlice()
}

func FlattenAllNodesToMap(nodes map[string]*data.Node) map[string]*data.Node {
	output := make(map[string]*data.Node)

	var nodesSlice []*data.Node
	for _, node := range nodes {
		nodesSlice = append(nodesSlice, node)
	}

	var node *data.Node
	for len(nodesSlice) > 0 {
		node = nodesSlice[len(nodesSlice)-1]
		nodesSlice = nodesSlice[:len(nodesSlice)-1]

		output[node.Name] = node

		for _, child := range node.Next {
			nodesSlice = append(nodesSlice, child)
		}
	}
	return output
}

func FlattenAllNodesToSlice(nodes map[string]*data.Node) []*data.Node {
	flattenedMap := FlattenAllNodesToMap(nodes)

	var flattenedSlice []*data.Node
	for _, node := range flattenedMap {
		flattenedSlice = append(flattenedSlice, node)
	}
	return flattenedSlice
}

func FlattenAllNodeNames(nodes map[string]*data.Node) []string {
	flattenedNodes := FlattenAllNodesToMap(nodes)

	var flattenedNames []string
	for _, node := range flattenedNodes {
		flattenedNames = append(flattenedNames, node.Name)
	}
	return flattenedNames
}
