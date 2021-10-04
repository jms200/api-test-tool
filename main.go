package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	endpoint, payload, requestType, targetURL, token *string
	required                                         = []string{"endpoint", "requestType", "targetURL", "token"}
)

func init() {
	endpoint = flag.String("endpoint", "", "Target API endpoint, example: /organizations/:organization_name/workspaces")
	payload = flag.String("payload", "./payload.json", "Path to the payload file.")
	requestType = flag.String("requestType", "", "Type of API request to perform: GET|PUT")
	targetURL = flag.String("targetURL", "https://app.terraform.io/api/v2", "Intended target URL for API, defaults to: https://app.terraform.io/api/v2/")
	token = flag.String("token", os.Getenv("TOKEN"), "API token, defaults to pulling from TOKEN envronment variable")
}

func checkReqFlags() []string {
	var missing []string

	flag.Parse()

	seen := make(map[string]bool)
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() != "" {
			seen[f.Name] = true
		}
	})

	for _, req := range required {
		if !seen[req] {
			// or possibly use `log.Fatalf` instead of:
			missing = append(missing, req)
		}
	}

	return missing
}

func buildRequest() (*http.Request, error) {
	if *requestType == "POST" {
		return buildPostRequest()
	}

	return buildGetRequest()
}

func buildGetRequest() (*http.Request, error) {
	req, err := http.NewRequest("GET", *targetURL+*endpoint, nil)
	req.Header.Add("Authorization", "Bearer "+*token)

	return req, err
}

func buildPostRequest() (*http.Request, error) {
	data, err := os.Open(*payload)

	if err != nil {
		fmt.Printf("Building PUT request encountered an error: %s", err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", *targetURL+*endpoint, data)
	req.Header.Add("Authorization", "Bearer "+*token)
	req.Header.Add("Content-Type", "application/vnd.api+json")

	return req, err
}

func getResponse(req *http.Request) *http.Response {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error on response.\n[ERROR] -", err)
		os.Exit(2)
	}

	return resp
}

func callAPI(req *http.Request) ([]byte, error) {
	var body []byte
	var err error

	body, err = ioutil.ReadAll(getResponse(req).Body)
	return body, err
}

func main() {
	var request *http.Request
	var err error

	flag.Parse()

	missingFlags := checkReqFlags()

	if len(missingFlags) != 0 {
		fmt.Println("The following flag(s) are required, but missing: " + strings.Join(missingFlags, ", "))
		os.Exit(1)
	}

	request, err = buildRequest()

	responseBody, err := callAPI(request)
	if err != nil {
		log.Println(string([]byte(responseBody)))
	}
	log.Println(string([]byte(responseBody)))
}
