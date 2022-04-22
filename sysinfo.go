package sysinfo

import (
	"bufio"
	"bytes"
	"encoding"
	"fmt"
	"math/bits"
	"sort"
	"strconv"
	"strings"
)

type Info struct {
	// CPUs is per-cpu information.
	//
	// CPUs is sorted by the Proc field in asending order.
	CPUs []CPU
	// Misc is any unknown information.
	//
	// Misc is sorted by the Key field in asending order.
	Misc []Pair
}

// Detect finds the current host information.
//
// Note that each call to Detect might return different
// information.
func Detect() Info {
	return detect()
}

// CPU describes a single CPU.
//
// On Linux, this information is read from /proc/cpuinfo.
// On BSDs (including macOS), this information is read from
// sysctl.
//
// Each read from /proc/cpuinfo has a "Matches:" comment
// describing the key used.
type CPU struct {
	// General

	// Proc is the processor number, usually zero-indexed.
	//
	// Matches: processor
	Proc int `json:"processor"`
	// BogoMIPS is a Linux-specific rough measurement of CPU
	// performance.
	//
	// Matches: BogoMIPS, bogomips
	BogoMIPS float64 `json:"bogomips,omitempty"`
	// Features is the set of supported CPU features or flags.
	//
	// Matches: Features, flags
	Features []string `json:"features,omitempty"`
	// Rev is the CPU "stepping" or revision.
	//
	// Matches: CPU revision, stepping
	Rev int `json:"revision,omitempty"`
	// ModelName is the human-readable CPU model name.
	//
	// Matches: model name
	ModelName string `json:"model_name,omitempty"`
	// MicroArch is the CPU's microarchitecture.
	MicroArch string `json:"micro_arch"`

	// ARM

	// Impl is the CPU implementer.
	//
	// Matches: CPU implementer
	Impl Implementer `json:"implementer,omitempty"`
	// Arch is the CPU architecture.
	//
	// Matches: CPU architecture
	Arch int `json:"arch,omitempty"`
	// Variant is the CPU variant.
	//
	// Matches: CPU variant
	Variant int `json:"variant,omitempty"`
	// Part identifies the specific CPU.
	//
	// Matches: CPU part
	Part Part `json:"part_number,omitempty"`

	// Intel/AMD (x86)

	// VendorID identifies the CPU vendor.
	//
	// Typically is "GenuineIntel" for Intel CPUs and
	// "AuthenticAMD" for AMD CPUs.
	//
	// Matches: vendor_id
	VendorID string `json:"vendor_id,omitempty"`
	// Family is the CPU family.
	//
	// Matches: cpu family
	Family int `json:"family,omitempty"`
	// Model is the CPU model number.
	//
	// Matches: model
	Model int `json:"model_number,omitempty"`
	// Microcode is the CPU's microcode version.
	//
	// Matches: microcode
	Microcode int `json:"microcode_version,omitempty"`
	// Freq is the CPU frequency in MHz.
	//
	// Matches: cpu MHz
	Freq float64 `json:"frequency_mhz,omitempty"`
	// Cache is the CPU's cache information.
	Cache Cache `json:"cache,omitempty"`
	// PhysID is the physical ID of this CPU core.
	//
	// Matches: physical id
	PhysID int `json:"physical_id,omitempty"`
	// Siblings is the number of sibling CPUs.
	//
	// Matches: siblings
	Siblings int `json:"siblings,omitempty"`
	// CoreID is the unique ID of this CPU core.
	//
	// Matches: core id
	CoreID int `json:"core_id,omitempty"`
	// Cores is the number of CPU cores.
	//
	// Matches: cores
	Cores int `json:"num_cores,omitempty"`
	// APICID is the APIC system's unique ID.
	//
	// Matches: apicid
	APICID int `json:"apic_id,omitempty"`
	// InitAPICID is the APIC system's unique ID as assigned at
	// startup.
	//
	// Matches: initial apicid
	InitAPICID int `json:"initial_apic_id,omitempty"`
	// FPU is whether the CPU has floating-point unit.
	//
	// Matches: fpu
	FPU bool `json:"fpu,omitempty"`
	// FPUExceptions is whether the CPU supports floating-point
	// unit exceptions.
	//
	// Matches: fpu_exceptions
	FPUExceptions bool `json:"fpu_exceptions,omitempty"`
	// CPUID level is the maximum CPUID level that can be used
	// when querying the CPU for information via the CPUID
	// instruction.
	//
	// Matches: cpuid level
	CPUIDLevel int `json:"cpuid_level,omitempty"`
	// WP is whether the CPU supports write protection.
	//
	// Matches: wp
	WP bool `json:"write_protection,omitempty"`
	// Bugs the set of bugs that have been detected or worked
	// around.
	//
	// Matches: bugs
	Bugs []string `json:"bugs,omitempty"`
	// AddrSizes are the CPU's memory address sizes.
	//
	// Matches: address sizes
	AddrSizes struct {
		// Phys is the number of bits in a physical memory
		// address.
		Phys int `json:"physical_bits,omitempty"`
		// Virt is the number of bits in a virtual memory
		// address.
		Virt int `json:"virtual_bits,omitempty"`
	} `json:"address_sizes,omitempty"`
	// PowerMgmt is the supported power management features.
	//
	// Matches: power management.
	PowerMgmt string `json:"power_management,omitempty"`

	// AMD

	// TLB is the CPU's Translation Lookaside Buffer.
	//
	// Matches: TLB size
	TLB struct {
		// N is the number of TLB pages.
		N int `json:"num_pages,omitempty"`
		// PageSize is the size in bytes of each page.
		PageSize int `json:"page_size,omitempty"`
	} `json:"tlb,omitempty"`
}

