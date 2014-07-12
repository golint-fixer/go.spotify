package webapi

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/pblaszczyk/gophtu"
	"github.com/pblaszczyk/sscc/model"
)

func Test_artists_data(t *testing.T) {
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
	for i, cfg := range cfg {
		gophtu.Check(t, reflect.DeepEqual(cfg.a.data(), cfg.exp), cfg.a.data(),
			cfg.exp, i)
	}
}

func Test_albums_data(t *testing.T) {
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
	for i, cfg := range cfg {
		gophtu.Check(t, reflect.DeepEqual(cfg.a.data(), cfg.exp), cfg.a.data(),
			cfg.exp, i)
	}
}

func Test_albums_track(t *testing.T) {
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
	for i, cfg := range cfg {
		gophtu.Check(t, reflect.DeepEqual(cfg.a.data(), cfg.exp), cfg.a.data(),
			cfg.exp, i)
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

func Test_respF(t *testing.T) {
	defer func() func() {
		gF := getF
		return func() {
			getF = gF
		}
	}()()
	mockJson(`
{"info": {"num_results": 1, "limit": 100, "offset": 0, "page": 1}, "sths": [{
"name": "mocked name"}, {"name": "mocked name 2"}]}`)
	rur := respMock{}
	eof, err := respF("", "", 0, &rur)
	gophtu.Assert(t, err == nil, nil, err)
	gophtu.Assert(t, eof, true, eof)
	gophtu.Assert(t, len(rur.Sths) == 2, 2, len(rur.Sths))
	gophtu.Check(t, rur.Sths[0].Name == "mocked name", "mocked name",
		rur.Sths[0].Name)
	gophtu.Check(t, rur.Sths[1].Name == "mocked name 2", "mocked name 2",
		rur.Sths[1].Name)
	eof, err = respF("", "", 0, &http.Response{})
	gophtu.Assert(t, err == errInvResp, errInvResp, err)
}

type searchStr []struct {
	search string
	err    error
	gt1    bool
	isnil  bool
}

func Test_SearchArtist(t *testing.T) {
	cfg := searchStr{
		{"In This Moment", nil, true, false}, {"łąśðəæóœę", nil, false, true}}
	for i, cfg := range cfg {
		res, err := SearchArtist(cfg.search)
		gophtu.Assert(t, reflect.DeepEqual(err, cfg.err), cfg.err, err, i)
		if cfg.isnil {
			gophtu.Assert(t, res == nil, res, nil, i)
		} else {
			gophtu.AssertFalse(t, res == nil, nil, i)
		}
		gophtu.CheckE(t, (len(res) >= 1) == cfg.gt1, cfg.gt1, len(res) >= 1,
			fmt.Sprintf("%d should be > 0", len(res)), i)
	}
}

func Test_SearchAlbum(t *testing.T) {
	cfg := searchStr{
		{"Piece by Piece", nil, true, false}, {"łąśðəæóœę", nil, false, true}}
	for i, cfg := range cfg {
		res, err := SearchAlbum(cfg.search)
		gophtu.Assert(t, reflect.DeepEqual(err, cfg.err), cfg.err, err, i)
		if cfg.isnil {
			gophtu.Assert(t, res == nil, res, nil, i)
		} else {
			gophtu.AssertFalse(t, res == nil, nil, i)
		}
		gophtu.CheckE(t, (len(res) >= 1) == cfg.gt1, cfg.gt1, len(res) >= 1,
			fmt.Sprintf("%d should be > 0", len(res)), i)
	}
}

func Test_SearchTrack(t *testing.T) {
	cfg := searchStr{
		{"Run To Your Mama", nil, true, false}, {"łąśðəæóœę", nil, false, true}}
	for i, cfg := range cfg {
		res, err := SearchTrack(cfg.search)
		gophtu.Assert(t, reflect.DeepEqual(err, cfg.err), cfg.err, err, i)
		if cfg.isnil {
			gophtu.Assert(t, res == nil, res, nil, i)
		} else {
			gophtu.AssertFalse(t, res == nil, nil, i)
		}
		gophtu.CheckE(t, (len(res) >= 1) == cfg.gt1, cfg.gt1, len(res) >= 1,
			fmt.Sprintf("%d should be > 0", len(res)), i)
	}
}
