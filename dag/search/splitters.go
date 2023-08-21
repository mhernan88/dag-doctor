package search

import (
    "math"
    "slices"
    "github.com/mhernan88/dag-bisect/dag"
    "github.com/sirupsen/logrus"
)

// Container to Pair a NodeName and it's Score together.
type NodeScore struct {
    Name string
    Score float64
}

// Generic Splitter Interface
type Splitter interface {
    FindCandidate(roots []*dag.Node) (string, error)
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
func (s FamilyTreeSplitter) FindCandidate(roots []*dag.Node) (string, error) {
    var candidates []NodeScore

    var nd *dag.Node
    for len(roots) > 0 {
        nd = roots[len(roots)-1]
        roots = roots[:len(roots)-1]
        s.l.Tracef("popped node %s from queue", nd.Name)

        numAncestors := getNumAncestors(*nd, s.l)
        numDescendants := getNumDescendants(*nd, s.l)

        s.l.Tracef("calculating node %s split score", nd.Name)
        diff := math.Abs(float64(numAncestors - numDescendants))
        mean := float64(numAncestors + numDescendants) / 2
        score := mean - diff

        candidates = append(candidates, NodeScore{Name: nd.Name, Score: score})

        for _, child := range nd.Next {
            s.l.Tracef("adding node %s child (%s) to stack", nd.Name, child.Name)
            roots = append(roots, child)
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
    return candidates[len(candidates)-1].Name, nil
}
