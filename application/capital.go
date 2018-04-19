package application

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/timshannon/bolthold"
)

func (a *Application) submitMetricsToCapital() {
	var result []MetricData
	json := ""

	err := a.db.data.Find(&result, bolthold.Where("Submitted").Eq(false).And("UUID").Ne(""))

	if err == nil {
		a.Logger.WithFields(logrus.Fields{"Results": len(result)}).Debug("There are some unsubmitted results to be sent to Capital.")

		for i := 0; i < len(result); i++ {
			json += a.createJSONFromResult(result[i])
			creationTime := result[i].CollectionTime
			a.dumpJSONToSendToFile(DATAOUTPUTFOLDER+creationTime.String()+".json", json)
			a.submitInformationToCapital(result[i].UUID, json)
		}
	} else {
		a.Logger.WithFields(logrus.Fields{"Error": err}).Debug("There has been an error.")
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

func (a *Application) submitInformationToCapital(id string, json string) {
	transcode := ""
	code := 201
	if override {
		code = 200
	}

	res,err := a.HX.SendDataToCaptial(json)
	if err != nil {
		a.Logger.Debug("We were unable to send data to Capital at this time.")
		a.LastError = err
	}

	if a.HX.GetResponseOK(res) {
		if a.HX.GetResponseCode(res) == code {
			a.Logger.Debug("Successfully sent metrics to Capital.")
			transcode = a.HX.GetResponseItemString(res,"transactioncode")
			a.Logger.WithFields(logrus.Fields{"Transaction Code": transcode}).Debug("Recevied acknowledgement from Capital.")
		}
	}

	if transcode != "" {
		if a.updateMetricRecordAfterSubmission(transcode, id) == nil {
			a.Logger.WithFields(logrus.Fields{"Metric ID": id, "Transaction Code": transcode}).Debug("Successfully submitted metrics to Capital.")
		} else {
			a.Logger.WithFields(logrus.Fields{"Metric ID": id, "Transaction Code": transcode}).Warn("Failed to submit metrics to Capital.")
		}
	}
}

func (a *Application) updateMetricRecordAfterSubmission(code string, id string) error {
	err := a.db.data.UpdateMatching(&MetricData{}, bolthold.Where("UUID").Eq(id), func(record interface{}) error {
		update, ok := record.(*MetricData)
		if !ok {
			a.Logger.Warn("We recevied a record type that we were not expecting.")
			return fmt.Errorf("We recevied a record type that we were not expecting.")
		}
		update.Submitted = true
		update.SubmittedOn = time.Now()
		update.SubmittedTransactionCode = code
		return nil
	})
	return err
}
