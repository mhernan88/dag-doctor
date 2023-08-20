package dag

// import (
//     "strings"
//     "github.com/sirupsen/logrus"
// )
//
// func NewPipeline(nodes []PipelineNode, l *logrus.Logger) *Pipeline {
//     return &Pipeline{
//         Nodes: nodes,
//         l: l,
//     }
// }

// type Pipeline struct {
//   Nodes []PipelineNode `json:"nodes"`
//   l *logrus.Logger `json:"-"`
// } 
//
// func (p *Pipeline) GetNodes() []PipelineNode {
//     return p.Nodes
// }
//
// func (p *Pipeline) SetLogger(l *logrus.Logger) {
//     p.l = l
// }
//
// func (p Pipeline) FindRoots(
//     nodeInputs map[string][]*PipelineNode,
//     nodeOutputs map[string][]*PipelineNode,
// ) []string {
//     var rawDataInputs map[string]interface{}
//     for inputDataSet, inputNodeSlice := range nodeInputs {
//         inputDataSetFoundInOutputDataSets := false
//
//         for outputDataSet, _ := range nodeOutputs {
//             if (inputDataSet == outputDataSet) {
//                 inputDataSetFoundInOutputDataSets = true
//             }
//             break
//         }
//
//         if !inputDataSetFoundInOutputDataSets {
//             for _, inputNode := range inputNodeSlice {
//                 var x interface{}
//                 rawDataInputs[inputNode.Name] = x
//             }
//         }
//     }
//
//     var rawInputNodeSlice []string
//     for rawInputNode, _ := range rawDataInputs {
//         rawInputNodeSlice = append(rawInputNodeSlice, rawInputNode)
//     }
//     return rawInputNodeSlice
// }
//
// func (p Pipeline) FindLeaves(
//     nodeInputs map[string]*PipelineNode,
//     nodeOutputs map[string]*PipelineNode,
// ) []string {
//     var terminalDataOutputs []string
//     for _, nodeOutput := range nodeOutputs {
//         nodeOutputInNodeInputs := false
//         for _, nodeInput := range nodeInputs {
//             if (nodeOutput == nodeInput) {
//                 nodeOutputInNodeInputs = true
//             }
//             break
//         }
//         if !nodeOutputInNodeInputs {
//             terminalDataOutputs = append(terminalDataOutputs, nodeOutput.Name)
//         }
//     }
//     return terminalDataOutputs
// }
//
// // Returns up-to-date hashmap, and a boolean indicating whether an update happened.
// func (p Pipeline) AddToMap(
//     key string, 
//     nd *PipelineNode, 
//     hashmap map[string][]*PipelineNode,
// ) (
//     map[string][]*PipelineNode,
//     bool,
// ) {
//     _, exists := hashmap[key]
//     if !exists {
//         hashmap[key] = []*PipelineNode{}
//     }
//
//     values := hashmap[key]
//     for _, value := range values {
//         if value.Name == nd.Name {
//             return hashmap, false
//         }
//     }
//     values = append(values, nd)
//     hashmap[key] = values
//     return hashmap, true
// }
//
// // Creates map[<input_name>][]<PipelineNodes that use input_name>
// func (p Pipeline) GetNodeInputsMap() map[string][]*PipelineNode {
//     nodeInputs := make(map[string][]*PipelineNode)
//
//     var ok bool
//     for _, nd := range p.Nodes {
//         inputs := nd.Inputs
//         for _, input := range inputs {
//             nodeInputs, ok = p.AddToMap(input, &nd, nodeInputs)
//             if !ok {
//                 p.l.Tracef("skipped adding duplicate input %s", nd.Name)
//             }
//         }
//     }
//     return nodeInputs
// }
//
// // Creates map[<output_name>][]<PipelineNodes that use output_name>
// func (p Pipeline) GetNodeOutputsMap() map[string][]*PipelineNode {
//     nodeOutputs := make(map[string][]*PipelineNode)
//
//     var ok bool
//     for _, nd := range p.Nodes {
//         outputs := nd.Outputs
//         for _, output := range outputs {
//             nodeOutputs, ok = p.AddToMap(output, &nd, nodeOutputs)
//             if !ok {
//                 p.l.Tracef("skipped adding duplicate output %s", nd.Name)
//             }
//         }
//     }
//     return nodeOutputs
// }
//
// // USE STACKS AND QUEUES OR RECURSION TO TRAVERSE INSTEAD
// func (p Pipeline) Link() error {
//     // Maps dataset name to all nodes that use it as an input.
//     nodeInputs := p.GetNodeInputsMap()
//     // Maps dataset name to all nodes that use it as an output.
//     nodeOutputs := p.GetNodeOutputsMap()
//
//     roots := p.FindRoots(nodeInputs, nodeOutputs)
//     var rootNodes []*PipelineNode
//     for _, root := range roots {
//         rootNodes = append(rootNodes, p.Find(root))
//     }
//     // leaves := FindLeaves(nodeInputs, nodeOutputs)
//
//
//     //MAYBE START AT LEAVES AND GO UP TO ROOT, then we can just set p.Nodes to our queue.
//     // Queue
//     for len(rootsNodes) > 0 {
//         rootNode, rootNodes := rootNodes[len(rootNodes)-1], rootNodes[:len(rootNodes)-1] 
//         for _, output := range rootNode.Outputs {
//             // Mapping outputs of this node to all of the nodes that use them as inputs.
//             nodeChildren, ok := nodeInputs[output]
//             if !ok {
//                 // In this case, we're dealing witha  leaf node.
//                 // Logic to check if node is actually in "leaves"
//             }
//
//             // The Next items are just the children.
//             rootNode.Next = nodeChildren
//
//             for _, child := range nodeChildren {
//                 // Need to link forward and backwards.
//                 child.Prev = rootNode
//                 // Add children to the queue to do the same to them.
//                 rootNodes = append(rootNodes, child)
//             }
//
//         }
//
//     }
// }
//
// func (p Pipeline) Link() error {
// 	nodeMap := make(map[string]*PipelineNode)
//
//     p.l.Trace("constructing node map")
// 	for i := range p.Nodes {
// 		nodeMap[p.Nodes[i].Name] = &p.Nodes[i]
// 	}
//
//     p.l.Trace("creating linked tree")
// 	for i := range p.Nodes {
// 		p.Nodes[i].Next = []*PipelineNode{}
// 		p.Nodes[i].Prev = []*PipelineNode{}
// 	}
//
// 	for i := range p.Nodes {
// 		for _, output := range p.Nodes[i].Outputs {
// 			nextNode, exists := nodeMap[output]
// 			if !exists {
// 				continue
// 			}
//             p.l.Trace("linking %s ---> %s", p.Nodes[i].Name, nextNode.Name)
// 			p.Nodes[i].Next = append(p.Nodes[i].Next, nextNode)
// 			nextNode.Prev = append(nextNode.Prev, &p.Nodes[i])
// 		}
// 	}
//     p.l.Trace("successfully linked nodes")
// 	return nil
// }
//
// func (p Pipeline) Find(name string) *PipelineNode {
// 	for i := range p.Nodes {
// 		if p.Nodes[i].Name == name {
// 			return &p.Nodes[i]
// 		}
// 	}
// 	return nil
// }
//
// func (p *Pipeline) calculateDistances() (map[string]int, map[string]int) {
// 	distanceToStart := make(map[string]int)
// 	distanceToEnd := make(map[string]int)
//
// 	for _, node := range p.Nodes {
// 		start, end := node.Distance()
// 		distanceToStart[node.Name] = start
// 		distanceToEnd[node.Name] = end
// 	}
//
// 	return distanceToStart, distanceToEnd
// }
//
// // Finds and returns the nodes in the middle of the pipeline
// func (p *Pipeline) FindMidpoint() []PipelineNode {
//     p.l.Tracef("finding midpoint node(s) from %d nodes", len(p.Nodes))
// 	distanceToStart, distanceToEnd := p.calculateDistances()
// 	var middleNodes []PipelineNode
//
//     var names []string
// 	for _, node := range p.Nodes {
//         p.l.Tracef(
//             "node %s: distanceToStart=%d, distanceToEnd=%d",
//             node.Name,
//             distanceToStart[node.Name],
//             distanceToEnd[node.Name],
//         )
//         names = append(names, node.Name)
// 		if distanceToStart[node.Name] == distanceToEnd[node.Name] {
// 			middleNodes = append(middleNodes, node)
// 		}
// 	}
//     p.l.Tracef(strings.Join(names, " ---> "))
//
//     p.l.Tracef("found %d midpoint node(s)", len(middleNodes))
// 	return middleNodes
// }
//
// // traverseBefore recursively marks the ancestor nodes of the given node
// // to be pruned unless they are also dependencies for other nodes.
// func (p *Pipeline) traverseBefore(node *PipelineNode, toBePruned map[string]bool) {
// 	if node == nil || toBePruned[node.Name] {
// 		return
// 	}
// 	toBePruned[node.Name] = true
// 	for _, prevNode := range node.Prev {
// 		// If this previous node is also a dependency for other nodes, skip it
// 		if len(prevNode.Next) > 1 {
// 			continue
// 		}
// 		p.traverseBefore(prevNode, toBePruned)
// 	}
// }
//
// // traverseAfter recursively marks all the descendant nodes of the given node to be pruned.
// func (p *Pipeline) traverseAfter(node *PipelineNode, toBePruned map[string]bool) {
// 	if node == nil || toBePruned[node.Name] {
// 		return
// 	}
// 	toBePruned[node.Name] = true
// 	for _, nextNode := range node.Next {
// 		p.traverseAfter(nextNode, toBePruned)
// 	}
// }
//
// // PruneNodes prunes nodes in the pipeline based on the given flag and target node.
// // If the boolean flag is true, all nodes before the provided node are pruned,
// // unless the ancestor node is a dependency for another branch.
// // If the boolean flag is false, all of the descendant nodes are pruned from the DAG.
// func (p *Pipeline) PruneNodes(targetNode *PipelineNode) {
// 	toBePruned := map[string]bool{}
//
// 	if *targetNode.IsValid {
// 		p.traverseBefore(targetNode, toBePruned)
// 	} else {
// 		p.traverseAfter(targetNode, toBePruned)
// 	}
//
// 	newNodes := []PipelineNode{}
// 	for _, node := range p.Nodes {
// 		if !toBePruned[node.Name] {
// 			newNodes = append(newNodes, node)
// 		}
// 	}
// 	p.Nodes = newNodes
// }
//
//
// func (p *Pipeline) AllNodesInspected() bool {
//     for _, node := range p.Nodes {
//         if node.IsValid == nil {
//             return false
//         }
//     }
//     return true
// }
//
