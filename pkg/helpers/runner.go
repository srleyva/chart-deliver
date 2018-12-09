package helpers

import (
	"os/exec"
)

// Runner is the interface to describe the shell executor
type Runner interface {
	Run(string, ...string) ([]byte, error)
}

// HelmHandler Actually runs shell commands
type HelmHandler struct{}

// NewHelmHandler returns a HelmHandler
func NewHelmHandler() *HelmHandler {
	return &HelmHandler{}
}

// Run will take a command and args and return the output and err returned from shell command
func (h *HelmHandler) Run(cmd string, args ...string) ([]byte, error) {
	return exec.Command(cmd, args...).CombinedOutput()
}
