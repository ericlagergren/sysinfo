package sysinfo

import "os"

func detect() Info {
	buf, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return Info{}
	}
	var v Info
	scanProc(&v, buf)
	return v
}
