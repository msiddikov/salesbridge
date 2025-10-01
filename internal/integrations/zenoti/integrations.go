package integrations_zenoti

import (
	"client-runaway-zenoti/internal/runway"
	svc_jpmreport "client-runaway-zenoti/internal/services/svc_jpmReport"
	"client-runaway-zenoti/internal/tgbot"
	"time"

	"github.com/go-co-op/gocron"
)

const (
	appointmentsFetchFromDays = -5
	appointmentsFetchToDays   = 30
)

func StartScheduledJobs() {
	tgbot.Notify("Scheduled jobs", "Starting scheduled jobs", false)
	s := gocron.NewScheduler(time.UTC)
	s.Every(2).Hours().Do(runFrequentJobs)
	s.StartBlocking()
}

func runFrequentJobs() {
	tgbot.Notify("Scheduled jobs", "Updating all tokens", false)
	runway.UpdateAllTokens()

	tgbot.Notify("Scheduled jobs", "Updating stages", false)
	UpdateStages()

	tgbot.Notify("Scheduled jobs", "Force checking all", false)
	runway.ForceCheckAllV2()

	tgbot.Notify("Scheduled jobs", "Syncing calendars", false)
	SyncCalendars()

	// run daily jobs if time is <= 2am
	if true || time.Now().Hour() <= 2 {
		runDailyJobs()
	}
	tgbot.Notify("Scheduled jobs", "Frequent jobs ended", false)
}

func runDailyJobs() {
	tgbot.Notify("Scheduled jobs", "Daily jobs started", false)
	UpdateNotesV2(false)
	runway.SetRanksBySales()
	runway.CheckForNewLeads()
	tgbot.Notify("Scheduled jobs", "Daily jobs ended", false)

	svc_jpmreport.UpdateAllLocationsReportDataForLargePeriod(
		time.Now().Add(-36*time.Hour),
		time.Now(),
	)
}
