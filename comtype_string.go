// Code generated by "stringer -type ComType"; DO NOT EDIT

package main

import "fmt"

const _ComType_name = "ComEnumComRecordComModuleComInterfaceComDispatchComCoclassComAliasComUnionComMax"

var _ComType_index = [...]uint8{0, 7, 16, 25, 37, 48, 58, 66, 74, 80}

func (i ComType) String() string {
	if i < 0 || i >= ComType(len(_ComType_index)-1) {
		return fmt.Sprintf("ComType(%d)", i)
	}
	return _ComType_name[_ComType_index[i]:_ComType_index[i+1]]
}
