package enumeration_test

import (
	"testing"

	"github.com/faris-arifiansyah/fws-rsvp/enumeration"
	"github.com/stretchr/testify/assert"
)

func TestAttendanceTypeString(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		attendanceType enumeration.AttendanceType
		expectedStr    string
	}{
		{
			attendanceType: enumeration.AttendanceType(0),
			expectedStr:    "No",
		},
		{
			attendanceType: enumeration.AttendanceType(1),
			expectedStr:    "Yes",
		},
		{
			attendanceType: enumeration.AttendanceType(2),
			expectedStr:    "Maybe",
		},
		{
			attendanceType: enumeration.AttendanceType(-1),
			expectedStr:    "AttendanceType(-1)",
		},
	}

	for _, tc := range testCases {
		assert.Equal(tc.expectedStr, tc.attendanceType.String())
	}
}
