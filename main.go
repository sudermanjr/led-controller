package main

import (
	"github.com/markbates/pkger"
	"github.com/sudermanjr/led-controller/cmd"
)

var (
	// version is set during build
	version = "development"
	// commit is set during build
	commit = "n/a"
)

func main() {
	pkger.Include("/pkg/dashboard/assets")
	pkger.Include("/pkg/dashboard/templates")
	pkger.Include("/pkg/screen/gifs")

	cmd.Execute(version, commit)
}
