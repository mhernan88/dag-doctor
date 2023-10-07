package interactive

import (
	"fmt"
	"log/slog"

	"github.com/enescakir/emoji"
	"github.com/mhernan88/dag-bisect/models"
	"github.com/mhernan88/dag-bisect/resolver/pruners"
	"github.com/mhernan88/dag-bisect/resolver/splitters"
)

func Terminate(state *models.State) {
	if state.LastFailedNode == "" {
		fmt.Printf(
			"%v dag ok\n",
			emoji.GrinningFace,
		)
	} else {
		fmt.Printf(
			"%v source of error: '%s'\n",
			emoji.Skull,
			state.LastFailedNode,
		)
	}
}

func CheckDAGIter(state *models.State, pruner pruners.DefaultPruner, splitter splitters.DefaultSplitter, l *slog.Logger) (bool, error) {
	node, err := splitter.FindCandidate(state.DAG, l)
	if err != nil {
		return false, err
	}
	l.Info("selected split candidate", "candidate", node.Name)

	prunedNodes, err := CheckNode(state, pruner, node, l)
	if err != nil {
		return false, err
	}

	if len(prunedNodes) == 0 {
		return true, nil
	}

	l.Info(
		"completed pruning nodes",
		"pruned nodes", models.SliceMapKeys(prunedNodes),
		"ok nodes", len(state.OKNodes),
		"err nodes", len(state.ERRNodes),
		"remaining nodes", len(state.DAG.Nodes),
	)
	return false, nil
}

func CheckDAG(state *models.State, pruner pruners.DefaultPruner, splitter splitters.DefaultSplitter, l *slog.Logger) (int, error) {
	fmt.Println("inspecting DAG")
	var err error

	abort := false
	i := 0
	for (len(state.DAG.Nodes) > 0) && (len(state.DAG.Roots) > 0) {
		abort, err = CheckDAGIter(state, pruner, splitter, l)
		if err != nil {
			return i, fmt.Errorf("failed to check dag | %v", err)
		}

		i++
		if abort {
			fmt.Println("exiting...")
			return i, nil
		}
	}

	// Automatically terminates on the last-viewed node.
	// Not necessarily the error node (esp if the last node was OK)
	return i, nil
}
