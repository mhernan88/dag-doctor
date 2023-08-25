package pruners

import(
    "slices"
    "github.com/sirupsen/logrus"
    "github.com/mhernan88/dag-bisect/data"
    "github.com/mhernan88/dag-bisect/utils"
    mapset "github.com/deckarep/golang-set/v2"
)

func NewDefaultPruner(iterationLimit int, l *logrus.Logger) DefaultPruner {
    return DefaultPruner{
        iterationLimit: iterationLimit,
        l: l,
    }
}

type DefaultPruner struct {
    iterationLimit int
    l *logrus.Logger
}

// Finds all nodes that can be pruned before `source`
// assuming `source` is error-free.
func (p DefaultPruner) findUpstreamPruneableNodes(source *data.Node) (map[string]*data.Node, []string) {
    p.l.Tracef("finding pruneable nodes before %s", source.Name)
    if source.Prev == nil {
        return make(map[string]*data.Node), []string{}
    }

    ancestorsMap := data.GetNodeAncestors([]*data.Node{source})
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


// Finds all nodes that can be pruned after `source`
// assuming `source` has an error.
func (p DefaultPruner) findDownstreamPruneableNodes(
    source *data.Node,
) (map[string]*data.Node, []string) {
    p.l.Tracef("finding pruneable nodes after %s", source.Name)
    if source.Next == nil {
        return make(map[string]*data.Node), []string{}
    }

    pruneableDescendantsMap := data.GetNodeDescendants([]*data.Node{source})
    p.l.Tracef("pulling pruneable nodes from %d descendants", len(pruneableDescendantsMap))

    pruneableDescendantNames, _ := data.SliceNodeMap(pruneableDescendantsMap)
    return pruneableDescendantsMap, pruneableDescendantNames
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
    roots []*data.Node,
) ([]string, error) {
    _, pruneableNodes := p.findUpstreamPruneableNodes(source)
    rootsMap := make(map[string]*data.Node)
    for _, root := range roots {
        rootsMap[root.Name] = root
    }

    p.unlinkNext(rootsMap, pruneableNodes)
    p.unlinkPrev(rootsMap, pruneableNodes)

    p.l.Tracef("pruned %v", pruneableNodes)
    return pruneableNodes, nil
}

func (p DefaultPruner) PruneAfter(
    source *data.Node, 
    roots []*data.Node,
) ([]string, error) {
    _, pruneableNodes := p.findDownstreamPruneableNodes(source)
    descendantsBefore := utils.FlattenAllNodeNames(source.Next)
    p.l.Tracef(
        "before pruning '%s' had %d descendants",
        source.Name, 
        len(descendantsBefore) + len(source.Next),
    )

    source.Next = nil
    
    descendantsAfter := utils.FlattenAllNodeNames(source.Next)
    p.l.Tracef(
        "after pruning '%s' had %d descendants",
        source.Name,
        len(descendantsAfter) + len(source.Next),
    )

    p.l.Tracef("pruned %v", pruneableNodes)
    return pruneableNodes, nil
}
