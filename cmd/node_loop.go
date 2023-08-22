package cmd

import (
    "fmt"
    "github.com/mhernan88/dag-bisect/data"
    "github.com/enescakir/emoji"
)

func (ui *UI) CheckNode(node *data.Node) (bool, []string, error) {
    fmt.Printf("inspecting node: %s\n", node.Name)
    allDatasetsOK := true

    for _, output := range node.Outputs {
        ok, err := ui.CheckDataset(output)
        if err != nil {
            return false, nil, err
        }
        if !ok {
            allDatasetsOK = false
        }
    }

    var nodeList []*data.Node
    for _, node := range ui.nodes {
        nodeList = append(nodeList, node)
    }

    var prunedNodes []string
    var err error
    if allDatasetsOK {
        prunedNodes, err = ui.pruner.PruneBefore(node, nodeList)
        if err != nil {
            return allDatasetsOK, nil, err
        }
        fmt.Printf("%v node %s cleared OK\n", emoji.CheckMarkButton, node.Name)
        fmt.Printf("|---> pruned nodes: %v\n", prunedNodes)
        return true, prunedNodes, nil
    } else {
        prunedNodes, err := ui.pruner.PruneAfter(node, nodeList)
        if err != nil {
            return allDatasetsOK, nil, err
        }
        fmt.Printf("%v node %s has ERR\n", emoji.CrossMarkButton, node.Name)
        fmt.Printf("|---> pruned nodes: %v\n", prunedNodes)
        return false, prunedNodes, nil
    }
}
