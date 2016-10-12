package frontendControllers

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	log "github.com/Sirupsen/logrus"
	"app/adscoops.caches"
	"app/structs"
	"github.com/garyburd/redigo/redis"
	"github.com/patrickmn/go-cache"


)



func startCachesTimer() {
	updateCaches()
	ticker := time.NewTicker(time.Second * 60)
	for range ticker.C {
		// log.Println("Going to automatically update caches")
		updateCaches()
	}
}

func updateCaches() {
	log.Println("updating Caches")
	if xu, found := gc.Get("urls_map"); found {
		hashes := xu.(map[uint]bool)

		for idx, y := range hashes {
			if !y {
				continue
			}
			log.Println("Updating cache for urls", idx)
			getAdscoopCampaignUrls(idx, true)
		}
	}
	if xr, found := gc.Get("redirs_map"); found {
		hashes := xr.(map[string]bool)

		log.Printf("redirs_map length: %v\r\n", len(hashes))

		for idx, y := range hashes {
			if !y {
				continue
			}
			log.Println("Updating cache for redir", idx)
			getAdscoopRedirs(idx, true)
		}
	}
	// if xr, found := gc.Get("feeds_map"); found {
	// 	hashes := xr.(map[string]bool)
	//
	// 	for idx, y := range hashes {
	// 		if !y {
	// 			continue
	// 		}
	// 		log.Println("Updating cache for feed", idx)
	// 		getAdscoopFeeds(idx)
	// 	}
	// }
	if xq, found := gc.Get("querystrings_map"); found {
		hashes := xq.(map[uint]bool)

		for idx, y := range hashes {
			if !y {
				continue
			}
			log.Println("Updating cache for querystrings", idx)
			getAdscoopRedirsQueryStrings(idx, true)
		}
	}
	if xc, found := gc.Get("campaigns_map"); found {
		hashes := xc.(map[uint]bool)

		for idx, y := range hashes {
			if !y {
				continue
			}
			log.Println("Updating cache for campaigns", idx)
			getAdscoopRedirsCampaign(idx, true)
		}
	}
}

func updateAdscoopFeeds(hash string) {
	if x, found := gc.Get("feeds_map"); found {
		hashes := x.(map[string]bool)
		hashes[hash] = true
	} else {
		hashes := make(map[string]bool)
		hashes[hash] = true
		gc.Set("redirs_map", hashes, cache.NoExpiration)
	}
}

func getAdscoopFeeds(hash string, fromCache bool) (asf structs.AdscoopFeed) {
	if x, found := gc.Get(fmt.Sprintf(adscoopsCaches.FEEDS_KEY, hash)); found && !fromCache {
		asf = x.(structs.AdscoopFeed)
	} else {
		rp := adscoopsCaches.RedisPool.Get()
		feedString, err := redis.String(rp.Do("GET", fmt.Sprintf(adscoopsCaches.FEEDS_KEY, hash)))
		updateAdscoopFeeds(hash)

		if err != nil {
			asf = adscoopsCaches.LoadAdscoopFeeds(hash)
		} else {
			json.Unmarshal([]byte(feedString), &asf)
			gc.Set(fmt.Sprintf(adscoopsCaches.FEEDS_KEY, hash), asf, cache.NoExpiration)
		}
	}
	return
}

func updateRedirsMap(hash string) {
	if x, found := gc.Get("redirs_map"); found {
		hashes := x.(map[string]bool)
		hashes[hash] = true
	} else {
		log.Println("redir map not found")
		hashes := make(map[string]bool)
		hashes[hash] = true
		gc.Set("redirs_map", hashes, cache.NoExpiration)
	}
}

