package pruners

import(
    "github.com/sirupsen/logrus"
    "github.com/mhernan88/dag-bisect/data"
)

func NewDefaultPruner(iterationLimit int, l *logrus.Logger) DefaultPruner {
    return DefaultPruner{
        iterationLimit: iterationLimit,
        l: l,
    }
}

type DefaultPruner struct {
    iterationLimit int
    l *logrus.Logger
}

func (p DefaultPruner) PruneBefore(
    source *data.Node, 
    roots[]*data.Node,
) ([]string, error) {
    return nil, nil
}

func (p DefaultPruner) PruneAfter(
    source *data.Node, 
    roots[]*data.Node,
) ([]string, error) {
    return nil, nil
}
