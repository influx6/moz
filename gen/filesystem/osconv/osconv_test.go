package osconv_test

import (
	"bytes"
	"testing"

	"github.com/influx6/faux/tests"
	"github.com/influx6/moz/gen/filesystem"
	"github.com/influx6/moz/gen/filesystem/osconv"
)

func TestOSConv(t *testing.T) {
	fs, err := osconv.ConvertDir("./fixtures", false, func(rel string) bool {
		return true
	})
	if err != nil {
		tests.Failed("Should have successsfully run through dir: %+q", err)
	}
	tests.Passed("Should have successsfully run through dir")

	var content bytes.Buffer
	defer content.Reset()

	jsonMemFS := filesystem.JSONFS(fs, true)
	written, err := jsonMemFS.WriteTo(&content)
	if err != nil {
		tests.Failed("Should have successfully written file system into json: %+q.", err)
	}
	tests.Passed("Should have successfully written file system into json.")

	if written == 0 {
		tests.Failed("Should have received written total above 0")
	}
	tests.Passed("Should have received written total above 0")
}
