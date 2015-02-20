package spotify

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func init() {
	if os.Getenv(testEnv) != "" {
		<-time.After(10 * time.Minute)
		os.Exit(0)
	}
}

func testStart(t *testing.T, name string, isnil, exec bool, i int) {
	td, name, err := copyexec(t, name, os.Args[0], i)
	if err != nil {
		return
	}
	defer td()
	app, err := NewApp(name)
	if err != nil {
		t.Errorf("want err=nil; got %q (%d)", err, i)
		return
	}
	if !exec {
		if runtime.GOOS == "linux" {
			if err = os.Chmod(name, 0666); err != nil {
				t.Errorf("want err=nil; got %q (%d)", err, i)
				return
			}
		} else {
			if err = os.Rename(name, name+"nonexec"); err != nil {
				t.Errorf("want err=nil; got %q (%d)", err, i)
				return
			}
			name += "nonexec"
		}
	}
	if err = os.Setenv(testEnv, "1"); err != nil {
		t.Errorf("want err=nil; got %q (%d)", err, i)
		return
	}
	defer os.Unsetenv(testEnv)
	if err = app.Start(); (err == nil) != isnil {
		t.Errorf("want (er=nil)=isnil; err: %v, isnil: %t (%d)", err, isnil, i)
		return
	}
	if err != nil {
		return
	}
	if err = app.Kill(); err != nil {
		t.Errorf("want err=nil; got %q (%d)", err, i)
	}
}

func TestStart(t *testing.T) {
	cases := []struct {
		exec  bool
		isnil bool
		name  string
	}{
		{
			exec:  true,
			isnil: true,
			name:  "mock1",
		},
		{
			exec:  false,
			isnil: false,
			name:  "mock2",
		},
		{
			exec:  false,
			isnil: false,
			name:  filepath.Base(os.Args[0]),
		},
	}
	for i, cas := range cases {
		testStart(t, cas.name, cas.isnil, cas.exec, i)
	}
}

func TestIsRunning(t *testing.T) {
	t.Parallel()
	cases := []struct {
		err error
		res bool
	}{
		{
			err: ErrIsRunning,
			res: true,
		},
		{
			err: errors.New(""),
			res: false,
		},
		{
			err: nil,
			res: false,
		},
	}
	for i, cas := range cases {
		if res := IsRunning(cas.err); res != cas.res {
			t.Errorf("want res=cas.res; got %t=%t (%d)", res, cas.res, i)
		}
	}
}

func testKill(t *testing.T, start, cop, isnil bool, i int) {
	td, n, err := copyexec(t, "spotifymock", os.Args[0], i)
	if err != nil {
		t.Errorf("want err=nil; got %q (%d)", err, i)
		return
	}
	defer td()
	app, err := NewApp(n)
	if err != nil {
		t.Errorf("want err=nil; got %q (%d)", err, i)
		return
	}
	if start {
		if err = os.Setenv(testEnv, "1"); err != nil {
			t.Errorf("want err=nil; got %q (%d)", err, i)
			return
		}
		defer os.Unsetenv(testEnv)
		if err = app.Start(); err != nil {
			t.Errorf("want err=nil; got %q (%d)", err, i)
			return
		}
	}
	if err = app.Kill(); (err == nil) != isnil {
		t.Errorf("want (err=nil)=isnil; err %v, isnil: %t (%d)", err, isnil, i)
	}
}

func TestKill(t *testing.T) {
	cases := []struct {
		start bool
		cop   bool
		isnil bool
	}{
		{
			start: true,
			cop:   true,
			isnil: true,
		},
		{
			start: false,
			cop:   false,
			isnil: false,
		},
	}
	for i, cas := range cases {
		testKill(t, cas.start, cas.cop, cas.isnil, i)
	}
}

func testAttach(t *testing.T, start, cop, isnil bool, i int) {
	td, n, err := copyexec(t, "spotifymock", os.Args[0], i)
	if err != nil {
		t.Errorf("want err=nil; got %q (%d)", err, i)
		return
	}
	defer td()
	app, err := NewApp(n)
	if err != nil {
		t.Errorf("want err=nil; got %q (%d)", err, i)
		return
	}
	if start {
		if err = os.Setenv(testEnv, "1"); err != nil {
			t.Errorf("want err=nil; got %q (%d)", err, i)
			return
		}
		defer os.Unsetenv(testEnv)
		if err = app.Start(); err != nil {
			t.Errorf("want err=nil; got %q (%d)", err, i)
			return
		}
	}
	if err = app.Attach(); (err == nil) != isnil {
		t.Errorf("want (err=nil)=isnil; err %v, isnil: %t (%d)", err, isnil, i)
	}
	if start {
		if err = app.Kill(); err != nil {
			t.Errorf("want err=nil; got %q (%d)", err, i)
		}
	}
}

