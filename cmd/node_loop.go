package cmd

import (
	"fmt"

	"github.com/enescakir/emoji"
	"github.com/mhernan88/dag-bisect/data"
)

func (ui *UI) CheckNode(node *data.Node) (bool, []string, error) {
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

	var roots map[string]*data.Node
	var prunedNodes []string
	var err error
	if allDatasetsOK {
		roots, prunedNodes, err = ui.pruner.PruneBefore(node, ui.nodes)
		if err != nil {
			return allDatasetsOK, nil, err
		}

		if roots != nil {
			ui.l.Debugf(
				"updating roots from %d roots to %d roots",
				len(ui.nodes),
				len(roots),
			)
			ui.nodes = roots
		} else {
			ui.l.Debugf("roots remain unchanged")
		}

		fmt.Printf("|-> %v node %s cleared OK\n", emoji.CheckMarkButton, node.Name)
		fmt.Printf("|---> pruned nodes: %v\n", prunedNodes)
		uniquePrunedNodes := data.LoadNodesIntoSet(prunedNodes)
		allNodes := data.UniqueNodes(roots)
		remainingNodes := allNodes.Difference(uniquePrunedNodes).ToSlice()
		fmt.Printf("|---> %d nodes remaining\n", len(remainingNodes))
		return true, prunedNodes, nil
	} else {
		roots, prunedNodes, err = ui.pruner.PruneAfter(node, ui.nodes)
		if err != nil {
			return allDatasetsOK, nil, err
		}
		if roots != nil {
			ui.l.Debugf(
				"updating roots from %d roots to %d roots",
				len(ui.nodes),
				len(roots),
			)
			ui.nodes = roots
		} else {
			ui.l.Debugf("roots remain unchanged")
		}
		fmt.Printf("|-> %v node %s has ERR\n", emoji.CrossMarkButton, node.Name)
		fmt.Printf("|---> pruned nodes: %v\n", prunedNodes)
		uniquePrunedNodes := data.LoadNodesIntoSet(prunedNodes)
		allNodes := data.UniqueNodes(roots)
		remainingNodes := allNodes.Difference(uniquePrunedNodes).ToSlice()
		fmt.Printf("|---> %d nodes remaining\n", len(remainingNodes))
		return false, prunedNodes, nil
	}
}
