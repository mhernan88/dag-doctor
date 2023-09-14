package pruners

import (
	"log/slog"

	"github.com/mhernan88/dag-bisect/data"
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
	dag data.DAG,
	l *slog.Logger,
) map[string]data.Node {
	if len(dag.Nodes[node].Prev) == 0 {
		return make(map[string]data.Node)
	}

	ancestorsMap := dag.Ancestors(node)

	pruneableAncestorsMap := make(map[string]data.Node)
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
	dag data.DAG,
	l *slog.Logger,
) (data.DAG, map[string]data.Node) {
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
	dag data.DAG,
	l *slog.Logger,
) (data.DAG, map[string]data.Node) {
	pruneableNodes := make(map[string]data.Node)
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
