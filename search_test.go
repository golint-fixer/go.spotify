package spotify

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestIsEOF(t *testing.T) {
	t.Parallel()
	for i, c := range []struct {
		e error
		r bool
	}{
		{errEOF, true}, {errors.New(""), false}, {nil, false},
	} {
		if ok := IsEOF(c.e); ok != c.r {
			t.Errorf("want %t==%t (%d)", ok, c.r, i)
		}
	}
}

func TestSearchArtist(t *testing.T) {
	t.Parallel()
	ctrl := NewControl(&Context{
		nil,
		nil,
		web{
			g: &getMock{
				d: []string{
					jsonData(t, "artist_1.json"),
					jsonData(t, "artist_2.json"),
				},
			},
			rl: 5,
		},
		"",
	})
	ch, err := make(chan []Artist), make(chan error, 1)
	ctrl.SearchArtist("", ch, err)
	i := 0
	for c := range ch {
		length := len(searchArtistFixt.res)
		if i >= length {
			t.Errorf("want %d < %d (%d)", i, length, i)
		}
		if !reflect.DeepEqual(c, searchArtistFixt.res[i]) {
			t.Errorf("want %v==%v (%d)", c, searchArtistFixt.res[i], i)
		}
		i++
	}
	expectErr(t, err, i)
}

func expectErr(t *testing.T, err chan error, i int) {
	select {
	case err := <-err:
		if ok := IsEOF(err); !ok {
			t.Errorf("want %t==false (%d)", ok, i)
		}
	default:
		t.Fail()
	}
}

func TestSearchAlbum(t *testing.T) {
	t.Parallel()
	ctrl := NewControl(&Context{
		nil,
		nil,
		web{
			&getMock{
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
			5,
		},
		"",
	})
	ch, err := make(chan []Album), make(chan error, 1)
	ctrl.SearchAlbum("", ch, err)
	i := 0
	for c := range ch {
		length := len(searchAlbumFixt.res)
		if i >= length {
			t.Errorf("want %d < %d (%d)", i, length, i)
		}
		for j := range c {
			if !reflect.DeepEqual(c[j], searchAlbumFixt.res[i][j]) {
				t.Errorf("want %v=%v (%d)", c[j], searchAlbumFixt.res[i][j], i)
			}
		}
		i++
	}
	expectErr(t, err, i)
}

func TestSearchTrack(t *testing.T) {
	t.Parallel()
	ctrl := NewControl(&Context{
		nil,
		nil,
		web{
			&getMock{
				d: []string{
					jsonData(t, "track_1.json"),
					jsonData(t, "track_2.json"),
				},
			},
			2,
		},
		"",
	})
	ch, err := make(chan []Track), make(chan error, 1)
	ctrl.SearchTrack("", ch, err)
	i := 0
	for c := range ch {
		length := len(searchTrackFixt.res)
		if i >= length {
			t.Errorf("want %d < %d (%d)", i, length, i)
		}
		if !reflect.DeepEqual(c, searchTrackFixt.res[i]) {
			t.Errorf("want %v==%v (%d)", c, searchTrackFixt.res[i], i)
		}
		i++
	}
	expectErr(t, err, i)
}

func TestSearchArtistError(t *testing.T) {
	t.Parallel()
	ctrl := NewControl(&Context{
		nil,
		nil,
		web{
			g: &getMock{
				d: []string{
					jsonData(t, "error_1.json"),
				},
			},
		},
		"",
	})
	ch, err := make(chan []Artist), make(chan error, 1)
	ctrl.SearchArtist("", ch, err)
	select {
	case err := <-err:
		if err == nil {
			t.Errorf("want %v != nil", err)
		}
		if ok := strings.Contains(err.Error(), "code"); !ok {
			t.Errorf("want %t==true", ok)
		}
	}
}
