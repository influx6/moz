package mock

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
func (impl MofInitableImpl) Configuration() {
	panic("Not yet implemented")
}

// Location implements the MofInitable.Location() method for the MofInitable.
func (impl MofInitableImpl) Location(var1 string) (GPSLoc, error) {
	panic("Not yet implemented")
}

// WriterTo implements the MofInitable.WriterTo() method for the MofInitable.
func (impl MofInitableImpl) WriterTo() (int64, error) {
	panic("Not yet implemented")
}
