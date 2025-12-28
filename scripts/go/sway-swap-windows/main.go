package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/joshuarubin/go-sway"
)

func main() {
	swapDirection, err := getSwapDirection()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	client, err := sway.New(ctx)
	if err != nil {

		log.Fatal("no client", err)
	}

	tree, err := client.GetTree(ctx)
	if err != nil {
		log.Fatal(err)
	}

	workspaces, err := client.GetWorkspaces(ctx)
	if err != nil {
		log.Fatal(err)
	}

	focusedWorkspaceIndex := slices.IndexFunc(workspaces, func(workspace sway.Workspace) bool {
		return workspace.Focused
	})
	if focusedWorkspaceIndex == -1 {
		log.Fatal("could not locate workspace")
	}

	focusedWorkspace := workspaces[focusedWorkspaceIndex]

	fmt.Print("focusedWorkspace")
	fmt.Println(focusedWorkspace.Name)
	focusedWorkspaceName := focusedWorkspace.Name

	realOutputs, err := getSortedOutputs(client, ctx)
	if err != nil {
		log.Fatal(err)
	}

	workspaceToSwapWith, err := getWorkspaceToSwapWith(tree, realOutputs, swapDirection, focusedWorkspace)
	if err != nil {
		log.Fatal(err)
	}

	windowToSwapWith := workspaceToSwapWith.Nodes[0].ID

	client.RunCommand(ctx, fmt.Sprintf("move container to workspace %s", workspaceToSwapWith.Name))
	client.RunCommand(ctx, fmt.Sprintf("[con_id=\"%d\"] focus", windowToSwapWith))
	client.RunCommand(ctx, fmt.Sprintf("move container to workspace %s", focusedWorkspaceName))
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

func getSortedOutputs(client sway.Client, ctx context.Context) ([]*sway.Output, error) {
	outputs, err := client.GetOutputs(ctx)
	if err != nil {
		return nil, err
	}

	outputSlice := make([]*sway.Output, len(outputs))
	for i := range outputs {
		outputSlice[i] = &outputs[i]
	}

	// sort outputs by rect.x
	sort.Slice(outputSlice, func(i, j int) bool {
		return outputSlice[i].Rect.X < outputSlice[j].Rect.X
	})

	realOutputs := slices.DeleteFunc(outputSlice, func(output *sway.Output) bool {
		return strings.Contains(output.Name, "xroot")
	})

	if len(realOutputs) == 0 {
		return nil, fmt.Errorf("no outputs found")
	}

	return realOutputs, nil
}

func getWorkspaceToSwapWith(tree *sway.Node, outputs []*sway.Output, swapDirection string, focusedWorkspace sway.Workspace) (*sway.Node, error) {
	fmt.Printf("Current workspace output: %s", focusedWorkspace.Output)
	fmt.Println()

	for _, output := range outputs {
		fmt.Println(output.Name)
		fmt.Println(output.CurrentWorkspace)
	}

	focusedOutputIndex := slices.IndexFunc(outputs, func(output *sway.Output) bool {
		return output.Name == focusedWorkspace.Output
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

	workspaceToSwapWith := tree.TraverseNodes(func(node *sway.Node) bool {
		return node.Type == "workspace" && node.Name == outputToSwapWith.CurrentWorkspace
	})

	if workspaceToSwapWith == nil {
		return nil, fmt.Errorf("could not locate workspace to swap with")
	}

	return workspaceToSwapWith, nil
}
