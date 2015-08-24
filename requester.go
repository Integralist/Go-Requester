package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/integralist/go-requester/Godeps/_workspace/src/gopkg.in/yaml.v2"
)

type component struct {
	ID        string `yaml:"id"`
	URL       string `yaml:"url"`
	Mandatory bool   `yaml:"mandatory"`
}

type componentsList struct {
	Components []component `yaml:"components"`
}

type componentResponse struct {
	ID        string `json:"id"`
	Status    int    `json:"status"`
	Body      string `json:"body"`
	Summary   string `json:"summary"`
	Mandatory bool   `json:"mandatory"`
}

type result struct {
	Summary    string              `json:"summary"`
	Components []componentResponse `json:"components"`
}

func checkError(msg string) int {
	timeout, _ := regexp.MatchString("Timeout", msg)

	if timeout {
		return 408
	}
	return 500
}

func getSummary(status int) string {
	if status == 200 || status == 304 {
		return "success"
	}
	return "failure"
}

func getComponent(wg *sync.WaitGroup, client *http.Client, i int, v component, ch chan componentResponse) {
	defer wg.Done()

	resp, err := client.Get(v.URL)

	if err != nil {
		fmt.Printf("Problem getting the response: %s\n\n", err)
		status := checkError(err.Error())

		ch <- componentResponse{
			v.ID, status, err.Error(), getSummary(status), v.Mandatory,
		}
	} else {
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Problem reading the body for %s -> %s\n", v.ID, err)
		}

		ch <- componentResponse{
			v.ID, resp.StatusCode, string(contents), getSummary(resp.StatusCode), v.Mandatory,
		}
	}
}

func getComponents() []byte {
	config, err := ioutil.ReadFile("./page.yaml")

	if err != nil {
		fmt.Printf("Problem reading the page config: %s\n", err)
		panic(err)
	}

	return config
}

func finalSummary(components []componentResponse) string {
	for _, c := range components {
		if c.Mandatory == true && c.Summary == "failure" {
			return "failure"
		}
	}

	return "success"
}

func handler(w http.ResponseWriter, r *http.Request) {
	var cr []componentResponse
	var y componentsList

	ch := make(chan componentResponse)

	b := getComponents()
	yaml.Unmarshal(b, &y)

	timeout := time.Duration(1 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	var wg sync.WaitGroup
	for i, v := range y.Components {
		wg.Add(1)
		go getComponent(&wg, &client, i, v, ch)
		cr = append(cr, <-ch)
	}
	wg.Wait()

	j, err := json.Marshal(result{finalSummary(cr), cr})
	if err != nil {
		fmt.Printf("Problem converting to JSON: %s\n", err)
		return
	}

	// fmt.Println(string(j))

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
