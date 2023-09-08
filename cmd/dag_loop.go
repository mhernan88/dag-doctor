package cmd

import (
	"fmt"

	"github.com/enescakir/emoji"
	"github.com/mhernan88/dag-bisect/data"
)

func (ui *UI) terminate(ok bool, nodeName string) {
	if ok {
		fmt.Printf(
			"%v dag ok\n",
			emoji.GrinningFace,
		)
	} else {
		fmt.Printf(
			"%v source of error: '%s'\n",
			emoji.Skull,
			nodeName,
		)
	}
}

func (ui *UI) CheckDAG() error {
	fmt.Println("inspecting DAG")

	dagOK := true
	var node data.Node
	var err error

	for (len(ui.dag.Nodes) > 0) && (len(ui.dag.Roots) > 0) {
		node, err = ui.splitter.FindCandidate(ui.dag)
		if err != nil {
			return err
		}

		ui.l.Tracef("selected split candidate: %s", node.Name)

		ok, prunedNodes, err := ui.CheckNode(node)
		if err != nil {
			return err
		}
		ui.l.Tracef("prunedNodes = %v", data.SliceMapKeys(prunedNodes))
        ui.l.Debugf(
            "%d ok nodes; %d err nodes", 
            len(ui.errNodes), 
            len(ui.okNodes),
        )

		if !ok {
			dagOK = false
		}
	}
	ui.terminate(dagOK, node.Name)
	return nil
}