func getAdscoopRedirs(hash string, fromCache bool) (asc structs.Redirect) {
	// log.Println("looking for redir in cache")
	log.Println("getAdscoopRedirs: ", fmt.Sprintf(adscoopsCaches.REDIRS_KEY, hash))
	if x, found := gc.Get(fmt.Sprintf(adscoopsCaches.REDIRS_KEY, hash)); found && !fromCache {
		log.Println("Loaded redir from cache: ", hash)
		asc = x.(structs.Redirect)
	} else {
		rp := adscoopsCaches.RedisPool.Get()
		defer rp.Close()
		redirString, err := redis.String(rp.Do("GET", fmt.Sprintf(adscoopsCaches.REDIRS_KEY, hash)))
		updateRedirsMap(hash)

		if err != nil {
			asc = adscoopsCaches.LoadAdscoopRedirs(hash)
		} else {
			log.Printf("Loaded redir from redis: %s", hash)
			json.Unmarshal([]byte(redirString), &asc)
			gc.Set(fmt.Sprintf(adscoopsCaches.REDIRS_KEY, hash), asc, cache.NoExpiration)
		}
	}
	return
}

func updateQueryStringsHash(id uint) {
	if x, found := gc.Get("querystrings_map"); found {
		hashes := x.(map[uint]bool)
		hashes[id] = true
	} else {
		hashes := make(map[uint]bool)
		hashes[id] = true
		gc.Set("querystrings_map", hashes, cache.NoExpiration)
	}
}

func getAdscoopRedirsQueryStrings(id uint, fromCache bool) (asqs []structs.RedirectQuerystring) {
	// log.Println("looking for query strings in cache")
	if x, found := gc.Get(fmt.Sprintf(adscoopsCaches.QUERYSTRING_KEY, id)); found && !fromCache {
		asqs = x.([]structs.RedirectQuerystring)
	} else {
		rp := adscoopsCaches.RedisPool.Get()
		defer rp.Close()
		redirsQueryStrings, err := redis.String(rp.Do("GET", fmt.Sprintf(adscoopsCaches.QUERYSTRING_KEY, id)))

		updateQueryStringsHash(id)

		if err != nil {
			asqs = adscoopsCaches.LoadAdscoopQueryStrings(id)
		} else {
			json.Unmarshal([]byte(redirsQueryStrings), &asqs)
			gc.Set(fmt.Sprintf(adscoopsCaches.QUERYSTRING_KEY, id), asqs, cache.NoExpiration)
		}
	}
	return
}

func updateRedirsCampaignHash(id uint) {
	if x, found := gc.Get("campaigns_map"); found {
		hashes := x.(map[uint]bool)
		hashes[id] = true
	} else {
		hashes := make(map[uint]bool)
		hashes[id] = true
		gc.Set("campaigns_map", hashes, cache.NoExpiration)
	}
}

func getAdscoopHostById(id string, fromCache bool) (host AdscoopHost) {
	if x, found := gc.Get(fmt.Sprintf(adscoopsCaches.HOSTS_BY_ID_KEY, id)); found && !fromCache {
		log.Println("getAdscoopHostById in memory")
		host = x.(AdscoopHost)
	} else {
		log.Println("getAdscoopHostById in DB")
		rp := adscoopsCaches.RedisPool.Get()
		defer rp.Close()
		redisHostString, err := redis.String(rp.Do("GET", fmt.Sprintf(adscoopsCaches.HOSTS_BY_ID_KEY), id))

		if err != nil {
			host = loadAdscoopHostById(id)
		} else {
			json.Unmarshal([]byte(redisHostString), &host)
			gc.Set(fmt.Sprintf(adscoopsCaches.HOSTS_BY_ID_KEY, id), host, cache.NoExpiration)
		}
	}
	return
}

func loadAdscoopHostById(id string) (host AdscoopHost) {
	db.Find(&host, id)

	b, err := json.Marshal(host)

	if err == nil {
		rp := adscoopsCaches.RedisPool.Get()
		defer rp.Close()

		rp.Do("SET", fmt.Sprintf(adscoopsCaches.HOSTS_BY_ID_KEY, id), b)
	}

	gc.Set(fmt.Sprintf(adscoopsCaches.HOSTS_BY_ID_KEY, id), host, cache.NoExpiration)
	return
}

