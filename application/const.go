package application

import (
	"time"
)

const (
	VERSION                  = "0.5.0"
	METRICKEY                = "NNDOBIXNSfWeKjE7TSyvPbgAFHXL6xSR"
	METRICSERVER             = "http://localhost:5003"
	DBPATH                   = "data.db"
	PORT                     = 1323
	TIMERINVENTORY           = 1
	TIMERINVENTORYSUBMISSION = 2
	TIMEROUT                 = 10
	TIMEREGENERATION         = 5
	TIMERREGENERATEINDEX     = time.Duration(TIMEREGENERATION) * time.Minute
	TIMERGATHERSTAT          = time.Duration(TIMERINVENTORY) * time.Minute
	TIMERSUBMITMETRICS       = time.Duration(TIMERINVENTORYSUBMISSION) * time.Minute //time.Duration(TIMERINVENTORYSUBMISSION) * time.Hour
	DATAOUTPUTFOLDER         = "output/"
)
