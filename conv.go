package spotify

// conv converts data structures from one format to another in a following way:
// - artistResp -> []Artist
// - albumResp  -> []Album
// - trackResp  -> []Track
// If different type is provided as argument, function panics.
func conv(d interface{}) interface{} {
	switch d := d.(type) {
	case *artistResp:
		var res []Artist
		for i := range d.Artists.Items {
			res = append(res, Artist{
				URI: d.Artists.Items[i].URI, Name: d.Artists.Items[i].Name,
			})
		}
		return res
	case *albumResp:
		var res []Album
		for i := range d.Albums.Items {
			res = append(res, Album{
				URI: d.Albums.Items[i].URI, Name: d.Albums.Items[i].Name,
			})
		}
		return res
	case *trackResp:
		var res []Track
		for i := range d.Tracks.Items {
			var a []Artist
			for j := range d.Tracks.Items[i].Artists {
				a = append(a, Artist{
					URI:  d.Tracks.Items[i].Artists[j].URI,
					Name: d.Tracks.Items[i].Artists[j].Name,
				})
			}
			res = append(res, Track{
				URI: d.Tracks.Items[i].URI, Name: d.Tracks.Items[i].Name,
				AlbumURI:  d.Tracks.Items[i].Album.URI,
				AlbumName: d.Tracks.Items[i].Album.Name, Artists: a,
			})
		}
		return res
	default:
		panic("sscc: unsupported data format")
	}
}
