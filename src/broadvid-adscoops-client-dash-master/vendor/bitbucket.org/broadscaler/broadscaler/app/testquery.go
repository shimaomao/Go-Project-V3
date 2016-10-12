package app

import (
	"time"

	"bitbucket.org/broadscaler/broadscaler/app/janitor"
)

func getDayStart() time.Time {
	loc, _ := time.LoadLocation("America/Los_Angeles")

	local := time.Now().In(loc)

	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
}
func TestQuery() {
	AdscoopsDB.LogMode(true)

	janitor.AdscoopsDB = AdscoopsDB

	janitor.TestRun()

}
