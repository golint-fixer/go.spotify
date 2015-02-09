package spotify

import (
	"io"
	"net/http"
)

type getMock struct {
	d []string
	i uint
}

func (g *getMock) get(req string) (r *http.Response, err error) {
	r = &http.Response{Body: &rcMock{data: g.d[g.i]}}
	g.i++
	return
}

type rcMock struct {
	data   string
	ready  bool
	offset int
}

func (*rcMock) Close() (_ error) { return }

func (rc *rcMock) Read(p []byte) (n int, err error) {
	l := []byte(rc.data)
	cnt := 0
	for i := rc.offset; i < rc.offset+len(p) && i < len(l); i++ {
		p[i-rc.offset] = l[i]
		cnt++
	}
	rc.offset += len(p)
	if rc.ready {
		return 0, io.EOF
	}
	if rc.offset >= len(l) {
		rc.ready = true
		return cnt, nil
	}
	return cnt, nil
}
