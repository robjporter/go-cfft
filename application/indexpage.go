package application

import (
	"fmt"

	"../packages/carbon"

	"github.com/Sirupsen/logrus"
)

func (a *Application) updateOnsiteIndexPage() {
	a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Beginning to update Index page data.")
	// Get start of day time
	// Get start of week
	// Get stat of month
	// Get start of previous month
	// Query DB for data from start of day
	// Query DB for data from current week
	// Query DB for data from current month
	// Query DB for data from previous month

	index.startOfDay = carbon.Now().StartOfDay().DateTimeString()
	index.endOfDay = carbon.Now().EndOfDay().DateTimeString()
	index.startOfWeek = carbon.Now().StartOfWeek().DateTimeString()
	index.endOfWeek = carbon.Now().EndOfWeek().DateTimeString()
	index.startOfMonth = carbon.Now().StartOfMonth().DateTimeString()
	index.endOfMonth = carbon.Now().EndOfMonth().DateTimeString()
	index.startOfYear = carbon.Now().StartOfYear().DateTimeString()
	index.endOfYear = carbon.Now().EndOfYear().DateTimeString()
	index.lastOfPreviousMonth = carbon.Now().PreviousMonthLastDay().EndOfDay().DateTimeString()
	index.startOfPreviousMonth = carbon.Now().PreviousMonthLastDay().StartOfMonth().DateTimeString()

	fmt.Println("START OF DAY: " + index.startOfDay)
	fmt.Println("START OF DAY: " + index.endOfDay)
	fmt.Println("START OF WEEK: " + index.startOfWeek)
	fmt.Println("START OF MONTH: " + index.startOfMonth)
	fmt.Println("LAST DAY OF MONTH: " + index.lastOfPreviousMonth)
	fmt.Println("START OF PREVIOUS MONTH: " + index.startOfPreviousMonth)

	a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Finished updating index page data.")
}
