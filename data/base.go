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

type Dataset struct {
    Filename string `json:"filename"`
}

