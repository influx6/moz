// Package dap defines a sample package which is used to test out our annotation parser.
// @assets(dapassets, ".tml : .js", dappers)
package dap

import (
	"fmt"
	"strconv"
)

//go:generate moz annotation

// IgnitionFn defines a functiion type for a ignition function.
type IgnitionFn func(string) string

// Ignitable defines a struct which is used to ignite the package.
type Ignitable interface {
	Ignite() string
}

// IgnitionDescription defines the description giving to a ignition key.
type IgnitionDescription string

// Ignitor defines a struct which is used to ignite the package.
//@httpapi
type Ignitor struct {
	PublicID string `json:"public_id" toml:"public_id"`
	Name     string `json:"name" toml:"name"`
	Repo     Repo   `json:"repo" toml:"repo"`
}

// Ignite returns the ignited string related to the struct.
func (i Ignitor) Ignite() string {
	return fmt.Sprintf("%s@%s", i.Name, strconv.Quote(i.Repo.URL))
}

// Repo defines a struct which defines a object pointing to a specific repo.
type Repo struct {
	URL string `json:"url"`
}
