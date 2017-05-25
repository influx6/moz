// Package dap defines a sample package which is used to test out our annotation parser.
//@plugin(./annotations/buga.sol)
//@plugin(./annotations/fude.sol)
//@plugin(./annotations/nick.sol)
package dap

import (
	"fmt"
	"strconv"
)

//go:generate moz

// IgnitionFn defines a functiion type for a ignition function.
//@homa
type IgnitionFn func(string) string

// Ignitable defines a struct which is used to ignite the package.
//@buma
type Ignitable interface {
	Ignite() string
}

// IgnitionDescription defines the description giving to a ignition key.
//@homa
type IgnitionDescription string

// Ignitor defines a struct which is used to ignite the package.
//@homa
//@mongo
type Ignitor struct {
	Name string `json:"name"`
	Repo Repo   `json:"repo"`
}

// Ignite returns the ignited string related to the struct.
func (i Ignitor) Ignite() string {
	return fmt.Sprintf("%s@%s", i.Name, strconv.Quote(i.Repo.URL))
}

// Repo defines a struct which defines a object pointing to a specific repo.
//@oauth
type Repo struct {
	URL string `json:"url"`
}
