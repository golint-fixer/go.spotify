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
			{URI: "some_uri", Name: " name1"},
			{URI: " second_uri", Name: "name2"},
			{URI: "uri 3", Name: "name 3"},
		},
			[]model.Artist{
				{URI: "some_uri", Name: " name1"},
				{URI: " second_uri", Name: "name2"},
				{URI: "uri 3", Name: "name 3"},
			},
		},
		{},
		{
			artists{
				{URI: "sth", Name: "urk"},
			},
			[]model.Artist{
				{URI: "sth", Name: "urk"},
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
			{URI: "some_uri", Name: " name1",
				Artists: []artist{
					{URI: "u1", Name: "n1"}, {URI: "u2", Name: "n2"}}},
			{URI: " second_uri", Name: "name2",
				Artists: []artist{
					{URI: "u3", Name: "n2"},
					{URI: "uri 3", Name: "name 3"}}},
		},
			[]model.Album{
				{URI: "some_uri", Name: " name1",
					Artists: []model.Artist{
						{URI: "u1", Name: "n1"}, {URI: "u2", Name: "n2"},
					}},
				{URI: " second_uri", Name: "name2",
					Artists: []model.Artist{
						{URI: "u3", Name: "n2"}, {URI: "uri 3", Name: "name 3"},
					}},
			},
		},
		{},
		{
			albums{
				{URI: "sth", Name: "urk"},
			},
			[]model.Album{
				{URI: "sth", Name: "urk"},
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
			{URI: "some_uri", Name: " name1",
				Album: album{Name: "sur", URI: "kur"},
				Artists: []artist{
					{URI: "u1", Name: "n1"}, {URI: "u2", Name: "n2"}}},
			{URI: " second_uri", Name: "name2",
				Album: album{Name: "rem", URI: " drem"},
				Artists: []artist{
					{URI: "u3", Name: "n2"},
					{URI: "uri 3", Name: "name 3"}}},
		},
			[]model.Track{
				{URI: "some_uri", Name: " name1",
					AlbumName: "sur", AlbumURI: "kur",
					Artists: []model.Artist{
						{URI: "u1", Name: "n1"}, {URI: "u2", Name: "n2"},
					}},
				{URI: " second_uri", Name: "name2",
					AlbumName: "rem", AlbumURI: " drem",
					Artists: []model.Artist{
						{URI: "u3", Name: "n2"}, {URI: "uri 3", Name: "name 3"},
					}},
			},
		},
		{},
		{
			tracks{
				{URI: "sth", Name: "urk"},
			},
			[]model.Track{
				{URI: "sth", Name: "urk"},
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

func mockJSON(data string) {
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
	mockJSON(`
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
