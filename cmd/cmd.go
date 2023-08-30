package cmd

import (
	"github.com/mhernan88/dag-bisect/data"
	"github.com/mhernan88/dag-bisect/pruners"
	"github.com/mhernan88/dag-bisect/splitters"
	"github.com/sirupsen/logrus"
)

func NewUI(
	dag data.DAG,
	splitter splitters.Splitter,
	pruner pruners.Pruner,
	iterationLimit int,
	l *logrus.Logger,
) UI {
	return UI{
		dag:            dag,
		okNodes:        make(map[string]data.Node),
		errNodes:       make(map[string]data.Node),
		splitter:       splitter,
		pruner:         pruner,
		iterationLimit: iterationLimit,
		l:              l,
	}
}

type UI struct {
	dag            data.DAG
	okNodes        map[string]data.Node
	errNodes       map[string]data.Node
	splitter       splitters.Splitter
	pruner         pruners.Pruner
	iterationLimit int
	l              *logrus.Logger
}

func (ui *UI) Run() error {
	ui.l.Debug("running ui loop")
	err := ui.CheckDAG()
	if err != nil {
		return err
	}
	ui.l.Debug("terminating ui")
	return nil
}
