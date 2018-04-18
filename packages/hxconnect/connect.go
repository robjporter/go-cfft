package hxconnect

// GET /cluster/savings
// POST /clusterVersionDetails
// GET /virtplatform/cluster
// GET /virtplatform/vms
// GET /datastores
// GET /summary
// GET /appliances

import (
  "time"
  "errors"
  "strings"

  "../request"
)

type Connection struct {
  token               string
  timeout             int
  hxurls              *hxurls
  capurls             *capurls
  Request             *request.Client
  Metrics             *metrics
  Capital             *capital
  Credentials         *Creds
}

type Creds struct {
  Url                 string
  Username            string
  Password            string
  Client_id           string
  Client_secret       string
}

type hxurls struct {
  Authentication      string
  About               string
  ClusterInfo         string
  ClusterSavings      string
  ClusterPlatform     string
  ClusterVM           string 
  ClusterDatastores   string 
  ClusterSummary      string
  ClusterAppliances   string 
  ClusterVersion      string
}

type metrics struct {
	Server					    string
	Key							    string
}

func init() {}
func New() *Connection {
  return &Connection{
      timeout: 30,
      hxurls: getHXURLS(),
      capurls: getCapURLS(),
      Request: request.New(),
      Metrics: &metrics{},
      Capital: &capital{},
      Credentials: &Creds{
        Url: "https://",
        Username: "",
        Password: "",
        Client_id: "HxGuiClient",
        Client_secret: "Sunnyvale",
      },
  }
}

func getHXURLS() *hxurls {
  return &hxurls{
    Authentication: "/aaa/v1/auth?grant_type=password", // POST
    About: "/rest/about", // GET
    ClusterInfo: "/rest/clusters", // GET
    ClusterSavings: "/rest/cluster/savings", // GET 
    ClusterPlatform: "/rest/virtplatform/cluster", // GET 
    ClusterVM: "/rest/virtplatform/vms", // GET 
    ClusterDatastores: "/rest/datastores", // GET 
    ClusterSummary: "/rest/summary", // GET 
    ClusterAppliances: "/rest/appliances", // GET 
    ClusterVersion: "/rest/clusterVersionDetails", // POST
  }
}

func (c *Connection) SetToken(token string) {
  c.token = token
}

func (c *Connection) GetToken() string {
  return "Bearer " + strings.TrimSpace(c.token)
}

func (c *Connection) SetUsername(username string) {
  c.Credentials.Username = username
}

func (c *Connection) SetPassword(password string) {
  c.Credentials.Password = password
}

func (c *Connection) SetClientID(clientid string) {
  c.Credentials.Client_id = clientid
}

func (c *Connection) SetClientSecret(clientsecret string) {
  c.Credentials.Client_secret = clientsecret
}

func (c *Connection) SetUrl(url string) {
  if strings.HasPrefix(url,"http://") || strings.HasPrefix(url,"https://") {
    c.Credentials.Url = url
  } else {
    c.Credentials.Url = c.Credentials.Url + url
  }
}

func (c *Connection) SetTimeout(timeout int) {
  c.timeout = timeout
}

func (c *Connection) GetUsername() string {
  return c.Credentials.Username
}

func (c *Connection) GetPassword() string {
  return c.Credentials.Password
}

func (c *Connection) GetClientID() string {
  return c.Credentials.Client_id
}

func (c *Connection) GetClientSecret() string {
  return c.Credentials.Client_secret
}

func (c *Connection) GetUrl() string {
  return c.Credentials.Url
}

func (c *Connection) GetTimeout() time.Duration {
  return time.Duration(c.timeout)*time.Second
}

func (c *Connection) settingsMade() bool {
  if c.Credentials.Url != "https://" && c.Credentials.Username != "" && c.Credentials.Password != "" {
    return true
  }
  return false
}

func (c *Connection) sendPostRequest(address string, url string, payload map[string]string) (error){
  c.Request = request.New()
  _, err := c.Request.
    Post(address + url).
    Timeout(c.GetTimeout()).
    Send(payload).
    Set("content-type", "application/json").
    Set("cache-control", "no-cache").
    Accept("application/json").
    JSON()
  return err
}

func (c *Connection) sendSecurePostRequest(address string, url string, payload map[string]string) (error){
  c.Request = request.New()
  _, err := c.Request.
    Post(address + url).
    Timeout(c.GetTimeout()).
    Send(payload).
		Set("Authorization", c.GetToken()).
    Set("content-type", "application/json").
    Set("cache-control", "no-cache").
    Accept("application/json").
    JSON()
  return err
}

func (c *Connection) sendGetRequest(address string, url string) (error){
  c.Request = request.New()
  _, err := c.Request.
    Get(address + url).
    Timeout(c.GetTimeout()).
    Set("content-type", "application/json").
    Set("cache-control", "no-cache").
    Accept("application/json").
    JSON()
  return err
}

