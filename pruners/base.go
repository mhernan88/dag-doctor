package pruners

import (
	"github.com/mhernan88/dag-bisect/data"
)

type Pruner interface {
	PruneBefore(node string, dag data.DAG) data.DAG
	PruneAfter(node string, dag data.DAG) data.DAG
}
