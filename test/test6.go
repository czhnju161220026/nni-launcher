package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func fetchData(ip string) (string, error) {
	routes := []string{
		"/api/v1/nni/metric-data",
		"/api/v1/nni/trial-jobs",
		"/api/v1/nni/check-status",
		"/api/v1/nni/experiment",
	}
	done := make(chan int)
	res := make([]string, 4)
	for i, route := range routes {
		go func(url string, index int) {
			resp, err := http.Get(url)
			if err != nil {
				done <- -1
				return
			}
			content, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				done <- -1
				return
			}
			res[index] = string(content)
			done <- 1
		}(ip+route, i)
	}

	for i := 0; i < 4; i++ {
		x := <-done
		if x != 1 {
			return "", errors.New("get url failed")
		}
	}

	data := res[3][:len(res[3])-1]
	data = data + ",\n" + `"check_status":` + res[2] + ",\n" + `"trial_jobs":` +
		res[1] + ",\n" + `"metric_data":` + res[0] + "}"

	return data, nil
}

func main() {
	res, err := fetchData("http://localhost:1234")
	if err != nil {
		fmt.Printf("Error:%v\n", err)
	}
	fmt.Println(res)
}
