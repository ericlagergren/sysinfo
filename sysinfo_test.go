package sysinfo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/r3labs/diff/v3"
)

func TestX(t *testing.T) {
	v := Detect()
	println(sprint(v))
}

func TestReadProc(t *testing.T) {
	split := func(s string) []string {
		return strings.Split(s, " ")
	}
	p6feats := split("fp asimd evtstrm aes pmull sha1 sha2 crc32 atomics fphp asimdhp cpuid asimdrdm lrcpc dcpop asimddp")
	rp64feats := split("fp asimd evtstrm aes pmull sha1 sha2 crc32 cpuid")
	rpifeats := split("half thumb fastmult vfp edsp neon vfpv3 tls vfpv4 idiva idivt vfpd32 lpae evtstrm crc32")

	for _, tc := range []struct {
		name string
		info Info
	}{
		{
			name: "google_pixel_6",
			info: Info{CPUs: []CPU{
				{Proc: 0, BogoMIPS: 49.15, Features: p6feats, Impl: ARMLtd, Arch: 8, Variant: 0x2, Part: CortexA55},
				{Proc: 1, BogoMIPS: 49.15, Features: p6feats, Impl: ARMLtd, Arch: 8, Variant: 0x2, Part: CortexA55},
				{Proc: 2, BogoMIPS: 49.15, Features: p6feats, Impl: ARMLtd, Arch: 8, Variant: 0x2, Part: CortexA55},
				{Proc: 3, BogoMIPS: 49.15, Features: p6feats, Impl: ARMLtd, Arch: 8, Variant: 0x2, Part: CortexA55},
				{Proc: 4, BogoMIPS: 49.15, Features: p6feats, Impl: ARMLtd, Arch: 8, Variant: 0x4, Part: CortexA76},
				{Proc: 5, BogoMIPS: 49.15, Features: p6feats, Impl: ARMLtd, Arch: 8, Variant: 0x4, Part: CortexA76},
				{Proc: 6, BogoMIPS: 49.15, Features: p6feats, Impl: ARMLtd, Arch: 8, Variant: 0x1, Part: CortexX1},
				{Proc: 7, BogoMIPS: 49.15, Features: p6feats, Impl: ARMLtd, Arch: 8, Variant: 0x1, Part: CortexX1},
			}},
		},
		{
			name: "rockpro64",
			info: Info{CPUs: []CPU{
				{Proc: 0, BogoMIPS: 48, Features: rp64feats, Impl: ARMLtd, Arch: 8, Part: CortexA53, Rev: 4},
				{Proc: 1, BogoMIPS: 48, Features: rp64feats, Impl: ARMLtd, Arch: 8, Part: CortexA53, Rev: 4},
				{Proc: 2, BogoMIPS: 48, Features: rp64feats, Impl: ARMLtd, Arch: 8, Part: CortexA53, Rev: 4},
				{Proc: 3, BogoMIPS: 48, Features: rp64feats, Impl: ARMLtd, Arch: 8, Part: CortexA53, Rev: 4},
				{Proc: 4, BogoMIPS: 48, Features: rp64feats, Impl: ARMLtd, Arch: 8, Part: CortexA72, Rev: 2},
				{Proc: 5, BogoMIPS: 48, Features: rp64feats, Impl: ARMLtd, Arch: 8, Part: CortexA72, Rev: 2},
			}},
		},
		{
			name: "raspberry_pi_4b",
			info: Info{
				CPUs: []CPU{
					{Proc: 0, BogoMIPS: 108, Features: rpifeats, Impl: ARMLtd, Arch: 7, Part: CortexA72, Rev: 3, ModelName: "ARMv7 Processor rev 3 (v7l)"},
					{Proc: 1, BogoMIPS: 108, Features: rpifeats, Impl: ARMLtd, Arch: 7, Part: CortexA72, Rev: 3, ModelName: "ARMv7 Processor rev 3 (v7l)"},
					{Proc: 2, BogoMIPS: 108, Features: rpifeats, Impl: ARMLtd, Arch: 7, Part: CortexA72, Rev: 3, ModelName: "ARMv7 Processor rev 3 (v7l)"},
					{Proc: 3, BogoMIPS: 108, Features: rpifeats, Impl: ARMLtd, Arch: 7, Part: CortexA72, Rev: 3, ModelName: "ARMv7 Processor rev 3 (v7l)"},
				},
				Misc: []Pair{
					{"Hardware", "BCM2711"},
					{"Model", "Raspberry Pi 4 Model B Rev 1.1"},
					{"Revision", "c03111"},
					{"Serial", "10000000771c4af4"},
				},
			},
		},
		{
			name: "intel_skylake_ubuntu",
			info: Info{
				CPUs: []CPU{
					{
						VendorID:      "GenuineIntel",
						Family:        6,
						Model:         94,
						ModelName:     "Intel Core Processor (Skylake, IBRS)",
						Rev:           3,
						Microcode:     0x1,
						Freq:          3791.976,
						Cache:         Cache{0, 0, 16384 * 1024, 0, 64, 64},
						Siblings:      1,
						Cores:         1,
						FPU:           true,
						FPUExceptions: true,
						CPUIDLevel:    13,
						WP:            true,
						Features:      split("fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 syscall nx rdtscp lm constant_tsc rep_good nopl xtopology cpuid tsc_known_freq pni pclmulqdq ssse3 fma cx16 pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand hypervisor lahf_lm abm cpuid_fault invpcid_single pti ssbd ibrs ibpb fsgsbase bmi1 avx2 smep bmi2 erms invpcid xsaveopt arat"),
						Bugs:          split("cpu_meltdown spectre_v1 spectre_v2 spec_store_bypass l1tf mds swapgs itlb_multihit srbds"),
						BogoMIPS:      7583.95,
						AddrSizes: struct {
							Phys int `json:"physical_bits,omitempty"`
							Virt int `json:"virtual_bits,omitempty"`
						}{40, 48},
					},
				},
			},
		},
		{
			name: "intel_cascadelake_ubuntu",
			info: Info{
				CPUs: []CPU{
					{
						VendorID:      "GenuineIntel",
						Family:        6,
						Model:         85,
						ModelName:     "Intel Xeon Processor (Cascadelake)",
						Rev:           6,
						Microcode:     0x1,
						Freq:          2992.968,
						Cache:         Cache{0, 0, 16384 * 1024, 0, 64, 64},
						Siblings:      1,
						Cores:         1,
						FPU:           true,
						FPUExceptions: true,
						CPUIDLevel:    13,
						WP:            true,
						Features:      split("fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 syscall nx pdpe1gb rdtscp lm constant_tsc rep_good nopl xtopology cpuid tsc_known_freq pni pclmulqdq ssse3 fma cx16 pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand hypervisor lahf_lm abm cpuid_fault invpcid_single pti ssbd ibrs ibpb fsgsbase bmi1 avx2 smep bmi2 erms invpcid avx512f avx512dq clflushopt clwb avx512cd avx512bw avx512vl xsaveopt arat pku ospke avx512_vnni"),
						Bugs:          split("cpu_meltdown spectre_v1 spectre_v2 spec_store_bypass l1tf mds swapgs itlb_multihit"),
						BogoMIPS:      5985.93,
						AddrSizes: struct {
							Phys int `json:"physical_bits,omitempty"`
							Virt int `json:"virtual_bits,omitempty"`
						}{40, 48},
					},
				},
			},
		},
		{
			name: "amd_epyc_centos",
			info: Info{
				CPUs: []CPU{
					{
						VendorID:      "AuthenticAMD",
						Family:        23,
						Model:         1,
						ModelName:     "AMD EPYC 7551 32-Core Processor",
						Rev:           2,
						Microcode:     0x1000065,
						Freq:          1996.245,
						Cache:         Cache{0, 0, 512 * 1024, 0, 64, 64},
						Siblings:      2,
						Cores:         1,
						FPU:           true,
						FPUExceptions: true,
						CPUIDLevel:    13,
						WP:            true,
						Features:      split("fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm rep_good nopl cpuid extd_apicid tsc_known_freq pni pclmulqdq ssse3 fma cx16 sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand hypervisor lahf_lm cmp_legacy cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw topoext perfctr_core ssbd ibpb vmmcall fsgsbase tsc_adjust bmi1 avx2 smep bmi2 rdseed adx smap clflushopt sha_ni xsaveopt xsavec xgetbv1 xsaves clzero xsaveerptr virt_ssbd arat arch_capabilities"),
						Bugs:          split("sysret_ss_attrs null_seg spectre_v1 spectre_v2 spec_store_bypass"),
						BogoMIPS:      3992.49,
						AddrSizes: struct {
							Phys int `json:"physical_bits,omitempty"`
							Virt int `json:"virtual_bits,omitempty"`
						}{40, 48},
						TLB: struct {
							N        int `json:"num_pages,omitempty"`
							PageSize int `json:"page_size,omitempty"`
						}{1024, 4096},
					},
					{
						Proc:          1,
						VendorID:      "AuthenticAMD",
						Family:        23,
						Model:         1,
						ModelName:     "AMD EPYC 7551 32-Core Processor",
						Rev:           2,
						Microcode:     0x1000065,
						Freq:          1996.245,
						Cache:         Cache{0, 0, 512 * 1024, 0, 64, 64},
						Siblings:      2,
						Cores:         1,
						APICID:        1,
						InitAPICID:    1,
						FPU:           true,
						FPUExceptions: true,
						CPUIDLevel:    13,
						WP:            true,
						Features:      split("fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm rep_good nopl cpuid extd_apicid tsc_known_freq pni pclmulqdq ssse3 fma cx16 sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand hypervisor lahf_lm cmp_legacy cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw topoext perfctr_core ssbd ibpb vmmcall fsgsbase tsc_adjust bmi1 avx2 smep bmi2 rdseed adx smap clflushopt sha_ni xsaveopt xsavec xgetbv1 xsaves clzero xsaveerptr virt_ssbd arat arch_capabilities"),
						Bugs:          split("sysret_ss_attrs null_seg spectre_v1 spectre_v2 spec_store_bypass"),
						BogoMIPS:      3992.49,
						AddrSizes: struct {
							Phys int `json:"physical_bits,omitempty"`
							Virt int `json:"virtual_bits,omitempty"`
						}{40, 48},
						TLB: struct {
							N        int `json:"num_pages,omitempty"`
							PageSize int `json:"page_size,omitempty"`
						}{1024, 4096},
					},
				},
			},
		},
		{
			name: "amd_epyc_ubuntu",
			info: Info{
				CPUs: []CPU{
					{
						VendorID:      "AuthenticAMD",
						Family:        23,
						Model:         1,
						ModelName:     "AMD EPYC 7551 32-Core Processor",
						Rev:           2,
						Microcode:     0x1000065,
						Freq:          1996.244,
						Cache:         Cache{0, 0, 512 * 1024, 0, 64, 64},
						Siblings:      2,
						Cores:         1,
						FPU:           true,
						FPUExceptions: true,
						CPUIDLevel:    13,
						WP:            true,
						Features:      split("fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm rep_good nopl cpuid extd_apicid tsc_known_freq pni pclmulqdq ssse3 fma cx16 sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand hypervisor lahf_lm cmp_legacy cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw topoext perfctr_core ssbd ibpb vmmcall fsgsbase tsc_adjust bmi1 avx2 smep bmi2 rdseed adx smap clflushopt sha_ni xsaveopt xsavec xgetbv1 xsaves clzero xsaveerptr virt_ssbd arat arch_capabilities"),
						Bugs:          split("sysret_ss_attrs null_seg spectre_v1 spectre_v2 spec_store_bypass"),
						BogoMIPS:      3992.48,
						AddrSizes: struct {
							Phys int `json:"physical_bits,omitempty"`
							Virt int `json:"virtual_bits,omitempty"`
						}{40, 48},
						TLB: struct {
							N        int `json:"num_pages,omitempty"`
							PageSize int `json:"page_size,omitempty"`
						}{1024, 4096},
					},
					{
						Proc:          1,
						VendorID:      "AuthenticAMD",
						Family:        23,
						Model:         1,
						ModelName:     "AMD EPYC 7551 32-Core Processor",
						Rev:           2,
						Microcode:     0x1000065,
						Freq:          1996.244,
						Cache:         Cache{0, 0, 512 * 1024, 0, 64, 64},
						Siblings:      2,
						Cores:         1,
						APICID:        1,
						InitAPICID:    1,
						FPU:           true,
						FPUExceptions: true,
						CPUIDLevel:    13,
						WP:            true,
						Features:      split("fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx mmxext fxsr_opt pdpe1gb rdtscp lm rep_good nopl cpuid extd_apicid tsc_known_freq pni pclmulqdq ssse3 fma cx16 sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand hypervisor lahf_lm cmp_legacy cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw topoext perfctr_core ssbd ibpb vmmcall fsgsbase tsc_adjust bmi1 avx2 smep bmi2 rdseed adx smap clflushopt sha_ni xsaveopt xsavec xgetbv1 xsaves clzero xsaveerptr virt_ssbd arat arch_capabilities"),
						Bugs:          split("sysret_ss_attrs null_seg spectre_v1 spectre_v2 spec_store_bypass"),
						BogoMIPS:      3992.48,
						AddrSizes: struct {
							Phys int `json:"physical_bits,omitempty"`
							Virt int `json:"virtual_bits,omitempty"`
						}{40, 48},
						TLB: struct {
							N        int `json:"num_pages,omitempty"`
							PageSize int `json:"page_size,omitempty"`
						}{1024, 4096},
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testReadProc(t, tc.name, tc.info)
		})
	}
}

func testReadProc(t *testing.T, name string, want Info) {
	buf, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	var got Info
	scanProc(&got, buf)

	if len(got.CPUs) != len(want.CPUs) {
		t.Fatalf("expected %d, got %d", len(want.CPUs), len(got.CPUs))
	}
	for i := range got.CPUs {
		// Compare one at a time so we limit what we barf to the
		// screen on failure.
		if !reflect.DeepEqual(got.CPUs[i], want.CPUs[i]) {
			c, _ := diff.Diff(got.CPUs[i], want.CPUs[i])
			t.Fatalf("#%d: values differ: %s", i, sprint(c))
		}
	}

	// Now check everything else.
	//
	// Drop the CPUs since we know they match.
	want.CPUs = nil
	got.CPUs = nil

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expected %#v, got %#v", want, got)
	}
}

func sprint(v interface{}) string {
	buf, err := json.MarshalIndent(v, " ", "  ")
	if err != nil {
		panic(err)
	}
	return string(buf)
}
