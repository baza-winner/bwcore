// Code generated by "stringer -type FuncGenType"; DO NOT EDIT.

package helper

import "strconv"

const _FuncGenType_name = "SetMethodSliceMethodSimpleFunc"

var _FuncGenType_index = [...]uint8{0, 9, 20, 30}

func (i FuncGenType) String() string {
	if i >= FuncGenType(len(_FuncGenType_index)-1) {
		return "FuncGenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FuncGenType_name[_FuncGenType_index[i]:_FuncGenType_index[i+1]]
}