type Cache struct {
	// Inst is the size in bytes of the CPU's instruction
	// cache.
	Inst int `json:"instruction,omitempty"`
	// L1Data is the size in bytes of the CPU's L1 data
	// cache.
	L1 int `json:"l1,omitempty"`
	// L2 is the size in bytes of the CPU's L2 cache.
	//
	// Matches: cache size
	L2 int `json:"l2,omitempty"`
	// L3 is the size in bytes of the CPU's L3 cache.
	L3 int `json:"l3,omitempty"`
	// Alignment is how the CPU caches are aligned.
	//
	// Matches: cache_alignment
	Alignment int `json:"alignment,omitempty"`
	// Flush is the size of a cache line flush (CLFLUSH).
	//
	// Matches: clflush size
	Flush int `json:"flush,omitempty"`
}

// Pair is a miscellaneous piece of data reported by the host.
type Pair struct {
	Key, Value string
}

func (c CPU) String() string {
	return fmt.Sprintf("%s %s", c.Impl, c.Name())
}

func (c CPU) Name() string {
	switch c.Impl {
	case ARMLtd:
		return armPartName(c.Part)
	case Broadcom, Cavium:
		return broadcomPartName(c.Part)
	case Fujitsu:
		return fujitsuPartName(c.Part)
	case NVIDIA:
		return nvidiaPartName(c.Part)
	case HiSilicon:
		return hiSiliconPartName(c.Part)
	case Qualcomm:
		return qualcommPartName(c.Part)
	case Samsung:
		return samsungPartName(c.Part)
	default:
		return "generic"
	}
}

const (
	ARMLtd    Implementer = 'A' // ARM Ltd
	Broadcom  Implementer = 'B' // Broadcom
	Cavium    Implementer = 'C' // Cavium
	Fujitsu   Implementer = 'F' // Fujitsu Ltd
	NVIDIA    Implementer = 'N' // NVIDIA Corporation
	HiSilicon Implementer = 'H' // HiSilicon Technologies Inc
	Qualcomm  Implementer = 'Q' // Qualcomm Technologies Inc
	Samsung   Implementer = 'S' // Samsung Technologies Inc
	Intel     Implementer = 'i' // Intel ARM parts
	Apple     Implementer = 'a' // Apple Inc
)

