package search

import(
    "fmt"
    "slices"
    "github.com/mhernan88/dag-bisect/dag"
)

func removeParent(child *dag.Node, parentName string) error {
    if child.Prev == nil {
        return nil
    }

    var parent *dag.Node
    for _, parentCandidate := range child.Prev {
        if parentCandidate.Name == parentName {
            parent = parentCandidate
        }
    }

    // Delete Parent from Grandparent (if Grandparent Exists)
    if parent.Prev != nil {
        for _, grandparentCandidate := range parent.Prev {
            if grandparentCandidate.Next == nil {
                return fmt.Errorf("expected grandparent %s to have children", grandparentCandidate.Name)
            }
            for i, grandparentChild := range grandparentCandidate.Next {
                if grandparentChild.Name == parentName {
                    grandparentCandidate.Next = append(grandparentCandidate.Next[:i], grandparentCandidate.Next[i+1:]...)
                }
            }
        }
    }

    // Delete Parent from Child
    for i, parent := range child.Prev {
        if parent.Name == parentName {
            child.Prev = append(
                child.Prev[:i],
                child.Prev[i+1:]...
            )
        }
    }
    return nil
}

func PruneBefore(node *dag.Node) (dag.Node, []*dag.Node) {
    // nodes := node.Prev

    // for len(nodes) > 0 {
    //     nd := nodes[len(nodes)-1]
    //     nodes = nodes[:len(nodes)-1]
    //
    //     if len(nd.Next) == 1 {
    //         // Remove pointers to this node
    //         nd.Prev.N
    //     }
    // }
    //
    // node.Prev = nil
    return *node, []*dag.Node{node}
}

func PruneAfter(node dag.Node, childrenToPrune []string) (dag.Node, []*dag.Node) {
    var updatedChildren []*dag.Node
    for _, child := range node.Next {
        if !slices.Contains(childrenToPrune, child.Name) {
            updatedChildren = append(updatedChildren, child)
        }
    }
    node.Next = updatedChildren
    return node, []*dag.Node{&node} // TODO: Fix this
}
