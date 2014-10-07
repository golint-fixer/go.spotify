package sscc

// Artist is a model for holding results of searching for artist.
type Artist struct {
	URI  string
	Name string
}

// Album is a model for holding results of searching for album.
type Album struct {
	URI     string
	Name    string
	Artists []Artist
}

// Track is a model for holding results of searching for track.
type Track struct {
	URI       string
	Name      string
	AlbumURI  string
	AlbumName string
	Artists   []Artist
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
