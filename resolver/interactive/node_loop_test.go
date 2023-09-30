package interactive

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/mhernan88/dag-bisect/models"
	"github.com/mhernan88/dag-bisect/resolver/pruners"
)

func TestUI__PruneNodes(t *testing.T) {
	l := slog.Default()

	dagPtr, err := models.LoadDAG("../dag.json")
	if err != nil {
		t.Error(err)
		return
	}
	dag := *dagPtr

	p := pruners.NewDefaultPruner()

    state := models.NewState(dag)
    prunedNodes := pruneNodes(
		state,
		p,
        dag.Nodes["preprocess_companies_and_employees"],
        false,
		l,
    )

    const expected = 5
    if len(prunedNodes) != expected {
        t.Errorf("prunedNodes; want=%d, got=%d", expected, len(prunedNodes))
        prunedNodesString := models.SliceMapKeys(prunedNodes)
        prunedNodesFormattedString := strings.Join(prunedNodesString, "\n- ")
        t.Logf("Nodes:\n- %v", prunedNodesFormattedString)
        return
    }
}
