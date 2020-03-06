package main

import (
	"fmt"
	myrest "github.com/cuizihan/launcher/handler"
	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"net/http"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	nnilauncher := myrest.NNILauncher{Clientset: clientset}
	r := mux.NewRouter()
	// SubmitExperiment和GetLog是主要的接口
	// 分别用于提交调参任务和获取任务的数据
	api := r.PathPrefix("/api/v1/nni-exp").Subrouter()
	api.HandleFunc("/submit", nnilauncher.SubmitExperiment).Methods(http.MethodPost)
	api.HandleFunc("/logs", nnilauncher.GetLog).Methods(http.MethodGet)
	api.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "%s\n", "Hello world")
	})

	addr := ":8000"
	fmt.Printf("Listen on %s.\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
