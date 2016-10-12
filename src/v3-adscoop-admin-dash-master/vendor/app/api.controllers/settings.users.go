package apiControllers

import (
	"app/helpers"
	"app/structs"
	"github.com/gin-gonic/gin"
)

func settingsUsersViewallCtrl(c *gin.Context) {
	helpers.FindAll(&structs.Users{}, c)
}

func settingsUsersViewCtrl(c *gin.Context) {
	helpers.FindOne(&structs.User{}, c)
}

func settingsUsersSaveCtrl(c *gin.Context) {
	uid, err := helpers.GetUserID(c)

	if err != nil {
		c.JSON(500, err)
		return
	}

	var ad = &structs.User{}

	err = helpers.SaveEntity(ad, c, "User has been saved", ad.Name+" has been saved", "success", 4, uid)

	if err != nil {
		c.JSON(500, err)
		return
	}

}
