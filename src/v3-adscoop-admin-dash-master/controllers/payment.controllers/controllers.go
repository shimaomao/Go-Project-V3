package paymentController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"app/configSettting"
	"app/structs"

	"encoding/csv"
	"github.com/mailgun/mailgun-go"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"log"
	"os"
	"strings"
	"time"
)

//var db  structs.AdscoopsDB
func controllerPaymentHash(c *gin.Context) {
	db := structs.AdscoopsDB

	var retData struct {
		AdscoopClient AdscoopClient
		AdscoopPaymentHash
		Config TomlConfig
		Error  string
	}

	retData.Config = TomlConfig{}
	hashKey := c.Param("hash")
	db.Where("hash = ?", hashKey).Find(&retData.AdscoopPaymentHash)
	if retData.AdscoopPaymentHash.Id == 0 {
		fmt.Println("Payment hash not valid")
		//http.Error(w, fmt.Sprintf("Payment hash not valid"), http.StatusInternalServerError)
		return
	}

	db.Select("adscoop_clients.*, (adscoop_clients.charge_amount * 100) as charge_amount").Find(&retData.AdscoopClient, retData.AdscoopPaymentHash.ClientId)

	if retData.AdscoopClient.Id == 0 {
		fmt.Println("No Valid Client Found")
		//http.Error(w, fmt.Sprintf("No Valid Client Found"), http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "payment.tmpl", retData)

}

func controllerPaymentHashPost(c *gin.Context) {

	fmt.Println("Hey sdsad")
	fmt.Println(StripeData{})
	db := structs.AdscoopsDB
	sd := StripeData{}
	fmt.Println("PRINTED STRIPED")
	fmt.Println(sd)
	err := c.Bind(&sd)
	hashKey := c.Param("hash")

	log.Println("PAYMENT" + sd.StripeToken)
	if sd.PaymentHash != hashKey {
		log.Println("Hashes do not match")
		return
	}
	var asph AdscoopPaymentHash
	db.Where("hash = ? AND client_id = ?", sd.PaymentHash, sd.ClientId).Find(&asph)
	if asph.Id == 0 {
		log.Println("no payment hash found")
		return
	}
	var asc AdscoopClient

	db.Find(&asc, asph.ClientId)

	if asc.Id == 0 {
		log.Println("no client found")
		return
	}
	var found int64

	fmt.Println()
	db.Select("1").Table("adscoop_client_emails").Where("email = ? AND client_id = ?", sd.StripeEmail, asc.Id).Row().Scan(&found)
	if found != 1 {
		var retData struct {
			AdscoopClient AdscoopClient
			AdscoopPaymentHash
			Config TomlConfig
			Error  string
		}
		retData.Config = TomlConfig{}
		retData.AdscoopClient = asc
		retData.AdscoopPaymentHash = asph
		retData.Error = "E-mail does not match with what we have on record, please try again.  No payment has been made"
		c.HTML(http.StatusOK, "payment.tmpl", retData)
		return
	}

	stripe.Key = configSettting.StripeSecretKey

	customerParams := &stripe.CustomerParams{
		Desc:  "AdScoop Client:" + asc.Name,
		Email: sd.StripeEmail,
	}

	customerParams.SetSource(sd.StripeToken) // obtained with Stripe.js
	cu, err := customer.New(customerParams)

	if err != nil {

		fmt.Println("There has been an error, no payment has been made.  Please contact support to resolve. Err Message")
		//w :=http.ResponseWriter()
		//http.Error(w, fmt.Sprintf("There has been an error, no payment has been made.  Please contact support to resolve. Err Message: %s", err), http.StatusInternalServerError)
		//return
	}

	db.DB().Exec("UPDATE adscoop_clients SET stripe_token = ? WHERE id = ?", cu.ID, asc.Id)

	db.Delete(AdscoopPaymentHash{}, asph.Id)

	go func() {
		gun := mailgun.NewMailgun(configSettting.MailGunDomain, configSettting.MailGunApiKey, configSettting.MailGunPublicApiKey)
		m := mailgun.NewMessage("donotnreply <donotreply@mg.adscoops.com>", fmt.Sprintf("Adscoops: %s has updated their payment information", asc.Name), fmt.Sprintf("%s has updated their payment information", asc.Name), configSettting.AdminEmail)
		gun.Send(m)

	}()

	c.HTML(http.StatusOK, "payment_success.tmpl", nil)

}

func controllerCheckPayments() {
	db := structs.AdscoopsDB
	var asc []AdscoopClient

	db.Select("adscoop_clients.*").
		Joins(`LEFT OUTER JOIN (
			SELECT DISTINCT client_id, SUM(amount_charged) as amount_charged
			FROM adscoop_client_transactions
			GROUP BY adscoop_client_transactions.client_id
		)
		B2 on adscoop_clients.id = B2.client_id
		LEFT OUTER JOIN (
			SELECT DISTINCT client_id, SUM(cost) as cost
			FROM adscoop_campaigns
			LEFT OUTER JOIN (
				SELECT DISTINCT campaign_id, SUM(cost) as cost
					FROM adscoop_urls
					LEFT OUTER JOIN(
						SELECT DISTINCT url_id, SUM((CASE WHEN adscoop_campaigns.tracking_method = 0 THEN count
						 WHEN adscoop_campaigns.tracking_method = 1 THEN engagement
						 WHEN adscoop_campaigns.tracking_method = 2 THEN adscoop_trackings.load
					END) * adscoop_trackings.cpc) as cost
						FROM adscoop_trackings
						JOIN adscoop_urls ON adscoop_urls.id = adscoop_trackings.url_id
						JOIN adscoop_campaigns ON adscoop_campaigns.id = adscoop_urls.campaign_id
						GROUP BY adscoop_trackings.url_id
					) C1 on adscoop_urls.id = C1.url_id
				GROUP BY adscoop_urls.campaign_id
			) B3 on adscoop_campaigns.id = B3.campaign_id
			GROUP BY adscoop_campaigns.client_id
		) A1 on adscoop_clients.id = A1.client_id`).
		Where(`(adscoop_clients.charge_amount * .15 > (IFNULL(B2.amount_charged,0) - IFNULL(A1.cost,0)))
			AND adscoop_clients.approved_to_charge = 1
			AND adscoop_clients.stripe_token != ""`).
		Find(&asc)

	for _, y := range asc {
		y.ChargeClient()
	}
}

func controllerCheckExpiringClients() {
	var asc []AdscoopClient
	db := structs.AdscoopsDB
	db.Select("adscoop_clients.*").
		Joins(`LEFT OUTER JOIN (
			SELECT DISTINCT client_id, SUM(amount_charged) as amount_charged
			FROM adscoop_client_transactions
			GROUP BY adscoop_client_transactions.client_id
		)
		B2 on adscoop_clients.id = B2.client_id
		LEFT OUTER JOIN (
			SELECT DISTINCT client_id, SUM(cost) as cost
			FROM adscoop_campaigns
			LEFT OUTER JOIN (
				SELECT DISTINCT campaign_id, SUM(cost) as cost
					FROM adscoop_urls
					LEFT OUTER JOIN(
						SELECT DISTINCT url_id, SUM((CASE WHEN adscoop_campaigns.tracking_method = 0 THEN count
						 WHEN adscoop_campaigns.tracking_method = 1 THEN engagement
						 WHEN adscoop_campaigns.tracking_method = 2 THEN adscoop_trackings.load
					END) * adscoop_trackings.cpc) as cost
						FROM adscoop_trackings
						JOIN adscoop_urls ON adscoop_urls.id = adscoop_trackings.url_id
						JOIN adscoop_campaigns ON adscoop_campaigns.id = adscoop_urls.campaign_id
						GROUP BY adscoop_trackings.url_id
					) C1 on adscoop_urls.id = C1.url_id
				GROUP BY adscoop_urls.campaign_id
			) B3 on adscoop_campaigns.id = B3.campaign_id
			GROUP BY adscoop_campaigns.client_id
		) A1 on adscoop_clients.id = A1.client_id`).
		Where(`(adscoop_clients.charge_amount * .15 > (IFNULL(B2.amount_charged,0) - IFNULL(A1.cost,0)))
			AND adscoop_clients.approved_to_charge = 0
			AND adscoop_clients.expiration_warning = 0`).
		Find(&asc)

	for _, y := range asc {
		y.NotifyClientExpiring()
	}
}
func controllerCheckExpiredClients() {
	var asc []AdscoopClient

	db.Select("adscoop_clients.*").
		Joins(`LEFT OUTER JOIN (
			SELECT DISTINCT client_id, SUM(amount_charged) as amount_charged
			FROM adscoop_client_transactions
			GROUP BY adscoop_client_transactions.client_id
		)
		B2 on adscoop_clients.id = B2.client_id
		LEFT OUTER JOIN (
			SELECT DISTINCT client_id, SUM(cost) as cost
			FROM adscoop_campaigns
			LEFT OUTER JOIN (
				SELECT DISTINCT campaign_id, SUM(cost) as cost
					FROM adscoop_urls
					LEFT OUTER JOIN(
						SELECT DISTINCT url_id, SUM((CASE WHEN adscoop_campaigns.tracking_method = 0 THEN count
						 WHEN adscoop_campaigns.tracking_method = 1 THEN engagement
						 WHEN adscoop_campaigns.tracking_method = 2 THEN adscoop_trackings.load
					END) * adscoop_trackings.cpc) as cost
						FROM adscoop_trackings
						JOIN adscoop_urls ON adscoop_urls.id = adscoop_trackings.url_id
						JOIN adscoop_campaigns ON adscoop_campaigns.id = adscoop_urls.campaign_id
						GROUP BY adscoop_trackings.url_id
					) C1 on adscoop_urls.id = C1.url_id
				GROUP BY adscoop_urls.campaign_id
			) B3 on adscoop_campaigns.id = B3.campaign_id
			GROUP BY adscoop_campaigns.client_id
		) A1 on adscoop_clients.id = A1.client_id`).
		Where(`(0 > (IFNULL(B2.amount_charged,0) - IFNULL(A1.cost,0)))
			AND adscoop_clients.approved_to_charge = 0
			AND adscoop_clients.expiration_notice = 0`).
		Find(&asc)

	for _, y := range asc {
		y.NotifyClientExpired()
	}
}

func controllerCsv() {
	var adscoopClients []AdscoopClient

	location, _ := time.LoadLocation("America/Los_Angeles")

	today := time.Now()

	today = today.In(location)

	today = time.Date(today.Year(),
		today.Month(),
		today.Day(), 0, 0, 0, 0, location)

	today = today.In(time.UTC)

	hour := time.Now()

	hour = hour.In(location)
	hour = time.Date(hour.Year(),
		hour.Month(), hour.Day(),
		hour.Hour(), 0, 0, 0, location)

	hour = hour.In(time.UTC)

	db.Where("enable_reporting = 1 AND ((hourly_reporting = 0 AND report_last_sent < ?) OR (hourly_reporting = 1 AND report_last_sent < ?))", today, hour).Find(&adscoopClients)

	for _, asc := range adscoopClients {
		if asc.Id == 0 {
			continue
		}

		resetDate := today

		if asc.HourlyReporting {
			resetDate = hour
		}
		db.DB().Exec("UPDATE adscoop_clients SET report_last_sent = ? WHERE id = ?", resetDate, asc.Id)

		yesterday := today.Add(time.Duration(-24 * time.Hour))

		var ascs []AdscoopCsvStats

		var timeFormat = "01/02/06"
		var breakByHourGroup string

		latime, _ := time.LoadLocation("America/Los_Angeles")
		currentTime := time.Now().In(latime)

		tzstring := currentTime.Format("-07:00")

		tzstring = strings.Replace(tzstring, "-", "+", 1)

		if asc.BreakDownReportByHour {
			timeFormat = "01/02/06 03:00 PM"
			log.Println("breakdown by hour")
			breakByHourGroup = fmt.Sprintf(", HOUR(CONVERT_TZ(adscoop_trackings.timeslice, '%s', '+00:00'))", tzstring)
		}

		if asc.EnhancedReporting == 1 {
			query := db.Select(fmt.Sprintf(`adscoop_campaigns.name as campaign_name,
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
					adscoop_trackings.unique_identifier as unique_identifier`, tzstring)).
				Table(`adscoop_clients`).
				Joins(`JOIN adscoop_campaigns ON adscoop_campaigns.client_id = adscoop_clients.id
					   JOIN adscoop_urls ON adscoop_urls.campaign_id = adscoop_campaigns.id
					   JOIN adscoop_trackings ON adscoop_trackings.url_id = adscoop_urls.id`)
			if asc.HourlyReporting {
				query = query.Where(`adscoop_clients.id = ? AND adscoop_trackings.timeslice >= ?`,
					asc.Id, today)
			} else {
				query = query.Where(`adscoop_clients.id = ? AND adscoop_trackings.timeslice >= ? AND adscoop_trackings.timeslice < ?`,
					asc.Id, yesterday, today)
			}
			query.Group(fmt.Sprintf(`adscoop_campaigns.id, adscoop_urls.id, adscoop_trackings.unique_identifier %s`, breakByHourGroup)).
				Having(`(adscoop_campaigns.tracking_method = 0 AND SUM(adscoop_trackings.count) != 0)
					   	OR (adscoop_campaigns.tracking_method = 1 AND SUM(adscoop_trackings.engagement) != 0)
					   	OR (adscoop_campaigns.tracking_method = 2 AND SUM(adscoop_trackings.load) != 0)`).
				Find(&ascs)
		} else {
			query := db.Select(fmt.Sprintf(`adscoop_campaigns.name as campaign_name,
				   SUM(adscoop_trackings.count) as impressions,
				   SUM(adscoop_trackings.engagement) as verifieds,
				   SUM(adscoop_trackings.load) as loads,
					SUM(adscoop_trackings.count * adscoop_trackings.cpc) as cost,
					SUM(adscoop_trackings.engagement * adscoop_trackings.cpc) as cost_verified,
					SUM(adscoop_trackings.load * adscoop_trackings.cpc) as cost_load,
					adscoop_campaigns.cpc AS cpc,
					adscoop_campaigns.tracking_method AS tracking_method,
							CONVERT_TZ(adscoop_trackings.timeslice, '%s', '+00:00') as timeslice,
					adscoop_trackings.unique_identifier as unique_identifier`, tzstring)).
				Table(`adscoop_clients`).
				Joins(`JOIN adscoop_campaigns ON adscoop_campaigns.client_id = adscoop_clients.id
					   JOIN adscoop_urls ON adscoop_urls.campaign_id = adscoop_campaigns.id
					   JOIN adscoop_trackings ON adscoop_trackings.url_id = adscoop_urls.id`)

			if asc.HourlyReporting {
				query = query.Where(`adscoop_clients.id = ? AND adscoop_trackings.timeslice >= ?`,
					asc.Id, today)
			} else {
				query = query.Where(`adscoop_clients.id = ? AND adscoop_trackings.timeslice >= ? AND adscoop_trackings.timeslice < ?`,
					asc.Id, yesterday, today)
			}

			query.Group(fmt.Sprintf(`adscoop_campaigns.id, adscoop_trackings.unique_identifier %s`, breakByHourGroup)).
				Having(`(adscoop_campaigns.tracking_method = 0 AND SUM(adscoop_trackings.count) != 0)
					   	OR (adscoop_campaigns.tracking_method = 1 AND SUM(adscoop_trackings.engagement) != 0)
					   	OR (adscoop_campaigns.tracking_method = 2 AND SUM(adscoop_trackings.load) != 0)`).
				Find(&ascs)
		}

		if len(ascs) == 0 {
			log.Println("No campaigns ran, so not sending an e-mail")
			continue
		}

		csvfile, err := os.Create("report.csv")
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		defer csvfile.Close()

		records := [][]string{}

		if asc.EnhancedReporting == 1 {
			records = [][]string{{"Adscoop Stats for:", asc.Name},
				{"Date Range", yesterday.String() + " to " + today.String()},
				{""},
				{"Date", "Campaign Name", "Url", "Title", "Impressions", "Unique Identifier", "Cost", "CPC", "Tracked By"}}
		} else {
			records = [][]string{{"Adscoop Stats for:", asc.Name},
				{"Date Range", yesterday.String() + " to " + today.String()},
				{""},
				{"Date", "Campaign Name", "Impressions", "Unique Identifier", "Cost", "CPC", "Tracked By"}}
		}

		if asc.EnableReportAccountBalance {

			var amountCharged float64
			var totalSpend float64

			row := db.Table("adscoop_client_transactions").Where("client_id = ?", asc.Id).Select("SUM(amount_charged)").Row()
			row.Scan(&amountCharged)

			row2 := db.Select(`SUM((CASE WHEN adscoop_campaigns.tracking_method = 0 THEN count
							 WHEN adscoop_campaigns.tracking_method = 1 THEN engagement
							 WHEN adscoop_campaigns.tracking_method = 2 THEN adscoop_trackings.load
						END) * adscoop_trackings.cpc) as total_spend`).
				Table("adscoop_campaigns").
				Joins(`JOIN adscoop_urls ON adscoop_urls.campaign_id = adscoop_campaigns.ID
					JOIN adscoop_trackings ON adscoop_trackings.url_id = adscoop_urls.ID`).
				Where("adscoop_campaigns.client_id = ?", asc.Id).Row()
			row2.Scan(&totalSpend)

			records = append(records, []string{""})
			records = append(records, []string{"Account Balance:", fmt.Sprintf("$%s", RenderFloat("#,###.##", amountCharged-totalSpend))})
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

			db.Select("FORMAT(SUM(adscoop_trackings.count * adscoop_trackings.cpc), 2)").
				Table("adscoop_campaigns").
				Joins(`JOIN adscoop_urls ON adscoop_urls.campaign_id = adscoop_campaigns.id
				   JOIN adscoop_trackings ON adscoop_trackings.url_id = adscoop_urls.id`).
				Where("adscoop_campaigns.client_id = ? AND adscoop_trackings.timeslice >= ?", asc.Id, thisMonth).
				Row().Scan(&totalCharged)

			records = append(records, []string{"Spent This Month:", "$" + totalCharged})
		}

		for _, y := range ascs {
			var trackedBy string
			if y.TrackingMethod == 0 {
				trackedBy = "Click"
				y.ImpressionsString = RenderInteger("#,###.", y.Impressions)
				y.CostString = RenderFloat("#,###.##", y.Cost)
			}
			if y.TrackingMethod == 1 {
				trackedBy = "Verify"
				y.ImpressionsString = RenderInteger("#,###.", y.Verifieds)
				y.CostString = RenderFloat("#,###.##", y.CostVerified)
			}
			if y.TrackingMethod == 2 {
				trackedBy = "Load"
				y.ImpressionsString = RenderInteger("#,###.", y.Loads)
				y.CostString = RenderFloat("#,###.##", y.CostLoad)
			}
			if asc.EnhancedReporting == 1 {
				records = append(records, []string{y.Timeslice.Format(timeFormat), y.CampaignName, y.Url, y.Title, fmt.Sprintf("%s", y.ImpressionsString), y.UniqueIdentifier, fmt.Sprintf("$%s", y.CostString), y.Cpc, trackedBy})
			} else {
				records = append(records, []string{y.Timeslice.Format(timeFormat), y.CampaignName, fmt.Sprintf("%s", y.ImpressionsString), y.UniqueIdentifier, fmt.Sprintf("$%s", y.CostString), y.Cpc, trackedBy})

			}
		}

		writer := csv.NewWriter(csvfile)
		for _, record := range records {
			err := writer.Write(record)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
		}
		writer.Flush()
		gun := mailgun.NewMailgun(configSettting.MailGunDomain, configSettting.MailGunApiKey, configSettting.MailGunPublicApiKey)

		m := mailgun.NewMessage("donotnreply <donotreply@mg.adscoops.com>", fmt.Sprintf("Adscoops: New Report for %s", asc.Name), fmt.Sprintf("Attached is the report for %s", asc.Name))

		var asce []AdscoopClientEmail

		db.Where("client_id = ?", asc.Id).Find(&asce)

		for _, x := range asce {
			m.AddRecipient(x.Email)
		}

		m.AddAttachment("report.csv")
		gun.Send(m)

		os.Remove("report.csv")
	}
}
