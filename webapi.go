package sscc

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

const (
	endPointURL = "https://api.spotify.com/v1/"
	searchURL   = endPointURL + "search?q=%s&type=%s&offset=%d&limit=%d"
	lookupURL   = endPointURL + "%s/%s"

	searchArtist = "artist"
	searchAlbum  = "album"
	searchTrack  = "track"

	lookupAlbum    = "albums"
	albumURIPrefix = "spotify:album:"
)

const (
	// resLimit is maximum number of returned elements by search calls used
	// to obtain artists, albums and tracks.
	resLimit = 50
	// timeout is timeout for http call.
	timeout = 30 * time.Second
)

const (
	resCount = "Total"
	next     = "Next"
)

var (
	client *http.Client
)

func init() {
	client = &http.Client{Transport: &http.Transport{
		Dial: func(n, a string) (net.Conn, error) {
			return net.DialTimeout(n, a, timeout)
		},
		TLSClientConfig: &tls.Config{},
	},
	}
}

var (
	errInvResp     = errors.New("webapi: creating response failed")
	errUnsupSearch = errors.New("webapi: unsupported search keyword")
)

// responser is an interface for json models allowing to convert them
// to models used in application.
type responser interface {
	data() interface{}
}

type respHeader struct {
	Total int     `json:"total"`
	Next  *string `json:"next"`
}

type (
	artist struct {
		URI  string `json:"uri"`
		Name string `json:"name"`
	}
	artists    []artist
	artistResp struct {
		Artists struct {
			Items artists `json:"items"`
			respHeader
		} `json:"artists"`
	}
)

type (
	album struct {
		URI  string `json:"uri"`
		Name string `json:"name"`
	}
	albums    []album
	albumResp struct {
		Albums struct {
			Items albums `json:"items"`
			respHeader
		} `json:"albums"`
	}
	albumArtist struct {
		Artists []struct {
			URI  string `json:"uri"`
			Name string `json:"name"`
		} `json:"artists"`
	}
)

type (
	track struct {
		URI  string `json:"uri"`
		Name string `json:"name"`
	}
	trackData struct {
		Album   album   `json:"album"`
		Artists artists `json:"artists"`
		track
	}
	tracks    []trackData
	trackResp struct {
		Tracks struct {
			Items tracks `json:"items"`
			respHeader
		} `json:"tracks"`
	}
)

func (a *artists) data() interface{} {
	var res []Artist
	for _, a := range []artist(*a) {
		res = append(res, Artist{URI: a.URI, Name: a.Name})
	}
	return res
}

func (a *artistResp) data() interface{} {
	return a.Artists.Items.data()
}

func (a *albums) data() interface{} {
	var res []Album
	for _, a := range []album(*a) {
		res = append(res, Album{URI: a.URI, Name: a.Name})
	}
	return res
}

func (a *albumResp) data() interface{} {
	return a.Albums.Items.data()
}

func (a *tracks) data() interface{} {
	var res []Track
	for _, a := range []trackData(*a) {
		var arts []Artist
		for _, art := range a.Artists {
			arts = append(arts, Artist{URI: art.URI, Name: art.Name})
		}
		res = append(res, Track{URI: a.URI, Name: a.Name,
			AlbumURI: a.Album.URI, AlbumName: a.Album.Name, Artists: arts})
	}
	return res
}

func (a *trackResp) data() interface{} {
	return a.Tracks.Items.data()
}

var getF = func(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	return client.Do(req)
}

var respF = func(s, val string, off, lim int, resp interface{}) (bool, error) {
	r, err := getF(fmt.Sprintf(searchURL, url.QueryEscape(val), s, off, lim))
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
	respV := reflect.ValueOf(resp).Elem()
	if respV.Kind() == reflect.Invalid || respV.NumField() != 1 {
		return false, errInvResp
	}
	if respV.Field(0).Kind() == reflect.Invalid ||
		respV.Field(0).FieldByName(next).Kind() == reflect.Invalid {
		return false, errInvResp
	}
	return respV.Field(0).FieldByName(next).IsNil(), nil
}

// lookupAlbums goes through all obtained albums by search of album
// and fills in data structure with information about their artists.
func lookupAlbums(res *[]Album) error {
	for i := range *res {
		r, err := getF(fmt.Sprintf(lookupURL, lookupAlbum,
			strings.TrimPrefix((*res)[i].URI, "spotify:album:")))
		if err != nil {
			return err
		}
		defer r.Body.Close()
		var body []byte
		if body, err = ioutil.ReadAll(r.Body); err != nil {
			return err
		}
		var resp albumArtist
		if err = json.Unmarshal(body, &resp); err != nil {
			return err
		}
		for j := range resp.Artists {
			(*res)[i].Artists = append((*res)[i].Artists,
				Artist{URI: resp.Artists[j].URI, Name: resp.Artists[j].Name})
		}
	}
	return nil
}

// search runs query used to obtain information about artists, albums or tracks.
func search(search, val string) (interface{}, error) {
	var (
		v   responser
		res reflect.Value
	)
	switch search {
	case searchArtist:
		v = &artistResp{}
		res = reflect.New(reflect.TypeOf([]Artist{}))
	case searchAlbum:
		v = &albumResp{}
		res = reflect.New(reflect.TypeOf([]Album{}))
	case searchTrack:
		v = &trackResp{}
		res = reflect.New(reflect.TypeOf([]Track{}))
	default:
		return nil, errUnsupSearch
	}
	p := 0
	for {
		eof, err := respF(search, val, p, resLimit, v)
		if err != nil {
			return nil, err
		}
		r := v.data()
		res.Elem().Set(reflect.AppendSlice(res.Elem(), reflect.ValueOf(r)))
		if eof {
			return res.Elem().Interface(), nil
		}
		p += resLimit
	}
}

// SearchArtist searches for artist.
func SearchArtist(artist string) ([]Artist, error) {
	r, err := search(searchArtist, artist)
	if err != nil {
		return nil, err
	}
	return r.([]Artist), nil
}

// SearchAlbum searches for album.
func SearchAlbum(album string) ([]Album, error) {
	r, err := search(searchAlbum, album)
	if err != nil {
		return nil, err
	}
	res := r.([]Album)
	err = lookupAlbums(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// SearchTrack searches for track.
func SearchTrack(track string) ([]Track, error) {
	r, err := search(searchTrack, track)
	if err != nil {
		return nil, err
	}
	return r.([]Track), nil
}
