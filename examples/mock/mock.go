package mock

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
	Location(string) (GPSLoc, error)
}
