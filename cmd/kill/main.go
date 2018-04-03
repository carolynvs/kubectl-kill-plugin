package main

import (
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	c       *kubernetes.Clientset
	ns      string
	podName string
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("kubectl plugin kill POD_NAME")
		os.Exit(1)
	}

	podName = os.Args[1]

	kill()
}

func kill() {
	loadClient()
	removeFinalizers(podName)
	deletePod(podName)
}

func getPod(podName string) *corev1.Pod {
	pod, err := c.CoreV1().Pods(ns).Get(podName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	return pod
}

func removeFinalizers(podName string) {
	pod := getPod(podName)
	pod.SetFinalizers([]string{})

	_, err := c.CoreV1().Pods(ns).Update(pod)
	if err != nil {
		panic(err)
	}

	fmt.Printf("removed finalizer from pod %s\n", podName)

}

func deletePod(podName string) {
	var noMercy int64

	opts := &metav1.DeleteOptions{GracePeriodSeconds: &noMercy}
	err := c.CoreV1().Pods(ns).Delete(podName, opts)
	if err != nil {
		panic(err)
	}

	fmt.Printf("deleted pod %s\n", podName)
}

func loadClient() {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	overrides := &clientcmd.ConfigOverrides{}

	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)

	restConfig, err := config.ClientConfig()
	if err != nil {
		panic(err)
	}

	c = kubernetes.NewForConfigOrDie(restConfig)
	ns, _, _ = config.Namespace()
}
