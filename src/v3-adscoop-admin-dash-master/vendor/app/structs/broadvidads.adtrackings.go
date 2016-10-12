package structs

import "time"

type AdTracking struct {
	Timeslice     time.Time
	AdID          int64
	EmbedID       int64
	Action        string
	Count         int64
	Limit         int64
	AdLabel       string
	EmbedLabel    string
	TrackReferrer int64
}

func (a AdTracking) TableName() string {
	return "ad_tracking"
}
