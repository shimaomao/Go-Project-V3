package paymentController

import (

	"github.com/gin-gonic/gin"
	"time"

)
func Setup( m *gin.Engine) {

	//   m.SetHTMLTemplate("layout")

	go processPayments()
	m.LoadHTMLGlob("./public/views/templates/**/*")

	m.GET("/payment/:hash", controllerPaymentHash)
	m.POST("/payment/:hash", controllerPaymentHashPost)
}

func processPayments() {
	controllerCheckPayments()
	controllerCheckExpiringClients()
	controllerCheckExpiredClients()
	controllerCsv()
	ticker := time.NewTicker(time.Minute * 1)
	for _= range ticker.C {
		controllerCheckPayments()
		controllerCheckExpiringClients()
		controllerCheckExpiredClients()
		controllerCsv()
	}

}
