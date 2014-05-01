package webapi

import (
	"io"
	"net/http"
	"testing"

	"github.com/pblaszczyk/sscc/model"

	. "gopkg.in/check.v1"
)

type webapiSuite struct {
	getF func(url string) (resp *http.Response, err error)
}

var _ = Suite(&webapiSuite{})

func TestWebapi(t *testing.T) { TestingT(t) }

func (s *webapiSuite) SetUpTest(c *C) {
	s.getF = getF
}

func (s *webapiSuite) TearDownTest(c *C) {
	getF = s.getF
}

func (s *webapiSuite) Test_artists_data(c *C) {
	cfg := []struct {
		a   artists
		exp []model.Artist
	}{
		{artists{
			{Uri: "some_uri", Name: " name1"},
			{Uri: " second_uri", Name: "name2"},
			{Uri: "uri 3", Name: "name 3"},
		},
			[]model.Artist{
				{Uri: "some_uri", Name: " name1"},
				{Uri: " second_uri", Name: "name2"},
				{Uri: "uri 3", Name: "name 3"},
			},
		},
		{},
		{
			artists{
				{Uri: "sth", Name: "urk"},
			},
			[]model.Artist{
				{Uri: "sth", Name: "urk"},
			},
		},
	}
	for _, cfg := range cfg {
		c.Check(cfg.a.data(), DeepEquals, cfg.exp)
	}
}

func (s *webapiSuite) Test_albums_data(c *C) {
	cfg := []struct {
		a   albums
		exp []model.Album
	}{
		{albums{
			{Uri: "some_uri", Name: " name1",
				Artists: []artist{
					{Uri: "u1", Name: "n1"}, {Uri: "u2", Name: "n2"}}},
			{Uri: " second_uri", Name: "name2",
				Artists: []artist{
					{Uri: "u3", Name: "n2"},
					{Uri: "uri 3", Name: "name 3"}}},
		},
			[]model.Album{
				{Uri: "some_uri", Name: " name1",
					Artists: []model.Artist{
						{Uri: "u1", Name: "n1"}, {Uri: "u2", Name: "n2"},
					}},
				{Uri: " second_uri", Name: "name2",
					Artists: []model.Artist{
						{Uri: "u3", Name: "n2"}, {Uri: "uri 3", Name: "name 3"},
					}},
			},
		},
		{},
		{
			albums{
				{Uri: "sth", Name: "urk"},
			},
			[]model.Album{
				{Uri: "sth", Name: "urk"},
			},
		},
	}
	for _, cfg := range cfg {
		c.Check(cfg.a.data(), DeepEquals, cfg.exp)
	}
}

func (s *webapiSuite) Test_albums_track(c *C) {
	cfg := []struct {
		a   tracks
		exp []model.Track
	}{
		{tracks{
			{Uri: "some_uri", Name: " name1",
				Album: album{Name: "sur", Uri: "kur"},
				Artists: []artist{
					{Uri: "u1", Name: "n1"}, {Uri: "u2", Name: "n2"}}},
			{Uri: " second_uri", Name: "name2",
				Album: album{Name: "rem", Uri: " drem"},
				Artists: []artist{
					{Uri: "u3", Name: "n2"},
					{Uri: "uri 3", Name: "name 3"}}},
		},
			[]model.Track{
				{Uri: "some_uri", Name: " name1",
					AlbumName: "sur", AlbumUri: "kur",
					Artists: []model.Artist{
						{Uri: "u1", Name: "n1"}, {Uri: "u2", Name: "n2"},
					}},
				{Uri: " second_uri", Name: "name2",
					AlbumName: "rem", AlbumUri: " drem",
					Artists: []model.Artist{
						{Uri: "u3", Name: "n2"}, {Uri: "uri 3", Name: "name 3"},
					}},
			},
		},
		{},
		{
			tracks{
				{Uri: "sth", Name: "urk"},
			},
			[]model.Track{
				{Uri: "sth", Name: "urk"},
			},
		},
	}
	for _, cfg := range cfg {
		c.Check(cfg.a.data(), DeepEquals, cfg.exp)
	}
}

type ReadClose struct {
	data   string
	ready  bool
	offset int
}

func (rc *ReadClose) Close() error {
	return nil
}

func (rc *ReadClose) Read(p []byte) (n int, err error) {
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

type respMock struct {
	Info respHeader `json:"info"`
	Sths []struct {
		Name string `json:"name"`
	} `json:"sths"`
}

func mockJson(data string) {
	getF = func(url string) (resp *http.Response, err error) {
		resp = &http.Response{Body: &ReadClose{data: data}}
		return
	}
	return
}

func (s *webapiSuite) Test_respF(c *C) {
	mockJson(`
{"info": {"num_results": 1, "limit": 100, "offset": 0, "page": 1}, "sths": [{
"name": "mocked name"}, {"name": "mocked name 2"}]}`)
	rur := respMock{}
	eof, err := respF("", "", 0, &rur)
	c.Assert(err, IsNil)
	c.Assert(eof, Equals, true)
	c.Assert(rur.Sths, HasLen, 2)
	c.Check(rur.Sths[0].Name, Equals, "mocked name")
	c.Check(rur.Sths[1].Name, Equals, "mocked name 2")
	eof, err = respF("", "", 0, &http.Response{})
	c.Assert(err, Equals, errInvResp)
}

type searchStr []struct {
	search string
	err    error
	gt1    bool
	resCh  Checker
}

func (s *webapiSuite) Test_SearchArtist(c *C) {
	cfg := searchStr{
		{"In This Moment", nil, true, NotNil}, {"łąśðəæóœę", nil, false, IsNil}}
	for _, cfg := range cfg {
		res, err := SearchArtist(cfg.search)
		c.Assert(err, DeepEquals, cfg.err)
		c.Assert(res, cfg.resCh)
		c.Check(len(res) >= 1, Equals, cfg.gt1)
	}
}

func (s *webapiSuite) Test_SearchAlbum(c *C) {
	cfg := searchStr{
		{"Piece by Piece", nil, true, NotNil}, {"łąśðəæóœę", nil, false, IsNil}}
	for _, cfg := range cfg {
		res, err := SearchAlbum(cfg.search)
		c.Assert(err, DeepEquals, cfg.err)
		c.Assert(res, cfg.resCh)
		c.Check(len(res) >= 1, Equals, cfg.gt1)
	}
}

func (s *webapiSuite) Test_SearchTrack(c *C) {
	cfg := searchStr{
		{"Run To Your Mama", nil, true, NotNil}, {"łąśðəæóœę", nil, false, IsNil}}
	for _, cfg := range cfg {
		res, err := SearchTrack(cfg.search)
		c.Assert(err, DeepEquals, cfg.err)
		c.Assert(res, cfg.resCh)
		c.Check(len(res) >= 1, Equals, cfg.gt1)
	}
}
