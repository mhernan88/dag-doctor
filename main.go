package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
    "github.com/mhernan88/dag-bisect/dag"
)

func main() {
	app := &cli.App{
		Name:  "Pipeline Inspector",
		Usage: "Inspect a data pipeline for errors",
		Action: func(c *cli.Context) error {
            pipeline, err := dag.Load("dag2.json")
            if err != nil {
                return err
            }

            if pipeline == nil {
                return fmt.Errorf("pipeline was nil")
            }
            fmt.Printf("loaded pipeline with %d nodes\n", len(pipeline.Nodes))

            err = pipeline.Link()
            if err != nil {
                return err
            }
            fmt.Println("linked pipeline nodes")


			// Iterate through the inspection process
			errorNodes := inspectPipeline(pipeline)

			// Handle the results
			if len(errorNodes) == 0 {
				fmt.Println("No errors found in the pipeline.")
			} else {
				fmt.Printf("Errors found in the following nodes: %v\n", errorNodes)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func inspectPipeline(pipeline *dag.Pipeline) []string {
	reader := bufio.NewReader(os.Stdin)
	var errorNodes []string

	for !allNodesInspected(pipeline) {
		midpoints := pipeline.FindMidpoint()
		node := midpoints[0] // You might want to handle multiple midpoints here

		// Show filename to inspect
		fmt.Printf("Please inspect the file: %s\n", node.Name)

		// Ask the user if the file looks correct
		for {
			fmt.Println("Does the file look correct? (yes/no):")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			isValid := false
			if input == "yes" {
				isValid = true
				pipeline.PruneNodes(false, &node)
			} else if input == "no" {
				pipeline.PruneNodes(true, &node)
			} else {
				fmt.Println("Invalid input, please enter 'yes' or 'no'.")
				continue
			}

			node.IsValid = &isValid

			// Confirm error if all input nodes are correct
			if !isValid && allInputsCorrect(&node) {
				errorNodes = append(errorNodes, node.Name)
			}

			break
		}
	}

	return errorNodes
}

func allNodesInspected(pipeline *dag.Pipeline) bool {
	for _, node := range pipeline.Nodes {
		if node.IsValid == nil {
			return false
		}
	}
	return true
}

func allInputsCorrect(node *dag.PipelineNode) bool {
	for _, prevNode := range node.Prev {
		if prevNode.IsValid == nil || *prevNode.IsValid == false {
			return false
		}
	}
	return true
}

