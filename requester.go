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

func main() {
	var cr []ComponentResponse
	var c ComponentsList

	b := getComponents()

	json.Unmarshal(b, &c)

	// fmt.Println(c)
	// fmt.Println("First Id:", c.Components[0].Id)
	// fmt.Println("First Url:", c.Components[0].Url)
	// fmt.Println("Number of components", len(c.Components))

	var wg sync.WaitGroup

	timeout := time.Duration(1 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	for i, v := range c.Components {
		wg.Add(1)

		go func(i int, v Component) {
			defer wg.Done()
			// fmt.Printf("index = %d; id = %s; value = %s\n", i, v.Id, v.Url)

			resp, err := client.Get(v.Url)

			// fmt.Printf("%+v", resp)

			if err != nil {
				fmt.Printf("Problem getting the response: %s\n", err)

				cr = append(cr, ComponentResponse{
					v.Id,
					404,
					err.Error(),
				})
			} else {
				defer resp.Body.Close()
				contents, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("Problem reading the body: %s\n", err)
				}
				// fmt.Printf("Response body: %s\n", string(contents))

				cr = append(cr, ComponentResponse{
					v.Id,
					resp.StatusCode,
					string(contents),
				})
			}
		}(i, v)
	}
	wg.Wait()

	j, err := json.Marshal(Result{overallStatus, cr})
	if err != nil {
		fmt.Printf("Problem converting to JSON: %s\n", err)
		return
	}

	fmt.Println(string(j))
}
