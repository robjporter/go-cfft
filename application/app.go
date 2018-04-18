package application

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"time"

	"../packages/cron"
	"../packages/hxconnect"
	"../packages/health"
	"../packages/health/url"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"

)

var (
	once     sync.Once
	instance *Application
	override bool
)

func init() {
	override = false
}

func GetInstance() *Application {
	once.Do(func() {
		instance = &Application{
			HX:        hxconnect.New(),
			Logger:    logrus.New(),
			Server:    echo.New(),
			Crons:     cron.New(),
			db:        &database{},
			Port:      PORT,
			StartTime: time.Now().Unix(),
			Stats:     Stat{counters: make(map[string]int64)},
			Checkers:  Checks{URLS: health.NewCompositeChecker(), Handler: health.NewHandler()},
		}
	})
	return instance
}

func New() *Application {
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := GetInstance()
	app.Logger.SetLevel(logrus.DebugLevel)

	app.Versions.appVersion = VERSION
	app.Versions.goVersion = runtime.Version()
	app.Versions.osVersion = runtime.GOOS
	app.Versions.archVersion = runtime.GOARCH
	app.Versions.cpuVersion = strconv.Itoa(runtime.NumCPU())

	app.HX.Metrics.Server = METRICSERVER
	app.HX.Metrics.Key = METRICKEY
	app.Logger.Debug("Initialisation complete.")

	app.setupServer()
	app.setupErrorHandler()
	app.setupRoutes()
	app.setupTemplates()
	return app
}

func (a *Application) setupCheckers() {
	a.Checkers.URLS.AddChecker("Cisco", url.NewChecker("https://www.cisco.com/"))
	a.Checkers.URLS.AddChecker("Capital", url.NewChecker(a.HX.Credentials.Url +"/health"))
	a.Checkers.Handler.AddChecker("Tests", a.Checkers.URLS)
}

func (a *Application) Start() {
	a.db.dbpath = DBPATH
	if !isFile(a.db.dbpath) {
		a.Logger.Debug("This looks like the first time the application has been run or is being reinitialised.")
		a.setupSetupRoutes()
	} else {
		a.Logger.Debug("This looks like the application has been setup.")
		if a.connectToDB(a.db.dbpath) {
			a.loadCredentialInformationFromDB()
			a.setupCronJobs()
		} else {
			a.Logger.Debug("Failed to connect to DB.")
		}
	}

	a.setupCheckers()
	a.Crons.Start()
	a.Server.Debug = true
	a.Logger.Debug("Application running and ready at http://<IP>:." + a.GetServerPort())
	// Start server
	go func() {
		if err := a.Server.Start(":" + a.GetServerPort()); err != nil {
			a.Logger.Info("shutting down the server")
		}
	}()

	fmt.Println(a.applicationBanner())

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.Server.Shutdown(ctx); err != nil {
		a.Logger.Fatal(err)
	}
	a.Stop()
}

func (a *Application) Stop() {
	a.Logger.Debug("Stopping all services.")
	a.db.data.Close()
	a.Crons.Stop()
	a.Logger.Debug("Successfully finished stopping all services.")
}



func (a *Application) GetServerPort() string {
	return strconv.Itoa(a.Port)
}

func (a *Application) SetServerPort(port int) {
	a.Port = port
}







