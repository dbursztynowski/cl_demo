package tools

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Podpatch(namespace string, podName string, cpu string, ctx context.Context, client client.Client, l VerbosityLog, req ctrl.Request) {

	secret := &v1.Secret{}
	err := client.Get(ctx, types.NamespacedName{Name: "kubeconfig", Namespace: req.Namespace}, secret)
	if err != nil {
		l.V(1).Info("Error getting secret: " + err.Error())
	}

	cfg, err := clientcmd.NewClientConfigFromBytes(secret.Data["kubeconfig"])
	if err != nil {
		l.V(1).Info("Error building kubeconfig: " + err.Error())
	}
	restClientConfig, err := cfg.ClientConfig()
	if err != nil {
		l.V(1).Info("Error creating restClientConfig: " + err.Error())
	}

	clientset, err := kubernetes.NewForConfig(restClientConfig)
	if err != nil {
		l.V(1).Info("Error creating clientset: " + err.Error())
	}


	patchBytes := []byte(`{
	"spec":
		{
			"containers":
			[
				{
					"name":"open5gs-upf", 
					"resources":
						{
							"limits":
								{
									"cpu": "` + cpu + `"
								}
						}
				}
			]
		}
	}`)

	pod := &v1.Pod{}
	pod, err = clientset.CoreV1().Pods(namespace).Patch(context.TODO(), podName, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
	if err != nil {
		l.V(1).Info("Error patching pod: " + err.Error())
	} else {
		l.V(1).Info("Patch applied for pod: " + podName)
	}
	fmt.Printf("Pod spec: %v", pod.Spec)

}
