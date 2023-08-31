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
	for key = range dag.Roots {
		if key == "" {
			return nil, fmt.Errorf("dag contained blank key")
		}
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("dag contained no keys")
	}

	var nd data.Node
	var ok bool

	for len(keys) > 0 {
		fmt.Printf("len(keys)=%d\n", len(keys))
		fmt.Println(keys)
		key = keys[len(keys)-1]
		keys = keys[:len(keys)-1]
		nd, ok = dag.Nodes[key]

		if key == "" {
			return nil, fmt.Errorf("found empty key in map")
		}

		if !ok {
			return nil, fmt.Errorf("failed to pull node %s from map", key)
		}

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
			if child == "" {
				return nil, fmt.Errorf(
					"node '%s' child had corrupt name",
					nd.Name,
				)
			}

			if dag.Nodes[child].Name == "" {
				return nil, fmt.Errorf(
					"node '%s' child obj had corrupt name (key=%s)",
					nd.Name,
					child,
				)
			}
			keys = append(keys, dag.Nodes[child].Name)
		}
	}

	if candidate == nil {
		return nil, fmt.Errorf("failed to select a split candidate")
	}
	s.l.Debugf("selected candidate '%s' has best split score: %f", candidate.Name, bestScore)
	return candidate, nil
}
