package io

import (
    "slices"
    "strings"
    "testing"
    "github.com/sirupsen/logrus"
    "github.com/mhernan88/dag-bisect/dag"
)

func setupNodeList() map[string]dag.Node {
    out := make(map[string]dag.Node)
    out["node1"] = dag.Node{
        Name: "node1",
        Inputs: []string{"dataset1", "dataset2"},
        Outputs: []string{"dataset3", "dataset4"},
        Prev: make(map[string]*dag.Node),
        Next: make(map[string]*dag.Node),
    }
    out["node2"] = dag.Node{
        Name: "node2",
        Inputs: []string{"dataset2", "dataset5"},
        Outputs: []string{"dataset6", "dataset7"},
        Prev: make(map[string]*dag.Node),
        Next: make(map[string]*dag.Node),
    }
    out["node3"] = dag.Node{
        Name: "node3",
        Inputs: []string{"dataset3", "dataset6"},
        Outputs: []string{"dataset8"},
        Prev: make(map[string]*dag.Node),
        Next: make(map[string]*dag.Node),
    }
    out["node4"] = dag.Node{
        Name: "node4",
        Inputs: []string{"dataset1", "dataset4"},
        Outputs: []string{"dataset9"},
        Prev: make(map[string]*dag.Node),
        Next: make(map[string]*dag.Node),
    }
    out["node5"] = dag.Node{
        Name: "node5",
        Inputs: []string{"dataset7"},
        Outputs: []string{"dataset10"},
        Prev: make(map[string]*dag.Node),
        Next: make(map[string]*dag.Node),
    }
    return out
}

func TestFindNonRoots(t *testing.T) {
    l := logrus.New()
    l.SetLevel(logrus.TraceLevel)

    nodes := setupNodeList()

    inputs, outputs := gatherDatasets(nodes, l)
    intersection := inputs.Intersect(outputs)
    nonRoots := findNonRoots(nodes, intersection, l).ToSlice()
    if len(nonRoots) != 3 {
        t.Errorf("expected 3 non-roots, got %d", len(nonRoots))
    }
    nonRootsString := strings.Join(nonRoots, ", ")

    for _, expectedNonRoot := range []string{"node3", "node4", "node5"} {
        if !slices.Contains(nonRoots, expectedNonRoot) {
            t.Errorf(
                "expected '%s' in nonRoots, but got '%s'",
                expectedNonRoot, 
                nonRootsString,
            )
        }
    }
}

func TestFindRoots(t *testing.T) {
    l :=  logrus.New()
    l.SetLevel(logrus.TraceLevel)

    nodes := setupNodeList()

    inputs, outputs := gatherDatasets(nodes, l)
    intersection := inputs.Intersect(outputs)
    nonRoots := findNonRoots(nodes, intersection, l)

    roots := findRoots(nodes, nonRoots, l).ToSlice()
    if len(roots) != 2 {
        t.Errorf("expected 2 roots, got %d", len(roots))
    }
    
    rootsString := strings.Join(roots, ", ")

    for _, expectedRoot := range []string{"node1", "node2"} {
        if !slices.Contains(roots, expectedRoot) {
            t.Errorf(
                "expected '%s' in roots, but got '%s'",
                expectedRoot, 
                rootsString,
            )
        }
    }
}


func TestProcessNodes(t *testing.T) {
    l := logrus.New()
    l.SetLevel(logrus.TraceLevel)

    pipeline := dag.Pipeline{Nodes: setupNodeList()}
    maxRecursionDepth := 99

    roots, err := processNodes(&pipeline, maxRecursionDepth, l)
    if err != nil {
        t.Error(err)
    }
    if len(roots) != 2 {
        t.Errorf("expected 2 roots, got %d", len(roots))
    }

    var names []string
    for _, root := range roots {
        names = append(names, root.Name)
    }

    rootsString := strings.Join(names, ", ")

    for _, expectedRoot := range []string{"node1", "node2"} {
        if !slices.Contains(names, expectedRoot) {
            t.Errorf(
                "expected '%s' in roots, but got '%s'",
                expectedRoot, 
                rootsString,
            )
        }
    }
}
