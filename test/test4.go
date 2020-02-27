package main

import (
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func homeDir2() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
func main() {
	var kubeconfig *string
	if home := homeDir2(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	_, err = clientset.CoreV1().Pods("default").Create(&v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test2",
			Namespace: "default",
			Labels: map[string]string{
				"workspace": "test-ws",
				"user":      "hahaha",
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    "test",
					Image:   "ubuntu:latest",
					Command: []string{"/bin/bash", "-c", "--"},
					Args:    []string{"while true; do sleep 30; done;"},
				},
			},
			RestartPolicy: "OnFailure",
		},
	})

	if err != nil {
		fmt.Println("Create failed")
	}
}
