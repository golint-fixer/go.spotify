package spotify

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

// NewExecer returns default implementation of `Execer`.
func NewExecer(ctx *Context) Execer {
	n := exename(ctx.name())
	return &execer{
		cmd:  exec.Command(n),
		name: n,
	}
}

// Execer represents available operations on Spotify process.
type Execer interface {
	Start() error  // Start starts Spotify.
	Kill() error   // Kill stops Spotify.
	Attach() error // Attach connects to an already running Spotify.
	Ping() error   // Ping checks if Spotify is already running.
}

// execer is a default implementation of `Execer`.
type execer struct {
	sync.Mutex
	cmd  *exec.Cmd // cmd is used to control execeress.
	name string    // name is the name of executable.
}

// ErrIsRunning is returned if Spotify is already running.
var ErrIsRunning = errors.New("spotify: app is already running")

// IsRunning returns a boolean indicating whether the error is known to report
// that requested execeress is already running.
func IsRunning(err error) bool {
	return err == ErrIsRunning
}

// Start implements `Execer`.
func (e *execer) Start() error {
	e.Lock()
	defer e.Unlock()
	if err := e.Ping(); err == nil {
		return ErrIsRunning
	}
	return e.start()
}

// Kill implements `Execer`.
func (e *execer) Kill() error {
	e.Lock()
	defer e.Unlock()
	if err := e.attach(); err != nil {
		return err
	}
	return e.kill()
}

// Attach implements `Execer`.
func (e *execer) Attach() error {
	e.Lock()
	defer e.Unlock()
	return e.attach()
}

// Ping implements `Execer`.
func (e *execer) Ping() (err error) {
	_, err = pid(e.name)
	return
}

func (e *execer) attach() error {
	pid, err := pid(e.name)
	if err != nil {
		return err
	}
	e.cmd.Process = &os.Process{
		Pid: int(pid),
	}
	return nil
}

func (e *execer) start() error {
	if err := e.cmd.Start(); err != nil {
		return fmt.Errorf("spotify: failed to start: %q", err)
	}
	return nil
}

func (e *execer) kill() error {
	if err := e.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("spotify: failed to stop: %q", err)
	}
	return nil
}

func min(x int32, y int64) int32 {
	if x < int32(y) {
		return x
	}
	return int32(y)
}

func outerr(out []byte, err error) (string, error) {
	o := "<no output>"
	if out != nil {
		o = string(out)
	}
	return o, err
}
