// webapi stores function for interacting with Spotify Metadata API
package webapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/pblaszczyk/sscc/model"

	"github.com/cheggaaa/pb"
)

const (
	searchUrl    = "http://ws.spotify.com/search/1/%s.json?q=%s&page=%d"
	searchArtist = "artist"
	searchAlbum  = "album"
	searchTrack  = "track"
)

var (
	bar *pb.ProgressBar
	// Bar specifies if progress bar should be displayed.
	Bar bool = false
)

var (
	errInvResp     = errors.New("webapi: creating response failed!")
	errUnsupSearch = errors.New("webapi: unsupported search keyword")
)

type responser interface {
	data() interface{}
}

type respHeader struct {
	ResCount int `json:"num_results"`
	Limit    int `json:"limit"`
	Offset   int `json:"offset"`
	Page     int `json:"page"`
}

type (
	artist struct {
		Uri  string `json:"href"`
		Name string `json:"name"`
	}
	artistResp struct {
		Info     respHeader `json:"info"`
		*artists `json:"artists"`
	}
	artists []artist
)

type (
	album struct {
		Uri     string   `json:"href"`
		Name    string   `json:"name"`
		Artists []artist `json:"artists"`
	}
	albumResp struct {
		Info    respHeader `json:"info"`
		*albums `json:"albums"`
	}
	albums []album
)

type (
	track struct {
		Uri     string   `json:"href"`
		Name    string   `json:"name"`
		Album   album    `json:"album"`
		Artists []artist `json:"artists"`
	}
	trackResp struct {
		Info    respHeader `json:"info"`
		*tracks `json:"tracks"`
	}
	tracks []track
)

func (a *artists) data() interface{} {
	var res []model.Artist
	for _, a := range []artist(*a) {
		res = append(res, model.Artist{Uri: a.Uri, Name: a.Name})
	}
	return res
}

func (a *albums) data() interface{} {
	var res []model.Album
	for _, a := range []album(*a) {
		var arts []model.Artist
		for _, art := range a.Artists {
			arts = append(arts, model.Artist{Uri: art.Uri, Name: art.Name})
		}
		res = append(res, model.Album{Uri: a.Uri, Name: a.Name, Artists: arts})
	}
	return res
}

func (a *tracks) data() interface{} {
	var res []model.Track
	for _, a := range []track(*a) {
		var arts []model.Artist
		for _, art := range a.Artists {
			arts = append(arts, model.Artist{Uri: art.Uri, Name: art.Name})
		}
		res = append(res, model.Track{Uri: a.Uri, Name: a.Name,
			AlbumUri: a.Album.Uri, AlbumName: a.Album.Name, Artists: arts})
	}
	return res
}

const (
	resCount = "ResCount"
	limit    = "Limit"
	info     = "Info"
	offset   = "Offset"
)

var getF = http.Get

var barF = func(inf reflect.Value, p int) {
	if p == 1 {
		bar = pb.StartNew((int)(inf.FieldByName(resCount).Int()))
	}
	if p*(int)(inf.FieldByName(limit).Int()) >
		(int)(inf.FieldByName(resCount).Int()) {
		bar.Set((int)(inf.FieldByName(resCount).Int()))
	} else {
		bar.Set(p * (int)(inf.FieldByName(limit).Int()))
	}
}

var respF = func(search, val string, p int, resp interface{}) (bool, error) {
	r, err := getF(fmt.Sprintf(searchUrl, search, url.QueryEscape(val), p))
	if err != nil {
		return false, err
	}
	defer r.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return false, err
	}
	if err = json.Unmarshal(body, &resp); err != nil {
		return false, err
	}
	inf := reflect.ValueOf(resp).Elem().FieldByName(info)
	if inf.Kind() == reflect.Invalid {
		return false, errInvResp
	}
	if Bar {
		barF(inf, p)
	}
	return inf.FieldByName(offset).Int()+inf.FieldByName(limit).Int() >=
		inf.FieldByName(resCount).Int(), nil
}

func search(search, val string) (interface{}, error) {
	var (
		v   responser
		res reflect.Value
	)
	switch search {
	case searchArtist:
		v = &artistResp{}
		res = reflect.New(reflect.TypeOf([]model.Artist{}))
	case searchAlbum:
		v = &albumResp{}
		res = reflect.New(reflect.TypeOf([]model.Album{}))
	case searchTrack:
		v = &trackResp{}
		res = reflect.New(reflect.TypeOf([]model.Track{}))
	default:
		return nil, errUnsupSearch
	}
	p := 1
	for {
		if eof, err := respF(search, val, p, v); err != nil {
			return nil, err
		} else {
			r := v.data()
			res.Elem().Set(reflect.AppendSlice(res.Elem(), reflect.ValueOf(r)))
			if eof {
				if Bar {
					bar.Finish()
				}
				return res.Elem().Interface(), nil
			}
			p++
		}
	}
}

// Search for artist.
func SearchArtist(artist string) ([]model.Artist, error) {
	if r, err := search(searchArtist, artist); err != nil {
		return nil, err
	} else {
		return r.([]model.Artist), nil
	}
}

// Search for album.
func SearchAlbum(album string) ([]model.Album, error) {
	if r, err := search(searchAlbum, album); err != nil {
		return nil, err
	} else {
		return r.([]model.Album), nil
	}
}

// Search for track.
func SearchTrack(track string) ([]model.Track, error) {
	if r, err := search(searchTrack, track); err != nil {
		return nil, err
	} else {
		return r.([]model.Track), nil
	}
}
