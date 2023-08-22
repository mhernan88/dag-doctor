package search

import (
    "math"
    "slices"
    "github.com/mhernan88/dag-bisect/dag"
    "github.com/sirupsen/logrus"
)

// Container to Pair a NodeName and it's Score together.
type NodeScore struct {
    Node *dag.Node
    Score float64
}

// Generic Splitter Interface
type Splitter interface {
    FindCandidate(roots map[string]*dag.Node) (*dag.Node, error)
}

func NewFamilyTreeSplitter(recursionLimit int, l *logrus.Logger) FamilyTreeSplitter {
    return FamilyTreeSplitter{
        recursionLimit: recursionLimit,
        l: l,
    }
}

// Splitter Based on Node Distance
type FamilyTreeSplitter struct {
    recursionLimit int
    l *logrus.Logger
}

// FindCandidate gets the best node to split a DAG along based on
// the number of ancestors and number of descendants.
func (s FamilyTreeSplitter) FindCandidate(roots map[string]*dag.Node) (*dag.Node, error) {
    var candidates []NodeScore
    var key string
    var keys []string
    for key, _ := range roots {
        keys = append(keys, key)
    }

    var nd *dag.Node
    for len(roots) > 0 {
        key = keys[len(keys)-1]
        keys = keys[:len(keys)-1]
        nd = roots[key]
        s.l.Tracef("popped node %s from queue", nd.Name)

        numAncestors := getNumAncestors(*nd, s.l)
        numDescendants := getNumDescendants(*nd, s.l)

        s.l.Tracef("calculating node %s split score", nd.Name)
        diff := math.Abs(float64(numAncestors - numDescendants))
        mean := float64(numAncestors + numDescendants) / 2
        score := mean - diff

        candidates = append(candidates, NodeScore{Node: nd, Score: score})

        for _, child := range nd.Next {
            s.l.Tracef("adding node %s child (%s) to stack", nd.Name, child.Name)
            keys = append(keys, child.Name)
            roots[child.Name] = child
        }
    }

    // Need to sort by NodeScore.Score to get best candidate.
    slices.SortFunc(
        candidates, 
        func(a, b NodeScore) int {
            if a.Score > b.Score {
                return 1
            } else if a.Score < b.Score {
                return -1
            } else {
                return 0
            }
        })
    return candidates[len(candidates)-1].Node, nil
}
