package dag

type Pipeline struct {
  Nodes []PipelineNode `json:"nodes"`
} 

func (p Pipeline) Link() error {
	nodeMap := make(map[string]*PipelineNode)
	for i := range p.Nodes {
		nodeMap[p.Nodes[i].Name] = &p.Nodes[i]
	}

	for i := range p.Nodes {
		p.Nodes[i].Next = []*PipelineNode{}
		p.Nodes[i].Prev = []*PipelineNode{}
	}

	for i := range p.Nodes {
		for _, output := range p.Nodes[i].Outputs {
			nextNode, exists := nodeMap[output]
			if !exists {
				return fmt.Errorf("output node not found: %s", output)
			}
			p.Nodes[i].Next = append(p.Nodes[i].Next, nextNode)
			nextNode.Prev = append(nextNode.Prev, &p.Nodes[i])
		}
	}
	return nil
}

func (p Pipeline) Find(name string) *PipelineNode {
	for i := range p.Nodes {
		if p.Nodes[i].Name == name {
			return &p.Nodes[i]
		}
	}
	return nil
}

func (p *Pipeline) calculateDistances() (map[string]int, map[string]int) {
	distanceToStart := make(map[string]int)
	distanceToEnd := make(map[string]int)

	for _, node := range p.Nodes {
		start, end := node.distance()
		distanceToStart[node.Name] = start
		distanceToEnd[node.Name] = end
	}

	return distanceToStart, distanceToEnd
}

// Finds and returns the nodes in the middle of the pipeline
func (p *Pipeline) FindMidpoint() []PipelineNode {
	distanceToStart, distanceToEnd := p.calculateDistances()
	var middleNodes []PipelineNode

	for _, node := range p.Nodes {
		if distanceToStart[node.Name] == distanceToEnd[node.Name] {
			middleNodes = append(middleNodes, node)
		}
	}

	return middleNodes
}
