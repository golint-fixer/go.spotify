package sscc

// URI is a type representing Spotify URI.
type URI string

// Artist is a model for artist's data.
type Artist struct {
	URI  string // URI is a Spotify URI of the artist.
	Name string // Name of the artist.
}

// Album is a model for album's data.
type Album struct {
	URI     string   // URI is a Spotify URI of the album.
	Name    string   // Name is the name of the album.
	Artists []Artist // Artists is a list of artists of the album.
}

// Track is a model for track's data.
type Track struct {
	URI       string   // URI is a Spotify URI of the track.
	Name      string   // Name is the name of the track.
	AlbumURI  string   // AlbumURI is a URI of album containing track.
	AlbumName string   // AlbumName is the name of album containing track.
	Artists   []Artist // Artists is a list of artists of the track.
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
