package structs

import (
	log "github.com/Sirupsen/logrus"
)

/*Ads is the struct to grab a group of Ad types together */
type Ads []Ad

/*FindAll is the method to find all available broadvid ads */
func (a *Ads) FindAll() error {
	return BroadvidDB.Table("ads").Find(&a).Error
}

/*GetID is to grab the primary key of the ad */
func (a *Ad) GetID() int64 {
	return a.ID
}

/*Find grabs all the ad and ad embed information */
func (a *Ad) Find(id string) error {
	err := BroadvidDB.Table("ads").Find(&a, id).Error

	if err != nil {
		return err
	}

	var adEmbeds []AdEmbed

	err = BroadvidDB.Select("ads_embeds.*, embeds.type, embeds.embed_label").
		Joins("left outer join embeds on embeds.id = ads_embeds.embed_id").
		Where("ad_id = ?", a.ID).Find(&adEmbeds).Error

	if err != nil {
		return err
	}

	for _, ae := range adEmbeds {

		switch ae.Type {
		case "1":
			a.AdDesktop = append(a.AdDesktop, ae)
			log.Printf("added %s to desktop, type: %s", ae.EmbedLabel, ae.Type)
		case "2":
			a.AdHTML5 = append(a.AdHTML5, ae)
			log.Printf("added %s to html5, type: %s", ae.EmbedLabel, ae.Type)
		case "3":
			a.PlayerDesktop = append(a.PlayerDesktop, ae)
			log.Printf("added %s to player desktop, type: %s", ae.EmbedLabel, ae.Type)
		case "4":
			a.PlayerHTML5 = append(a.PlayerHTML5, ae)
			log.Printf("added %s to player html5, type: %s", ae.EmbedLabel, ae.Type)
		case "5":
			a.DefaultTag = append(a.DefaultTag, ae)
			log.Printf("added %s to default tag, type: %s", ae.EmbedLabel, ae.Type)
		}
	}

	return nil
}

/*Save saves the ad*/
func (a *Ad) Save() error {
	return BroadvidDB.Table("ads").Save(&a).Error
}

/*Ad is the type of object that is outputted at go.broadvid.com/ad/PRIKEY*/
type Ad struct {
	ID                      int64     `form:"id"`
	Label                   string    `form:"label"`
	AdHTML5                 []AdEmbed `sql:"-"`
	AdDesktop               []AdEmbed `sql:"-"`
	PlayerHTML5             []AdEmbed `sql:"-"`
	PlayerDesktop           []AdEmbed `sql:"-"`
	DefaultTag              []AdEmbed `sql:"-"`
	Type                    string    `form:"type"`
	TargetDiv               string    `form:"targetDiv"`
	Width                   string    `form:"width"`
	Height                  string    `form:"height"`
	TrackAnalytics          string    `form:"track_analytics"`
	Mute                    string    `form:"mute"`
	HideBorder              string    `form:"hideBorder"`
	BackgroundColor         string    `form:"backgroundColor"`
	BorderColor             string    `form:"borderColor"`
	AjaxPageNav             string    `form:"ajaxPageNav"`
	AdonTracking            string    `form:"adonTracking"`
	VsliderBottom           string    `form:"vsliderBottom"`
	ForceReload             string    `form:"forceReload"`
	ForceReloadKey          string    `form:"forceReloadKey"`
	ForceReloadValue        string    `form:"forceReloadValue"`
	ForceReloadTimeout      string    `form:"forceReloadTimeout"`
	NoLoad                  string    `form:"noLoad"`
	NoLoadKey               string    `form:"noLoadKey"`
	NoLoadValue             string    `form:"noLoadValue"`
	FireTwsPixel            int64     `form:"fire_tws_pixel"`
	DesktopOverride         int64     `form:"desktop_override"`
	MobileOverride          int64     `form:"mobile_override"`
	OverrideScripts         []Embed   `sql:"-"`
	LockWhitelistID         int64     `form:"lock_whitelist_id"`
	LockUseragentID         int64     `form:"lock_useragent_id"`
	HideCloseButton         uint      `form:"hide_close_button"`
	AvailableAdDesktops     []Embed
	AvailableAdHTML5s       []Embed
	AvailablePlayerDesktops []Embed
	AvailablePlayerHTML5s   []Embed
	AvailableDefaultTags    []Embed
}
