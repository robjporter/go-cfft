package application

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

func (a *Application) routesHomeSetup1(c echo.Context) error {
	lerror := ""
	if a.LastError != nil {
		lerror = a.LastError.Error()
		a.LastError = nil
	}
	return c.Render(200, "setup.html", map[string]interface{}{"appname": "APPNAME", "title": "TITLE", "error": lerror})
}

func (a *Application) routesHomeHXSetup1(c echo.Context) error {
	agreement := ""
	agreementdate := ""
	res, err := a.HX.GetLicenseAgreement()
	if err != nil {
		a.Logger.WithFields(logrus.Fields{"Error": err}).Warning("Unable to retrieve license agreement.")
	}
	if a.HX.GetResponseOK(res) {
		if a.HX.GetResponseCode(res) == 200 {
			agreement = a.HX.GetResponseItemString(res, "agreementMessage")
			agreementdate = a.HX.GetResponseItemString(res, "agreementDate")
		}
	}
	return c.Render(200, "hxsetup1.html", map[string]interface{}{"appname": "APPNAME", "title": "TITLE", "agreement": agreement, "agreementdate": agreementdate})
}

func (a *Application) routesHomeHXSetup2(c echo.Context) error {
	if c.FormValue("secure") != "on" {
		c.Redirect(301, "/hxsetup1")
	}
	return c.Render(200, "hxsetup2.html", map[string]interface{}{"appname": "APPNAME", "title": "TITLE"})
}

func (a *Application) routesHomeHXSetup3(c echo.Context) error {
	url := c.FormValue("url")
	username := c.FormValue("username")
	password := c.FormValue("password")
	if url != "" && username != "" && password != "" {
		a.HX.SetUrl(url)
		a.HX.SetUsername(username)
		a.HX.SetPassword(password)
		return c.Render(200, "hxsetup3.html", map[string]interface{}{"appname": "APPNAME", "title": "TITLE"})
	}
	return c.Redirect(301, "/setup")
}

func (a *Application) routesHomeHXSetup4(c echo.Context) error {
	contract := c.FormValue("contract")
	if contract != "" {
		res, err := a.HX.ValidateContactNumberWithCisco(Encrypt([]byte(a.HX.Metrics.Key), []byte(contract)))
		if err != nil {
			a.LastError = err
			c.Redirect(301, "/setup")
		}
		if a.HX.GetResponseItemBool(res, "valid") {
			a.HX.Capital.Company = a.HX.GetResponseItemString(res, "company")
			a.HX.Capital.Start = a.HX.GetResponseItemString(res, "contractStartDate")
			a.HX.Capital.Duration = a.HX.GetResponseItemString(res, "contractDuration")
			a.HX.Capital.EncryptionKey = a.HX.GetResponseItemString(res, "encryptionKey")
			a.HX.Capital.Nodes = a.HX.GetResponseItemInt(res, "nodes")
			a.HX.Capital.Costs = a.HX.GetResponseItem(res, "costs")
			return c.Render(200, "hxsetup4.html", map[string]interface{}{"appname": "APPNAME", "title": "TITLE", "company": a.HX.Capital.Company, "start": a.HX.Capital.Start, "duration": a.HX.Capital.Duration, "nodes": a.HX.Capital.Nodes, "url": a.HX.GetUrl(), "username": a.HX.GetUsername()})
		}
		a.LastError = errors.New("The server was unable to idenitfy your contract number.")
		c.Redirect(301, "/setup")
	}
	a.LastError = errors.New("There was an issue with the supplied contract number.")
	return c.Redirect(301, "/setup")
}

func (a *Application) routesHomeHXSetup5(c echo.Context) error {
	res, err := a.HX.Authenticate()
	success := false
	a.LastError = err
	if err == nil {
		if a.HX.GetResponseOK(res) {
			if a.HX.GetResponseCode(res) == 201 {
				success = true
				a.LastError = nil
				a.saveCredentialsDataToDB()
			} else {
				a.LastError = errors.New("We recevied a response other than 201 as expected.")
				a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks"), "Error": err, "Code": a.HX.GetResponseCode(res)}).Debug("We received an error and cannot continue.")
			}
		} else {
			a.LastError = errors.New("We recevied a response which was invalid.")
			a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks"), "Error": err, "OK": a.HX.GetResponseOK(res)}).Debug("We received an error and cannot continue.")
		}
	} else {
		a.Logger.WithFields(logrus.Fields{"Task Number": a.Stats.GetCounter("tasks"), "Error": err}).Debug("We received an error and cannot continue.")
	}
	return c.Render(200, "hxsetup5.html", map[string]interface{}{"appname": "APPNAME", "title": "TITLE", "success": success})
}
