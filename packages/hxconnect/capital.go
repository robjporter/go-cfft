package hxconnect

type capurls struct {
	LicenseAgreement string
	ValidateContract string
}

type capital struct {
	Company 	string
	Start		string
	Duration	string
	EncryptionKey string
	Nodes		int
	Costs		interface{}
}

func getCapURLS() *capurls {
	return &capurls{
		LicenseAgreement: "/licenseagreement",
		ValidateContract: "/validatecontract",
	}
}

func (c *Connection) GetLicenseAgreement() (error) {
	e := c.sendGetRequest(c.Metrics.Server, c.capurls.LicenseAgreement)
	return e
}

func (c *Connection) ValidateContactNumberWithCisco(contract string) error {
	data := make(map[string]string)
	data["contractnumber"] = contract
	e := c.sendPostRequest(c.Metrics.Server, c.capurls.ValidateContract, data)
	return e
}