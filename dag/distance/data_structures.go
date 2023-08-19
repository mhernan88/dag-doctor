package distance

// TODO: Move to another package
type Node struct {
	Name             string          `json:"name"`
	Inputs           []string        `json:"inputs"`
	Outputs          []string        `json:"outputs"`
	Next             []*Node `json:"-"`
	Prev             []*Node `json:"-"`
    IsValid *bool `json:"-"`
}

type NodeWrapper struct {
    Node Node
    Distance int
}

type Distance struct {
    ToStart int
    ToEnd int
}
