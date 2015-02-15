package spotify

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// Search implements operations for searching through Spotify Web API.
type Search struct {
	get   geter // get is used for http GET requests.
	batch uint  // batch represents number of read positions.
}

// NewSearch returns Search instance.
func NewSearch() *Search {
	return &Search{
		get:   newGet(),
		batch: 50,
	}
}

// errEOF is returned when there is no more data to be returned.
var errEOF = errorf("end of response")

// IsEOF returns a boolean indicating if err is known to be returned when there
// is no more data to return.
func IsEOF(err error) bool {
	return err == errEOF
}

// Artist searches for requested artists. name is the name of searched artist,
// c chan is used to return found artists and err i used to return
// search errors.
func (s *Search) Artist(name string, c chan<- []Artist, errch chan<- error) {
	go s.search("artist", name, c, &artistResp{}, errch, custom)
}

// Album searches for requested albums. name is the name of searched album,
// c chan is used to return found albums and err i used to return
// search errors.
func (s *Search) Album(name string, c chan<- []Album, errch chan<- error) {
	go s.search("album", name, c, &albumResp{}, errch, (*Search).lookupAlbums)
}

// Track searches for requested tracks. name is the name of searched track,
// c chan is used to return found tracks and err i used to return
// search errors.
func (s *Search) Track(name string, c chan<- []Track, errch chan<- error) {
	go s.search("track", name, c, &trackResp{}, errch, custom)
}

var custom = func(_ *Search, _ interface{}) (_ error) {
	return
}

// search searches for requested artist/album/track and sends results through
// channel when they are available.
func (s *Search) search(tag, value string, r, resp interface{},
	errch chan<- error, f func(*Search, interface{}) error) {
	p, e, m := uint(0), error(nil), resp
	for {
		if e = s.read(tag, value, p, s.batch, resp); e != nil && !IsEOF(e) {
			errch <- e
			return
		}
		u := conv(resp)
		if err := f(s, u); err != nil {
			errch <- err
		}
		p += s.batch
		sendRes(r, u)
		if IsEOF(e) {
			errch <- e
			return
		}
		resp = reflect.New(reflect.TypeOf(m).Elem()).Interface()
	}
}

var errInvResp = errorf("creating response failed")

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

// read creates HTTP request based on provided search keyword t, value val,
// limit of elements to obtain lim and stores result in resp.
// If no more data is available to return, it returns errEOF and stores
// remaining data in resp.
func (s *Search) read(t, val string, off, lim uint, resp interface{}) error {
	r, err := s.get.get(fmt.Sprintf(queryURL, url.QueryEscape(val), t, off, lim))
	if err != nil {
		return err
	}
	defer r.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return err
	}
	if err = unmarshal(body, resp); err != nil {
		return err
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

func unmarshal(body []byte, resp interface{}) error {
	var e webError
	if err := json.Unmarshal(body, &e); err == nil && e.Err.Status != 0 {
		return e
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return err
	}
	return nil
}

// lookupAlbums goes through all obtained albums by query of album
// and fills in data structure with information about their artists.
func (s *Search) lookupAlbums(d interface{}) (err error) {
	r, body, resp := &http.Response{}, []byte(nil), albumArtist{}
	res := d.([]Album)
	for i := range res {
		if r, err = s.get.get(fmt.Sprintf(lookupURL, lookupAlbum,
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
				Artist{
					URI:  resp.Artists[j].URI,
					Name: resp.Artists[j].Name,
				})
		}
		r.Body.Close()
	}
	return
}

// sendRes sends partial results through channel.
func sendRes(res, v interface{}) {
	r := reflect.New(reflect.TypeOf(res).Elem())
	r.Elem().Set(reflect.ValueOf(v))
	reflect.ValueOf(res).Send(r.Elem())
}

const timeout = 30 * time.Second // timeout for HTTP requests.

// geter is an interface for HTTP GET requests.
type geter interface {
	get(string) (*http.Response, error)
}

// get is a control structure implementing geter.
type get struct {
	c *http.Client
}

// get implements geter.
func (g get) get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return g.c.Do(req)
}

// newGet returns a default implementation of geter.
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
