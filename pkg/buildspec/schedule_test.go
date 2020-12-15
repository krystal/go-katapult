package buildspec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScheduleIntervals(t *testing.T) {
	tests := []struct {
		name   string
		obj    ScheduleInterval
		string string
	}{
		{
			name:   "ScheduledHourly",
			obj:    ScheduledHourly,
			string: "hourly",
		},
		{
			name:   "ScheduledDaily",
			obj:    ScheduledDaily,
			string: "daily",
		},
		{
			name:   "ScheduledWeekly",
			obj:    ScheduledWeekly,
			string: "weekly",
		},
		{
			name:   "ScheduledMonthly",
			obj:    ScheduledMonthly,
			string: "monthly",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.IsType(t, ScheduleInterval(""), tt.obj)
			assert.Equal(t, tt.string, string(tt.obj))
		})
	}
}

func TestSchedule_Marshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *Schedule
	}{
		{
			name: "empty",
			obj:  &Schedule{},
		},
		{
			name: "full",
			obj: &Schedule{
				Interval:  ScheduledDaily,
				Frequency: 1,
				Time:      13,
			},
		},
	}
	for _, tt := range tests {
		t.Run("json_"+tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
		t.Run("xml_"+tt.name, func(t *testing.T) {
			testXMLMarshaling(t, tt.obj)
		})
		t.Run("yaml_"+tt.name, func(t *testing.T) {
			testYAMLMarshaling(t, tt.obj)
		})
	}
}
