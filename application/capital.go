package application

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/timshannon/bolthold"
)

func (a *Application) submitMetricsToCapital() {
	var result []MetricData
	json := ""
	err := a.db.data.Find(&result, bolthold.Where("Submitted").Eq(false))
	//a.testPrintAll(result)
	if err == nil {
		a.Logger.WithFields(logrus.Fields{"Results": len(result)}).Debug("There are some unsubmitted results to be sent to Capital.")
		
		for i := 0; i < len(result); i++ {
			json += a.createJSONFromResult(result[i])
			creationTime := result[i].CollectionTime
			a.dumpJSONToSendToFile(DATAOUTPUTFOLDER+creationTime.String()+".json", json)
			a.submitInformationToCapital(json)
		}
	} else {
		fmt.Println(err)
	}
}

func (a *Application) createJSONFromResult(m MetricData) string {
	tmp, err := json.Marshal(m)
	if err != nil {
		a.Logger.WithFields(logrus.Fields{"Error": err}).Debug("There has been an error.")
		return ""
	}
	return string(tmp)
}

func (a *Application) dumpJSONToSendToFile(filename string, json string) bool {
	pos := strings.LastIndex(filename, "/")

	dir := filename[:pos]
	filename = strings.TrimSpace(filename)
	os.MkdirAll(dir, os.ModePerm)

	err := ioutil.WriteFile(filename, []byte(json), os.ModePerm)
	if err != nil {
		a.Logger.WithFields(logrus.Fields{"File": filename, "Error": err}).Debug("Failed to save data to file.")
		return false
	}
	a.Logger.WithFields(logrus.Fields{"File": filename}).Debug("Successfully saved data to file.")
	return true
}

func (a *Application) updateOnsiteIndexPage() {

}

func (a *Application) submitInformationToCapital(json string) {

}

