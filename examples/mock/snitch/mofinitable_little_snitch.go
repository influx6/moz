<standard input>:243:6: expected statement, found ':='
<standard input>:428:6: expected statement, found ':='

-----------------------
package snitch

import (
     "time"


     "runtime"


     "github.com/influx6/moz/examples/mock"

)




// MethodCallForIgnite defines a type which holds meta-details about the giving calls associated 
// with the MofInitable.Ignite() method.
type MethodCallForIgnite struct{
    When time.Time
    Start time.Time
    End time.Time

    // Details of panic if such occurs.
    PanicStack []byte
    PanicError interface{}

    // Argument values.
    

    // Return values.
    
    Var1 string
    
}

// MethodCallForCrunch defines a type which holds meta-details about the giving calls associated 
// with the MofInitable.Crunch() method.
type MethodCallForCrunch struct{
    When time.Time
    Start time.Time
    End time.Time

    // Details of panic if such occurs.
    PanicStack []byte
    PanicError interface{}

    // Argument values.
    

    // Return values.
    
    Cr string
    
}

// MethodCallForConfiguration defines a type which holds meta-details about the giving calls associated 
// with the MofInitable.Configuration() method.
type MethodCallForConfiguration struct{
    When time.Time
    Start time.Time
    End time.Time

    // Details of panic if such occurs.
    PanicStack []byte
    PanicError interface{}

    // Argument values.
    

    // Return values.
    
}

// MethodCallForLocation defines a type which holds meta-details about the giving calls associated 
// with the MofInitable.Location() method.
type MethodCallForLocation struct{
    When time.Time
    Start time.Time
    End time.Time

    // Details of panic if such occurs.
    PanicStack []byte
    PanicError interface{}

    // Argument values.
    
    Var1 string
    

    // Return values.
    
    Var1 GPSLoc
    
    Var2 error
    
}

// MethodCallForWriterTo defines a type which holds meta-details about the giving calls associated 
// with the MofInitable.WriterTo() method.
type MethodCallForWriterTo struct{
    When time.Time
    Start time.Time
    End time.Time

    // Details of panic if such occurs.
    PanicStack []byte
    PanicError interface{}

    // Argument values.
    

    // Return values.
    
    Var3 int64
    
    Var4 error
    
}


// MofInitableSnitch is an implementation of the MofInitable which will wrap the 
// MofInitableImpl and allow you to record all arguments and return values from calls to the 
// provided implementation which it wraps.
type MofInitableLittleSnitch struct{
    Implementer mock.MofInitable
    
    IgniteMethodCalls []MethodCallForIgnite
    
    CrunchMethodCalls []MethodCallForCrunch
    
    ConfigurationMethodCalls []MethodCallForConfiguration
    
    LocationMethodCalls []MethodCallForLocation
    
    WriterToMethodCalls []MethodCallForWriterTo
    
}

// NewLittleSnitch returns a new instance of a MofInitableLittleSnitch for recording.
func NewLittleSnitch(impl mock.MofInitable) *MofInitableLittleSnitch {
    var snitch MofInitableLittleSnitch
    
    snitch.IgniteMethodCalls = make([]MethodCallForIgnite, 0)
    
    snitch.CrunchMethodCalls = make([]MethodCallForCrunch, 0)
    
    snitch.ConfigurationMethodCalls = make([]MethodCallForConfiguration, 0)
    
    snitch.LocationMethodCalls = make([]MethodCallForLocation, 0)
    
    snitch.WriterToMethodCalls = make([]MethodCallForWriterTo, 0)
    

    return &snitch
}

 
// Ignite implements the MofInitable.Ignite() method for the MofInitable.
func (impl *MofInitableLittleSnitch) Ignite() (string){
    var caller MethodCallForIgnite

    defer func(){
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
func (impl *MofInitableLittleSnitch) Crunch() (string){
    var caller MethodCallForCrunch

    defer func(){
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
func (impl *MofInitableLittleSnitch) Configuration() (){
    var caller MethodCallForConfiguration

    defer func(){
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

    

     := impl.Implementer.Configuration()

    

    return 
}
 
// Location implements the MofInitable.Location() method for the MofInitable.
func (impl *MofInitableLittleSnitch) Location(var1 string) (GPSLoc,error){
    var caller MethodCallForLocation

    defer func(){
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
    

    var1,var2 := impl.Implementer.Location(var1)

    
    caller.Var1 = var1
    
    caller.Var2 = var2
    

    return var1,var2
}
 
// WriterTo implements the MofInitable.WriterTo() method for the MofInitable.
func (impl *MofInitableLittleSnitch) WriterTo() (int64,error){
    var caller MethodCallForWriterTo

    defer func(){
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

    

    var3,var4 := impl.Implementer.WriterTo()

    
    caller.Var3 = var3
    
    caller.Var4 = var4
    

    return var3,var4
}


//==============================================================================================================

// MofInitableMockSnitch defines a function type which implements a struct with the 
// methods for the MofInitable as fields which allows you provide implementations of 
// these functions to provide flexible testing.
type MofInitableMockSnitch struct{
    
    IgniteMethodCalls []MethodCallForIgnite
    IgniteFunc func() (string)
    
    CrunchMethodCalls []MethodCallForCrunch
    CrunchFunc func() (string)
    
    ConfigurationMethodCalls []MethodCallForConfiguration
    ConfigurationFunc func() ()
    
    LocationMethodCalls []MethodCallForLocation
    LocationFunc func(var1 string) (GPSLoc,error)
    
    WriterToMethodCalls []MethodCallForWriterTo
    WriterToFunc func() (int64,error)
    
}

 
// Ignite implements the MofInitable.Ignite() method for the MofInitable.
func (impl *MofInitableMockSnitch) Ignite() (string){
    var caller MethodCallForIgnite

    defer func(){
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
func (impl *MofInitableMockSnitch) Crunch() (string){
    var caller MethodCallForCrunch

    defer func(){
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
func (impl *MofInitableMockSnitch) Configuration() (){
    var caller MethodCallForConfiguration

    defer func(){
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

    

     := impl.ConfigurationFunc()

    

    return 
}
 
// Location implements the MofInitable.Location() method for the MofInitable.
func (impl *MofInitableMockSnitch) Location(var1 string) (GPSLoc,error){
    var caller MethodCallForLocation

    defer func(){
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
    

    var1,var2 := impl.LocationFunc(var1)

    
    caller.Var1 = var1
    
    caller.Var2 = var2
    

    return var1,var2
}
 
// WriterTo implements the MofInitable.WriterTo() method for the MofInitable.
func (impl *MofInitableMockSnitch) WriterTo() (int64,error){
    var caller MethodCallForWriterTo

    defer func(){
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

    

    var3,var4 := impl.WriterToFunc()

    
    caller.Var3 = var3
    
    caller.Var4 = var4
    

    return var3,var4
}
