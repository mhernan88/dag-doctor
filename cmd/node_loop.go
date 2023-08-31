package cmd

import (
	"fmt"

	"github.com/enescakir/emoji"
	"github.com/mhernan88/dag-bisect/data"
)

func (ui *UI) CheckNode(node data.Node) (bool, []string, error) {
	fmt.Printf("|-> %v inspecting node: %s\n", emoji.Microscope, node.Name)
	allDatasetsOK := true

	for _, output := range node.Outputs {
		ok, err := ui.CheckDataset(output)
		if err != nil {
			return false, nil, err
		}
		if !ok {
			allDatasetsOK = false
			break
		}
	}

	var prunedNodes []string
	if allDatasetsOK {
		ui.dag = ui.pruner.PruneBefore(node.Name, ui.dag)
		fmt.Printf("|-> %v node %s cleared OK\n", emoji.CheckMarkButton, node.Name)
		// TODO: Add pruned nodes to other list.
		// fmt.Printf("|---> pruned nodes: %v\n", prunedNodes)
		return true, nil, nil
	} else {
		ui.dag = ui.pruner.PruneAfter(node.Name, ui.dag)
		fmt.Printf("|-> %v node %s has ERR\n", emoji.CrossMarkButton, node.Name)
		fmt.Printf("|---> pruned nodes: %v\n", prunedNodes)
		return false, nil, nil
	}
}
