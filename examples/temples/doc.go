// Package temples defines a series of structures which are auto-generated based on
// a template and series of type declerations.
//
// @templater(id => Mob, gen => Partial.Go, {
//
//  // Add brings the new level into the system.
//  func Add(m {{ sel "Type1"}}, n {{ sel "Type2"}}) {{ sel "Type3" }} {
//      return {{sel "Type3"}}(m * n)
//  }
//
// })
//
// @templaterTypesFor(id => Mob, filename => temples_add.go, Type1 => int32, Type2 => int32, Type3 => int64)
//
package temples
