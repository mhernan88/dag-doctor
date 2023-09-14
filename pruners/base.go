package pruners

import (
	"log/slog"

	"github.com/mhernan88/dag-bisect/data"
)

type Pruner interface {
	PruneBefore(node string, dag data.DAG, l *slog.Logger) (data.DAG, map[string]data.Node)
	PruneAfter(node string, dag data.DAG, l *slog.Logger) (data.DAG, map[string]data.Node)
	GetName() string
}
