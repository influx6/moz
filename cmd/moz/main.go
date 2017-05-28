package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/fatih/color"
	"github.com/influx6/faux/sink"
	"github.com/influx6/faux/sink/sinks"
	"github.com/influx6/moz"
	"github.com/influx6/moz/ast"
	"github.com/minio/cli"
)

var events = sink.New(sinks.Stdout{})

// Version defines the version number for the cli.
var Version = "0.1"

var helpTemplate = `NAME:
{{.Name}} - {{.Usage}}

DESCRIPTION:
{{.Description}}

USAGE:
{{.Name}} {{if .Flags}}[flags] {{end}}command{{if .Flags}}{{end}} [arguments...]

COMMANDS:
	{{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
	{{end}}{{if .Flags}}
FLAGS:
	{{range .Flags}}{{.}}
	{{end}}{{end}}
VERSION:
` + Version +
	`{{ "\n"}}`

// Cmd defines a struct for defining a command.
type Cmd struct {
	*cli.App
}

// VersionAction defines the action called when seeking the Version detail.
func VersionAction(c *cli.Context) {
	fmt.Println(color.BlueString(fmt.Sprintf("moz version %s %s/%s", Version, runtime.GOOS, runtime.GOARCH)))
}

func main() {
	app := cli.NewApp()
	app.Name = "moz"
	app.Author = ""
	app.Usage = "moz {{command}}"
	app.Flags = []cli.Flag{}
	app.Description = "moz: CLI tooling for the go language generator."
	app.CustomAppHelpTemplate = helpTemplate

	app.Commands = []cli.Command{
		{
			Name:   "version",
			Action: VersionAction,
		},
		{
			Name:   "annotation",
			Action: annotationCLI,
		},
	}

	app.Before = func(c *cli.Context) error {
		return nil
	}

	app.RunAndExitOnError()
}

func annotationCLI(c *cli.Context) {
	cdir, err := os.Getwd()
	if err != nil {
		events.Emit(sinks.Error(err).With("dir", cdir).With("message", "Failed to retrieve current directory"))
		return
	}

	events.Emit(sinks.Info("Using Dir: %s", cdir).With("dir", cdir))

	pkgs, err := ast.ParseAnnotations(cdir)
	if err != nil {
		events.Emit(sinks.Error(err).With("dir", cdir).With("message", "Failed to parse package annotations"))
		return
	}

	if err := moz.Parse(events, pkgs...); err != nil {
		events.Emit(sinks.Error(err).With("dir", cdir).With("message", "Failed to parse package declarations"))
	}
}
