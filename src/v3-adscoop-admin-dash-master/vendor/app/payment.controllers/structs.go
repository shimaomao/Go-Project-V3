package paymentController

import (
	"bytes"
	"html/template"
	"log"
	"time"

	"github.com/mailgun/mailgun-go"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"app/structs"
	"app/configSettting"
	
)

var db =structs.AdscoopsDB
type TomlConfig struct {
	MysqlConfig         string
	StripeSecretKey     string
	StripePublishKey    string
	MailGunDomain       string
	MailGunApiKey       string
	MailGunPublicApikey string
	AdminEmail          string
}

type AdscoopClient struct {
	Id                         int64   `form:"id"`
	Name                       string  `form:"name"`
	DefaultCpc                 string  `form:"default_cpc"`
	DailyImpsLimit             string  `form:"daily_imps_limit"`
	Paused                     int64   `form:"paused"`
	ChargeAmount               float64 `form:"charge_amount"`
	ApprovedToCharge           int64   `form:"approved_to_charge"`
	EnhancedReporting          int64
	StripeToken                string
	HourlyReporting            bool
	EnableReportAccountBalance bool
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
	DeletedAt                  time.Time
	BreakDownReportByHour      bool
	ShowMtdSpendInReport       bool
	Emails                     []string `form:"email[]" json:"emails" sql:"-"`
}

func (a AdscoopClient) NotifyClientExpiring() {
	var retData struct {
		Message string
	}

	var doc bytes.Buffer

	var messageEmail *template.Template
	var err error

	if messageEmail, err = template.ParseFiles("templates/email.tmpl"); err != nil {
		return
	}
	retData.Message = "Your account has less than 15% of it's funds available, your campaigns will end soon"

	messageEmail.Execute(&doc, retData)

	go func() {
		gun := mailgun.NewMailgun(configSettting.MailGunDomain, configSettting.MailGunApiKey, configSettting.MailGunPublicApiKey)

		m := mailgun.NewMessage("donotreply <donotreply@mg.adscoops.com>", retData.Message, "Message from Adscoops: "+retData.Message, configSettting.AdminEmail)

		var emails []AdscoopClientEmail

		db.Where("client_id = ?", a.Id).Find(&emails)

		for _, x := range emails {
			m.AddRecipient(x.Email)
		}
		m.SetHtml(doc.String())
		gun.Send(m)
		db.DB().Exec("UPDATE adscoop_clients SET expiration_warning = 1 WHERE id = ?", a.Id)
	}()
}

func (a AdscoopClient) NotifyClientExpired() {
	var retData struct {
		Message string
	}

	var doc bytes.Buffer

	var messageEmail *template.Template
	var err error

	if messageEmail, err = template.ParseFiles("templates/email.tmpl"); err != nil {
		return
	}
	retData.Message = "Your account has ran out of funds available, all active campaigns are paused"

	messageEmail.Execute(&doc, retData)

	go func() {
		gun := mailgun.NewMailgun(configSettting.MailGunDomain, configSettting.MailGunApiKey, configSettting.MailGunPublicApiKey)

		m := mailgun.NewMessage("donotreply <donotreply@mg.adscoops.com>", retData.Message, "Message from Adscoops: "+retData.Message, configSettting.AdminEmail)

		var emails []AdscoopClientEmail

		db.Where("client_id = ?", a.Id).Find(&emails)

		for _, x := range emails {
			m.AddRecipient(x.Email)
		}
		m.SetHtml(doc.String())
		gun.Send(m)
		db.DB().Exec("UPDATE adscoop_clients SET expiration_notice = 1 WHERE id = ?", a.Id)
	}()
}

func (a AdscoopClient) ChargeClient() {

	stripe.Key = configSettting.StripeSecretKey
	chargeParams := &stripe.ChargeParams{
		Amount:   uint64(a.ChargeAmount * 100),
		Currency: "usd",
		Customer: a.StripeToken,
		Desc:     "Charge for " + a.Name,
	}
	ch, err := charge.New(chargeParams)

	if err == nil {
		var act AdscoopClientTransaction
		act.ClientId = a.Id
		act.AmountCharged = a.ChargeAmount
		act.TransactionId = ch.ID
		act.Successful = 1
		act.Attempts = 1
		db.Save(&act)
		db.DB().Exec("UPDATE adscoop_clients SET expiration_warning = 0, expiration_notice = 0 WHERE id = ?", a.Id)
	} else {
		log.Println("Charge could not go through. err: %s", err)
		var retData struct {
			Message string
		}

		var doc bytes.Buffer

		var messageEmail *template.Template
		var err error

		if messageEmail, err = template.ParseFiles("templates/email.tmpl"); err != nil {
			return
		}
		retData.Message = "Card on file could not be charged"

		db.DB().Exec("UPDATE adscoop_clients SET approved_to_charge = 0")

		messageEmail.Execute(&doc, retData)

		go func() {
			gun := mailgun.NewMailgun(configSettting.MailGunDomain, configSettting.MailGunApiKey, configSettting.MailGunPublicApiKey)

			m := mailgun.NewMessage("donotreply <donotreply@mg.adscoops.com>", retData.Message, "Message from Adscoops: "+retData.Message, configSettting.AdminEmail)

			var emails []AdscoopClientEmail

			db.Where("client_id = ?", a.Id).Find(&emails)

			for _, x := range emails {
				m.AddRecipient(x.Email)
			}
			m.SetHtml(doc.String())
			gun.Send(m)
		}()
	}
}

type AdscoopClientTransaction struct {
	Id            int64 `form:"id"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
	ClientId      int64 `form:"client_id"`
	AmountCharged float64
	TransactionId string
	Successful    int64
	Attempts      int64
}

type AdscoopPaymentHash struct {
	Id        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Hash      string
	ClientId  int64
}

type StripeData struct {
	ClientId        int64  `form:"client_id"`
	PaymentHash     string `form:"hash"`
	StripeToken     string `form:"stripeToken"`
	StripeTokenType string `form:"stripeTokenType"`
	StripeEmail     string `form:"stripeEmail"`
}


type AdscoopClientEmail struct {
	Id        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Email     string
	ClientId  int64
}

type AdscoopCsvStats struct {
	Timeslice         time.Time
	CampaignName      string
	Title             string
	Url               string
	Impressions       int
	Verifieds         int
	Loads             int
	TrackingMethod    int64
	ImpressionsString string
	Cost              float64
	CostVerified      float64
	CostLoad          float64
	CostString        string
	Cpc               string
	UniqueIdentifier  string
}
