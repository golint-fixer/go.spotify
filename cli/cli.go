// Package cli provides command line interface for ssc application
// including: parsing arguments, handling specified arguments and printint out
// obtained results.
package cli

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/pblaszczyk/sscc"
)

type impl struct {
	*cli.App
}

// NewApp returns initialized instance of ssc struct.
func NewApp() (app *impl) {
	app = &impl{cli.NewApp()}
	app.App.Name = "sscc"
	app.App.Version = "0.0.1"
	app.App.Usage = "commandline controller of Spotify desktop app."
	app.App.Commands = []cli.Command{
		{Name: "run", Usage: "Run Spotify desktop app.", Action: app.Start},
		{Name: "kill", Usage: "Kill Spotify desktop app.", Action: app.Kill},
		{Name: "raise", Usage: "Raise Spotify desktop app.", Action: app.Raise},
		{Name: "next", Usage: "Play next track.", Action: app.Next},
		{Name: "prev", Usage: "Play prev track.", Action: app.Prev},
		{Name: "open", Usage: "Play music identified by uri.", Action: app.Open},
		{Name: "seek", Usage: "Seek.", Action: app.Seek},
		{Name: "play", Usage: "Play current track/uri/pos.", Action: app.Play},
		{Name: "stop", Usage: "Stop.", Action: app.Stop},
		{Name: "toggle", Usage: "Play/Pause.", Action: app.Toggle},
		{Name: "search", Usage: "Search for artist/album/track.",
			Subcommands: []cli.Command{
				{Name: "artist", Usage: "Search for artist.", Action: app.Artist,
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "i", Usage: "Quasi-interactive mode"}},
				},
				{Name: "album", Usage: "Search for album.", Action: app.Album,
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "i", Usage: "Quasi-interactive mode"}},
				},
				{Name: "track", Usage: "Search for track.", Action: app.Track,
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "i", Usage: "Quasi-interactive mode"}},
				},
			},
		},
	}
	return
}

var handleErr = func(err error) {
	if err != nil {
		fmt.Printf("sscc: %q\n", err)
		os.Exit(1)
	}
}

// Start starts spotify app.
func (s *impl) Start(ctx *cli.Context) {
	handleErr(sscc.Start())
}

// Raise raises spotify app's window.
func (s *impl) Raise(ctx *cli.Context) {
	handleErr(sscc.Raise())
}

// Kill stops spotify app.
func (s *impl) Kill(ctx *cli.Context) {
	handleErr(sscc.Kill())
}

// Next starts playing next track.
func (s *impl) Next(ctx *cli.Context) {
	handleErr(sscc.Next())
}

// Prev starts playing prev track.
func (s *impl) Prev(ctx *cli.Context) {
	handleErr(sscc.Prev())
}

// Open starts playing specified uri.
func (s *impl) Open(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	handleErr(sscc.OpenURI(ctx.Args().First()))
}

// Open starts playing.
func (s *impl) Play(ctx *cli.Context) {
	handleErr(sscc.Play())
}

// Seek pos.
func (s *impl) Seek(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	n, err := strconv.ParseInt(ctx.Args().First(), 10, 64)
	handleErr(err)
	handleErr(sscc.Seek(n))
}

// Stop playing current track.
func (s *impl) Stop(ctx *cli.Context) {
	handleErr(sscc.Stop())
}

// Play/Pause current track.
func (s *impl) Toggle(ctx *cli.Context) {
	handleErr(sscc.PlayPause())
}

// interactive runs in limited interactive mode if configured.
func interactive(ctx *cli.Context) {
	if ctx.Bool("i") {
		fmt.Print("Play: ")
		r := bufio.NewReader(os.Stdin)
		uri, _, err := r.ReadLine()
		handleErr(err)
		handleErr(sscc.OpenURI(string(uri)))
	}
}

// Search for artist.
func (s *impl) Artist(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	r, err := sscc.SearchArtist(ctx.Args().First())
	handleErr(err)
	disp(r)
	interactive(ctx)
}

// Search for album.
func (s *impl) Album(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	r, err := sscc.SearchAlbum(ctx.Args().First())
	handleErr(err)
	disp(r)
	interactive(ctx)
}

// Search for track.
func (s *impl) Track(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	r, err := sscc.SearchTrack(ctx.Args().First())
	handleErr(err)
	disp(r)
	interactive(ctx)
}

func validateSingle(args cli.Args) error {
	n := len(args)
	if n == 0 {
		return fmt.Errorf("please specify valid argument")
	}
	if n > 1 {
		return fmt.Errorf("invalid number of arguments")
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
