package snitch

import (
	"io"
	"time"

	"runtime"

	"github.com/influx6/moz/examples/mock"

	toml "github.com/BurntSushi/toml"
)

// MethodCallForIgnite defines a type which holds meta-details about the giving calls associated
// with the MofInitable.Ignite() method.
type MethodCallForIgnite struct {
	When  time.Time
	Start time.Time
	End   time.Time

	// Details of panic if such occurs.
	PanicStack []byte
	PanicError interface{}

	// Argument values.

	// Return values.

	Var1 string
}

// MethodCallForCrunch defines a type which holds meta-details about the giving calls associated
// with the MofInitable.Crunch() method.
type MethodCallForCrunch struct {
	When  time.Time
	Start time.Time
	End   time.Time

	// Details of panic if such occurs.
	PanicStack []byte
	PanicError interface{}

	// Argument values.

	// Return values.

	Cr string
}

// MethodCallForConfiguration defines a type which holds meta-details about the giving calls associated
// with the MofInitable.Configuration() method.
type MethodCallForConfiguration struct {
	When  time.Time
	Start time.Time
	End   time.Time

	// Details of panic if such occurs.
	PanicStack []byte
	PanicError interface{}

	// Argument values.

	// Return values.

	Var1 toml.Primitive
}

// MethodCallForLocation defines a type which holds meta-details about the giving calls associated
// with the MofInitable.Location() method.
type MethodCallForLocation struct {
	When  time.Time
	Start time.Time
	End   time.Time

	// Details of panic if such occurs.
	PanicStack []byte
	PanicError interface{}

	// Argument values.

	Var1 string

	// Return values.

	Var2 mock.GPSLoc

	Var3 error
}

// MethodCallForWriterTo defines a type which holds meta-details about the giving calls associated
// with the MofInitable.WriterTo() method.
type MethodCallForWriterTo struct {
	When  time.Time
	Start time.Time
	End   time.Time

	// Details of panic if such occurs.
	PanicStack []byte
	PanicError interface{}

	// Argument values.

	Var2 io.Writer

	// Return values.

	Var4 int64

	Var5 error
}

// MethodCallForMaps defines a type which holds meta-details about the giving calls associated
// with the MofInitable.Maps() method.
type MethodCallForMaps struct {
	When  time.Time
	Start time.Time
	End   time.Time

	// Details of panic if such occurs.
	PanicStack []byte
	PanicError interface{}

	// Argument values.

	Var3 string

	// Return values.

	Var6 map[string]mock.GPSLoc

	Var7 error
}

// MethodCallForMapsIn defines a type which holds meta-details about the giving calls associated
// with the MofInitable.MapsIn() method.
type MethodCallForMapsIn struct {
	When  time.Time
	Start time.Time
	End   time.Time

	// Details of panic if such occurs.
	PanicStack []byte
	PanicError interface{}

	// Argument values.

	Var4 string

	// Return values.

	Var8 map[string]*mock.GPSLoc

	Var9 error
}

// MethodCallForMapsOut defines a type which holds meta-details about the giving calls associated
// with the MofInitable.MapsOut() method.
type MethodCallForMapsOut struct {
	When  time.Time
	Start time.Time
	End   time.Time

	// Details of panic if such occurs.
	PanicStack []byte
	PanicError interface{}

	// Argument values.

	Var5 string

	// Return values.

	Var10 map[*mock.GPSLoc]string

	Var11 error
}

// MethodCallForDrop defines a type which holds meta-details about the giving calls associated
// with the MofInitable.Drop() method.
type MethodCallForDrop struct {
	When  time.Time
	Start time.Time
	End   time.Time

	// Details of panic if such occurs.
	PanicStack []byte
	PanicError interface{}

	// Argument values.

	// Return values.

	Var12 *mock.GPSLoc

	Var13 *toml.Primitive

	Var14 *[]byte

	Var15 *[5]byte
}

// MethodCallForClose defines a type which holds meta-details about the giving calls associated
// with the MofInitable.Close() method.
type MethodCallForClose struct {
	When  time.Time
	Start time.Time
	End   time.Time

	// Details of panic if such occurs.
	PanicStack []byte
	PanicError interface{}

	// Argument values.

	// Return values.

	Var16 chan struct{}

	Var17 chan toml.Primitive

	Var18 chan string

	Var19 chan []byte

	Var20 chan *[]string
}

