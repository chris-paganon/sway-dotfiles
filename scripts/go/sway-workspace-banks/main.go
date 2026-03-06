package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/joshuarubin/go-sway"
)

const (
	activeStart = 1
	activeEnd   = 9
	bankStride  = 10
	bankWidth   = 9
	statePath   = "/tmp/sway-workspace-banks-state.json"
)

type workspaceInfo struct {
	Exists   bool
	NonEmpty bool
}

type focusedWorkspace struct {
	Num    int
	Output string
}

type state struct {
	CurrentBankStart int   `json:"current_bank_start"`
	UpdatedAt        int64 `json:"updated_at"`
}

func main() {
	cmd, err := parseCommand(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	client, err := sway.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	focused, err := getFocusedWorkspace(client, ctx)
	if err != nil {
		log.Fatal(err)
	}

	if !isActiveWorkspaceNumber(focused.Num) {
		notify("Workspace banks", "Run only from workspaces 1-9")
		return
	}

	infoByNum, err := getWorkspaceInfo(client, ctx)
	if err != nil {
		log.Fatal(err)
	}

	st, err := loadState(statePath)
	if err != nil {
		log.Fatal(err)
	}

	switch cmd {
	case "swap-next":
		err = doSwap(ctx, client, focused, infoByNum, st, true)
	case "swap-prev":
		err = doSwap(ctx, client, focused, infoByNum, st, false)
	case "push-next":
		err = doPush(ctx, client, focused, infoByNum, st)
	default:
		err = fmt.Errorf("unsupported command: %s", cmd)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func parseCommand(args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("usage: sway-workspace-banks <swap-next|swap-prev|push-next>")
	}

	switch args[1] {
	case "swap-next", "swap-prev", "push-next":
		return args[1], nil
	default:
		return "", fmt.Errorf("invalid command: %s", args[1])
	}
}

func doSwap(
	ctx context.Context,
	client sway.Client,
	focused focusedWorkspace,
	infoByNum map[int]workspaceInfo,
	st state,
	next bool,
) error {
	banks := nonEmptyBanks(infoByNum)
	if len(banks) == 0 {
		notify("Workspace banks", "No non-empty banks found")
		return nil
	}

	target := selectBank(banks, st.CurrentBankStart, next)
	if target == 0 {
		notify("Workspace banks", "No target bank found")
		return nil
	}

	if err := swapWithBank(ctx, client, infoByNum, target); err != nil {
		return err
	}

	if err := restoreFocus(ctx, client, focused); err != nil {
		return err
	}

	st.CurrentBankStart = target
	st.UpdatedAt = time.Now().Unix()
	if err := saveState(statePath, st); err != nil {
		return err
	}

	notify("Workspace banks", fmt.Sprintf("Swapped with bank %d-%d", target, target+bankWidth-1))
	return nil
}

func doPush(
	ctx context.Context,
	client sway.Client,
	focused focusedWorkspace,
	infoByNum map[int]workspaceInfo,
	st state,
) error {
	if activeSetIsEmpty(infoByNum) {
		notify("Workspace banks", "Active set 1-9 is empty")
		return nil
	}

	target := findNextAvailableBank(infoByNum)
	if target == 0 {
		return errors.New("could not find an available bank")
	}

	for slot := activeStart; slot <= activeEnd; slot++ {
		if !infoByNum[slot].Exists {
			continue
		}

		targetSlot := bankSlot(target, slot)
		if err := runCommand(client, ctx, fmt.Sprintf("rename workspace number %d to %d", slot, targetSlot)); err != nil {
			return err
		}
	}

	if err := restoreFocus(ctx, client, focused); err != nil {
		return err
	}

	st.CurrentBankStart = target
	st.UpdatedAt = time.Now().Unix()
	if err := saveState(statePath, st); err != nil {
		return err
	}

	notify("Workspace banks", fmt.Sprintf("Pushed active set to %d-%d", target, target+bankWidth-1))
	return nil
}

func getFocusedWorkspace(client sway.Client, ctx context.Context) (focusedWorkspace, error) {
	workspaces, err := client.GetWorkspaces(ctx)
	if err != nil {
		return focusedWorkspace{}, err
	}

	focusedIndex := slices.IndexFunc(workspaces, func(workspace sway.Workspace) bool {
		return workspace.Focused
	})
	if focusedIndex == -1 {
		return focusedWorkspace{}, errors.New("could not locate focused workspace")
	}

	focused := workspaces[focusedIndex]
	return focusedWorkspace{
		Num:    int(focused.Num),
		Output: focused.Output,
	}, nil
}

func isActiveWorkspaceNumber(num int) bool {
	return num >= activeStart && num <= activeEnd
}

func getWorkspaceInfo(client sway.Client, ctx context.Context) (map[int]workspaceInfo, error) {
	tree, err := client.GetTree(ctx)
	if err != nil {
		return nil, err
	}

	infoByNum := make(map[int]workspaceInfo)
	queue := []*sway.Node{tree}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node == nil {
			continue
		}

		if node.Type == "workspace" {
			num, ok := workspaceNumberFromName(node.Name)
			if ok && num > 0 {
				infoByNum[num] = workspaceInfo{
					Exists:   true,
					NonEmpty: len(node.Nodes) > 0 || len(node.FloatingNodes) > 0,
				}
			}
		}

		queue = append(queue, node.Nodes...)
		queue = append(queue, node.FloatingNodes...)
	}

	return infoByNum, nil
}

func workspaceNumberFromName(name string) (int, bool) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return 0, false
	}

	firstToken := trimmed
	if strings.Contains(trimmed, ":") {
		firstToken = strings.SplitN(trimmed, ":", 2)[0]
	} else {
		parts := strings.Fields(trimmed)
		if len(parts) == 0 {
			return 0, false
		}
		firstToken = parts[0]
	}

	num, err := strconv.Atoi(firstToken)
	if err != nil || num < 0 {
		return 0, false
	}

	return num, true
}

