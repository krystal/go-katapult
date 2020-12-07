package buildspec

type ScheduleInterval string

const (
	ScheduledHourly  ScheduleInterval = "hourly"
	ScheduledDaily   ScheduleInterval = "daily"
	ScheduledWeekly  ScheduleInterval = "weekly"
	ScheduledMonthly ScheduleInterval = "monthly"
)

type Schedule struct {
	Interval  ScheduleInterval `xml:",omitempty" json:"interval,omitempty" yaml:"interval,omitempty"`
	Frequency int              `xml:",omitempty" json:"frequency,omitempty" yaml:"frequency,omitempty"`
	Time      int              `xml:",omitempty" json:"time,omitempty" yaml:"time,omitempty"`
}