// MethodCallForBob defines a type which holds meta-details about the giving calls associated
// with the MofInitable.Bob() method.
type MethodCallForBob struct {
	When  time.Time
	Start time.Time
	End   time.Time

	// Details of panic if such occurs.
	PanicStack []byte
	PanicError interface{}

	// Argument values.

	// Return values.

	Var21 chan chan struct{}
}

// MofInitableSnitch is an implementation of the MofInitable which will wrap the
// MofInitableImpl and allow you to record all arguments and return values from calls to the
// provided implementation which it wraps.
type MofInitableLittleSnitch struct {
	Implementer mock.MofInitable

	IgniteMethodCalls []MethodCallForIgnite

	CrunchMethodCalls []MethodCallForCrunch

	ConfigurationMethodCalls []MethodCallForConfiguration

	LocationMethodCalls []MethodCallForLocation

	WriterToMethodCalls []MethodCallForWriterTo

	MapsMethodCalls []MethodCallForMaps

	MapsInMethodCalls []MethodCallForMapsIn

	MapsOutMethodCalls []MethodCallForMapsOut

	DropMethodCalls []MethodCallForDrop

	CloseMethodCalls []MethodCallForClose

	BobMethodCalls []MethodCallForBob
}

// NewLittleSnitch returns a new instance of a MofInitableLittleSnitch for recording.
func NewLittleSnitch(impl mock.MofInitable) *MofInitableLittleSnitch {
	var snitch MofInitableLittleSnitch

	snitch.IgniteMethodCalls = make([]MethodCallForIgnite, 0)

	snitch.CrunchMethodCalls = make([]MethodCallForCrunch, 0)

	snitch.ConfigurationMethodCalls = make([]MethodCallForConfiguration, 0)

	snitch.LocationMethodCalls = make([]MethodCallForLocation, 0)

	snitch.WriterToMethodCalls = make([]MethodCallForWriterTo, 0)

	snitch.MapsMethodCalls = make([]MethodCallForMaps, 0)

	snitch.MapsInMethodCalls = make([]MethodCallForMapsIn, 0)

	snitch.MapsOutMethodCalls = make([]MethodCallForMapsOut, 0)

	snitch.DropMethodCalls = make([]MethodCallForDrop, 0)

	snitch.CloseMethodCalls = make([]MethodCallForClose, 0)

	snitch.BobMethodCalls = make([]MethodCallForBob, 0)

	return &snitch
}

