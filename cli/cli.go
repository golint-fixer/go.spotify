package cli

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/codegangsta/cli"
	"github.com/pblaszczyk/sscc"
)

// App is a structure controlling sscc executable.
type App struct {
	*cli.App
	ctrl sscc.Controller
}

// NewApp returns initialized instance of ssc struct.
func NewApp() (app *App) {
	ctrl := sscc.NewControl(&sscc.Context{})
	app = &App{cli.NewApp(), ctrl}
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
		{Name: "seek", Usage: "Goto.", Action: app.Goto},
		{Name: "play", Usage: "Play current track/uri/pos.", Action: app.Play},
		{Name: "stop", Usage: "Stop.", Action: app.Stop},
		{Name: "toggle", Usage: "Play/Pause.", Action: app.Toggle},
		{Name: "status", Usage: "Status.", Action: app.Status},
		{Name: "track", Usage: "Current track.", Action: app.CurTrack},
		{Name: "setpos", Usage: "Sets position.", Action: app.SetPos},
		{Name: "length", Usage: "Length of current track.", Action: app.Length},
		{Name: "pos", Usage: "Current position.", Action: app.Pos},
		{Name: "canplay", Usage: "Can play.", Action: app.CanPlay},
		{Name: "cannext", Usage: "Can next.", Action: app.CanNext},
		{Name: "canprev", Usage: "Can prev.", Action: app.CanPrev},
		{Name: "canctrl", Usage: "Can control.", Action: app.CanControl},
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
	switch {
	case sscc.IsEOF(err) || err == nil:
		return
	default:
		fmt.Printf("sscc: %q\n", err)
		os.Exit(1)
	}
}

// Start starts spotify app.
func (a *App) Start(ctx *cli.Context) {
	handleErr(a.ctrl.Run())
}

// Raise raises spotify app's window.
func (a *App) Raise(ctx *cli.Context) {
	handleErr(a.ctrl.Raise())
}

// Kill stops spotify app.
func (a *App) Kill(ctx *cli.Context) {
	handleErr(a.ctrl.Kill())
}

// Next starts playing next track.
func (a *App) Next(ctx *cli.Context) {
	handleErr(a.ctrl.Next())
}

// Prev starts playing prev track.
func (a *App) Prev(ctx *cli.Context) {
	handleErr(a.ctrl.Prev())
}

// Open starts playing specified uri.
func (a *App) Open(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	handleErr(a.ctrl.Open(sscc.URI(ctx.Args().First())))
}

// Play starts playing.
func (a *App) Play(ctx *cli.Context) {
	handleErr(a.ctrl.Play())
}

// Goto pos.
func (a *App) Goto(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	d, err := time.ParseDuration(ctx.Args().First())
	handleErr(err)
	handleErr(a.ctrl.Goto(d))
}

// Stop playing current track.
func (a *App) Stop(ctx *cli.Context) {
	handleErr(a.ctrl.Stop())
}

// Toggle plays/pauses current track.
func (a *App) Toggle(ctx *cli.Context) {
	handleErr(a.ctrl.Toggle())
}

// CurTrack displays info about current track.
func (a *App) CurTrack(ctx *cli.Context) {
	track, err := a.ctrl.CurTrack()
	handleErr(err)
	fmt.Println(track)
}

// SetPos moves to requested pos.
func (a *App) SetPos(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	d, err := time.ParseDuration(ctx.Args().First())
	handleErr(err)
	handleErr(a.ctrl.SetPos(d))
}

// Length display current track's length.
func (a *App) Length(ctx *cli.Context) {
	l, err := a.ctrl.Length()
	handleErr(err)
	fmt.Println(l)
}

// Pos returns current position.
func (a *App) Pos(ctx *cli.Context) {
	d, err := a.ctrl.Pos()
	handleErr(err)
	fmt.Println(d)
}

// CanPlay returns info if playing is possible.
func (a *App) CanPlay(ctx *cli.Context) {
	b, err := a.ctrl.CanPlay()
	handleErr(err)
	fmt.Println(b)
}

// CanNext returns info if going to next track is possible.
func (a *App) CanNext(ctx *cli.Context) {
	b, err := a.ctrl.CanNext()
	handleErr(err)
	fmt.Println(b)
}

// CanPrev returns info if going to prev track is possible.
func (a *App) CanPrev(ctx *cli.Context) {
	b, err := a.ctrl.CanPrev()
	handleErr(err)
	fmt.Println(b)
}

// CanControl returns info if controlling is possible.
func (a *App) CanControl(ctx *cli.Context) {
	b, err := a.ctrl.CanControl()
	handleErr(err)
	fmt.Println(b)
}

// interactive runs in limited interactive mode if configured.
func (a *App) interactive(ctx *cli.Context) {
	if ctx.Bool("i") {
		fmt.Print("Play: ")
		r := bufio.NewReader(os.Stdin)
		uri, _, err := r.ReadLine()
		handleErr(err)
		handleErr(a.ctrl.Open(sscc.URI(uri)))
	}
}

// Artist searches for artist.
func (a *App) Artist(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	res, err := make(chan []sscc.Artist), make(chan error)
	a.ctrl.SearchArtist(ctx.Args().First(), res, err)
LOOP:
	for {
		select {
		case res := <-res:
			disp(res)
		case err := <-err:
			handleErr(err)
			break LOOP
		}
	}
	fmt.Println("")
	a.interactive(ctx)
}

// Album searches for album.
func (a *App) Album(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	res, err := make(chan []sscc.Album), make(chan error, 1)
	a.ctrl.SearchAlbum(ctx.Args().First(), res, err)
LOOP:
	for {
		select {
		case res := <-res:
			disp(res)
		case err := <-err:
			handleErr(err)
			break LOOP
		}
	}
	fmt.Println("")
	a.interactive(ctx)
}

// Track searches for track.
func (a *App) Track(ctx *cli.Context) {
	handleErr(validateSingle(ctx.Args()))
	res, err := make(chan []sscc.Track), make(chan error)
	a.ctrl.SearchTrack(ctx.Args().First(), res, err)
LOOP:
	for {
		select {
		case res := <-res:
			disp(res)
		case err := <-err:
			handleErr(err)
			break LOOP
		}
	}
	a.interactive(ctx)
}

// Status returns current status.
func (a *App) Status(ctx *cli.Context) {
	status, err := a.ctrl.Status()
	handleErr(err)
	fmt.Println(status)
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
		for j, l := 0, reflect.ValueOf(r).Index(i).NumField(); j < l; j++ {
			f := reflect.ValueOf(r).Index(i).Field(j)
			if f.Kind() == reflect.Slice {
				fmt.Printf("%q\n", reflect.ValueOf(r).Index(i).Type().Field(j).Name)
				disp(f.Interface())
			} else {
				fmt.Printf("%q: %q",
					reflect.ValueOf(r).Index(i).Type().Field(j).Name, f.String())
			}
			if j < l-1 {
				fmt.Println("")
			}
		}
		fmt.Println("")
	}
}
