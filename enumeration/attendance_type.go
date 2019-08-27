package enumeration

import (
	"fmt"
)

type AttendanceType int16

const (
	AttendanceTypeNo AttendanceType = iota
	AttendanceTypeYes
	AttendanceTypeMaybe
)

var atMap = map[AttendanceType]string{
	AttendanceTypeNo:    "No",
	AttendanceTypeYes:   "Yes",
	AttendanceTypeMaybe: "Maybe",
}

func (at AttendanceType) String() string {
	if str, ok := atMap[at]; ok {
		return str
	}
	return fmt.Sprintf("AttendanceType(%d)", at)
}
