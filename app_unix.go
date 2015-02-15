// +build !windows

package spotify

import (
	"math"
	"os/exec"
	"strconv"
	"strings"
)

func pid(name string) (int32, error) {
	out, err := outerr(exec.Command("pidof", name).Output())
	if err != nil {
		return 0, errorf("failed to get PID: %q; out: %q", err, out)
	}
	l := strings.Split(strings.TrimSpace(out), " ")
	if len(l) < 1 {
		return 0, errorf("failed to parse PID: %q; out: %q", err, out)
	}
	pid, p := int32(math.MaxInt32), int64(0)
	for _, v := range l {
		if p, err = strconv.ParseInt(v, 10, 32); err != nil {
			return 0, errorf("PID is invalid: %q; data: %q", err, p)
		}
		pid = min(pid, int32(p))
	}
	return pid, nil
}
