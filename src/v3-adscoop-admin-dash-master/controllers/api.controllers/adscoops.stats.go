package apiControllers

import (
	"structs"

	"github.com/gin-gonic/gin"
)

func statsRealTimeClientCtrl(c *gin.Context) {

	type FlotData struct {
		Label string      `json:"label"`
		Data  [][]float64 `json:"data"`
		Color string      `json:"color"`
	}
	type MakeData map[uint]FlotData

	var retData struct {
		Clients []FlotData
	}

	var clients structs.MultiClientTempStats

	clients.Today()

	retDataClients := make(map[uint]FlotData)

	for _, x := range clients {
		client := retDataClients[x.ClientID]

		client.Color = x.ChartColor

		client.Label = x.ClientName
		client.Data = append(client.Data, []float64{float64(x.Timeslice.Unix()), x.Revenue})

		retDataClients[x.ClientID] = client
	}

	for _, x := range retDataClients {
		retData.Clients = append(retData.Clients, x)
	}

	c.JSON(200, retData)
}

func statsRealTimeYesterdayCtrl(c *gin.Context) {
	var retData struct {
		Data structs.MultiTempStats
	}

	retData.Data.Yesterday()

	c.JSON(200, retData)
}

func statsRealTimeTodayCtrl(c *gin.Context) {
	var retData struct {
		Data structs.MultiTempStats
	}

	retData.Data.Today()

	c.JSON(200, retData)
}

func statsDailyImpressionCountCtrl(c *gin.Context) {
	var retData struct {
		Impressions uint
		Limit       uint
		Revenue     float64
	}
	var trackings structs.Trackings
	trackings.GetDailyImpressionsCount()
	retData.Impressions = trackings.Count
	retData.Limit = trackings.Engagement
	retData.Revenue = trackings.Cpc
	c.JSON(200, retData)
}

func statsDailyImpressionCountByVerticalCtrl(c *gin.Context) {
	var retData struct {
		Load struct {
			Impressions uint
			Limit       uint
			Revenue     float64
			Breakdown   structs.TrackingRows
		}
		Verification struct {
			Impressions uint
			Limit       uint
			Revenue     float64
			Breakdown   structs.TrackingRows
		}
		Clicks struct {
			Impressions uint
			Limit       uint
			Revenue     float64
			Breakdown   structs.TrackingRows
		}
	}
	var trackings structs.Trackings
	trackings.GetDailyImpressionsCountByVertical(0)
	retData.Clicks.Impressions = trackings.Count
	retData.Clicks.Limit = trackings.Engagement
	retData.Clicks.Revenue = trackings.Cpc
	retData.Clicks.Breakdown.GetDailyImpressionsCountByVertical(0)

	trackings = structs.Trackings{}
	trackings.GetDailyImpressionsCountByVertical(1)
	retData.Verification.Impressions = trackings.Count
	retData.Verification.Limit = trackings.Engagement
	retData.Verification.Revenue = trackings.Cpc
	retData.Verification.Breakdown.GetDailyImpressionsCountByVertical(1)

	trackings = structs.Trackings{}
	trackings.GetDailyImpressionsCountByVertical(2)
	retData.Load.Impressions = trackings.Count
	retData.Load.Limit = trackings.Engagement
	retData.Load.Revenue = trackings.Cpc
	retData.Load.Breakdown.GetDailyImpressionsCountByVertical(2)

	c.JSON(200, retData)
}
