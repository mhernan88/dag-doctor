package pruners

import (
    "slices"
    "testing"
    "github.com/sirupsen/logrus"
    "github.com/mhernan88/dag-bisect/data"
    "github.com/mhernan88/dag-bisect/utils"
)

func TestFindUpstreamPruneableNodes(t *testing.T) {
    l := logrus.New()
    l.SetLevel(logrus.TraceLevel)

    dag, err := data.LoadDAG("../dag.json")
    if err != nil {
        t.Error(err)
    }
    roots := utils.FlattenAllNodesToSlice(dag)
    nodes := utils.FlattenAllNodesToMap(dag)

    p := NewDefaultPruner(99, l)
    pruneableNodes, err := p.findUpstreamPruneableNodes(
        nodes["create_wide_table"],
        roots)
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
