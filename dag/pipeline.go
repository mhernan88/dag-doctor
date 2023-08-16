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
				continue
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
		start, end := node.Distance()
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

// traverseBefore recursively marks the ancestor nodes of the given node
// to be pruned unless they are also dependencies for other nodes.
func (p *Pipeline) traverseBefore(node *PipelineNode, toBePruned map[string]bool) {
	if node == nil || toBePruned[node.Name] {
		return
	}
	toBePruned[node.Name] = true
	for _, prevNode := range node.Prev {
		// If this previous node is also a dependency for other nodes, skip it
		if len(prevNode.Next) > 1 {
			continue
		}
		p.traverseBefore(prevNode, toBePruned)
	}
}

// traverseAfter recursively marks all the descendant nodes of the given node to be pruned.
func (p *Pipeline) traverseAfter(node *PipelineNode, toBePruned map[string]bool) {
	if node == nil || toBePruned[node.Name] {
		return
	}
	toBePruned[node.Name] = true
	for _, nextNode := range node.Next {
		p.traverseAfter(nextNode, toBePruned)
	}
}

// PruneNodes prunes nodes in the pipeline based on the given flag and target node.
// If the boolean flag is true, all nodes before the provided node are pruned,
// unless the ancestor node is a dependency for another branch.
// If the boolean flag is false, all of the descendant nodes are pruned from the DAG.
func (p *Pipeline) PruneNodes(pruneBefore bool, targetNode *PipelineNode) {
	toBePruned := map[string]bool{}

	if pruneBefore {
		p.traverseBefore(targetNode, toBePruned)
	} else {
		p.traverseAfter(targetNode, toBePruned)
	}

	newNodes := []PipelineNode{}
	for _, node := range p.Nodes {
		if !toBePruned[node.Name] {
			newNodes = append(newNodes, node)
		}
	}
	p.Nodes = newNodes
}

