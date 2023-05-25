package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dghubble/oauth1"
	twitter "github.com/g8rswimmer/go-twitter/v2"
)

type Config struct {
	Timer int    `json:"timer"`
	Dir   string `json:"dir"`
}

type Secret struct {
	BearerToken    string `json:"bearer_token"`
	ConsumerToken  string `json:"consumer_token"`
	ConsumerSecret string `json:"consumer_secret"`
	UserToken      string `json:"user_token"`
	UserSecret     string `json:"user_secret"`
	Tweet          string `json:"tweet"`
}

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func tweet() {
	// Read the configuration file
	file, err := ioutil.ReadFile("twitter-secrets.json")
	if err != nil {
		fmt.Println("Error reading configuration file:", err)
		return
	}

	// Unmarshal the JSON data into a Config struct
	var secret Secret
	err = json.Unmarshal(file, &secret)
	if err != nil {
		fmt.Println("Error unmarshaling configuration file:", err)
		return
	}

	twitterConfig := oauth1.NewConfig(secret.ConsumerToken, secret.ConsumerSecret)
	httpClient := twitterConfig.Client(oauth1.NoContext, &oauth1.Token{
		Token:       secret.UserToken,
		TokenSecret: secret.UserSecret,
	})

	client := &twitter.Client{
		Authorizer: authorize{
			Token: *&secret.BearerToken,
		},
		Client: httpClient,
		Host:   "https://api.twitter.com",
	}

	req := twitter.CreateTweetRequest{
		Text: *&secret.Tweet,
	}

	tweetResponse, err := client.CreateTweet(context.Background(), req)
	if err != nil {
		log.Panicf("create tweet error: %v", err)
	}

	enc, err := json.MarshalIndent(tweetResponse, "", "    ")
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(string(enc))
}

func main() {

	// Read the JSON file
	filePath := "config.json" // Replace with the actual file path
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Parse JSON into Config struct
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	usbDir := "dummy/dev/bus/usb"
	var initialCount int

	// Get initial count of files in usbDir
	initialCount = countFiles(usbDir)

	fmt.Println("dms active")

	// Loop every 5 seconds to check if the number of files in usbDir has changed
	for {
		time.Sleep(time.Duration(config.Timer) * time.Second)
		currentCount := countFiles(usbDir)
		if currentCount != initialCount {
			// Delete secretsDir if the number of files in usbDir has changed
			err := os.RemoveAll(config.Dir)
			tweet()
			fmt.Println("dms triggered")
			if err != nil {
				fmt.Println(err)
			}
			break
		}
	}

}

// countFiles counts the number of files in a directory
func countFiles(dir string) int {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		fmt.Println(err)
	}
	return len(files)
}
