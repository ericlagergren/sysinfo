package sysinfo

import (
	"encoding/binary"
	"fmt"

	"golang.org/x/sys/unix"
)

// system_profiler SPHardwareDataType

const (
	famFireIce = 0x1b588bb3
)

func detect() Info {
	v := Info{
		Misc: []Pair{
			{"Kernel Version", sysctl("kern.version")},
			{"OS Version", sysctl("kern.osversion")},
		},
	}
	switch sysctl32("hw.cpufamily") {
	case famFireIce:
		detectM1(&v)
	}
	return v
}

func detectM1(o *Info) {
	brand := sysctl("machdep.cpu.brand_string") // Apple M1
	switch sysctl32("hw.cpusubfamily") {
	case 2:
		// OK
	case 4:
		brand += " Pro"
	case 5:
		brand += " Max"
	}

	vaddr := sysctl32("machdep.virtual_address_size")

	// TODO(eric): update golang.org/x/sys/cpu to 128
	align := sysctl64("hw.cachelinesize")
	lvls := int(sysctl32("hw.nperflevels"))

	for lvl := 0; lvl < lvls; lvl++ {
		cache := Cache{
			Inst:      int(sysctl32(fmt.Sprintf("hw.perflevel%d.l1icachesize", lvl))),
			L1:        int(sysctl32(fmt.Sprintf("hw.perflevel%d.l1dcachesize", lvl))),
			L2:        int(sysctl32(fmt.Sprintf("hw.perflevel%d.l2cachesize", lvl))),
			Alignment: int(align),
		}
		cores := int(sysctl32(fmt.Sprintf("hw.perflevel%d.physicalcpu", lvl)))
		for i := 0; i < cores; i++ {
			c := CPU{
				Proc:      len(o.CPUs),
				Impl:      Apple,
				Model:     famFireIce,
				ModelName: brand,
				Cache:     cache,
				Arch:      8,
			}
			if lvl == 0 {
				c.MicroArch = "Firestorm"
			} else {
				c.MicroArch = "Icestorm"
			}
			c.AddrSizes.Virt = int(vaddr)
			o.CPUs = append(o.CPUs, c)
		}
	}

	o.Misc = append(o.Misc,
		Pair{Key: "Model", Value: sysctl("hw.model")},
	)
}

const debug = false

func sysctl(s string) string {
	v, err := unix.Sysctl(s)
	if debug && err != nil {
		fmt.Printf("%q: %v\n", s, err)
	}
	return v
}

func sysctl32(s string) uint32 {
	v, err := unix.SysctlUint32(s)
	if debug && err != nil {
		fmt.Printf("%q: %v\n", s, err)
	}
	return v
}

func sysctl64(s string) uint64 {
	v, err := unix.SysctlUint64(s)
	if debug && err != nil {
		fmt.Printf("%q: %v\n", s, err)
	}
	return v
}

func cacheconfig() []uint64 {
	v, _ := unix.SysctlRaw("hw.cacheconfig")
	s := make([]uint64, len(v)/8)
	for i := range s {
		s[i] = binary.LittleEndian.Uint64(v[i*8:])
	}
	return s
}
