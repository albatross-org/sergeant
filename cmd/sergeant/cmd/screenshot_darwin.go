package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

func screenshot(dest string) error {
	maimCmd := exec.Command("screencapture", "-s", dest)

	bytes, err := maimCmd.CombinedOutput()
	output := string(bytes)

	if err != nil && !strings.Contains(output, "Selection was cancelled by keystroke or right-click.") {
		return fmt.Errorf("maim command '%s' exited with message '%s', error: %w", maimCmd.String(), output, err)
	}

	return nil
}
