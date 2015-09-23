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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	baseUrl, _ := baseUrl()
	// Shout out to https://www.snip2code.com/Snippet/551369/Example-usage-of-https---godoc-org-golan
	clientID, _ := clientID()
	clientSecret, _ := clientSecret()
	config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     baseUrl + "/v1/auth/token",
	}
	client := config.Client(context.Background())

	// the client will update its token if it's expired
	flightSearchUrl := baseUrl + "/v1/shop/flights"
	params := url.Values{}
	params.Add("origin", "DFW")
	params.Add("destination", "NYC")
	params.Add("departuredate", "2015-10-01")
	params.Add("returndate", "2015-10-04")
	params.Add("limit", "500")
	params.Add("outboundflightstops", "2")
	params.Add("inboundflightstops", "2")
	params.Add("excludedcarriers", "NK")
	flightSearchUrl = flightSearchUrl + "?" + params.Encode()
	fmt.Printf("+%v\n", flightSearchUrl)
	resp, err := client.Get(flightSearchUrl)
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
