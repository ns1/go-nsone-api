package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"gopkg.in/ns1/ns1-go.v2/rest/model/pulsar"

	api "gopkg.in/ns1/ns1-go.v2/rest"
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
	// create applications
	var app *pulsar.Application
	var apps []*pulsar.Application
	for i := 0; i < 2; i++ {
		app = pulsar.NewApplication("MyAPP_" + strconv.Itoa(i))
		app.DefaultConfig = pulsar.DefaultConfig{
			Http:                 false,
			Https:                false,
			RequestTimeoutMillis: 100 + 10*i,
			JobTimeoutMillis:     100 + 10*i,
			UseXhr:               false,
			StaticValues:         false,
		}
		app.Active = false
		app.BrowserWaitMillis = 100
		app.JobsPerTransaction = 100

		_, err := client.Applications.Create(app)
		if err != nil {
			panic("could not create " + app.Name)
			return
		}
		apps = append(apps, app)
	}

	// list all
	err := printAllApplications()
	if err != nil {
		return
	}

	// print all by id
	err = printAllByID(apps)
	if err != nil {
		return
	}

	// update
	for i, v := range apps {
		v.Name = "MyAPP_Updated_" + strconv.Itoa(i)
		_, err = client.Applications.Update(v)

		if err != nil {
			panic("could not get application" + err.Error())
			return
		}
	}

	// print all by id
	err = printAllByID(apps)
	if err != nil {
		return
	}

	//delete all
	for _, v := range apps {
		_, err := client.Applications.Delete(v.ID)
		if err != nil {
			panic("could not delete " + v.ID + " " + err.Error())
			return
		}
	}
}

func printAllApplications() error {
	appsList, _, err := client.Applications.List()
	fmt.Println("\n#### printing all ####")
	if err != nil {
		panic("could not list all applications " + err.Error())
		return err
	}
	for _, v := range appsList {
		fmt.Println("%w", v)
	}
	fmt.Println("")
	return nil
}

func printAllByID(apps []*pulsar.Application) error {
	var appsReturned []*pulsar.Application
	for _, v := range apps {
		app, _, err := client.Applications.Get(v.ID)
		if err != nil {
			panic("could not get application" + err.Error())
			return err
		}
		appsReturned = append(appsReturned, app)
	}
	fmt.Println("\n#### printing all by id ####")
	for _, v := range appsReturned {
		fmt.Println("%w", v)
	}
	fmt.Println("")
	return nil
}
