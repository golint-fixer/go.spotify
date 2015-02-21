// +build linux

package main

import (
	"fmt"
	"os"

	"github.com/pblaszczyk/go.spotify"
)

func newDbus() *spotify.Dbus {
	d, err := spotify.NewDbus()
	handlerr(err)
	return d
}

func open() {
	if len(os.Args) != 3 {
		usage()
	}
	handlerr(newDbus().Open(spotify.URI(os.Args[2])))
}

func length() {
	length, err := newDbus().Length()
	handlerr(err)
	fmt.Println(length)
}

func status() {
	status, err := newDbus().Status()
	handlerr(err)
	fmt.Println(status)
}

func track() {
	track, err := newDbus().Track()
	handlerr(err)
	fmt.Println(track)
}

func next() {
	handlerr(newDbus().Next())
}

func prev() {
	handlerr(newDbus().Prev())
}

func play() {
	handlerr(newDbus().Play())
}

func stop() {
	handlerr(newDbus().Stop())
}

func toggle() {
	handlerr(newDbus().Toggle())
}

func raise() {
	handlerr(newDbus().Raise())
}

func platform() {
	if f, ok := cmd2func[os.Args[1]]; ok {
		f()
	} else {
		usage()
	}
}

var cmd2func = map[string]func(){
	"next":   next,
	"prev":   prev,
	"open":   open,
	"play":   play,
	"stop":   stop,
	"toggle": toggle,
	"status": status,
	"track":  track,
	"length": length,
	"raise":  raise,
}

func platfusage() {
	fmt.Printf(`  play <URI>- Play URI.
  stop      - Stop playing.
  toggle    - Toggle playing.
  status    - Current Status.
  track     - Current track.
  length    - Length of a current track.
  raise     - Raise the Spotify desktop app.
`)
}
