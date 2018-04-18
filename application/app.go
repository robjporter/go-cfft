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
	"../packages/health"
	"../packages/health/url"
	"../packages/hxconnect"

	static "github.com/Code-Hex/echo-static"
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/timshannon/bolthold"
)

const (
	VERSION                  = "0.5.0"
	METRICKEY                = "NNDOBIXNSfWeKjE7TSyvPbgAFHXL6xSR"
	METRICSERVER             = "http://localhost:5003"
	DBPATH                   = "data.db"
	PORT                     = 1323
	TIMERINVENTORY           = 1
	TIMERINVENTORYSUBMISSION = 1
	TIMEROUT                 = 10
	TIMEREGENERATION = 5
	TIMERREGENERATEINDEX 	 = time.Duration(TIMEREGENERATION) * time.Minute
	TIMERGATHERSTAT          = time.Duration(TIMERINVENTORY) * time.Minute
	TIMERSUBMITMETRICS       = time.Duration(2) * time.Minute //time.Duration(TIMERINVENTORYSUBMISSION) * time.Hour
	DATAOUTPUTFOLDER         = "output/"
)

var (
	once     sync.Once
	instance *Application
	override bool
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

func (a *Application) DEBUGOverrideLocalHXServer(server string) {
	a.HX.Credentials.Url = server
	override = true
}

func (a *Application) Start() {
	a.db.dbpath = DBPATH
	if !isFile(a.db.dbpath) {
		a.Logger.Debug("This looks like the first time the application has been run or is being reinitialised.")
		a.setupSetupRoutes()
	} else {
		a.Logger.Debug("This looks like the application has been setup.")
		if a.connectToDB(a.db.dbpath) {
			a.loadCredentialInformation()
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

func isFile(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}
	return true
}

func (a *Application) connectToDB(file string) bool {
	var err error
	a.db.data, err = bolthold.Open(file, 0666, nil)
	if err != nil {
		a.db.data = nil
		a.Logger.WithFields(logrus.Fields{"DB File": file, "Error": err}).Debug("The database was not located.")
		return false
	}
	a.Logger.WithFields(logrus.Fields{"DB File": file}).Debug("Connected to DB successfully.")
	return true

}

func (a *Application) GetServerPort() string {
	return strconv.Itoa(a.Port)
}

func (a *Application) SetServerPort(port int) {
	a.Port = port
}

func (a *Application) setupServer() {
	a.Server.Pre(middleware.RemoveTrailingSlash())
	a.Server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	a.Server.Use(middleware.Recover())
	a.Server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		CustomTimeFormat: "02-01-2006 | 03:04:05.00000",
		Format:           `[CFFT] ${time_custom} | ${method} | ${status} | ${uri} -> ${latency_human}` + "\n",
	}))
	a.Server.HideBanner = true
	a.Server.HidePort = true
	a.Logger.Debug("Web Server middleware setup complete.")
}

func (a *Application) setupErrorHandler() {
	a.Server.HTTPErrorHandler = a.customHTTPErrorHandler
	a.Logger.Debug("Web Server error handler setup complete.")
}

func (a *Application) setupRoutes() {
	a.Server.GET("/ping", a.routesHomePing)
	a.Server.GET("/health", func(e echo.Context) error {
		a.Checkers.Handler.ServeHTTP(e.Response().Writer,e.Request())
		return nil
	})
	a.Server.GET("/", a.routesHomeIndex)
	a.Logger.Debug("Server core routes initialisation complete.")
}

func (a *Application) setupSetupRoutes() {
	a.Server.GET("/setup", a.routesHomeSetup1)
	a.Server.GET("/hxsetup1", a.routesHomeHXSetup1)
	a.Server.POST("/hxsetup2", a.routesHomeHXSetup2)
	a.Server.POST("/hxsetup3", a.routesHomeHXSetup3)
	a.Server.POST("/hxsetup4", a.routesHomeHXSetup4)
	a.Server.POST("/hxsetup5", a.routesHomeHXSetup5)
	a.Logger.Debug("Server setup routes initialisation complete.")
}

func (a *Application) setupTemplates() {
	a.Server.Use(static.ServeRoot("/", NewAssets("public")))
	a.Server.Renderer = NewTemplate()
	a.Logger.Debug("Server Renderer initialised successfully.")
}

func (a *Application) saveAllData() {
	if a.connectToDB(a.db.dbpath) {
		err := a.db.data.Insert("credentials", a.HX.Credentials)
		if err == nil {
			a.Logger.Debug("Data has been saved successfully.")
		} else {
			a.Logger.Warning("There was an error writing to the Database.  No data has been saved.")
		}
	} else {
		a.Logger.Warning("There was an error connecting to the Database.  No data has been saved.")
	}
}

func (a *Application) loadCredentialInformation() {
	var b []hxconnect.Creds
	a.db.data.Find(&b, bolthold.Where(bolthold.Key).Eq("credentials"))
	if len(b) != 0 {
		if !override {
			a.HX.Credentials.Url = b[0].Url
		}
		a.HX.Credentials.Username = b[0].Username
		a.HX.Credentials.Password = b[0].Password
		a.HX.Credentials.Client_id = b[0].Client_id
		a.HX.Credentials.Client_secret = b[0].Client_secret
		a.Logger.Debug("Successfully loaded HX Connect credentials from database.")
	} else {
		a.Logger.Warning("There has been an error reading the credentials from the database.")
		os.Exit(1)
	}
}

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
	a.Logger.WithFields(logrus.Fields{"Task Jobs": 3}).Debug("Completed setting up Cron jobs.")
}

func (a *Application) applicationBanner() string {
	tmp := "  _____       ______    ______  _____" + "\n"
	tmp += " /  __ \\      |  ___|   |  ___||_   _|" + "\n"
	tmp += " | /  \\/      | |_      | |_     | |" + "\n"
	tmp += " | |          |  _|     |  _|    | |" + "\n"
	tmp += " | \\__/\\      | |       | |      | |" + "\n"
	tmp += "  \\____/apital\\_|lexible\\_|inance\\_/ool" + "\n"
	tmp += "===========================================" + "\n"
	tmp += "App Version: " + a.Versions.appVersion + " | Server Version: " + echo.Version + "\n"
	tmp += "===========================================" + "\n"
	tmp += time.Now().Format("Monday, 02-Jan-06 15:04:05")
	tmp += a.Server.Server.Addr + "\n"
	// TODO: FIX Server port not being displayed.
	return tmp
}
