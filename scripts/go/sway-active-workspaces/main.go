package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"sort"

	"github.com/joshuarubin/go-sway"
)

const (
	activeStart = 1
	activeEnd   = 9
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: sway-active-workspaces <command>")
	}

	ctx := context.Background()
	client, err := sway.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "open-next-free":
		err = openNextFree(ctx, client)
	case "start-on-next-free":
		err = startOnNextFree(ctx, client, args)
	case "move-to-next-free":
		err = moveToNextFree(ctx, client)
	case "focus-next-on-output":
		err = focusNextOnOutput(ctx, client)
	case "move-follow-next-on-output":
		err = moveFollowNextOnOutput(ctx, client)
	case "move-to-output-cycle":
		err = moveToOutputCycle(ctx, client, args)
	case "move-to-output-new":
		err = moveToOutputNew(ctx, client, args)
	default:
		err = fmt.Errorf("invalid command: %s", command)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func openNextFree(ctx context.Context, client sway.Client) error {
	focused, workspaces, err := getFocusedAndWorkspaces(ctx, client)
	if err != nil {
		return err
	}

	if !isActiveWorkspace(focused.Num) {
		notifyError("Run only from workspaces 1-9")
		return nil
	}

	nextFree, ok := firstFreeSlot(workspaces)
	if !ok {
		notifyError("No free workspaces in 1-9")
		return nil
	}

	return runCommand(client, ctx, fmt.Sprintf("workspace number %d", nextFree))
}

func startOnNextFree(ctx context.Context, client sway.Client, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: start-on-next-free <command>")
	}

	focused, workspaces, err := getFocusedAndWorkspaces(ctx, client)
	if err != nil {
		return err
	}

	if !isActiveWorkspace(focused.Num) {
		notifyError("Run only from workspaces 1-9")
		return nil
	}

	nextFree, ok := firstFreeSlot(workspaces)
	if !ok {
		notifyError("No free workspaces in 1-9")
		return nil
	}

	if err := runCommand(client, ctx, fmt.Sprintf("workspace number %d", nextFree)); err != nil {
		return err
	}

	return runCommand(client, ctx, "exec "+joinArgs(args))
}

func moveToNextFree(ctx context.Context, client sway.Client) error {
	focused, workspaces, err := getFocusedAndWorkspaces(ctx, client)
	if err != nil {
		return err
	}

	if !isActiveWorkspace(focused.Num) {
		notifyError("Run only from workspaces 1-9")
		return nil
	}

	nextFree, ok := firstFreeSlot(workspaces)
	if !ok {
		notifyError("No free workspaces in 1-9")
		return nil
	}

	if err := runCommand(client, ctx, fmt.Sprintf("move container to workspace number %d", nextFree)); err != nil {
		return err
	}

	return runCommand(client, ctx, fmt.Sprintf("workspace number %d", nextFree))
}

func focusNextOnOutput(ctx context.Context, client sway.Client) error {
	focused, workspaces, err := getFocusedAndWorkspaces(ctx, client)
	if err != nil {
		return err
	}

	if !isActiveWorkspace(focused.Num) {
		notifyError("Run only from workspaces 1-9")
		return nil
	}

	next, ok := nextWorkspaceOnOutput(workspaces, focused.Output, focused.Num)
	if !ok {
		return nil
	}

	return runCommand(client, ctx, fmt.Sprintf("workspace number %d", next))
}

func moveFollowNextOnOutput(ctx context.Context, client sway.Client) error {
	focused, workspaces, err := getFocusedAndWorkspaces(ctx, client)
	if err != nil {
		return err
	}

	if !isActiveWorkspace(focused.Num) {
		notifyError("Run only from workspaces 1-9")
		return nil
	}

	next, ok := nextWorkspaceOnOutput(workspaces, focused.Output, focused.Num)
	if !ok {
		return nil
	}

	if err := runCommand(client, ctx, fmt.Sprintf("move container to workspace number %d", next)); err != nil {
		return err
	}

	return runCommand(client, ctx, fmt.Sprintf("workspace number %d", next))
}

