package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var autoUpdate bool
var excludeFile string

func init() {
	flag.BoolVar(&autoUpdate, "update", false, "Automatically update VMs with available updates")
	flag.StringVar(&excludeFile, "exclude", "exclude.txt", "Path to file containing excluded VM names")
}

func listVMs() ([]string, error) {
	cmd := exec.Command("qvm-ls", "--raw-list")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
}

func checkUpdates(vm string) (bool, string, error) {
	cmd := exec.Command("qvm-run", "-p", vm, "sudo dnf updateinfo list available")
	output, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "No update notices found") {
		return false, "", err
	}
	hasUpdates := !strings.Contains(string(output), "No update notices found")
	return hasUpdates, string(output), nil
}

func updateVM(vm string) error {
	cmd := exec.Command("qvm-run", "-p", vm, "sudo dnf update -y")
	output, err := cmd.CombinedOutput()
	fmt.Printf("🔧 Update output for %s:\n%s\n", vm, string(output))
	return err
}

func loadExclusions(filename string) (map[string]bool, error) {
	excluded := make(map[string]bool)
	data, err := os.ReadFile(filename)
	if err != nil {
		return excluded, err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			excluded[line] = true
		}
	}
	return excluded, nil

}

func main() {
	logFile, err := os.OpenFile("qubes-updater.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Unable to create log file :", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags)

	flag.Parse()

	excluded, err := loadExclusions(excludeFile)
	if err != nil {
		fmt.Printf(" Any exclude found  (%s), we continue without.\n", excludeFile)
	}

	vms, err := listVMs()
	if err != nil {
		fmt.Println(" VMs loading error:", err)
		return
	}

	for _, vm := range vms {
		if vm == "dom0" || excluded[vm] {
			fmt.Printf(" Skip VM: %s\n", vm)
			continue
		}

		fmt.Printf("Checking VM: %s\n", vm)
		updates, output, err := checkUpdates(vm)
		if err != nil {
			fmt.Printf(" Error with %s: %s\n", vm, err)
			continue
		}

		if updates {
			fmt.Printf("Updates available for %s:\n%s\n", vm, output)
			if autoUpdate {
				fmt.Printf(" Updating %s...\n", vm)
				if err := updateVM(vm); err != nil {
					fmt.Printf("Update failed: %s\n", err)
				} else {
					fmt.Printf("%s updated successfully.\n", vm)
				}
			}
		} else {
			fmt.Printf("%s is up-to-date.\n", vm)
		}
	}
}
