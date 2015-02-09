package spotify

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/pblaszczyk/gophtu/asserts"
)

func jsonData(t *testing.T, name string) string {
	b, err := ioutil.ReadFile(filepath.Join("testdata", name))
	asserts.Assert(t, err == nil, err, nil)
	return string(b)
}

var searchArtistFixt = struct {
	res [][]Artist
	err error
}{
	[][]Artist{{
		{"spotify:artist:1XpDYCrUJnvCo9Ez6yeMWh", "Tenacious D"},
		{"spotify:artist:5sgprVkYi5OjM4nxKI8ZWg", "Tenacious"},
		{"spotify:artist:2Sf0QliiNtuNTJe51TgalE", "Young Tenacious"},
		{"spotify:artist:6snWJ93BNH3JIbLGWczD1D", "BO, TENACIOUS BREED"},
		{"spotify:artist:7mtIirvrKV5SUE90cPeUnR", "Tenacious Da Terrist"},
	},
		{
			{"spotify:artist:4FUej2oub0ZMSAfSpVpc4H", "M.T.T.S. (Ty Bru, Medic, Tenacious)"},
		},
		[]Artist(nil),
	}, errEOF,
}

var searchAlbumFixt = struct {
	res [][]Album
	err error
}{
	[][]Album{
		{{
			"spotify:album:4LJbsUCNTcNNNHNiX6qES1", "POD",
			[]Artist{
				{"spotify:artist:1XpDYCrUJnvCo9Ez6yeMWh", "Tenacious D"},
			},
		},
			{
				"spotify:album:33LXyaRjDrMZILnvp1umPU", "Tenacious",
				[]Artist{
					{"spotify:artist:1XpDYCrUJnvCo9Ez6yeMUU", "Tenacious DR"},
				},
			},
			{
				"spotify:album:7mv1ciCld5Bp1y6TDGtjQY", "Tenacious D",
				[]Artist{
					{"spotify:artist:1XpDYCrUJnvCo9Ez6yeMWh", "Tenacious D"},
				},
			},
			{
				"spotify:album:0zPvqiP3ZmCyYgXdupvdBi", "Tenacious D",
				[]Artist{
					{"spotify:artist:1XpDYCrUJnvCo9Ez6yeMUU", "Tenacious DR"},
				},
			},
			{
				"spotify:album:6PjFFuDv6tnIlwyT33ugdj", "Best In Da State, Vol. 1",
				[]Artist{
					{"spotify:artist:1XpDYCrUJnvCo9Ez6yeMWh", "Tenacious D"},
				},
			},
		},
		{
			{
				"spotify:album:4LJbsUCNTcNNNHNiX6qES1", "POD",
				[]Artist{
					{"spotify:artist:1XpDYCrUJnvCo9Ez6yeMUU", "Tenacious DR"},
				},
			},
		},
		[]Album(nil),
	}, errEOF,
}

var searchTrackFixt = struct {
	res [][]Track
	err error
}{
	[][]Track{
		{
			{
				"spotify:track:6crBy2sODw2HS53xquM6us", "Tribute",
				"spotify:album:1AckkxSo39144vOBrJ1GkS", "Tenacious D",
				[]Artist{
					{"spotify:artist:1XpDYCrUJnvCo9Ez6yeMWh", "Tenacious D"},
				},
			},
			{
				"spotify:track:3USOTQdPtZlgVurwaqsdI5", "Fuck Her Gently",
				"spotify:album:1AckkxSo39144vOBrJ1GkS", "Tenacious D",
				[]Artist{
					{"spotify:artist:1XpDYCrUJnvCo9Ez6yeMWh", "Tenacious D"},
				},
			},
		},
		{
			{"spotify:track:7jxSwduTONaRn00SQQbInK", "White and Nerdy (Karaoke Version)",
				"spotify:album:4zb12kro15pJa2YG1uEcmR", "Drew's Famous #1 Karaoke Hits: Sing Like Tenacious D, Flight of the Conchords, & Friends",
				[]Artist{
					{"spotify:artist:0GRelLAVHzwasg0Ja7gJUy", "The Karaoke Crew"},
				},
			},
		},
		[]Track(nil),
	}, errEOF,
}
