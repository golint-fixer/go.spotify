package spotify

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

// NewApp returns new instance of App.
func NewApp(name string, args ...string) (app *App, err error) {
	if name == "" {
		name = "spotify"
	}
	if name, err = exec.LookPath(name); err != nil {
		return
	}
	app = &App{
		cmd:  exec.Command(name, args...),
		name: filepath.Base(name),
	}
	return
}

// App is a representation of Spotify desktop application.
type App struct {
	sync.Mutex
	cmd  *exec.Cmd
	name string
}

// ErrIsRunning is returned if application is already running.
var ErrIsRunning = errorf("app is already running")

// IsRunning returns a boolean indicating whether the error is known to report
// that requested process is already running.
func IsRunning(err error) bool {
	return err == ErrIsRunning
}

// Start starts Spotify desktop application.
func (a *App) Start() (err error) {
	a.Lock()
	if err = a.Ping(); err == nil {
		a.Unlock()
		return ErrIsRunning
	}
	err = a.start()
	a.Unlock()
	return
}

// Kill kills Spotify desktop application.
func (a *App) Kill() (err error) {
	a.Lock()
	if !a.connected() {
		if err = a.attach(); err != nil {
			a.Unlock()
			return err
		}
	}
	err = a.kill()
	a.Unlock()
	return
}

// Attach binds data structure with already running Spotify desktop application.
func (a *App) Attach() (err error) {
	a.Lock()
	err = a.attach()
	a.Unlock()
	return
}

// Ping checks if Spotify desktop application is running.
func (a *App) Ping() (err error) {
	_, err = pid(a.name)
	return
}

// Connected returns a boolean indicating of a is connected to a Spotify
// instance.
func (a *App) Connected() (ok bool) {
	a.Lock()
	ok = a.connected()
	a.Unlock()
	return
}

func (a *App) connected() bool {
	return a.cmd != nil && a.cmd.Process != nil && a.cmd.Process.Pid != 0
}

func (a *App) attach() error {
	pid, err := pid(a.name)
	if err != nil {
		return err
	}
	a.cmd.Process = &os.Process{
		Pid: int(pid),
	}
	return nil
}

func (a *App) start() error {
	if err := a.cmd.Start(); err != nil {
		return errorf("failed to start: %q", err)
	}
	return nil
}
