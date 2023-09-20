package cmd

import (
	"fmt"
	"log/slog"

	"github.com/enescakir/emoji"
	"github.com/mhernan88/dag-bisect/data"
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

func (ui *UI) CheckDAGIter(l *slog.Logger) error {
	fmt.Println("DAG")
	fmt.Println(ui.DAG)
	fmt.Println("ERR")
	fmt.Println(ui.ERRNodes)
	node, err := ui.Splitter.FindCandidate(ui.DAG, l)
	if err != nil {
		return err
	}

	l.Debug("selected split candidate", "candidate", node.Name)

	prunedNodes, err := ui.CheckNode(node, l)
	if err != nil {
		return err
	}

	l.Debug(
		"completed pruning nodes",
		"pruned nodes", data.SliceMapKeys(prunedNodes),
		"ok nodes", len(ui.OKNodes),
		"err nodes", len(ui.ERRNodes),
		"remaining nodes", len(ui.DAG.Nodes),
	)
	return nil
}

func (ui *UI) CheckDAG(l *slog.Logger) error {
	fmt.Println("inspecting DAG")

	var node data.Node
	var err error

	for (len(ui.DAG.Nodes) > 0) && (len(ui.DAG.Roots) > 0) {
		node, err = ui.Splitter.FindCandidate(ui.DAG, l)
		if err != nil {
			return err
		}

		l.Debug("selected split candidate", "candidate", node.Name)

		prunedNodes, err := ui.CheckNode(node, l)
		if err != nil {
			return err
		}

		l.Debug(
			"completed pruning nodes",
			"pruned nodes", data.SliceMapKeys(prunedNodes),
			"ok nodes", len(ui.OKNodes),
			"err nodes", len(ui.ERRNodes),
			"remaining nodes", len(ui.DAG.Nodes),
		)
	}

	// Automatically terminates on the last-viewed node.
	// Not necessarily the error node (esp if the last node was OK)
	ui.Terminate()
	return nil
}