type Implementer uint8

var _ encoding.TextMarshaler = Implementer(0)

func (i Implementer) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// ARM
const (
	ARM926EJS   Part = 0x926 // ARM926EJ-S
	ARM11MPCore Part = 0xb02 // ARM11-MPCore
	ARM1136JS   Part = 0xb36 // ARM1136J-S
	ARM1156T2S  Part = 0xb56 // ARM1156T2-S
	ARM1176JZS  Part = 0xb76 // ARM1176JZ-S
	CortexA8    Part = 0xc08 // Cortex-A8
	CortexA9    Part = 0xc09 // Cortex-A9
	CortexA15   Part = 0xc0f // Cortex-A15
	CortexM0    Part = 0xc20 // Cortex-M0
	CortexM3    Part = 0xc23 // Cortex-M3
	CortexM4    Part = 0xc24 // Cortex-M4
	CortexM55   Part = 0xd22 // Cortex-M55
	CortexA34   Part = 0xd02 // Cortex-A34
	CortexA35   Part = 0xd04 // Cortex-A35
	CortexA53   Part = 0xd03 // Cortex-A53
	CortexA55   Part = 0xd05 // Cortex-A55
	CortexA57   Part = 0xd07 // Cortex-A57
	CortexA72   Part = 0xd08 // Cortex-A72
	CortexA73   Part = 0xd09 // Cortex-A73
	CortexA75   Part = 0xd0a // Cortex-A75
	CortexA76   Part = 0xd0b // Cortex-A76
	CortexA77   Part = 0xd0d // Cortex-A77
	CortexA78   Part = 0xd41 // Cortex-A78
	CortexX1    Part = 0xd44 // Cortex-X1
	CortexX1C   Part = 0xd4c // Cortex-X1C
	NeoverseN1  Part = 0xd0c // neoverse N1
	NeoverseN2  Part = 0xd49 // Neoverse N2
	NeoverseV1  Part = 0xd40 // Neoverse V1
	Firestorm   Part = 0x23  // M1 Firestorm
	Icestorm    Part = 0x22  // M1 Icestorm
)

type Part uint16

func armPartName(p Part) string {
	switch p {
	case ARM926EJS:
		return "ARM926EJ-S"
	case ARM11MPCore:
		return "ARM11 MPCore"
	case ARM1136JS:
		return "ARM1136J-S"
	case ARM1156T2S:
		return "ARM1156T2-S"
	case ARM1176JZS:
		return "ARM1176JZ-S"
	case CortexA8:
		return "Cortex-A8"
	case CortexA9:
		return "Cortex-A9"
	case CortexA15:
		return "Cortex-A15"
	case CortexM0:
		return "Cortex-M0"
	case CortexM3:
		return "Cortex-M3"
	case CortexM4:
		return "Cortex-M4"
	case CortexM55:
		return "Cortex-M55"
	case CortexA34:
		return "Cortex-A34"
	case CortexA35:
		return "Cortex-A35"
	case CortexA53:
		return "Cortex-A53"
	case CortexA55:
		return "Cortex-A55"
	case CortexA57:
		return "Cortex-A57"
	case CortexA72:
		return "Cortex-A72"
	case CortexA73:
		return "Cortex-A73"
	case CortexA75:
		return "Cortex-A75"
	case CortexA76:
		return "Cortex-A76"
	case CortexA77:
		return "Cortex-A77"
	case CortexA78:
		return "Cortex-A78"
	case CortexX1:
		return "Cortex-X1"
	case CortexX1C:
		return "Cortex-X1C"
	case NeoverseN1:
		return "neoverse-n1"
	case NeoverseN2:
		return "neoverse-n2"
	case NeoverseV1:
		return "neoverse-v1"
	case Firestorm:
		return "M1 Firestorm"
	case Icestorm:
		return "M1 Icestorm"
	default:
		return "generic"
	}
}

