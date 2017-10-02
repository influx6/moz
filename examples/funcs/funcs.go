// +build maze

package funcs

import (
	"github.com/influx6/moz/examples/funcs/balls"
)

// DoRoll does something with balls.Roll.
func DoRoll() {
	return balls.Roll("jerk")
}

// RunDown does something.
func RunDown() string {
	return "running..."
}
