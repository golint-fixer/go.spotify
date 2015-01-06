package sscc

import (
	"reflect"
	"testing"

	"github.com/pblaszczyk/gophtu/asserts"
)

func TestIsEOF(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		e error
		r bool
	}{
		{errEOF, true}, {errUnsupported, false}, {nil, false},
	} {
		asserts.Check(t, IsEOF(c.e) == c.r, IsEOF(c.e), c.r)
	}
}

func TestSearch_Artist(t *testing.T) {
	t.Parallel()
	w := web{
		g: &getMock{
			d: []string{
				jsonData(t, "artist_1.json"),
				jsonData(t, "artist_2.json"),
			},
		},
		rl: 5,
	}
	ch, err := make(chan []Artist), make(chan error, 1)
	w.SearchArtist("", ch, err)
	i := 0
	for c := range ch {
		length := len(searchArtistFixt.res)
		asserts.AssertE(t, i < length, i, length, "want i < length")
		asserts.Check(t, reflect.DeepEqual(c, searchArtistFixt.res[i]),
			c, searchArtistFixt.res[i], i)
		i++
	}
	expectErr(t, err, i)
}

func expectErr(t *testing.T, err chan error, i int) {
	select {
	case err := <-err:
		if !IsEOF(err) {
			asserts.Check(t, IsEOF(err), err, errEOF, i)
		}
	default:
		t.Fail()
	}
}

func TestSearch_Album(t *testing.T) {
	t.Parallel()
	w := web{
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
	}
	ch, err := make(chan []Album), make(chan error, 1)
	w.SearchAlbum("", ch, err)
	i := 0
	for c := range ch {
		length := len(searchAlbumFixt.res)
		asserts.AssertE(t, i < length, i, length, "want i < length")
		for j := range c {
			asserts.Check(t, reflect.DeepEqual(c[j], searchAlbumFixt.res[i][j]),
				c[j], searchAlbumFixt.res[i][j], i, j)
		}
		i++
	}
	expectErr(t, err, i)
}

func TestSearch_Track(t *testing.T) {
	t.Parallel()
	w := web{
		&getMock{
			d: []string{
				jsonData(t, "track_1.json"),
				jsonData(t, "track_2.json"),
			},
		},
		2,
	}
	ch, err := make(chan []Track), make(chan error, 1)
	w.SearchTrack("", ch, err)
	i := 0
	for c := range ch {
		length := len(searchTrackFixt.res)
		asserts.AssertE(t, i < length, i, length, "want i < length")
		asserts.Check(t, reflect.DeepEqual(c, searchTrackFixt.res[i]),
			c, searchTrackFixt.res[i], i)
		i++
	}
	expectErr(t, err, i)
}
