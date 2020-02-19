package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cuizihan/launcher/typed"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"os"
	"os/exec"
)

// PodIP给出了任务名到Pod ip到映射
type NNILauncher struct {
	Clientset *kubernetes.Clientset
}

func (l *NNILauncher) SubmitExperiment(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Submit called")
	w.Header().Set("Content-type", "application/json")
	body, _ := ioutil.ReadAll(r.Body)
	var experiment typed.NNIExperiment
	err := json.Unmarshal(body, &experiment)
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
	item := r.URL.Query().Get("id")
	pod, err := l.Clientset.CoreV1().Pods("nnijob").Get(item, metav1.GetOptions{})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"message":"get pod failed"}`)
		return
	}
	fmt.Printf("Find pod:%s, podIp:%s\n", pod.Name, pod.Status.PodIP)
	if pod.Status.PodIP == "" {
		fmt.Println("Waiting")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message":"Waiting"}`)
		return
	}
	url := "http://" + pod.Status.PodIP + ":8000/api/v1/nni/metric-data/"
	content, err := fetchURL(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Get log failed:%v\n", err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"message":"Get log failed"}`)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"message":"GetLog succeed", "content":`)
	fmt.Fprint(w, content)
	fmt.Fprint(w, "}")
}

func (l *NNILauncher) DeleteExperiment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(body))
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Delete called"}`)
}

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
