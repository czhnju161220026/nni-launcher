package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cuizihan/launcher/typed"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// PodIPCache给出了任务名到Pod ip到映射
type NNILauncher struct {
	Clientset *kubernetes.Clientset
}

func (l *NNILauncher) SubmitExperiment(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Submit called")
	w.Header().Set("Content-type", "application/json")
	body, _ := ioutil.ReadAll(r.Body)
	var experiment typed.NNIExperiment
	err := json.Unmarshal(body, &experiment)
	experiment.ExpID = generateId()
	if err != nil {
		fmt.Fprintf(os.Stderr, "job info error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"Bad request"}`)
		return
	}
	_, err = experiment.CreatePod(l.Clientset)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pod create error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"Pod create error"}`)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Submit called"}`)
}

func (l *NNILauncher) GetLog(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetLog called")
	w.Header().Set("Content-type", "application/json")
	// workspaceName, userName := r.URL.Query().Get("workspace"), r.URL.Query().Get("user")
	body, _ := ioutil.ReadAll(r.Body)
	var info typed.QueryInfo
	err := json.Unmarshal(body, &info)
	if err != nil {
		fmt.Fprintf(os.Stderr, "json unmarshal error:%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expPods, err := l.Clientset.CoreV1().Pods("nni-exp").List(metav1.ListOptions{
		LabelSelector: labels.Set{"user": info.User, "workspace": info.Workspace}.String(),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "pod list error:%v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var expInfo []string
	for _, pod := range expPods.Items {
		if pod.Status.Phase == "Running" && pod.Status.PodIP != "" {
			data, err := fetchData("http://" + pod.Status.PodIP + ":8000")
			if err != nil {
				fmt.Fprintf(os.Stderr, "fetch pod:%s error: %v\n", pod.Name, err)
				data = "{}"
			}
			expInfo = append(expInfo, data)
		}
	}
	res := "["
	for i, data := range expInfo {
		res += data
		if i != len(expInfo)-1 {
			res += ","
		}
	}
	res += "]"
	fmt.Fprint(w, res)
}

// func (l *NNILauncher) DeleteExperiment(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Delete called")
// 	w.Header().Set("Content-type", "application/json")
// 	podName := r.URL.Query().Get("workspace") + "-" + r.URL.Query().Get("id")
// 	err := l.Clientset.CoreV1().Pods("nnijob").Delete(podName, &metav1.DeleteOptions{})
// 	if err != nil {
// 		w.WriteHeader(http.StatusNotFound)
// 		fmt.Fprintf(os.Stderr, "delete error: %v\n", err)
// 		fmt.Fprint(w, `{"message":"delete pod failed"}`)
// 		return
// 	}
// 	fmt.Printf("Delete pod:%s\n", podName)
// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprint(w, `{"message":"Delete succeed"}`)
// }

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

var alphabet = "abcdefghijklmnopqrstuvwxyz1234567890"

func generateId() string {
	rand.Seed(time.Now().Unix())
	id := make([]byte, 8)
	for i := 0; i < 8; i++ {
		id[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(id)
}
