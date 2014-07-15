sscc [![GoDoc](https://godoc.org/github.com/pblaszczyk/sscc?status.png)](https://godoc.org/github.com/pblaszczyk/sscc) [![Build Status](https://travis-ci.org/pblaszczyk/sscc.svg?branch=master)](https://travis-ci.org/pblaszczyk/sscc)
========

sscc is a set of tools used by cmd/ssc in order to control Spotify desktop
application on Linux and use Web API in order to search for artists/albums/tracks.


## cmd/sscc [![GoDoc](https://godoc.org/github.com/pblaszczyk/sscc/cmd/sscc?status.png)](https://godoc.org/github.com/pblaszczyk/sscc/cmd/sscc)

The cmd/sscc is a commdn-line tool for controlling Spotify desktop application.

#### Installation

In order to install sscc application in your environment, Go compiler is required.
If environment is configured, please run:

```
~ $ go get -u github.com/pblaszczyk/sscc/cmd/sscc
```

#### Usage

```
NAME:
   sscc - commandline controller of Spotify desktop app.

USAGE:
   sscc [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   run		Run Spotify desktop app.
   kill		Kill Spotify desktop app.
   next		Play next track.
   prev		Play prev track..
   open		Play music identified by uri.
   seek		Seek.
   play		Play current track/uri/pos.
   stop		Stop.
   toggle	Play/Pause.
   search	Search for artist/album/track.
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --version, -v	print the version
   --help, -h		show help
```
