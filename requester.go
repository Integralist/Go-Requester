package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

type Component struct {
	Id  string `yaml:"id"`
	Url string `yaml:"url"`
}

type ComponentsList struct {
	Components []Component `yaml:"components"`
}

type ComponentResponse struct {
	Id     string `json:"id"`
	Status int    `json:"status"`
	Body   string `json:"body"`
}

type Result struct {
	Status     string              `json:"status"`
	Components []ComponentResponse `json:"components"`
}

var overallStatus string = "success"

func checkError(msg string) int {
	timeout, _ := regexp.MatchString("Timeout", msg)

	if timeout {
		return 408
	} else {
		return 500
	}
}

func getComponent(wg *sync.WaitGroup, client *http.Client, i int, v Component, ch chan ComponentResponse) {
	defer wg.Done()

	resp, err := client.Get(v.Url)

	if err != nil {
		fmt.Printf("Problem getting the response: %s\n\n", err)

		ch <- ComponentResponse{
			v.Id, checkError(err.Error()), err.Error(),
		}
	} else {
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Problem reading the body for %s -> %s\n", v.Id, err)
		}

		ch <- ComponentResponse{
			v.Id, resp.StatusCode, string(contents),
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

func handler(w http.ResponseWriter, r *http.Request) {
	var cr []ComponentResponse
	var y ComponentsList

	ch := make(chan ComponentResponse)

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

	j, err := json.Marshal(Result{overallStatus, cr})
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
