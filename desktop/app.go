// Package desktop stores function for controlling Spotify desktop application
package desktop

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const (
	procName = "spotify"
)

var (
	cmd = exec.Command(procName)
)

// Start spotify process.
func Start() error {
	return cmd.Start()
}

// Kill spotify process.
func Kill() error {
	if cmd.Process == nil {
		if err := attach(); err != nil {
			return err
		}
	}
	return cmd.Process.Kill()
}

func attach() (err error) {
	var (
		pid  int
		pids []int
	)
	out, err := exec.Command("pidof", procName).Output()
	if err != nil {
		return err
	}
	l := strings.Split(strings.Trim(string(out), "\n"), " ")
	if len(l) < 1 {
		return fmt.Errorf("desktop: failed to get PID of %s", procName)
	}
	for _, v := range l {
		if pid, err = strconv.Atoi(v); err != nil {
			return fmt.Errorf("desktop: retrieved PID of %s is invalid: %s",
				procName, v)
		}
		pids = append(pids, pid)
		sort.Ints(pids)
	}
	cmd = exec.Command(procName)
	cmd.Process = &os.Process{Pid: pids[0]}
	return nil
}
