package dag

type Pipeline struct {
    Nodes []Node `json:"nodes"`
}

type Node struct {
	Name             string          `json:"name"`
	Inputs           []string        `json:"inputs"`
	Outputs          []string        `json:"outputs"`
	Next             map[string]*Node `json:"-"` // map[nodeName]nodePointer
	Prev             map[string]*Node `json:"-"` // map[nodeName]nodePointer
    IsValid *bool `json:"-"` // TODO: Can we remove this?
}

type NodeWrapper struct {
    Node Node
    Distance int
}

type Dataset struct {
    Filename string `json:"filename"`
}

type Catalog struct {
    Datasets map[string]Dataset `json:"catalog"`
}
