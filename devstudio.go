package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/clientcredentials"
	"io/ioutil"
	"log"
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
func url() (string, error) {
	url := os.Getenv("URL")
	return url, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	url, _ := url()
	// Shout out to https://www.snip2code.com/Snippet/551369/Example-usage-of-https---godoc-org-golan
	clientID, _ := clientID()
	clientSecret, _ := clientSecret()
	config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     url + "/v1/auth/token",
	}
	// you can modify the client (for example ignoring bad certs or otherwise)
	// by modifying the context
	client := config.Client(context.Background())

	// the client will update its token if it's expired
	resp, err := client.Get(url + "/v1/lists/supported/shop/themes")
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
