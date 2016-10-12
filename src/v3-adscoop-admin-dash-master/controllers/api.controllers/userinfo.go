package apiControllers

import (
	ginsessions "github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"app/structs"
)

func UserInfoCtrl(c *gin.Context) {
	ses := ginsessions.Default(c)

	v := ses.Get("UserID")

	if v == nil {
		return
	}

	var u structs.User
	u.FindById(v.(uint))

	if u.ID == 0 {
		return
	}

	c.JSON(200, u)
}
