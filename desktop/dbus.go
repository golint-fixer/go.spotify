package desktop

import (
	"fmt"

	idbus "github.com/guelfey/go.dbus"
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
	methodOpenUri      = "org.mpris.MediaPlayer2.Player.OpenUri"
	methodQuit         = "org.mpris.MediaPlayer2.Quit"
	methodRaise        = "org.mpris.MediaPlayer2.Raise"
	propTrackList      = "org.mpris.MediaPlayer2.HasTrackList"
	propIdentity       = "org.mpris.MediaPlayer2.Identity"
	propDesktopEntry   = "org.mpris.MediaPlayer2.DesktopEntry"
	propSupUriSchemes  = "org.mpris.MediaPlayer2.SupportedUriSchemes"
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
	methodMachineId    = "org.freedesktop.DBus.Peer.GetMachineId"
)

var initDbus = func() (*idbus.Object, error) {
	c, err := idbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("dbus: failed to init dbus session: %q", err.Error())
	}
	return c.Object(dest, objPath), nil
}

// Play next track from currently played playlist.
func Next() error {
	return noArgsMethod(methodNext)
}

// Play previous track from currently played playlist.
func Prev() error {
	return noArgsMethod(methodPrev)
}

// Pause currently played track.
func Pause() error {
	return noArgsMethod(methodPause)
}

// Toggle between play & pause state.
func PlayPause() error {
	return noArgsMethod(methodPlayPause)
}

// Stop playing current track.
func Stop() error {
	return noArgsMethod(methodStop)
}

// Play currently set track.
func Play() error {
	return noArgsMethod(methodPlay)
}

// Seek for offset µs. Negative offset means seeking back.
func Seek(offset int64) error {
	o, err := initDbus()
	if err != nil {
		return err
	}
	call := o.Call(methodSeek, 0, offset)
	return call.Err
}

// Go to specified position in currently set track.
// trackId must be equal to currently set track.
func SetPos(trackId string, pos int64) error {
	o, err := initDbus()
	if err != nil {
		return err
	}
	call := o.Call(methodSetPos, 0, idbus.ObjectPath(trackId), pos)
	return call.Err
}

// Start playing track with provided uri.
func OpenUri(uri string) error {
	o, err := initDbus()
	if err != nil {
		return err
	}
	call := o.Call(methodOpenUri, 0, uri)
	return call.Err
}

// Quit spotify process.
func Quit() error {
	return noArgsMethod(methodQuit)
}

// Raise spotify player.
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
