package spotify

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pblaszczyk/go.spotify/Godeps/_workspace/src/github.com/pblaszczyk/go.utils"
)

type tdf func()

// copyexec is a util function used to copy executable file to d file placed
// in a temporary directory. It returns teardown function and absolute path
// to newly created file. If e is false, destination file has permissions
// set to 666.
func copyexec(t *testing.T, d, s string, i int) (tdf, string, error) {
	td, name, err := func() {}, "", error(nil)
	tmp, err := ioutil.TempDir("", "go.spotify")
	if err != nil {
		t.Errorf("copyexec: TempDir failed: %q (%d)", err, i)
		return td, name, err
	}
	name = exext(filepath.Join(tmp, filepath.Base(d)), s)
	if err = utils.CopyFile(name, s); err != nil {
		t.Errorf("copyexec: copying %q to %q failed: %q (%d)", s, d, err, i)
		return td, name, err
	}
	td = func() {
		if err := os.RemoveAll(filepath.Dir(name)); err != nil {
			t.Logf("copyexec: failed to remove %q: %q (%d)", d, err, i)
		}
	}
	return td, name, err
}

func exext(name, source string) string {
	if strings.HasSuffix(source, ".exe") {
		return name + ".exe"
	}
	return name
}
