package cmd

import (
	"fmt"

	"github.com/enescakir/emoji"
)

func (ui *UI) CheckDAG() error {
	fmt.Println("inspecting DAG")

	dagOK := true

	for len(ui.dag.Nodes) > 0 {
		node, err := ui.splitter.FindCandidate(ui.dag)
		if err != nil {
			return err
		}

		if node == nil {
			return fmt.Errorf("failed to find split candidate")
		}

		ui.l.Tracef("selected split candidate: %s", node.Name)

		ok, prunedNodes, err := ui.CheckNode(*node)
		ui.l.Tracef("prunedNodes = %v", prunedNodes)
		if err != nil {
			return err
		}

		if !ok {
			dagOK = false
		}

		// if !dagOK && (len(node.Next) == 0) {
		// 	fmt.Printf(
		// 		"%v source of error: '%s'\n",
		// 		emoji.Skull,
		// 		node.Name,
		// 	)
		// }

		if len(prunedNodes) == 0 {
			if dagOK {
				fmt.Printf(
					"%v dag ok\n",
					emoji.GrinningFace,
				)
				return nil
			} else {
				fmt.Printf(
					"%v source of error: '%s'\n",
					emoji.Skull,
					node.Name,
				)
				return nil
			}
		}
	}
	return nil
}
