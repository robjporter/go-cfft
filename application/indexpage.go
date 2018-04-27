package application

import (
	"fmt"
	"os"

	"../packages/xTools/carbon"
	"../packages/xTools/conditions"

	"github.com/timshannon/bolthold"
)

func (a *Application) updateOnsiteIndexPage() {
	count := a.Stats.GetCounter("tasks")
	a.addToLogDebug(count, nil, "Beginning to update Index page data.")
	a.loadProcessedMetricsDataFromDB()
	calculateNewDates()
	a.updateProcessedMetricDataFromDB()
	a.saveProcessedMetricsDataToDB()
	a.addToLogDebug(count, nil, "Finished updating index page data.")
}

func (a *Application) loadProcessedMetricsDataFromDB() {
	count := a.Stats.GetCounter("tasks")
	a.addToLogDebug(count, nil, "Loading metrics from Database.")
	if a.db.data != nil {
		var b []Page
		a.db.data.Find(&b, bolthold.Where(bolthold.Key).Eq("processedIndexData"))
		if len(b) != 0 {
			a.addToLogDebug(count, nil, "Found metrics to load from Database.")
			index = b[0]
		} else {
			a.addToLogDebug(count, nil, "No metrics currently stored in Database.")
		}
	} else {
		a.addToLogDebug(count, nil, "Unfortunately we are not connected to a DB and at this point we should be.")
		os.Exit(1)
	}
	a.addToLogDebug(count, nil, "Loaded metrics from Database.")
}

func (a *Application) updateProcessedMetricDataFromDB() {
	count := a.Stats.GetCounter("tasks")
	a.addToLogDebug(count, nil, "Updating metrics from Database.")
	a.processCurrentData()
	a.processCurrentDayData()
	a.processCurrentWeekData()
	a.processCurrentMonthData()
	a.processCurrentYearData()
	a.addToLogDebug(count, nil, "Updated metrics from Database.")
}

func (a *Application) saveProcessedMetricsDataToDB() {
	count := a.Stats.GetCounter("tasks")
	a.addToLogDebug(count, nil, "Saving metrics from Database.")
	if a.db.data != nil {
		err := a.db.data.Insert("processedIndexData", index)
		if err == nil {
			a.addToLogDebug(count, nil, "Data has been saved successfully.")
		} else {
			a.addToLogWarning(count, nil, "There was an error writing to the Database.  No data has been saved.Saving metrics from Database.")
		}
	} else {
		a.addToLogWarning(count, nil, "There was an error connecting to the Database.  No data has been saved.")
	}
	a.addToLogDebug(count, nil, "Saved metrics from Database.")
}

func calculateNewDates() {
	index.StartOfDay = conditions.IfThenElse(index.StartOfDay == 0, carbon.Now().StartOfDay().ToTimeStamp(), index.StartOfDay).(int64)
	index.EndOfDay = conditions.IfThenElse(index.EndOfDay == 0, carbon.Now().EndOfDay().ToTimeStamp(), index.EndOfDay).(int64)
	index.StartOfWeek = conditions.IfThenElse(index.StartOfWeek == 0, carbon.Now().StartOfWeek().ToTimeStamp(), index.StartOfWeek).(int64)
	index.EndOfWeek = conditions.IfThenElse(index.EndOfWeek == 0, carbon.Now().EndOfWeek().ToTimeStamp(), index.EndOfWeek).(int64)
	index.CurrentMonthName = conditions.IfThenElse(index.CurrentMonthName == "", carbon.Now().MonthName(), index.CurrentMonthName).(string)
	index.StartOfMonth = conditions.IfThenElse(index.StartOfMonth == 0, carbon.Now().StartOfMonth().ToTimeStamp(), index.StartOfMonth).(int64)
	index.EndOfMonth = conditions.IfThenElse(index.EndOfMonth == 0, carbon.Now().EndOfMonth().ToTimeStamp(), index.EndOfMonth).(int64)
	index.StartOfPreviousMonth = conditions.IfThenElse(index.StartOfPreviousMonth == 0, carbon.Now().PreviousMonthStartDay().ToTimeStamp(), index.StartOfPreviousMonth).(int64)
	index.EndOfPreviousMonth = conditions.IfThenElse(index.EndOfPreviousMonth == 0, carbon.Now().PreviousMonthLastDay().ToTimeStamp(), index.EndOfPreviousMonth).(int64)
	index.PreviousMonthName = conditions.IfThenElse(index.PreviousMonthName == "", carbon.Now().PreviousMonthStartDay().MonthName(), index.PreviousMonthName).(string)
	index.StartOfYear = conditions.IfThenElse(index.StartOfYear == 0, carbon.Now().StartOfYear().ToTimeStamp(), index.StartOfYear).(int64)
	index.EndOfYear = conditions.IfThenElse(index.EndOfYear == 0, carbon.Now().EndOfYear().ToTimeStamp(), index.EndOfYear).(int64)
	index.CurrentDay = conditions.IfThenElse(index.CurrentDay == 0, carbon.Now().DayNumber(), index.CurrentDay).(int)
	index.CurrentMonth = conditions.IfThenElse(index.CurrentMonth == 0, carbon.Now().MonthNumber(), index.CurrentMonth).(int)
	index.CurrentYear = conditions.IfThenElse(index.CurrentYear == 0, carbon.Now().YearNumber(), index.CurrentYear).(int)
}

func (a *Application) processCurrentData() {
	count := a.Stats.GetCounter("tasks")
	a.addToLogDebug(count, nil, "Processing data for current day.")

	var data []MetricData
	a.db.data.Find(&data, bolthold.Where("CollectionTime").Gt(index.StartOfDay).SortBy("CollectionTime").Reverse())

	if len(data) > 0 {
		a.addToLogDebug(count, map[string]interface{}{"Datapoints": len(data)}, "Found data to process.")
		fmt.Println(data[0])
		// TODO
	} else {
		a.addToLogDebug(count, map[string]interface{}{"Datapoints": 0}, "No data was found to process.")
	}

	a.addToLogDebug(count, nil, "Finished processing data for current day.")
}

func (a *Application) processCurrentDayData() {
	count := a.Stats.GetCounter("tasks")
	a.addToLogDebug(count, nil, "Processing data for current day.")

	var data []MetricData
	a.db.data.Find(&data, bolthold.Where("CollectionTime").Gt(index.StartOfDay))

	if len(data) > 0 {
		a.addToLogDebug(count, map[string]interface{}{"Datapoints": len(data)}, "Found data to process.")

		submittedCount := 0
		for i := 0; i < len(data); i++ {
			if data[i].Submitted {
				submittedCount++
			}
		}
	} else {
		a.addToLogDebug(count, map[string]interface{}{"Datapoints": 0}, "No data found to process.")
	}

	a.addToLogDebug(count, nil, "Finished processing data for current day.")
}

func (a *Application) processCurrentWeekData()  {}
func (a *Application) processCurrentMonthData() {}
func (a *Application) processCurrentYearData()  {}
