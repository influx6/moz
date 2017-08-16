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

	Var2 GPSLoc

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

	Var6 *GPSLoc

	Var7 *toml.Primitive

	Var8 *[]byte

	Var9 *[5]byte
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

	Var10 chan struct{}

	Var11 chan toml.Primitive

	Var12 chan string

	Var13 chan []byte

	Var14 chan *[]string
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

	Var15 chan chan struct{}
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
			trace = trace[:runtime.Stack(trace, all)]

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
			trace = trace[:runtime.Stack(trace, all)]

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
			trace = trace[:runtime.Stack(trace, all)]

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
func (impl *MofInitableLittleSnitch) Location(var1 string) (GPSLoc, error) {
	var caller MethodCallForLocation

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, all)]

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
			trace = trace[:runtime.Stack(trace, all)]

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

// Drop implements the MofInitable.Drop() method for the MofInitable.
func (impl *MofInitableLittleSnitch) Drop() (*GPSLoc, *toml.Primitive, *[]byte, *[5]byte) {
	var caller MethodCallForDrop

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, all)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.DropMethodCalls = append(impl.DropMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	var6, var7, var8, var9 := impl.Implementer.Drop()

	caller.Var6 = var6

	caller.Var7 = var7

	caller.Var8 = var8

	caller.Var9 = var9

	return var6, var7, var8, var9
}

// Close implements the MofInitable.Close() method for the MofInitable.
func (impl *MofInitableLittleSnitch) Close() (chan struct{}, chan toml.Primitive, chan string, chan []byte, chan *[]string) {
	var caller MethodCallForClose

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, all)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.CloseMethodCalls = append(impl.CloseMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	var10, var11, var12, var13, var14 := impl.Implementer.Close()

	caller.Var10 = var10

	caller.Var11 = var11

	caller.Var12 = var12

	caller.Var13 = var13

	caller.Var14 = var14

	return var10, var11, var12, var13, var14
}

// Bob implements the MofInitable.Bob() method for the MofInitable.
func (impl *MofInitableLittleSnitch) Bob() chan chan struct{} {
	var caller MethodCallForBob

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, all)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.BobMethodCalls = append(impl.BobMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	var15 := impl.Implementer.Bob()

	caller.Var15 = var15

	return var15
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
	LocationFunc        func(var1 string) (GPSLoc, error)

	WriterToMethodCalls []MethodCallForWriterTo
	WriterToFunc        func(var2 io.Writer) (int64, error)

	DropMethodCalls []MethodCallForDrop
	DropFunc        func() (*GPSLoc, *toml.Primitive, *[]byte, *[5]byte)

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
			trace = trace[:runtime.Stack(trace, all)]

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
			trace = trace[:runtime.Stack(trace, all)]

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
			trace = trace[:runtime.Stack(trace, all)]

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
func (impl *MofInitableMockSnitch) Location(var1 string) (GPSLoc, error) {
	var caller MethodCallForLocation

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, all)]

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
			trace = trace[:runtime.Stack(trace, all)]

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

// Drop implements the MofInitable.Drop() method for the MofInitable.
func (impl *MofInitableMockSnitch) Drop() (*GPSLoc, *toml.Primitive, *[]byte, *[5]byte) {
	var caller MethodCallForDrop

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, all)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.DropMethodCalls = append(impl.DropMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	var6, var7, var8, var9 := impl.DropFunc()

	caller.Var6 = var6

	caller.Var7 = var7

	caller.Var8 = var8

	caller.Var9 = var9

	return var6, var7, var8, var9
}

// Close implements the MofInitable.Close() method for the MofInitable.
func (impl *MofInitableMockSnitch) Close() (chan struct{}, chan toml.Primitive, chan string, chan []byte, chan *[]string) {
	var caller MethodCallForClose

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, all)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.CloseMethodCalls = append(impl.CloseMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	var10, var11, var12, var13, var14 := impl.CloseFunc()

	caller.Var10 = var10

	caller.Var11 = var11

	caller.Var12 = var12

	caller.Var13 = var13

	caller.Var14 = var14

	return var10, var11, var12, var13, var14
}

// Bob implements the MofInitable.Bob() method for the MofInitable.
func (impl *MofInitableMockSnitch) Bob() chan chan struct{} {
	var caller MethodCallForBob

	defer func() {
		if err := recover(); err != nil {
			trace := make([]byte, 1000)
			trace = trace[:runtime.Stack(trace, all)]

			caller.PanicError = err
			caller.PanicStack = trace
		}

		caller.End = time.Now()
		impl.BobMethodCalls = append(impl.BobMethodCalls, caller)
	}()

	caller.When = time.Now()
	caller.Start = caller.When

	var15 := impl.BobFunc()

	caller.Var15 = var15

	return var15
}
