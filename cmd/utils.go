/*
Copyright © 2025 Christophe Jauffret <reg-github@geo6.net>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// DisplaySpinner shows a spinning animation with the given message.
// It returns a channel that should be closed to stop the spinner.
// When closed, it displays a green checkmark and the message.
func DisplaySpinner(message string) chan struct{} {
	spinChars := []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
	done := make(chan struct{})

	go func() {
		i := 0
		for {
			select {
			case <-done:
				// Replace spinner with green checkmark
				fmt.Printf("\r \033[32m✓\033[0m %s\n", message)
				return
			default:
				fmt.Printf("\r%c  %s", spinChars[i], message)
				i = (i + 1) % len(spinChars)
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()

	return done
}

func runCommand(cmd *exec.Cmd, showOutput bool) error {
	if showOutput {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}
