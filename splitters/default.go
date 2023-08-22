package splitters

import (
    "fmt"
    "math"
    "github.com/mhernan88/dag-bisect/data"
    "github.com/sirupsen/logrus"
)

func NewDefaultSplitter(iterationLimit int, l *logrus.Logger) DefaultSplitter {
    return DefaultSplitter{
        iterationLimit: iterationLimit,
        l:l,
    }
}

type DefaultSplitter struct {
    iterationLimit int
    l *logrus.Logger
}


func (s DefaultSplitter) FindCandidate(roots map[string]*data.Node) (*data.Node, error) {
    var candidate *data.Node = nil
    var bestScore = math.Inf(-1)

    var key string
    var keys []string
    for key, _ = range roots {
        keys = append(keys, key)
    }

    var nd *data.Node
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

        if score > bestScore {
            bestScore = score
            candidate = nd
        }

        for _, child := range nd.Next {
            s.l.Tracef("adding node %s child (%s) to stack", nd.Name, child.Name)
            keys = append(keys, child.Name)
            roots[child.Name] = child
        }
    }

    if candidate == nil {
        return nil, fmt.Errorf("failed to select a split candidate")
    }
    return candidate, nil
}
