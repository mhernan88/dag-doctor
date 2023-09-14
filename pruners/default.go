package pruners

import (
	"github.com/mhernan88/dag-bisect/data"
	"github.com/sirupsen/logrus"
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
	l *logrus.Logger,
) map[string]data.Node {
	l.Tracef("finding pruneable nodes before %s", node)
	if len(dag.Nodes[node].Prev) == 0 {
		return make(map[string]data.Node)
	}

	ancestorsMap := dag.Ancestors(node)
	l.Tracef("pulling pruneable nodes from %d ancestors", len(ancestorsMap))

	pruneableAncestorsMap := make(map[string]data.Node)
	for ancestorName, ancestor := range ancestorsMap {
		isPruneable := true
		for _, childName := range ancestor.Next {
			_, ok := ancestorsMap[childName]
			if !ok {
				l.Tracef("%s is not pruneable (it is not an ancestor)", ancestorName)
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
	l *logrus.Logger,
) (data.DAG, map[string]data.Node) {
	l.Tracef("pruning nodes before %s", node)
	pruneableNodes := p.findUpstreamPruneableNodes(node, dag, l)
	pruneableNodes[node] = dag.Nodes[node]

	for name := range pruneableNodes {
		l.Tracef("PruneBefore popping node %s", name)
		dag.Pop(name)
	}

	dag.ReconcileInputsAndOutputsWithNodes()
	return dag, pruneableNodes
}

func (p DefaultPruner) PruneAfter(
	node string,
	dag data.DAG,
	l *logrus.Logger,
) (data.DAG, map[string]data.Node) {
	l.Tracef("pruning nodes after %s", node)
	pruneableNodes := make(map[string]data.Node)
	pruneableNodes[node] = dag.Nodes[node]

	for descendantName, descendant := range dag.Descendants(node) {
		pruneableNodes[descendantName] = descendant
	}

	for name := range pruneableNodes {
		l.Tracef("PruneAfter popping node %s", name)
		dag.Pop(name)
	}

	dag.ReconcileInputsAndOutputsWithNodes()
	return dag, pruneableNodes
}