func moveToOutputCycle(ctx context.Context, client sway.Client, args []string) error {
	direction, err := parseDirectionArg(args)
	if err != nil {
		return err
	}

	if err := runCommand(client, ctx, fmt.Sprintf("move container to output %s", direction)); err != nil {
		return err
	}
	if err := runCommand(client, ctx, fmt.Sprintf("focus output %s", direction)); err != nil {
		return err
	}

	focused, workspaces, err := getFocusedAndWorkspaces(ctx, client)
	if err != nil {
		return err
	}

	next, ok := nextWorkspaceOnOutput(workspaces, focused.Output, focused.Num)
	if !ok {
		return nil
	}

	return runCommand(client, ctx, fmt.Sprintf("workspace number %d", next))
}

func moveToOutputNew(ctx context.Context, client sway.Client, args []string) error {
	direction, err := parseDirectionArg(args)
	if err != nil {
		return err
	}

	if err := runCommand(client, ctx, fmt.Sprintf("move container to output %s", direction)); err != nil {
		return err
	}
	if err := runCommand(client, ctx, fmt.Sprintf("focus output %s", direction)); err != nil {
		return err
	}

	_, workspaces, err := getFocusedAndWorkspaces(ctx, client)
	if err != nil {
		return err
	}

	nextFree, ok := firstFreeSlot(workspaces)
	if !ok {
		notifyError("No free workspaces in 1-9")
		return nil
	}

	if err := runCommand(client, ctx, fmt.Sprintf("move container to workspace number %d", nextFree)); err != nil {
		return err
	}

	return runCommand(client, ctx, fmt.Sprintf("workspace number %d", nextFree))
}

func getFocusedAndWorkspaces(ctx context.Context, client sway.Client) (sway.Workspace, []sway.Workspace, error) {
	workspaces, err := client.GetWorkspaces(ctx)
	if err != nil {
		return sway.Workspace{}, nil, err
	}

	focusedIndex := slices.IndexFunc(workspaces, func(workspace sway.Workspace) bool {
		return workspace.Focused
	})
	if focusedIndex == -1 {
		return sway.Workspace{}, nil, fmt.Errorf("could not locate focused workspace")
	}

	return workspaces[focusedIndex], workspaces, nil
}

func isActiveWorkspace(num int64) bool {
	return num >= activeStart && num <= activeEnd
}

func firstFreeSlot(workspaces []sway.Workspace) (int, bool) {
	used := make(map[int64]bool)
	for _, ws := range workspaces {
		if isActiveWorkspace(ws.Num) {
			used[ws.Num] = true
		}
	}

	for slot := activeStart; slot <= activeEnd; slot++ {
		if !used[int64(slot)] {
			return slot, true
		}
	}

	return 0, false
}

func nextWorkspaceOnOutput(workspaces []sway.Workspace, output string, current int64) (int, bool) {
	if output == "" {
		return 0, false
	}

	var onOutput []int
	for _, ws := range workspaces {
		if ws.Output != output {
			continue
		}
		if !isActiveWorkspace(ws.Num) {
			continue
		}
		onOutput = append(onOutput, int(ws.Num))
	}

	if len(onOutput) == 0 {
		return 0, false
	}

	sort.Ints(onOutput)
	onOutput = slices.Compact(onOutput)
	if len(onOutput) <= 1 {
		return 0, false
	}

	if !slices.Contains(onOutput, int(current)) {
		return onOutput[0], true
	}

	for _, slot := range onOutput {
		if slot > int(current) {
			return slot, true
		}
	}

	return onOutput[0], true
}

func parseDirectionArg(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("expected direction argument: left|right")
	}

	if args[0] != "left" && args[0] != "right" {
		return "", fmt.Errorf("invalid direction: %s", args[0])
	}

	return args[0], nil
}

func joinArgs(args []string) string {
	if len(args) == 0 {
		return ""
	}

	cmd := args[0]
	for i := 1; i < len(args); i++ {
		cmd += " " + args[i]
	}
	return cmd
}

func runCommand(client sway.Client, ctx context.Context, command string) error {
	replies, err := client.RunCommand(ctx, command)
	if err != nil {
		return err
	}

	for _, reply := range replies {
		if !reply.Success {
			if reply.Error != "" {
				return fmt.Errorf("sway command failed (%s): %s", command, reply.Error)
			}
			return fmt.Errorf("sway command failed (%s)", command)
		}
	}

	return nil
}

func notifyError(message string) {
	cmd := exec.Command("notify-send", "Workspace navigation", message)
	_ = cmd.Run()
}
