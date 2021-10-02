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
	endpointURL, token, requestType *string
)

func init() {
	endpointURL = flag.String("endpointURL", "", "URL for the API endpoint")
	requestType = flag.String("requestType", "", "Type of API request to perform: GET|PUT")
	token = flag.String("token", os.Getenv("TOKEN"), "API token, defaults to pulling from TOKEN envronment variable")
}

func checkReqFlags() []string {
	required := []string{"endpointURL", "requestType", "token"}
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
	req, err := http.NewRequest(*requestType, *endpointURL, nil)
	req.Header.Add("Authorization", "Bearer "+*token)

	return req, err
}

func callAPI(req *http.Request) ([]byte, error) {
	var body []byte
	var err error

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	return body, err
}

func main() {
	var request *http.Request
	var err error

	flag.Parse()

	missingFlags := checkReqFlags()

	if len(missingFlags) != 0 {
		fmt.Println("The following flag(s) are required, but missing: " + strings.Join(missingFlags, ", "))
		os.Exit(2)
	}

	request, err = buildRequest()

	responseBody, err := callAPI(request)
	if err != nil {
		log.Println(string([]byte(responseBody)))
	}
	log.Println(string([]byte(responseBody)))

	fmt.Printf("Request is type: %T\n", request)
	fmt.Println(err)
}
