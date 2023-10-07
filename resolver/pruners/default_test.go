package pruners

import (
	"log/slog"
	"testing"

	"github.com/mhernan88/dag-bisect/models"
)

// DAG
//                       -> prep_shut_and_rout
//                     /                     \
// prep_comp_and_emp--<                       \                  ->postprocess_table
//                     \                       \                /                   \
//                      -------------------------> wide_table--<                     \
//                                            /                \                      \
// prep_rev_and_rat--------------------------/                  ------------------------->create_final_table
//
// In pruning ancestors of wide_table, remove:
// - prep_shut_and_rout
// - prep_comp_and_emp
// - prep_rev_and_rat
//
// In pruning ancestors of prep_shut_and_rout, remove:
// - NOTHING
// |
// ----> prep_comp_and_emp also feeds into a non-ancestor node.

func TestFindUpstreamPruneableNodes(t *testing.T) {
	l := slog.Default()

	dagPtr, err := models.LoadDAG("../dag.json")
	if err != nil {
		t.Error(err)
		return
	}
	dag := *dagPtr

	p := NewDefaultPruner()
	pruneableNodes := p.findUpstreamPruneableNodes("create_wide_table", dag, l)
	t.Logf("pruneable nodes: %v", pruneableNodes)

	expectedPruneableNodes := []string{
		"preprocess_companies_and_employees",
		"preprocess_shuttles_and_routes",
		"preprocess_reviews_and_ratings",
	}

	for _, expectedPruneableNode := range expectedPruneableNodes {
		_, ok := pruneableNodes[expectedPruneableNode]
		if !ok {
			t.Errorf(
				"expected node '%s' to be in pruneable nodes (%v)",
				expectedPruneableNode,
				pruneableNodes,
			)
		}
	}
}

func TestFindUpstreamPruneableNodes2(t *testing.T) {
	l := slog.Default()

	dagPtr, err := models.LoadDAG("../dag.json")
	if err != nil {
		t.Error(err)
		return
	}
	dag := *dagPtr

	p := NewDefaultPruner()
	pruneableNodes := p.findUpstreamPruneableNodes("preprocess_shuttles_and_routes", dag, l)
	t.Logf("pruneable nodes: %v", pruneableNodes)

	expectedPruneableNodes := []string{}

	for _, expectedPruneableNode := range expectedPruneableNodes {
		_, ok := pruneableNodes[expectedPruneableNode]
		if !ok {
			t.Errorf(
				"expected node '%s' to be in pruneable nodes (%v)",
				expectedPruneableNode,
				pruneableNodes,
			)
		}
	}
}

func TestPruneAfter(t *testing.T) {
	l := slog.Default()

	dagPtr, err := models.LoadDAG("../dag.json")
	if err != nil {
		t.Error(err)
		return
	}
	dag := *dagPtr
	if len(dag.Nodes) != 6 {
		t.Error("dag loaded improperly")
		return
	}

	p := NewDefaultPruner()

    dag, prunedNodes := p.PruneAfter("create_wide_table", dag, l)
	if err != nil {
		t.Error(err)
		return
	}

    expected := 3
	if len(dag.Nodes) != expected {
		t.Errorf(
			"nodes after PruneAfter, expected %d, got %d",
			expected, len(dag.Nodes))
		t.Logf("nodes: %v", dag.Nodes)
		return
	}

    expected = 3
    if len(prunedNodes) != expected {
        t.Errorf(
            "prunedNodes after PruneAfter, expected %d, got %d",
            expected, len(prunedNodes))
        t.Logf("prunedNodes: %v", models.SliceMapKeys(prunedNodes))
        return
    }
}

func TestPruneBefore(t *testing.T) {
	l := slog.Default()

	dagPtr, err := models.LoadDAG("../dag.json")
	if err != nil {
		t.Error(err)
		return
	}
	dag := *dagPtr

	p := NewDefaultPruner()

	const nodesBefore = 6
	if len(dag.Nodes) != nodesBefore {
		t.Errorf(
			"nodes before PruneBefore, expected %d, got %d",
			nodesBefore, len(dag.Nodes),
		)
		t.Logf("nodes: %v", dag.Nodes)
		return
	}

	dag, prunedNodes := p.PruneBefore("create_wide_table", dag, l)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("pruned: %v", prunedNodes)

    expected := 2
	if len(dag.Nodes) != expected {
		// Should be pruned:
		// - prep_shuttles_and_routes
		// - prep_companies_and_employees
		// - prep_reviews_and_ratings
		// - create_wide_table
		// Should remain:
		// - postprocess_table
		// - create_final_model
		t.Errorf(
			"nodes after PruneBefore: expected %d, got %d",
			expected, len(dag.Nodes),
		)
		t.Logf("nodes: %v", dag.Nodes)
		return
	}

    expected = 4
    if len(prunedNodes) != expected {
        t.Errorf(
            "prunedNodes after PruneBefore, expected %d, got %d",
            expected, len(prunedNodes))
        t.Logf("prunedNodes: %v", models.SliceMapKeys(prunedNodes))
        return
    }
}
