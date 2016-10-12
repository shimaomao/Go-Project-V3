package apiControllers

import (
	"app/helpers"
	"app/structs"
	ginsessions "github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func broadvidadsAdEmbedsRemoveCtrl(c *gin.Context) {
	sess := ginsessions.Default(c)

	v := sess.Get("UserID")

	if v == nil {
		return
	}

	var ae structs.AdEmbedSave

	err := c.BindJSON(&ae)

	if err == nil {
		ae.Remove(v.(uint))
	}
}

func broadvidadsAdEmbedsCopyCtrl(c *gin.Context) {
	sess := ginsessions.Default(c)

	v := sess.Get("UserID")

	if v == nil {
		return
	}

	var ae structs.AdEmbedSave

	err := c.BindJSON(&ae)

	if err == nil {
		ae.Copy(v.(uint))
	}
}

func broadvidadsAdEmbedsPauseCtrl(c *gin.Context) {
	sess := ginsessions.Default(c)

	v := sess.Get("UserID")

	if v == nil {
		return
	}

	var ae structs.AdEmbedSave

	err := c.BindJSON(&ae)

	if err == nil {
		ae.PauseToggle(v.(uint))
	}
}

func broadvidAdEmbedsViewallCtrl(c *gin.Context) {
	helpers.FindAll(&structs.AdEmbeds{}, c)
}

func broadvidAdEmbedsViewCtrl(c *gin.Context) {
	helpers.FindOne(&structs.AdEmbed{}, c)
}

func broadvidAdEmbedsSaveCtrl(c *gin.Context) {
	uid, err := helpers.GetUserID(c)

	if err != nil {
		c.JSON(500, err)
		return
	}

	var ad = &structs.AdEmbed{}

	err = helpers.SaveEntity(ad, c, "Ad Embed has been saved", ad.EmbedLabel+" has been saved", "success", 2, uid)

	if err != nil {
		c.JSON(500, err)
		return
	}

}
