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
	_, leaves := GetNodeDescendants(nodes)
	if len(leaves) != 1 {
		t.Errorf("leaves: expected %d, got %d", 1, len(leaves))
	}
}
