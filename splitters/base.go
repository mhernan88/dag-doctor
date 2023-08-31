package splitters

import (
	"github.com/mhernan88/dag-bisect/data"
)

type Splitter interface {
	FindCandidate(dag data.DAG) (*data.Node, error)
}
