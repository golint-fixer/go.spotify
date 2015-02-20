// +build linux

package spotify

import (
	"sync"
	"time"

	dbs "github.com/guelfey/go.dbus"
)

// Dbus is a structure implementing Dbus logic controlling Spotify
// desktop application.
type Dbus struct {
	sync.Mutex
	o *dbs.Object // o is a dbus control object.
}

// NewDbus returns a new instance of Dbus.
func NewDbus() (*Dbus, error) {
	d := &Dbus{}
	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Dbus) init() error {
	d.Lock()
	defer d.Unlock()
	if d.o != nil {
		return nil
	}
	c, err := dbs.SessionBus()
	if err != nil {
		return errorf("failed to init dbus session: %q", err)
	}
	d.o = c.Object(dest, objPath)
	return nil
}

// Next plays next track.
func (d *Dbus) Next() error {
	return d.noArgsMethod(methodNext)
}

// Prev plays previous track.
func (d *Dbus) Prev() error {
	return d.noArgsMethod(methodPrev)
}

// Pause pauses currently played track.
func (d *Dbus) Pause() error {
	return d.noArgsMethod(methodPause)
}

// Toggle toggles between play & pause state.
func (d *Dbus) Toggle() error {
	return d.noArgsMethod(methodPlayPause)
}

// Stop stops playing current track.
func (d *Dbus) Stop() error {
	return d.noArgsMethod(methodStop)
}

// Play plays currently active track.
func (d *Dbus) Play() error {
	return d.noArgsMethod(methodPlay)
}

// Goto seeks for offset µs.
func (d *Dbus) Goto(offset time.Duration) error {
	return d.o.Call(methodSeek, 0, offset.Nanoseconds()*1000).Err
}

// SetPos goes to specified positionk.
func (d *Dbus) SetPos(pos time.Duration) error {
	track, err := d.Track()
	if err != nil {
		return err
	}
	return d.o.Call(methodSetPos, 0, dbs.ObjectPath(string(track.URI)),
		pos.Nanoseconds()*1000).Err
}

// Open starts playing track with URI.
func (d *Dbus) Open(uri URI) error {
	return d.o.Call(methodOpenURI, 0, string(uri)).Err
}

// Quit quits Spotify app.
func (d *Dbus) Quit() error {
	return d.noArgsMethod(methodQuit)
}

// Raise raises Spotify app.
func (d *Dbus) Raise() error {
	return d.noArgsMethod(methodRaise)
}

// Track returns currently played track.
func (d *Dbus) Track() (track Track, err error) {
	v, err := d.o.GetProperty(propMetadata)
	if err != nil {
		return
	}
	m, ok := v.Value().(map[string]dbs.Variant)
	defer func() {
		if !ok {
			err = errorf(invDbusResp, v.Value())
			return
		}
	}()
	if !ok {
		return
	}
	if track.Name, ok = m["xesam:title"].Value().(string); !ok {
		return
	}
	if track.URI, ok = m["xesam:url"].Value().(string); !ok {
		return
	}
	if track.AlbumName, ok = m["xesam:album"].Value().(string); !ok {
		return
	}
	artists, ok := m["xesam:artist"].Value().([]string)
	if ok = ok && len(artists) > 0; !ok {
		return
	}
	track.Artists = append(track.Artists, Artist{Name: artists[0]})
	return
}

// Status returns current status of an app.
func (d *Dbus) Status() (Status, error) {
	v, err := d.o.GetProperty(propPlaybackStatus)
	if err != nil {
		return Status(""), err
	}
	status, ok := v.Value().(string)
	if !ok {
		return Status(""), errorf(invDbusResp, v.Value())
	}
	return makeStatus(status)
}

// Length returns length of current track.
func (d *Dbus) Length() (l time.Duration, err error) {
	v, err := d.o.GetProperty(propMetadata)
	if err != nil {
		return
	}
	m, ok := v.Value().(map[string]dbs.Variant)
	defer func() {
		if !ok {
			err = errorf(invDbusResp, v.Value())
			return
		}
	}()
	if !ok {
		return
	}
	length, ok := m["mpris:length"].Value().(uint64)
	if !ok {
		return
	}
	l = time.Duration(length * 1000)
	return
}

// Pos returns current position.
func (d *Dbus) Pos() (time.Duration, error) {
	v, err := d.o.GetProperty(propPos)
	if err != nil {
		return 0, err
	}
	pos, ok := v.Value().(int64)
	if !ok {
		return 0, errorf(invDbusResp, v.Value())
	}
	return time.Duration(pos * 1000), nil
}

// CanPlay returns info if you can play.
func (d *Dbus) CanPlay() (bool, error) {
	return d.boolOpt(propCanPlay)
}

// boolOpt is a helper func retrieving value of boolean property.
func (d *Dbus) boolOpt(prop string) (bool, error) {
	v, err := d.o.GetProperty(prop)
	if err != nil {
		return false, err
	}
	res, ok := v.Value().(bool)
	if !ok {
		return false, errorf(invDbusResp, v.Value())
	}
	return res, nil
}

// CanNext checks if next is available.
func (d *Dbus) CanNext() (bool, error) {
	return d.boolOpt(propCanGoNext)
}

// CanPrev checks if prev is available.
func (d *Dbus) CanPrev() (bool, error) {
	return d.boolOpt(propCanGoPrev)
}

// CanControl checks if control is available.
func (d *Dbus) CanControl() (bool, error) {
	return d.boolOpt(propCanControl)
}

// noArgsMethod is a helper function initializing `*dbuser` and calling
// a provided method.
func (d *Dbus) noArgsMethod(method string) error {
	return d.o.Call(method, 0).Err
}

// invDbusResp is a format of an error message for an invalid dbus response.
const invDbusResp = "invalid dbus response: %v"

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
