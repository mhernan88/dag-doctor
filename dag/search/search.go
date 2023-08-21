package search

import (
    "github.com/mhernan88/dag-bisect/dag"
    mapset "github.com/deckarep/golang-set/v2"
)

// Given root nodes, finds all leaf nodes.
func findLeaves(roots []*dag.Node) []*dag.Node {
    var leaves []*dag.Node

    var root *dag.Node
    for len(roots) > 0 {
        root = roots[len(roots)-1]
        roots = roots[:len(roots)-1]

        if len(root.Next) == 0 {
            leaves = append(leaves, root)
        } else {
            for _, nextNode := range root.Next {
                roots = append(roots, nextNode)
            }
        }
    }

    seen := mapset.NewSet[string]()
    var dedupLeaves []*dag.Node
    for _, leaf := range leaves {
        if seen.Contains(leaf.Name) {
            continue
        } else {
            seen.Add(leaf.Name)
            dedupLeaves = append(dedupLeaves, leaf)
        }
    }
    return dedupLeaves
}
