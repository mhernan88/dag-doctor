package data

type Dataset struct {
    Filename string `json:"filename"`
}

type Node struct {
	Name             string          `json:"name"`
	Inputs           []string        `json:"inputs"`
	Outputs          []string        `json:"outputs"`
    NextLabels []string `json:"next"` // []nodeName
    PrevLabels []string `json:"prev"` // []nodeName
	Next             map[string]*Node `json:"-"` // map[nodeName]nodePointer
	Prev             map[string]*Node `json:"-"` // map[nodeName]nodePointer
    State string `json:"-"`
}

func NewNodeMap(nodeslice []*Node) map[string]*Node {
    dag := make(map[string]*Node)
    for _, node := range nodeslice {
        dag[node.Name] = node
    }
    return dag
}

func SliceNodeMap(nodemap map[string]*Node) ([]string, []*Node) {
    var names []string
    var nodes []*Node
    for name, node := range nodemap {
        names = append(names, name)
        nodes = append(nodes, node)
    }
    return names, nodes
}

func GetNodeDescendants(sources []*Node) map[string]*Node {
    var nodes []*Node
    for _, source := range sources {
        nodes = append(nodes, source)
    }

    descendants := make(map[string]*Node)

    for len(nodes) > 0 {
        node := nodes[len(nodes)-1]
        nodes = nodes[:len(nodes)-1]

        for childName, child := range node.Next {
            descendants[childName] = child
            nodes = append(nodes, child)
        }
    }
    return descendants
}

func GetNodeAncestors(sources []*Node) map[string]*Node {
    var nodes []*Node
    for _, source := range sources {
        nodes = append(nodes, source)
    }

    ancestors := make(map[string]*Node)

    for len(nodes) > 0 {
        node := nodes[len(nodes)-1]
        nodes = nodes[:len(nodes)-1]

        for childName, child := range node.Prev {
            ancestors[childName] = child
            nodes = append(nodes, child)
        }
    }
    return ancestors
}
