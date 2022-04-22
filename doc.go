// Package sysinfo provides host detection information.
//
// It is intended to be complimentary to
// github.com/klauspost/cpuid and golang.org/x/sys/cpu.
package sysinfo

//go:generate go run golang.org/x/tools/cmd/stringer -type Implementer -linecomment
