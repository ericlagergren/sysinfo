// Code generated by "stringer -type Implementer -linecomment"; DO NOT EDIT.

package sysinfo

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ARMLtd-65]
	_ = x[Broadcom-66]
	_ = x[Cavium-67]
	_ = x[Fujitsu-70]
	_ = x[NVIDIA-78]
	_ = x[HiSilicon-72]
	_ = x[Qualcomm-81]
	_ = x[Samsung-83]
	_ = x[Intel-105]
	_ = x[Apple-97]
}

const (
	_Implementer_name_0 = "ARM LtdBroadcomCavium"
	_Implementer_name_1 = "Fujitsu Ltd"
	_Implementer_name_2 = "HiSilicon Technologies Inc"
	_Implementer_name_3 = "NVIDIA Corporation"
	_Implementer_name_4 = "Qualcomm Technologies Inc"
	_Implementer_name_5 = "Samsung Technologies Inc"
	_Implementer_name_6 = "Apple Inc"
	_Implementer_name_7 = "Intel ARM parts"
)

var _Implementer_index_0 = [...]uint8{0, 7, 15, 21}

func (i Implementer) String() string {
	switch {
	case 65 <= i && i <= 67:
		i -= 65
		return _Implementer_name_0[_Implementer_index_0[i]:_Implementer_index_0[i+1]]
	case i == 70:
		return _Implementer_name_1
	case i == 72:
		return _Implementer_name_2
	case i == 78:
		return _Implementer_name_3
	case i == 81:
		return _Implementer_name_4
	case i == 83:
		return _Implementer_name_5
	case i == 97:
		return _Implementer_name_6
	case i == 105:
		return _Implementer_name_7
	default:
		return "Implementer(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}