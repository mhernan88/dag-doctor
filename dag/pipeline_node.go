package dag

import (
    "fmt"
)

type PipelineNode struct {
	Name             string          `json:"name"`
	Inputs           []string        `json:"inputs"`
	Outputs          []string        `json:"outputs"`
	Next             []*PipelineNode `json:"-"`
	Prev             []*PipelineNode `json:"-"`
    IsValid *bool `json:"-"`
}

type PipelineNodeWrapper struct {
    Node PipelineNode
    Distance int
}

func (node *PipelineNode) DistanceToStart() int {
	maxDist := -1
	for _, prevNode := range node.Prev {
		dist := prevNode.DistanceToStart()
		if dist > maxDist {
			maxDist = dist
		}
	}
	return maxDist + 1
}

func (node *PipelineNode) DistanceToEnd() int {
    if (node.Next == nil) || (len(node.Next) == 0) {
        fmt.Printf("node %s next=nil\n", node.Name)
        return 0
    }

    // Go Eqvuialent of Tuple (PipelineNode, int)
    var nodes []PipelineNodeWrapper
    for _, nodePtr := range node.Next {
        nodes = append(nodes, PipelineNodeWrapper{Node: *nodePtr, Distance: 1})
    }

    maxDistance := 1

    // Queue To Iterate Over All Cases
    fmt.Printf("len(nodes)=%d\n", len(nodes))
    for len(nodes) > 0 {
        nd, nodes := nodes[len(nodes)-1], nodes[:len(nodes)-1]
        fmt.Printf("popping node %s from stack\n", nd.Node.Name)

        for _, childNode := range nd.Node.Next {
            fmt.Printf("adding childNode of %s to queue; queue len=%d\n", childNode.Name, len(nodes))
            maxDistance = max(nd.Distance + 1, maxDistance)
            childNodeWrapper := PipelineNodeWrapper{Node: *childNode, Distance: nd.Distance + 1}
            nodes = append(nodes, childNodeWrapper)
        }
    }
    fmt.Printf("node %s has dist of %d\n", node.Name, maxDistance)
    return maxDistance
}

// func (node *PipelineNode) DistanceToEnd() int {
// 	maxDist := -1
// 	for _, nextNode := range node.Next {
// 		dist := nextNode.DistanceToEnd()
// 		if dist > maxDist {
// 			maxDist = dist
// 		}
// 	}
// 	return maxDist + 1
// }

func (node *PipelineNode) Distance() (int, int) {
	beginningDist := node.DistanceToStart() 
	endDist := node.DistanceToEnd()
	return beginningDist, endDist
}

func (node *PipelineNode) InputsValid() bool {
    for _, prevNode := range node.Prev {
        if prevNode.IsValid == nil || *prevNode.IsValid == false {
            return false
        }
    }
    return true
}
