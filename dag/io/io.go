package io

import (
    "os"
    "fmt"
    "slices"
    "strings"
    goio "io"
    "encoding/json"
    mapset "github.com/deckarep/golang-set/v2"
    "github.com/mhernan88/dag-bisect/dag"
    "github.com/sirupsen/logrus"
)

func readNodes(filename string, l *logrus.Logger) (*dag.Pipeline, error) {
    l.Tracef("opening dag %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

    l.Tracef("reading bytes from dag %s", filename)
	bytes, err := goio.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

    l.Trace("marhsaling byte into Pipeline")
	var pipeline dag.Pipeline
	if err := json.Unmarshal(bytes, &pipeline); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

    l.Tracef("successfully loaded dag")
	return &pipeline, nil
}

func gatherDatasets(nodes map[string]dag.Node, l *logrus.Logger) (
    mapset.Set[string],
    mapset.Set[string],
) {
    // Get all Input Datasets
    l.Trace("gathering pipeline inputs")
    inputs := mapset.NewSet[string]()
    for _, node := range nodes {
        for _, input := range node.Inputs {
            inputs.Add(input)
        }
    }

    // Get all Output Datasets
    l.Trace("gathering pipeline outputs")
    outputs := mapset.NewSet[string]()
    for _, node := range nodes {
        for _, output := range node.Outputs {
            outputs.Add(output)
        }
    }

    return inputs, outputs
}

func findNonRoots(
    nodes map[string]dag.Node, 
    nonRootDatasets mapset.Set[string],
    l *logrus.Logger,
) mapset.Set[string] {
    // Find all nodes associated with non-Root Datasets
    l.Trace("finding non-root nodes")
    nonRootNodes := mapset.NewSet[string]()
    for _, nonRootDataset := range nonRootDatasets.ToSlice() {
        for name, node := range nodes {
            if slices.Contains(node.Inputs, nonRootDataset) {
                nonRootNodes.Add(name)
            }
        }
    }

    l.Tracef("found %d non-root nodes", len(nonRootNodes.ToSlice()))
    return nonRootNodes
}


func findRoots(
    nodes map[string]dag.Node,
    nonRootNodes mapset.Set[string],
    l *logrus.Logger,
) mapset.Set[string] {
    allNodes := mapset.NewSet[string]()
    for name, _ := range nodes {
        allNodes.Add(name)
    }
    rootNodes := allNodes.Difference(nonRootNodes)
    l.Tracef("found %d root nodes", len(rootNodes.ToSlice()))
    return rootNodes
}

func linkNode(root *dag.Node, nodes map[string]dag.Node) *dag.Node {
    // Look through all outputs...
    for _, output := range root.Outputs{
        // And then look through entire pipeline...
        for nodeName, node := range nodes {
            if nodeName == root.Name {
                continue
            }

            // Each node where root.Output == node.Input 
            // is an immediate descendant of root.
            if slices.Contains(node.Inputs, output) {
                nodePtr := linkNode(&node, nodes)
                if root.Next == nil {
                    root.Next = make(map[string]*dag.Node)
                }
                root.Next[nodePtr.Name] = nodePtr
            }
        }
    }
    return root
}

func backlinkNodes(
    nodes map[string]*dag.Node, 
    upstream *dag.Node,
    recursionDepth,
    maxRecursionDepth int,
    l *logrus.Logger,
) (map[string]*dag.Node, error) {
    if maxRecursionDepth > 0 {
        if recursionDepth > maxRecursionDepth {
            return nil, fmt.Errorf("backlinkNodes reached maxRecursionDepth")
        }
    }
    // For each node...
    for name, node := range nodes {
        l.Tracef("backlinking node %s", name)
        // If an upstream is provided, then add it to node.Prev
        if upstream != nil {
            if node.Prev == nil {
                node.Prev = make(map[string]*dag.Node)
            }
            nodes[name].Prev[upstream.Name] = upstream
        }

        if nodes[name].Next != nil {
            // Go to the next node and repeat.
            l.Tracef("recursively running backlinkNodes() on %v", nodes[name].Next)
            linkedNodePtrs, err := backlinkNodes(
                nodes[name].Next, 
                nodes[name],
                recursionDepth + 1,
                maxRecursionDepth,
                l,
            )
            if err != nil {
                return nil, err
            }
            for _, linkedNodePtr := range linkedNodePtrs {
                nodes[name].Next[linkedNodePtr.Name] = linkedNodePtr
            }
        }
    }
    return nodes, nil
}

func linkNodes(
    roots map[string]*dag.Node, 
    nodes map[string]dag.Node, 
    maxRecursionDepth int,
    l *logrus.Logger,
) (map[string]*dag.Node, error) {
    l.Trace("forward linking nodes")

    for name, root := range roots {
        roots[name] = linkNode(root, nodes)
    }

    l.Trace("backward linking nodes")
    roots, err := backlinkNodes(roots, nil, 0, maxRecursionDepth, l)
    if err != nil {
        return nil, err
    }

    l.Tracef("successfully linked nodes (%d roots)", len(roots))
    return roots, nil
}

func processNodes(
    pipeline *dag.Pipeline, 
    maxRecursionDepth int,
    l *logrus.Logger,
) (map[string]*dag.Node, error) {
    inputs, outputs := gatherDatasets(pipeline.Nodes, l)
    intersection := inputs.Intersect(outputs)

    nonRoots := findNonRoots(pipeline.Nodes, intersection, l)
    roots := findRoots(pipeline.Nodes, nonRoots, l)

    rootNodes := make(map[string]*dag.Node)
    for name, node := range pipeline.Nodes {
        if roots.Contains(name) {
            nodeCopy := node
            rootNodes[name] = &nodeCopy
        }
    }

    l.Tracef("linking roots: %s", strings.Join(roots.ToSlice(), ", "))
    nodes, err := linkNodes(rootNodes, pipeline.Nodes, maxRecursionDepth, l)
    if err != nil {
        return nil, err
    }

    return nodes, nil
}

func LoadAndProcessNodes(
    filename string, 
    maxRecursionDepth int,
    l *logrus.Logger,
) (map[string]*dag.Node, error) {
    pipeline, err := readNodes(filename, l)
    if err != nil {
        return nil, err
    }
    nodes, err := processNodes(pipeline, maxRecursionDepth, l)
    if err != nil {
        return nil, err
    }
    if len(nodes) == 0 {
        return nil, fmt.Errorf("0 nodes loaded!")
    }
    return nodes, nil
}

func LoadCatalog(filename string) (*map[string]dag.Dataset, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %v", err)
    }
    defer file.Close()

    bytes, err := goio.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read file: %v", err)
    }

    var catalog dag.Catalog
    if err := json.Unmarshal(bytes, &catalog); err != nil {
        return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
    }

    return &catalog.Datasets, nil
}
