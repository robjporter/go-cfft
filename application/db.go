package application

import (
	"os"

	"../packages/xTools/hxconnect"
	"github.com/timshannon/bolthold"
)

func (a *Application) loadCredentialInformationFromDB() {
	var b []hxconnect.Creds
	count := a.Stats.GetCounter("tasks")
	a.addToLogDebug(count, nil, "Beginning to read data from DB.")

	a.db.data.Find(&b, bolthold.Where(bolthold.Key).Eq("credentials"))
	if len(b) != 0 {
		if !override {
			a.HX.Credentials.Url = b[0].Url
		}
		a.HX.Credentials.Username = b[0].Username
		a.HX.Credentials.Password = b[0].Password
		a.HX.Credentials.Client_id = b[0].Client_id
		a.HX.Credentials.Client_secret = b[0].Client_secret
		a.addToLogDebug(count, nil, "Successfully loaded HX Connect credentials from database.")
	} else {
		a.addToLogWarning(count, nil, "There has been an error reading the credentials from the database.")
		os.Exit(1)
	}
	a.addToLogDebug(count, nil, "Finished reading data from DB.")
}

func (a *Application) saveCredentialsDataToDB() {
	count := a.Stats.GetCounter("tasks")
	a.addToLogDebug(count, nil, "Beginning to save credentials to DB.")

	if a.connectToDB(a.db.dbpath) {
		err := a.db.data.Insert("credentials", a.HX.Credentials)
		if err == nil {
			a.addToLogDebug(count, nil, "Data has been saved successfully.")
		} else {
			a.addToLogWarning(count, map[string]interface{}{"Error": err}, "There was an error writing to the Database.  No data has been saved.")
		}
	} else {
		a.addToLogDebug(count, nil, "There was an error connecting to the Database.  No data has been saved.")
	}
	a.addToLogDebug(count, nil, "Finished saving credentials to DB.")
}

func (a *Application) connectToDB(file string) bool {
	var err error
	count := a.Stats.GetCounter("tasks")
	a.addToLogDebug(count, nil, "Beginning connection to DB.")

	a.db.data, err = bolthold.Open(file, 0666, nil)
	if err != nil {
		a.db.data = nil
		a.addToLogDebug(count, map[string]interface{}{"DB File": file, "Error": err}, "The database was not located.")
		return false
	}
	a.addToLogDebug(count, map[string]interface{}{"DB File": file}, "Connected to DB successfully.")
	return true
}
