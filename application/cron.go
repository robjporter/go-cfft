package application

import (
	"time"

	"../packages/cron"
	"github.com/Sirupsen/logrus"
)

func (a *Application) setupCronJobs() {
	a.Logger.Debug("Setting up Cron jobs.")
	a.Crons.Schedule(cron.Every(TIMERGATHERSTAT), cron.FuncJob(func() {
		counter := a.Stats.IncreaseCounter("tasks")
		a.Logger.WithFields(logrus.Fields{"Task Number": counter, "Task Title": "Gather HX Metrics"}).Debug("Task starting now.")
		a.gatherAndRecordMetrics()
		a.Logger.WithFields(logrus.Fields{"Task Number": counter, "Task Title": "Gather HX Metrics"}).Debug("Task finished now.")
	}))
	a.Crons.Schedule(cron.Every(TIMERSUBMITMETRICS), cron.FuncJob(func() {
		counter := a.Stats.IncreaseCounter("tasks")
		a.Logger.WithFields(logrus.Fields{"Task Number": counter, "Task Title": "Submit HX Metrics"}).Debug("Task starting now.")
		a.submitMetricsToCapital()
		a.Logger.WithFields(logrus.Fields{"Task Number": counter, "Task Title": "Submit HX Metrics"}).Debug("Task finished now.")
	}))
	a.Crons.Schedule(cron.Every(TIMERREGENERATEINDEX), cron.FuncJob(func() {
		counter := a.Stats.IncreaseCounter("tasks")
		a.Logger.WithFields(logrus.Fields{"Task Number": counter, "Task Title": "Updating Index Page with latest metrics."}).Debug("Task starting now.")
		a.updateOnsiteIndexPage()
		a.Logger.WithFields(logrus.Fields{"Task Number": counter, "Task Title": "Updating Index Page with latest metrics."}).Debug("Task finished now.")

	}))
	a.Crons.Schedule(cron.Every(1*time.Hour), cron.FuncJob(func() {
		// Summerise day
		counter := a.Stats.IncreaseCounter("tasks")
		a.Logger.WithFields(logrus.Fields{"Task Number": counter, "Task Title": "Summerising metrics for current 24 hour period."}).Debug("Task starting now.")
		a.Logger.WithFields(logrus.Fields{"Task Number": counter, "Task Title": "Summerising metrics for current 24 hour period."}).Debug("Task finished now.")
	}))
	a.Crons.Schedule(cron.Every(24*time.Hour), cron.FuncJob(func() {
		// Summerise week
		counter := a.Stats.IncreaseCounter("tasks")
		a.Logger.WithFields(logrus.Fields{"Task Number": counter, "Task Title": "Summerising metrics for current week period."}).Debug("Task starting now.")
		a.Logger.WithFields(logrus.Fields{"Task Number": counter, "Task Title": "Summerising metrics for current week period."}).Debug("Task finished now.")
	}))
	a.Crons.Schedule(cron.Every(96*time.Hour), cron.FuncJob(func() {
		// Summerise Month
		counter := a.Stats.IncreaseCounter("tasks")
		a.Logger.WithFields(logrus.Fields{"Task Number": counter, "Task Title": "Summerising metrics for current month period."}).Debug("Task starting now.")
		a.Logger.WithFields(logrus.Fields{"Task Number": counter, "Task Title": "Summerising metrics for current month period."}).Debug("Task finished now.")
	}))
	a.Logger.WithFields(logrus.Fields{"Task Jobs": len(a.Crons.Entries())}).Debug("Completed setting up Cron jobs.")
}