// Ignite implements the MofInitable.Ignite() method for the MofInitable.
func (impl *MofInitableLittleSnitch) Ignite() string {
	var caller MethodCallForIgnite

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.IgniteMethodCalls = append(impl.IgniteMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	var1 := impl.Implementer.Ignite()

	caller.Var1 = var1

	return var1
}

// Crunch implements the MofInitable.Crunch() method for the MofInitable.
func (impl *MofInitableLittleSnitch) Crunch() string {
	var caller MethodCallForCrunch

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.CrunchMethodCalls = append(impl.CrunchMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	cr := impl.Implementer.Crunch()

	caller.Cr = cr

	return cr
}

// Configuration implements the MofInitable.Configuration() method for the MofInitable.
func (impl *MofInitableLittleSnitch) Configuration() toml.Primitive {
	var caller MethodCallForConfiguration

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.ConfigurationMethodCalls = append(impl.ConfigurationMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	var1 := impl.Implementer.Configuration()

	caller.Var1 = var1

	return var1
}

// Location implements the MofInitable.Location() method for the MofInitable.
func (impl *MofInitableLittleSnitch) Location(var1 string) (mock.GPSLoc, error) {
	var caller MethodCallForLocation

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.LocationMethodCalls = append(impl.LocationMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	caller.Var1 = var1

	var2, var3 := impl.Implementer.Location(var1)

	caller.Var2 = var2

	caller.Var3 = var3

	return var2, var3
}

// WriterTo implements the MofInitable.WriterTo() method for the MofInitable.
func (impl *MofInitableLittleSnitch) WriterTo(var2 io.Writer) (int64, error) {
	var caller MethodCallForWriterTo

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.WriterToMethodCalls = append(impl.WriterToMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	caller.Var2 = var2

	var4, var5 := impl.Implementer.WriterTo(var2)

	caller.Var4 = var4

	caller.Var5 = var5

	return var4, var5
}

// Maps implements the MofInitable.Maps() method for the MofInitable.
func (impl *MofInitableLittleSnitch) Maps(var3 string) (map[string]mock.GPSLoc, error) {
	var caller MethodCallForMaps

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.MapsMethodCalls = append(impl.MapsMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	caller.Var3 = var3

	var6, var7 := impl.Implementer.Maps(var3)

	caller.Var6 = var6

	caller.Var7 = var7

	return var6, var7
}

// MapsIn implements the MofInitable.MapsIn() method for the MofInitable.
func (impl *MofInitableLittleSnitch) MapsIn(var4 string) (map[string]*mock.GPSLoc, error) {
	var caller MethodCallForMapsIn

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.MapsInMethodCalls = append(impl.MapsInMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	caller.Var4 = var4

	var8, var9 := impl.Implementer.MapsIn(var4)

	caller.Var8 = var8

	caller.Var9 = var9

	return var8, var9
}

// MapsOut implements the MofInitable.MapsOut() method for the MofInitable.
func (impl *MofInitableLittleSnitch) MapsOut(var5 string) (map[*mock.GPSLoc]string, error) {
	var caller MethodCallForMapsOut

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.MapsOutMethodCalls = append(impl.MapsOutMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	caller.Var5 = var5

	var10, var11 := impl.Implementer.MapsOut(var5)

	caller.Var10 = var10

	caller.Var11 = var11

	return var10, var11
}

// Drop implements the MofInitable.Drop() method for the MofInitable.
func (impl *MofInitableLittleSnitch) Drop() (*mock.GPSLoc, *toml.Primitive, *[]byte, *[5]byte) {
	var caller MethodCallForDrop

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.DropMethodCalls = append(impl.DropMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	var12, var13, var14, var15 := impl.Implementer.Drop()

	caller.Var12 = var12

	caller.Var13 = var13

	caller.Var14 = var14

	caller.Var15 = var15

	return var12, var13, var14, var15
}

// Close implements the MofInitable.Close() method for the MofInitable.
func (impl *MofInitableLittleSnitch) Close() (chan struct{}, chan toml.Primitive, chan string, chan []byte, chan *[]string) {
	var caller MethodCallForClose

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.CloseMethodCalls = append(impl.CloseMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	var16, var17, var18, var19, var20 := impl.Implementer.Close()

	caller.Var16 = var16

	caller.Var17 = var17

	caller.Var18 = var18

	caller.Var19 = var19

	caller.Var20 = var20

	return var16, var17, var18, var19, var20
}

// Bob implements the MofInitable.Bob() method for the MofInitable.
func (impl *MofInitableLittleSnitch) Bob() chan chan struct{} {
	var caller MethodCallForBob

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.BobMethodCalls = append(impl.BobMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	var21 := impl.Implementer.Bob()

	caller.Var21 = var21

	return var21
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

	ConfigurationMethodCalls []MethodCallForConfiguration
	ConfigurationFunc        func() toml.Primitive

	LocationMethodCalls []MethodCallForLocation
	LocationFunc        func(var1 string) (mock.GPSLoc, error)

	WriterToMethodCalls []MethodCallForWriterTo
	WriterToFunc        func(var2 io.Writer) (int64, error)

	MapsMethodCalls []MethodCallForMaps
	MapsFunc        func(var3 string) (map[string]mock.GPSLoc, error)

	MapsInMethodCalls []MethodCallForMapsIn
	MapsInFunc        func(var4 string) (map[string]*mock.GPSLoc, error)

	MapsOutMethodCalls []MethodCallForMapsOut
	MapsOutFunc        func(var5 string) (map[*mock.GPSLoc]string, error)

	DropMethodCalls []MethodCallForDrop
	DropFunc        func() (*mock.GPSLoc, *toml.Primitive, *[]byte, *[5]byte)

	CloseMethodCalls []MethodCallForClose
	CloseFunc        func() (chan struct{}, chan toml.Primitive, chan string, chan []byte, chan *[]string)

	BobMethodCalls []MethodCallForBob
	BobFunc        func() chan chan struct{}
}

// Ignite implements the MofInitable.Ignite() method for the MofInitable.
func (impl *MofInitableMockSnitch) Ignite() string {
	var caller MethodCallForIgnite

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.IgniteMethodCalls = append(impl.IgniteMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	var1 := impl.IgniteFunc()

	caller.Var1 = var1

	return var1
}

// Crunch implements the MofInitable.Crunch() method for the MofInitable.
func (impl *MofInitableMockSnitch) Crunch() string {
	var caller MethodCallForCrunch

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.CrunchMethodCalls = append(impl.CrunchMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	cr := impl.CrunchFunc()

	caller.Cr = cr

	return cr
}

// Configuration implements the MofInitable.Configuration() method for the MofInitable.
func (impl *MofInitableMockSnitch) Configuration() toml.Primitive {
	var caller MethodCallForConfiguration

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.ConfigurationMethodCalls = append(impl.ConfigurationMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	var1 := impl.ConfigurationFunc()

	caller.Var1 = var1

	return var1
}

// Location implements the MofInitable.Location() method for the MofInitable.
func (impl *MofInitableMockSnitch) Location(var1 string) (mock.GPSLoc, error) {
	var caller MethodCallForLocation

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.LocationMethodCalls = append(impl.LocationMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	caller.Var1 = var1

	var2, var3 := impl.LocationFunc(var1)

	caller.Var2 = var2

	caller.Var3 = var3

	return var2, var3
}

// WriterTo implements the MofInitable.WriterTo() method for the MofInitable.
func (impl *MofInitableMockSnitch) WriterTo(var2 io.Writer) (int64, error) {
	var caller MethodCallForWriterTo

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.WriterToMethodCalls = append(impl.WriterToMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	caller.Var2 = var2

	var4, var5 := impl.WriterToFunc(var2)

	caller.Var4 = var4

	caller.Var5 = var5

	return var4, var5
}

// Maps implements the MofInitable.Maps() method for the MofInitable.
func (impl *MofInitableMockSnitch) Maps(var3 string) (map[string]mock.GPSLoc, error) {
	var caller MethodCallForMaps

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.MapsMethodCalls = append(impl.MapsMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	caller.Var3 = var3

	var6, var7 := impl.MapsFunc(var3)

	caller.Var6 = var6

	caller.Var7 = var7

	return var6, var7
}

// MapsIn implements the MofInitable.MapsIn() method for the MofInitable.
func (impl *MofInitableMockSnitch) MapsIn(var4 string) (map[string]*mock.GPSLoc, error) {
	var caller MethodCallForMapsIn

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.MapsInMethodCalls = append(impl.MapsInMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	caller.Var4 = var4

	var8, var9 := impl.MapsInFunc(var4)

	caller.Var8 = var8

	caller.Var9 = var9

	return var8, var9
}

// MapsOut implements the MofInitable.MapsOut() method for the MofInitable.
func (impl *MofInitableMockSnitch) MapsOut(var5 string) (map[*mock.GPSLoc]string, error) {
	var caller MethodCallForMapsOut

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.MapsOutMethodCalls = append(impl.MapsOutMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	caller.Var5 = var5

	var10, var11 := impl.MapsOutFunc(var5)

	caller.Var10 = var10

	caller.Var11 = var11

	return var10, var11
}

// Drop implements the MofInitable.Drop() method for the MofInitable.
func (impl *MofInitableMockSnitch) Drop() (*mock.GPSLoc, *toml.Primitive, *[]byte, *[5]byte) {
	var caller MethodCallForDrop

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.DropMethodCalls = append(impl.DropMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	var12, var13, var14, var15 := impl.DropFunc()

	caller.Var12 = var12

	caller.Var13 = var13

	caller.Var14 = var14

	caller.Var15 = var15

	return var12, var13, var14, var15
}

// Close implements the MofInitable.Close() method for the MofInitable.
func (impl *MofInitableMockSnitch) Close() (chan struct{}, chan toml.Primitive, chan string, chan []byte, chan *[]string) {
	var caller MethodCallForClose

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.CloseMethodCalls = append(impl.CloseMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	var16, var17, var18, var19, var20 := impl.CloseFunc()

	caller.Var16 = var16

	caller.Var17 = var17

	caller.Var18 = var18

	caller.Var19 = var19

	caller.Var20 = var20

	return var16, var17, var18, var19, var20
}

// Bob implements the MofInitable.Bob() method for the MofInitable.
func (impl *MofInitableMockSnitch) Bob() chan chan struct{} {
	var caller MethodCallForBob

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, true)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.BobMethodCalls = append(impl.BobMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	var21 := impl.BobFunc()

	caller.Var21 = var21

	return var21
}
