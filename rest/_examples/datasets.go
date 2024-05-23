package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/ns1/ns1-go.v2/rest/model/dataset"

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
	dt, _, err := client.Datasets.Create(&dataset.Dataset{
		Name: "My dataset",
		Datatype: &dataset.Datatype{
			Type:  dataset.DatatypeTypeNumQueries,
			Scope: dataset.DatatypeScopeAccount,
			Data:  nil,
		},
		Repeat: nil,
		Timeframe: &dataset.Timeframe{
			Aggregation: dataset.TimeframeAggregationMontly,
			Cycles:      func() *int32 { i := int32(1); return &i }(),
		},
		ExportType:      dataset.ExportTypeCSV,
		RecipientEmails: nil,
	})
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.MarshalIndent(dt, "", "  ")
	fmt.Println(string(b))

	for {
		fmt.Println("waiting for report to be generated...")

		dt, _, err = client.Datasets.Get(dt.ID)
		if err != nil {
			log.Fatal(err)
		}
		if len(dt.Reports) == 0 {
			time.Sleep(time.Second * 5)
			continue
		}

		dtStatus := dt.Reports[0].Status

		if dtStatus == dataset.ReportStatusGenerating || dtStatus == dataset.ReportStatusQueued {
			time.Sleep(time.Second * 5)
			continue
		}

		if dtStatus == dataset.ReportStatusFailed {
			log.Fatal(fmt.Errorf("dataset[%s] report[%s] failed to generate", dt.ID, dt.Reports[0].ID))
		}

		// report generated
		break
	}

	reportBuf, reportResp, err := client.Datasets.GetReport(dt.ID, dt.Reports[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	var fileName string
	contentDisposition := reportResp.Header.Get("Content-Disposition")
	params := strings.Split(contentDisposition, ";")
	for _, param := range params {
		param = strings.TrimSpace(param)
		if strings.HasPrefix(param, "filename=") {
			fileName = strings.TrimPrefix(param, "filename=")
			fileName = strings.TrimSpace(fileName)
			fileName = strings.Trim(fileName, "\"")
			break
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(fmt.Sprintf("%s/%s", cwd, fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, reportBuf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("dataset report saved to", fileName)
}
