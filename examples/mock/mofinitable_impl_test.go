package mock_test

import (
	"testing"

	"github.com/influx6/faux/tests"

	"github.com/influx6/moz/examples/mock/snitch"
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

		if lastCall.PanicErr != nil {
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

		if lastCall.PanicErr != nil {
			tests.Failed("Should have successfully executed Crunch method without panic.")
		}
		tests.Passed("Should have successfully executed Crunch method without panic.")
	}
}

func testMethodCallForConfiguration(t *testing.T) {
	t.Logf("\tWhen the method Configuration is called MofInitableImpl")
	{
		impl := &snitch.MofInitableMockSnitch{
			ConfigurationFunc: func() {
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

		if lastCall.PanicErr != nil {
			tests.Failed("Should have successfully executed Configuration method without panic.")
		}
		tests.Passed("Should have successfully executed Configuration method without panic.")
	}
}

func testMethodCallForLocation(t *testing.T) {
	t.Logf("\tWhen the method Location is called MofInitableImpl")
	{
		impl := &snitch.MofInitableMockSnitch{
			LocationFunc: func(var1 string) (GPSLoc, error) {
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

		if lastCall.PanicErr != nil {
			tests.Failed("Should have successfully executed Location method without panic.")
		}
		tests.Passed("Should have successfully executed Location method without panic.")
	}
}

func testMethodCallForWriterTo(t *testing.T) {
	t.Logf("\tWhen the method WriterTo is called MofInitableImpl")
	{
		impl := &snitch.MofInitableMockSnitch{
			WriterToFunc: func() (int64, error) {
				// Add implementation logic.
				panic("Please write your implementation logic in here for WriterTo")
			},
		}

		// Stub variables for method.
		// TODO: Replace this stubs with real values for method

		// Call WriterTo method with arguments
		impl.WriterTo()

		if len(impl.WriterToMethodCalls) == 0 {
			tests.Failed("Should have received new method call record for WriterTo.")
		}
		tests.Passed("Should have received new method call record for WriterTo.")

		lastCall := impl.WriterToMethodCalls[len(impl.WriterToMethodCalls)-1]

		if lastCall.PanicErr != nil {
			tests.Failed("Should have successfully executed WriterTo method without panic.")
		}
		tests.Passed("Should have successfully executed WriterTo method without panic.")
	}
}
