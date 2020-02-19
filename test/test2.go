package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func fetchURL(url string) (string, error) {
	cmd := exec.Command("curl", url)
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func main() {
	res, err := fetchURL("http://localhost:8000/api/v1/nni/metric-data")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
