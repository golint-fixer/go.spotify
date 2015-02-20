// +build windows

package spotify

import (
	"encoding/csv"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

func pid(name string) (int32, error) {
	out, err := outerr(exec.Command("tasklist.exe",
		"/FI", "IMAGENAME eq "+name, "/FO", "CSV").CombinedOutput())
	if err != nil {
		return 0, errorf("failed to get PID: %q; out: %q", err, out)
	}
	if strings.Contains(out, "No tasks are running") {
		return 0, errorf("process %q is not running", name)
	}
	r, err := csv.NewReader(strings.NewReader(out)).ReadAll()
	if err != nil {
		return 0, errorf("failed to parse PID: %q; out: %q", err, out)
	}
	pid, p := int32(math.MaxInt32), int64(0)
	for i := 1; i < len(r); i++ {
		if p, err = strconv.ParseInt(r[i][1], 10, 32); err != nil {
			return 0, errorf("PID is invalid: %q; data: %q", err, r[i][1])
		}
		pid = min(pid, int32(p))
	}
	return pid, nil
}

func (a *App) kill() (err error) {
	_, err = exec.Command("taskkill.exe", "/PID",
		strconv.Itoa(a.cmd.Process.Pid), "/F", "/T").CombinedOutput()
	return
}
