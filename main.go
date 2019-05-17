package main

import (
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"strings"
)

type podStructItem struct {
	podName string
	podStatus v1.PodPhase
}

type podStruct struct {
	items []podStructItem
}

func (pod *podStruct) AddItem(item podStructItem) {
	pod.items = append(pod.items, item)
}

func connect(hostname string)(clientset *kubernetes.Clientset, namespace *string){

	// Load Config
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	namespace = flag.String("namespace", "kube-system", "Namespace identifier")
	flag.Parse()

	// Configure Connection
	var configInfo *rest.Config
	if hostname != "" {
		configInfo, _ = clientcmd.BuildConfigFromFlags(hostname, *kubeconfig)
	} else {
	    configInfo, _ = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	}

	// Return Connection Object
	clientset, _ = kubernetes.NewForConfig(configInfo)

	return clientset, namespace
}

func execute(clientset *kubernetes.Clientset, namespace *string) int {
	// List Pods
	podClient := clientset.CoreV1().Pods(*namespace)
	fmt.Printf("Listing pods in namespace kube-system\n\n")
	list, err := podClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	pod := podStruct{}

	// Populate Pod struct
	for _, d := range list.Items {
		p := podStructItem{podName: d.Name, podStatus: d.Status.Phase}
		pod.AddItem(p)
	}
	// Operate based on criteria
	for _, podinfo := range pod.items {
		if strings.HasPrefix(podinfo.podName,"core") {
			fmt.Println(podinfo.podName)
		}
	}
	return 0
}

func main() {

	clientset, namespace := connect("")
	execute(clientset, namespace)
}


