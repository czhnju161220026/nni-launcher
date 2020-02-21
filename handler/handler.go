package handler

import (
	"encoding/json"
	"fmt"
	"github.com/cuizihan/launcher/typed"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// PodIPCache给出了任务名到Pod ip到映射
type NNILauncher struct {
	Clientset  *kubernetes.Clientset
	PodIPCache map[string]string
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
	fmt.Fprintf(w, `{"message":"Submit called","id":`)
	fmt.Fprintf(w, "%q}", experiment.ExpID)
}

func (l *NNILauncher) GetLog(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetLog called")
	w.Header().Set("Content-type", "application/json")
	podName := r.URL.Query().Get("workspace") + "-" + r.URL.Query().Get("id")
	if ip, ok := l.PodIPCache[podName]; ok {
		url := "http://" + ip + ":8000/api/v1/nni/trial-jobs/"
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
		return
	}
	pod, err := l.Clientset.CoreV1().Pods("nnijob").Get(podName, metav1.GetOptions{})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"message":"get pod failed"}`)
		return
	}
	fmt.Printf("Find pod:%s, podIp:%s\n", pod.Name, pod.Status.PodIP)
	if pod.Status.PodIP == "" {
		fmt.Println("Waiting")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"message":"Waiting"}`)
		return
	}
	l.PodIPCache[podName] = pod.Status.PodIP
	url := "http://" + pod.Status.PodIP + ":8000/api/v1/nni/trial-jobs/"
	content, err := fetchURL(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Get log failed:%v\n", err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"message":"Get log failed"}`)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"message":"GetLog succeed", "content":`)
	fmt.Fprint(w, content)
	fmt.Fprint(w, "}")
}

func (l *NNILauncher) DeleteExperiment(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete called")
	w.Header().Set("Content-type", "application/json")
	podName := r.URL.Query().Get("workspace") + "-" + r.URL.Query().Get("id")
	err := l.Clientset.CoreV1().Pods("nnijob").Delete(podName, &metav1.DeleteOptions{})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(os.Stderr, "delete error: %v\n", err)
		fmt.Fprint(w, `{"message":"delete pod failed"}`)
		return
	}
	fmt.Printf("Delete pod:%s\n", podName)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"message":"Delete succeed"}`)
}

func fetchURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadAll(resp.Body)
	// fmt.Printf("Raw: %s\n", string(content))
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	var trails []typed.NNITrial
	err = json.Unmarshal(content, &trails)
	if err != nil {
		return "", err
	}
	// fmt.Printf("Trails:%v\n", trails)
	res := "["
	for i, trail := range trails {
		temp := trail.ToJSON()
		// fmt.Printf("Trail:%s\n", temp)
		res += temp
		if i != len(trails)-1 {
			res += ","
		}
	}
	res += "]"
	return res, nil
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
