package emailer

import (
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/mailgun/mailgun-go"
)

const MailGunDomain = "mg.adscoops.com"
const MailGunApiKey = "key-22a1cbbf6bbbec6d1f4a3fd7219f1486"
const MailGunPublicApiKey = "pubkey-40277136acd57b2aea204afab1ccba0c"

const Address = "donotreply <donotreply@mg.adscoops.com>"

type Emailer struct {
	Title   string
	Message string
	Emails  []string
}

func (e *Emailer) Send(title string, message string, emails []string) (err error) {
	e.Title = title
	e.Message = message
	e.Emails = emails

	log.Infof("Going to send email.  Emails: %#s, Title: %s, Message: %s", e.Emails, e.Title, e.Message)

	if strings.ToLower(os.Getenv("GO_ENVIRONMENT")) == "development" {
		return nil
	}

	gun := mailgun.NewMailgun(MailGunDomain, MailGunApiKey, MailGunPublicApiKey)
	m := mailgun.NewMessage(Address, title, message)

	for _, e := range emails {
		m.AddRecipient(e)
	}

	_, _, err = gun.Send(m)
	return
}
