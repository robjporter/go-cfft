package application

import (
	"sync"

	"../packages/cron"
	"../packages/health"
	"../packages/hxconnect"
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/timshannon/bolthold"
)

type Page struct {
	startOfDay           string // Get start of day time
	endOfDay             string
	startOfWeek          string // Get start of week
	endOfWeek            string
	startOfMonth         string // Get stat of month
	endOfMonth           string
	startOfYear          string
	endOfYear            string
	startOfPreviousMonth string // Get start of previous month
	endOfPreviousMonth   string
}

/*
type Page struct {
	startOfDay           time.Time // Get start of day time
	startOfWeek          time.Time // Get start of week
	startOfMonth         time.Time // Get stat of month
	startOfPreviousMonth time.Time // Get start of previous month
	lastOfPreviousMonth  time.Time
}*/

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
