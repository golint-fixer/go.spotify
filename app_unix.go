// +build linux darwin

package spotify

import (
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

func exename(name string) string {
	return name
}

func pid(name string) (int32, error) {
	out, err := outerr(exec.Command("pidof", name).Output())
	if err != nil {
		return 0, fmt.Errorf("spotify: failed to get PID: %q. out: %q", err, out)
	}
	l := strings.Split(strings.TrimSpace(out), " ")
	if len(l) < 1 {
		return 0, fmt.Errorf("spotify: failed to parse PID data: %q", out)
	}
	pid, p := int32(math.MaxInt32), int64(0)
	for _, v := range l {
		if p, err = strconv.ParseInt(v, 10, 32); err != nil {
			return 0, fmt.Errorf("spotify: retrieved PID is invalid: %s", v)
		}
		pid = min(pid, p)
	}
	return pid, nil
}
