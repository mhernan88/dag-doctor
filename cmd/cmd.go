package cmd

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/mhernan88/dag-bisect/data"
	"github.com/mhernan88/dag-bisect/pruners"
	"github.com/mhernan88/dag-bisect/splitters"
)

func NewUI(
	dag data.DAG,
	splitter splitters.Splitter,
	pruner pruners.Pruner,
) UI {
	return UI{
		DAG:            dag,
		OKNodes:        make(map[string]data.Node),
		ERRNodes:       make(map[string]data.Node),
		LastFailedNode: "",
		Splitter:       splitter,
		Pruner:         pruner,
	}
}

func SaveState(filename string, ui UI) error {
	f, err := os.OpenFile(filename, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.Encode(ui)
	return nil
}

func LoadState(filename string) (*UI, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ui UI
	decoder := json.NewDecoder(f)
	decoder.Decode(&ui)
	return &ui, nil
}

type UI struct {
	DAG data.DAG `json:"dag"`
	OKNodes        map[string]data.Node `json:"ok_nodes"`
	ERRNodes       map[string]data.Node `json:"err_nodes"`
	LastFailedNode string `json:"last_failed_node"`
	Splitter       splitters.Splitter `json:"splitter"`
	Pruner         pruners.Pruner `json:"pruner"`
}

func (ui *UI) Run(l *slog.Logger) error {
	l.Debug("running ui loop")
	err := ui.CheckDAG(l)
	if err != nil {
		return err
	}
	l.Debug("terminating ui")
	return nil
}
