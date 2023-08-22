package cmd

import (
    "github.com/dag-bisect/data
)

type UI struct {
    nodes map[string]*data.Node
    okNodes map[string]*data.Node
    errNodes map[string]*data.Node
}

func Run(l *logrus.Logger) {
}
