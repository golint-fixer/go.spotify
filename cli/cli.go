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
		{Name: "open", Usage: "Play music identified by uri.", Action: nil},
		{Name: "play", Usage: "Play current track/uri/pos.",
			Flags: []cli.Flag{cli.StringFlag{Name: "uri", Usage: "Play uri."},
				cli.StringFlag{Name: "pos", Usage: "Seek pos."}},
			Action: app.Play},
		{Name: "stop", Usage: "Stop.", Action: app.Stop},
		{Name: "toggle", Usage: "Play/Pause.", Action: app.Toggle},
		{Name: "search", Usage: "Search for artist/album/track.",
			Flags: []cli.Flag{cli.StringFlag{Name: "type", Usage: "artist/album/track"},
				cli.StringFlag{Name: "value", Usage: "Name of artist/album/track."}},
			Action: app.Search},
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

func (s *sscc) Start(ctx *cli.Context) {
	handleErr(desktop.Start())
}

func (s *sscc) Kill(ctx *cli.Context) {
	handleErr(desktop.Kill())
}

func (s *sscc) Next(ctx *cli.Context) {
	handleErr(desktop.Next())
}

func (s *sscc) Prev(ctx *cli.Context) {
	handleErr(desktop.Prev())
}

func (s *sscc) Play(ctx *cli.Context) {
	uri := ctx.String("uri")
	if uri == "" {
		handleErr(desktop.Play())
	} else {
		handleErr(desktop.OpenUri(uri))
	}
	posS := ctx.String("pos")
	if posS == "" {
		return
	}
	pos, err := strconv.ParseInt(posS, 10, 64)
	handleErr(err)
	if pos != 0 {
		handleErr(desktop.Seek(pos))
	}
}

func (s *sscc) Stop(ctx *cli.Context) {
	handleErr(desktop.Stop())
}

func (s *sscc) Toggle(ctx *cli.Context) {
	handleErr(desktop.PlayPause())
}

func (s *sscc) Search(ctx *cli.Context) {
	typ, val := ctx.String("type"), ctx.String("value")
	handleErr(search(typ, val))
}

func search(typ, val string) error {
	if val == "" {
		return fmt.Errorf("spotifycli: searched value not specified")
	}
	switch typ {
	case "artist":
		r, err := webapi.SearchArtist(val)
		if err != nil {
			return err
		}
		disp(r)
	case "album":
		r, err := webapi.SearchAlbum(val)
		if err != nil {
			return err
		}
		disp(r)
	case "track":
		r, err := webapi.SearchTrack(val)
		if err != nil {
			return err
		}
		disp(r)
	default:
		return fmt.Errorf("spotifycli: invalid search command")
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
