package pruners

import (
	"log/slog"

	"github.com/mhernan88/dag-bisect/models"
)

type Pruner interface {
	PruneBefore(node string, dag models.DAG, l *slog.Logger) (models.DAG, map[string]models.Node)
	PruneAfter(node string, dag models.DAG, l *slog.Logger) (models.DAG, map[string]models.Node)
	GetName() string
}
