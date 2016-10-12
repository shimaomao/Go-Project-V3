// thank you calmh and __john_doe__ very much for helping me the past few hours, this will save me so much time by not having to write 60 controllers that are essentially the same exact logic

package helpers

import (
	"errors"

	log "github.com/Sirupsen/logrus"

	"bitbucket.org/broadscaler/broadscaler/app/structs"
	ginsessions "github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) (uid uint, err error) {

	sess := ginsessions.Default(c)

	v := sess.Get("UserID")

	if v == nil {
		return 0, errors.New("User not logged in")
	}

	return v.(uint), nil
}

func SaveEntity(object structs.SingleIntf, c *gin.Context, title string, message string, messageType string, productType uint, userID uint) (err error) {
	err = c.BindJSON(&object)

	if err != nil {
		log.Errorf("Cannot bind JSON: %s", err)
		c.JSON(500, err)
		return
	}

	err = object.Save()

	if err != nil {
		log.Errorf("Cannot save: %s", err)
		c.JSON(500, err)
		return
	}

	pushUpdate(title, message, messageType, productType, userID)
	return nil
}

func SaveEntityByRedir(object structs.SingleIntfByRedir, c *gin.Context, title string, message string, messageType string, productType uint, userID uint) (err error) {
	err = c.BindJSON(&object)

	if err != nil {
		log.Errorf("Cannot bind JSON: %s", err)
		c.JSON(500, err)
		return
	}

	err = object.Save(c.Param("id"))

	if err != nil {
		log.Errorf("Cannot save: %s", err)
		c.JSON(500, err)
		return
	}

	pushUpdate(title, message, messageType, productType, userID)
	return nil
}

func BasicSaveEntity(object structs.BasicIntf, c *gin.Context, title string, message string, messageType string, productType uint, userID uint) (err error) {
	err = c.BindJSON(&object)

	if err != nil {
		log.Errorf("Cannot bind JSON: %s", err)
		c.JSON(500, err)
		return
	}

	err = object.BasicSave()

	if err != nil {
		log.Errorf("Cannot save: %s", err)
		c.JSON(500, err)
		return
	}

	pushUpdate(title, message, messageType, productType, userID)
	return nil
}

func FindAll(object structs.GroupIntf, c *gin.Context) {
	err := object.FindAll()

	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, object)
}

func FindVisible(object structs.GroupVisbleIntf, c *gin.Context) {
	userID, err := GetUserID(c)

	if err != nil {
		c.JSON(500, err)
		return
	}
	err = object.FindVisible(userID)

	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, object)
}

func FindAllSchedules(object structs.GroupScheduleIntf, c *gin.Context) {
	err := object.FindAll(c.Param("id"))

	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, object)
}

func FindOne(object structs.SingleIntf, c *gin.Context) {
	err := object.Find(c.Param("id"))

	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, object)
}

func LinkToStructs() {
	structs.AdscoopsDB = AdscoopsDB
	structs.AdscoopsRealtimeDB = AdscoopsRealtimeDB
	structs.BroadvidDB = BroadvidDB
	structs.RedisPool = RedisPool

	structs.LinkToStructs()

}
