package distance

import (
    "testing"
    "github.com/sirupsen/logrus"
)

func setupDAG() map[string]*Node {
    node5a := Node{Name: "node5a"}
    node4a := Node{Name: "node4a"}
    node3a := Node{Name: "node3a"}
    node3b := Node{Name: "node3b"}
    node2 := Node{Name: "node2"}
    node1 := Node{Name: "node1"}

    node5a.Prev = []*Node{&node4a}

    node4a.Next = []*Node{&node5a}
    node4a.Prev = []*Node{&node3a}

    node3a.Next = []*Node{&node4a}
    node3a.Prev = []*Node{&node2}

    node3b.Prev = []*Node{&node2}

    node2.Next = []*Node{&node3a, &node3b}
    node2.Prev = []*Node{&node1}

    node1.Next = []*Node{&node2}

    nodes := make(map[string]*Node)
    nodes["node1"] = &node1
    nodes["node2"] = &node2
    nodes["node3a"] = &node3a
    nodes["node3b"] = &node3b
    nodes["node4a"] = &node4a
    nodes["node5a"] = &node5a

    return nodes
}

func TestDistanceCalculator__DistanceToStart(t *testing.T) {
    l := logrus.New()
    l.SetLevel(logrus.TraceLevel)
    dc := NewDistanceCalculator(l)

    nodes := setupDAG()
    endPtr := nodes["node5a"]

    result, err := dc.DistanceToStart(endPtr, 99)
    if err != nil {
        t.Error(err)
    }
    if result != 4 {
        t.Errorf("want=4, got=%d", result)
    }
}

func TestDistanceCalculator__DistanceToEnd(t *testing.T) {
    l := logrus.New()
    l.SetLevel(logrus.TraceLevel)
    dc := NewDistanceCalculator(l)

    nodes := setupDAG()
    startPtr := nodes["node1"]

    result, err := dc.DistanceToEnd(startPtr, 99)
    if err != nil {
        t.Error(err)
    }
    if result != 4 {
        t.Errorf("want=4, got=%d", result)
    }
}

func TestDistanceCalculator__Midpoints(t *testing.T) {
    l := logrus.New()
    l.SetLevel(logrus.TraceLevel)
    dc := NewDistanceCalculator(l)

    nodes := setupDAG()
    startPtr := nodes["node1"]

    midpoints, err := dc.Midpoints([]*Node{startPtr}, 99)
    if err != nil {
        t.Error(err)
    }
    if len(midpoints) == 0 {
        t.Error("dc found no midpoints")
    }
    if len(midpoints) != 1 {
        t.Errorf("dc found wrong number of midpoints; got=%d, want=1", len(midpoints))
    }
    if midpoints[0] != "node3a" {
        t.Errorf("dc found wrong midpoint; got='%s', want='node3a'", midpoints[0])
    }
}
