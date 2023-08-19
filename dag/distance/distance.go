package distance

import (
    "errors"
    "github.com/sirupsen/logrus"
)

func NewDistanceCalculator(l *logrus.Logger) *DistanceCalculator {
    return &DistanceCalculator{l}
}

type DistanceCalculator struct {
    l *logrus.Logger
}

func (dc DistanceCalculator) Midpoints(nodes []*Node, limit int) ([]string, error) {
    var midpoints []string

    var nd *Node
    for len(nodes) > 0 {
        nd, nodes = nodes[len(nodes)-1], nodes[:len(nodes)-1]
        dc.l.Tracef("popped node %s from queue", nd.Name)

        distanceToStart, err := dc.DistanceToStart(nd, limit)
        if err != nil {
            return nil, err
        }

        distanceToEnd, err := dc.DistanceToEnd(nd, limit)
        if err != nil {
            return nil, err
        }

        if distanceToStart == distanceToEnd {
            midpoints = append(midpoints, nd.Name) 
        }

        for _, child := range nd.Next {
            nodes = append(nodes, child)
        }
    }
    return midpoints, nil
}

// DistanceToStart calculates the longest path to the start of the dag.
func (dc DistanceCalculator) DistanceToStart(node *Node, limit int) (int, error) {
    dc.l.Trace("calculating DistanceToStart()")

    // Edge cases where node is leaf node.
    if (node.Prev == nil) || (len(node.Prev) == 0) {
        dc.l.Tracef("node %s has no parents, DistanceToStart() is 0", node.Name)
        return 0, nil
    }

    // Go Eqvuialent of Tuple (Node, int)
    var nodes []NodeWrapper
    var nd NodeWrapper

    // Assemble a queue of NodeWrappers
    for _, nodePtr := range node.Prev {
        if nodePtr == nil {
            continue
        }
        nodes = append(nodes, NodeWrapper{Node: *nodePtr, Distance: 1})
    }

    maxDistance := 1
    i := 0

    // Queue To Iterate Over All Cases
    for len(nodes) > 0 {
        nd, nodes = nodes[len(nodes)-1], nodes[:len(nodes)-1]
        dc.l.Tracef("popped node %s from queue", nd.Node.Name)

        for _, childNode := range nd.Node.Prev {
            dc.l.Tracef("analyzing childNode %s", childNode.Name)
            maxDistance = max(nd.Distance + 1, maxDistance)
            childNodeWrapper := NodeWrapper{Node: *childNode, Distance: nd.Distance + 1}
            nodes = append(nodes, childNodeWrapper)
            dc.l.Tracef("calculated node %s distance and added it to queue", childNode.Name)
        }

        i++
        if limit > 0 {
            if i > limit{
                return 0, errors.New("DistanceToStart() exceeded iteration limit")
            }
        }
    }

    dc.l.Tracef("node %s has DistanceToStart() of %d", node.Name, maxDistance)
    return maxDistance, nil
}

// DistanceToEnd calculates the longest path to the end of the dag.
func (dc DistanceCalculator) DistanceToEnd(node *Node, limit int) (int, error) {
    dc.l.Trace("calculating DistanceToEnd()")

    // Edge cases where node is leaf node.
    if (node.Next == nil) || (len(node.Next) == 0) {
        dc.l.Tracef("node %s has no children, DistanceToEnd() is 0", node.Name)
        return 0, nil
    }

    // Go Eqvuialent of Tuple (Node, int)
    var nodes []NodeWrapper
    var nd NodeWrapper

    // Assemble a queue of NodeWrappers
    for _, nodePtr := range node.Next {
        if nodePtr == nil {
            continue
        }
        nodes = append(nodes, NodeWrapper{Node: *nodePtr, Distance: 1})
    }

    maxDistance := 1
    i := 0

    // Queue To Iterate Over All Cases
    for len(nodes) > 0 {
        nd, nodes = nodes[len(nodes)-1], nodes[:len(nodes)-1]
        dc.l.Tracef("popped node %s from queue", nd.Node.Name)

        for _, childNode := range nd.Node.Next {
            dc.l.Tracef("analyzing childNode %s", childNode.Name)
            maxDistance = max(nd.Distance + 1, maxDistance)
            childNodeWrapper := NodeWrapper{Node: *childNode, Distance: nd.Distance + 1}
            nodes = append(nodes, childNodeWrapper)
            dc.l.Tracef("calculated node %s distance and added it to queue", nd.Node.Name)
        }

        i++
        if limit > 0 {
            if i > limit {
                return 0, errors.New("DistanceToEnd() exceeded iteration limit")
            }
        }
    }

    dc.l.Tracef("node %s has DistanceToEnd() of %d", node.Name, maxDistance)
    return maxDistance, nil
}
