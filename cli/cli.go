package cli

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/pblaszczyk/sscc/desktop"
	"github.com/pblaszczyk/sscc/webapi"

	"github.com/codegangsta/cli"
)

type sscc struct {
	*cli.App
}

// NewApp returns initialized instance of ssc struct.
func NewApp() (app *sscc) {
	app = &sscc{cli.NewApp()}
	app.App.Name = "sscc"
	app.App.Version = "0.0.1"
	app.App.Usage = "commandline controller of Spotify desktop app."
	app.App.Commands = []cli.Command{
		{Name: "run", Usage: "Run Spotify desktop app.", Action: app.Start},
		{Name: "kill", Usage: "Kill Spotify desktop app.", Action: app.Kill},
		{Name: "next", Usage: "Play next track.", Action: app.Next},
		{Name: "prev", Usage: "Play prev track..", Action: app.Prev},
		{Name: "open", Usage: "Play music identified by uri.", Action: app.Open},
		{Name: "seek", Usage: "Seek.", Action: app.Seek},
		{Name: "play", Usage: "Play current track/uri/pos.", Action: app.Play},
		{Name: "stop", Usage: "Stop.", Action: app.Stop},
		{Name: "toggle", Usage: "Play/Pause.", Action: app.Toggle},
		{Name: "search", Usage: "Search for artist/album/track.",
			Subcommands: []cli.Command{
				{Name: "artist", Usage: "Search for artist.", Action: app.Artist},
				{Name: "album", Usage: "Search for album.", Action: app.Album},
				{Name: "track", Usage: "Search for track.", Action: app.Track},
			}},
	}
	app.App.EnableBashCompletion = true
	webapi.Bar = true
	return
}

var handleErr = func(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Start starts spotify app.
func (s *sscc) Start(ctx *cli.Context) {
	handleErr(desktop.Start())
}

// Kill stops spotify app.
func (s *sscc) Kill(ctx *cli.Context) {
	handleErr(desktop.Kill())
}

// Next starts playing next track.
func (s *sscc) Next(ctx *cli.Context) {
	handleErr(desktop.Next())
}

// Prev starts playing prev track.
func (s *sscc) Prev(ctx *cli.Context) {
	handleErr(desktop.Prev())
}

// Open starts playing specified uri.
func (s *sscc) Open(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	handleErr(desktop.OpenUri(ctx.Args().First()))
}

// Open starts playing.
func (s *sscc) Play(ctx *cli.Context) {
	handleErr(desktop.Play())
}

// Seek pos.
func (s *sscc) Seek(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	n, err := strconv.ParseInt(ctx.Args().First(), 10, 64)
	handleErr(err)
	handleErr(desktop.Seek(n))
}

// Stop playing current track.
func (s *sscc) Stop(ctx *cli.Context) {
	handleErr(desktop.Stop())
}

// Play/Pause current track.
func (s *sscc) Toggle(ctx *cli.Context) {
	handleErr(desktop.PlayPause())
}

// Search for artist.
func (s *sscc) Artist(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	r, err := webapi.SearchArtist(ctx.Args().First())
	handleErr(err)
	disp(r)
}

// Search for album.
func (s *sscc) Album(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	r, err := webapi.SearchAlbum(ctx.Args().First())
	handleErr(err)
	disp(r)
}

// Search for track.
func (s *sscc) Track(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	r, err := webapi.SearchTrack(ctx.Args().First())
	handleErr(err)
	disp(r)
}

func validateSingle(args cli.Args) error {
	n := len(args)
	if n == 0 {
		return fmt.Errorf("Please specify valid argument.")
	}
	if n > 1 {
		return fmt.Errorf("Invalid number of arguments.")
	}
	return nil
}

func disp(r interface{}) {
	for i := reflect.ValueOf(r).Len() - 1; i >= 0; i-- {
		rec := false
		for j := 0; j < reflect.ValueOf(r).Index(i).NumField(); j++ {
			f := reflect.ValueOf(r).Index(i).Field(j)
			if f.Kind() == reflect.Slice {
				rec = true
				fmt.Printf("\n%q\n",
					reflect.ValueOf(r).Index(i).Type().Field(j).Name)
				disp(f.Interface())
			} else {
				fmt.Printf("%q: %q ",
					reflect.ValueOf(r).Index(i).Type().Field(j).Name, f.String())
			}
		}
		if !rec {
			fmt.Printf("\n\n\n")
		}
	}
}
