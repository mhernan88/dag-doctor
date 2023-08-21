package search

import (
    "github.com/mhernan88/dag-bisect/dag"
    "github.com/sirupsen/logrus"
    mapset "github.com/deckarep/golang-set/v2"
)

func getNumAncestors(node dag.Node, l *logrus.Logger) int {
    ancestors := mapset.NewSet[string]()

    nodes := []*dag.Node{&node}
    for len(nodes) > 0 {
        nd, nodes := nodes[len(nodes)-1], nodes[:len(nodes)-1]
        for _, parent := range(nd.Prev) {
            ancestors.Add(parent.Name)
            nodes = append(nodes, parent)
        }
    }
    return len(ancestors.ToSlice())
}

func getNumDescendants(node dag.Node, l *logrus.Logger) int {
    descendants := mapset.NewSet[string]()

    nodes := []*dag.Node{&node}
    for len(nodes) > 0 {
        nd, nodes := nodes[len(nodes)-1], nodes[:len(nodes)-1]
        for _, child := range(nd.Next) {
            descendants.Add(child.Name)
            nodes = append(nodes, child)
        }
    }
    return len(descendants.ToSlice())
}
