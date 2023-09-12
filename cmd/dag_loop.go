package cmd

import (
	"fmt"

	"github.com/enescakir/emoji"
	"github.com/mhernan88/dag-bisect/data"
)

func (ui *UI) terminate() {
	if ui.lastFailedNode == "" {
		fmt.Printf(
			"%v dag ok\n",
			emoji.GrinningFace,
		)
	} else {
		fmt.Printf(
			"%v source of error: '%s'\n",
			emoji.Skull,
			ui.lastFailedNode,
		)
	}
}

func (ui *UI) CheckDAG() error {
	fmt.Println("inspecting DAG")

	var node data.Node
	var err error

	for (len(ui.dag.Nodes) > 0) && (len(ui.dag.Roots) > 0) {
		node, err = ui.splitter.FindCandidate(ui.dag)
		if err != nil {
			return err
		}

		ui.l.Tracef("selected split candidate: %s", node.Name)

		prunedNodes, err := ui.CheckNode(node)
		if err != nil {
			return err
		}
		ui.l.Tracef("prunedNodes = %v", data.SliceMapKeys(prunedNodes))
        ui.l.Debugf(
            "%d ok nodes; %d err nodes; %d remaining nodes", 
            len(ui.okNodes), 
            len(ui.errNodes),
            len(ui.dag.Nodes),
        )

	}

	// Automatically terminates on the last-viewed node.
	// Not necessarily the error node (esp if the last node was OK)
	ui.terminate()
	return nil
}
