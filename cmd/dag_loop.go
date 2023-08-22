package cmd

import (
    "fmt"
)

func (ui *UI) CheckDAG() error {
    fmt.Println("inspecting DAG")

    for len(ui.nodes) > 0 {
        node, err := ui.splitter.FindCandidate(ui.nodes)
        if err != nil {
            return err
        }
        ui.l.Tracef("selected split candidate: %s", node.Name)

        ui.CheckNode(node)
    }
    return nil
}
