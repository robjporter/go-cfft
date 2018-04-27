package application

import (
	"time"

	"../packages/cron"
)

func (a *Application) setupCronJobs() {
	counter := a.Stats.GetCounter("tasks")
	a.addToLogDebug(counter, nil, "Setting up Cron jobs.")
	a.Crons.Schedule(cron.Every(TIMERGATHERSTAT), cron.FuncJob(func() {
		counter := a.Stats.IncreaseCounter("tasks")
		a.addToLogDebug(counter, nil, "Gathering HX Metrics starting now.")
		a.gatherAndRecordMetrics()
		a.addToLogDebug(counter, nil, "Gathering HX Metrics finished now.")
	}))
	a.Crons.Schedule(cron.Every(TIMERSUBMITMETRICS), cron.FuncJob(func() {
		counter := a.Stats.IncreaseCounter("tasks")
		a.addToLogDebug(counter, nil, "Submitting HX Metrics starting now.")
		a.submitMetricsToCapital()
		a.addToLogDebug(counter, nil, "Subimitted HX Metrics finished now.")
	}))
	a.Crons.Schedule(cron.Every(TIMERREGENERATEINDEX), cron.FuncJob(func() {
		counter := a.Stats.IncreaseCounter("tasks")
		a.addToLogDebug(counter, nil, "Updating Index Page with latests HX Metrics starting now.")
		a.updateOnsiteIndexPage()
		a.addToLogDebug(counter, nil, "Updating Index Page with latests HX Metrics finished now.")

	}))
	a.Crons.Schedule(cron.Every(1*time.Hour), cron.FuncJob(func() {
		// Summerise day
		counter := a.Stats.IncreaseCounter("tasks")
		a.addToLogDebug(counter, nil, "Summerising latests HX Metrics for 24 hour period starting now.")
		a.addToLogDebug(counter, nil, "Summerising latests HX Metrics for 24 hour period finished now.")
	}))
	a.Crons.Schedule(cron.Every(24*time.Hour), cron.FuncJob(func() {
		// Summerise week
		counter := a.Stats.IncreaseCounter("tasks")
		a.addToLogDebug(counter, nil, "Summerising latests HX Metrics for current week period starting now.")
		a.addToLogDebug(counter, nil, "Summerising latests HX Metrics for current week period finished now.")
	}))
	a.Crons.Schedule(cron.Every(96*time.Hour), cron.FuncJob(func() {
		// Summerise Month
		counter := a.Stats.IncreaseCounter("tasks")
		a.addToLogDebug(counter, nil, "Summerising latests HX Metrics for current month period starting now.")
		a.addToLogDebug(counter, nil, "Summerising latests HX Metrics for current month period finished now.")
	}))
	a.addToLogDebug(counter, map[string]interface{}{"Task Jobs": len(a.Crons.Entries())}, "Completed setting up Cron jobs.")
}
