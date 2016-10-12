package apiControllers

import (
	"encoding/json"

	"app/structs"
	ginsessions "github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func AuthCallbackCtrl(c *gin.Context) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)

	if err != nil {
		c.Data(500, "text/html", []byte("This resource is not available1"))
		return
	}

	var adminUser structs.User
	b, _ := json.Marshal(user)

	adminUser.SaveJSON(b)

	if adminUser.ID == 0 {
		c.Data(500, "text/html", []byte("This resource is not available2"))
		return
	}

	ses := ginsessions.Default(c)

	ses.Set("UserID", adminUser.ID)
	ses.Save()

	c.Redirect(302, "/")
}

func RequireAuth(c *gin.Context) {
	sess := ginsessions.Default(c)

	v := sess.Get("UserID")

	if v == nil {
		c.Redirect(302, "/auth/login?provider=gplus")
		return
	}

	var u structs.User
	u.FindById(v.(uint))

	if u.ID == 0 {
		ses := ginsessions.Default(c)

		ses.Delete("UserID")
		ses.Save()

		c.Data(500, "text/html", []byte("This resource is not available3"))
		return
	}
}
