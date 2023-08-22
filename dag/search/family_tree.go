package search

import (
    "github.com/mhernan88/dag-bisect/dag"
    "github.com/sirupsen/logrus"
    mapset "github.com/deckarep/golang-set/v2"
)

func getNumAncestors(node dag.Node, l *logrus.Logger) int {
    ancestors := mapset.NewSet[string]()

    l.Tracef("locating ancestors")
    nodes := []*dag.Node{&node}
    for len(nodes) > 0 {
        nd := nodes[len(nodes)-1]
        nodes = nodes[:len(nodes)-1]
        l.Tracef("popping node %s", nd.Name)
        for _, parent := range(nd.Prev) {
            l.Tracef("|---> looking at parent %s", parent.Name)
            ancestors.Add(parent.Name)
            nodes = append(nodes, parent)
        }
    }
    l.Tracef("found %d ancestors", len(ancestors.ToSlice()))
    return len(ancestors.ToSlice())
}

func getNumDescendants(node dag.Node, l *logrus.Logger) int {
    descendants := mapset.NewSet[string]()

    l.Tracef("locating descendants")
    nodes := []*dag.Node{&node}
    for len(nodes) > 0 {
        nd := nodes[len(nodes)-1]
        nodes = nodes[:len(nodes)-1]
        l.Tracef("popping node %s", nd.Name)
        for _, child := range(nd.Next) {
            l.Tracef("|---> looking at child %s", child.Name)
            descendants.Add(child.Name)
            nodes = append(nodes, child)
        }
    }
    l.Tracef("found %d descendants", len(descendants.ToSlice()))
    return len(descendants.ToSlice())
}
