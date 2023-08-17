package dag

type PipelineNode struct {
	Name             string          `json:"name"`
	Inputs           []string        `json:"inputs"`
	Outputs          []string        `json:"outputs"`
	Next             []*PipelineNode `json:"-"`
	Prev             []*PipelineNode `json:"-"`
    IsValid *bool `json:"-"`
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
	maxDist := -1
	for _, nextNode := range node.Next {
		dist := nextNode.DistanceToEnd()
		if dist > maxDist {
			maxDist = dist
		}
	}
	return maxDist + 1
}

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
