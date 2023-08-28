package data

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestGetNodeDescendants(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)
	dag, err := LoadDAG("../dag.json")
	if err != nil {
		t.Error(err)
		return
	}

	_, nodes := SliceNodeMap(dag)
	descendants, leaves := GetNodeDescendants(nodes)
	if len(descendants) != 4 {
		t.Errorf("descendants: expected %d, got %d", 4, len(descendants))
	}
	if len(leaves) != 1 {
		t.Errorf("leaves: expected %d, got %d", 1, len(leaves))
	}
}

func TestGetNodeAncestors(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)
	dag, err := LoadDAG("../dag.json")
	if err != nil {
		t.Error(err)
		return
	}

	_, nodes := SliceNodeMap(dag)
	descendants, roots := GetNodeAncestors(nodes)
	if len(descendants) != 0 {
		t.Errorf("ancestors: expected %d, got %d", 0, len(descendants))
	}
	if len(roots) != 0 {
		t.Errorf("leaves: expected %d, got %d", 0, len(roots))
	}
}

func TestUniqueNodes(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)
	dag, err := LoadDAG("../dag.json")
	if err != nil {
		t.Error(err)
		return
	}

	uniqueNodes := UniqueNodes(dag).ToSlice()
	if len(uniqueNodes) != 6 {
		t.Errorf("uniqueNodes: expected %d, got %d", 6, len(uniqueNodes))
	}
}
