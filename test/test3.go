package main

import (
	"encoding/json"
	"fmt"
	"github.com/cuizihan/launcher/typed"
	"io/ioutil"
	"net/http"
)

func fetchURL2(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	var trails []typed.NNITrial
	err = json.Unmarshal(content, &trails)
	if err != nil {
		return "", err
	}
	res := "["
	for i, trail := range trails {
		res += trail.ToJSON()
		if i != len(trails)-1 {
			res += ","
		}
	}
	res += "]"
	return res, nil
}
func main() {
	res, _ := fetchURL2("http://localhost:8000/api/v1/nni/trial-jobs/")
	fmt.Println(res)
}
