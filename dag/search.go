package dag

import (
    "bufio"
    "os"
    "fmt"
    "slices"
    "github.com/sirupsen/logrus"
)

func NewInspector(
    catalog map[string]PipelineFile,
    pipeline Pipeline,
    l *logrus.Logger,
) Inspector {
    r := bufio.NewReader(os.Stdin)
    return Inspector{
        catalog: catalog,
        pipeline: pipeline,
        l: l,
        r: r,
    }
}

type Inspector struct {
    catalog map[string]PipelineFile
    pipeline Pipeline
    l *logrus.Logger
    r *bufio.Reader
}

func (i *Inspector) IsDatasetOK(dataset string) (bool, error) {
    fmt.Printf("\t- inspecting dataset: %s\n", dataset)
    pipelineFile, ok := i.catalog[dataset]
    if !ok {
        return false, fmt.Errorf("dataset '%s' not found in catalog", dataset)
    }
    err := pipelineFile.LoadAndDisplay(10)
    if err != nil {
        return false, err
    }

    for {
        fmt.Println("\t- output correct? (y/n):")
        input, _, err := i.r.ReadRune()
        if err != nil {
            return false, err
        }

        if slices.Contains([]rune{'y', 'Y'}, input) {
            fmt.Printf("\t- dataset %s maked OK\n", dataset)
            return true, nil
        } else if slices.Contains([]rune{'n', 'N'}, input) {
            fmt.Printf("\t- dataset %s maked ERR\n", dataset)
            return false, nil
        }
    }
}

func (i *Inspector) IsNodeOK(node PipelineNode) (bool, error) {
    fmt.Printf("inspecting node: %s\n", node.Name)
    for _, output := range node.Outputs {
        ok, err := i.IsDatasetOK(output)
        if err != nil {
            return false, err
        }
        if !ok {
            return false, nil
        }
    }
    return true, nil
}

func (i *Inspector) IsPipelineOK() (bool, error) {
    var errNodes []string
    for !i.pipeline.AllNodesInspected() {
        midpoints := i.pipeline.FindMidpoint()
        node := midpoints[0]
        ok, err := i.IsNodeOK(node)
        if err != nil {
            return false, err
        }

        if !ok {
            errNodes = append(errNodes, node.Name)
        }
    }

    if len(errNodes) == 0 {
        fmt.Println("all clear! no errors found in pipeline")
    } else {
        fmt.Printf(
            "pipeline data issues were cuased by the following nodes: %v", 
            errNodes,
        )
    }
    return true, nil
}
