package apiControllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"

	"app/sockets"
	"app/structs"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func PingCtrl(c *gin.Context) {
	var message structs.RealTimeUpdate
	message.Event = "message"
	message.Data.Title = "this is a title"
	message.Data.Message = time.Now().String()

	sockets.BroadcastMessage(1, message.JSONify(), "updates")
}

func AdscoopsRealtimeUpdatesCtrl(c *gin.Context) {
	ws, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(c.Writer, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	client := ws.RemoteAddr()
	sockCli := sockets.ClientConn{"adscoops_updates", ws, client}
	sockets.AddClient(sockCli)
	sendAllStats()

	for {
		log.Println("Client list", len(sockets.ActiveClients), sockets.ActiveClients)
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			sockets.DeleteClient(sockCli)
			log.Println("bye")
			return
		}

		log.Println("Message type: ", messageType, " p: ", fmt.Sprintf("%s", p))
		// broadcastMessage(messageType, p, "updates")
	}
}

func RealtimeUpdatesCtrl(c *gin.Context) {
	ws, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(c.Writer, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	client := ws.RemoteAddr()
	sockCli := sockets.ClientConn{"updates", ws, client}
	sockets.AddClient(sockCli)

	for {
		log.Println("Client list", len(sockets.ActiveClients), sockets.ActiveClients)
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			sockets.DeleteClient(sockCli)
			log.Println("bye")
			return
		}

		log.Println("Message type: ", messageType, " p: ", fmt.Sprintf("%s", p))
		// broadcastMessage(messageType, p, "updates")
	}
}

func statsUpdates() {
	for {
		sendAllStats()
		time.Sleep(10 * time.Second)
	}
}

func sendAllStats() {
	// Daily Impressions stats
	sendDailyImpsStats()

	// Vertical stats
	sendVerticalStats()

	// Realtime stats Yesterday
	sendRealtimeStatsRedir()

	// Realtime stats clients
	sendRealtimeStatsClients()
}

func sendRealtimeStatsClients() {

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

	payload, err := json.Marshal(&retData)

	if err != nil {
		log.Errorf("Cannot prep data to send: %s", err)
		return
	}

	var message structs.RealTimeUpdate
	message.Event = "getRealTimeClientStats"
	message.Data.Message = fmt.Sprintf("%s", payload)

	sockets.BroadcastMessage(1, message.JSONify(), "adscoops_updates")

}

func sendRealtimeStatsRedir() {
	var retData struct {
		Data structs.MultiTempStats
	}

	retData.Data.Yesterday()

	payload, err := json.Marshal(&retData)

	if err != nil {
		log.Errorf("Cannot prep data to send: %s", err)
		return
	}

	var retDataToday struct {
		Data structs.MultiTempStats
	}

	retDataToday.Data.Today()

	payloadTitle, err := json.Marshal(&retDataToday)

	if err != nil {
		log.Errorf("Cannot prep data to send: %s", err)
		return
	}

	var message structs.RealTimeUpdate
	message.Event = "getRealtimeRedirStats"
	message.Data.Message = fmt.Sprintf("%s", payload)
	message.Data.Title = fmt.Sprintf("%s", payloadTitle)

	sockets.BroadcastMessage(1, message.JSONify(), "adscoops_updates")
}

func sendDailyImpsStats() {

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

	payload, err := json.Marshal(&retData)

	if err != nil {
		log.Errorf("Cannot prep data to send: %s", err)
		return
	}

	var message structs.RealTimeUpdate
	message.Event = "dailyImpressionsStats"
	message.Data.Message = fmt.Sprintf("%s", payload)

	sockets.BroadcastMessage(1, message.JSONify(), "adscoops_updates")
}

func sendVerticalStats() {
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

	payload, err := json.Marshal(&retData)

	if err != nil {
		log.Errorf("Cannot prep data to send: %s", err)
		return
	}

	var message structs.RealTimeUpdate
	message.Event = "getVerticalStats"
	message.Data.Message = fmt.Sprintf("%s", payload)

	sockets.BroadcastMessage(1, message.JSONify(), "adscoops_updates")
}
