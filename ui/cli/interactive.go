package distance

import (
    "fmt"
    "github.com/mhernan88/dag-bisect/dag"
    "github.com/mhernan88/dag-bisect/dag/search"
    "github.com/sirupsen/logrus"
)

func NewInspector(
    roots []*dag.Node, 
    splitter search.Splitter,
    catalog map[string]dag.Dataset, 
    l *logrus.Logger
) Inspector {
    return Inspector{
        roots: roots,
        goodNodes: []*dag.Node{},
        badNodes: []*dag.Node{},
        splitter: search.Splitter,
        catalog: catalog,
        l: l,
    }
}

type Inspector struct {
    roots []*dag.Node
    goodNodes []*dag.Node
    badNodes []*dag.Node
    splitter search.Splitter
    catalog map[string]dag.Dataset
    l *logrus.Logger
}

func (i *Inspector) IsDatasetOK(dataset string) (bool, error) {
    fmt.Printf("----> inspecting dataset: %s\n\n", dataset)

    i.l.Tracef("requesting dataset %s from catalog", dataset)
    pipelineFile, ok := catalog[dataset]
    if !ok {
        return false, fmt.Errorf("dataset '%s' not found in catalog", dataset)
    }
    i.l.Tracef("successfully pulled dataset %s from catalog", dataset)

    i.l.Tracef("rendering dataset %s", dataset)
    err := pipelineFile.LoadAndDisplay(10)
    if err != nil {
        return false, err
    }
    i.l.Tracef("successfully rendered dataset %s", dataset)

    err = keyboard.Open()
    if err != nil {
        return false, err
    }
    defer keyboard.Close()

    for {
        fmt.Println("\n----> output correct? (y/n):")
        i.l.Trace("reading keyboard input")
        char, _, err := keyboard.GetKey()
        if err != nil {
            return false, err
        }

        switch char {
        case 'y', 'Y':
            fmt.Printf(
                "----> %v dataset '%s' maked OK\n", 
                emoji.CheckMarkButton, 
                dataset)
            i.l.Trace("input was 'y', returning true, nil")
            return true, nil
        case 'n', 'N':
            fmt.Printf(
                "----> %v dataset '%s' maked ERR\n", 
                emoji.CrossMarkButton, 
                dataset)
            i.l.Trace("input was 'n', returning false, nil")
            return false, nil
        case 'c', 'C', 'q', 'Q':
            os.Exit(0)
        default:
            fmt.Printf("----> invalid input; options are: y, Y, n, N")
        }
    }
}

func (i *Inspector) IsNodeOK(node *dag.Node) (bool, error) {
) (*dag.Dataset, error) {
    fmt.Printf("inspecting node: %s\n", node.Name)
    for _, output := range node.Outputs {
        ok, err := IsDatasetOK(output)
        if err != nil {
            return nil, err
        }

        if !ok {
            i.
        }

        var isValid bool
        if !ok {
            isValid = false
            node.IsValid = &isValid
            i.pipeline.PruneNodes(node)
        } else {
            isValid = true
            node.IsValid = &isValid
            i.pipeline.PruneNodes(node)
        }
    }
    fmt.Printf("node %s cleared OK\n", node.Name)
    return node, nil
}

func (i *Inspector) IsDAGOK(
    nodes []dag.Node,
    catalog map[string]dag.Dataset,
    l *logrus.Logger,
) (bool, error) {
    var errNodes []string

    i.l.Trace("searching for available nodes")
    for !i.pipeline.AllNodesInspected() {
        // midpoints := i.pipeline.FindMidpoint()
        node := midpoints[0]

        i.l.Trace("evaluating node")
        nodePtr, err := i.IsNodeOK(&node)
        if err != nil {
            return false, err
        }
        node = *nodePtr

        if !*node.IsValid {
            i.l.Tracef("appending node %s to errNode", node.Name)
            errNodes = append(errNodes, node.Name)
        }
    }

    if len(errNodes) == 0 {
        fmt.Println("all clear! no errors found in pipeline")
    } else {
        fmt.Printf(
            "pipeline data issues were cuased by the following nodes: %v", 
            errNodes,
        )
    }
    return true, nil
}