func getAdscoopHostByHost(id string, fromCache bool) (host AdscoopHost) {
	if x, found := gc.Get(fmt.Sprintf(adscoopsCaches.HOSTS_BY_HOST_KEY, id)); found && !fromCache {
		host = x.(AdscoopHost)
	} else {
		rp := adscoopsCaches.RedisPool.Get()
		defer rp.Close()
		redisHostString, err := redis.String(rp.Do("GET", fmt.Sprintf(adscoopsCaches.HOSTS_BY_HOST_KEY), id))

		if err != nil {
			host = loadAdscoopHostByHost(id)
		} else {
			json.Unmarshal([]byte(redisHostString), &host)
			gc.Set(fmt.Sprintf(adscoopsCaches.HOSTS_BY_HOST_KEY, id), host, cache.NoExpiration)
		}
	}
	return
}

func loadAdscoopHostByHost(id string) (host AdscoopHost) {
	db.Where("host = ?", id).Find(&host)

	b, err := json.Marshal(host)

	if err == nil {
		rp := adscoopsCaches.RedisPool.Get()
		defer rp.Close()

		rp.Do("SET", fmt.Sprintf(adscoopsCaches.HOSTS_BY_HOST_KEY, id), b)
	}

	gc.Set(fmt.Sprintf(adscoopsCaches.HOSTS_BY_HOST_KEY, id), host, cache.NoExpiration)
	return
}

func getAdscoopRedirsCampaign(id uint, fromCache bool) (asc []structs.Campaign) {
	// log.Println("looking for campaign in cache")
	log.Println("key", fmt.Sprintf(adscoopsCaches.CAMPAIGNS_KEY, id))
	if x, found := gc.Get(fmt.Sprintf(adscoopsCaches.CAMPAIGNS_KEY, id)); found && !fromCache {
		asc = x.([]structs.Campaign)
	} else {
		rp := adscoopsCaches.RedisPool.Get()
		defer rp.Close()

		redirsCampaignString, err := redis.String(rp.Do("GET", fmt.Sprintf(adscoopsCaches.CAMPAIGNS_KEY, id)))
		updateRedirsCampaignHash(id)

		if err != nil {
			asc = adscoopsCaches.LoadAdscoopRedirsCampaign(id)
		} else {
			log.Printf("Got redir from REDIS: %v\r\n", id)
			json.Unmarshal([]byte(redirsCampaignString), &asc)
			gc.Set(fmt.Sprintf(adscoopsCaches.CAMPAIGNS_KEY, id), asc, cache.NoExpiration)
		}
	}
	return
}

func updateCampaignUrlsHash(id uint) {
	if x, found := gc.Get("urls_map"); found {
		hashes := x.(map[uint]bool)
		hashes[id] = true
	} else {
		hashes := make(map[uint]bool)
		hashes[id] = true
		gc.Set("urls_map", hashes, cache.NoExpiration)
	}
}

func getAdscoopCampaignUrls(id uint, fromCache bool) (asu []structs.CampaignUrl) {
	// log.Println("looking for URLs in cache")
	if x, found := gc.Get(fmt.Sprintf(adscoopsCaches.URLS_KEY, id)); found && !fromCache {
		asu = x.([]structs.CampaignUrl)
	} else {

		rp := adscoopsCaches.RedisPool.Get()
		defer rp.Close()

		campaignUrlString, err := redis.String(rp.Do("GET", fmt.Sprintf(adscoopsCaches.URLS_KEY, id)))
		updateCampaignUrlsHash(id)

		if err != nil {
			asu = adscoopsCaches.LoadAdscoopCampaignUrls(id)
		} else {
			json.Unmarshal([]byte(campaignUrlString), &asu)
			gc.Set(fmt.Sprintf(adscoopsCaches.URLS_KEY, id), asu, cache.NoExpiration)
		}
	}
	return
}

type CacheCounter struct {
	Count int
}

type JemaAds struct {
	XMLName    xml.Name      `xml:"Ads"`
	Successful string        `xml:"Successful,attr"`
	Listings   []JemaListing `xml:"Listing"`
}

type JemaListing struct {
	XMLName xml.Name `xml:"Listing"`
	URI     string   `xml:"URI"`
	Cost    string   `xml:"Cost"`
}

