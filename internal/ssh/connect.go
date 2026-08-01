package ssh

import (
	"os/exec"
	"strconv"

	"github.com/luynrs/justssh/internal/models"
)

// just shell out to ssh, no point reimplementing ssh-agent/ProxyJump/ControlMaster
func Command(s models.Server) *exec.Cmd {
	// auto-trust new hosts so first connect doesn't stop on a yes/no prompt,
	// still warns (and refuses) if a known host's key ever changes
	args := []string{"-o", "StrictHostKeyChecking=accept-new"}

	if s.Port != 0 && s.Port != 22 {
		args = append(args, "-p", strconv.Itoa(s.Port))
	}
	if s.Key != "" {
		args = append(args, "-i", s.Key)
	}

	target := s.Host
	if s.User != "" {
		target = s.User + "@" + s.Host
	}
	args = append(args, target)

	return exec.Command("ssh", args...)
}
