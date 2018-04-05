package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/pluginutils"

	_ "github.com/golang/glog"
)

var (
	c           *kubernetes.Clientset
	ns          string
	podName     string
	gracePeriod int64
)

func init() {
	// Initialize glog flags
	flag.CommandLine.Set("logtostderr", "true")
	flag.CommandLine.Set("v", os.Getenv("KUBECTL_PLUGINS_GLOBAL_FLAG_V"))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: kubectl plugin kill POD_NAME [--grace-period]")
		os.Exit(1)
	}

	podName = os.Args[1]

	gracePeriodFlag := os.Getenv("KUBECTL_PLUGINS_LOCAL_FLAG_GRACE_PERIOD")
	if g, err := strconv.ParseInt(gracePeriodFlag, 10, 64); err == nil {
		gracePeriod = g
	}

	kill()
}

func kill() {
	loadConfig()
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

	fmt.Printf("removed finalizers from pod %s\n", podName)

}

func deletePod(podName string) {
	fmt.Printf("killing %s/%s with a grace period of %ds...\n", ns, podName, gracePeriod)

	opts := &metav1.DeleteOptions{GracePeriodSeconds: &gracePeriod}
	err := c.CoreV1().Pods(ns).Delete(podName, opts)
	if err != nil && !errors.IsNotFound(err) {
		panic(err)
	}

	fmt.Printf("deleted pod %s\n", podName)
}

func loadConfig() {
	restConfig, kubeConfig, err := pluginutils.InitClientAndConfig()
	if err != nil {
		panic(err)
	}
	c = kubernetes.NewForConfigOrDie(restConfig)
	ns, _, _ = kubeConfig.Namespace()
}
