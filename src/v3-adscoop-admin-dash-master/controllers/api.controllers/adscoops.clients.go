package apiControllers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/olekukonko/tablewriter"

	"github.com/gin-gonic/gin"
	"app/helpers"
	"app/structs"
)

func clientsViewRedirStatsCtrl(c *gin.Context) {
	_, err := helpers.GetUserID(c)

	if err != nil {
		log.Errorf("Cannot get user ID, so cancelling request: %s", err)
		return
	}

	var client structs.Client
	if err := client.Find(c.Param("id")); err != nil {
		log.Errorf("Client cannot be found")
		return
	}

	c.JSON(http.StatusOK, client.RedirRealtimeStats())
}

func clientsViewAssociatedRedirectsCtrl(c *gin.Context) {
	_, err := helpers.GetUserID(c)

	if err != nil {
		log.Errorf("Cannot get user ID, so cancelling request: %s", err)
		return
	}

	var client structs.Client
	if err := client.Find(c.Param("id")); err != nil {
		log.Errorf("Client cannot be found")
		return
	}

	c.JSON(http.StatusOK, client.AssociatedRedirects())
}

func clientsViewallCtrl(c *gin.Context) {
	helpers.FindAll(&structs.Clients{}, c)
}

func clientsViewVisibleCtrl(c *gin.Context) {
	helpers.FindVisible(&structs.Clients{}, c)
}

func clientsUpdateCampaignSortPerUserCtrl(c *gin.Context) {
	uid, err := helpers.GetUserID(c)

	if err != nil {
		log.Errorf("Cannot get user ID, so canceling request: %s", err)
		return
	}

	var userSettings structs.UserAdscoopsClientSetting

	if err := c.BindJSON(&userSettings); err != nil {
		log.Errorf("Cannot read posted data: %s", err)
		return
	}

	userSettings.UserID = uid
	if err := userSettings.Save(); err != nil {
		log.Errorf("Cannot save user settings: %s", err)
		return
	}
}

func clientsAutoChargeCtrl(c *gin.Context) {
	var postData struct {
		Charge uint
	}

	if err := c.BindJSON(&postData); err != nil {
		c.JSON(500, err)
		return
	}

	var client structs.Client

	if err := client.Find(c.Param("id")); err != nil {
		c.JSON(500, err)
		return
	}

	if err := client.Charge(postData.Charge); err != nil {
		c.JSON(500, err)
		return
	}

}

func clientsManualChargeCtrl(c *gin.Context) {
	var postData struct {
		Charge uint
	}

	if err := c.BindJSON(&postData); err != nil {
		c.JSON(500, err)
		return
	}

	var client structs.Client

	if err := client.Find(c.Param("id")); err != nil {
		c.JSON(500, err)
		return
	}

	if err := client.ManualCharge(postData.Charge); err != nil {
		c.JSON(500, err)
		return
	}
}

func clientsNewPaymentHashCtrl(c *gin.Context) {
	var client structs.Client

	if err := client.Find(c.Param("id")); err != nil {
		c.JSON(500, err)
		return
	}

	var aph structs.AdscoopPaymentHash
	aph, err := client.CreateAPH()

	if err != nil {
		c.JSON(500, err)
		return
	}

	c.Redirect(302, "http://localhost:8080/payment/"+aph.Hash)
}

func clientsViewCtrl(c *gin.Context) {
	helpers.FindOne(&structs.Client{}, c)
}

func clientsSaveCtrl(c *gin.Context) {

	uid, err := helpers.GetUserID(c)

	if err != nil {
		c.JSON(500, err)
		return
	}

	var client = &structs.Client{}
	helpers.SaveEntity(client, c, "Client has been saved", client.Name+" has been saved", "success", 1, uid)
}

