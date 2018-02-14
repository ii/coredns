# Installation

CoreDNS is written in Go

## Downloading

## Docker

## Source

To compile CoreDNS we assume you have a working Go setup, see various tutorials (/links) if you
don't have that already configured. The Go version that comes with your OS is probably too old to
compile CoreDNS as we require Go 1.9.x at the moment (Feb 2018).

With CoreDNS we try to vendor all our dependencies, but because of various reasons (mostly making it
possible for external plugins to compile), we can not vendor all our dependencies. Hence to compile
CoreDNS, you still need to `go get`. The `Makefile` we include handles all of these steps. So if you
downloaded the source of CoreDNS (possible with a `go get`), you'll need to perform the following
steps to compile CoreDNS:

## External Plugins

You'll need to have the source of CoreDNS installed for this to work.

