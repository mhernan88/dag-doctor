package cmd

import (
	"fmt"
	"log/slog"

	"github.com/enescakir/emoji"
	"github.com/mhernan88/dag-bisect/models"
)

func (ui *UI) Terminate() {
	if ui.LastFailedNode == "" {
		fmt.Printf(
			"%v dag ok\n",
			emoji.GrinningFace,
		)
	} else {
		fmt.Printf(
			"%v source of error: '%s'\n",
			emoji.Skull,
			ui.LastFailedNode,
		)
	}
}

func (ui *UI) CheckDAGIter(l *slog.Logger) (bool, error) {
	node, err := ui.Splitter.FindCandidate(ui.DAG, l)
	if err != nil {
		return false, err
	}
	l.Debug("selected split candidate", "candidate", node.Name)

	prunedNodes, err := ui.CheckNode(node, l)
	if err != nil {
		return false, err
	}

	if len(prunedNodes) == 0 {
		return true, nil
	}

	l.Debug(
		"completed pruning nodes",
		"pruned nodes", models.SliceMapKeys(prunedNodes),
		"ok nodes", len(ui.OKNodes),
		"err nodes", len(ui.ERRNodes),
		"remaining nodes", len(ui.DAG.Nodes),
	)
	return false, nil
}

func (ui *UI) CheckDAG(l *slog.Logger) (int, error) {
	fmt.Println("inspecting DAG")
	var err error

	abort := false
	i := 0
	for (len(ui.DAG.Nodes) > 0) && (len(ui.DAG.Roots) > 0) {
		abort, err = ui.CheckDAGIter(l)
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
