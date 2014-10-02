package sscc

import (
	"os"
	"os/exec"
	"reflect"
	"testing"
	"time"

	"github.com/pblaszczyk/gophtu"
)

const testEnv = "sscc_mocked_app"

func TestMockApp(t *testing.T) {
	if os.Getenv(testEnv) == "" {
		t.Skip("helper test only")
	}
	<-time.After(time.Second)
}

func TestStart(t *testing.T) {
	defer func() func() {
		c := cmd
		return func() {
			cmd = c
		}
	}()()
	cfg := []struct {
		cmd *exec.Cmd
		err error
	}{
		{exec.Command(os.Args[0], "-test.run", "TestMockApp"), nil},
	}
	os.Setenv(testEnv, testEnv)
	for _, cfg := range cfg {
		cmd = cfg.cmd
		err := Start()
		gophtu.Check(t, reflect.DeepEqual(err, cfg.err), cfg.err, err)
		if cmd != nil {
			err = Kill()
			gophtu.Assert(t, err == nil, nil, err)
		}
	}
}

func TestStart_attach_Kill(t *testing.T) {
	if os.Getenv("sscc_app_test") == "" {
		t.Skip("sscc_app_test var not set")
	}
	err := Start()
	gophtu.Assert(t, err == nil, nil, err)
	<-time.After(5 * time.Second)
	err = Kill()
	if !gophtu.Check(t, err == nil, nil, err) {
		panic("Failed to kill spotify app")
	}
}
