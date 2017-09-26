package utils

import (
	"io"
	"os"
	"path/filepath"

	"github.com/influx6/faux/metrics"
)

// WriteFile copies the data from the writer into the desired path, creating
// any directory as needed.
func WriteFile(events metrics.Metrics, writer io.WriterTo, toPath string) error {
	dirPath := filepath.Dir(toPath)

	if err := os.MkdirAll(dirPath, 0700); err != nil && !os.IsExist(err) {
		events.Emit(metrics.Error(err).With("dir", dirPath).With("targetPath", toPath).
			With("message", "Failed to create new package directory: generate.go"))
		return err
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		events.Emit(metrics.Error(err).With("targetPath", toPath).With("dir", dirPath).
			With("message", "Failed to create new source file"))
		return err
	}

	defer toFile.Close()

	if _, err = writer.WriteTo(toFile); err != nil {
		events.Emit(metrics.Error(err).With("targetPath", toPath).With("dir", dirPath).
			With("message", "Failed to write new source file"))
		return err
	}

	return nil
}
