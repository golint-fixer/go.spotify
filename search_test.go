package spotify

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestIsEOF(t *testing.T) {
	t.Parallel()
	cases := []struct {
		err error
		res bool
	}{
		{
			err: errEOF,
			res: true,
		},
		{
			err: errors.New(""),
			res: false,
		},
		{
			err: nil,
			res: false,
		},
	}
	for i, cas := range cases {
		if ok := IsEOF(cas.err); ok != cas.res {
			t.Errorf("want ok=cas.res; %t==%t (%d)", ok, cas.res, i)
		}
	}
}

func TestArtist(t *testing.T) {
	t.Parallel()
	s := &Search{
		get: &getMock{
			d: []string{
				jsonData(t, "artist_1.json"),
				jsonData(t, "artist_2.json"),
			},
		},
		batch: 5,
	}
	ch, err := make(chan []Artist), make(chan error, 1)
	s.Artist("", ch, err)
	for i := 0; ; i++ {
		l := len(searchArtistFixt.res)
		select {
		case c := <-ch:
			if l := l - 1; i >= l {
				t.Errorf("want i<l; %d<%d (%d)", i, l, i)
			}
			if !reflect.DeepEqual(c, searchArtistFixt.res[i]) {
				t.Errorf("want c=searchArtistFixt.res[i]; got %v==%v (%d)",
					c, searchArtistFixt.res[i], i)
			}
		case e := <-err:
			if l := l - 1; i != l {
				t.Errorf("error expected for i=%d; got %d", i, l)
			}
			if !IsEOF(e) {
				t.Errorf("want e=errEOF; err: %q (%d)", e, i)
			}
			return
		}
	}
}

func TestAlbum(t *testing.T) {
	t.Parallel()
	s := NewSearch()
	s = &Search{
		get: &getMock{
			d: []string{
				jsonData(t, "album_1.json"),
				jsonData(t, "album_1_artist_1.json"),
				jsonData(t, "album_1_artist_2.json"),
				jsonData(t, "album_1_artist_1.json"),
				jsonData(t, "album_1_artist_2.json"),
				jsonData(t, "album_1_artist_1.json"),
				jsonData(t, "album_2.json"),
				jsonData(t, "album_1_artist_2.json"),
			},
		},
		batch: 5,
	}
	ch, err := make(chan []Album), make(chan error, 1)
	s.Album("", ch, err)
	for i := 0; ; i++ {
		l := len(searchAlbumFixt.res)
		select {
		case c := <-ch:
			for j := range c {
				if !reflect.DeepEqual(c[j], searchAlbumFixt.res[i][j]) {
					t.Errorf("want %v=%v (%d)", c[j], searchAlbumFixt.res[i][j], i)
				}
			}
		case e := <-err:
			if l := l - 1; i != l {
				t.Errorf("error expected for i=%d; got %d", i, l)
			}
			if !IsEOF(e) {
				t.Errorf("want e=errEOF; err: %q (%d)", e, i)
			}
			return
		}
	}
}

func TestTrack(t *testing.T) {
	t.Parallel()
	s := &Search{
		get: &getMock{
			d: []string{
				jsonData(t, "track_1.json"),
				jsonData(t, "track_2.json"),
			},
		},
		batch: 5,
	}
	ch, err := make(chan []Track), make(chan error, 1)
	s.Track("", ch, err)
	for i := 0; ; i++ {
		l := len(searchTrackFixt.res)
		select {
		case c := <-ch:
			if l := l - 1; i >= l {
				t.Errorf("want i<l; %d<%d (%d)", i, l, i)
			}
			if !reflect.DeepEqual(c, searchTrackFixt.res[i]) {
				t.Errorf("want c=searchTrackFixt.res[i]; got %v==%v (%d)",
					c, searchTrackFixt.res[i], i)
			}
		case e := <-err:
			if l := l - 1; i != l {
				t.Errorf("error expected for i=%d; got %d", i, l)
			}
			if !IsEOF(e) {
				t.Errorf("want e=errEOF; err: %q (%d)", e, i)
			}
			return
		}
	}
}

func TestArtistError(t *testing.T) {
	t.Parallel()
	s := &Search{
		get: &getMock{
			d: []string{
				jsonData(t, "error_1.json"),
			},
		},
	}
	ch, err := make(chan []Artist), make(chan error, 1)
	s.Artist("", ch, err)
	select {
	case err := <-err:
		if err == nil {
			t.Errorf("want %v != nil", err)
		}
		if ok := strings.Contains(err.Error(), "code"); !ok {
			t.Errorf("want %t==true", ok)
		}
	case <-time.After(5 * time.Second):
		t.Error("timeout")
		return
	}
}
