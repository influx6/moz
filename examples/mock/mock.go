package mock

import (
	"io"

	toml "github.com/BurntSushi/toml"
)

// Ignitable defines a struct which is used to ignite the package.
type Ignitable interface {
	Ignite() string
}

// GPSLoc defines a struct to hold long and lat values for a gps location.
type GPSLoc struct {
	Lat  float64
	Long float64
}

// MofInitable defines a interface for a Mof.
// @iface
type MofInitable interface {
	Ignitable
	Crunch() (cr string)
	Configuration() toml.Primitive
	Location(string) (GPSLoc, error)
	WriterTo(io.Writer) (int64, error)
	Maps(string) (map[string]GPSLoc, error)
	MapsIn(string) (map[string]*GPSLoc, error)
	MapsOut(string) (map[*GPSLoc]string, error)
	Drop() (*GPSLoc, *toml.Primitive, *[]byte, *[5]byte)
	Close() (chan struct{}, chan toml.Primitive, chan string, chan []byte, chan *[]string)
	Bob() chan chan struct{}
}
