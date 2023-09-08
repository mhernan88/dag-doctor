package cmd

import (
	"fmt"

	"github.com/enescakir/emoji"
	"github.com/mhernan88/dag-bisect/data"
)

func (ui *UI) checkDatasets(node data.Node) (bool, error) {
	for _, output := range node.Outputs {
		ok, err := ui.CheckDataset(output)
		if err != nil {
			return false, err
		}
		if !ok {
            return false, nil
		}
	}
    return true, nil
}

func (ui *UI) pruneNodes(node data.Node, allDatasetsOK bool) map[string]data.Node {
	var prunedNodes map[string]data.Node
	if allDatasetsOK {
		ui.dag, prunedNodes = ui.pruner.PruneBefore(node.Name, ui.dag)
        for name, node := range prunedNodes {
            ui.okNodes[name] = node
        }
		fmt.Printf("|-> %v node %s cleared OK\n", emoji.CheckMarkButton, node.Name)
		fmt.Printf("|---> pruned upstream nodes: %v\n", data.SliceMapKeys(prunedNodes))
		return prunedNodes
	} 

    ui.dag, prunedNodes = ui.pruner.PruneAfter(node.Name, ui.dag)
    for name, node := range prunedNodes {
        ui.errNodes[name] = node
    }
    fmt.Printf("|-> %v node %s has ERR\n", emoji.CrossMarkButton, node.Name)
    fmt.Printf("|---> pruned downstream nodes: %v\n", data.SliceMapKeys(prunedNodes))
    return prunedNodes
}

func (ui *UI) CheckNode(node data.Node) (bool, map[string]data.Node, error) {
	fmt.Printf("|-> %v inspecting node: %s\n", emoji.Microscope, node.Name)
    allDatasetsOK, err := ui.checkDatasets(node)
    if err != nil {
        return false, nil, err
    }
    prunedNodes := ui.pruneNodes(node, allDatasetsOK)
    return allDatasetsOK, prunedNodes, nil
}
