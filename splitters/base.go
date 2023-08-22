package splitters

import (
    "github.com/mhernan88/dag-bisect/data"
)

type Splitter interface {
    FindCandidate(roots map[string]*data.Node) (*data.Node, error)
}
