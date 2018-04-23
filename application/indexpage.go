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

	index.startOfDay = carbon.Now().StartOfDay().ToDateTimeString()
	index.endOfDay = carbon.Now().EndOfDay().ToDateTimeString()
	index.startOfWeek = carbon.Now().StartOfWeek().ToDateTimeString()
	index.endOfWeek = carbon.Now().EndOfWeek().ToDateTimeString()
	index.startOfMonth = carbon.Now().StartOfMonth().ToDateTimeString()
	index.endOfMonth = carbon.Now().EndOfMonth().ToDateTimeString()
	index.startOfPreviousMonth = carbon.Now().PreviousMonthStartDay().ToDateTimeString()
	index.endOfPreviousMonth = carbon.Now().PreviousMonthLastDay().ToDateTimeString()
	index.startOfYear = carbon.Now().StartOfYear().ToDateTimeString()
	index.endOfYear = carbon.Now().EndOfYear().ToDateTimeString()

	fmt.Println("START OF DAY: " + index.startOfDay)
	fmt.Println("END OF DAY: " + index.endOfDay)
	fmt.Println("START OF WEEK: " + index.startOfWeek)
	fmt.Println("END OF WEEK: " + index.endOfWeek)
	fmt.Println("START OF MONTH: " + index.startOfMonth)
	fmt.Println("END OF MONTH: " + index.endOfMonth)
	fmt.Println("START OF PREVIOUS MONTH: " + index.startOfPreviousMonth)
	fmt.Println("END OF PREVIOUS MONTH: " + index.endOfPreviousMonth)
	fmt.Println("START OF YEAR: " + index.startOfYear)
	fmt.Println("END OF YEAR: " + index.endOfYear)

	a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Finished updating index page data.")
}
