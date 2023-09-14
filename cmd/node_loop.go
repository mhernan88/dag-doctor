package cmd

import (
	"fmt"
	"log/slog"

	"github.com/enescakir/emoji"
	"github.com/mhernan88/dag-bisect/data"
)

func (ui *UI) checkDatasets(node data.Node, l *slog.Logger) (bool, error) {
	for _, output := range node.Outputs {
		ok, err := ui.CheckDataset(output, l)
		if err != nil {
			return false, err
		}
		if !ok {
            return false, nil
		}
	}
    return true, nil
}

func (ui *UI) pruneNodes(
	node data.Node, 
	allDatasetsOK bool,
	l *slog.Logger,
) (
	map[string]data.Node,
) {
	var prunedNodes map[string]data.Node
	if allDatasetsOK {
		ui.DAG, prunedNodes = ui.Pruner.PruneBefore(node.Name, ui.DAG, l)
        for name, node := range prunedNodes {
            ui.OKNodes[name] = node
        }
		fmt.Printf("|-> %v node %s cleared OK\n", emoji.CheckMarkButton, node.Name)
		fmt.Printf("|---> pruned upstream nodes: %v\n", data.SliceMapKeys(prunedNodes))
		return prunedNodes
	} 

    ui.DAG, prunedNodes = ui.Pruner.PruneAfter(node.Name, ui.DAG, l)
    for name, node := range prunedNodes {
        ui.ERRNodes[name] = node
    }
    fmt.Printf("|-> %v node %s has ERR\n", emoji.CrossMarkButton, node.Name)
    fmt.Printf("|---> pruned downstream nodes: %v\n", data.SliceMapKeys(prunedNodes))
	ui.LastFailedNode = node.Name
    return prunedNodes
}

func (ui *UI) CheckNode(node data.Node, l *slog.Logger) (map[string]data.Node, error) {
	fmt.Printf("|-> %v inspecting node: %s\n", emoji.Microscope, node.Name)
    allDatasetsOK, err := ui.checkDatasets(node, l)
    if err != nil {
        return nil, err
    }
    prunedNodes := ui.pruneNodes(node, allDatasetsOK, l)
    return prunedNodes, nil
}
