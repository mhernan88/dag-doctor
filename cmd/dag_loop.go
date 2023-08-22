package cmd

func (ui *UI) CheckDAG() error {
    ui.l.Trace("inspecting DAG")

    for len(ui.nodes) > 0 {
        node, err := ui.splitter.FindCandidate(ui.nodes)
        if err != nil {
            return err
        }
        ui.l.Tracef("selected split candidate: %s", node.Name)
        break
    }
    return nil
}
