package pruners

import (
	"log/slog"

	"github.com/mhernan88/dag-bisect/models"
)

func NewDefaultPruner() DefaultPruner {
	return DefaultPruner{
		Name: "default",
	}
}

type DefaultPruner struct {
	Name string `json:"name"`
}

func (p DefaultPruner) GetName() string {
	return p.Name
}

// Finds all nodes that can be pruned before `source`
// assuming `source` is error-free.
func (p DefaultPruner) findUpstreamPruneableNodes(
	node string,
	dag models.DAG,
	l *slog.Logger,
) map[string]models.Node {
	if len(dag.Nodes[node].Prev) == 0 {
		return make(map[string]models.Node)
	}

	ancestorsMap := dag.Ancestors(node)

	pruneableAncestorsMap := make(map[string]models.Node)
	for ancestorName, ancestor := range ancestorsMap {
		isPruneable := true
		for _, childName := range ancestor.Next {
			_, ok := ancestorsMap[childName]
			if !ok {
				isPruneable = false
				break
			}
		}

		if !isPruneable {
			pruneableAncestorsMap[ancestorName] = ancestor
		}
	}

	// Remaining nodes are ones we can safely prune.
	return pruneableAncestorsMap
}

func (p DefaultPruner) PruneBefore(
	node string,
	dag models.DAG,
	l *slog.Logger,
) (models.DAG, map[string]models.Node) {
	pruneableNodes := p.findUpstreamPruneableNodes(node, dag, l)
	pruneableNodes[node] = dag.Nodes[node]

	for name := range pruneableNodes {
		dag.Pop(name)
	}

	dag.ReconcileInputsAndOutputsWithNodes()
	return dag, pruneableNodes
}

func (p DefaultPruner) PruneAfter(
	node string,
	dag models.DAG,
	l *slog.Logger,
) (models.DAG, map[string]models.Node) {
	pruneableNodes := make(map[string]models.Node)
	pruneableNodes[node] = dag.Nodes[node]

	for descendantName, descendant := range dag.Descendants(node) {
		pruneableNodes[descendantName] = descendant
	}

	for name := range pruneableNodes {
		dag.Pop(name)
	}

	dag.ReconcileInputsAndOutputsWithNodes()
	return dag, pruneableNodes
}
