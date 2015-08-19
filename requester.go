package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Component struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type ComponentsList struct {
	Components []Component `json:"components"`
}

type ComponentResponse struct {
	Id     string
	Status int
	Body   string
}

type Result struct {
	Status     string
	Components []ComponentResponse
}

var overallStatus string = "success"

func getComponents() []byte {
	return []byte(`{"components":[{"id":"local","url":"http://localhost:8080/pugs"},{"id":"google","url":"http://google.com/"},{"id":"integralist","url":"http://integralist.co.uk/"},{"id":"sloooow","url":"http://stevesouders.com/cuzillion/?c0=hj1hfff30_5_f&t=1439194716962"}]}`)
}

func getComponent(wg *sync.WaitGroup, client *http.Client, i int, v Component, ch chan ComponentResponse) {
	defer wg.Done()

	resp, err := client.Get(v.Url)

	if err != nil {
		fmt.Printf("Problem getting the response: %s\n\n", err)

		ch <- ComponentResponse{
			v.Id, 500, err.Error(),
		}
	} else {
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Problem reading the body: %s\n", err)
		}

		ch <- ComponentResponse{
			v.Id, resp.StatusCode, string(contents),
		}
	}
}

func main() {
	var cr []ComponentResponse
	var c ComponentsList

	ch := make(chan ComponentResponse)
	b := getComponents()

	json.Unmarshal(b, &c)

	timeout := time.Duration(1 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	var wg sync.WaitGroup
	for i, v := range c.Components {
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

	fmt.Println(string(j))
}