type AdNetEntry struct {
	XMLName     xml.Name `xml:"entry"`
	Link        string   `xml:"click_url"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Domain      string   `xml:"site_url"`
	Cost        string   `xml:"bid"`
}

type AdNetEntries struct {
	XMLName xml.Name     `xml:"entries"`
	Entry   []AdNetEntry `xml:"entry"`
}

type AdNetResponse struct {
	XMLName xml.Name     `xml:"response"`
	Entries AdNetEntries `xml:"entries"`
}

type EzangaResponse struct {
	XMLName xml.Name      `xml:"dsxout"`
	Results EzangaResults `xml:"results"`
}

type EzangaResults struct {
	XMLName  xml.Name        `xml:"results"`
	Listings []EzangaListing `xml:"listing"`
}

type EzangaListing struct {
	XMLName xml.Name `xml:"listing"`
	Url     string   `xml:"url"`
	Bid     string   `xml:"bid"`
}

func getJemaXmlUrl(urls []structs.CampaignUrl, id uint, req *http.Request, asr structs.Redirect) (redirUrl *url.URL, uID uint, cpc string) {
	redirUrl = nil
	if len(urls) != 1 {
		fmt.Println("not enough urls")
		return
	}
	timeout := time.Duration(4 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	ipInfo := req.RemoteAddr
	ipInfoSplit := strings.SplitN(string(ipInfo), ":", -1)
	ipSlice := ipInfoSplit[0 : len(ipInfoSplit)-1]
	ip := strings.Join(ipSlice, "")

	if req.Header.Get("Fastly-Client-IP") != "" {
		ip = req.Header.Get("Fastly-Client-IP")
	}

	ua := req.UserAgent()

	remoteUrl, _ := url.Parse(urls[0].Url)

	qsr, err := url.Parse(req.URL.String())

	if err == nil && qsr.Query().Encode() != "" && remoteUrl != nil {

		if asr.StripQueryString {

			qs := remoteUrl.Query()

			var asqs []structs.RedirectQuerystring

			asqs = getAdscoopRedirsQueryStrings(asr.ID, false)

			for _, y := range asqs {
				if req.FormValue(y.QueryStringKey) == "" {
					continue
				}
				qs.Set(y.QueryStringKey, req.FormValue(y.QueryStringKey))
			}

			remoteUrl.RawQuery = qs.Encode()

		} else {

			pRedirUrl, _ := url.ParseQuery(remoteUrl.RawQuery)
			qry := qsr.Query()
			for x, y := range pRedirUrl {
				for _, yy := range y {
					fmt.Println("value", yy)
					qry.Set(x, yy)
				}
			}
			qsr.RawQuery = qry.Encode()

			remoteUrl.RawQuery = qsr.Query().Encode()
		}
	}

	qs := remoteUrl.Query()

	qs.Add("uip", ip)
	qs.Add("ua", ua)

	remoteUrl.RawQuery = qs.Encode()

	fmt.Println("remoteUrl", remoteUrl.String())

	reqHTTP, err := http.NewRequest("GET", remoteUrl.String(), nil)

	if err != nil {
		fmt.Println("err", err)
		return
	}

	res, err := client.Do(reqHTTP)

	if err != nil {
		fmt.Println("err", err)
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println("err", err)
		return
	}

	var retXML JemaAds

	err = xml.Unmarshal(body, &retXML)

	if err != nil {
		fmt.Println("err", err)
		return
	}

	if len(retXML.Listings) == 0 {
		fmt.Println("not enough urls for " + remoteUrl.String())
		return
	}

	fmt.Println("url", retXML.Listings[0].URI)

	redirUrl, err = url.Parse(retXML.Listings[0].URI)
	if err != nil {
		// log.Println("url is invalid")
		redirUrl = nil
		return
	}

	uID = urls[0].ID
	cpc = retXML.Listings[0].Cost
	return
}

func getAdNetXmlUrl(urls []structs.CampaignUrl, id uint, req *http.Request, asr structs.Redirect) (redirUrl *url.URL, uID uint, cpc string) {
	redirUrl = nil
	if len(urls) != 1 {
		fmt.Println("not enough urls")
		return
	}
	timeout := time.Duration(4 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	ipInfo := req.RemoteAddr
	ipInfoSplit := strings.SplitN(string(ipInfo), ":", -1)
	ipSlice := ipInfoSplit[0 : len(ipInfoSplit)-1]
	ip := strings.Join(ipSlice, "")

	if req.Header.Get("Fastly-Client-IP") != "" {
		ip = req.Header.Get("Fastly-Client-IP")
	}

	ua := req.UserAgent()

	remoteUrl, _ := url.Parse(urls[0].Url)

	qsr, err := url.Parse(req.URL.String())

	if err == nil && qsr.Query().Encode() != "" && remoteUrl != nil {

		if asr.StripQueryString {

			qs := remoteUrl.Query()

			var asqs []structs.RedirectQuerystring

			asqs = getAdscoopRedirsQueryStrings(asr.ID, false)

			for _, y := range asqs {
				if req.FormValue(y.QueryStringKey) == "" {
					continue
				}
				qs.Set(y.QueryStringKey, req.FormValue(y.QueryStringKey))
			}

			remoteUrl.RawQuery = qs.Encode()

		} else {

			pRedirUrl, _ := url.ParseQuery(remoteUrl.RawQuery)
			qry := qsr.Query()
			for x, y := range pRedirUrl {
				for _, yy := range y {
					fmt.Println("value", yy)
					qry.Set(x, yy)
				}
			}
			qsr.RawQuery = qry.Encode()

			remoteUrl.RawQuery = qsr.Query().Encode()
		}
	}

	qs := remoteUrl.Query()

	qs.Add("client_ip", ip)
	qs.Add("ua", ua)

	remoteUrl.RawQuery = qs.Encode()

	fmt.Println("remoteUrl", remoteUrl.String())

	reqHTTP, err := http.NewRequest("GET", remoteUrl.String(), nil)

	if err != nil {
		fmt.Println("err", err)
		return
	}

	res, err := client.Do(reqHTTP)

	if err != nil {
		fmt.Println("err", err)
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	body = bytes.Replace(body, []byte("encoding=\"ISO-8859-1\""), []byte("encoding=\"utf-8\""), -1)

	if err != nil {
		fmt.Println("err", err)
		return
	}

	var retXML AdNetResponse

	err = xml.Unmarshal(body, &retXML)

	if err != nil {
		fmt.Println("err", err)
		return
	}

	if len(retXML.Entries.Entry) == 0 {
		fmt.Println("not enough urls for " + remoteUrl.String())
		return
	}

	fmt.Println("url", retXML.Entries.Entry[0].Link)

	redirUrl, err = url.Parse(retXML.Entries.Entry[0].Link)
	if err != nil {
		// log.Println("url is invalid")
		redirUrl = nil
		return
	}

	uID = urls[0].ID
	cpc = retXML.Entries.Entry[0].Cost
	return
}

func getEzangaXmlUrl(urls []structs.CampaignUrl, id uint, req *http.Request, asr structs.Redirect) (redirUrl *url.URL, uID uint, cpc string) {
	redirUrl = nil
	if len(urls) != 1 {
		fmt.Println("not enough urls")
		return
	}
	timeout := time.Duration(4 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	ipInfo := req.RemoteAddr
	ipInfoSplit := strings.SplitN(string(ipInfo), ":", -1)
	ipSlice := ipInfoSplit[0 : len(ipInfoSplit)-1]
	ip := strings.Join(ipSlice, "")

	if req.Header.Get("Fastly-Client-IP") != "" {
		ip = req.Header.Get("Fastly-Client-IP")
	}

	ua := req.UserAgent()

	remoteUrl, _ := url.Parse(urls[0].Url)

	qsr, err := url.Parse(req.URL.String())

	if err == nil && qsr.Query().Encode() != "" && remoteUrl != nil {

		if asr.StripQueryString {

			qs := remoteUrl.Query()

			var asqs []structs.RedirectQuerystring

			asqs = getAdscoopRedirsQueryStrings(asr.ID, false)

			for _, y := range asqs {
				if req.FormValue(y.QueryStringKey) == "" {
					continue
				}
				qs.Set(y.QueryStringKey, req.FormValue(y.QueryStringKey))
			}

			remoteUrl.RawQuery = qs.Encode()

		} else {

			pRedirUrl, _ := url.ParseQuery(remoteUrl.RawQuery)
			qry := qsr.Query()
			for x, y := range pRedirUrl {
				for _, yy := range y {
					fmt.Println("value", yy)
					qry.Set(x, yy)
				}
			}
			qsr.RawQuery = qry.Encode()

			remoteUrl.RawQuery = qsr.Query().Encode()
		}
	}

	qs := remoteUrl.Query()

	qs.Set("ip", ip)
	qs.Set("ua", ua)

	qs.Del("rf")

	remoteUrl.RawQuery = qs.Encode()

	qsone := url.URL{}

	qryone := qsone.Query()
	qryone.Add("rf", req.Referer())

	rfString := "&" + qryone.Encode()

	fmt.Println("encode", rfString)

	remoteUrl.RawQuery += rfString

	fmt.Println("remoteUrl", remoteUrl.String())

	reqHTTP, err := http.NewRequest("GET", remoteUrl.String(), nil)

	if err != nil {
		fmt.Println("err", err)
		return
	}

	res, err := client.Do(reqHTTP)

	if err != nil {
		fmt.Println("err", err)
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println("err", err)
		return
	}

	var retXML EzangaResponse

	err = xml.Unmarshal(body, &retXML)

	if err != nil {
		fmt.Println("err", err)
		return
	}

	if len(retXML.Results.Listings[0].Url) == 0 {
		fmt.Println("not enough urls for " + remoteUrl.String())
		return
	}

	fmt.Println("url", retXML.Results.Listings[0].Url)

	redirUrl, err = url.Parse(retXML.Results.Listings[0].Url)
	if err != nil {
		redirUrl = nil
		return
	}

	uID = urls[0].ID
	cpc = retXML.Results.Listings[0].Bid
	return
}

func getCurrentUrl(urls []structs.CampaignUrl, id uint, req *http.Request, asr structs.Redirect) (redirUrl *url.URL, uID uint) {
	uIdx := &CacheCounter{Count: 0}
	if x, found := gc.Get(fmt.Sprintf(adscoopsCaches.URLS_KEY+"index", id)); found {
		uIdx = x.(*CacheCounter)
	} else {
		gc.Set(fmt.Sprintf(adscoopsCaches.URLS_KEY+"index", id), uIdx, cache.NoExpiration)
	}

	// log.Println("uIdx", uIdx.Count)
	// log.Println("len(urls)", len(urls))

	uIdx.Count++

	if uIdx.Count >= len(urls) {
		// log.Println("resetting index")
		uIdx.Count = 0
	}

	// log.Println("uIdx", uIdx.Count)
	// log.Println("len(urls)", len(urls))

	if len(urls) == 0 {
		log.Println("no URL's available in cache:", fmt.Sprintf(adscoopsCaches.URLS_KEY+"index", id))
		return
	}

	y := urls[uIdx.Count]

	var err error

	redirUrl, err = url.Parse(y.Url)
	if err != nil {
		// log.Println("url is invalid")
		redirUrl = nil
		return
	}

	qsr, err := url.Parse(req.URL.String())

	if err == nil && qsr.Query().Encode() != "" {

		if asr.StripQueryString {

			qs := redirUrl.Query()

			var asqs []structs.RedirectQuerystring

			asqs = getAdscoopRedirsQueryStrings(asr.ID, false)

			for _, y := range asqs {
				if req.FormValue(y.QueryStringKey) == "" {
					continue
				}
				qs.Set(y.QueryStringKey, req.FormValue(y.QueryStringKey))
			}

			redirUrl.RawQuery = qs.Encode()

		} else {

			pRedirUrl, _ := url.ParseQuery(redirUrl.RawQuery)
			qry := qsr.Query()
			for x, y := range pRedirUrl {
				fmt.Println("key", x)
				for _, yy := range y {
					fmt.Println("value", yy)
					qry.Set(x, yy)
				}
			}
			qsr.RawQuery = qry.Encode()

			redirUrl.RawQuery = qsr.Query().Encode()
		}
	}

	qry := redirUrl.Query()
	qry.Del("_ast")
	qry.Del("_asts")
	redirUrl.RawQuery = qry.Encode()
	// log.Println("url", redirUrl.String())
	uID = y.ID
	return
}
