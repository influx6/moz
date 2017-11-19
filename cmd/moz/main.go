package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/metrics/custom"
	"github.com/influx6/gobuild/build"
	"github.com/influx6/moz"
	"github.com/influx6/moz/ast"
	"github.com/minio/cli"

	annons "github.com/influx6/moz/annotations"
)

var (
	events      = metrics.New(custom.StackDisplay(os.Stdout))
	annotations = moz.CopyAnnotationsTo(ast.NewAnnotationRegistryWith(events))
)

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
			Name:        "generate-file",
			Action:      generateFileCLI,
			Description: "Runs the moz parser to parse and generate code for all annotations found in the file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "f,fromFile",
					Value: "",
					Usage: "-f=./",
				},
				cli.BoolFlag{
					Name:  "fw,forceWrite",
					Usage: "-fw=true",
				},
				cli.StringFlag{
					Name:  "t,toDir",
					Value: "",
					Usage: "-t=./",
				},
			},
		},
		{
			Name:        "generate-tag",
			Action:      generatePackageCLIWithTag,
			Description: "Runs the moz parser to parse and generate code for all annotations in the package directory where files have the provided build tag",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "tg,tag",
					Value: "",
					Usage: "-tg=maze",
				},
				cli.StringFlag{
					Name:  "f,fromDir",
					Value: "",
					Usage: "-f=./",
				},
				cli.BoolFlag{
					Name:  "fw,forceWrite",
					Usage: "-fw=true",
				},
				cli.StringFlag{
					Name:  "t,toDir",
					Value: "",
					Usage: "-t=./",
				},
			},
		},
		{
			Name:        "generate",
			Action:      generatePackageCLI,
			Description: "Runs the moz parser to parse and generate code for all annotations in the package directory",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "f,fromDir",
					Value: "",
					Usage: "-f=./",
				},
				cli.BoolFlag{
					Name:  "fw,forceWrite",
					Usage: "-fw=true",
				},
				cli.StringFlag{
					Name:  "t,toDir",
					Value: "",
					Usage: "-t=./",
				},
			},
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

	rootDir := c.String("root")
	if rootDir == "" {
		rootDir = "./"
	}

	inDir := "files"

	directives, err := annons.AssetsAnnotationGenerator(rootDir,
		ast.AnnotationDeclaration{Arguments: []string{pkgName, strings.Join(extensions, ":"), inDir}},
		ast.PackageDeclaration{FilePath: rootDir}, ast.Package{})

	if err != nil {
		events.Emit(metrics.Error(err), metrics.With("dir", rootDir), metrics.With("message", "Failed to generation assets"))
		return
	}

	for _, directive := range directives {
		if directive.Dir != "" {
			coDir := filepath.Join(rootDir, directive.Dir)

			if _, err := os.Stat(coDir); err != nil {
				fmt.Printf("- Creating package directory: %q\n", coDir)
				if err := os.MkdirAll(coDir, 0700); err != nil && err != os.ErrExist {
					events.Emit(metrics.Error(err), metrics.With("dir", coDir), metrics.With("message", "Failed to create directory"))
					return
				}
			}

		}

		if directive.Writer == nil {
			continue
		}

		coFile := filepath.Join(rootDir, directive.Dir, directive.FileName)

		if _, err := os.Stat(coFile); err == nil {
			if directive.DontOverride {
				continue
			}
		}

		dFile, err := os.Create(coFile)
		if err != nil {
			events.Emit(metrics.Error(err), metrics.With("dir", rootDir), metrics.With("file", coFile), metrics.With("message", "Failed to create file"))
			return
		}

		if _, err := directive.Writer.WriteTo(dFile); err != nil {
			events.Emit(metrics.Error(err), metrics.With("dir", rootDir), metrics.With("file", coFile), metrics.With("message", "Failed to write to file"))
			return
		}

		rel, _ := filepath.Rel(rootDir, coFile)
		fmt.Printf("- Add file to package directory: %q\n", rel)

		dFile.Close()
	}
}

func generateFileCLI(c *cli.Context) {
	var err error

	forceWrite := c.Bool("forceWrite")
	fromFile := c.String("fromFile")
	if fromFile == "" {
		err = fmt.Errorf("file target not provided, use the -fromfile flag")
		events.Emit(metrics.Error(err), metrics.With("dir", fromFile), metrics.With("message", "Failed to retrieve current directory"))
		return
	}

	toDir := c.String("toDir")
	if filepath.IsAbs(toDir) {
		err = fmt.Errorf("-toDir flag can not be a absolute path but a relative path to the directory")
		events.Emit(metrics.Error(err), metrics.With("dir", fromFile), metrics.With("toDir", toDir), metrics.With("message", "Failed to retrieve current directory"))
		return
	}

	// If its not an absolute path then get real absolute
	if !filepath.IsAbs(fromFile) {
		pwd, perr := os.Getwd()
		if perr != nil {
			events.Emit(metrics.Error(perr), metrics.With("file", fromFile), metrics.With("toDir", toDir), metrics.With("message", "Failed to retrieve current fileectory"))
			return
		}

		fromFile = filepath.Join(pwd, fromFile)
	}

	events.Emit(metrics.Info("Using FromFile: %s", fromFile), metrics.With("file", fromFile))
	events.Emit(metrics.Info("Using ToDir: %s", toDir), metrics.With("Dir", toDir))

	pkg, err := ast.ParseFileAnnotations(events, fromFile)
	if err != nil {
		events.Emit(metrics.Error(err), metrics.With("file", fromFile), metrics.With("toDir", toDir), metrics.With("message", "Failed to parse package annotations"))
		return
	}

	events.Emit(metrics.Info("Begin Annotation Execution"), metrics.With("toDir", toDir), metrics.With("fromFile", fromFile))

	if err := moz.ParseWith(toDir, events, annotations, forceWrite, pkg); err != nil {
		events.Emit(metrics.Error(err), metrics.With("file", fromFile), metrics.With("toDir", toDir), metrics.With("message", "Failed to parse package declarations"))
	}

	events.Emit(metrics.Info("Finished"), metrics.With("toDir", toDir), metrics.With("fromFile", fromFile))
}

