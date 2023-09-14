package cmd

import (
	"testing"

	"github.com/mhernan88/dag-bisect/data"
	"github.com/mhernan88/dag-bisect/pruners"
	"github.com/mhernan88/dag-bisect/splitters"
	"github.com/sirupsen/logrus"
)

func TestSaveState(t *testing.T) {
	l := logrus.New()
	dag, err := data.LoadDAG("dag.json")
	if err != nil {
		t.Error(err)
	}

	pruner := pruners.NewDefaultPruner()
	splitter := splitters.NewDefaultSplitter()

	ui := NewUI(
		dag, 
		splitter, 
		pruner, 99, l,
	)
}
