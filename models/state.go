package models


func NewDefaultState(
	dag DAG,
) *State {
	return &State{
		DAG:            dag,
		OKNodes:        make(map[string]Node),
		ERRNodes:       make(map[string]Node),
		LastFailedNode: "",
	}
}

func NewState(
	dag DAG,
) *State {
	return &State{
		DAG:            dag,
		OKNodes:        make(map[string]Node),
		ERRNodes:       make(map[string]Node),
		LastFailedNode: "",
	}
}

type State struct {
	DAG DAG `json:"dag"`
	OKNodes        map[string]Node `json:"ok_nodes"`
	ERRNodes       map[string]Node `json:"err_nodes"`
	LastFailedNode string `json:"last_failed_node"`
}
