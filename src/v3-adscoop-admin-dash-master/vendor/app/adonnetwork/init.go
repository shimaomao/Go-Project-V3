package adonnetwork

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

const (
	accountEmail = "daniel.aharonoff@broadscaler.com"
	apiToken     = "pkjrGU5OuupHm3_fbt2DNQi3BBzE5Csx"

	apiJsonReturn = "%s%s?format=json&username=%s&api_token=%s"
	apiPutReturn  = "%s%s?username=%s&api_token=%s"

	// Base URL for API
	apiBaseUrl = "https://account.adonnetwork.com/api/v2/"

	// Endpoints
	apiGetAllCampaigns   = "campaign/"
	apiGetSingleCampaign = "campaign/%v/"
)

var (
	AdscoopsDB *gorm.DB
)

func makeUrl(endpoint string) string {
	return fmt.Sprintf(apiJsonReturn, apiBaseUrl, endpoint, accountEmail, apiToken)
}

func makePutUrl(endpoint string) string {
	return fmt.Sprintf(apiPutReturn, apiBaseUrl, endpoint, accountEmail, apiToken)
}

// https://account.adonnetwork.com/api/v2/campaign/?format=json&username=user@somedomain.com&api_token=jkcB5wjaNGEjrVrAgZQhsEO64usjoEdS

func AutoMigrate() error {
	return AdscoopsDB.AutoMigrate(&Campaign{}, &Geo{}).Error
}
