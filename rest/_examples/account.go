package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/account"
)

var client *api.Client

// Helper that initializes rest api client from environment variable.
func init() {
	k := os.Getenv("NS1_APIKEY")
	if k == "" {
		fmt.Println("NS1_APIKEY environment variable is not set, giving up")
	}

	httpClient := &http.Client{Timeout: time.Second * 10}
	// Adds logging to each http request.
	doer := api.Decorate(httpClient, api.Logging(log.New(os.Stdout, "", log.LstdFlags)))
	client = api.NewClient(doer, api.SetAPIKey(k))
}

func main() {
	teams, _, err := client.Teams.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range teams {
		b, _ := json.MarshalIndent(t, "", "  ")
		fmt.Println(string(b))
	}

	users, _, err := client.Users.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range users {
		b, _ := json.MarshalIndent(u, "", "  ")
		fmt.Println(string(b))
	}

	keys, _, err := client.APIKeys.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, k := range keys {
		b, _ := json.MarshalIndent(k, "", "  ")
		fmt.Println(string(b))
	}

	key := account.APIKey{
		Name:              "new-example",
		IPWhitelist:       []string{"1.1.1.1", "2.2.2.2"},
		IPWhitelistStrict: true,
	}

	if _, err := client.APIKeys.Create(&key); err != nil {
		log.Fatal(err)
	}

	activity, _, err := client.Activity.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, a := range activity {
		b, _ := json.MarshalIndent(a, "", "  ")
		fmt.Println(string(b))
	}
}
