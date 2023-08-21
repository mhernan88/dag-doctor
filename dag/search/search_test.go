package search

import (
    "testing"
    "github.com/sirupsen/logrus"
    "github.com/mhernan88/dag-bisect/dag"
)

func TestFindLeaves(t *testing.T) {
    nodes := setupDAG()
    l := logrus.New()
    l.SetLevel(logrus.TraceLevel)

    leaves := findLeaves([]*dag.Node{nodes["node1"]})
    if len(leaves) != 2 {
        t.Errorf("len(leaves) incorrect; want=2, got=%d (%v)", len(leaves), leaves)
    }
}

