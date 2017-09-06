package filesystem_test

import (
	"bytes"
	"testing"

	"github.com/influx6/faux/tests"
	"github.com/influx6/moz/gen/filesystem"
)

func TestFileSystemGetFile(t *testing.T) {
	memFS := filesystem.FileSystem(
		filesystem.Meta("version", "1.0"),
		filesystem.Description("FileSystem for the faas docker system"),
		filesystem.Dir(
			"app",
			filesystem.File(
				"readme.md",
				filesystem.Content("# App v1.0"),
			),
			filesystem.Dir(
				"src",
				filesystem.File(
					"main.go",
					filesystem.Content(`package main`),
				),
				filesystem.Dir(
					"dest",
					filesystem.File(
						"main.js",
						filesystem.Content(`strict;`),
					),
				),
			),
		),
		filesystem.File(
			"dockerfile",
			filesystem.Content(`
        FROM alpine:latest

        RUN sudo apt-get update
        RUN sudo apt-get install git golang

        CMD ["/bin/bomb"]
      `),
		),
	)

	_, dirErr := memFS.GetDir("app/src")
	if dirErr != nil {
		tests.Failed("Should have successfully retrieved dir %q but got %+q", "app/src", dirErr)
	}
	tests.Passed("Should have successfully retrieved dir %q", "app/src")

	_, fileErr := memFS.GetFile("app/src/dest/main.js")
	if dirErr != nil {
		tests.Failed("Should have successfully retrieved file %q but got %+q", "app/src/dest/main.js", fileErr)
	}
	tests.Passed("Should have successfully retrieved file %q", "app/src/dest/main.js")

	_, fakeErr := memFS.GetFile("app/src/dest/bust.js")
	if fakeErr == nil {
		tests.Failed("Should have successfully failed to retrieve file %q but got %+q", "app/src/dest/bust.js", fakeErr)
	}
	tests.Passed("Should have successfully failed to retrieve file %q", "app/src/dest/bust.js")
}

func TestFileSystem(t *testing.T) {
	expected := []byte("{\n\n\t\".meta\": \"{\\\"description\\\":\\\"FileSystem for the faas docker system\\\",\\\"version\\\":\\\"1.0\\\"}\",\n\n\t\"app/readme.md\": \"# App v1.0\",\n\n\t\"app/src/dest/main.js\": \"strict;\",\n\n\t\"app/src/main.go\": \"package main\",\n\n\t\"dockerfile\": \"\\n        FROM alpine:latest\\n\\n        RUN sudo apt-get update\\n        RUN sudo apt-get install git golang\\n\\n        CMD [\\\"/bin/bomb\\\"]\\n      \"\n\n}\n")
	memFS := filesystem.FileSystem(
		filesystem.Meta("version", "1.0"),
		filesystem.Description("FileSystem for the faas docker system"),
		filesystem.Dir(
			"app",
			filesystem.File(
				"readme.md",
				filesystem.Content("# App v1.0"),
			),
			filesystem.Dir(
				"src",
				filesystem.File(
					"main.go",
					filesystem.Content(`package main`),
				),
				filesystem.Dir(
					"dest",
					filesystem.File(
						"main.js",
						filesystem.Content(`strict;`),
					),
				),
			),
		),
		filesystem.File(
			"dockerfile",
			filesystem.Content(`
        FROM alpine:latest

        RUN sudo apt-get update
        RUN sudo apt-get install git golang

        CMD ["/bin/bomb"]
      `),
		),
	)

	var content bytes.Buffer

	jsonMemFS := filesystem.JSONFS(memFS, true)
	written, err := jsonMemFS.WriteTo(&content)
	if err != nil {
		tests.Failed("Should have successfully written file system into json: %+q.", err)
	}
	tests.Passed("Should have successfully written file system into json.")

	if written == 0 {
		tests.Failed("Should have received written total above 0")
	}
	tests.Passed("Should have received written total above 0")

	if !bytes.Equal(expected, content.Bytes()) {
		tests.Info("Expected: %+q\n", expected)
		tests.Info("Recieved: %+q\n", content.Bytes())
		tests.Failed("Should have received expected result")
	}
	tests.Passed("Should have received expected result")
}
