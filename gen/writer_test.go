package gen_test

import (
	"bytes"
	"testing"

	"github.com/influx6/faux/tests"
	"github.com/influx6/moz/gen"
)

func TestFromReader(t *testing.T) {
	var in, out bytes.Buffer
	in.WriteString("Awating the liberty cong")
	in.WriteString("of the west library rang")
	in.WriteString("to become the wrector of worlds.")

	readerLen := in.Len()
	readerBytes := in.Bytes()

	writer := gen.FromReader{R: &in}

	written, err := writer.WriteTo(&out)
	if err != nil {
		tests.Failed("Should have written content without error: %+q.", err)
	}
	tests.Passed("Should have written content without error.")

	if int(written) != readerLen {
		tests.Info("ReaderLen: %d", readerLen)
		tests.Info("WriterLen: %d", written)
		tests.Failed("Should have written content length the same as reader.")
	}
	tests.Passed("Should have written content length the same as reader.")

	if !bytes.Equal(out.Bytes(), readerBytes) {
		tests.Failed("Should have written content exactly the same as reader.")
	}
	tests.Passed("Should have written content exactly the same as reader.")
}
