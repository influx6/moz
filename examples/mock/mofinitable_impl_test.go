package mock_test

import (
	"io"
	"testing"

	"github.com/influx6/moz/examples/mock"

	"github.com/influx6/faux/tests"

	"github.com/influx6/moz/examples/mock/snitch"

	toml "github.com/BurntSushi/toml"
)

// TestImplementationForMofInitable defines the test for asserting the behaviour of
// the implementation for the MofInitable interface methods.
func TestImplementationForMofInitable(t *testing.T) {
	t.Logf("Given the need to validate the behaviour of elements implementing the MofInitable interface")
	{

		testMethodCallForIgnite(t)

		testMethodCallForCrunch(t)

		testMethodCallForConfiguration(t)

		testMethodCallForLocation(t)

		testMethodCallForWriterTo(t)

		testMethodCallForMaps(t)

		testMethodCallForMapsIn(t)

		testMethodCallForMapsOut(t)

		testMethodCallForDrop(t)

		testMethodCallForClose(t)

		testMethodCallForBob(t)

	}
}

func testMethodCallForIgnite(t *testing.T) {
	t.Logf("\tWhen the method Ignite is called MofInitableImpl")
	{
		impl := &snitch.MofInitableMockSnitch{
			IgniteFunc: func() string {
				// Add implementation logic.
				panic("Please write your implementation logic in here for Ignite")
			},
		}

		// Stub variables for method.
		// TODO: Replace this stubs with real values for method

		// Call Ignite method with arguments
		impl.Ignite()

		if len(impl.IgniteMethodCalls) == 0 {
			tests.Failed("Should have received new method call record for Ignite.")
		}
		tests.Passed("Should have received new method call record for Ignite.")

		lastCall := impl.IgniteMethodCalls[len(impl.IgniteMethodCalls)-1]

		if lastCall.PanicError != nil {
			tests.Failed("Should have successfully executed Ignite method without panic.")
		}
		tests.Passed("Should have successfully executed Ignite method without panic.")
	}
}

func testMethodCallForCrunch(t *testing.T) {
	t.Logf("\tWhen the method Crunch is called MofInitableImpl")
	{
		impl := &snitch.MofInitableMockSnitch{
			CrunchFunc: func() string {
				// Add implementation logic.
				panic("Please write your implementation logic in here for Crunch")
			},
		}

		// Stub variables for method.
		// TODO: Replace this stubs with real values for method

		// Call Crunch method with arguments
		impl.Crunch()

		if len(impl.CrunchMethodCalls) == 0 {
			tests.Failed("Should have received new method call record for Crunch.")
		}
		tests.Passed("Should have received new method call record for Crunch.")

		lastCall := impl.CrunchMethodCalls[len(impl.CrunchMethodCalls)-1]

		if lastCall.PanicError != nil {
			tests.Failed("Should have successfully executed Crunch method without panic.")
		}
		tests.Passed("Should have successfully executed Crunch method without panic.")
	}
}

func testMethodCallForConfiguration(t *testing.T) {
	t.Logf("\tWhen the method Configuration is called MofInitableImpl")
	{
		impl := &snitch.MofInitableMockSnitch{
			ConfigurationFunc: func() toml.Primitive {
				// Add implementation logic.
				panic("Please write your implementation logic in here for Configuration")
			},
		}

		// Stub variables for method.
		// TODO: Replace this stubs with real values for method

		// Call Configuration method with arguments
		impl.Configuration()

		if len(impl.ConfigurationMethodCalls) == 0 {
			tests.Failed("Should have received new method call record for Configuration.")
		}
		tests.Passed("Should have received new method call record for Configuration.")

		lastCall := impl.ConfigurationMethodCalls[len(impl.ConfigurationMethodCalls)-1]

		if lastCall.PanicError != nil {
			tests.Failed("Should have successfully executed Configuration method without panic.")
		}
		tests.Passed("Should have successfully executed Configuration method without panic.")
	}
}

