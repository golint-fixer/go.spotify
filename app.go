package sscc

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// execName is the name of Spotify process.
const execName = "spotify"

// Execer represents Spotify process.
type Execer interface {
	Run() error    // Run starts Spotify.
	Kill() error   // Kill stops Spotify.
	Attach() error // Attach connects to an already running Spotify.
	Ping() error   // Ping checks if Spotify is already running.
}

// errIsRunning is returned if Spotify is already running.
var errIsRunning = errors.New("sscc: spotify is already running")

// IsRunning returns a boolean indicating whether the error is known to report
// that requested process is already running.
func IsRunning(err error) bool {
	return err == errIsRunning
}

// Run implements `Execer`.
func (p *proc) Run() error {
	if err := p.Ping(); err == nil {
		return errIsRunning
	}
	return p.start()
}

// Kill implements `Execer`.
func (p *proc) Kill() error {
	if err := p.Attach(); err != nil {
		return err
	}
	return p.kill()
}

// Attach implements `Execer`.
func (p *proc) Attach() error {
	pid, err := pid(p.name)
	if err != nil {
		return err
	}
	p.cmd.Process = &os.Process{Pid: int(pid)}
	return nil
}

// Ping implements `Execer`.
func (p *proc) Ping() (err error) {
	_, err = pid(p.name)
	return
}

// start starts process.
func (p *proc) start() error {
	if err := p.cmd.Start(); err != nil {
		return fmt.Errorf("sscc: failed to start: %q", err)
	}
	return nil
}

// kill kills process.
func (p *proc) kill() error {
	if err := p.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("sscc: failed to kill: %q", err)
	}
	return nil
}

// proc is a default implementation of `Execer`.
type proc struct {
	cmd  *exec.Cmd // cmd is used to control process.
	name string    // name is the name of executable.
}

// newExecer returns default implementation of `Execer`.
func newExecer() Execer {
	return &proc{exec.Command(execName), execName}
}

// pid returns PID of running Spotify.
func pid(name string) (int32, error) {
	out, err := exec.Command("pidof", name).Output()
	if err != nil {
		o := "<no output>"
		if out != nil {
			o = string(out)
		}
		return 0, fmt.Errorf("sscc: failed to get PID: %v %q", err, o)
	}
	l := strings.Split(strings.TrimSpace(string(out)), " ")
	if len(l) < 1 {
		return 0, fmt.Errorf("sscc: failed to parse PID data: %q", string(out))
	}
	// Returned PID is chosen based on assumption that the lowest PID is the one.
	pid, p := int32(math.MaxInt32), int64(0)
	for _, v := range l {
		if p, err = strconv.ParseInt(v, 10, 32); err != nil {
			return 0, fmt.Errorf("sscc: retrieved PID is invalid: %s", v)
		}
		pid = min(pid, p)
	}
	return pid, nil
}

// min is a helper func returning minimal of 2 values.
func min(x int32, y int64) int32 {
	if x < int32(y) {
		return x
	}
	return int32(y)
}
