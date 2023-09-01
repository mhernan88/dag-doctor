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

func (s DefaultSplitter) FindCandidate(dag data.DAG) (data.Node, error) {
	s.l.Debug("selecting best split candidate")
	s.l.Tracef("options: %v", dag.SliceKeys())
	var candidate data.Node
	var candidateFound bool
	var bestScore = math.Inf(-1)

	var key string
	var keys []string
	for key = range dag.Roots {
		if key == "" {
			return data.Node{}, fmt.Errorf("dag contained blank key")
		}
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		return data.Node{}, fmt.Errorf("dag contained no keys")
	}

	var nd data.Node
	var ok bool

	for len(keys) > 0 {
		key = keys[len(keys)-1]
		keys = keys[:len(keys)-1]
		nd, ok = dag.Nodes[key]

		if key == "" {
			return data.Node{}, fmt.Errorf("found empty key in map")
		}

		if !ok {
			return data.Node{}, fmt.Errorf("failed to pull node %s from map", key)
		}

		isReconciled, labelsAndNotNodes, nodesAndNotlabels := dag.IsReconciled()
		if !isReconciled {
			return data.Node{}, fmt.Errorf(
				"dag not reconciled, excess labels = %v, excess nodes = %v",
				labelsAndNotNodes, nodesAndNotlabels,
			)
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
			candidate = nd
			candidateFound = true
		}

		for _, child := range nd.Next {
			s.l.Tracef("adding node %s child (%s) to stack", nd.Name, child)
			if child == "" {
				return data.Node{}, fmt.Errorf(
					"node '%s' child had corrupt name",
					nd.Name,
				)
			}

			if dag.Nodes[child].Name == "" {
				return data.Node{}, fmt.Errorf(
					"node '%s' child obj had corrupt name (key=%s)",
					nd.Name,
					child,
				)
			}
			keys = append(keys, dag.Nodes[child].Name)
		}
	}

	if candidateFound == false {
		return data.Node{}, fmt.Errorf("failed to select a split candidate")
	}
	s.l.Debugf("selected candidate '%s' has best split score: %f", candidate.Name, bestScore)
	return candidate, nil
}
