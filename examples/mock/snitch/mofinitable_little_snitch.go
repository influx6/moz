package littlesnitch

import (
	"github.com/influx6/moz/examples/mock"

	"time"
)

// MethodCallForIgnite defines a type which holds meta-details about the giving calls associated
// with the MofInitable.Ignite() method.
type MethodCallForIgnite struct {
	When  time.Time
	Start time.Time
	End   time.Time

	// Argument values.

	// Return values.

	Ret1 string
}

// MethodCallForCrunch defines a type which holds meta-details about the giving calls associated
// with the MofInitable.Crunch() method.
type MethodCallForCrunch struct {
	When  time.Time
	Start time.Time
	End   time.Time

	// Argument values.

	// Return values.

	Cr string
}

// MethodCallForLocation defines a type which holds meta-details about the giving calls associated
// with the MofInitable.Location() method.
type MethodCallForLocation struct {
	When  time.Time
	Start time.Time
	End   time.Time

	// Argument values.

	Var1 string

	// Return values.

	Ret1 GPSLoc

	Ret2 error
}

// MofInitableSnitch is an implementation of the MofInitable which will wrap the
// MofInitableImpl and allow you to record all arguments and return values from calls to the
// provided implementation which it wraps.
type MofInitableLittleSnitch struct {
	Implementer mock.MofInitable

	IgniteMethodCalls []MethodCallForIgnite

	CrunchMethodCalls []MethodCallForCrunch

	LocationMethodCalls []MethodCallForLocation
}

// New returns a new instance of a MofInitableLittleSnitch for recording.
func New(impl mock.MofInitable) *MofInitableLittleSnitch {
	var snitch MofInitableLittleSnitch

	snitch.IgniteMethodCalls = make([]MethodCallForIgnite, 0)

	snitch.CrunchMethodCalls = make([]MethodCallForCrunch, 0)

	snitch.LocationMethodCalls = make([]MethodCallForLocation, 0)

	return &snitch
}

// Ignite implements the MofInitable.Ignite() method for the MofInitable.
func (impl *MofInitableLittleSnitch) Ignite() string {
	var caller MethodCallForIgnite
	caller.When = time.Now()
	caller.Start = caller.When

	ret1 := impl.Implementer.Ignite()

	caller.Ret1 = ret1

	caller.End = time.Now()

	impl.IgniteMethodCalls = append(impl.IgniteMethodCalls, caller)

	return ret1
}

// Crunch implements the MofInitable.Crunch() method for the MofInitable.
func (impl *MofInitableLittleSnitch) Crunch() string {
	var caller MethodCallForCrunch
	caller.When = time.Now()
	caller.Start = caller.When

	cr := impl.Implementer.Crunch()

	caller.Cr = cr

	caller.End = time.Now()

	impl.CrunchMethodCalls = append(impl.CrunchMethodCalls, caller)

	return cr
}

// Location implements the MofInitable.Location() method for the MofInitable.
func (impl *MofInitableLittleSnitch) Location(var1 string) (GPSLoc, error) {
	var caller MethodCallForLocation
	caller.When = time.Now()
	caller.Start = caller.When

	caller.Var1 = var1

	ret1, ret2 := impl.Implementer.Location(var1)

	caller.Ret1 = ret1

	caller.Ret2 = ret2

	caller.End = time.Now()

	impl.LocationMethodCalls = append(impl.LocationMethodCalls, caller)

	return ret1, ret2
}

//==============================================================================================================

// MofInitableMockSnitch defines a function type which implements a struct with the
// methods for the MofInitable as fields which allows you provide implementations of
// these functions to provide flexible testing.
type MofInitableMockSnitch struct {
	IgniteMethodCalls []MethodCallForIgnite
	IgniteFunc        func() string

	CrunchMethodCalls []MethodCallForCrunch
	CrunchFunc        func() string

	LocationMethodCalls []MethodCallForLocation
	LocationFunc        func(var1 string) (GPSLoc, error)
}

// Ignite implements the MofInitable.Ignite() method for the MofInitable.
func (impl *MofInitableMockSnitch) Ignite() string {
	var caller MethodCallForIgnite
	caller.When = time.Now()
	caller.Start = caller.When

	ret1 := impl.IgniteFunc()

	caller.Ret1 = ret1

	caller.End = time.Now()

	impl.IgniteMethodCalls = append(impl.IgniteMethodCalls, caller)

	return ret1
}

// Crunch implements the MofInitable.Crunch() method for the MofInitable.
func (impl *MofInitableMockSnitch) Crunch() string {
	var caller MethodCallForCrunch
	caller.When = time.Now()
	caller.Start = caller.When

	cr := impl.CrunchFunc()

	caller.Cr = cr

	caller.End = time.Now()

	impl.CrunchMethodCalls = append(impl.CrunchMethodCalls, caller)

	return cr
}

// Location implements the MofInitable.Location() method for the MofInitable.
func (impl *MofInitableMockSnitch) Location(var1 string) (GPSLoc, error) {
	var caller MethodCallForLocation
	caller.When = time.Now()
	caller.Start = caller.When

	caller.Var1 = var1

	ret1, ret2 := impl.LocationFunc(var1)

	caller.Ret1 = ret1

	caller.Ret2 = ret2

	caller.End = time.Now()

	impl.LocationMethodCalls = append(impl.LocationMethodCalls, caller)

	return ret1, ret2
}
