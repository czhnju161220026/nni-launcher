package typed

import (
	"encoding/json"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
)

const (
	NfsServer = "210.28.132.167"
	NfsPath   = "/data/nfs/nni_storage"
	IMAGE     = "registry.cn-hangzhou.aliyuncs.com/cuizihan/registry-cuizihan/nni:0.18.1-tf1.14.0-torch1.2.0-mxnet1.5.0-py3.6-nni1.3"
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
	Trainer     string      `json:"trainer"`
	ExpID       string
}

func (e NNIExperiment) String() string {
	res := fmt.Sprintf("User:%s\nWorkspace:%s\nGPU:%s\nConcurrcy:%s\nTarget:%s\nCMD:%s\nNUM:%s\nSearchSpace:%s\n",
		e.User, e.WorkSpace, string(e.GPU), string(e.Concurrency), e.Target, e.CMD, e.Num, e.GetSearchSpaceJson())
	return res
}

func (e NNIExperiment) GetSearchSpaceJson() string {
	res, err := json.MarshalIndent(e.SearchSpace, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "get searchspacejson error:%v\n", err)
	}
	return string(res)
}

func (e NNIExperiment) CreatePod(clientset *kubernetes.Clientset) (*apiv1.Pod, error) {
	podsClient := clientset.CoreV1().Pods("nni-resource")
	newPod := &apiv1.Pod{
		TypeMeta: v1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "manager-" + e.User + "-" + e.WorkSpace + "-" + e.ExpID,
			Namespace: "nni-resource",
			Labels: map[string]string{
				"user":      e.User,
				"workspace": e.WorkSpace,
			},
		},
		Spec: apiv1.PodSpec{
			Volumes: []apiv1.Volume{
				{
					Name: "nni-storage",
					VolumeSource: apiv1.VolumeSource{
						NFS: &apiv1.NFSVolumeSource{
							Server:   NfsServer,
							Path:     NfsPath,
							ReadOnly: false,
						},
					},
				},
			},
			Containers: []apiv1.Container{
				{
					Ports:   []apiv1.ContainerPort{{ContainerPort: 8000}},
					Name:    "nni-manager",
					Image:   "registry.cn-hangzhou.aliyuncs.com/cuizihan/nni-manager:test",
					Command: []string{"python", "-u", "entry.py"},
					Env: []apiv1.EnvVar{
						{
							Name:  "USER",
							Value: e.User,
						},
						{
							Name:  "WORKSPACE",
							Value: e.WorkSpace,
						},
						{
							Name:  "TRAINER",
							Value: e.Trainer,
						},
						{
							Name:  "GPU_NUM",
							Value: fmt.Sprintf("%d", e.GPU),
						},
						{
							Name:  "COMMAND",
							Value: e.CMD,
						},
						{
							Name:  "TARGET",
							Value: e.Target,
						},
						{
							Name:  "USER",
							Value: e.User,
						},
						{
							Name:  "SEARCH_SPACE",
							Value: e.GetSearchSpaceJson(),
						},
						{
							Name:  "TRIAL_CONCURRENCY",
							Value: fmt.Sprintf("%d", e.Concurrency),
						},
						{
							Name:  "MAX_TRIAL_NUM",
							Value: fmt.Sprintf("%d", e.Num),
						},
						{
							Name:  "PYTHONUNBUFFERED",
							Value: "0",
						},
						{
							Name:  "EXP_ID",
							Value: e.ExpID,
						},
						{
							Name:  "NFS_SERVER",
							Value: NfsServer,
						},
						{
							Name:  "NFS_PATH",
							Value: NfsPath,
						},
						{
							Name:  "IMAGE",
							Value: IMAGE,
						},
					},
					VolumeMounts: []apiv1.VolumeMount{
						{
							Name:      "nni-storage",
							ReadOnly:  false,
							MountPath: "/nfs/nni-log-dir/" + e.ExpID + "/trials-nfs-tmp",
						},
					},
				},
			},
			RestartPolicy: "OnFailure",
		},
	}

	return podsClient.Create(newPod)
}
