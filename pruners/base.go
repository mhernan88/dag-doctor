package pruners

import (
    "github.com/mhernan88/dag-bisect/data"
)

type Pruner interface {
    PruneBefore(source *data.Node, roots []*data.Node) ([]string, error)
    PruneAfter(source *data.Node, roots []*data.Node) ([]string, error)
}
