package utils

import (
	"errors"
	"reflect"
	"testing"
)

func TestIsInvalidEnv(t *testing.T) {
	t.Parallel()
	cases := []struct {
		err error
		ok  bool
	}{
		{err: errors.New(""), ok: false},
		{err: &ErrInvalidEnv{}, ok: true},
	}
	for i := range cases {
		if ok := IsInvalidEnv(cases[i].err); ok != cases[i].ok {
			t.Errorf("want %t=%t (%d)", ok, cases[i].ok, i)
		}
	}
}

func TestPrependPathEnv(t *testing.T) {
	t.Parallel()
	cases := []struct {
		env  string
		args []string
		res  string
	}{
		{
			env:  "p1" + pls + "p2",
			args: []string{"p3", "p4"},
			res:  "p3" + pls + "p4" + pls + "p1" + pls + "p2",
		},
		{
			env:  "",
			args: []string{"p1"},
			res:  "p1",
		},
		{},
		{
			env: "u",
			res: "u",
		},
	}
	for i := range cases {
		if p := PrependPathEnv(cases[i].env, cases[i].args...); p != cases[i].res {
			t.Errorf("want %q==%q (%d)", p, cases[i].res, i)
		}
	}
}

func TestAppendPathEnv(t *testing.T) {
	t.Parallel()
	cases := []struct {
		env  string
		args []string
		res  string
	}{
		{
			env:  "p1" + pls + "p2",
			args: []string{"p3", "p4"},
			res:  "p1" + pls + "p2" + pls + "p3" + pls + "p4",
		},
		{
			env:  "",
			args: []string{"p1"},
			res:  "p1",
		},
		{},
		{
			env: "u",
			res: "u",
		},
	}
	for i := range cases {
		if p := AppendPathEnv(cases[i].env, cases[i].args...); p != cases[i].res {
			t.Errorf("want %q==%q (%d)", p, cases[i].res, i)
		}
	}
}

func TestPrependPathEnvs(t *testing.T) {
	t.Parallel()
	cases := []struct {
		envs  []string
		name  string
		args  []string
		res   []string
		isnil bool
	}{
		{
			envs:  []string{"=:C=s", "u=", "a=c"},
			name:  "a",
			args:  []string{"j", "z"},
			res:   []string{"=:C=s", "u=", "a=j" + pls + "z" + pls + "c"},
			isnil: true,
		},
		{
			envs:  []string{"u=", "a=c"},
			name:  "u",
			args:  nil,
			res:   []string{"u=", "a=c"},
			isnil: true,
		},
		{
			envs:  []string{"u="},
			name:  "j",
			args:  nil,
			res:   []string{"u=", "j="},
			isnil: true,
		},
		{
			envs:  []string{"u="},
			name:  "j",
			args:  []string{"s"},
			res:   []string{"u=", "j=s"},
			isnil: true,
		},
		{
			envs:  []string{"virko"},
			name:  "s",
			isnil: false,
		},
	}
	for i := range cases {
		cop := make([]string, len(cases[i].envs))
		copy(cop, cases[i].envs)
		err := PrependPathEnvs(&cop, cases[i].name, cases[i].args...)
		if (err == nil) != cases[i].isnil {
			t.Errorf("want (%v==nil)=%t (%d)", err, cases[i].isnil, i)
		}
		if err != nil {
			continue
		}
		if !reflect.DeepEqual(cop, cases[i].res) {
			t.Errorf("want %v==%v (%d)", cop, cases[i].res, i)
		}
	}
}

func TestAppendPathEnvs(t *testing.T) {
	t.Parallel()
	cases := []struct {
		envs  []string
		name  string
		args  []string
		res   []string
		isnil bool
	}{
		{
			envs:  []string{"u=", "a=c"},
			name:  "a",
			args:  []string{"j", "z"},
			res:   []string{"u=", "a=c" + pls + "j" + pls + "z"},
			isnil: true,
		},
		{
			envs:  []string{"u=", "a=c"},
			name:  "u",
			args:  nil,
			res:   []string{"u=", "a=c"},
			isnil: true,
		},
		{
			envs:  []string{"u="},
			name:  "j",
			args:  nil,
			res:   []string{"u=", "j="},
			isnil: true,
		},
		{
			envs:  []string{"u="},
			name:  "j",
			args:  []string{"s"},
			res:   []string{"u=", "j=s"},
			isnil: true,
		},
		{
			envs:  []string{"virko"},
			name:  "s",
			isnil: false,
		},
	}
	for i := range cases {
		cop := make([]string, len(cases[i].envs))
		copy(cop, cases[i].envs)
		err := AppendPathEnvs(&cop, cases[i].name, cases[i].args...)
		if (err == nil) != cases[i].isnil {
			t.Errorf("want (%v==nil)=%t (%d)", err, cases[i].isnil, i)
		}
		if err != nil {
			continue
		}
		if !reflect.DeepEqual(cop, cases[i].res) {
			t.Errorf("want %v==%v (%d)", cop, cases[i].res, i)
		}
	}
}

func TestReplacePathEnvs(t *testing.T) {
	t.Parallel()
	cases := []struct {
		envs  []string
		name  string
		args  []string
		res   []string
		isnil bool
	}{
		{
			envs:  []string{"u=", "a=c"},
			name:  "a",
			args:  []string{"j", "z"},
			res:   []string{"u=", "a=j" + pls + "z"},
			isnil: true,
		},
		{
			envs:  []string{"u=", "a=c"},
			name:  "u",
			args:  nil,
			res:   []string{"u=", "a=c"},
			isnil: true,
		},
		{
			envs:  []string{"u="},
			name:  "j",
			args:  nil,
			res:   []string{"u=", "j="},
			isnil: true,
		},
		{
			envs:  []string{"u="},
			name:  "j",
			args:  []string{"s"},
			res:   []string{"u=", "j=s"},
			isnil: true,
		},
		{
			envs:  []string{"virko"},
			name:  "s",
			isnil: false,
		},
	}
	for i := range cases {
		cop := make([]string, len(cases[i].envs))
		copy(cop, cases[i].envs)
		err := ReplacePathEnvs(&cop, cases[i].name, cases[i].args...)
		if (err == nil) != cases[i].isnil {
			t.Errorf("want (%v==nil)=%t (%d)", err, cases[i].isnil, i)
		}
		if err != nil {
			continue
		}
		if !reflect.DeepEqual(cop, cases[i].res) {
			t.Errorf("want %v==%v (%d)", cop, cases[i].res, i)
		}
	}
}
