package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const pls = string(os.PathListSeparator)

var re = regexp.MustCompile(`^(\w+)=(.*)$`)

// ErrInvalidEnv is returned when environment variable to be modified
// is not valid.
type ErrInvalidEnv struct {
	Msg string
	Env string
}

// Error implements `error`.
func (e *ErrInvalidEnv) Error() string {
	return fmt.Sprintf("utils: %q is an invalid env: %q", e.Env, e.Msg)
}

// IsInvalidEnv returns a boolean indicating whether `error` is known
// to be returned when environment variable is not valid.
func IsInvalidEnv(err error) (ok bool) {
	_, ok = err.(*ErrInvalidEnv)
	return
}

// PrependPathEnv returns path like environment variable (PATH, GOPATH, ...)
// with provided values prepended to variable.
func PrependPathEnv(env string, v ...string) (path string) {
	if path = strings.Join(v[:], pls); env != "" {
		if path != "" {
			path += pls
		}
		path += env
	}
	return
}

// AppendPathEnv returns path like environment variable (PATH, GOPATH, ...)
// with provided values appended to variable.
func AppendPathEnv(env string, v ...string) (path string) {
	if path = strings.Join(v[:], pls); env != "" {
		if path != "" {
			path = pls + path
		}
		path = env + path
	}
	return
}

type singleproc func(string, ...string) string

func pp(pr singleproc, env *[]string, name string, v ...string) (err error) {
	var sp []string
	for i, p := range *env {
		// http://blogs.msdn.com/b/oldnewthing/archive/2010/05/06/10008132.aspx
		if strings.Index(p, "=") == 0 {
			continue
		}
		if sp, err = envsplit(p); err != nil {
			return
		}
		if sp[1] != name {
			continue
		}
		(*env)[i] = name + "=" + pr(sp[2], v...)
		return
	}
	*env = append(*env, name+"="+pr("", v...))
	return
}

// PrependPathEnvs gets set of environment variables and prepends to provided
// variable requested arguments if it already exists. If it does not exist,
// adds new variable to existing ones.
func PrependPathEnvs(env *[]string, name string, v ...string) error {
	return pp(PrependPathEnv, env, name, v...)
}

// AppendPathEnvs gets set of environement variables and appends to provided
// variable requested arguments if it already exists. If it does not exist,
// adds new variable to existing ones.
func AppendPathEnvs(env *[]string, name string, v ...string) error {
	return pp(AppendPathEnv, env, name, v...)
}

// ReplacePathEnvs gets set of environment variables and replaces existing
// value of variable with requested arguments if it already exists. If it does
// not exist, adds new variable to existing ones.
func ReplacePathEnvs(env *[]string, name string, v ...string) error {
	return pp(func(_ string, v ...string) string {
		return strings.Join(v, pls)
	}, env, name, v...)
}

func envsplit(env string) (sp []string, err error) {
	if sp = re.FindStringSubmatch(env); len(sp) != 3 {
		return nil, &ErrInvalidEnv{
			Msg: fmt.Sprintf("splitting var failed"),
			Env: env,
		}
	}
	return
}
