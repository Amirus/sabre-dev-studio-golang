package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/clientcredentials"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

func clientID() (string, error) {
	clientID := os.Getenv("CLIENT_ID")
	clientID = base64.StdEncoding.EncodeToString([]byte(clientID))
	return clientID, nil
}
func clientSecret() (string, error) {
	clientSecret := os.Getenv("CLIENT_SECRET")
	clientSecret = base64.StdEncoding.EncodeToString([]byte(clientSecret))
	return clientSecret, nil
}
func baseUrl() (string, error) {
	baseUrl := os.Getenv("URL")
	return baseUrl, nil
}

type DevStudioApiClient struct {
	*http.Client
	BaseUrl string
}

func NewClient() *DevStudioApiClient {
	// Shout out to https://www.snip2code.com/Snippet/551369/Example-usage-of-https---godoc-org-golan
	baseUrl, _ := baseUrl()
	clientID, _ := clientID()
	clientSecret, _ := clientSecret()
	config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     baseUrl + "/v1/auth/token",
	}
	// the client will update its token if it's expired
	client := config.Client(context.Background())
	return &DevStudioApiClient{Client: client, BaseUrl: baseUrl}
}
func (c *DevStudioApiClient) Request(requestUrl string) {
	fmt.Printf("+%v\n", requestUrl)
	resp, err := c.Get(requestUrl)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	// do something with resp
	content, _ := ioutil.ReadAll(resp.Body)
	var f interface{}
	_ = json.Unmarshal(content, &f)
	prettyJSON, _ := json.MarshalIndent(f, "", "  ")
	os.Stdout.Write(prettyJSON)
}
func (c *DevStudioApiClient) RequestWithParams(requestUrl string, params map[string]string) {
	q := url.Values{}
	for key, value := range params {
		q.Add(key, value)
	}
	requestUrl = requestUrl + "?" + q.Encode()
	c.Request(requestUrl)
}
func (c *DevStudioApiClient) GetTravelThemes() {
	travelThemesUrl := c.BaseUrl + "/v1/lists/supported/shop/themes"
	c.Request(travelThemesUrl)
}
func (c *DevStudioApiClient) GetFlightSearch() {
	flightSearchUrl := c.BaseUrl + "/v1/shop/flights"
	params := map[string]string{
		"origin":              "DFW",
		"destination":         "NYC",
		"departuredate":       "2015-10-01",
		"returndate":          "2015-10-04",
		"limit":               "1",
		"outboundflightstops": "2",
		"inboundflightstops":  "2",
		"excludecarriers":     "NK",
	}
	c.RequestWithParams(flightSearchUrl, params)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	client := NewClient()
	//client.GetTravelThemes()
	client.GetFlightSearch()
}
