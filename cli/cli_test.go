package cli

import (
	"testing"

	. "gopkg.in/check.v1"
)

type cliSuite struct {
}

var _ = Suite(&cliSuite{})

func TestCli(t *testing.T) { TestingT(t) }
