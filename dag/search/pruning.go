package search

import(
    "fmt"
    "slices"
    "github.com/sirupsen/logrus"
    "github.com/mhernan88/dag-bisect/dag"
    mapset "github.com/deckarep/golang-set/v2"
)

type Pruner interface {
    PruneBefore(start *dag.Node) ([]string, error)
    PruneAfter(start *dag.Node) ([]string, error)
}

func NewDefaultPruner(iterationLimit int, l *logrus.Logger) DefaultPruner {
    return DefaultPruner{iterationLimit: iterationLimit, l: l}
}

type DefaultPruner struct {
    iterationLimit int
    l *logrus.Logger
}

// Given a list of targets, PruneForward removes those from the DAG going forward.
func (p DefaultPruner) PruneForward(roots map[string]*dag.Node, targets []string) {
    for _, root := range roots {
        if len(root.Next) > 0 {
            p.PruneForward(root.Next, targets)
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

// Given a list of targets, PruneForward removes those from the DAG going backward.
func (p DefaultPruner) PruneBackward(roots map[string]*dag.Node, targets []string) {
    for _, root := range roots {
        if len(root.Next) > 0 {
            p.PruneBackward(root.Prev, targets)
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
    start *dag.Node, 
    roots map[string]*dag.Node,
) ([]string, error) {
    p.l.Tracef("pruning nodes before %s", start.Name)
    prunedNodes := mapset.NewSet[string]()
    if start.Prev == nil {
        return prunedNodes.ToSlice(), nil
    }

    var keys []string
    var key string
    nodes := make(map[string]*dag.Node)
    for k, v := range start.Prev {
        nodes[k] = v
        keys = append(keys, k)
    }

    // Traverse Backwards
    i := 0
    for len(keys) > 0 {
        p.l.Tracef("picking from %d keys", len(keys))
        key = keys[len(keys)-1]
        keys = keys[:len(keys)-1]
        node := nodes[key]
        p.l.Tracef("evaluating node %s", key)

        if len(node.Next) == 1 {
            // If node only proceeds to pruneable node, then remove the entire node.
            for _, parent := range node.Prev {
                p.l.Tracef("---> adding parent %s to stack", parent.Name)
                nodes[parent.Name] = parent
                keys = append(keys, parent.Name)
            }
            p.l.Tracef("---> adding node %s to prunedNodes", node.Name)
            prunedNodes.Add(node.Name)
        } 

        i++
        if p.iterationLimit > 0 {
            if i > p.iterationLimit {
                return nil, fmt.Errorf("reached iteration limit")
            }
        }
    }

    p.PruneForward(roots, prunedNodes.ToSlice())

    var rootsList []*dag.Node
    for _, v := range roots {
        rootsList = append(rootsList, v)
    }
    leavesList := findLeaves(rootsList)

    leaves := make(map[string]*dag.Node)
    for _, leaf := range leavesList {
        leaves[leaf.Name] = leaf
    }

    p.PruneBackward(leaves, prunedNodes.ToSlice())
    return prunedNodes.ToSlice(), nil
}


func (p DefaultPruner) PruneAfter(start *dag.Node) ([]string, error) {
    p.l.Tracef("pruning nodes after %s", start.Name)
    prunedNodes := mapset.NewSet[string]()
    if start.Next == nil {
        return prunedNodes.ToSlice(), nil
    }

    var keys []string
    nodes := make(map[string]*dag.Node)
    for k, v := range start.Next {
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
            p.l.Tracef("---> adding child %s to stack", child.Name)
        }
        p.l.Tracef("---> adding node %s to prunedNodes", node.Name)
        prunedNodes.Add(node.Name)

        i++
        if p.iterationLimit > 0 {
            if i > p.iterationLimit {
                return nil, fmt.Errorf("reached iteration limit")
            }
        }
    }

    start.Next = nil
    return prunedNodes.ToSlice(), nil
}
