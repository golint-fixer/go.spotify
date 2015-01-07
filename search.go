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

// errUnsupported is an error returned if search kind is not supported.
var errUnsupported = errors.New("sscc: unsupported type of Search result")

// errEOF is returned when there is no more data to be returned.
var errEOF = errors.New("sscc: end of response")

// IsEOF returns a boolean indicating if err is known to be returned when there
// is no more data to return.
func IsEOF(err error) bool {
	return err == errEOF
}

// IsNotSupported returns a boolean indicating if err is known to be returned
// when there search type is not supported.
func IsNotSupported(err error) bool {
	return err == errUnsupported
}

// defaultSearch is a default implementaiton of `Searcher`.
var defaultSearch = newWeb()

// Searcher is an interface for searching for requested artists/albums/tracks.
type Searcher interface {
	// SearchArtist searches for requested artists.
	SearchArtist(string, chan<- []Artist, chan<- error)
	// SearchAlbum searches for requested albums.
	SearchAlbum(string, chan<- []Album, chan<- error)
	// SearchTrack searches for requested tracks.
	SearchTrack(string, chan<- []Track, chan<- error)
}

// sendRes sends partial results through channel.
func sendRes(res interface{}, v interface{}) {
	r := reflect.New(reflect.TypeOf(res).Elem())
	r.Elem().Set(reflect.ValueOf(v))
	reflect.ValueOf(res).Send(r.Elem())
}

type cstm func(web, interface{}) error

// respParam is a type o
type respParam struct {
	t string
	r interface{}
	f cstm
}

const (
	timeout = 30 * time.Second // timeout for http call.
)

func senderr(err error, c chan<- error) {
	if c != nil {
		c <- err
	}
}

func (w web) SearchArtist(name string, c chan<- []Artist, err chan<- error) {
	go w.s("artist", name, c, &artistResp{}, err, nil)
}

func (w web) SearchAlbum(name string, c chan<- []Album, err chan<- error) {
	go w.s("album", name, c, &albumResp{}, err, web.lookupAlbums)
}

func (w web) SearchTrack(name string, c chan<- []Track, err chan<- error) {
	go w.s("track", name, c, &trackResp{}, err, nil)
}

// s searches for requested artist/album/track and sends results through
// channel when they are available.
func (w web) s(s, v string, r, resp interface{}, err chan<- error, f cstm) {
	p, e, m := uint(0), error(nil), resp
	defer func() {
		reflect.ValueOf(r).Close()
	}()
	for {
		if e = w.read(s, v, p, w.rl, resp); e != nil && !IsEOF(e) {
			err <- e
			return
		}
		u := conv(resp)
		if f != nil {
			if err2 := f(w, u); err2 != nil {
				err <- err2
			}
		}
		p += w.rl
		sendRes(r, u)
		if IsEOF(e) {
			err <- e
			return
		}
		resp = reflect.New(reflect.TypeOf(m).Elem()).Interface()
	}
}

var errInvResp = errors.New("sscc: creating response failed")

// strings used for interacting with Spotify Web API
const (
	endPointURL    = "https://api.spotify.com/v1/"
	queryURL       = endPointURL + "search?q=%s&type=%s&offset=%d&limit=%d"
	lookupURL      = endPointURL + "%s/%s"
	queryArtist    = "artist"
	queryAlbum     = "album"
	queryTrack     = "track"
	lookupAlbum    = "albums"
	albumURIPrefix = "spotify:album:"
)

const (
	resCount = "Total"
	next     = "Next"
)

// read creates http request based on provided search keyword `s`, value `val`,
// limit of elements to obtain `lim` and stores result in `resp`.
// If no more data is available to return, it returns errEOF and stores
// remaining data in `resp`.
func (w web) read(s, val string, off, lim uint, resp interface{}) error {
	r, err := w.g.get(fmt.Sprintf(queryURL, url.QueryEscape(val), s, off, lim))
	if err != nil {
		return err
	}
	defer r.Body.Close()
	{
		var body []byte
		if body, err = ioutil.ReadAll(r.Body); err != nil {
			return err
		}
		{
			var e webError
			if err = json.Unmarshal(body, &e); err == nil && e.Err.Status != 0 {
				return e
			}
		}
		if err = json.Unmarshal(body, &resp); err != nil {
			return err
		}
	}
	v := reflect.ValueOf(resp).Elem()
	if v.Kind() == reflect.Invalid || v.NumField() != 1 {
		return errInvResp
	}
	if v.Field(0).Kind() == reflect.Invalid ||
		v.Field(0).FieldByName(next).Kind() == reflect.Invalid {
		return errInvResp
	}
	if v.Field(0).FieldByName(next).IsNil() {
		return errEOF
	}
	return nil
}

// lookupAlbums goes through all obtained albums by query of album
// and fills in data structure with information about their artists.
func (w web) lookupAlbums(d interface{}) (err error) {
	res := d.([]Album)
	var r *http.Response
	var body []byte
	var resp albumArtist
	for i := range res {
		if r, err = w.g.get(fmt.Sprintf(lookupURL, lookupAlbum,
			strings.TrimPrefix(res[i].URI, "spotify:album:"))); err != nil {
			return
		}
		if body, err = ioutil.ReadAll(r.Body); err != nil {
			r.Body.Close()
			return
		}
		if err = json.Unmarshal(body, &resp); err != nil {
			r.Body.Close()
			return
		}
		for j := range resp.Artists {
			res[i].Artists = append(res[i].Artists,
				Artist{URI: resp.Artists[j].URI, Name: resp.Artists[j].Name})
		}
		r.Body.Close()
	}
	return
}

type web struct {
	g  geter
	rl uint
}

func newWeb() *web {
	return &web{newGet(), 50}
}

type geter interface {
	get(string) (*http.Response, error)
}

type get struct {
	c *http.Client
}

func (g get) get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return g.c.Do(req)
}

func newGet() geter {
	return get{
		&http.Client{
			Transport: &http.Transport{
				Dial: func(n, a string) (net.Conn, error) {
					return net.DialTimeout(n, a, timeout)
				},
				TLSClientConfig: &tls.Config{},
			},
		},
	}
}
