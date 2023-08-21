package search

import (
    "testing"
    "github.com/mhernan88/dag-bisect/dag"
)

func setupDAG() map[string]*dag.Node {
    node5a := dag.Node{Name: "node5a"}
    node4a := dag.Node{Name: "node4a"}
    node3a := dag.Node{Name: "node3a"}
    node3b := dag.Node{Name: "node3b"}
    node2 := dag.Node{Name: "node2"}
    node1 := dag.Node{Name: "node1"}

    node5a.Prev = []*dag.Node{&node4a}

    node4a.Next = []*dag.Node{&node5a}
    node4a.Prev = []*dag.Node{&node3a}

    node3a.Next = []*dag.Node{&node4a}
    node3a.Prev = []*dag.Node{&node2}

    node3b.Prev = []*dag.Node{&node2}

    node2.Next = []*dag.Node{&node3a, &node3b}
    node2.Prev = []*dag.Node{&node1}

    node1.Next = []*dag.Node{&node2}

    nodes := make(map[string]*dag.Node)
    nodes["node1"] = &node1
    nodes["node2"] = &node2
    nodes["node3a"] = &node3a
    nodes["node3b"] = &node3b
    nodes["node4a"] = &node4a
    nodes["node5a"] = &node5a

    return nodes
}


func TestRemoveParent1(t *testing.T) {
    nodes := setupDAG()
    grandchild := nodes["node3a"]
    grandparent := nodes["node1"]

    err := removeParent(grandchild, "node2")
    if err != nil {
        t.Error(err)
    }

    if len(grandchild.Prev) > 0 {
        t.Errorf("want len=0, got %d", len(grandchild.Prev))
    }

    if len(grandparent.Next) > 0 {
        t.Errorf("want len=0, got %d", len(grandparent.Next))
    }
}

func TestRemoveParent2(t *testing.T) {
    nodes := setupDAG()
    grandchild := nodes["node4a"]
    grandparent := nodes["node2"]

    err := removeParent(grandchild, "node3a")
    if err != nil {
        t.Error(err)
    }

    if len(grandchild.Prev) > 0 {
        t.Errorf("want len=0, got %d", len(grandchild.Prev))
    }

    if len(grandparent.Next) != 1 {
        t.Errorf("want len=1, got %d", len(grandparent.Next))
    }

    if grandparent.Next[0].Name != "node3b" {
        t.Errorf("want name='node3b', got '%s'", grandparent.Next[0].Name)
    }
}
