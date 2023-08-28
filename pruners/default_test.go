package pruners

import (
	"slices"
	"testing"

	"github.com/mhernan88/dag-bisect/data"
	"github.com/mhernan88/dag-bisect/utils"
	"github.com/sirupsen/logrus"
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
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)

	dag, err := data.LoadDAG("../dag.json")
	if err != nil {
		t.Error(err)
		return
	}
	nodes := utils.FlattenAllNodesToMap(dag)

	p := NewDefaultPruner(99, l)
	_, pruneableNodes := p.findUpstreamPruneableNodes(nodes["create_wide_table"])
	t.Logf("pruneable nodes: %v", pruneableNodes)

	expectedPruneableNodes := []string{
		"preprocess_companies_and_employees",
		"preprocess_shuttles_and_routes",
		"preprocess_reviews_and_ratings",
	}

	for _, expectedPruneableNode := range expectedPruneableNodes {
		if !slices.Contains(pruneableNodes, expectedPruneableNode) {
			t.Errorf(
				"expected node '%s' to be in pruneable nodes (%v)",
				expectedPruneableNode,
				pruneableNodes,
			)
		}
	}
}

func TestFindUpstreamPruneableNodes2(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)

	dag, err := data.LoadDAG("../dag.json")
	if err != nil {
		t.Error(err)
		return
	}
	nodes := utils.FlattenAllNodesToMap(dag)

	p := NewDefaultPruner(99, l)
	_, pruneableNodes := p.findUpstreamPruneableNodes(nodes["preprocess_shuttles_and_routes"])
	t.Logf("pruneable nodes: %v", pruneableNodes)

	expectedPruneableNodes := []string{}

	for _, expectedPruneableNode := range expectedPruneableNodes {
		if !slices.Contains(pruneableNodes, expectedPruneableNode) {
			t.Errorf(
				"expected node '%s' to be in pruneable nodes (%v)",
				expectedPruneableNode,
				pruneableNodes,
			)
		}
	}
}

func TestUnlinkNext(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)

	dag, err := data.LoadDAG("../dag.json")
	if err != nil {
		t.Error(err)
		return
	}

	p := NewDefaultPruner(99, l)

	if len(data.UniqueNodes(dag).ToSlice()) != 6 {
		t.Errorf(
			"nodes before PruneBefore, expected %d, got %d",
			4, len(data.UniqueNodes(dag).ToSlice()),
		)
		t.Logf("nodes: %v", data.UniqueNodes(dag))
		return
	}

	p.unlinkNext(dag, []string{"create_wide_table"})

	if len(data.UniqueNodes(dag).ToSlice()) != 3 {
		t.Errorf(
			"nodes before PruneBefore, expected %d, got %d",
			4, len(data.UniqueNodes(dag).ToSlice()),
		)
		t.Logf("nodes: %v", data.UniqueNodes(dag))
		return
	}
}

func TestPruneAfter(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)

	dag, err := data.LoadDAG("../dag.json")
	if err != nil {
		t.Error(err)
		return
	}
	nodes := utils.FlattenAllNodesToMap(dag)

	p := NewDefaultPruner(99, l)
	if len(data.UniqueNodes(dag).ToSlice()) != 6 {
		t.Errorf(
			"nodes before PruneAfter, expected %d, got %d",
			6, len(data.UniqueNodes(dag).ToSlice()),
		)
		t.Logf("nodes: %v", data.UniqueNodes(dag))
		return
	}

	_, err = p.PruneAfter(nodes["create_wide_table"], dag)
	if err != nil {
		t.Error(err)
		return
	}

	if len(data.UniqueNodes(dag).ToSlice()) != 4 {
		t.Errorf(
			"nodes after PruneAfter, expected %d, got %d",
			4, len(data.UniqueNodes(dag).ToSlice()),
		)
		t.Logf("nodes: %v", data.UniqueNodes(dag))
		return
	}
}

func TestPruneBefore(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)

	dag, err := data.LoadDAG("../dag.json")
	if err != nil {
		t.Error(err)
		return
	}
	nodes := utils.FlattenAllNodesToMap(dag)

	p := NewDefaultPruner(99, l)

	const nodesBefore = 6
	if len(data.UniqueNodes(dag).ToSlice()) != nodesBefore {
		t.Errorf(
			"nodes before PruneBefore, expected %d, got %d",
			nodesBefore, len(data.UniqueNodes(dag).ToSlice()),
		)
		t.Logf("nodes: %v", data.UniqueNodes(dag))
		return
	}

	dag, prunedNodes, err := p.PruneBefore(nodes["create_wide_table"], dag)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("pruned: %v", prunedNodes)

	const nodesAfter = 3
	if len(data.UniqueNodes(dag).ToSlice()) != nodesAfter {
		// Should be pruned:
		// - prep_shuttles_and_routes
		// - prep_companies_and_employees
		// - prep_reviews_and_ratings
		// Should remain:
		// - create_wide_table
		// - postprocess_table
		// - create_final_model
		t.Errorf(
			"nodes after PruneBefore: expected %d, got %d",
			nodesAfter, len(data.UniqueNodes(dag).ToSlice()),
		)
		t.Logf("nodes: %v", data.UniqueNodes(dag))
		return
	}
}
