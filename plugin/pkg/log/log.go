// Package log implements a small wrapper around the std lib log package.
// All logging can be done as-is, but in addition debug logging is available.
//
// I.e. normal use: log.Print("[INFO] this is some logging").
//
// Debug logging: log.Debug("this is debug output"). Note "[DEBUG] " is prepended
// to the output.
//
// These last one only show output when the *debug* plugin is loaded.
package log

import (
	"fmt"
	golog "log"
)

// D controls whether we should ouput debug logs. If true, we do.
var D bool

// Debug is equivalent to log.Print() but prefixed with "[DEBUG] ".
func Debug(v ...interface{}) {
	if !D {
		return
	}
	s := debug + fmt.Sprint(v...)
	golog.Print(s)
}

// Debugf is equivalent to log.Printf() but prefixed with "[DEBUG] ".
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
