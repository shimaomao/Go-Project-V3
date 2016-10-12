package structs

import "encoding/xml"

type RSSAol struct {
	XMLName xml.Name `xml:"rss"`
	Items   ItemsAol `xml:"channel"`
}
type ItemsAol struct {
	XMLName  xml.Name  `xml:"channel"`
	ItemList []ItemAol `xml:"item"`
}
type ItemAol struct {
	ID          string  `xml:"id"`
	Guid        string  `xml:"guid"`
	Title       string  `xml:"title"`
	Description string  `xml:"description"`
	Thumbnail   Image   `xml:"thumbnail"`
	Player      Image   `xml:"player"`
	Group       Group   `xml:"group"`
	Content     Content `xml:"content"`
}
type Content struct {
	Url         string `xml:"url,attr"`
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Thumbnail   Image  `xml:"thumbnail"`
}

type Group struct {
	Content []struct {
		Url string `xml:"url,attr"`
	} `xml:"content"`
}

type Image struct {
	Url string `xml:"url,attr"`
}

type VideoItem struct {
	ID          int64
	Title       string `form:"title"`
	Description string `form:"description"`
	Thumbnail   string `form:"thumbnail"`
	Type        int64  `form:"type"`
	Enabled     int64
	VideoID     string `form:"videoId"`
	FeedID      int64
}

func (v VideoItem) TableName() string {
	return "items"
}

type RssItem struct {
	RssID  int64
	ItemID int64
}

func (v RssItem) TableName() string {
	return "rss_items"
}
