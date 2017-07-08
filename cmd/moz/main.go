package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/influx6/faux/fmtwriter"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/metrics/sentries/stdout"
	"github.com/influx6/moz"
	"github.com/influx6/moz/ast"
	"github.com/influx6/moz/cmd/moz/templates"
	"github.com/influx6/moz/gen"
	"github.com/influx6/moz/utils"
	"github.com/minio/cli"

	_ "github.com/influx6/moz/annotations"
)

var events = metrics.New(stdout.Stdout{})

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
			Flags:  []cli.Flag{},
		},
		{
			Name:        "generate",
			Action:      annotationCLI,
			Description: "Runs the moz parser to parse and generate code for all annotations",
			Flags:       []cli.Flag{},
		},
		{
			Name:        "assets",
			Action:      assetsCLI,
			Description: "Generates a package to build files into go source",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "n,name",
					Value: "assets",
					Usage: "-n=assets",
				},
				cli.StringFlag{
					Name:  "e,exts",
					Value: "",
					Usage: "-e='.tml, .ball, .css'",
				},
				cli.StringFlag{
					Name:  "r,root",
					Value: "./",
					Usage: "-r=/tmp/bob",
				},
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		return nil
	}

	app.RunAndExitOnError()
}

func assetsCLI(c *cli.Context) {
	extensions := strings.Split(c.String("exts"), ",")

	pkgName := c.String("name")

	rootCMD := c.String("root")
	if rootCMD == "" {
		rootCMD = "./"
	}

	inDir := "files"
	rootDir := filepath.Join(rootCMD, pkgName)
	assetDir := filepath.Join(rootDir, inDir)

	genFile := gen.Package(
		gen.Name(pkgName),
		gen.Text("//go:generate go run generate.go"),
	)

	mainFile := gen.Block(
		gen.Commentary(
			gen.Text("+build ignore"),
		),
		gen.Text("\n"),
		gen.Package(
			gen.Name("main"),
			gen.Imports(
				gen.Import("fmt", ""),
				gen.Import("path/filepath", ""),
				gen.Import("github.com/influx6/moz/gen", ""),
				gen.Import("github.com/influx6/moz/utils", ""),
				gen.Import("github.com/influx6/faux/vfiles", ""),
				gen.Import("github.com/influx6/faux/fmtwriter", ""),
				gen.Import("github.com/influx6/faux/metrics", ""),
				gen.Import("github.com/influx6/faux/metrics/sentries/stdout", ""),
			),
			gen.Function(
				gen.Name("main"),
				gen.Constructor(),
				gen.Returns(),
				gen.Block(
					gen.SourceText(
						templates.Must("main.tml"),
						struct {
							Extensions       []string
							TargetDir        string
							Package          string
							GenerateTemplate string
						}{
							TargetDir:  inDir,
							Extensions: extensions,
							Package:    pkgName,
							GenerateTemplate: `{{range $key, $value := .Files}}
								files[{{quote $key}}] = []byte("{{$value}}")
							{{end}}`,
						},
					),
				),
			),
		),
	)

	if err := os.MkdirAll(assetDir, 0700); err != nil && !os.IsExist(err) {
		events.Emit(stdout.Error(err).With("dir", rootCMD).
			With("targetDir", rootDir).
			With("message", "Failed to create new package directory"))
		panic(err)
	}

	genDir := filepath.Join(rootDir, pkgName+".go")
	if err := utils.WriteFile(events, fmtwriter.New(genFile, true, true), genDir); err != nil {
		events.Emit(stdout.Error(err).With("dir", rootCMD).
			With("targetDir", rootDir).
			With("message", "Failed to create new package directory: generate.go"))
		panic(err)
	}

	dir := filepath.Join(rootDir, "generate.go")
	if err := utils.WriteFile(events, fmtwriter.New(mainFile, true, true), dir); err != nil {
		events.Emit(stdout.Error(err).With("dir", rootCMD).
			With("targetDir", rootDir).
			With("message", "Failed to create new package directory: generate.go"))
		panic(err)
	}
}

func annotationCLI(c *cli.Context) {
	cdir, err := os.Getwd()
	if err != nil {
		events.Emit(stdout.Error(err).With("dir", cdir).With("message", "Failed to retrieve current directory"))
		return
	}

	events.Emit(stdout.Info("Using Dir: %s", cdir).With("dir", cdir))

	pkgs, err := ast.ParseAnnotations(events, cdir)
	if err != nil {
		events.Emit(stdout.Error(err).With("dir", cdir).With("message", "Failed to parse package annotations"))
		return
	}

	if err := moz.Parse(events, pkgs...); err != nil {
		events.Emit(stdout.Error(err).With("dir", cdir).With("message", "Failed to parse package declarations"))
	}
}
