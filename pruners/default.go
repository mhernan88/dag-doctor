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
	if len(dag.Nodes[node].Prev) == 0 {
		p.l.Debugf("node %s had no parents to prune", node)
		dag.Pop(node)
		dag.Reconcile()
		return dag, nil
	}

	pruneableNodes := p.findUpstreamPruneableNodes(node, dag)
	for rootName := range dag.Roots {
		_, ok := pruneableNodes[rootName]
		if ok {
			delete(dag.Roots, rootName)
		}
	}

	for nodeName, node := range dag.Nodes {
		_, ok := pruneableNodes[nodeName]
		if ok {
			delete(dag.Nodes, nodeName)
		}

		var newNext []string
		for _, elem := range node.Next {
			_, ok := pruneableNodes[elem]
			if !ok {
				newNext = append(newNext, elem)
			}
		}
		node.Next = newNext

		var newPrev []string
		for _, elem := range node.Prev {
			_, ok := pruneableNodes[elem]
			if !ok {
				newPrev = append(newPrev, elem)
			}
		}
		node.Prev = newPrev
	}

	pruneableNodes[node] = dag.Nodes[node]
	dag.Pop(node)
	dag.Reconcile()
	return dag, pruneableNodes
}

func (p DefaultPruner) PruneAfter(
	node string,
	dag data.DAG,
) (data.DAG, map[string]data.Node) {
	p.l.Tracef("pruning nodes after %s", node)
	if len(dag.Nodes[node].Next) == 0 {
		p.l.Debugf("node %s had no children to prune", node)
		dag.Pop(node)
		dag.Reconcile()
		return dag, make(map[string]data.Node)
	}
	pruneableNodes := make(map[string]data.Node)

	// Fault has to be before this point.
	descendants := dag.Descendants(node)
	for descendantName, descendant := range descendants {
		pruneableNodes[descendantName] = descendant
		delete(dag.Nodes, descendantName)
		_, ok := dag.Roots[descendantName]
		if !ok {
			delete(dag.Roots, descendantName)
		}
	}

	newNode := dag.Nodes[node]
	newNode.Next = []string{}
	dag.Nodes[node] = newNode

	pruneableNodes[node] = dag.Nodes[node]
	dag.Pop(node)
	dag.Reconcile()
	return dag, pruneableNodes
}
