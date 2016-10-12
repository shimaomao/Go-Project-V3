package structs

type AdFiltering struct {
	AdID    int64  `form:"ad_id"`
	EmbedID int64  `form:"embed_id"`
	Action  string `form:"action"`
}
