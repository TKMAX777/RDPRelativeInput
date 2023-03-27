package linuxapi

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func GetDisplay() (display string) {
	b, _ := exec.Command("bash", "-c", "ps -x|grep Xorg").Output()

	var out = string(b)
	var disps = strings.Split(out, "/usr/lib/xorg/Xorg ")
	if len(disps) > 1 {
		return strings.Split(disps[1], " ")[0]
	}

	b, _ = exec.Command("bash", "-c", "ps -x|grep Xvnc").Output()
	out = string(b)
	disps = strings.Split(out, "/usr/bin/Xvnc ")
	if len(disps) > 0 {
		return strings.Split(disps[1], " ")[0]
	}

	return ""
}

func GetDisplaySize(display string) (x, y int) {
	cmd := exec.Command("bash", "-c", "/usr/bin/xrandr", display)
	cmd.Env = append(os.Environ(), "DISPLAY="+display)
	b, _ := cmd.Output()

	var out = string(b)
	if !strings.Contains(out, ",") {
		return
	}

	out = strings.TrimSpace(strings.Split(out, ",")[1])
	if !strings.Contains(out, "current") {
		return
	}

	var axis = strings.Split(out, " ")

	X, _ := strconv.Atoi(axis[1])
	Y, _ := strconv.Atoi(axis[3])

	return X, Y
}
