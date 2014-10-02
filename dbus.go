package sscc

import "fmt"
import "github.com/guelfey/go.dbus"

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

var initDbus = func() (*dbus.Object, error) {
	c, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("dbus: failed to init dbus session: %q", err.Error())
	}
	return c.Object(dest, objPath), nil
}

// Next plays next track from currently played playlist.
func Next() error {
	return noArgsMethod(methodNext)
}

// Prev plays previous track from currently played playlist.
func Prev() error {
	return noArgsMethod(methodPrev)
}

// Pause pauses currently played track.
func Pause() error {
	return noArgsMethod(methodPause)
}

// PlayPause toggles between play & pause state.
func PlayPause() error {
	return noArgsMethod(methodPlayPause)
}

// Stop stops playing current track.
func Stop() error {
	return noArgsMethod(methodStop)
}

// Play plays currently set track.
func Play() error {
	return noArgsMethod(methodPlay)
}

// Seek seeks for offset µs. Negative offset means seeking back.
func Seek(offset int64) error {
	o, err := initDbus()
	if err != nil {
		return err
	}
	call := o.Call(methodSeek, 0, offset)
	return call.Err
}

// SetPos goes to specified position in currently set track.
// trackID must be equal to currently set track.
func SetPos(trackID string, pos int64) error {
	o, err := initDbus()
	if err != nil {
		return err
	}
	call := o.Call(methodSetPos, 0, dbus.ObjectPath(trackID), pos)
	return call.Err
}

// OpenURI starts playing track with provided uri.
func OpenURI(uri string) error {
	o, err := initDbus()
	if err != nil {
		return err
	}
	call := o.Call(methodOpenURI, 0, uri)
	return call.Err
}

// Quit quits spotify process.
func Quit() error {
	return noArgsMethod(methodQuit)
}

// Raise raises spotify player.
func Raise() error {
	return noArgsMethod(methodRaise)
}

func noArgsMethod(method string) error {
	o, err := initDbus()
	if err != nil {
		return err
	}
	call := o.Call(method, 0)
	return call.Err
}
