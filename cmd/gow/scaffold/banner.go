package scaffold

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

const gowBanner = `
  ██████╗  ██████╗  ██╗    ██╗
 ██╔════╝ ██╔═══██╗ ██║    ██║
 ██║  ███╗██║   ██║ ██║ █╗ ██║
 ██║   ██║██║   ██║ ██║███╗██║
 ╚██████╔╝╚██████╔╝ ╚███╔███╔╝
  ╚═════╝  ╚═════╝   ╚══╝╚══╝ `

var (
	spinnerStop = make(chan struct{})
	spinnerDone = make(chan struct{})
)

// ShowBanner displays the GoW ASCII art banner.
func ShowBanner() {
	fmt.Println()
	lines := strings.Split(strings.TrimSpace(gowBanner), "\n")
	for _, line := range lines {
		fmt.Printf("  \033[36m%s\033[0m\n", line)
	}
	fmt.Println()
}

// startSpinner starts an animated spinner with the given message.
func startSpinner(message string) {
	spinnerStop = make(chan struct{})
	spinnerDone = make(chan struct{})
	spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	go func() {
		defer close(spinnerDone)
		for {
			select {
			case <-spinnerStop:
				fmt.Printf("\r\033[2K")
				return
			default:
				fmt.Printf("\r  \033[36m%s\033[0m %s...", spinners[i%len(spinners)], message)
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
}

// stopSpinner stops the current spinner.
func stopSpinner() {
	select {
	case <-spinnerStop:
		// already stopped
	default:
		close(spinnerStop)
	}
	<-spinnerDone
}

// ShowProgress displays an animated progress bar with steps.
// It returns a function that should be called when each step completes.
func ShowProgress(totalSteps int) func(string) {
	fmt.Println()

	completed := 0
	mu := sync.Mutex{}

	done := func(description string) {
		mu.Lock()
		defer mu.Unlock()

		// Stop any active spinner
		stopSpinner()

		completed++
		percent := (completed * 100) / totalSteps
		barWidth := 30
		filled := (completed * barWidth) / totalSteps
		empty := barWidth - filled

		bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)

		fmt.Printf("  \033[32m✓\033[0m %-25s \033[36m[%s]\033[0m \033[33m%d%%\033[0m\n", description, bar, percent)

		// Start spinner for next step
		if completed < totalSteps {
			nextSteps := []string{
				"Scaffolding project",
				"Fetching template",
				"Configuring project",
				"Applying fixes",
				"Installing dependencies",
				"Preparing environment",
				"Initializing git",
			}
			if completed < len(nextSteps) {
				startSpinner(nextSteps[completed])
			}
		}
	}

	// Start first spinner
	startSpinner("Scaffolding project")

	return done
}
