package search

import (
    "errors"
    "github.com/mhernan88/dag-bisect/dag"
    "github.com/sirupsen/logrus"
)


// DistanceToStart calculates the longest path to the start of the dag.
func distanceToStart(node *dag.Node, limit int, l *logrus.Logger) (int, error) {
    l.Trace("calculating DistanceToStart()")

    // Edge cases where node is leaf node.
    if (node.Prev == nil) || (len(node.Prev) == 0) {
        l.Tracef("node %s has no parents, DistanceToStart() is 0", node.Name)
        return 0, nil
    }

    // Go Eqvuialent of Tuple (Node, int)
    var nodes []dag.NodeWrapper
    var nd dag.NodeWrapper

    // Assemble a queue of NodeWrappers
    for _, nodePtr := range node.Prev {
        if nodePtr == nil {
            continue
        }
        nodes = append(nodes, dag.NodeWrapper{Node: *nodePtr, Distance: 1})
    }

    maxDistance := 1
    i := 0

    // Queue To Iterate Over All Cases
    for len(nodes) > 0 {
        nd, nodes = nodes[len(nodes)-1], nodes[:len(nodes)-1]
        l.Tracef("popped node %s from queue", nd.Node.Name)

        for _, childNode := range nd.Node.Prev {
            l.Tracef("analyzing childNode %s", childNode.Name)
            maxDistance = max(nd.Distance + 1, maxDistance)
            childNodeWrapper := dag.NodeWrapper{Node: *childNode, Distance: nd.Distance + 1}
            nodes = append(nodes, childNodeWrapper)
            l.Tracef("calculated node %s distance and added it to queue", childNode.Name)
        }

        i++
        if limit > 0 {
            if i > limit{
                return 0, errors.New("DistanceToStart() exceeded iteration limit")
            }
        }
    }

    l.Tracef("node %s has DistanceToStart() of %d", node.Name, maxDistance)
    return maxDistance, nil
}

// DistanceToEnd calculates the longest path to the end of the dag.
func distanceToEnd(node *dag.Node, limit int, l *logrus.Logger) (int, error) {
    l.Trace("calculating DistanceToEnd()")

    // Edge cases where node is leaf node.
    if (node.Next == nil) || (len(node.Next) == 0) {
        l.Tracef("node %s has no children, DistanceToEnd() is 0", node.Name)
        return 0, nil
    }

    // Go Eqvuialent of Tuple (Node, int)
    var nodes []dag.NodeWrapper
    var nd dag.NodeWrapper

    // Assemble a queue of NodeWrappers
    for _, nodePtr := range node.Next {
        if nodePtr == nil {
            continue
        }
        nodes = append(nodes, dag.NodeWrapper{Node: *nodePtr, Distance: 1})
    }

    maxDistance := 1
    i := 0

    // Queue To Iterate Over All Cases
    for len(nodes) > 0 {
        nd, nodes = nodes[len(nodes)-1], nodes[:len(nodes)-1]
        l.Tracef("popped node %s from queue", nd.Node.Name)

        for _, childNode := range nd.Node.Next {
            l.Tracef("analyzing childNode %s", childNode.Name)
            maxDistance = max(nd.Distance + 1, maxDistance)
            childNodeWrapper := dag.NodeWrapper{Node: *childNode, Distance: nd.Distance + 1}
            nodes = append(nodes, childNodeWrapper)
            l.Tracef("calculated node %s distance and added it to queue", nd.Node.Name)
        }

        i++
        if limit > 0 {
            if i > limit {
                return 0, errors.New("DistanceToEnd() exceeded iteration limit")
            }
        }
    }

    l.Tracef("node %s has DistanceToEnd() of %d", node.Name, maxDistance)
    return maxDistance, nil
}
