package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testCopyFile(src, dst string, isnil, rm, indir bool, i int, t *testing.T) {
	var srcnt []byte
	var err error
	var srcfi os.FileInfo
	if isnil {
		if srcnt, err = ioutil.ReadFile(src); err != nil {
			t.Errorf("want %v==nil (%d)", err, i)
		}
		if srcfi, err = os.Stat(src); err != nil {
			t.Errorf("want %v==nil (%d)", err, i)
		}
	}
	if err = CopyFile(dst, src); (err == nil) != isnil {
		t.Errorf("want (%v==nil)==%t (%d)", err, isnil, i)
	}
	if err != nil {
		return
	}
	dest := dst
	if indir {
		dest = filepath.Join(dst, filepath.Base(src))
	}
	checkfile(t, dest, srcnt, srcfi, i)
	if rm {
		os.Remove(dest)
	}
}

func checkfile(t *testing.T, dest string, srcnt []byte,
	srcfi os.FileInfo, i int) {
	dstfi, err := os.Stat(dest)
	if err != nil {
		t.Errorf("want %v==nil (%d)", err, i)
	}
	if d, s := dstfi.Mode(), srcfi.Mode(); d != s {
		t.Errorf("want %d=%d (%d)", d, s, i)
	}
	dstcnt, err := ioutil.ReadFile(dest)
	if err != nil || string(srcnt) != string(dstcnt) {
		t.Errorf("want %v==nil && %q=%q (%d)",
			string(srcnt), string(dstcnt), err, i)
	}
}

func TestCopyFile(t *testing.T) {
	cases := []struct {
		src   string
		dst   string
		isnil bool
		rm    bool
		indir bool
	}{
		{os.Args[0], os.Args[0], true, false, false},
		{os.Args[0], "temp_file", true, true, false},
		{"not_exit", "temp_file", false, false, false},
		{os.Args[0], "", false, false, false},
		{os.Args[0], ".", true, true, true},
		{filepath.Dir(os.Args[0]), "t", false, false, false},
	}
	for i, cas := range cases {
		testCopyFile(cas.src, cas.dst, cas.isnil, cas.rm, cas.indir, i, t)
	}
}

func TestExists(t *testing.T) {
	cases := []struct {
		path  string
		ok    bool
		isnil bool
	}{
		{os.Args[0], true, true}, {"not_existing23", false, true},
	}
	for i, cas := range cases {
		ok, err := Exists(cas.path)
		if (err == nil) != cas.isnil {
			t.Errorf("want (%v==nil)==%t (%d)", err, cas.isnil, i)
		}
		if err != nil {
			continue
		}
		if ok != cas.ok {
			t.Errorf("want %t==%t (%d)", ok, cas.ok, i)
		}
	}
}

func TestIsDir(t *testing.T) {
	cases := []struct {
		path  string
		ok    bool
		isnil bool
	}{
		{os.Args[0], false, true}, {filepath.Dir(os.Args[0]), true, true},
		{"not_exist", false, false},
	}
	for i, cas := range cases {
		ok, err := IsDir(cas.path)
		if (err == nil) != cas.isnil {
			t.Errorf("want (%v==nil)==%t (%d)", err, cas.isnil, i)
		}
		if err != nil {
			continue
		}
		if ok != cas.ok {
			t.Errorf("want %t==%t (%d)", ok, cas.ok, i)
		}
	}
}

func TestIsTheSame(t *testing.T) {
	cases := []struct {
		lf    string
		rf    string
		ok    bool
		isnil bool
	}{
		{os.Args[0], os.Args[0], true, true},
		{"", ".", true, true},
		{"dd", "u", false, true},
		{"dir", filepath.Join("dir", "..", "dir"), true, true},
	}
	for i, cas := range cases {
		ok, err := IsTheSame(cas.lf, cas.rf)
		if (err == nil) != cas.isnil {
			t.Errorf("want (%v==nil)==%t (%d)", err, cas.isnil, i)
		}
		if err != nil {
			continue
		}
		if ok != cas.ok {
			t.Errorf("want %t==%t for %q and %q (%d)", ok, cas.ok, cas.lf, cas.rf, i)
		}
	}
}

func TestCopyDirShort(t *testing.T) {
	cases := []struct {
		src   string
		dst   string
		isnil bool
	}{
		{os.Args[0], filepath.Dir(os.Args[0]), false},
		{filepath.Dir(os.Args[0]), os.Args[0], false},
		{filepath.Dir(os.Args[0]), filepath.Dir(os.Args[0]), true},
		{"not_exit", os.Args[0], false},
		{filepath.Dir(os.Args[0]), "not_exist", false},
	}
	for i, cas := range cases {
		if err := CopyDir(cas.dst, cas.src); (err == nil) != cas.isnil {
			t.Errorf("want (err=nil)==cas.isnil; err: %v, cas.isnil: %t (%d)",
				err, cas.isnil, i)
		}
	}
}

func TestCopyDir(t *testing.T) {
	tmp, err := ioutil.TempDir("", "goutils")
	if err != nil {
		t.Fatalf("want err==nil; err: %q", err)
	}
	defer os.RemoveAll(tmp)
	dirs := []string{
		filepath.Join("p1", "p2", "p3"),
		filepath.Join("p1", "wlazl", "kotek"),
		filepath.Join("na", "plotek"),
		filepath.Join("p1", "wlazl", "zaba"),
	}
	for _, dir := range dirs {
		if err = os.MkdirAll(filepath.Join(tmp, dir), 0700); err != nil {
			t.Fatalf("want err==nil; err: %q", err)
		}
	}
	if err = filepath.Walk(tmp, func(path string, fi os.FileInfo, err error) error {
		if !fi.IsDir() {
			return nil
		}
		return CopyFile(filepath.Join(path, filepath.Base(os.Args[0])), os.Args[0])
	}); err != nil {
		t.Fatalf("want err==nil; err: %q", err)
	}
	tmpdst, err := ioutil.TempDir("", "goutils")
	if err != nil {
		t.Fatalf("want err=nil; err %q", err)
	}
	defer os.RemoveAll(tmpdst)
	if err = CopyDir(tmpdst, tmp); err != nil {
		t.Fatalf("want err=nil; err: %q", err)
	}
	walk(tmp, tmpdst, t)
	walk(tmpdst, tmp, t)
}

func walk(pth, pth2 string, t *testing.T) error {
	return filepath.Walk(pth, func(path string, fi os.FileInfo, err error) error {
		dstfi, err := os.Stat(strings.Replace(path, pth, pth2, 1))
		if err != nil {
			t.Fatalf("want err=nil; %v", err)
		}
		if dstfi.IsDir() != fi.IsDir() {
			t.Fatalf("want dstfi.IsDir()==fi.IsDir(); %t!=%t for %q",
				dstfi.IsDir(), fi.IsDir(), path)
		}
		return nil
	})
}
