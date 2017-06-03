// +build ignore
//
//

package main

import (
	"fmt"

	"path/filepath"

	"github.com/influx6/moz/gen"

	"github.com/influx6/moz/utils"

	"github.com/influx6/faux/vfiles"

	"github.com/influx6/faux/fmtwriter"

	"github.com/influx6/faux/sink"

	"github.com/influx6/faux/sink/sinks"
)

func main() {
	// go:generate go run generate.go

	events := sink.New(sinks.Stdout{})

	items, err := vfiles.ParseDir("files", []string{})
	if err != nil {
		panic(fmt.Sprintf("Failed to walk directory properly: %+q", err))
	}

	assetGen := gen.Package(
		gen.Name("templates"),
		gen.AssignVar(
			gen.Name("files"),
			gen.Type("make(map[string][]byte)"),
		),
		gen.Text(`
    // Must attempts to retrieve the file data if found else panics.
    func Must(file string) []byte {
        data, err := Get(file)
        if err != nil {
            panic(err)
        }

        return data
    }

    // Get retrieves the giving file data from the map store if it exists.
    func Get(file string) ([]byte, error){
        data, ok := files[file]
        if !ok {
            return nil, fmt.Errorf("File data for %q not found", file)
        }

        return data, nil
    }
    `),
		gen.Function(
			gen.Name("init"),
			gen.Constructor(),
			gen.Returns(),
			gen.Block(
				gen.SourceText(`
                {{range $key, $value := .Files}}
								files[{{quote $key}}] = []byte("{{$value}}")
							{{end}}
            `, struct {
					Files map[string]string
				}{
					Files: items,
				}),
			),
		),
	)

	dir := filepath.Join(".", "templates.go")
	if err := utils.WriteFile(events, fmtwriter.New(assetGen, true), dir); err != nil {
		events.Emit(sinks.Error(err).With("dir", dir).
			With("message", "Failed to create new package file: templates.go"))
		panic(err)
	}
}
