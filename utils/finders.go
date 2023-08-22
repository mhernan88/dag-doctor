package utils

import (
    "github.com/mhernan88/dag-bisect/data"
    mapset "github.com/deckarep/golang-set/v2"
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

        for _, child := range(node.Next) {
            nodes = append(nodes, child)
        }
    }
    return output.ToSlice()
}

func FlattenAllNodes(nodes map[string]*data.Node) []string {
    output := mapset.NewSet[string]()

    var nodesSlice []*data.Node
    for _, node := range nodes {
        nodesSlice = append(nodesSlice, node)
    }

    var node *data.Node
    for len(nodesSlice) > 0 {
        node = nodesSlice[len(nodesSlice)-1]
        nodesSlice = nodesSlice[:len(nodesSlice)-1]
        output.Add(node.Name)

        for _, child := range(node.Next) {
            nodesSlice = append(nodesSlice, child)
        }
    }
    return output.ToSlice()
}
