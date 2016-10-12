package structs

type AutomatedAction struct {
	ID            int64  `form:"id"`
	AdID          int64  `form:"ad_id"`
	EmbedID       int64  `form:"embed_id"`
	Action        string `form:"action"`
	Limit         int64  `form:"limit"`
	TrackReferrer int64  `form:"track_referrer"`
}
