package pruners

import(
    "fmt"
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
func (p DefaultPruner) findUpstreamPruneableNodes(source *data.Node) (map[string]*data.Node, []string, error) {
    p.l.Tracef("finding pruneable nodes before %s", source.Name)
    if source.Prev == nil {
        return make(map[string]*data.Node), []string{}, nil
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
    return pruneableAncestorsMap, pruneableAncestorNames.ToSlice(), nil
}


// Finds all nodes that can be pruned after `source`
// assuming `source` has an error.
func (p DefaultPruner) findDownstreamPruneableNodes(
    source *data.Node,
    roots[]*data.Node,
) ([]string, error) {
    p.l.Tracef("finding pruneable nodes after %s", source.Name)
    prunedNodes := mapset.NewSet[string]()
    if source.Next == nil {
        return prunedNodes.ToSlice(), nil
    }

    var keys []string
    nodes := make(map[string]*data.Node)
    for k, v := range source.Next {
        nodes[k] = v
        keys = append(keys, k)
        p.l.Tracef("evaluating node %s", k)
    }

    // Traverse Forwards 
    i := 0
    for len(keys) > 0 {
        key := keys[len(keys)-1]
        keys = keys[:len(keys)-1]
        node := nodes[key]

        for _, child := range node.Next {
            nodes[child.Name] = child
            keys = append(keys, child.Name)
            p.l.Tracef("|---> adding child %s to stack", child.Name)
        }
        p.l.Tracef("|---> adding node %s to prunedNodes", node.Name)
        prunedNodes.Add(node.Name)

        i++
        if p.iterationLimit > 0 {
            if i > p.iterationLimit {
                return nil, fmt.Errorf("reached iteration limit")
            }
        }
    }
    return prunedNodes.ToSlice(), nil
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
    _, pruneableNodes, err := p.findUpstreamPruneableNodes(source)
    if err != nil {
        return nil, err
    }

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
    pruneableNodes, err := p.findDownstreamPruneableNodes(source, roots)
    if err != nil {
        return nil, err
    }

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
