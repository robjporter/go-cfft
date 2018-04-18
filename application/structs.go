package application

import (
	"sync"
	"../packages/cron"
	"../packages/health"
	"../packages/hxconnect"
	"github.com/labstack/echo"
	"github.com/Sirupsen/logrus"
	"github.com/timshannon/bolthold"
)

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