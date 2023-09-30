package models

type Node struct {
	Name    string   `json:"name"`
	Inputs  []string `json:"inputs"`
	Outputs []string `json:"outputs"`
	Next    []string `json:"next"` // []nodeName
	Prev    []string `json:"prev"` // []nodeName
}
