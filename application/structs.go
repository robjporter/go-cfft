package application

import (
	"sync"

	"../packages/cron"
	"../packages/health"
	"../packages/xTools/hxconnect"
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/timshannon/bolthold"
)

type ProcessedMetricData struct {
}
type Page struct {
	CurrentDay        int
	CurrentMonth      int
	CurrentYear       int
	CurrentMonthName  string
	PreviousMonthName string

	StartOfDay   int64 // Get start of day time
	StartOfMonth int64 // Get stat of month
	StartOfWeek  int64 // Get start of week
	StartOfYear  int64

	EndOfDay   int64
	EndOfMonth int64
	EndOfWeek  int64
	EndOfYear  int64

	StartOfPreviousMonth int64 // Get start of previous month
	EndOfPreviousMonth   int64

	Current          ProcessedMetricData
	CurrentDayData   ProcessedMetricData
	CurrentWeekData  ProcessedMetricData
	CurrentMonthData ProcessedMetricData
	CurrentYearData  ProcessedMetricData
}

type TaskCounter struct {
	taskCounter     int64
	taskcounterLock sync.RWMutex
}

type Checks struct {
	URLS    health.CompositeChecker
	Handler health.Handler
}

type Application struct {
	db        *database
	Versions  version
	HX        *hxconnect.Connection
	Logger    *logrus.Logger
	Server    *echo.Echo
	LastError error
	Flags     *flags
	Crons     *cron.Cron
	Port      int
	StartTime int64
	Stats     Stat
	Checkers  Checks
}

type database struct {
	dbpath string
	data   *bolthold.Store
}

type flags struct {
	FirstRun bool
}

type version struct {
	appVersion  string
	goVersion   string
	osVersion   string
	archVersion string
	cpuVersion  string
}
