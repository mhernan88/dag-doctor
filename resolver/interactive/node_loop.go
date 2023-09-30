package interactive

import (
	"fmt"
	"log/slog"

	"github.com/enescakir/emoji"
	"github.com/mhernan88/dag-bisect/models"
	"github.com/mhernan88/dag-bisect/resolver/pruners"
)

func checkDatasets(node models.Node, l *slog.Logger) (string, error) {
	for _, output := range node.Outputs {
		status, err := CheckDataset(output, l)
		if (err != nil) || (status == "err") || (status == "aborted") {
			return status, err
		}
	}
    return "ok", nil
}

func pruneNodes(
	state *models.State,
	pruner pruners.DefaultPruner,
	node models.Node, 
	allDatasetsOK bool,
	l *slog.Logger,
) (
	map[string]models.Node,
) {
	var prunedNodes map[string]models.Node
	if allDatasetsOK {
		state.DAG, prunedNodes = pruner.PruneBefore(node.Name, state.DAG, l)
        for name, node := range prunedNodes {
            state.OKNodes[name] = node
        }
		fmt.Printf("|-> %v node %s cleared OK\n", emoji.CheckMarkButton, node.Name)
		fmt.Printf("|---> pruned upstream nodes: %v\n", models.SliceMapKeys(prunedNodes))
		return prunedNodes
	} 

    state.DAG, prunedNodes = pruner.PruneAfter(node.Name, state.DAG, l)
    for name, node := range prunedNodes {
        state.ERRNodes[name] = node
    }
    fmt.Printf("|-> %v node %s has ERR\n", emoji.CrossMarkButton, node.Name)
    fmt.Printf("|---> pruned downstream nodes: %v\n", models.SliceMapKeys(prunedNodes))
	state.LastFailedNode = node.Name
    return prunedNodes
}

func CheckNode(state *models.State, pruner pruners.DefaultPruner, node models.Node, l *slog.Logger) (map[string]models.Node, error) {
	fmt.Printf("|-> %v inspecting node: %s\n", emoji.Microscope, node.Name)
    datasetStatus, err := checkDatasets(node, l)

	var allDatasetsOK bool
    if (err != nil) || (datasetStatus == "aborted") {
        return nil, err
    } else if datasetStatus == "ok" {
		allDatasetsOK = true
	} else if datasetStatus == "err" {
		allDatasetsOK = false
	} else {
		return nil, fmt.Errorf("invalid dataset status")
	}

    prunedNodes := pruneNodes(state, pruner, node, allDatasetsOK, l)
    return prunedNodes, nil
}
