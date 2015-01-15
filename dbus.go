package sscc

import (
	"fmt"
	"sync"
	"time"

	dbs "github.com/guelfey/go.dbus"
)

// Dbuser describe operations on Spotify's dbus interface.
type Dbuser interface {
	Open(URI) error                 // Open starts playing track with `URI`.
	Play() error                    // Play plays currently active track.
	Stop() error                    // Stop stops playing current track.
	Pause() error                   // Pause pauses currently played track.
	Toggle() error                  // Toggle toggles between play & pause state.
	Next() error                    // Next plays next track.
	Prev() error                    // Prev plays previous track.
	Goto(time.Duration) error       // Goto seeks for offset µs.
	SetPos(time.Duration) error     // SetPos goes to specified positionk.
	Length() (time.Duration, error) // Length returns length of current track.
	Raise() error                   // Raise raises Spotify app.
	Quit() error                    // Quit quits Spotify app.
	CurTrack() (Track, error)       // CurTrack returns currently played track.
	Status() (Status, error)        // Status returns current status of an app.
	Pos() (time.Duration, error)    // Pos returns current position.
	CanPlay() (bool, error)         // CanPlay returns info if you can play.
	CanNext() (bool, error)         // CanNext checks if next is available.
	CanPrev() (bool, error)         // CanPrev checks if prev is available.
	CanControl() (bool, error)      // CanControl checks if control is available.
}

// dbus is a default implementation of `Dbuser`.
type dbus struct {
	sync.Mutex
	o *dbs.Object // o is a dbus control object.
}

// newDbuser returns a default implementation of `Dbuser`.
func newDbuser() Dbuser {
	return &dbus{}
}

// init initializes dbus if not yet done.
func (d *dbus) init() error {
	d.Lock()
	defer d.Unlock()
	if d.o != nil {
		return nil
	}
	c, err := dbs.SessionBus()
	if err != nil {
		return fmt.Errorf("sscc: failed to init dbus session: %q", err)
	}
	d.o = c.Object(dest, objPath)
	return nil
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
func (d *dbus) Goto(offset time.Duration) error {
	if err := d.init(); err != nil {
		return err
	}
	return d.o.Call(methodSeek, 0, offset.Nanoseconds()*1000).Err
}

// SetPos implements `Dbuser`.
func (d *dbus) SetPos(pos time.Duration) error {
	track, err := d.CurTrack()
	if err != nil {
		return err
	}
	return d.o.Call(methodSetPos, 0, dbs.ObjectPath(string(track.URI)),
		pos.Nanoseconds()*1000).Err
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

// CurTrack implements `Dbuser`.
func (d *dbus) CurTrack() (track Track, err error) {
	if err = d.init(); err != nil {
		return
	}
	v, err := d.o.GetProperty(propMetadata)
	if err != nil {
		return
	}
	m, ok := v.Value().(map[string]dbs.Variant)
	defer func() {
		if !ok {
			err = fmt.Errorf(invDbusResp, v.Value())
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

// Status implements `Dbuser`.
func (d *dbus) Status() (Status, error) {
	if err := d.init(); err != nil {
		return Status(""), err
	}
	v, err := d.o.GetProperty(propPlaybackStatus)
	if err != nil {
		return Status(""), err
	}
	status, ok := v.Value().(string)
	if !ok {
		return Status(""), fmt.Errorf(invDbusResp, v.Value())
	}
	return makeStatus(status)
}

// Length implements `Dbuser`.
func (d *dbus) Length() (l time.Duration, err error) {
	if err = d.init(); err != nil {
		return
	}
	v, err := d.o.GetProperty(propMetadata)
	if err != nil {
		return
	}
	m, ok := v.Value().(map[string]dbs.Variant)
	defer func() {
		if !ok {
			err = fmt.Errorf(invDbusResp, v.Value())
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

// Pos implements `Dbuser`.
func (d *dbus) Pos() (time.Duration, error) {
	if err := d.init(); err != nil {
		return 0, err
	}
	v, err := d.o.GetProperty(propPos)
	if err != nil {
		return 0, err
	}
	pos, ok := v.Value().(int64)
	if !ok {
		return 0, fmt.Errorf(invDbusResp, v.Value())
	}
	return time.Duration(pos * 1000), nil
}

// CanPlay implements `Dbuser`.
func (d *dbus) CanPlay() (bool, error) {
	return d.boolOpt(propCanPlay)
}

// boolOpt is a helper func retrieving value of boolean property.
func (d *dbus) boolOpt(prop string) (bool, error) {
	if err := d.init(); err != nil {
		return false, err
	}
	v, err := d.o.GetProperty(prop)
	if err != nil {
		return false, err
	}
	res, ok := v.Value().(bool)
	if !ok {
		return false, fmt.Errorf(invDbusResp, v.Value())
	}
	return res, nil
}

// CanNext implements `Dbuser`.
func (d *dbus) CanNext() (bool, error) {
	return d.boolOpt(propCanGoNext)
}

// CanPrev implements `Dbuser`.
func (d *dbus) CanPrev() (bool, error) {
	return d.boolOpt(propCanGoPrev)
}

// CanControl implements `Dbuser`.
func (d *dbus) CanControl() (bool, error) {
	return d.boolOpt(propCanControl)
}

// noArgsMethod is a helper function initializing `*dbuser` and calling
// a provided method.
func (d *dbus) noArgsMethod(method string) error {
	if err := d.init(); err != nil {
		return err
	}
	return d.o.Call(method, 0).Err
}

// dbusErr is an error returned for invalid dbus response.
type dbusErr struct {
	msg string
}

// Error implements `error`.
func (e *dbusErr) Error() string {
	return e.msg
}

// newDbusError returns instance of `*dbusErr`.
func newDbusError(format string, a ...interface{}) error {
	return &dbusErr{fmt.Sprintf(format, a...)}
}

// IsDbus returns a boolean indicating whether the error is known to report
// that dbus operation failed.
func IsDbus(err error) (ok bool) {
	_, ok = err.(*dbusErr)
	return
}

// invDbusResp is a format of an error message for an invalid dbus response.
const invDbusResp = "sscc: invalid dbus response: %v"

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
