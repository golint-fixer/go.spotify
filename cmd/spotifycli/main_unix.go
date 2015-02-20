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

func platform() {
	switch os.Args[1] {
	case "next":
		handlerr(newDbus().Next())
	case "prev":
		handlerr(newDbus().Prev())
	case "open":
		if len(os.Args) != 3 {
			usage()
		}
		handlerr(newDbus().Open(spotify.URI(os.Args[2])))
	case "play":
		handlerr(newDbus().Play())
	case "stop":
		handlerr(newDbus().Stop())
	case "toggle":
		handlerr(newDbus().Toggle())
	case "status":
		status, err := newDbus().Status()
		handlerr(err)
		fmt.Println(status)
	case "track":
		track, err := newDbus().Track()
		handlerr(err)
		fmt.Println(track)
	case "length":
		length, err := newDbus().Length()
		handlerr(err)
		fmt.Println(length)
	case "raise":
		handlerr(newDbus().Raise())
	default:
		usage()
	}
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
