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

// procName is a name of Spotify desktop application process.
const procName = "spotify"

// Procer is an interface for operations on Spotify desktop application.
type Procer interface {
	Run() error    // Run starts the app.
	Kill() error   // Kill stops the app.
	Attach() error // Attach connects to an already running app.
	Ping() error   // Ping checks if the app is currently running.
}

// errIsRunning is returned if Spotify desktop application is already running
// when it's expected to be down.
var errIsRunning = errors.New("sscc: spotify is already running")

// IsRunning returns a boolean indicating whether the error is known to report
// that requested process is already running.
func IsRunning(err error) bool {
	return err == errIsRunning
}

// Run implements `Procer`.
func (p *proc) Run() error {
	if err := p.Ping(); err == nil {
		return errIsRunning
	}
	return p.start()
}

// Kill implements `Procer`.
func (p *proc) Kill() error {
	if err := p.Attach(); err != nil {
		return err
	}
	return p.kill()
}

// Attach implements `Procer`.
func (p *proc) Attach() error {
	pid, err := pid(p.name)
	if err != nil {
		return err
	}
	p.cmd.Process = &os.Process{Pid: int(pid)}
	return nil
}

// Ping implements `Procer`.
func (p *proc) Ping() (err error) {
	_, err = pid(p.name)
	return
}

func (p *proc) start() error {
	if err := p.cmd.Start(); err != nil {
		return fmt.Errorf("sscc: failed to start: %q", err)
	}
	return nil
}

func (p *proc) kill() error {
	if err := p.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("sscc: failed to kill: %q", err)
	}
	return nil
}

type proc struct {
	cmd  *exec.Cmd
	name string
}

func newProc() *proc {
	return &proc{exec.Command(procName), procName}
}

var defaultProc = newProc()

// pid returns PID of running Spotify desktop application.
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
	pid, p := int32(math.MaxInt32), int64(0)
	for _, v := range l {
		if p, err = strconv.ParseInt(v, 10, 32); err != nil {
			return 0, fmt.Errorf("sscc: retrieved PID is invalid: %s", v)
		}
		pid = min(pid, p)
	}
	return pid, nil
}

func min(x int32, y int64) int32 {
	if x < int32(y) {
		return x
	}
	return int32(y)
}
