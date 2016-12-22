// Package typetree provides an API for operating on arbitrary types
// and their children.
//
// Allow public struct objects in this library MUST
// allow instantiation via 'name{}', var n name, n := new(name).
//
// This will be called the "zero value principle" and noted in
// the support code.
package typetree
