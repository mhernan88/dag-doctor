package dag

import "testing"

func TestDistanceToStart(t *testing.T) {
	// Creating a simple DAG structure
	node1 := &PipelineNode{Name: "Node1"}
	node2 := &PipelineNode{Name: "Node2", Prev: []*PipelineNode{node1}}
	node3 := &PipelineNode{Name: "Node3", Prev: []*PipelineNode{node2}}

	if dist := node1.DistanceToStart(); dist != 0 {
		t.Errorf("Expected distance to start: 0, got: %d", dist)
	}

	if dist := node2.DistanceToStart(); dist != 1 {
		t.Errorf("Expected distance to start: 1, got: %d", dist)
	}

	if dist := node3.DistanceToStart(); dist != 2 {
		t.Errorf("Expected distance to start: 2, got: %d", dist)
	}
}

func TestDistanceToEnd(t *testing.T) {
	// Creating a simple DAG structure
	node3 := &PipelineNode{Name: "Node3"}
	node2 := &PipelineNode{Name: "Node2", Next: []*PipelineNode{node3}}
	node1 := &PipelineNode{Name: "Node1", Next: []*PipelineNode{node2}}

	if dist := node3.DistanceToEnd(); dist != 0 {
		t.Errorf("Expected distance to end: 0, got: %d", dist)
	}

	if dist := node2.DistanceToEnd(); dist != 1 {
		t.Errorf("Expected distance to end: 1, got: %d", dist)
	}

	if dist := node1.DistanceToEnd(); dist != 2 {
		t.Errorf("Expected distance to end: 2, got: %d", dist)
	}
}

func TestDistance(t *testing.T) {
	// Creating a simple DAG structure
	node1 := &PipelineNode{Name: "Node1"}
	node2 := &PipelineNode{Name: "Node2", Prev: []*PipelineNode{node1}}
	node3 := &PipelineNode{Name: "Node3", Prev: []*PipelineNode{node2}, Next: []*PipelineNode{node1}}

	startDist, endDist := node3.Distance()
	if startDist != 1 || endDist != -1 {
		t.Errorf("Expected distance to start: 1 and end: -1, got start: %d, end: %d", startDist, endDist)
	}
}

