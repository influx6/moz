package utils

import (
	"bytes"
	"context"
	"io"

	"github.com/influx6/process"
)

// GoFmt returns a new io.WriterTo which contains a go formated source output.
func GoFmt(w io.WriterTo) (io.WriterTo, error) {
	cmd := process.Command{
		Name:  "gofmt",
		Level: process.RedAlert,
	}

	var input, output bytes.Buffer

	if _, err := w.WriteTo(&input); err != nil && err != io.EOF {
		return nil, err
	}

	if err := cmd.Run(context.Background(), &output, nil, &input); err != nil {
		return &input, nil
	}

	return &output, nil
}
