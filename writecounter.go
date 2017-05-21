package moz

import "io"

// WriteCounter defines a struct which collects write counts of
// a giving io.Writer
type WriteCounter struct {
	io.Writer
	written int64
}

// NewWriteCounter returns a new instance of the WriteCounter.
func NewWriteCounter(w io.Writer) *WriteCounter {
	return &WriteCounter{Writer: w}
}

// Written returns the total number of data writer to the underline writer.
func (w *WriteCounter) Written() int64 {
	return w.written
}

// Write calls the internal io.Writer.Write method and adds up
// the write counts.
func (w *WriteCounter) Write(data []byte) (int, error) {
	inc, err := w.Writer.Write(data)

	w.written += int64(inc)

	return inc, err
}

//======================================================================================================================

// ConstantWriter defines a writer that consistently writes a provided output.
type ConstantWriter struct {
	d []byte
}

// NewConstantWriter returns a new instance of ConstantWriter.
func NewConstantWriter(d []byte) ConstantWriter {
	return ConstantWriter{d: d}
}

// WriteTo writes the data provided into the writer.
func (cw ConstantWriter) WriteTo(w io.Writer) (int64, error) {
	total, err := w.Write(cw.d)
	return int64(total), err
}
