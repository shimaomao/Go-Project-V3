package apiControllers

import (
	"github.com/gin-gonic/gin"
	"app/structs"
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
