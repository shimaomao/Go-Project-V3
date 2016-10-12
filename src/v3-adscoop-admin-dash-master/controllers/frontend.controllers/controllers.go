package frontendControllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"
	//"configSettting"
	"app/structs"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
)

var store = sessions.NewCookieStore([]byte("349qg8nmzxd;glksagh938qaghas;glksahgd"))

var templates = template.Must(template.ParseFiles("templates/tracking.tmpl.js", "templates/loadClient.tmpl.js"))

var TSC structs.TempStatsContainer
var r *render.Render
var td = ClickStats{}

func controllerTrackingJs(c *gin.Context) {

	req := c.Request
	//log.Println("gin context",c)
	var retData struct {
		Host string
	}
	retData.Host = req.Host
	c.Header("Content-Type", "application/javascript")

	c.HTML(http.StatusOK, "tracking.tmpl.js", retData)

}

func controllerLoadClient(c *gin.Context) {

	var retData struct {
		Callback             string
		Host                 string
		ET                   bool
		MinTimeout           uint
		MaxTimeout           uint
		Redir                string
		EnableUnloadTracking bool
	}

	req := c.Request

	data := req.URL.Query().Get("_ast")
	ast := strings.Split(data, "_")
	host := getAdscoopHostById(ast[0], false)
	retData.Host = host.Host

	if retData.Host == "" {
		retData.Host = "send.adscoops.com"
	}

	session, _ := store.Get(req, "asdata")

	if len(ast) == 3 {
		if val, ok := session.Values["hash"]; ok && val.(string) == ast[2] {
			retData.ET = true
			redir := getAdscoopRedirs(val.(string), false)
			retData.MinTimeout = redir.Min
			retData.MaxTimeout = redir.Max
			log.Println("min", redir.Min)
			redirUrl, _ := url.Parse(req.Referer())
			retData.Redir = "http://" + retData.Host + "/u/" + val.(string) + "?" + redirUrl.RawQuery
		}
	}
	retData.Callback = req.FormValue("callback")

	if val, ok := session.Values["enable_unload_tracking"]; ok {
		retData.EnableUnloadTracking = val.(bool)
		log.Println("LOG LOG LOG enable_unload_tracking", val)
	}

	c.HTML(http.StatusOK, "loadClient.tmpl.js", retData)

}

func controllerTrackEngagement(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	c.Header("Expires", "0")
	c.Header("Content-type", "image/jpg")

	req := c.Request
	session, _ := store.Get(req, "asdata")

	var ast AdscoopTracking

	if val, ok := session.Values["cpc"]; !ok {
		return
	} else {
		ast.Cpc = val.(string)
	}

	if val, ok := session.Values["urlid"]; !ok {
		return
	} else {
		ast.UrlId = val.(uint)
	}

	if val, ok := session.Values["redirectid"]; !ok {
		return
	} else {
		ast.RedirectId = val.(uint)
	}

	ast.UniqueIdentifier = req.FormValue("player_instance_id")

	go ast.TrackEngagement()
	go td.Add(&ast, 1)
	return
}

func controllerTrackLoad(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	c.Header("Expires", "0")
	c.Header("Content-type", "image/jpg")

	req := c.Request

	session, _ := store.Get(req, "asdata")

	var ast AdscoopTracking

	if val, ok := session.Values["cpc"]; !ok {
		return
	} else {
		ast.Cpc = val.(string)
	}

	if val, ok := session.Values["urlid"]; !ok {
		return
	} else {
		ast.UrlId = val.(uint)
	}

	if val, ok := session.Values["redirectid"]; !ok {
		return
	} else {
		ast.RedirectId = val.(uint)
	}

	ast.UniqueIdentifier = req.FormValue("player_instance_id")

	go ast.TrackLoad()
	go td.Add(&ast, 2)
	return
}

func controllerTrackTimeOnSite(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	c.Header("Expires", "0")
	c.Header("Content-type", "image/jpg")

	req := c.Request

	session, _ := store.Get(req, "asdata")

	var ast AdscoopTracking

	if val, ok := session.Values["cpc"]; !ok {
		return
	} else {
		ast.Cpc = val.(string)
	}

	if val, ok := session.Values["urlid"]; !ok {
		return
	} else {
		ast.UrlId = val.(uint)
	}

	if val, ok := session.Values["redirectid"]; !ok {
		return
	} else {
		ast.RedirectId = val.(uint)
	}

	ast.UniqueIdentifier = req.FormValue("player_instance_id")

	go ast.TrackLoad()
	tos, _ := strconv.ParseFloat(req.FormValue("_tos"), 64)
	ast.TimeOnSite = tos
	go td.Add(&ast, 3)
	return
}

