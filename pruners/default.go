package pruners

import (
	"fmt"
	"slices"

	mapset "github.com/deckarep/golang-set/v2"
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
func (p DefaultPruner) findUpstreamPruneableNodes(source *data.Node) (map[string]*data.Node, []string) {
	p.l.Tracef("finding pruneable nodes before %s", source.Name)
	if source.Prev == nil {
		return make(map[string]*data.Node), []string{}
	}

	ancestorsMap, _ := data.GetNodeAncestors([]*data.Node{source})
	p.l.Tracef("pulling pruneable nodes from %d ancestors", len(ancestorsMap))

	pruneableAncestorsMap := make(map[string]*data.Node)
	pruneableAncestorNames := mapset.NewSet[string]()
	for ancestorName, ancestor := range ancestorsMap {
		isPruneable := true
		for childName := range ancestor.Next {
			_, ok := ancestorsMap[childName]
			if !ok {
				p.l.Tracef("%s is not pruneable (it is not an ancestor)", ancestorName)
				isPruneable = false
				break
			}
		}

		if !isPruneable {
			p.l.Tracef("%s is pruneable (it is an ancestor)", ancestorName)
			pruneableAncestorsMap[ancestorName] = ancestor
			pruneableAncestorNames.Add(ancestorName)
		}
	}

	// Remaining nodes are ones we can safely prune.
	return pruneableAncestorsMap, pruneableAncestorNames.ToSlice()
}

// Given a list of targets, unlinkNext removes those
// from the DAG going forward.
func (p DefaultPruner) unlinkNext(roots map[string]*data.Node, targets []string) {
	for _, root := range roots {
		if len(root.Next) > 0 {
			p.unlinkNext(root.Next, targets)
			for childName := range root.Next {
				if slices.Contains(targets, childName) {
					delete(root.Next, childName)
				}
			}
			for parentName := range root.Prev {
				if slices.Contains(targets, parentName) {
					delete(root.Prev, parentName)
				}
			}
		}
	}
}

// Given a list of targets, unlinkPrev removes those
// from the DAG going backward.
// func (p DefaultPruner) unlinkPrev(leaves map[string]*data.Node, targets []string) {
// 	for _, leaf := range leaves {
// 		if len(leaf.Prev) > 0 {
// 			p.unlinkPrev(leaf.Prev, targets)
// 			for childName := range leaf.Next {
// 				if slices.Contains(targets, childName) {
// 					delete(leaf.Next, childName)
// 				}
// 			}
// 			for parentName := range leaf.Prev {
// 				if slices.Contains(targets, parentName) {
// 					delete(leaf.Prev, parentName)
// 				}
// 			}
// 		}
// 	}
// }

func (p DefaultPruner) PruneBefore(
	source *data.Node,
	roots map[string]*data.Node,
) (map[string]*data.Node, []string, error) {
	p.l.Tracef("pruning nodes before %s", source.Name)
	pruneableNodes, pruneableNodeNames := p.findUpstreamPruneableNodes(source)

	numNodesBeforePrune := len(data.UniqueNodes(roots).ToSlice())

	p.unlinkNext(roots, pruneableNodeNames)
	// p.unlinkPrev(roots, pruneableNodeNames)

	newRoots := make(map[string]*data.Node)
	for rootName, root := range roots {
		_, ok := pruneableNodes[rootName]
		if !ok {
			// Root was not in pruneableNodes. Retain.
			p.l.Tracef("retained root node %s", rootName)
			newRoots[rootName] = root
		} else {
			for childName, child := range root.Next {
				newRoots[childName] = child
			}
			p.l.Tracef("pruned root node %s", rootName)
		}
	}

	numNodesAfterPrune := len(data.UniqueNodes(newRoots).ToSlice())

	if numNodesBeforePrune == numNodesAfterPrune {
		return nil, nil, fmt.Errorf("no nodes were pruned")
	}

	p.l.Debugf(
		"num nodes before pruning after %d; num nodes after pruning after %d",
		numNodesBeforePrune,
		numNodesAfterPrune,
	)

	// p.l.Tracef("%d possible faulty nodes remaining", len(ancestors))
	return newRoots, pruneableNodeNames, nil
}

func (p DefaultPruner) PruneAfter(
	source *data.Node,
	roots map[string]*data.Node,
) ([]string, error) {
	p.l.Tracef("pruning nodes after %s", source.Name)
	if source.Next == nil {
		return []string{}, nil
	}

	pruneableDescendantsMap, _ := data.GetNodeDescendants([]*data.Node{source})
	pruneableDescendantNames, _ := data.SliceNodeMap(pruneableDescendantsMap)

	// Roots are pruned to only ancestor roots.
	// ancestors, roots := data.GetNodeAncestors([]*data.Node{source})
	// p.l.Tracef("%d possible faulty nodes remaining", len(ancestors))

	numNodesBeforePrune := len(data.UniqueNodes(roots).ToSlice())

	// Fault has to be before this point.
	source.Next = nil

	numNodesAfterPrune := len(data.UniqueNodes(roots).ToSlice())

	if numNodesBeforePrune == numNodesAfterPrune {
		return nil, fmt.Errorf("no nodes were pruned")
	}

	p.l.Debugf(
		"num nodes before pruning after %d; num nodes after pruning after %d",
		numNodesBeforePrune,
		numNodesAfterPrune,
	)

	return pruneableDescendantNames, nil
}
