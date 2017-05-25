package main

import (
	"os"

	"github.com/influx6/faux/sink"
	"github.com/influx6/faux/sink/sinks"
	"github.com/influx6/moz/ast"
)

var events = sink.New(sinks.Stdout{})

func main() {
	cdir, err := os.Getwd()
	if err != nil {
		events.Emit(sinks.Error("Failed to retrieve current directory: %+q", err))
		return
	}

	events.Emit(sinks.Info("Using Dir: %s", cdir).With("dir", cdir))

	_, _ = ast.ParseAnnotations(cdir)
}