// Broadcom/Cavium
const (
	ThunderX2T99   Part = 0x516 // thunderx2t99
	ThunderX2T99_2 Part = 0xaf  // thunderx2t99
	ThunderXT88    Part = 0xa1  // thunderxt88
)

func broadcomPartName(p Part) string {
	switch p {
	case ThunderX2T99, ThunderX2T99_2:
		return "ThunderX2T99"
	case ThunderXT88:
		return "ThunderXT88"
	default:
		return "generic"
	}
}

// Fujitsu
const (
	A64FX Part = 0x001 // a64fx
)

func fujitsuPartName(p Part) string {
	switch p {
	case A64FX:
		return "A64FX"
	default:
		return "generic"
	}
}

// NVIDIA
const (
	Carmel Part = 0x004 // carmel
)

func nvidiaPartName(p Part) string {
	switch p {
	case Carmel:
		return "Carmel"
	default:
		return "generic"
	}
}

// HiSilicon
const (
	TSV110 Part = 0xd01 // tsv110
)

func hiSiliconPartName(p Part) string {
	switch p {
	case TSV110:
		return "TSV110"
	default:
		return "generic"
	}
}

const (
	Krait         Part = 0x06f // krait
	Kryo          Part = 0x201 // kryo
	Kryo_2        Part = 0x205 // kryo
	Kryo_3        Part = 0x211 // kryo
	Kryo2xxGold   Part = 0x800 // cortex-a73
	Kryo2xxSilver Part = 0x801 // cortex-a73
	Kryo3xxGold   Part = 0x802 // cortex-a75
	Kryo3xxSilver Part = 0x803 // cortex-a75
	Kryo4xxGold   Part = 0x804 // cortex-a76
	Kryo4xxSilver Part = 0x805 // cortex-a76
	Falkor        Part = 0xc00 // falkor
	Saphira       Part = 0xc01 // saphira
)

func qualcommPartName(p Part) string {
	switch p {
	case Krait:
		return "Krait"
	case Kryo, Kryo_2, Kryo_3:
		return "Kryo"
	case Kryo2xxGold, Kryo2xxSilver:
		return "Cortex-A73"
	case Kryo3xxGold, Kryo3xxSilver:
		return "Cortex-A75"
	case Kryo4xxGold, Kryo4xxSilver:
		return "Cortex-A76"
	case Falkor:
		return "Falkor"
	case Saphira:
		return "Saphira"
	default:
		return "generic"
	}
}

func samsungPartName(p Part) string {
	switch p {
	default:
		return "generic"
	}
}

