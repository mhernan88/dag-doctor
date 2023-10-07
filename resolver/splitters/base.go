package splitters

import (
	"log/slog"

	"github.com/mhernan88/dag-bisect/models"
)

type Splitter interface {
	FindCandidate(dag models.DAG, l *slog.Logger) (models.Node, error)
	GetName() string
}
