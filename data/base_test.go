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

func TestCompileInputsAndOutputs(t *testing.T) {
    l := logrus.New()
    l.SetLevel(logrus.TraceLevel)
    dag, err := LoadDAG("../dag.json")
    if err != nil {
        t.Error(err)
        return
    }

    inputsAndOutputs := dag.CompileInputsAndOutputs()
    if len(inputsAndOutputs.ToSlice()) != 6 {
        t.Errorf(
            "inputsAndOutputs len: expected 7, got %d", 
            len(inputsAndOutputs.ToSlice()),
        )
    }
}

func TestUnlink1(t *testing.T) {
    l := logrus.New()
    l.SetLevel(logrus.TraceLevel)
    dag, err := LoadDAG("../dag.json")
    if err != nil {
        t.Error(err)
        return
    }

    n := dag.Unlink("")
    if n != 0 {
        t.Errorf("TestUnlink n: expected 0, got %d", n)
    }
}

func TestUnlink2(t *testing.T) {
    l := logrus.New()
    l.SetLevel(logrus.TraceLevel)
    dag, err := LoadDAG("../dag.json")
    if err != nil {
        t.Error(err)
        return
    }

    n := dag.Unlink("create_final_model")
    if n != 2 {
        t.Errorf("TestUnlink n: expected 2, got %d", n)
    }
}

func TestReconNodesWithInputsAndOutputs1(t *testing.T) {
    l := logrus.New()
    l.SetLevel(logrus.TraceLevel)
    dag, err := LoadDAG("../dag.json")
    if err != nil {
        t.Error(err)
        return
    }

    n := dag.reconcileNodesWithInputsAndOutputs()
    if n != 0 {
        t.Errorf("nodesWithInputsAndOuputs n: expected 0, got %d", n)
    }
}

func TestReconNodesWithInputsAndOutputs2(t *testing.T) {
    l := logrus.New()
    l.SetLevel(logrus.TraceLevel)
    dag, err := LoadDAG("../dag.json")
    if err != nil {
        t.Error(err)
        return
    }

    const nodeToDelete = "create_final_model"
    n := dag.Unlink(nodeToDelete)
    if n != 2 {
        t.Errorf("TestUnlink n: expected 2, got %d", n)
    }

    n = dag.reconcileNodesWithInputsAndOutputs()
    if n != 1 {
        t.Errorf("nodesWithInputsAndOuputs n: expected 1, got %d", n)
    }
}

func TestReconInputsAndOutputsWithNodes1(t *testing.T) {
    l := logrus.New()
    l.SetLevel(logrus.TraceLevel)
    dag, err := LoadDAG("../dag.json")
    if err != nil {
        t.Error(err)
        return
    }

    n := dag.reconcileInputsAndOutputsWithNodes()
    if n != 0 {
        t.Errorf("nodesWithInputsAndOuputs n: expected 0, got %d", n)
    }
}

func TestReconInputsAndOutputsWithNodes2(t *testing.T) {
    l := logrus.New()
    l.SetLevel(logrus.TraceLevel)
    dag, err := LoadDAG("../dag.json")
    if err != nil {
        t.Error(err)
        return
    }

    const nodeToDelete = "create_final_model"
    delete(dag.Nodes, nodeToDelete)

    n := dag.reconcileInputsAndOutputsWithNodes()
    if n != 2 {
        t.Errorf("nodesWithInputsAndOuputs n: expected 2, got %d", n)
    }
}
