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

func gatherDatasets(nodes []dag.Node, l *logrus.Logger) (
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
    nodes []dag.Node, 
    nonRootDatasets mapset.Set[string],
    l *logrus.Logger,
) mapset.Set[string] {
    // Find all nodes associated with non-Root Datasets
    l.Trace("finding non-root nodes")
    nonRootNodes := mapset.NewSet[string]()
    for _, nonRootDataset := range nonRootDatasets.ToSlice() {
        for _, node := range nodes {
            if slices.Contains(node.Inputs, nonRootDataset) {
                nonRootNodes.Add(node.Name)
            }
        }
    }

    l.Tracef("found %d non-root nodes", len(nonRootNodes.ToSlice()))
    return nonRootNodes
}


func findRoots(
    nodes []dag.Node,
    nonRootNodes mapset.Set[string],
    l *logrus.Logger,
) mapset.Set[string] {
    allNodes := mapset.NewSet[string]()
    for _, node := range nodes {
        allNodes.Add(node.Name)
    }
    rootNodes := allNodes.Difference(nonRootNodes)
    l.Tracef("found %d root nodes", len(rootNodes.ToSlice()))
    return rootNodes
}

func linkNode(root *dag.Node, nodes []dag.Node) *dag.Node {
    // Look through all outputs...
    for _, output := range root.Outputs{
        // And then look through entire pipeline...
        for _, node := range nodes {
            if node.Name == root.Name {
                continue
            }

            // Each node where root.Output == node.Input 
            // is an immediate descendant of root.
            if slices.Contains(node.Inputs, output) {
                nodePtr := linkNode(&node, nodes)
                root.Next = append(root.Next, nodePtr)
            }
        }
    }
    return root
}

func backlinkNodes(nodes[]*dag.Node, upstream *dag.Node) []*dag.Node {
    // For each node...
    for i := range nodes {
        // If an upstream is provided, then add it to node.Prev
        if upstream != nil {
            nodes[i].Prev = append(nodes[i].Prev, upstream)
        }

        if nodes[i].Next != nil {
            // Go to the next node and repeat.
            nodes[i].Next = backlinkNodes(nodes[i].Next, nodes[i])
        }
    }
    return nodes
}

func linkNodes(roots []*dag.Node, nodes []dag.Node, l *logrus.Logger) []*dag.Node {
    l.Trace("forward linking nodes")

    for i := range roots {
        roots[i] = linkNode(roots[i], nodes)
    }

    l.Trace("backward linking nodes")
    roots = backlinkNodes(roots, nil)

    l.Tracef("successfully linked nodes (%d roots)", len(roots))
    return roots
}

func processNodes(pipeline *dag.Pipeline, l *logrus.Logger) ([]*dag.Node, error) {
    inputs, outputs := gatherDatasets(pipeline.Nodes, l)
    intersection := inputs.Intersect(outputs)

    nonRoots := findNonRoots(pipeline.Nodes, intersection, l)
    roots := findRoots(pipeline.Nodes, nonRoots, l)

    var rootNodes []*dag.Node
    for _, node := range pipeline.Nodes {
        if roots.Contains(node.Name) {
            nodeCopy := node
            rootNodes = append(rootNodes, &nodeCopy)
        }
    }

    l.Tracef("linking roots: %s", strings.Join(roots.ToSlice(), ", "))
    nodes := linkNodes(rootNodes, pipeline.Nodes, l)

    return nodes, nil
}

func LoadAndProcessNodes(filename string, l *logrus.Logger) ([]*dag.Node, error) {
    pipeline, err := readNodes(filename, l)
    if err != nil {
        return nil, err
    }
    return processNodes(pipeline, l)
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
