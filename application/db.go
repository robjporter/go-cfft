package application

import (
	"os"
	"../packages/hxconnect"
	"github.com/Sirupsen/logrus"
	"github.com/timshannon/bolthold"
)

func (a *Application) loadCredentialInformationFromDB() {
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

func (a *Application) saveCredentialsDataToDB() {
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