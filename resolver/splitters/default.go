package splitters

import (
	"fmt"
	"math"
	"log/slog"

	"github.com/mhernan88/dag-bisect/models"
)

func NewDefaultSplitter() DefaultSplitter {
	return DefaultSplitter{
		Name: "default",
	}
}

type DefaultSplitter struct {
	Name string `json:"default"`
}

func (s DefaultSplitter) GetName() string {
	return s.Name
}

func (s DefaultSplitter) FindCandidate(
	dag models.DAG, 
	l *slog.Logger,
) (models.Node, error) {
	l.Info("selecting best split candidate")
	var candidate models.Node
	var candidateFound bool
	var bestScore = math.Inf(-1)

	var key string
	var keys []string
	for key = range dag.Roots {
		if key == "" {
			return models.Node{}, fmt.Errorf("dag contained blank key")
		}
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		return models.Node{}, fmt.Errorf("dag contained no keys")
	}

	var nd models.Node
	var ok bool

	for len(keys) > 0 {
		key = keys[len(keys)-1]
		keys = keys[:len(keys)-1]
		nd, ok = dag.Nodes[key]

		if key == "" {
			return models.Node{}, fmt.Errorf("found empty key in map")
		}

		if !ok {
			return models.Node{}, fmt.Errorf("failed to pull node %s from map", key)
		}

		numAncestors := len(dag.Ancestors(key))
		numDescendants := len(dag.Descendants(key))

		diff := math.Abs(float64(numAncestors - numDescendants))
		mean := float64(numAncestors+numDescendants) / 2
		score := mean - diff

		if score > bestScore {
			bestScore = score
			candidate = nd
			candidateFound = true
		}

		for _, child := range nd.Next {
			if child == "" {
				return models.Node{}, fmt.Errorf(
					"node '%s' child had corrupt name",
					nd.Name,
				)
			}

			if dag.Nodes[child].Name == "" {
				return models.Node{}, fmt.Errorf(
					"node '%s' child obj had corrupt name (key=%s)",
					nd.Name,
					child,
				)
			}
			keys = append(keys, dag.Nodes[child].Name)
		}
	}

	if candidateFound == false {
		return models.Node{}, fmt.Errorf("failed to select a split candidate")
	}
	l.Info(
		"selected candidate has best split score", 
		"candidate", candidate.Name, 
		"score", bestScore,
	)
	return candidate, nil
}
