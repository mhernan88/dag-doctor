package cmd

import (
	"fmt"

	"github.com/enescakir/emoji"
)

func (ui *UI) CheckDAG() error {
	fmt.Println("inspecting DAG")

	dagOK := true

	for len(ui.nodes) > 0 {
		node, err := ui.splitter.FindCandidate(ui.nodes)
		if err != nil {
			return err
		}

		ui.l.Tracef("selected split candidate: %s", node.Name)

		ok, prunedNodes, err := ui.CheckNode(node)
		ui.l.Tracef("prunedNodes = %v", prunedNodes)
		if err != nil {
			return err
		}

		if !ok {
			dagOK = false
		}

		if !dagOK && (node.Next == nil) {
			fmt.Printf(
				"%v source of error: '%s'\n",
				emoji.Skull,
				node.Name,
			)
		}

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
