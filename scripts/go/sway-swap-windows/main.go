package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strings"

	"go.i3wm.org/i3"
)

func main() {
	swapDirection, err := getSwapDirection()
	if err != nil {
		log.Fatal(err)
	}

	tree, err := i3.GetTree()
	if err != nil {
		log.Fatal(err)
	}
	focusedWorkspace := tree.Root.FindFocused(func(node *i3.Node) bool {
		return node.Type == i3.WorkspaceNode
	})
	if focusedWorkspace == nil {
		log.Fatal("could not locate workspace")
	}
	focusedWorkspaceName := focusedWorkspace.Name

	realOutputs, err := getSortedOutputs()
	if err != nil {
		log.Fatal(err)
	}

	workspaceToSwapWith, err := getWorkspaceToSwapWith(tree, realOutputs, swapDirection, focusedWorkspaceName)
	if err != nil {
		log.Fatal(err)
	}

	windowToSwapWith := workspaceToSwapWith.Nodes[0].Window

	i3.RunCommand(fmt.Sprintf("move container to workspace %s", workspaceToSwapWith.Name))
	i3.RunCommand(fmt.Sprintf("[id=\"%d\"] focus", windowToSwapWith))
	i3.RunCommand(fmt.Sprintf("move container to workspace %s", focusedWorkspaceName))
}

func getSwapDirection() (string, error) {
	args := os.Args
	if len(args) == 1 {
		return "", fmt.Errorf("no argument provided")
	}

	swapDirection := args[1]
	if swapDirection != "left" && swapDirection != "right" {
		return "", fmt.Errorf("invalid argument provided, must be 'left' or 'right'")
	}

	return swapDirection, nil
}

func getSortedOutputs() ([]*i3.Output, error) {
	outputs, err := i3.GetOutputs()
	if err != nil {
		return nil, err
	}

	outputSlice := make([]*i3.Output, len(outputs))
	for i := range outputs {
		outputSlice[i] = &outputs[i]
	}

	// sort outputs by rect.x
	sort.Slice(outputSlice, func(i, j int) bool {
		return outputSlice[i].Rect.X < outputSlice[j].Rect.X
	})

	realOutputs := slices.DeleteFunc(outputSlice, func(output *i3.Output) bool {
		return strings.Contains(output.Name, "xroot")
	})

	if len(realOutputs) == 0 {
		return nil, fmt.Errorf("no outputs found")
	}

	return realOutputs, nil
}

func getWorkspaceToSwapWith(tree i3.Tree, outputs []*i3.Output, swapDirection string, focusedWorkspaceName string) (*i3.Node, error) {
	focusedOutputIndex := slices.IndexFunc(outputs, func(output *i3.Output) bool {
		return output.CurrentWorkspace == focusedWorkspaceName
	})

	if focusedOutputIndex == -1 {
		return nil, fmt.Errorf("could not locate focused output")
	}

	var outputIndexToSwapWith int
	if swapDirection == "left" {
		outputIndexToSwapWith = focusedOutputIndex - 1
		if outputIndexToSwapWith < 0 {
			outputIndexToSwapWith = len(outputs) - 1
		}
	} else {
		outputIndexToSwapWith = focusedOutputIndex + 1
		if outputIndexToSwapWith >= len(outputs) {
			outputIndexToSwapWith = 0
		}
	}

	outputToSwapWith := outputs[outputIndexToSwapWith]

	workspaceToSwapWith := tree.Root.FindChild(func(node *i3.Node) bool {
		return node.Type == i3.WorkspaceNode && node.Name == outputToSwapWith.CurrentWorkspace
	})

	if workspaceToSwapWith == nil {
		return nil, fmt.Errorf("could not locate workspace to swap with")
	}

	return workspaceToSwapWith, nil
}
