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

	descendants := dag.Descendants("create_wide_table")
	const expected = 2
	if len(descendants) != expected {
		t.Errorf("descendants: expected %d, got %d", expected, len(descendants))
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

	ancestors := dag.Ancestors("create_wide_table")
	const expected = 3
	if len(ancestors) != expected {
		t.Errorf("ancestors: expected %d, got %d", expected, len(ancestors))
	}
}
