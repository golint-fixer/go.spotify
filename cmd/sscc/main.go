/*NAME:
   sscc - commandline controller of Spotify desktop app.

USAGE:
   sscc [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   run          Run Spotify desktop app.
   kill         Kill Spotify desktop app.
   raise        Raise Spotify desktop app.
   next         Play next track.
   prev         Play prev track.
   open         Play music identified by uri.
   seek         Seek.
   play         Play current track/uri/pos.
   stop         Stop.
   toggle       Play/Pause.
   search       Search for artist/album/track.
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h           show help
   --version, -v        print the version
*/
package main

import (
	"os"

	"github.com/pblaszczyk/sscc/cli"
)

func main() {
	cli.NewApp().Run(os.Args)
}
