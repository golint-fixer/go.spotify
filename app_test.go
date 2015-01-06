package sscc

import (
	"errors"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/pblaszczyk/gophtu/asserts"
	"github.com/pblaszczyk/gophtu/times"
)

const testEnv = "sscc_mocked_proc"

func TestMockApp(t *testing.T) {
	if os.Getenv(testEnv) == "" {
		t.Skip("helper test only")
	}
	<-time.After(times.Timeout() * 1000)
}

func TestStart(t *testing.T) {
	cases := []struct {
		name  string
		isnil bool
	}{
		{os.Args[0], false}, {"non_exist", false},
	}
	for i, cas := range cases {
		a := proc{exec.Command(cas.name), cas.name}
		err := a.Run()
		asserts.Check(t, (err == nil) == cas.isnil, cas.isnil, err, i)
	}
}

func TestIsRunning(t *testing.T) {
	t.Parallel()
	cases := []struct {
		err error
		res bool
	}{
		{errIsRunning, true}, {errors.New(""), false}, {nil, false},
	}
	for i, cas := range cases {
		res := IsRunning(cas.err)
		asserts.Check(t, res == cas.res, res, cas.res, i)
	}
}

func TestKill(t *testing.T) {
	cases := []struct {
		name  string
		start bool
		isnil bool
	}{
		{"non_exist", false, false},
	}
	for i, cas := range cases {
		a := &proc{exec.Command(cas.name), cas.name}
		if cas.start {
			err := a.Run()
			asserts.Assert(t, err == nil, nil, err, i)
		}
		err := a.Kill()
		asserts.Check(t, (err == nil) == cas.isnil, cas.isnil, err, i)
	}
}

func TestRun(t *testing.T) {
	cases := []struct {
		name  string
		isnil bool
	}{
		{os.Args[0], false}, {"alsdj", false},
	}
	for i, cas := range cases {
		p := proc{exec.Command(cas.name, "test.run", "TestMockApp"), cas.name}
		err := p.Run()
		asserts.Check(t, (err == nil) == cas.isnil, cas.isnil, err, i)
		if cas.isnil {
			p.Kill()
		}
	}
}

func TestAttach(t *testing.T) {
	cases := []struct {
		name  string
		isnil bool
	}{
		{os.Args[0], true}, {"non_exist", false},
	}
	for i, cas := range cases {
		a := proc{exec.Command(cas.name), cas.name}
		err := a.Attach()
		asserts.Check(t, (err == nil) == cas.isnil, cas.isnil, err, i)
	}
}

func TestPing(t *testing.T) {
	cases := []struct {
		name  string
		isnil bool
	}{
		{"not_running", false}, {os.Args[0], true},
	}
	for i, cas := range cases {
		a := proc{exec.Command(cas.name), cas.name}
		err := a.Ping()
		asserts.Check(t, (err == nil) == cas.isnil, cas.isnil, err, i)
	}
}
