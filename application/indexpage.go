package application

import (
	"fmt"
	"os"

	"../packages/xTools/carbon"
	"../packages/xTools/conditions"

	"github.com/Sirupsen/logrus"
	"github.com/timshannon/bolthold"
)

func (a *Application) updateOnsiteIndexPage() {
	a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Beginning to update Index page data.")
	a.loadProcessedMetricsDataFromDB()
	calculateNewDates()
	a.updateProcessedMetricDataFromDB()
	a.saveProcessedMetricsDataToDB()
	a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Finished updating index page data.")
}

func (a *Application) loadProcessedMetricsDataFromDB() {
	a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Loading metrics from Database.")
	if a.db.data != nil {
		var b []Page
		a.db.data.Find(&b, bolthold.Where(bolthold.Key).Eq("processedIndexData"))
		if len(b) != 0 {
			a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Found metrics to load from Database.")
			index = b[0]
		} else {
			a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("No metrics currently stored in Database.")
		}
	} else {
		a.Logger.Fatal("Unfortunately we are not connected to a DB and at this point we should be.")
		os.Exit(1)
	}
	a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Loaded metrics from Database.")
}

func (a *Application) updateProcessedMetricDataFromDB() {
	a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Updating metrics from Database.")
	a.processCurrentData()
	a.processCurrentDayData()
	a.processCurrentWeekData()
	a.processCurrentMonthData()
	a.processCurrentYearData()
	a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Updated metrics from Database.")
}

func (a *Application) saveProcessedMetricsDataToDB() {
	a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Saving metrics from Database.")
	if a.db.data != nil {
		err := a.db.data.Insert("processedIndexData", index)
		if err == nil {
			a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Data has been saved successfully.")
		} else {
			a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks"), "Error": err}).Warning("There was an error writing to the Database.  No data has been saved.")
		}
	} else {
		a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Warning("There was an error connecting to the Database.  No data has been saved.")
	}
	a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Saved metrics from Database.")
}

func calculateNewDates() {
	index.StartOfDay = conditions.IfThen(index.StartOfDay == "", carbon.Now().StartOfDay().ToDateTimeString()).(string)
	index.EndOfDay = conditions.IfThen(index.EndOfDay == "", carbon.Now().EndOfDay().ToDateTimeString()).(string)
	index.StartOfWeek = conditions.IfThen(index.StartOfWeek == "", carbon.Now().StartOfWeek().ToDateTimeString()).(string)
	index.EndOfWeek = conditions.IfThen(index.EndOfWeek == "", carbon.Now().EndOfWeek().ToDateTimeString()).(string)
	index.CurrentMonthName = conditions.IfThen(index.CurrentMonthName == "", carbon.Now().MonthName()).(string)
	index.StartOfMonth = conditions.IfThen(index.StartOfMonth == "", carbon.Now().StartOfMonth().ToDateTimeString()).(string)
	index.EndOfMonth = conditions.IfThen(index.EndOfMonth == "", carbon.Now().EndOfMonth().ToDateTimeString()).(string)
	index.StartOfPreviousMonth = conditions.IfThen(index.StartOfPreviousMonth == "", carbon.Now().PreviousMonthStartDay().ToDateTimeString()).(string)
	index.EndOfPreviousMonth = conditions.IfThen(index.EndOfPreviousMonth == "", carbon.Now().PreviousMonthLastDay().ToDateTimeString()).(string)
	index.PreviousMonthName = conditions.IfThen(index.PreviousMonthName == "", carbon.Now().PreviousMonthStartDay().MonthName()).(string)
	index.StartOfYear = conditions.IfThen(index.StartOfYear == "", carbon.Now().StartOfYear().ToDateTimeString()).(string)
	index.EndOfYear = conditions.IfThen(index.EndOfYear == "", carbon.Now().EndOfYear().ToDateTimeString()).(string)
	index.CurrentDay = conditions.IfThen(index.CurrentDay == 0, carbon.Now().DayNumber()).(int)
	index.CurrentMonth = conditions.IfThen(index.CurrentMonth == 0, carbon.Now().MonthNumber()).(int)
	index.CurrentYear = conditions.IfThen(index.CurrentYear == 0, carbon.Now().YearNumber()).(int)

}

func (a *Application) processCurrentData() {
	a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Processing data for current day.")

	var data []MetricData
	a.db.data.Find(&data, bolthold.Where(bolthold.Key).Eq("processedIndexData"))
	a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks"), "Datapoints": len(data)}).Debug("Found data to process.")

	fmt.Println(data[0])

	a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Finished processing data for current day.")
}

func (a *Application) processCurrentDayData() {
	a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Processing data for current day.")

	var data []MetricData
	a.db.data.Find(&data, bolthold.Where(bolthold.Key).Eq("processedIndexData"))
	a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks"), "Datapoints": len(data)}).Debug("Found data to process.")

	submittedCount := 0
	for i := 0; i < len(data); i++ {
		if data[i].Submitted {
			submittedCount++
		}
	}

	a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks")}).Debug("Finished processing data for current day.")
}

func (a *Application) processCurrentWeekData()  {}
func (a *Application) processCurrentMonthData() {}
func (a *Application) processCurrentYearData()  {}
