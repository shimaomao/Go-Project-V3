package app

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"bitbucket.com/barrettbsi/broadvid-adscoops-shared/adscoopUtils"
	"bitbucket.com/barrettbsi/broadvid-adscoops-shared/structs"
	"github.com/BurntSushi/toml"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/secure"
	"github.com/martini-contrib/sessions"
)

var (
	dbConn   = "GO_DATABASE_CONN"
	dbTable  = "adscoops"
	tomlFile = "config.toml"
	asUtils  *adscoopUtils.UtilManager
	config   tomlConfig
	db       gorm.DB
)

func loadToml() {
	if _, err := toml.DecodeFile(tomlFile, &config); err != nil {
		fmt.Println(err)
		return
	}
}

func App(isTesting bool) (m *martini.ClassicMartini, err error) {
	loadToml()

	db, err = gorm.Open("mysql", config.SqlConnection)

	if err != nil {
		log.Println("err", err)
		return
	}

	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	go processFuncs()
	go keepWSAlive()

	asUtils = adscoopUtils.NewUtilManager(&db)

	m = martini.Classic()

	store := sessions.NewCookieStore([]byte("asdlg908dgslkasdgn"))
	m.Use(sessions.Sessions("ascp", store))

	if !isTesting {
		martini.Env = martini.Prod
		m.Use(secure.Secure(secure.Options{
			AllowedHosts:         []string{"localhost:3003", "client-dash.adscoops.com", "localhost:3009"},
			SSLProxyHeaders:      map[string]string{"X-Forwarded-Proto": "https"},
			STSSeconds:           315360000,
			STSIncludeSubdomains: true,
			FrameDeny:            true,
			ContentTypeNosniff:   true,
			BrowserXssFilter:     true,
			// ContentSecurityPolicy: "default-src 'self'",
		}))
	}

	priFuncs := template.FuncMap{
		"LoadTemplate": func(name string, data interface{}) (ret template.HTML, err error) {
			var buf bytes.Buffer
			t := template.Must(template.ParseFiles("templates/" + name + ".tmpl"))
			err = t.Execute(&buf, data)
			ret = template.HTML(buf.String())
			return
		},
		"FormatReadableTime": func(dateTime time.Time) (ret template.HTML, err error) {
			ret = template.HTML(dateTime.Format("01/02/2006 03:04 PM"))
			return
		},
		"FormatReadableTimeLosAngeles": func(dateTime time.Time) (ret template.HTML, err error) {
			pst := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(),
				dateTime.Hour(), dateTime.Minute(), 0, 0, time.UTC)
			location, _ := time.LoadLocation("America/Los_Angeles")

			dateTime = pst.In(location)
			ret = template.HTML(dateTime.Format("01/02/2006 03:04 PM"))
			return
		},
		"Addition": func(numone int, numtwo int) int {
			return numone + numtwo
		},
	}

	funcMap := []template.FuncMap{
		priFuncs,
	}
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
		Funcs:  funcMap,
	}))

	controllersSetup(m)

	m.Get("/wsupdates", requireLogin, func(w http.ResponseWriter, r *http.Request, user *UserWithPolicy) {
		uid := fmt.Sprintf("%v", user.ClientID)
		ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(w, "Not a websocket handshake", 400)
			return
		} else if err != nil {
			log.Println(err)
			return
		}
		client := ws.RemoteAddr()
		sockCli := ClientConn{uid, ws, client}
		addClient(sockCli)

		for {
			messageType, p, err := ws.ReadMessage()
			if err != nil {
				deleteClient(sockCli)
				return
			}
			broadcastMessage(messageType, p, uid)
		}
	})

	return
}

func processFuncs() {

	go setupCampaignCaches()
	ticker := time.NewTicker(time.Minute * 5)

	for range ticker.C {
		go setupCampaignCaches()
	}
}

func keepWSAlive() {
	ticker := time.NewTicker(time.Second * 20)

	for range ticker.C {
		broadcastMessageToAll(1, []byte("hi"))
	}
}

func requireLogin(s sessions.Session, r render.Render, c martini.Context) {
	v := s.Get("user_id")

	if v == nil {
		r.Redirect("/login")
		return
	}

	var user UserWithPolicy

	id := v.(uint)

	err := db.Select("adscoop_client_users.*").Joins("JOIN adscoop_clients ON adscoop_clients.id = adscoop_client_users.client_id").Where("adscoop_clients.enable_client_login = 1 AND adscoop_client_users.id = ?", id).Find(&user).Error

	if err != nil {

		ops := sessions.Options{
			MaxAge: -1,
		}
		s.Options(ops)
		s.Clear()
		r.Redirect("/login")
		return
	}

	var client structs.AdscoopClient
	db.Where("id = ?", user.ClientID).Find(&client)

	user.Name = client.Name

	db.Find(&user.Policy, user.UserPolicyID)

	c.Map(&user)
}
