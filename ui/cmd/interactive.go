package cmd

import (
    "os"
    "fmt"
    "github.com/mhernan88/dag-bisect/dag"
    "github.com/mhernan88/dag-bisect/dag/search"
    "github.com/sirupsen/logrus"
    "github.com/eiannone/keyboard"
    "github.com/enescakir/emoji"
)

func NewInspector(
    roots map[string]*dag.Node, 
    splitter search.Splitter,
    pruner search.Pruner,
    catalog map[string]dag.Dataset, 
    l *logrus.Logger,
) Inspector {
    return Inspector{
        roots: roots,
        goodNodes: []*dag.Node{},
        badNodes: []*dag.Node{},
        splitter: splitter,
        pruner: pruner,
        catalog: catalog,
        l: l,
    }
}

type Inspector struct {
    roots map[string]*dag.Node
    goodNodes []*dag.Node
    badNodes []*dag.Node
    splitter search.Splitter
    pruner search.Pruner
    catalog map[string]dag.Dataset
    l *logrus.Logger
}

func (i *Inspector) IsDatasetOK(dataset string) (bool, error) {
    fmt.Printf("----> inspecting dataset: %s\n\n", dataset)

    i.l.Tracef("requesting dataset %s from catalog", dataset)
    // pipelineFile, ok := i.catalog[dataset]
    // if !ok {
    //     return false, fmt.Errorf("dataset '%s' not found in catalog", dataset)
    // }
    // i.l.Tracef("successfully pulled dataset %s from catalog", dataset)
    //
    // i.l.Tracef("rendering dataset %s", dataset)
    // err := pipelineFile.LoadAndDisplay(10)
    // if err != nil {
    //     return false, err
    // }
    // i.l.Tracef("successfully rendered dataset %s", dataset)

    err := keyboard.Open()
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

func (i *Inspector) IsNodeOK(
    node *dag.Node, 
    roots map[string]*dag.Node,
) (*dag.Node, error) {
    var err error
    fmt.Printf("inspecting node: %s\n", node.Name)
    allOK := true
    for _, output := range node.Outputs {
        ok, err := i.IsDatasetOK(output)
        if err != nil {
            return nil, err
        }
        if ok == false {
            allOK = false
        }
    }

    var prunedNodes []string
    if !allOK {
        prunedNodes, err = i.pruner.PruneBefore(node, roots)
        if err != nil {
            return nil, err
        }

        fmt.Printf("node %s cleared OK\n", node.Name)
        fmt.Printf("|---> pruned nodes: %v", prunedNodes)
    } else {
        prunedNodes, err := i.pruner.PruneAfter(node)
        if err != nil {
            return nil, err
        }
        fmt.Printf("node %s has ERR\n", node.Name)
        fmt.Printf("|---> pruned nodes: %v", prunedNodes)
    }

    return node, nil
}

func (i *Inspector) IsDAGOK(
    iterationLimit int,
) (bool, error) {
    // var errNodes []string

    i.l.Trace("searching for available nodes")
    for len(i.roots) > 0 {
        node, err := i.splitter.FindCandidate(i.roots)
        if err != nil {
            return false, err
        }

        i.l.Trace("evaluating node")
        nodePtr, err := i.IsNodeOK(node, i.roots)
        if err != nil {
            return false, err
        }
        node = nodePtr

        // if !*node.IsValid {
        //     i.l.Tracef("appending node %s to errNode", node.Name)
        //     errNodes = append(errNodes, node.Name)
        // }
    }

    // if len(errNodes) == 0 {
    //     fmt.Println("all clear! no errors found in pipeline")
    // } else {
    //     fmt.Printf(
    //         "pipeline data issues were cuased by the following nodes: %v", 
    //         errNodes,
    //     )
    // }
    return true, nil
}
