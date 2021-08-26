package types

import "time"

type Job struct {
	Location    Address
	Description string
	Steps       []*JobStep
}

type JobStep struct {
	Type        string
	Timestamp   time.Time
	Description string
	User        string
}
