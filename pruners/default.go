package pruners

import (
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
		for childName, _ := range ancestor.Next {
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
			for childName, _ := range root.Next {
				if slices.Contains(targets, childName) {
					delete(root.Next, childName)
				}
			}
			for parentName, _ := range root.Prev {
				if slices.Contains(targets, parentName) {
					delete(root.Prev, parentName)
				}
			}
		}
	}
}

// Given a list of targets, unlinkPrev removes those
// from the DAG going backward.
func (p DefaultPruner) unlinkPrev(roots map[string]*data.Node, targets []string) {
	for _, root := range roots {
		if len(root.Next) > 0 {
			p.unlinkPrev(root.Prev, targets)
			for childName, _ := range root.Next {
				if slices.Contains(targets, childName) {
					delete(root.Next, childName)
				}
			}
			for parentName, _ := range root.Prev {
				if slices.Contains(targets, parentName) {
					delete(root.Prev, parentName)
				}
			}
		}
	}
}

func (p DefaultPruner) PruneBefore(
	source *data.Node,
	roots map[string]*data.Node,
) (map[string]*data.Node, []string, error) {
	p.l.Tracef("pruning nodes before %s", source.Name)
	_, pruneableNodes := p.findUpstreamPruneableNodes(source)

	p.unlinkNext(roots, pruneableNodes)
	p.unlinkPrev(roots, pruneableNodes)

	// p.l.Tracef("%d possible faulty nodes remaining", len(ancestors))
	return roots, pruneableNodes, nil
}

func (p DefaultPruner) PruneAfter(
	source *data.Node,
	roots map[string]*data.Node,
) (map[string]*data.Node, []string, error) {
	p.l.Tracef("pruning nodes after %s", source.Name)
	if source.Next == nil {
		return nil, []string{}, nil
	}

	pruneableDescendantsMap, _ := data.GetNodeDescendants([]*data.Node{source})
	pruneableDescendantNames, _ := data.SliceNodeMap(pruneableDescendantsMap)

	// Roots are pruned to only ancestor roots.
	// ancestors, roots := data.GetNodeAncestors([]*data.Node{source})
	// p.l.Tracef("%d possible faulty nodes remaining", len(ancestors))

	// Fault has to be before this point.
	source.Next = nil

	return roots, pruneableDescendantNames, nil
}
