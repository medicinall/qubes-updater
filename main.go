package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// listVMs returns all Qubes VMs names using `qvm-ls --raw-list`
func listVMs() ([]string, error) {
	cmd := exec.Command("qvm-ls", "--raw-list")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	return lines, nil
}

// checkUpdates checks if a VM has available updates using `dnf`
func checkUpdates(vm string) (bool, string, error) {
	cmd := exec.Command("qvm-run", "-p", vm, "sudo dnf updateinfo list available")
	output, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "No update notices found") {
		return false, "", err
	}
	hasUpdates := !strings.Contains(string(output), "No update notices found")
	return hasUpdates, string(output), nil
}

func main() {
	vms, err := listVMs()
	if err != nil {
		fmt.Println("VMs loading error:", err)
		return
	}

	for _, vm := range vms {
		if vm == "dom0" {
			continue // dom0 updates are handled differently
		}

		fmt.Printf("Checking VM: %s\n", vm)
		updates, output, err := checkUpdates(vm)
		if err != nil {
			fmt.Printf(" Error with %s: %s\n", vm, err)
			continue
		}
		if updates {
			fmt.Printf(" Update available for %s:\n%s\n", vm, output)
		} else {
			fmt.Printf(" %s is up-to-date.\n", vm)
		}
	}
}
