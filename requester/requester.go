package requester

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/integralist/go-requester/aggregator"
	"github.com/integralist/go-requester/config"
	"gopkg.in/yaml.v2"
)

var wg sync.WaitGroup
var ch chan aggregator.ComponentResponse

type component struct {
	ID        string `yaml:"id"`
	URL       string `yaml:"url"`
	Mandatory bool   `yaml:"mandatory"`
}

type componentsList struct {
	Components []component `yaml:"components"`
}

func checkError(err error) int {
	if e, ok := err.(net.Error); ok && e.Timeout() {
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

func getComponent(i int, v component) {
	defer wg.Done()

	timeout := time.Duration(1 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(v.URL)

	if err != nil {
		fmt.Printf("Problem getting the response: %s\n\n", err)
		status := checkError(err)

		ch <- aggregator.ComponentResponse{
			ID:        v.ID,
			Status:    status,
			Body:      err.Error(),
			Summary:   getSummary(status),
			Mandatory: v.Mandatory,
		}
	} else {
		defer resp.Body.Close()

		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Problem reading the body for %s -> %s\n", v.ID, err)
		}

		ch <- aggregator.ComponentResponse{
			ID:        v.ID,
			Status:    resp.StatusCode,
			Body:      string(contents),
			Summary:   getSummary(resp.StatusCode),
			Mandatory: v.Mandatory,
		}
	}
}

// Process function kick starts the concurrent requesting of components
func Process(w http.ResponseWriter, r *http.Request, configPath string) {
	var cr []aggregator.ComponentResponse
	var components componentsList

	config := config.Parse(configPath)
	yaml.Unmarshal(config, &components)

	ch = make(chan aggregator.ComponentResponse, len(components.Components))

	for i, v := range components.Components {
		wg.Add(1)
		go getComponent(i, v)
	}
	wg.Wait()
	close(ch)

	for component := range ch {
		cr = append(cr, component)
	}

	response, err := aggregator.Process(cr)
	if err != nil {
		fmt.Printf("There was an error aggregating a response: %s", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
