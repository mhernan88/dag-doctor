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
		l:              l,
	}
}

type DefaultSplitter struct {
	iterationLimit int
	l              *logrus.Logger
}

func (s DefaultSplitter) FindCandidate(dag data.DAG) (*data.Node, error) {
	var candidate *data.Node = nil
	var bestScore = math.Inf(-1)

	var key string
	var keys []string
	for key, _ = range dag.Roots {
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("failed to find dag keys")
	}

	var nd data.Node
	for len(keys) > 0 {
		key = keys[len(keys)-1]
		keys = keys[:len(keys)-1]
		nd = dag.Roots[key]
		s.l.Tracef("popped node %s from queue", nd.Name)

		numAncestors := len(dag.Ancestors(key))
		numDescendants := len(dag.Descendants(key))

		s.l.Tracef("calculating node %s split score", nd.Name)
		diff := math.Abs(float64(numAncestors - numDescendants))
		mean := float64(numAncestors+numDescendants) / 2
		score := mean - diff

		s.l.Tracef("node '%s' has split score: %f", key, score)
		if score > bestScore {
			bestScore = score
			candidate = &nd
		}

		for _, child := range nd.Next {
			s.l.Tracef("adding node %s child (%s) to stack", nd.Name, child)
			keys = append(keys, dag.Nodes[child].Name)
		}
	}

	if candidate == nil {
		return nil, fmt.Errorf("failed to select a split candidate")
	}
	s.l.Debugf("selected candidate '%s' has best split score: %f", candidate.Name, bestScore)
	return candidate, nil
}
