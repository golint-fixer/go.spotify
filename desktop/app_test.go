package desktop

import (
	"os"
	"os/exec"
	"testing"
	"time"

	. "github.com/go-check/check"
)

type appSuite struct {
	cmdB *exec.Cmd
}

var _ = Suite(&appSuite{})

func TestApp(t *testing.T) { TestingT(t) }

const testEnv = "sscc_mocked_app"

func TestMockApp(t *testing.T) {
	if os.Getenv(testEnv) == "" {
		return
	}
	<-time.After(time.Second)
}

func (s *appSuite) SetUpTest(c *C) {
	s.cmdB = cmd
}

func (s *appSuite) TearDownTest(c *C) {
	cmd = s.cmdB
}

func (s *appSuite) TestStart(c *C) {
	cfg := []struct {
		cmd *exec.Cmd
		err error
	}{
		{nil, errCmdInit},
		{exec.Command(os.Args[0], "-test.run", "TestMockApp"), nil},
	}
	os.Setenv(testEnv, testEnv)
	for _, cfg := range cfg {
		cmd = cfg.cmd
		c.Check(Start(), DeepEquals, cfg.err)
		if cmd != nil {
			c.Assert(Kill(), IsNil)
		}
	}
}

func (s *appSuite) TestStart_attach_Kill(c *C) {
	if os.Getenv("sscc_app_test") == "" {
		c.Skip("sscc_app_test var not set")
	}
	c.Assert(Start(), IsNil)
	<-time.After(5 * time.Second)
	cmd = nil
	if !c.Check(Kill(), IsNil) {
		panic("Failed to kill spotify app")
	}
}
