package core

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"sync"
)

// HookResult contains the result of an external command execution.
type HookResult struct {
	Command  string
	Stdout   string
	Stderr   string
	ExitCode int
	Success  bool
}

// HookError is returned when one or more hooks fail.
type HookError struct {
	Results []HookResult
}

func (e *HookError) Error() string {
	return "one or more hooks failed"
}

// ExecuteHooks runs multiple commands in parallel and returns their results.
// It returns a *HookError if any of the commands fail (non-zero exit code).
func ExecuteHooks(ctx context.Context, commands []string) ([]HookResult, error) {
	if len(commands) == 0 {
		return nil, nil
	}

	results := make([]HookResult, len(commands))
	var wg sync.WaitGroup
	var mu sync.Mutex
	var hasFailure bool

	for i, cmdStr := range commands {
		wg.Add(1)
		go func(i int, cmdStr string) {
			defer wg.Done()

			parts := strings.Fields(cmdStr)
			if len(parts) == 0 {
				return
			}

			var stdout, stderr bytes.Buffer
			// #nosec G204 - User-defined hooks require variable command execution
			cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			
			res := HookResult{
				Command: cmdStr,
				Stdout:  stdout.String(),
				Stderr:  stderr.String(),
				Success: true,
			}

			if err != nil {
				res.Success = false
				if exitErr, ok := err.(*exec.ExitError); ok {
					res.ExitCode = exitErr.ExitCode()
				} else {
					res.ExitCode = -1
				}
				mu.Lock()
				hasFailure = true
				mu.Unlock()
			}

			results[i] = res
		}(i, cmdStr)
	}

	wg.Wait()

	if hasFailure {
		return results, &HookError{Results: results}
	}

	return results, nil
}
