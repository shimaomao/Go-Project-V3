package structs

import "strings"

type TemplateHost struct {
	Host      string
	Blacklist map[int]bool `sql:"-"`
	Whitelist map[int]bool `sql:"-"`
}

type SaveList struct {
	Urls        []string `form:"url[]"`
	WhitelistID int64    `form:"whitelist_id"`
	BlacklistID int64    `form:"blacklist_id"`
}

func cleanUpMacros(tag string) string {
	tag = strings.Replace(tag, "%5breferrer_url%5d", "__page-url__", -1)
	tag = strings.Replace(tag, "%5bdescription_url%5d", "__description-url__", -1)
	tag = strings.Replace(tag, "%5btimestamp%5d", "__random-number__", -1)
	tag = strings.Replace(tag, "%5btimestamp", "__random-number__", -1)
	tag = strings.Replace(tag, "%5b", "", -1)

	tag = strings.Replace(tag, "[referrer_url]", "__page-url__", -1)
	tag = strings.Replace(tag, "[description_url]", "__page-url__", -1)
	tag = strings.Replace(tag, "[timestamp]", "__random-number__", -1)
	tag = strings.Replace(tag, "[TIMESTAMP]", "__random-number__", -1)

	tag = strings.Replace(tag, ":REQUIRED]", "]", -1)
	tag = strings.Replace(tag, ":RECOMMENDED]", "]", -1)
	tag = strings.Replace(tag, "[PLAYER_WIDTH]", "__player-width__", -1)
	tag = strings.Replace(tag, "[PLAYER_HEIGHT]", "__player-height__", -1)
	tag = strings.Replace(tag, "[PLAYER_POSITION]", "top", -1)
	tag = strings.Replace(tag, "[MEDIA_TITLE]", "__item-title__", -1)
	tag = strings.Replace(tag, "[MEDIA_DESCRIPTION]", "__item-description__", -1)
	tag = strings.Replace(tag, "[MEDIA_ID]", "__item-id__", -1)
	tag = strings.Replace(tag, "[CONTENT_MEDIA_URL]", "__item-url__", -1)
	tag = strings.Replace(tag, "[SOURCE_PAGE_URL]", "__page-url__", -1)
	tag = strings.Replace(tag, "[CONTENT_LENGTH]", "12", -1)
	tag = strings.Replace(tag, "[CONTENT_LENGTH", "120", -1)

	tag = strings.Replace(tag, "[INSERT_PAGE_URL]", "__page-url__", -1)
	tag = strings.Replace(tag, "[INSERT_PAGE_TITLE]", "__page-title__", -1)
	tag = strings.Replace(tag, "[CACHE_BUSTER]", "__random-number__", -1)

	tag = strings.Replace(tag, "[LR_URL]", "__item-url__", -1)
	tag = strings.Replace(tag, "[LR_URL_PATH]", "__page-url__", -1)
	tag = strings.Replace(tag, "[LR_TITLE]", "__item-title__", -1)
	tag = strings.Replace(tag, "[LR_DESCRIPTION]", "__item-description__", -1)
	tag = strings.Replace(tag, "[CACHEBUSTER]", "__random-number__", -1)
	tag = strings.Replace(tag, "[LR_DURATION]", "120", -1)
	tag = strings.Replace(tag, "[LR_AUTOPLAY]", "1", -1)
	tag = strings.Replace(tag, "[LR_WIDTH]", "__player-width__", -1)
	tag = strings.Replace(tag, "[LR_HEIGHT]", "__player-height__", -1)

	return tag
}