func TestAttach(t *testing.T) {
	cases := []struct {
		start bool
		cop   bool
		isnil bool
	}{
		{
			start: true,
			cop:   true,
			isnil: true,
		},
		{
			start: false,
			cop:   false,
			isnil: false,
		},
	}
	for i, cas := range cases {
		testAttach(t, cas.start, cas.cop, cas.isnil, i)
	}
}

func testPing(t *testing.T, start, cop, isnil bool, i int) {
	td, n, err := copyexec(t, "spotifymock", os.Args[0], i)
	if err != nil {
		t.Errorf("want err=nil; got %q (%d)", err, i)
		return
	}
	defer td()
	app, err := NewApp(n)
	if err != nil {
		t.Errorf("want err=nil; got %q (%d)", err, i)
		return
	}
	if start {
		if err = os.Setenv(testEnv, "1"); err != nil {
			t.Errorf("want err=nil; got %q (%d)", err, i)
			return
		}
		defer os.Unsetenv(testEnv)
		if err = app.Start(); err != nil {
			t.Errorf("want err=nil; got %q (%d)", err, i)
			return
		}
	}
	if err = app.Ping(); (err == nil) != isnil {
		t.Errorf("want (err=nil)=isnil; err %v, isnil: %t (%d)", err, isnil, i)
	}
	if start {
		if err = app.Kill(); err != nil {
			t.Errorf("want err=nil; got %q (%d)", err, i)
		}
	}
}

func TestPing(t *testing.T) {
	cases := []struct {
		start bool
		cop   bool
		isnil bool
	}{
		{
			start: true,
			cop:   true,
			isnil: true,
		},
		{
			start: false,
			cop:   false,
			isnil: false,
		},
	}
	for i, cas := range cases {
		testPing(t, cas.start, cas.cop, cas.isnil, i)
	}
}

func TestNewApp(t *testing.T) {
	cases := []struct {
		path  string
		isnil bool
	}{
		{
			path:  os.Args[0],
			isnil: true,
		},
		{
			path:  "not_exist",
			isnil: false,
		},
	}
	for i, cas := range cases {
		app, err := NewApp(cas.path)
		if (err == nil) != cas.isnil {
			t.Errorf("want (err=nil)=cas.isnil; err: %v, cas.isnil: %t (%d)",
				err, cas.isnil, i)
		}
		if cas.isnil {
			if app == nil {
				t.Errorf("want app!=nil (%d)", i)
			}
		}
	}
}

func testConnected(t *testing.T, start, cop, cn bool, i int) {
	td, n, err := copyexec(t, "spotifymock", os.Args[0], i)
	if err != nil {
		t.Errorf("want err=nil; got %q (%d)", err, i)
		return
	}
	defer td()
	app, err := NewApp(n)
	if err != nil {
		t.Errorf("want err=nil; got %q (%d)", err, i)
		return
	}
	if start {
		if err = os.Setenv(testEnv, "1"); err != nil {
			t.Errorf("want err=nil; got %q (%d)", err, i)
			return
		}
		defer os.Unsetenv(testEnv)
		if err = app.Start(); err != nil {
			t.Errorf("want err=nil; got %q (%d)", err, i)
			return
		}
	}
	if res := app.Connected(); res != cn {
		t.Errorf("want res=cn; got %t=%t (%d)", res, cn, i)
	}
	if start {
		if err = app.Kill(); err != nil {
			t.Errorf("want err=nil; got %q (%d)", err, i)
		}
	}
}

func TestConnected(t *testing.T) {
	cases := []struct {
		start bool
		cop   bool
		res   bool
	}{
		{
			start: true,
			cop:   true,
			res:   true,
		},
		{
			start: false,
			cop:   false,
			res:   false,
		},
	}
	for i, cas := range cases {
		testConnected(t, cas.start, cas.cop, cas.res, i)
	}
}
