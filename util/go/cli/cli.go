// Package cli contains command-line interface utilities.
package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/google/logger"
)

// RequestConfirmation requests a boolean response from a user in the command-line and returns it.
func RequestConfirmation(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/n]: ", prompt)
		response, err := reader.ReadString('\n')
		if err != nil {
			logger.Fatal("failed reading response: %v", err)
		}
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