func nonEmptyBanks(infoByNum map[int]workspaceInfo) []int {
	bankSet := make(map[int]bool)
	for num, info := range infoByNum {
		if !info.NonEmpty {
			continue
		}

		start, ok := bankStartForWorkspace(num)
		if !ok {
			continue
		}

		bankSet[start] = true
	}

	banks := make([]int, 0, len(bankSet))
	for start := range bankSet {
		banks = append(banks, start)
	}

	sort.Ints(banks)
	return banks
}

func bankStartForWorkspace(num int) (int, bool) {
	if num < 10 {
		return 0, false
	}
	if num%10 == 9 {
		return 0, false
	}

	start := (num / 10) * 10
	if num < start || num > start+bankWidth-1 {
		return 0, false
	}

	return start, true
}

func selectBank(banks []int, cursor int, next bool) int {
	if len(banks) == 0 {
		return 0
	}

	if !slices.Contains(banks, cursor) {
		if next {
			return banks[0]
		}
		return banks[len(banks)-1]
	}

	if next {
		for _, bank := range banks {
			if bank > cursor {
				return bank
			}
		}
		return banks[0]
	}

	for i := len(banks) - 1; i >= 0; i-- {
		if banks[i] < cursor {
			return banks[i]
		}
	}

	return banks[len(banks)-1]
}

func swapWithBank(ctx context.Context, client sway.Client, infoByNum map[int]workspaceInfo, bankStart int) error {
	tmpBySlot := make(map[int]string)
	token := strconv.FormatInt(time.Now().UnixNano(), 36)

	for slot := activeStart; slot <= activeEnd; slot++ {
		if !infoByNum[slot].Exists {
			continue
		}

		tmpName := fmt.Sprintf("__banktmp_%s_%d", token, slot)
		tmpBySlot[slot] = tmpName
		if err := runCommand(client, ctx, fmt.Sprintf("rename workspace number %d to %s", slot, tmpName)); err != nil {
			return err
		}
	}

	for slot := activeStart; slot <= activeEnd; slot++ {
		targetSlot := bankSlot(bankStart, slot)
		if !infoByNum[targetSlot].Exists {
			continue
		}

		if err := runCommand(client, ctx, fmt.Sprintf("rename workspace number %d to %d", targetSlot, slot)); err != nil {
			return err
		}
	}

	for slot := activeStart; slot <= activeEnd; slot++ {
		tmpName, ok := tmpBySlot[slot]
		if !ok {
			continue
		}

		targetSlot := bankSlot(bankStart, slot)
		if err := runCommand(client, ctx, fmt.Sprintf("rename workspace %s to %d", tmpName, targetSlot)); err != nil {
			return err
		}
	}

	return nil
}

func bankSlot(start int, activeSlot int) int {
	return start + (activeSlot - activeStart)
}

func findNextAvailableBank(infoByNum map[int]workspaceInfo) int {
	for start := 10; start < 10000; start += bankStride {
		available := true
		for offset := 0; offset < bankWidth; offset++ {
			num := start + offset
			info := infoByNum[num]
			if info.Exists && info.NonEmpty {
				available = false
				break
			}
		}

		if available {
			return start
		}
	}

	return 0
}

func activeSetIsEmpty(infoByNum map[int]workspaceInfo) bool {
	for slot := activeStart; slot <= activeEnd; slot++ {
		if infoByNum[slot].NonEmpty {
			return false
		}
	}

	return true
}

func restoreFocus(ctx context.Context, client sway.Client, focused focusedWorkspace) error {
	if err := runCommand(client, ctx, fmt.Sprintf("focus output %s", focused.Output)); err != nil {
		return err
	}

	if err := runCommand(client, ctx, fmt.Sprintf("workspace number %d", focused.Num)); err != nil {
		return err
	}

	return nil
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

func loadState(path string) (state, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return state{}, nil
		}
		return state{}, err
	}

	var st state
	if err := json.Unmarshal(data, &st); err != nil {
		return state{}, nil
	}

	return st, nil
}

func saveState(path string, st state) error {
	data, err := json.Marshal(st)
	if err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp("/tmp", "sway-workspace-banks-state-*.json")
	if err != nil {
		return err
	}

	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}

	return nil
}

func notify(summary string, body string) {
	cmd := exec.Command("notify-send", summary, body)
	_ = cmd.Run()
}
