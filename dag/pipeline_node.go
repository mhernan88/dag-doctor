package dag

type PipelineNode struct {
	Name             string          `json:"name"`
	Inputs           []string        `json:"inputs"`
	Outputs          []string        `json:"outputs"`
	Next             []*PipelineNode `json:"-"`
	Prev             []*PipelineNode `json:"-"`
}

func (node *PipelineNode) DistanceToStart() int {
	maxDist := -1
	for _, prevNode := range node.Prev {
		dist := prevNode.distanceToStart()
		if dist > maxDist {
			maxDist = dist
		}
	}
	return maxDist + 1
}

func (node *PipelineNode) DistanceToEnd() int {
	maxDist := -1
	for _, nextNode := range node.Next {
		dist := nextNode.distanceToEnd()
		if dist > maxDist {
			maxDist = dist
		}
	}
	return maxDist + 1
}

func (node *PipelineNode) Distance() (int, int) {
	beginningDist := node.distanceToStart() - 1 // Subtracting 1 to exclude the node itself
	endDist := node.distanceToEnd() - 1         // Subtracting 1 to exclude the node itself
	return beginningDist, endDist
}