// scanProc parses the output /proc/cpuinfo.
//
// It should look like
//
//    processor	: 5
//    BogoMIPS	: 48.00
//    Features	: fp asimd evtstrm aes pmull sha1 sha2 crc32 cpuid
//    CPU implementer	: 0x41
//    CPU architecture: 8
//    CPU variant	: 0x0
//    CPU part	: 0xd08
//    CPU revision	: 2
//
// See http://www.linfo.org/proc_cpuinfo.html
func scanProc(o *Info, buf []byte) {
	var c CPU
	s := bufio.NewScanner(bytes.NewReader(buf))
	for s.Scan() {
		k, v := split(s.Text())
		switch k {
		case "":
			o.CPUs = append(o.CPUs, c)
		case "processor":
			c.Proc = atoi(v)
		case "BogoMIPS", "bogomips":
			c.BogoMIPS = atof(v)
		case "Features":
			c.Features = strings.Split(v, " ")
		case "CPU implementer":
			c.Impl = Implementer(atoi(v))
		case "CPU architecture":
			c.Arch = atoi(v)
		case "CPU variant":
			c.Variant = atoi(v)
		case "CPU part":
			c.Part = Part(atoi(v))
		case "CPU revision", "stepping":
			c.Rev = atoi(v)
		case "model name":
			c.ModelName = v
		case "vendor_id":
			c.VendorID = v
		case "cpu family":
			c.Family = atoi(v)
		case "model":
			c.Model = atoi(v)
		case "microcode":
			c.Microcode = atoi(v)
		case "cpu MHz":
			c.Freq = atof(v)
		case "cache size":
			c.Cache.L2 = parseSize(v)
		case "physical id":
			c.PhysID = atoi(v)
		case "siblings":
			c.Siblings = atoi(v)
		case "core id":
			c.CoreID = atoi(v)
		case "cpu cores":
			c.Cores = atoi(v)
		case "apicid":
			c.APICID = atoi(v)
		case "initial apicid":
			c.InitAPICID = atoi(v)
		case "fpu":
			c.FPU = parseBool(v)
		case "fpu_exception":
			c.FPUExceptions = parseBool(v)
		case "cpuid level":
			c.CPUIDLevel = atoi(v)
		case "wp":
			c.WP = parseBool(v)
		case "flags":
			c.Features = strings.Split(v, " ")
		case "bugs":
			c.Bugs = strings.Split(v, " ")
		case "clflush size":
			c.Cache.Flush = atoi(v)
		case "cache_alignment":
			c.Cache.Alignment = atoi(v)
		case "address sizes":
			c.AddrSizes.Phys, c.AddrSizes.Virt = parseAddrSizes(v)
		case "power management":
			c.PowerMgmt = v
		case "TLB size":
			c.TLB.N, c.TLB.PageSize = parseTLB(v)
		default:
			o.Misc = append(o.Misc, Pair{Key: k, Value: v})
		}
	}
	sort.Slice(o.CPUs, func(i, j int) bool {
		return o.CPUs[i].Proc < o.CPUs[j].Proc
	})
	sort.Slice(o.Misc, func(i, j int) bool {
		return o.Misc[i].Key < o.Misc[j].Key
	})
}

func parseBool(s string) bool {
	return s == "yes"
}

func parseSize(s string) int {
	i := strings.IndexByte(s, ' ')
	if i < 0 || i == len(s) {
		return 0
	}
	x := atoi(s[:i])
	switch unit := s[i+1:]; unit {
	case "KB":
		x *= 1024
	case "MB":
		x *= 1024 * 1024
	default:
		return 0
	}
	return x
}

// addrSizes parses an "address sizes" string with the format
//
//    40 bits physical, 48 bits virtual
//
func parseAddrSizes(s string) (phys, virt int) {
	const (
		sep = ", "
	)
	i := strings.Index(s, sep)
	if i < 0 || i+len(sep) == len(s) {
		return 0, 0
	}
	lhs, rhs := s[:i], s[i+len(sep):]

	if !strings.HasSuffix(lhs, " bits physical") {
		return 0, 0
	}
	lhs = strings.TrimSuffix(lhs, " bits physical")

	if !strings.HasSuffix(rhs, " bits virtual") {
		return 0, 0
	}
	rhs = strings.TrimSuffix(rhs, " bits virtual")

	phys = atoi(lhs)
	virt = atoi(rhs)
	return phys, virt
}

// parseTLB parses a "TLB size" string with the format
//
//    1024 4K pages
//
func parseTLB(s string) (n, page int) {
	if !strings.HasSuffix(s, " pages") {
		return 0, 0
	}
	s = strings.TrimSuffix(s, " pages")

	var unit int
	switch {
	case strings.HasSuffix(s, "K"):
		unit = 1024
	case strings.HasSuffix(s, "M"):
		unit = 1024 * 1024
	default:
		return 0, 0
	}
	s = s[:len(s)-1]

	i := strings.IndexByte(s, ' ')
	if i < 0 || i == len(s) {
		return 0, 0
	}
	n = atoi(s[:i])
	page = atoi(s[i+1:]) * unit
	return n, page
}

func atof(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func atoi(s string) int {
	x, _ := strconv.ParseUint(s, 0, bits.UintSize)
	return int(x)
}

func split(s string) (key, value string) {
	i := strings.IndexByte(s, ':')
	if i >= 0 {
		key = strings.TrimSpace(s[:i])
	}
	if i+1 < len(s) {
		value = strings.TrimSpace(s[i+1:])
	}
	return
}
