package cmd

import (
	"os"
	"testing"

	"github.com/mhernan88/dag-bisect/data"
	"github.com/mhernan88/dag-bisect/pruners"
	"github.com/mhernan88/dag-bisect/splitters"
)

func TestSaveState(t *testing.T) {
	dag, err := data.LoadDAG("../dag.json")
	if err != nil {
		t.Error(err)
	}

	pruner := pruners.NewDefaultPruner()
	splitter := splitters.NewDefaultSplitter()

	ui := NewUI(
		*dag, 
		splitter, 
		pruner,
	)

	err = SaveState("test.json", ui)
	if err != nil {
		t.Error(err)
		return
	}

	newUI, err := LoadState("test.json")
	if err != nil {
		t.Error(err)
		return
	}

	if newUI.DAG.Nodes == nil {
		t.Error("dag nodes was nil!")
	}

	os.Remove("test.json")
}
