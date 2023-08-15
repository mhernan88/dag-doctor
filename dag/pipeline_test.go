package dag

import (
    "testing"
)

func TestLink(t *testing.T) {
    pipeline := Pipeline{
        Nodes: []PipelineNode{
            {Name: "Node1", Outputs: []string{"Node2"}},
            {Name: "Node2", Outputs: []string{"Node3"}},
            {Name: "Node3"},
        },
    }
    err := pipeline.Link()
    if err != nil {
        t.Fatalf("Failed to link nodes: %v", err)
    }
    if pipeline.Nodes[0].Next[0] != &pipeline.Nodes[1] || pipeline.Nodes[1].Prev[0] != &pipeline.Nodes[0] {
        t.Error("Nodes are not properly linked")
    }
}

func TestFind(t *testing.T) {
    pipeline := Pipeline{
        Nodes: []PipelineNode{{Name: "Node1"}},
    }
    node := pipeline.Find("Node1")
    if node != &pipeline.Nodes[0] {
        t.Error("Failed to find the correct node")
    }
}

func TestFindMidpoint(t *testing.T) {
    node1 := PipelineNode{Name: "Node1", Outputs: []string{"Node2"}}
    node2 := PipelineNode{Name: "Node2", Outputs: []string{"Node3"}}
    node3 := PipelineNode{Name: "Node3"}
    pipeline := Pipeline{
        Nodes: []PipelineNode{node1, node2, node3},
    }

    err := pipeline.Link()
    if err != nil {
        t.Fatalf("Failed to link nodes: %v", err)
    }

    midpoints := pipeline.FindMidpoint()
    if len(midpoints) != 1 || midpoints[0].Name != "Node2" {
        t.Error("Failed to find the correct midpoint node")
    }
}

