
package models

import (
	"os"
	"testing"
)

func TestSaveState(t *testing.T) {
	dag, err := LoadDAG("../dag.json")
	if err != nil {
		t.Error(err)
	}

	state := NewState(
		*dag, 
	)

	err = SaveState("test.json", state)
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
