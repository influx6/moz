package mock

import (
	"io"

	toml "github.com/BurntSushi/toml"
)

// MofInitableImpl defines a concrete struct which implements the methods for the
// MofInitable interface. All methods will panic, so add necessary internal logic.
type MofInitableImpl struct {
}

// Ignite implements the MofInitable.Ignite() method for the MofInitable.
func (impl MofInitableImpl) Ignite() string {
	panic("Not yet implemented")
}

// Crunch implements the MofInitable.Crunch() method for the MofInitable.
func (impl MofInitableImpl) Crunch() string {
	panic("Not yet implemented")
}

// Configuration implements the MofInitable.Configuration() method for the MofInitable.
func (impl MofInitableImpl) Configuration() toml.Primitive {
	panic("Not yet implemented")
}

// Location implements the MofInitable.Location() method for the MofInitable.
func (impl MofInitableImpl) Location(var1 string) (GPSLoc, error) {
	panic("Not yet implemented")
}

// WriterTo implements the MofInitable.WriterTo() method for the MofInitable.
func (impl MofInitableImpl) WriterTo(var2 io.Writer) (int64, error) {
	panic("Not yet implemented")
}

// Drop implements the MofInitable.Drop() method for the MofInitable.
func (impl MofInitableImpl) Drop() (*GPSLoc, *toml.Primitive, *[]byte, *[5]byte) {
	panic("Not yet implemented")
}

// Close implements the MofInitable.Close() method for the MofInitable.
func (impl MofInitableImpl) Close() (chan struct{}, chan toml.Primitive, chan string, chan []byte, chan *[]string) {
	panic("Not yet implemented")
}

// Bob implements the MofInitable.Bob() method for the MofInitable.
func (impl MofInitableImpl) Bob() chan chan struct{} {
	panic("Not yet implemented")
}