func clientsShowReportCtrl(c *gin.Context) {
	var ascvs structs.AdscoopClientCsvPost
	var err error

	if c.Query("download") == "true" {
		err = c.Bind(&ascvs)
	} else {
		err = c.BindJSON(&ascvs)
	}

	if err != nil {
		log.Errorf("Error reading json: %s", err)
		return
	}

	var asc structs.Client

	location, _ := time.LoadLocation("America/Los_Angeles")

	var layout = "2006-01-02"
	start, err := time.ParseInLocation(layout, ascvs.StartDate, location)

	if err != nil {
		c.Writer.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(c.Writer, "Start date not formatted correctly")
		return
	}
	end, err := time.ParseInLocation(layout, ascvs.EndDate, location)

	if err != nil {
		c.Writer.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(c.Writer, "End Date not formatted correctly")
		return
	}

	end = end.Add(time.Duration(24 * time.Hour))

	start = start.In(time.UTC)
	end = end.In(time.UTC)

	AdscoopsDB.Find(&asc, c.Param("id"))

	if asc.ID == 0 {
		c.Writer.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(c.Writer, "Client not found")
		return
	}

	var ascs []structs.AdscoopCsvStats

	var timeFormat = "01/02/06"
	var breakByHourWhere string
	var breakByHourGroup string

	latime, _ := time.LoadLocation("America/Los_Angeles")
	currentTime := time.Now().In(latime)

	tzstring := currentTime.Format("-07:00")

	tzstring = strings.Replace(tzstring, "-", "+", 1)

	if asc.BreakDownReportByHour {
		timeFormat = "01/02/06 03:00 PM"
		breakByHourGroup = fmt.Sprintf("HOUR(CONVERT_TZ(adscoop_trackings.timeslice, '%s', '+00:00')),", tzstring)
	}

	if asc.EnhancedReporting {
		AdscoopsDB.Select(fmt.Sprintf(`adscoop_campaigns.name as campaign_name,
					 adscoop_urls.title as title,
					 adscoop_urls.url as url,
					 SUM(adscoop_trackings.count) as impressions,
					 SUM(adscoop_trackings.engagement) as verifieds,
					 SUM(adscoop_trackings.load) as loads,
					SUM(adscoop_trackings.count * adscoop_trackings.cpc) as cost,
					SUM(adscoop_trackings.engagement * adscoop_trackings.cpc) as cost_verified,
					SUM(adscoop_trackings.load * adscoop_trackings.cpc) as cost_load,
					adscoop_campaigns.cpc AS cpc,
					adscoop_campaigns.tracking_method AS tracking_method,
					CONVERT_TZ(adscoop_trackings.timeslice, '%s', '+00:00') as timeslice,
					%s
					adscoop_trackings.unique_identifier as unique_identifier`, tzstring, breakByHourWhere)).
			Table(`adscoop_clients`).
			Joins(`JOIN adscoop_campaigns ON adscoop_campaigns.client_id = adscoop_clients.ID
						 JOIN adscoop_urls ON adscoop_urls.campaign_id = adscoop_campaigns.ID
						 JOIN adscoop_trackings ON adscoop_trackings.url_id = adscoop_urls.ID`).
			Where(`adscoop_clients.ID = ? AND adscoop_trackings.timeslice >= ? AND adscoop_trackings.timeslice < ?`,
				asc.ID, start, end).
			Group(fmt.Sprintf(`adscoop_campaigns.ID, adscoop_urls.ID, DATE(CONVERT_TZ(adscoop_trackings.timeslice, '%s', '+00:00')), %s adscoop_trackings.unique_identifier`, tzstring, breakByHourGroup)).
			Find(&ascs)
	} else {
		AdscoopsDB.Select(fmt.Sprintf(`adscoop_campaigns.name as campaign_name,
							 SUM(adscoop_trackings.count) as impressions,
							 SUM(adscoop_trackings.engagement) as verifieds,
							 SUM(adscoop_trackings.load) as loads,
							SUM(adscoop_trackings.count * adscoop_trackings.cpc) as cost,
							SUM(adscoop_trackings.engagement * adscoop_trackings.cpc) as cost_verified,
							SUM(adscoop_trackings.load * adscoop_trackings.cpc) as cost_load,
							adscoop_campaigns.cpc AS cpc,
							adscoop_campaigns.tracking_method AS tracking_method,
							CONVERT_TZ(adscoop_trackings.timeslice, '%s', '+00:00') as timeslice,
							%s
							adscoop_clients.id as client_id,
							adscoop_trackings.unique_identifier as unique_identifier`, tzstring, breakByHourWhere)).
			Table(`adscoop_clients`).
			Joins(`JOIN adscoop_campaigns ON adscoop_campaigns.client_id = adscoop_clients.ID
								 JOIN adscoop_urls ON adscoop_urls.campaign_id = adscoop_campaigns.ID
								 JOIN adscoop_trackings ON adscoop_trackings.url_id = adscoop_urls.ID`).
			Where(`adscoop_clients.ID = ? AND adscoop_trackings.timeslice >= ? AND adscoop_trackings.timeslice < ?`,
				asc.ID, start, end).
			Group(fmt.Sprintf(`adscoop_campaigns.ID, DATE(CONVERT_TZ(adscoop_trackings.timeslice, '%s', '+00:00')), %s adscoop_trackings.unique_identifier`, tzstring, breakByHourGroup)).
			Find(&ascs)
	}

	if len(ascs) == 0 {
		c.Writer.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(c.Writer, "No campaigns ran yesterday, so no report to show")
		return
	}

	var buffer bytes.Buffer

	records := [][]string{}

	if asc.EnableReportAccountBalance {

		var amountCharged float64
		var totalSpend float64

		row := AdscoopsDB.Table("adscoop_client_transactions").Where("client_id = ?", asc.ID).Select("SUM(amount_charged)").Row()
		row.Scan(&amountCharged)

		row2 := AdscoopsDB.Select(`SUM((CASE WHEN adscoop_campaigns.tracking_method = 0 THEN count
							 WHEN adscoop_campaigns.tracking_method = 1 THEN engagement
							 WHEN adscoop_campaigns.tracking_method = 2 THEN adscoop_trackings.load
						END) * adscoop_trackings.cpc) as total_spend`).
			Table("adscoop_campaigns").
			Joins(`JOIN adscoop_urls ON adscoop_urls.campaign_id = adscoop_campaigns.ID
					JOIN adscoop_trackings ON adscoop_trackings.url_id = adscoop_urls.ID`).
			Where("adscoop_campaigns.client_id = ?", asc.ID).Row()
		row2.Scan(&totalSpend)

		records = append(records, []string{"Account Balance:", fmt.Sprintf("$%s", helpers.RenderFloat("#,###.##", amountCharged-totalSpend))})
		records = append(records, []string{""})

	}

	if asc.ShowMtdSpendInReport {
		location, _ := time.LoadLocation("America/Los_Angeles")

		thisMonth := time.Now()

		thisMonth = thisMonth.In(location)

		thisMonth = time.Date(thisMonth.Year(),
			thisMonth.Month(), 1, 0, 0, 0, 0, location)

		thisMonth = thisMonth.In(time.UTC)

		var totalCharged string

		AdscoopsDB.Select("FORMAT(SUM(adscoop_trackings.count * adscoop_trackings.cpc), 2)").
			Table("adscoop_campaigns").
			Joins(`JOIN adscoop_urls ON adscoop_urls.campaign_id = adscoop_campaigns.id
					 JOIN adscoop_trackings ON adscoop_trackings.url_id = adscoop_urls.id`).
			Where("adscoop_campaigns.client_id = ? AND adscoop_trackings.timeslice >= ?", asc.ID, thisMonth).
			Row().Scan(&totalCharged)

		records = append(records, []string{"Spent This Month:", "$" + totalCharged})
	}

	if asc.EnhancedReporting {
		records = append(records, []string{"Adscoop Stats for:", asc.Name})
		records = append(records, []string{"Date Range", start.String() + " to " + end.String()})
		records = append(records, []string{""})
		records = append(records, []string{"Date", "Campaign Name", "Url", "Title", "Impressions", "Unique Identifier", "Cost", "CPC", "Tracked By"})
	} else {
		records = append(records, []string{"Adscoop Stats for:", asc.Name})
		records = append(records, []string{"Date Range", start.String() + " to " + end.String()})
		records = append(records, []string{""})
		records = append(records, []string{"Date", "Campaign Name", "Impressions", "Unique Identifier", "Cost", "CPC", "Tracked By"})
	}

	for _, y := range ascs {
		var trackedBy string
		if y.TrackingMethod == 0 {
			trackedBy = "Click"
			y.ImpressionsString = helpers.RenderInteger("#,###.", y.Impressions)
			y.CostString = helpers.RenderFloat("#,###.##", y.Cost)
		}
		if y.TrackingMethod == 1 {
			trackedBy = "Verify"
			y.ImpressionsString = helpers.RenderInteger("#,###.", y.Verifieds)
			y.CostString = helpers.RenderFloat("#,###.##", y.CostVerified)
		}
		if y.TrackingMethod == 2 {
			trackedBy = "Load"
			y.ImpressionsString = helpers.RenderInteger("#,###.", y.Loads)
			y.CostString = helpers.RenderFloat("#,###.##", y.CostLoad)
		}

		if asc.EnhancedReporting {
			records = append(records, []string{y.Timeslice.Format(timeFormat), y.CampaignName, y.Url, y.Title, fmt.Sprintf("%s", y.ImpressionsString), y.UniqueIdentifier, fmt.Sprintf("$%s", y.CostString), y.Cpc, trackedBy})
		} else {
			records = append(records, []string{y.Timeslice.Format(timeFormat), y.CampaignName, fmt.Sprintf("%s", y.ImpressionsString), y.UniqueIdentifier, fmt.Sprintf("$%s", y.CostString), y.Cpc, trackedBy})
		}
	}

	writer := csv.NewWriter(&buffer)
	for _, record := range records {
		err := writer.Write(record)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
	writer.Flush()

	if c.Query("download") == "true" {
		c.Writer.Header().Set("Content-Type", "text/csv")
		c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename='%s_from_%s_to_%s.csv'", asc.Name, ascvs.StartDate, ascvs.EndDate))
		fmt.Fprint(c.Writer, buffer.String())
		return
	}

	table := tablewriter.NewWriter(c.Writer)
	table.SetBorder(false)    // Set Border to false
	table.AppendBulk(records) // Add Bulk Data
	table.Render()
	return
}
