package main

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)
import "C"

func StartApplication(rw *VirtualChannelReadWriteCloser, version string) error {
	root := os.Getenv("ProgramW6432")
	if !filepath.IsAbs(root) {
		return syscall.EINVAL
	}
	cmd := exec.CommandContext(rw.received.Context(), filepath.Join(root, "RDPRelativeInput", "RelativeInputClient.exe"))
	// Only the validated version line belongs on stdin. Streaming the channel
	// here would leave cmd.Wait blocked waiting for an unused stdin copier.
	cmd.Stdin = strings.NewReader(version + "\n")
	cmd.Stdout = rw
	// Do not accumulate unbounded stderr or read mstsc's own stderr as commands.
	cmd.Stderr = io.Discard
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	cmd.Env = append(os.Environ(), "RDP_CLIENT_PID="+strconv.Itoa(os.Getpid()))
	return cmd.Run()
}
