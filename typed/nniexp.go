package typed

import (
	"encoding/json"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
)

type NNIExperiment struct {
	User        string      `json:"user"`
	WorkSpace   string      `json:"workspace"`
	GPU         int         `json:"gpuNum"`
	Concurrency int         `json:"trailConcurrency"`
	Target      string      `json:"target"`
	CMD         string      `json:"command"`
	SearchSpace interface{} `json:"search_space"`
	Num         int         `json:"num"`
	ExpID       string
}

func (j NNIExperiment) String() string {
	res := fmt.Sprintf("User:%s\nWorkspace:%s\nGPU:%s\nConcurrcy:%s\nTarget:%s\nCMD:%s\nNUM:%s\nSearchSpace:%s\n",
		j.User, j.WorkSpace, string(j.GPU), string(j.Concurrency), j.Target, j.CMD, j.Num, j.GetSearchSpaceJson())
	return res
}

func (j NNIExperiment) GetSearchSpaceJson() string {
	res, err := json.MarshalIndent(j.SearchSpace, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "get searchspacejson error:%v\n", err)
	}
	return string(res)
}

func (j NNIExperiment) CreatePod(clientset *kubernetes.Clientset) (*apiv1.Pod, error) {
	podsClient := clientset.CoreV1().Pods("nni-exp")
	newPod := &apiv1.Pod{
		TypeMeta: v1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      j.User + "-" + j.WorkSpace + "-" + j.ExpID,
			Namespace: "nni-exp",
			Labels: map[string]string{
				"user":      j.User,
				"workspace": j.WorkSpace,
			},
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Ports:   []apiv1.ContainerPort{{ContainerPort: 8000}},
					Name:    "worker",
					Image:   "czh1998/nnidemo:0.3",
					Command: []string{"python", "-u", "entry.py"},
					Env: []apiv1.EnvVar{
						{
							Name:  "GPU_NUM",
							Value: fmt.Sprintf("%d", j.GPU),
						},
						{
							Name:  "COMMAND",
							Value: j.CMD,
						},
						{
							Name:  "TARGET",
							Value: j.Target,
						},
						{
							Name:  "USER",
							Value: j.User,
						},
						{
							Name:  "SEARCH_SPACE",
							Value: j.GetSearchSpaceJson(),
						},
						{
							Name:  "CONCURRENCY",
							Value: fmt.Sprintf("%d", j.Concurrency),
						},
						{
							Name:  "NUM",
							Value: fmt.Sprintf("%d", j.Num),
						},
						{
							Name:  "PYTHONUNBUFFERED",
							Value: "0",
						},
					},
				},
			},
			RestartPolicy: "OnFailure",
		},
	}

	return podsClient.Create(newPod)
}
