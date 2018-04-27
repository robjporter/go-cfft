package application

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/timshannon/bolthold"
)

func (a *Application) submitMetricsToCapital() {
	var result []MetricData
	json := ""
	counter := a.Stats.GetCounter("tasks")
	err := a.db.data.Find(&result, bolthold.Where("Submitted").Eq(false).And("UUID").Ne(""))

	if err == nil {
		a.addToLogDebug(counter, map[string]interface{}{"Results": len(result)}, "There are some unsubmitted results to be sent to Capital.")
		for i := 0; i < len(result); i++ {
			json += a.createJSONFromResult(result[i])
			creationTime := strings.TrimSpace(time.Unix(result[i].CollectionTime, 0).String())
			a.dumpJSONToSendToFile(DATAOUTPUTFOLDER+creationTime+".json", json)
			a.submitInformationToCapital(result[i].UUID, json)
		}
	} else {
		a.addToLogDebug(counter, map[string]interface{}{"Error": err}, "There has been an error.")
	}
}

func (a *Application) createJSONFromResult(m MetricData) string {
	tmp, err := json.Marshal(m)
	if err != nil {
		a.addToLogDebug(a.Stats.GetCounter("tasks"), map[string]interface{}{"Error": err}, "There has been an error.")
		return ""
	}
	return string(tmp)
}

func (a *Application) dumpJSONToSendToFile(filename string, json string) bool {
	pos := strings.LastIndex(filename, "/")
	counter := a.Stats.GetCounter("tasks")

	dir := filename[:pos]
	filename = strings.TrimSpace(filename)
	os.MkdirAll(dir, os.ModePerm)

	err := ioutil.WriteFile(filename, []byte(json), os.ModePerm)
	if err != nil {
		a.addToLogDebug(counter, map[string]interface{}{"File": filename, "Error": err}, "Failed to save data to file.")
		return false
	}
	a.addToLogDebug(counter, nil, "Successfully saved data to file.")
	return true
}

func (a *Application) submitInformationToCapital(id string, json string) {
	transcode := ""
	counter := a.Stats.GetCounter("tasks")
	code := 201
	if override {
		code = 200
	}

	res, err := a.HX.SendDataToCaptial(json)
	if err != nil {
		a.addToLogDebug(counter, map[string]interface{}{"Error": err}, "We were unable to send data to Capital at this time.")
		a.LastError = err
	}

	if a.HX.GetResponseOK(res) {
		if a.HX.GetResponseCode(res) == code {
			a.addToLogDebug(counter, nil, "Successfully sent metrics to Capital.")
			transcode = a.HX.GetResponseItemString(res, "transactioncode")
			a.addToLogDebug(counter, map[string]interface{}{"Transaction Code": transcode}, "Recevied acknowledgement from Capital.")
		}
	}

	if transcode != "" {
		if a.updateMetricRecordAfterSubmission(transcode, id) == nil {
			a.addToLogDebug(counter, map[string]interface{}{"Metric ID": id, "Transaction Code": transcode}, "Successfully submitted metrics to Capital.")
		} else {
			a.addToLogWarning(counter, map[string]interface{}{"Metric ID": id, "Transaction Code": transcode}, "Failed to submit metrics to Capital.")
		}
	}
}

func (a *Application) updateMetricRecordAfterSubmission(code string, id string) error {
	counter := a.Stats.GetCounter("tasks")

	a.addToLogDebug(counter, map[string]interface{}{"Metric ID": id}, "Updating DB entries after successful submission.")
	err := a.db.data.UpdateMatching(&MetricData{}, bolthold.Where("UUID").Eq(id), func(record interface{}) error {
		update, ok := record.(*MetricData)
		if !ok {
			a.addToLogWarning(counter, nil, "We recevied a record type that we were not expecting.")
			return fmt.Errorf("We recevied a record type that we were not expecting")
		}
		update.Submitted = true
		update.SubmittedOn = time.Now()
		update.SubmittedTransactionCode = code
		a.addToLogDebug(counter, map[string]interface{}{"Metric ID": id}, "Updating DB entries after successful submission.")
		return nil
	})
	if err == nil {
		a.addToLogDebug(counter, map[string]interface{}{"Metric ID": id}, "Updating DB entries after successful submission.")
	} else {
		a.addToLogDebug(counter, map[string]interface{}{"Metric ID": id, "Error": err}, "Updating DB entries after successful submission.")
	}
	return err
}
