package sscc

import (
	"fmt"

	dbs "github.com/guelfey/go.dbus"
)

const (
	dest               = "org.mpris.MediaPlayer2.spotify"
	objPath            = "/org/mpris/MediaPlayer2"
	methodNext         = "org.mpris.MediaPlayer2.Player.Next"
	methodPrev         = "org.mpris.MediaPlayer2.Player.Previous"
	methodPause        = "org.mpris.MediaPlayer2.Player.Pause"
	methodPlayPause    = "org.mpris.MediaPlayer2.Player.PlayPause"
	methodStop         = "org.mpris.MediaPlayer2.Player.Stop"
	methodPlay         = "org.mpris.MediaPlayer2.Player.Play"
	methodSeek         = "org.mpris.MediaPlayer2.Player.Seek"
	methodSetPos       = "org.mpris.MediaPlayer2.Player.SetPosition"
	methodOpenURI      = "org.mpris.MediaPlayer2.Player.OpenUri"
	methodQuit         = "org.mpris.MediaPlayer2.Quit"
	methodRaise        = "org.mpris.MediaPlayer2.Raise"
	propTrackList      = "org.mpris.MediaPlayer2.HasTrackList"
	propIdentity       = "org.mpris.MediaPlayer2.Identity"
	propDesktopEntry   = "org.mpris.MediaPlayer2.DesktopEntry"
	propSupURISchemes  = "org.mpris.MediaPlayer2.SupportedUriSchemes"
	propSupMimeTypes   = "org.mpris.MediaPlayer2.SupportedMimeTypes"
	propPlaybackStatus = "org.mpris.MediaPlayer2.Player.PlaybackStatus"
	propLoopStatus     = "org.mpris.MediaPlayer2.Player.LoopStatus"
	propRate           = "org.mpris.MediaPlayer2.Player.Rate"
	propShuffle        = "org.mpris.MediaPlayer2.Player.Shuffle"
	propMetadata       = "org.mpris.MediaPlayer2.Player.Metadata"
	propVolume         = "org.mpris.MediaPlayer2.Player.Volume"
	propPos            = "org.mpris.MediaPlayer2.Player.Position"
	propMinRate        = "org.mpris.MediaPlayer2.Player.MinimumRate"
	propMaxRate        = "org.mpris.MediaPlayer2.Player.MaximumRate"
	propCanGoNext      = "org.mpris.MediaPlayer2.Player.CanGoNext"
	propCanGoPrev      = "org.mpris.MediaPlayer2.Player.CanGoPrevious"
	propCanPlay        = "org.mpris.MediaPlayer2.Player.CanPlay"
	propCanPause       = "org.mpris.MediaPlayer2.Player.CanPause"
	propCanSeek        = "org.mpris.MediaPlayer2.Player.CanSeek"
	propCanControl     = "org.mpris.MediaPlayer2.Player.CanControl"
	propSeeked         = "org.mpris.MediaPlayer2.Player.Seeked"
	methodGet          = "org.freedesktop.DBus.Properties.Get"
	methodSet          = "org.freedesktop.DBus.Properties.Set"
	methodGetAll       = "org.freedesktop.DBus.Properties.GetAll"
	methodIntrospect   = "org.freedesktop.DBus.Introspectable.Introspect"
	methodPing         = "org.freedesktop.DBus.Peer.Ping"
	methodMachineID    = "org.freedesktop.DBus.Peer.GetMachineId"
)

// Dbuser is an interface for operations on Spotify desktop application's dbus
// interface.
type Dbuser interface {
	Open(URI) error          // Open starts playing track with provided URI.
	Play() error             // Play plays currently set track.
	Stop() error             // Stop stops playing current track.
	Pause() error            // Pause pauses currently played track.
	Toggle() error           // Toggle toggles between play & pause state.
	Next() error             // Next plays next track.
	Prev() error             // Prev plays previous track.
	Goto(int64) error        // Goto seeks for offset µs.
	SetPos(URI, int64) error // SetPos goes to specified position for track.
	Raise() error            // Raise raises Spotify app.
	Quit() error             // Quit quits Spotify app.
	Track() (Track, error)   // Track returns currently played track.
}

type dbus struct {
	o      *dbs.Object
	initer dbusInit
}

func newDbus() *dbus {
	return &dbus{initer: initDbus}
}

var defaultDbus = newDbus()

func init() {
	if err := defaultDbus.init(); err != nil {
		panic(err)
	}
}

type dbusInit func() (*dbs.Object, error)

func initDbus() (*dbs.Object, error) {
	c, err := dbs.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("sscc: failed to init dbus session: %q", err.Error())
	}
	return c.Object(dest, objPath), nil
}

func (d *dbus) init() (err error) {
	if d.o != nil {
		return
	}
	d.o, err = d.initer()
	return
}

// Next implements `Dbuser`.
func (d *dbus) Next() error {
	return d.noArgsMethod(methodNext)
}

// Prev implements `Dbuser`.
func (d *dbus) Prev() error {
	return d.noArgsMethod(methodPrev)
}

// Pause implements `Dbuser`.
func (d *dbus) Pause() error {
	return d.noArgsMethod(methodPause)
}

// Toggle implements `Dbuser`.
func (d *dbus) Toggle() error {
	return d.noArgsMethod(methodPlayPause)
}

// Stop implements `Dbuser`.
func (d *dbus) Stop() error {
	return d.noArgsMethod(methodStop)
}

// Play implements `Dbuser`.
func (d *dbus) Play() error {
	return d.noArgsMethod(methodPlay)
}

// Goto implements `Dbuser`.
func (d *dbus) Goto(offset int64) error {
	if err := d.init(); err != nil {
		return err
	}
	return d.o.Call(methodSeek, 0, offset).Err
}

// SetPos implements `Dbuser`.
func (d *dbus) SetPos(trackID URI, pos int64) error {
	if err := d.init(); err != nil {
		return err
	}
	return d.o.Call(methodSetPos, 0, dbs.ObjectPath(string(trackID)), pos).Err
}

// Open implements `Dbuser`.
func (d *dbus) Open(uri URI) error {
	if err := d.init(); err != nil {
		return err
	}
	return d.o.Call(methodOpenURI, 0, string(uri)).Err
}

// Quit implements `Dbuser`.
func (d *dbus) Quit() error {
	return d.noArgsMethod(methodQuit)
}

// Raise implements `Dbuser`.
func (d *dbus) Raise() error {
	return d.noArgsMethod(methodRaise)
}

// Track implements `Dbuser`.
func (d *dbus) Track() (Track, error) {
	panic("sscc: not implemented")
}

func (d *dbus) noArgsMethod(method string) error {
	if err := d.init(); err != nil {
		return err
	}
	return d.o.Call(method, 0).Err
}
