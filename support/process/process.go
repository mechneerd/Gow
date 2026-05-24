package process

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"
)

// Result represents the output of a process.
type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Error    error
}

// Process provides a fluent wrapper for os/exec.
type Process struct {
	cmd     *exec.Cmd
	timeout time.Duration
	faked   bool
	result  *Result
}

var fakesEnabled bool
var fakeResults map[string]*Result

// Fake globally enables faking for testing.
func Fake(results map[string]*Result) {
	fakesEnabled = true
	fakeResults = results
}

// Command creates a new Process instance.
func Command(name string, args ...string) *Process {
	return &Process{
		cmd: exec.Command(name, args...),
	}
}

// Timeout sets a timeout for the process execution.
func (p *Process) Timeout(d time.Duration) *Process {
	p.timeout = d
	return p
}

// Run executes the command synchronously.
func (p *Process) Run() *Result {
	if fakesEnabled {
		cmdStr := p.cmd.Path + " " + strings.Join(p.cmd.Args[1:], " ")
		if res, ok := fakeResults[cmdStr]; ok {
			return res
		}
		if res, ok := fakeResults["*"]; ok {
			return res
		}
	}

	var ctx context.Context
	var cancel context.CancelFunc

	if p.timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), p.timeout)
		defer cancel()
		// Recreate cmd with context
		args := p.cmd.Args[1:]
		p.cmd = exec.CommandContext(ctx, p.cmd.Path, args...)
	}

	var stdout, stderr bytes.Buffer
	p.cmd.Stdout = &stdout
	p.cmd.Stderr = &stderr

	err := p.cmd.Run()

	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = -1 // generic error or timeout
		}
	}

	return &Result{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
		Error:    err,
	}
}

