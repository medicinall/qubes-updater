package qubes-updater
	import(
		"os"
		"strings"
	)
func listVMs() ([]string, error) {
	cmd := exec.COmmand("qvm-ls", "--raw-list")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	return lines, nil
}
