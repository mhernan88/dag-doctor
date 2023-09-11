package pruners

import (
	"github.com/mhernan88/dag-bisect/data"
	"github.com/sirupsen/logrus"
)

func NewDefaultPruner(iterationLimit int, l *logrus.Logger) DefaultPruner {
	return DefaultPruner{
		iterationLimit: iterationLimit,
		l:              l,
	}
}

type DefaultPruner struct {
	iterationLimit int
	l              *logrus.Logger
}

// Finds all nodes that can be pruned before `source`
// assuming `source` is error-free.
func (p DefaultPruner) findUpstreamPruneableNodes(
	node string,
	dag data.DAG,
) map[string]data.Node {
	p.l.Tracef("finding pruneable nodes before %s", node)
	if len(dag.Nodes[node].Prev) == 0 {
		return make(map[string]data.Node)
	}

	ancestorsMap := dag.Ancestors(node)
	p.l.Tracef("pulling pruneable nodes from %d ancestors", len(ancestorsMap))

	pruneableAncestorsMap := make(map[string]data.Node)
	for ancestorName, ancestor := range ancestorsMap {
		isPruneable := true
		for _, childName := range ancestor.Next {
			_, ok := ancestorsMap[childName]
			if !ok {
				p.l.Tracef("%s is not pruneable (it is not an ancestor)", ancestorName)
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
) (data.DAG, map[string]data.Node) {
	p.l.Tracef("pruning nodes before %s", node)
	pruneableNodes := p.findUpstreamPruneableNodes(node, dag)
	pruneableNodes[node] = dag.Nodes[node]

	for name := range pruneableNodes {
		p.l.Tracef("PruneBefore popping node %s", name)
		dag.Pop(name)
	}

	dag.ReconcileInputsAndOutputsWithNodes()
	return dag, pruneableNodes
}

func (p DefaultPruner) PruneAfter(
	node string,
	dag data.DAG,
) (data.DAG, map[string]data.Node) {
	p.l.Tracef("pruning nodes after %s", node)
	pruneableNodes := make(map[string]data.Node)
	pruneableNodes[node] = dag.Nodes[node]

	for descendantName, descendant := range dag.Descendants(node) {
		pruneableNodes[descendantName] = descendant
	}

	for name := range pruneableNodes {
		p.l.Tracef("PruneAfter popping node %s", name)
		dag.Pop(name)
	}

	dag.ReconcileInputsAndOutputsWithNodes()
	return dag, pruneableNodes
}
