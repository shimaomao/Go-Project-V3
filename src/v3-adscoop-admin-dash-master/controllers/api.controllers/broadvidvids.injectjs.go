package apiControllers

import (
	"helpers"
	"structs"
	"github.com/gin-gonic/gin"
)

func broadvidVideosInjectJssViewallCtrl(c *gin.Context) {
	helpers.FindAll(&structs.InjectJss{}, c)
}

func broadvidVideosInjectJssViewCtrl(c *gin.Context) {
	helpers.FindOne(&structs.InjectJs{}, c)
}

func broadvidVideosInjectJssSaveCtrl(c *gin.Context) {

	uid, err := helpers.GetUserID(c)

	if err != nil {
		c.JSON(500, err)
		return
	}

	var ad = &structs.InjectJs{}

	err = helpers.SaveEntity(ad, c, "Inject JS has been saved", ad.Name+" has been saved", "success", 2, uid)

	if err != nil {
		c.JSON(500, err)
		return
	}

}
