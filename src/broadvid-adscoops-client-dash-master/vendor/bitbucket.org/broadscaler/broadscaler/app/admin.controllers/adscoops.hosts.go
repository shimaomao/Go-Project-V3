package adminControllers

import (
	"bitbucket.org/broadscaler/broadscaler/app/structs"
	"github.com/gin-gonic/gin"
)

func hostsViewallCtrl(c *gin.Context) {
	var hosts structs.Hosts
	err := hosts.FindAll()

	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, hosts)
}