func testMethodCallForLocation(t *testing.T) {
	t.Logf("\tWhen the method Location is called MofInitableImpl")
	{
		impl := &snitch.MofInitableMockSnitch{
			LocationFunc: func(var1 string) (mock.GPSLoc, error) {
				// Add implementation logic.
				panic("Please write your implementation logic in here for Location")
			},
		}

		// Stub variables for method.
		// TODO: Replace this stubs with real values for method
		var var1 string

		// Call Location method with arguments
		impl.Location(var1)

		if len(impl.LocationMethodCalls) == 0 {
			tests.Failed("Should have received new method call record for Location.")
		}
		tests.Passed("Should have received new method call record for Location.")

		lastCall := impl.LocationMethodCalls[len(impl.LocationMethodCalls)-1]

		if lastCall.PanicError != nil {
			tests.Failed("Should have successfully executed Location method without panic.")
		}
		tests.Passed("Should have successfully executed Location method without panic.")
	}
}

func testMethodCallForWriterTo(t *testing.T) {
	t.Logf("\tWhen the method WriterTo is called MofInitableImpl")
	{
		impl := &snitch.MofInitableMockSnitch{
			WriterToFunc: func(var2 io.Writer) (int64, error) {
				// Add implementation logic.
				panic("Please write your implementation logic in here for WriterTo")
			},
		}

		// Stub variables for method.
		// TODO: Replace this stubs with real values for method
		var var2 io.Writer

		// Call WriterTo method with arguments
		impl.WriterTo(var2)

		if len(impl.WriterToMethodCalls) == 0 {
			tests.Failed("Should have received new method call record for WriterTo.")
		}
		tests.Passed("Should have received new method call record for WriterTo.")

		lastCall := impl.WriterToMethodCalls[len(impl.WriterToMethodCalls)-1]

		if lastCall.PanicError != nil {
			tests.Failed("Should have successfully executed WriterTo method without panic.")
		}
		tests.Passed("Should have successfully executed WriterTo method without panic.")
	}
}

func testMethodCallForMaps(t *testing.T) {
	t.Logf("\tWhen the method Maps is called MofInitableImpl")
	{
		impl := &snitch.MofInitableMockSnitch{
			MapsFunc: func(var3 string) (map[string]mock.GPSLoc, error) {
				// Add implementation logic.
				panic("Please write your implementation logic in here for Maps")
			},
		}

		// Stub variables for method.
		// TODO: Replace this stubs with real values for method
		var var3 string

		// Call Maps method with arguments
		impl.Maps(var3)

		if len(impl.MapsMethodCalls) == 0 {
			tests.Failed("Should have received new method call record for Maps.")
		}
		tests.Passed("Should have received new method call record for Maps.")

		lastCall := impl.MapsMethodCalls[len(impl.MapsMethodCalls)-1]

		if lastCall.PanicError != nil {
			tests.Failed("Should have successfully executed Maps method without panic.")
		}
		tests.Passed("Should have successfully executed Maps method without panic.")
	}
}

func testMethodCallForMapsIn(t *testing.T) {
	t.Logf("\tWhen the method MapsIn is called MofInitableImpl")
	{
		impl := &snitch.MofInitableMockSnitch{
			MapsInFunc: func(var4 string) (map[string]*mock.GPSLoc, error) {
				// Add implementation logic.
				panic("Please write your implementation logic in here for MapsIn")
			},
		}

		// Stub variables for method.
		// TODO: Replace this stubs with real values for method
		var var4 string

		// Call MapsIn method with arguments
		impl.MapsIn(var4)

		if len(impl.MapsInMethodCalls) == 0 {
			tests.Failed("Should have received new method call record for MapsIn.")
		}
		tests.Passed("Should have received new method call record for MapsIn.")

		lastCall := impl.MapsInMethodCalls[len(impl.MapsInMethodCalls)-1]

		if lastCall.PanicError != nil {
			tests.Failed("Should have successfully executed MapsIn method without panic.")
		}
		tests.Passed("Should have successfully executed MapsIn method without panic.")
	}
}

