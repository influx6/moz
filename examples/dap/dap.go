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
//@mongoapi
//@associates(@mongoapi, New, Ignitor)
//@associates(@httpapi, New, Ignitor)
//@associates(@mongoapi, Update, Ignitor)
//@associates(@httpapi, Update, Ignitor)
type Ignitor struct {
	Identity
	PublicID string `json:"public_id,omitempty" toml:"public_id"`
	Name     string `json:"name" toml:"name"`
	Rex      Repo   `json:"rex" toml:"rex"`
	Rack     int    `json:"rack"`
	version  string
}

// Repo defines a struct which defines a object pointing to a specific repo.
type Repo struct {
	URL string `json:"url"`
}

type Identity struct {
	Hash string `json:"hash"`
}

// UpdateIgnitor defines a struct for creating a Ignitor.
type UpdateIgnitor struct {
	ID   string
	Repo string
}

// NewIgnitor defines a struct for creating a Ignitor.
type NewIgnitor struct {
	ID   string
	Repo string
}

// Ignite returns the ignited string related to the struct.
func (i Ignitor) Ignite(in string) string {
	return fmt.Sprintf("%s@%s", i.Name, strconv.Quote(i.Rex.URL))
}

// Build run.
func Build() error {
	return nil
}
