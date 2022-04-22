//go:build freebsd || dragonfly || openbsd || netbsd

package sysinfo

func detect() Info {
	return Info{}
}
