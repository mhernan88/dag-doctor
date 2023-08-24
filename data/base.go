package data

type Node struct {
	Name             string          `json:"name"`
	Inputs           []string        `json:"inputs"`
	Outputs          []string        `json:"outputs"`
    NextLabels []string `json:"next"` // []nodeName
    PrevLabels []string `json:"prev"` // []nodeName
	Next             map[string]*Node `json:"-"` // map[nodeName]nodePointer
	Prev             map[string]*Node `json:"-"` // map[nodeName]nodePointer
}

type DAG struct {
    Nodes map[string]*Node
}

func NewDAG(nodes []*Node) DAG {
    dag := make(map[string]*Node)
    for _, node := range nodes {
        dag[node.Name] = node
    }
    return DAG{Nodes: dag}
}

func (d DAG) ToSlices() ([]string, []*Node){
    var names []string
    var nodes []*Node
    for name, node := range d.Nodes {
        names = append(names, name)
        nodes = append(nodes, node)
    }
    return names, nodes
}

type Dataset struct {
    Filename string `json:"filename"`
}