func controllerFeed(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	c.Header("Expires", "0")

	req := c.Request

	params := mux.Vars(req)

	feed := getAdscoopFeeds(params["hash"], false)

	if feed.Id == 0 {
		if params["type"] == "json" {
			r.JSON(c.Writer, http.StatusInternalServerError, ErrorXml{Message: "No links available at this time"})
		} else if params["type"] == "xml" {
			r.XML(c.Writer, http.StatusInternalServerError, ErrorXml{Message: "No links available at this time"})
		}

		return
	}

	var retXML ReturnXml

	for _, y := range feed.Redirects {
		redir := getAdscoopRedirs(y.Hash, false)

		if redir.ID == 0 {
			continue
		}

		campaign := getAdscoopRedirsCampaign(redir.ID, false)

		if len(campaign) == 0 {
			continue
		}

		var link ReturnLinkXml

		link.Cpc = campaign[0].Cpc
		link.Title = redir.Name
		link.Link = "http://" + req.Host + "/r/" + redir.Hash

		retXML.Link = append(retXML.Link, &link)
	}

	if params["type"] == "json" {
		r.JSON(c.Writer, http.StatusOK, retXML)
	} else if params["type"] == "xml" {
		r.XML(c.Writer, http.StatusOK, retXML)
	}
}

func controllerValidUser(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	c.Header("Expires", "0")
	req := c.Request
	retData := LoadRetData(c)

	session, _ := store.Get(req, "asdata_flash")
	session.AddFlash("valid_user_" + retData.AdscoopRedirect.Hash)

	session.Save(req, c.Writer)

	http.Redirect(c.Writer, req, req.Referer(), http.StatusTemporaryRedirect)
}

func controllerLastUpdated(c *gin.Context) {
	retData := LoadRetData(c)

	r.Text(c.Writer, http.StatusOK, fmt.Sprintf("%v", retData.AdscoopRedirect.UpdatedAt.Unix()))
}

func controllerRedirect(c *gin.Context) {
	var trackingValue uint32
	req := c.Request
	trackingValue = setTrackingValue(trackingValue, req)

	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	c.Writer.Header().Set("Expires", "0")
	retData := LoadRetData(c)
	fmt.Println(fmt.Sprintf("hash: %s referrer: %s", retData.AdscoopRedirect.Hash, req.Referer()))

	if retData.AdscoopRedirect.ForceHost != "0" && req.Host != retData.AdscoopRedirect.ForceHostString &&
		retData.AdscoopRedirect.ForceHostString != "" {
		url := req.URL
		url.Host = retData.AdscoopRedirect.ForceHostString
		http.Redirect(c.Writer, req, ("http://" + strings.Trim(url.String(), "//")), http.StatusTemporaryRedirect)
	}

	session, _ := store.Get(req, "asdata")

	sessionf, _ := store.Get(req, "asdata_flash")

	var flashData string

	if flashes := sessionf.Flashes(); len(flashes) > 0 {
		flashData = flashes[0].(string)
		sessionf.Save(req, c.Writer)
	}

	if retData.AdscoopRedirect.BapiScoring != 0 && flashData != "valid_user_"+retData.AdscoopRedirect.Hash {
		var rd struct {
			RedirUrl        string
			Name            string
			BapiScoring     uint
			AdscoopRedirect structs.Redirect
		}
		rd.AdscoopRedirect = retData.AdscoopRedirect
		rd.BapiScoring = retData.AdscoopRedirect.BapiScoring

		redirUrl := "/v/" + retData.AdscoopRedirect.Hash

		if req.URL.Query().Encode() != "" {
			redirUrl = redirUrl + "?" + req.URL.Query().Encode()
		}

		rd.Name = retData.AdscoopRedirect.Name
		rd.RedirUrl = redirUrl

		ref, err := url.Parse(req.Referer())

		fmt.Println("before scored http_ref", req.Referer(), retData.AdscoopRedirect.Hash, req.Header.Get("Fastly-Client-IP"), "flash data: ", flashData, "ua: ", req.UserAgent())

		if err == nil {
			session.Values["http_ref"] = ref.Host
			session.Save(req, c.Writer)
		}

		r.HTML(c.Writer, http.StatusOK, "score", rd)
		return
	}

	if retData.AdscoopRedirect.ID == 0 {
		return
	}

	log.Println("iframe", retData.AdscoopRedirect.Iframe)

	if retData.AdscoopRedirect.Iframe == 1 {

		if retData.AdscoopRedirect.LockWhitelistId != "0" {
			retData.AllowRefresh = false
			var httpRef string

			if val, ok := session.Values["http_ref"]; ok {
				httpRef = val.(string)
				fmt.Println("ref found in session", httpRef, retData.AdscoopRedirect.Hash)
			}

			if httpRef == "" || retData.AdscoopRedirect.BapiScoring == 0 {
				ref, err := url.Parse(req.Referer())
				if err == nil {
					httpRef = ref.Host
				}
			}

			fmt.Println("after scored http_ref", httpRef, retData.AdscoopRedirect.Hash, req.Header.Get("Fastly-Client-IP"))

			for _, y := range retData.AdscoopRedirect.LockWhitelistUrls {
				if strings.Contains(httpRef, y.Url) {
					if retData.AdscoopRedirect.LockWhitelistReverse {
						retData.AllowRefresh = false
						break
					}
					retData.AllowRefresh = true
				}
			}
		}

		if retData.AdscoopRedirect.LockUseragentId != "0" {
			for _, y := range retData.AdscoopRedirect.LockUseragents {
				if strings.Contains(req.UserAgent(), y.Useragent) {
					if retData.AdscoopRedirect.LockUseragentReverse {
						retData.AllowRefresh = false
						break
					}
					retData.AllowRefresh = true
				}
			}
		}

		qsr, err := url.Parse(req.URL.String())
		var qsi = url.Values{}

		if err == nil && qsr.Query().Encode() != "" {
			if retData.AdscoopRedirect.StripQueryString {
				var asqs []structs.RedirectQuerystring

				asqs = getAdscoopRedirsQueryStrings(retData.AdscoopRedirect.ID, false)

				for _, y := range asqs {
					if req.FormValue(y.QueryStringKey) == "" {
						continue
					}
					qsi.Add(y.QueryStringKey, req.FormValue(y.QueryStringKey))
				}
				retData.QueryString = qsi.Encode()
			} else {
				qsi = qsr.Query()
			}
		}
		tvString := ""

		if trackingValue != 0 {
			tvString = fmt.Sprintf("%v", trackingValue)
		}

		if retData.AdscoopRedirect.BbsiPath == "" {

			qsi.Set("_ast", fmt.Sprintf("%s_0", tvString))
		} else {
			qsi.Set("_ast", fmt.Sprintf("%s_1", tvString))
		}
		retData.QueryString = qsi.Encode()

		fmt.Println("trackingValue prior to print", trackingValue)

		r.HTML(c.Writer, http.StatusOK, "iframe", retData)
		return
	}

	redirectUrl(c.Writer, req, &retData, trackingValue)
}