func generatePackageCLI(c *cli.Context) {
	var err error

	forceWrite := c.Bool("forceWrite")
	fromDir := c.String("fromDir")
	toDir := c.String("toDir")

	if filepath.IsAbs(toDir) {
		err = fmt.Errorf("-toDir flag can not be a absolute path but a relative path to the directory")
		events.Emit(metrics.Error(err), metrics.With("dir", fromDir), metrics.With("toDir", toDir), metrics.With("message", "Failed to retrieve current directory"))
		return
	}

	if fromDir == "" {
		fromDir, err = os.Getwd()
		if err != nil {
			events.Emit(metrics.Error(err), metrics.With("dir", fromDir), metrics.With("toDir", toDir), metrics.With("message", "Failed to retrieve current directory"))
			return
		}
	}

	events.Emit(metrics.Info("Using FromDir: %s", fromDir), metrics.With("dir", fromDir))
	events.Emit(metrics.Info("Using ToDir: %s", toDir), metrics.With("dir", toDir))

	pkg, err := ast.ParseAnnotations(events, fromDir)
	if err != nil {
		events.Emit(metrics.Error(err), metrics.With("dir", fromDir), metrics.With("toDir", toDir), metrics.With("message", "Failed to parse package annotations"))
		return
	}

	events.Emit(metrics.Info("Begin Annotation Execution"), metrics.With("toDir", toDir), metrics.With("fromDir", fromDir), metrics.With("Packages", len(pkg)))

	if err := moz.ParseWith(toDir, events, annotations, forceWrite, pkg...); err != nil {
		events.Emit(metrics.Error(err), metrics.With("dir", fromDir), metrics.With("toDir", toDir), metrics.With("message", "Failed to parse package declarations"))
	}

	events.Emit(metrics.Info("Finished"), metrics.With("toDir", toDir), metrics.With("fromDir", fromDir))
}

func generatePackageCLIWithTag(c *cli.Context) {
	var err error

	forceWrite := c.Bool("forceWrite")
	fromDir := c.String("fromDir")
	toDir := c.String("toDir")
	buildTag := c.String("tag")

	if toDir == "" {
		toDir = "."
	}

	if buildTag == "" {
		buildTag = c.Args().First()
	}

	if filepath.IsAbs(toDir) {
		err = fmt.Errorf("-toDir flag can not be a absolute path but a relative path to the directory")
		events.Emit(metrics.Error(err), metrics.With("dir", fromDir), metrics.With("toDir", toDir), metrics.With("message", "Failed to retrieve current directory"))
		return
	}

	if fromDir == "" {
		fromDir, err = os.Getwd()
		if err != nil {
			events.Emit(metrics.Error(err), metrics.With("dir", fromDir), metrics.With("toDir", toDir), metrics.With("message", "Failed to retrieve current directory"))
			return
		}
	}

	events.Emit(metrics.Info("Using FromDir: %s", fromDir), metrics.With("dir", fromDir))
	events.Emit(metrics.Info("Using ToDir: %s", toDir), metrics.With("dir", toDir))
	events.Emit(metrics.Info("Using Build Tag"), metrics.With("tag", buildTag))

	ctx := build.Default
	ctx.BuildTags = append(ctx.BuildTags, buildTag)
	ctx.RequiredTags = append(ctx.RequiredTags, buildTag)

	pkg, err := ast.FilteredPackageWithBuildCtx(events, fromDir, ctx)
	if err != nil {
		events.Emit(metrics.Error(err), metrics.With("dir", fromDir), metrics.With("toDir", toDir), metrics.With("message", "Failed to parse package annotations"))
		return
	}

	events.Emit(metrics.Info("Begin Annotation Execution"), metrics.With("toDir", toDir), metrics.With("fromDir", fromDir), metrics.With("Packages", len(pkg)))

	if err := moz.ParseWith(toDir, events, annotations, forceWrite, pkg...); err != nil {
		events.Emit(metrics.Error(err), metrics.With("dir", fromDir), metrics.With("toDir", toDir), metrics.With("message", "Failed to parse package declarations"))
	}

	events.Emit(metrics.Info("Finished"), metrics.With("toDir", toDir), metrics.With("fromDir", fromDir))
}