func testMethodCallForMapsOut(t *testing.T) {
	t.Logf("\tWhen the method MapsOut is called MofInitableImpl")
	{
		impl := &snitch.MofInitableMockSnitch{
			MapsOutFunc: func(var5 string) (map[*mock.GPSLoc]string, error) {
				// Add implementation logic.
				panic("Please write your implementation logic in here for MapsOut")
			},
		}

		// Stub variables for method.
		// TODO: Replace this stubs with real values for method
		var var5 string

		// Call MapsOut method with arguments
		impl.MapsOut(var5)

		if len(impl.MapsOutMethodCalls) == 0 {
			tests.Failed("Should have received new method call record for MapsOut.")
		}
		tests.Passed("Should have received new method call record for MapsOut.")

		lastCall := impl.MapsOutMethodCalls[len(impl.MapsOutMethodCalls)-1]

		if lastCall.PanicError != nil {
			tests.Failed("Should have successfully executed MapsOut method without panic.")
		}
		tests.Passed("Should have successfully executed MapsOut method without panic.")
	}
}

func testMethodCallForDrop(t *testing.T) {
	t.Logf("\tWhen the method Drop is called MofInitableImpl")
	{
		impl := &snitch.MofInitableMockSnitch{
			DropFunc: func() (*mock.GPSLoc, *toml.Primitive, *[]byte, *[5]byte) {
				// Add implementation logic.
				panic("Please write your implementation logic in here for Drop")
			},
		}

		// Stub variables for method.
		// TODO: Replace this stubs with real values for method

		// Call Drop method with arguments
		impl.Drop()

		if len(impl.DropMethodCalls) == 0 {
			tests.Failed("Should have received new method call record for Drop.")
		}
		tests.Passed("Should have received new method call record for Drop.")

		lastCall := impl.DropMethodCalls[len(impl.DropMethodCalls)-1]

		if lastCall.PanicError != nil {
			tests.Failed("Should have successfully executed Drop method without panic.")
		}
		tests.Passed("Should have successfully executed Drop method without panic.")
	}
}

func testMethodCallForClose(t *testing.T) {
	t.Logf("\tWhen the method Close is called MofInitableImpl")
	{
		impl := &snitch.MofInitableMockSnitch{
			CloseFunc: func() (chan struct{}, chan toml.Primitive, chan string, chan []byte, chan *[]string) {
				// Add implementation logic.
				panic("Please write your implementation logic in here for Close")
			},
		}

		// Stub variables for method.
		// TODO: Replace this stubs with real values for method

		// Call Close method with arguments
		impl.Close()

		if len(impl.CloseMethodCalls) == 0 {
			tests.Failed("Should have received new method call record for Close.")
		}
		tests.Passed("Should have received new method call record for Close.")

		lastCall := impl.CloseMethodCalls[len(impl.CloseMethodCalls)-1]

		if lastCall.PanicError != nil {
			tests.Failed("Should have successfully executed Close method without panic.")
		}
		tests.Passed("Should have successfully executed Close method without panic.")
	}
}

func testMethodCallForBob(t *testing.T) {
	t.Logf("\tWhen the method Bob is called MofInitableImpl")
	{
		impl := &snitch.MofInitableMockSnitch{
			BobFunc: func() chan chan struct{} {
				// Add implementation logic.
				panic("Please write your implementation logic in here for Bob")
			},
		}

		// Stub variables for method.
		// TODO: Replace this stubs with real values for method

		// Call Bob method with arguments
		impl.Bob()

		if len(impl.BobMethodCalls) == 0 {
			tests.Failed("Should have received new method call record for Bob.")
		}
		tests.Passed("Should have received new method call record for Bob.")

		lastCall := impl.BobMethodCalls[len(impl.BobMethodCalls)-1]

		if lastCall.PanicError != nil {
			tests.Failed("Should have successfully executed Bob method without panic.")
		}
		tests.Passed("Should have successfully executed Bob method without panic.")
	}
}
