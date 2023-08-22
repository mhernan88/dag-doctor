package data

type Node struct {
	Name             string          `json:"name"`
	Inputs           []string        `json:"inputs"`
	Outputs          []string        `json:"outputs"`
	Next             map[string]*Node `json:"next"` // map[nodeName]nodePointer
	Prev             map[string]*Node `json:"prev"` // map[nodeName]nodePointer
}

type Dataset struct {
    Filename string `json:"filename"`
}

