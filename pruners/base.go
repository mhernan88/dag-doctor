package pruners

import (
	"github.com/mhernan88/dag-bisect/data"
)

type Pruner interface {
	PruneBefore(node string, dag data.DAG) (data.DAG, map[string]data.Node)
	PruneAfter(node string, dag data.DAG) (data.DAG, map[string]data.Node)
}
