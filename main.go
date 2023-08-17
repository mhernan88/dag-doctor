package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
    "github.com/sirupsen/logrus"
    "github.com/mhernan88/dag-bisect/dag"
)

var yesVariants = []rune{'y', 'Y'}
var noVariants = []rune{'n', 'N'}

var flags = []cli.Flag{
    &cli.BoolFlag{
        Name: "verbose",
        Aliases: []string{"v"},
        Value: false,
    },
    &cli.StringFlag{
        Name: "pipeline",
        Aliases: []string{"p"},
        Value: "dag.json",
    },
    &cli.StringFlag{
        Name: "catalog",
        Aliases: []string{"c"},
        Value: "catalog.json",
    },
}


func main() {
	app := &cli.App{
		Name:  "Pipeline Inspector",
		Usage: "Inspect a data pipeline for errors",
        Flags: flags,
		Action: func(c *cli.Context) error {
            log := logrus.New()
            if c.Bool("verbose") {
                log.SetLevel(logrus.TraceLevel)
            }

            log.Tracef("loading pipeline file %s", c.String("pipeline"))
            pipeline, err := dag.LoadPipeline(c.String("pipeline"))
            if err != nil {
                return err
            }

            if pipeline == nil {
                return fmt.Errorf("pipeline was nil")
            }
            log.Tracef("successfully loaded pipeline file %s", c.String("pipeline"))

            log.Tracef("linking pipeline nodes")
            err = pipeline.Link()
            if err != nil {
                return err
            }
            log.Tracef("successfully linked pipeline nodes")

            log.Tracef("loading catalog file %s", c.String("catalog"))
            catalog, err := dag.LoadCatalog(c.String("catalog"))
            if err != nil {
                return err
            }
            log.Tracef("successfully loaded catalog file %s", c.String("catalog"))
            
            inspector := dag.NewInspector(*catalog, *pipeline, log)

			// Iterate through the inspection process
            log.Tracef("performing binary error search on pipeline")
            _, err = inspector.IsPipelineOK()
            if err != nil {
                log.Errorf(err.Error())
            }
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

// func inspectPipeline(pipeline *dag.Pipeline, catalog map[string]dag.PipelineFile) ([]string, error) {
// 	reader := bufio.NewReader(os.Stdin)
// 	var errorNodes []string
//
// 	for !allNodesInspected(pipeline) {
// 		midpoints := pipeline.FindMidpoint()
// 		node := midpoints[0] // You might want to handle multiple midpoints here
//
// 		// Show filename to inspect
// 		fmt.Printf("Please inspect the file: %s\n", node.Name)
//         data, ok := catalog[node.Name]
//         if !ok {
//             return nil, nil
//         }
//
//         err := data.LoadAndDisplay(0)
//         if err != nil {
//             return nil, nil
//         }
//
// 		// Ask the user if the file looks correct
// 		for {
// 			fmt.Println("Does the file look correct? (y/n):")
// 			input, _, err := reader.ReadRune()
//             if err != nil {
//                 return nil, nil
//             }
//
// 			isValid := false
//             if slices.Contains(yesVariants, input) {
// 				isValid = true
// 				pipeline.PruneNodes(false, &node)
//             } else if slices.Contains(noVariants, input) {
// 				pipeline.PruneNodes(true, &node)
// 			} else {
//                 fmt.Println("Invalid input, please enter one of: 'y', 'n'.")
// 				continue
// 			}
//
// 			node.IsValid = &isValid
//
// 			// Confirm error if all input nodes are correct
// 			if !isValid && allInputsCorrect(&node) {
// 				errorNodes = append(errorNodes, node.Name)
// 			}
//
// 			break
// 		}
// 	}
//
// 	return errorNodes, nil
// }
//
// func allNodesInspected(pipeline *dag.Pipeline) bool {
// 	for _, node := range pipeline.Nodes {
// 		if node.IsValid == nil {
// 			return false
// 		}
// 	}
// 	return true
// }
//
// func allInputsCorrect(node *dag.PipelineNode) bool {
// 	for _, prevNode := range node.Prev {
// 		if prevNode.IsValid == nil || *prevNode.IsValid == false {
// 			return false
// 		}
// 	}
// 	return true
// }
//