func (c *Connection) sendSecureGetRequest(address string, url string) (error) {
  c.Request = request.New()
  _, err := c.Request.
    Get(address + url).
    Timeout(c.GetTimeout()).
		Set("Authorization", c.GetToken()).
    Set("content-type", "application/json").
    Set("cache-control", "no-cache").
    Accept("application/json").
    Text()
  return err
}

func (c *Connection) GetResponseReason() string {
  return c.Request.ResponseReason()
}

func (c *Connection) GetResponseOK() bool {
  return c.Request.ResponseOK()
}

func (c *Connection) GetResponseURL() string {
  return c.Request.ResponseURL()
}

func (c *Connection) GetResponseCode() int {
  return c.Request.ResponseCode()
}

func (c *Connection) GetResponseData() interface{} {
  return c.Request.ResponseData()
}

func (c *Connection) GetResponseItem(item string) interface{} {
  return c.Request.ResponseDataItem(item)
}

func (c *Connection) GetResponseItemInt(item string) int {
  return c.Request.ResponseDataItemInt(item)
}

func (c *Connection) GetResponseItemInt64(item string) int64 {
  return c.Request.ResponseDataItemInt64(item)
}

func (c *Connection) GetResponseItemFloat(item string) float64 {
  return c.Request.ResponseDataItemFloat(item)
}

func (c *Connection) GetResponseItemString(item string) string {
  return c.Request.ResponseDataItemString(item)
}

func (c *Connection) GetResponseItemBool(item string) bool {
  return c.Request.ResponseDataItemBool(item)
}

func (c *Connection) GetResponseItemTime(item string) time.Time {
  return c.Request.ResponseDataItemTime(item)
}

func (c *Connection) Authenticate() (error){
  if c.settingsMade() {
    data := map[string]string{"username": c.Credentials.Username,"password":c.Credentials.Password,"client_id": c.Credentials.Client_id,"client_secret": c.Credentials.Client_secret,"redirect_uri": c.Credentials.Url}
    e := c.sendPostRequest(c.Credentials.Url, c.hxurls.Authentication, data)
    return e
  }
  return errors.New("Settings need to be updated before requests can be made.")
}

func (c *Connection) About() (error) {
  if c.settingsMade() {
    e := c.sendGetRequest(c.Credentials.Url, c.hxurls.About)
    return e
  }
  return errors.New("Settings need to be updated before requests can be made.")
}

func (c *Connection) ClusterInfo() (error) {
  if c.settingsMade() {
    e := c.sendSecureGetRequest(c.Credentials.Url, c.hxurls.ClusterInfo)
    return e
  }
  return errors.New("Settings need to be updated before requests can be made.")
}

func (c *Connection) ClusterSavings() (error) {
  if c.settingsMade() {
    e := c.sendSecureGetRequest(c.Credentials.Url, c.hxurls.ClusterSavings)
    return e
  }
  return errors.New("Settings need to be updated before requests can be made.")
}

func (c *Connection) ClusterPlatform() (error) {
  if c.settingsMade() {
    e := c.sendSecureGetRequest(c.Credentials.Url, c.hxurls.ClusterPlatform)
    return e
  }
  return errors.New("Settings need to be updated before requests can be made.")
}

func (c *Connection) ClusterVM() (error) {
  if c.settingsMade() {
    e := c.sendSecureGetRequest(c.Credentials.Url, c.hxurls.ClusterVM)
    return e
  }
  return errors.New("Settings need to be updated before requests can be made.")
}

func (c *Connection) ClusterDatastores() (error) {
  if c.settingsMade() {
    e := c.sendSecureGetRequest(c.Credentials.Url, c.hxurls.ClusterDatastores)
    return e
  }
  return errors.New("Settings need to be updated before requests can be made.")
}

func (c *Connection) ClusterSummary() (error) {
  if c.settingsMade() {
    e := c.sendSecureGetRequest(c.Credentials.Url, c.hxurls.ClusterSummary)
    return e
  }
  return errors.New("Settings need to be updated before requests can be made.")
}

func (c *Connection) ClusterAppliances() (error) {
  if c.settingsMade() {
    e := c.sendSecureGetRequest(c.Credentials.Url, c.hxurls.ClusterAppliances)
    return e
  }
  return errors.New("Settings need to be updated before requests can be made.")
}

func (c *Connection) ClusterVersion() (error) {
  if c.settingsMade() {
    e := c.sendSecurePostRequest(c.Credentials.Url, c.hxurls.ClusterVersion, make(map[string]string))
    return e
  }
  return errors.New("Settings need to be updated before requests can be made.")
}