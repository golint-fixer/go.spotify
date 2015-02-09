package spotify

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

const testEnv = "SPOTIFY_APP_MOCK"

func TestMockApp(t *testing.T) {
	if os.Getenv(testEnv) == "" {
		t.Skip("helper test only")
	}
	<-time.After(time.Minute)
}

func newExecer(name string) Execer {
	return NewExecer(&Context{
		Name: name,
	})
}

func TestStart(t *testing.T) {
	cases := []struct {
		name  string
		isnil bool
	}{
		{name: os.Args[0], isnil: false},
		{name: "nonexist", isnil: false},
	}
	for i, cas := range cases {
		if err := newExecer(cas.name).Start(); (err == nil) != cas.isnil {
			t.Errorf("want (%v==nil)==%t (%d)", err, cas.isnil, i)
		}
	}
}

func TestIsRunning(t *testing.T) {
	t.Parallel()
	cases := []struct {
		err error
		res bool
	}{
		{err: ErrIsRunning, res: true},
		{err: errors.New(""), res: false},
		{err: nil, res: false},
	}
	for i, cas := range cases {
		if res := IsRunning(cas.err); res != cas.res {
			t.Errorf("want %t=%t (%d)", res, cas.res, i)
		}
	}
}
func copyfile(src, dst string) error {
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()
	w, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, r)
	return err
}

func TestKill(t *testing.T) {
	cases := []struct {
		name  string
		start bool
		args  []string
		cop   bool
		isnil bool
	}{
		{
			name:  os.Args[0],
			start: true,
			args:  []string{"-test.run", "TestMockApp"},
			cop:   true,
			isnil: true,
		},
		{
			name:  "not_exist",
			start: false,
			isnil: false,
		},
	}
	for i, cas := range cases {
		n := cas.name
		if cas.cop {
			dst := filepath.Join(os.TempDir(), "temp_sscc_test")
			if err := copyfile(n, dst); err != nil {
				t.Fatalf("want %v==nil (%d)", err, i)
			}
			n = dst
			if err := os.Chmod(n, 0777); err != nil {
				t.Errorf("want %v=nil (%d)", err, i)
			}
		}
		e := execer{
			cmd:  exec.Command(exename(n), cas.args...),
			name: filepath.Base(n),
		}
		if cas.start {
			e.cmd.Env = append(os.Environ(), testEnv+"=1")
			if err := e.Start(); err != nil {
				t.Errorf("want %v=nil (%d)", err, i)
			}
		}
		if err := e.Kill(); (err == nil) != cas.isnil {
			t.Errorf("want (%v==nil)==%t (%d)", err, cas.isnil, i)
		}
		if cas.cop {
			if err := os.Remove(n); err != nil {
				t.Errorf("want %v==nil (%d)", err, i)
			}
		}
	}
}

func TestAttach(t *testing.T) {
	cases := []struct {
		name  string
		isnil bool
	}{
		{name: os.Args[0], isnil: true},
		{name: "non_exist", isnil: false},
	}
	for i, cas := range cases {
		if err := newExecer(cas.name).Attach(); (err == nil) != cas.isnil {
			t.Errorf("want (%v==nil)==%t (%d)", err, cas.isnil, i)
		}
	}
}

func TestPing(t *testing.T) {
	cases := []struct {
		name  string
		isnil bool
	}{
		{name: "not_running", isnil: false},
		{name: os.Args[0], isnil: true},
	}
	for i, cas := range cases {
		if err := newExecer(cas.name).Ping(); (err == nil) != cas.isnil {
			t.Errorf("want (%v==nil)==%t (%d)", err, cas.isnil, i)
		}
	}
}
