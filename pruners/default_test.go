package pruners

import (
    "slices"
    "testing"
    "github.com/sirupsen/logrus"
    "github.com/mhernan88/dag-bisect/data"
    "github.com/mhernan88/dag-bisect/utils"
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
    }
    nodes := utils.FlattenAllNodesToMap(dag)

    p := NewDefaultPruner(99, l)
    _, pruneableNodes, err := p.findUpstreamPruneableNodes(nodes["create_wide_table"])
    if err != nil {
        t.Error(err)
    }

    expectedPruneableNodes := []string {
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
