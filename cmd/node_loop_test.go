package cmd

import (
	"strings"
	"testing"

	"github.com/mhernan88/dag-bisect/data"
	"github.com/mhernan88/dag-bisect/pruners"
	"github.com/mhernan88/dag-bisect/splitters"
	"github.com/sirupsen/logrus"
)

func TestUI__PruneNodes(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)

	dagPtr, err := data.LoadDAG("../dag.json")
	if err != nil {
		t.Error(err)
		return
	}
	dag := *dagPtr

	p := pruners.NewDefaultPruner(99, l)
    s := splitters.NewDefaultSplitter(99, l)

    ui := NewUI(dag, s, p, 99, l)
    prunedNodes := ui.pruneNodes(
        dag.Nodes["preprocess_companies_and_employees"],
        false,
    )

    const expected = 5
    if len(prunedNodes) != expected {
        t.Errorf("prunedNodes; want=%d, got=%d", expected, len(prunedNodes))
        prunedNodesString := data.SliceMapKeys(prunedNodes)
        prunedNodesFormattedString := strings.Join(prunedNodesString, "\n- ")
        t.Logf("Nodes:\n- %v", prunedNodesFormattedString)
        return
    }
}
