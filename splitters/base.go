package splitters

import (
	"log/slog"
	"github.com/mhernan88/dag-bisect/data"
)

type Splitter interface {
	FindCandidate(dag data.DAG, l *slog.Logger) (data.Node, error)
	GetName() string
}
