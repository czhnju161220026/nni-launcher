package main

import (
	"flag"
	"fmt"
	myrest "github.com/cuizihan/launcher/handler"
	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"net/http"
)

func main() {
	// var kubeconfig *string
	// if home := homeDir(); home != "" {
	// 	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	// } else {
	// 	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	// }
	var ip = flag.String("h", "", "specify the host ip")
	var port = flag.String("p", "8000", "sepcify the port")
	flag.Parse()
	// use the current context in kubeconfig
	config, err := rest.InClusterConfig()
	// config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
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
	api := r.PathPrefix("/api/v1/nni-exp").Subrouter()
	api.HandleFunc("/submit", nnilauncher.SubmitExperiment).Methods(http.MethodPost)
	api.HandleFunc("/logs", nnilauncher.GetLog).Methods(http.MethodPost)
	// api.HandleFunc("", nnilauncher.DeleteExperiment).Methods(http.MethodDelete)
	api.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "%s", "Hello world")
	})

	addr := *ip + ":" + *port
	fmt.Printf("Listen on %s.\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))

}