func controllerRedirectUrl(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	c.Header("Expires", "0")
	var trackingValue uint32
	req := c.Request
	trackingValue = setTrackingValue(trackingValue, req)
	retData := LoadRetData(c)
	retData.AdscoopRedirect.BustIframe = false
	fmt.Println(fmt.Sprintf("hash: %s referrer: %s", retData.AdscoopRedirect.Hash, req.Referer()))
	if retData.AdscoopRedirect.ID == 0 {
		return
	}

	redirectUrl(c.Writer, req, &retData, trackingValue)

}

func redirectUrl(w http.ResponseWriter, req *http.Request, retData *RetData, trackingValue uint32) {

	redirID := retData.AdscoopRedirect.ID

	var ascs []structs.Campaign

	ascs = getAdscoopRedirsCampaign(redirID, false)

	if len(ascs) == 0 {
		var ast AdscoopTracking
		ast.RedirectId = retData.AdscoopRedirect.ID
		ast.UrlId = 0
		ast.Cpc = "0"
		go td.Add(&ast, 0)
		http.Error(w,
			fmt.Sprint("Unable to serve content at this time"),
			http.StatusInternalServerError)
		return
	}
	var redirUrl *url.URL
	var uID uint
	var cpc string

	var xCpc string
	var enableUnloadTracking bool

	//var trackingMethod string

	for _, asc := range ascs {
		var asus []structs.CampaignUrl

		asus = getAdscoopCampaignUrls(asc.ID, false)

		if asc.Type == 0 || asc.Type == 2 {
			fmt.Println("is regular links campaign")
			redirUrl, uID = getCurrentUrl(asus, asc.ID, req, retData.AdscoopRedirect)
		}

		if asc.Type == 1 {
			if asc.XmlType == 0 {
				redirUrl, uID, xCpc = getJemaXmlUrl(asus, asc.ID, req, retData.AdscoopRedirect)
			}
			if asc.XmlType == 1 {
				redirUrl, uID, xCpc = getAdNetXmlUrl(asus, asc.ID, req, retData.AdscoopRedirect)
			}
			if asc.XmlType == 2 {
				redirUrl, uID, xCpc = getEzangaXmlUrl(asus, asc.ID, req, retData.AdscoopRedirect)
			}
		}

		if redirUrl != nil {
			//trackingMethod = asc.TrackingMethod
			enableUnloadTracking = asc.EnableUnloadTracking
			cpc = asc.Cpc
			if xCpc != "" {
				cpc = xCpc
			}
			break
		}
		fmt.Println("going onto the next campaign")
	}

	if redirUrl == nil {
		var ast AdscoopTracking
		ast.RedirectId = retData.AdscoopRedirect.ID
		ast.UrlId = 0
		ast.Cpc = "0"
		go td.Add(&ast, 0)
		http.Error(w,
			fmt.Sprintf("Cannot point to URL at this time"),
			http.StatusInternalServerError)
		return
	}

	var ast AdscoopTracking

	fmt.Println("host", req.Host)

	ast.UrlId = uID
	ast.Cpc = cpc
	ast.RedirectId = retData.AdscoopRedirect.ID

	//UrlTrackingMethod.Lock.Lock()
	//UrlTrackingMethod.List[uID] = trackingMethod
	//UrlTrackingMethod.Lock.Unlock()

	if req.FormValue("dnt") != "true" {
		go td.Add(&ast, 0)
	}

	go ast.Track()

	session, _ := store.Get(req, "asdata")
	session.Values["urlid"] = ast.UrlId
	session.Values["enable_unload_tracking"] = enableUnloadTracking
	session.Values["redirectid"] = retData.AdscoopRedirect.ID
	session.Values["cpc"] = cpc
	session.Save(req, w)

	var queries = redirUrl.Query()

	tvString := ""

	if trackingValue != 0 {
		tvString = fmt.Sprintf("%v", trackingValue)
	}

	if retData.AdscoopRedirect.BbsiPath == "" {
		queries.Set("_ast", fmt.Sprintf("%s_0", tvString))
	} else {
		queries.Set("_ast", fmt.Sprintf("%s_1", tvString))
	}

	if retData.AdscoopRedirect.Iframe == 2 {
		ast := queries.Get("_ast")

		ast += "_" + retData.AdscoopRedirect.Hash
		session, _ := store.Get(req, "asdata")
		session.Values["hash"] = retData.AdscoopRedirect.Hash
		session.Save(req, w)
		queries.Set("_ast", ast)
	}

	if queries.Get("url") != "" {
		queryUrl, _ := url.Parse(queries.Get("url"))

		quQs := queryUrl.Query()

		if quQs.Get("_ast") != "" {
			quQs.Set("_ast", queries.Get("_ast"))

			queryUrl.RawQuery = quQs.Encode()

			queries.Set("url", queryUrl.String())
		}
	}

	redirUrl.RawQuery = queries.Encode()

	fmt.Println("redirUrl", redirUrl.String())

	if retData.AdscoopRedirect.RedirType == 0 {
		// log.Println("Server-side redir")
		http.Redirect(w, req, redirUrl.String(), http.StatusTemporaryRedirect)
		return
	}

	if retData.AdscoopRedirect.RedirType != 0 {
		// log.Println("Client-side redir")
		redir := retData.AdscoopRedirect
		var retData struct {
			Url             *url.URL
			AdscoopRedirect structs.Redirect
		}

		retData.AdscoopRedirect = redir
		retData.Url = redirUrl
		r.HTML(w, http.StatusOK, "clientRedir", retData)
		return
	}
}

func LoadRetData(c *gin.Context) (retData RetData) {
	req := c.Request
	params := mux.Vars(req)

	retData.AllowRefresh = true

	if params["hash"] == "" {
		var ast AdscoopTracking
		ast.RedirectId = retData.AdscoopRedirect.ID
		ast.UrlId = 0
		ast.Cpc = "0"
		go td.Add(&ast, 0)
		http.Error(c.Writer,
			fmt.Sprint("Hash is empty"),
			http.StatusInternalServerError)
		return
	}

	retData.AdscoopRedirect = getAdscoopRedirs(params["hash"], false)

	if retData.AdscoopRedirect.ID == 0 {
		var ast AdscoopTracking
		ast.RedirectId = retData.AdscoopRedirect.ID
		ast.UrlId = 0
		ast.Cpc = "0"
		go td.Add(&ast, 0)
		http.Error(c.Writer,
			fmt.Sprint("Hash not found"),
			http.StatusNotFound)
		return
	}

	return
}

func setTrackingValue(trackingValue uint32, req *http.Request) uint32 {
	host := getAdscoopHostByHost(req.Host, false)
	trackingValue = uint32(host.Id)
	return trackingValue
}
