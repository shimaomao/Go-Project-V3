package structs

import "time"

type AdscoopFeed struct {
	Id        int64
	Hash      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Redirects []Redirect
}

func (f AdscoopFeed) TableName() string {
	return "adscoop_feeds"
}
