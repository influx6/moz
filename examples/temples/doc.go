// Package temples defines a series of structures which are auto-generated based on
// a template and series of type declerations.
//
// @templater(id => Mob, gen => Go, {
//
//  // Add brings the new level into the system.
//  func Add(m TYPE1, n TYPE2) TYPE3 {
//
//  }
//
// })
// @templaterTypesFor(id => Mob, TYPE1 => int32, TYPE2 => int32, TYPE3 => int64)
// @templaterTypesFor(id => Mob, TYPE1 => int, TYPE2 => int, TYPE3 => int64)
//
package temples
