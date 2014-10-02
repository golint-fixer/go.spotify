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
