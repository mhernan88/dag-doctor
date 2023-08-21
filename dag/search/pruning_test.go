package search

import (
    "testing"
    "github.com/sirupsen/logrus"
    "github.com/mhernan88/dag-bisect/dag"
)

func setupDAG() map[string]*dag.Node {
    node5a := dag.Node{
        Name: "node5a",
        Prev: make(map[string]*dag.Node),
        Next: make(map[string]*dag.Node)}
    node4a := dag.Node{
        Name: "node4a",
        Prev: make(map[string]*dag.Node),
        Next: make(map[string]*dag.Node)}
    node3a := dag.Node{
        Name: "node3a",
        Prev: make(map[string]*dag.Node),
        Next: make(map[string]*dag.Node)}
    node3b := dag.Node{
        Name: "node3b",
        Prev: make(map[string]*dag.Node),
        Next: make(map[string]*dag.Node)}
    node2 := dag.Node{
        Name: "node2",
        Prev: make(map[string]*dag.Node),
        Next: make(map[string]*dag.Node)}
    node1 := dag.Node{
        Name: "node1",
        Prev: make(map[string]*dag.Node),
        Next: make(map[string]*dag.Node)}

    node5a.Prev[node4a.Name] = &node4a

    node4a.Next[node5a.Name] = &node5a
    node4a.Prev[node3a.Name] = &node3a

    node3a.Next[node4a.Name] = &node4a
    node3a.Prev[node2.Name] = &node2

    node3b.Prev[node2.Name] = &node2

    node2.Next[node3a.Name] = &node3a
    node2.Next[node3b.Name] = &node3b
    node2.Prev[node1.Name] = &node1

    node1.Next[node2.Name] = &node2

    nodes := make(map[string]*dag.Node)
    nodes["node1"] = &node1
    nodes["node2"] = &node2
    nodes["node3a"] = &node3a
    nodes["node3b"] = &node3b
    nodes["node4a"] = &node4a
    nodes["node5a"] = &node5a

    return nodes
}

func TestPruneBefore1(t *testing.T) {
    nodes := setupDAG()
    roots := make(map[string]*dag.Node)
    roots["node1"] = nodes["node1"]
    l := logrus.New()
    l.SetLevel(logrus.TraceLevel)

    pruner := NewDefaultPruner(99, l)

    prunedNodes, err := pruner.PruneBefore(nodes["node2"], roots)
    if err != nil {
        t.Error(err)
    }
    if len(prunedNodes) != 1 {
        t.Errorf("len(prunedNodes) incorrect; want=1, got=%d (%v)", len(prunedNodes), prunedNodes)
    }
    t.Logf("pruned nodes: %v", prunedNodes)

    if len(nodes["node2"].Prev) != 0 {
        t.Errorf("expected node2 to have 0 parents, but it had %d", len(nodes["node2"].Prev))
    }
}

func TestPruneBefore2(t *testing.T) {
    nodes := setupDAG()
    roots := make(map[string]*dag.Node)
    roots["node1"] = nodes["node1"]
    l := logrus.New()
    l.SetLevel(logrus.TraceLevel)

    pruner := NewDefaultPruner(99, l)

    prunedNodes, err := pruner.PruneBefore(nodes["node5a"], roots)
    if err != nil {
        t.Error(err)
    }
    if len(prunedNodes) != 2 {
        t.Errorf("len(prunedNodes) incorrect; want=2, got=%d (%v)", len(prunedNodes), prunedNodes)
    }
}

func TestPruneAfter1(t *testing.T) {
    nodes := setupDAG()
    l := logrus.New()
    l.SetLevel(logrus.TraceLevel)

    pruner := NewDefaultPruner(99, l)

    prunedNodes, err := pruner.PruneAfter(nodes["node2"])
    if err != nil {
        t.Error(err)
    }
    if len(prunedNodes) != 4 {
        t.Errorf("len(prunedNodes) incorrect; want=4, got=%d (%v)", len(prunedNodes), prunedNodes)
    }
}

func TestPruneAfter(t *testing.T) {
    nodes := setupDAG()
    l := logrus.New()
    l.SetLevel(logrus.TraceLevel)

    pruner := NewDefaultPruner(99, l)

    prunedNodes, err := pruner.PruneAfter(nodes["node3a"])
    if err != nil {
        t.Error(err)
    }
    if len(prunedNodes) != 2 {
        t.Errorf("len(prunedNodes) incorrect; want=2, got=%d (%v)", len(prunedNodes), prunedNodes)
    }
}
