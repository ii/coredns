// The log package implements a small wrapper around the std lib log package.
// All logging can be done as-is, but in addition debug logging can be done
// as well.
//
// These last ones only show output when the *debug* plugin is loaded.
package log

import (
	"fmt"
	golog "log"
)

// When D is true, debugging output will be shown.
var D bool

// Debug is equivalent to log.Print() but prefixed with "[DEBUG] ".
func Debug(v ...interface{}) {
	if !D {
		return
	}
	s := debug + fmt.Sprint(v...)
	golog.Print(s)
}

// Debug is equivalent to log.Printf() but prefixed with "[DEBUG] ".
func Debugf(format string, v ...interface{}) {
	if !D {
		return
	}
	s := debug + fmt.Sprintf(format, v...)
	golog.Print(s)
}

// Debugln is equivalent to log.Println() but prefixed with "[DEBUG] ".
func Debugln(v ...interface{}) {
	if !D {
		return
	}
	s := debug + fmt.Sprintln(v...)
	golog.Print(s)
}

// Printf calls log.Printf.
func Printf(format string, v ...interface{}) { golog.Printf(format, v...) }

// Print calls log.Print.
func Print(v ...interface{}) { golog.Print(v...) }

// Println calls log.Println.
func Println(v ...interface{}) { golog.Println(v...) }

const debug = "[DEBUG] "
