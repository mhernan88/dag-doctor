package cmd

import (
    "github.com/mhernan88/dag-bisect/data"
    "github.com/mhernan88/dag-bisect/splitters"
    "github.com/mhernan88/dag-bisect/pruners"
    "github.com/sirupsen/logrus"
)

func NewUI(
    nodes map[string]*data.Node,
    catalog map[string]data.Dataset,
    splitter splitters.Splitter,
    pruner pruners.Pruner,
    l *logrus.Logger,
) UI {
    return UI{
        nodes: nodes,
        okNodes: make(map[string]*data.Node),
        errNodes: make(map[string]*data.Node),
        catalog: catalog,
        splitter: splitter,
        pruner: pruner,
        l: l,
    }
}

type UI struct {
    nodes map[string]*data.Node
    okNodes map[string]*data.Node
    errNodes map[string]*data.Node
    catalog map[string]data.Dataset
    splitter splitters.Splitter
    pruner pruners.Pruner
    l *logrus.Logger
}

func (ui UI) Run() {
    ui.l.Debug("running ui loop")
    ui.l.Debug("terminating ui")
}
