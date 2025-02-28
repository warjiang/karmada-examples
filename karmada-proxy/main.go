package main

import (
	"context"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

func main() {
	karmadaApiServer := os.Getenv("KARMADA_API_SERVER")
	kubeconfigPath := os.Getenv("KUBECONFIG_PATH")
	memberName := os.Getenv("MEMBER_NAME")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Fatal(err)
	}
	config.Host = fmt.Sprintf("%s/apis/cluster.karmada.io/v1alpha1/clusters/%s/proxy/", karmadaApiServer, memberName)
	// 使用 rest.Config 创建一个新的 Kubernetes 客户端
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	versionInfo, err := k8sClient.Discovery().ServerVersion()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("karmada proxy version: %s\n", versionInfo.String())

	list, err := k8sClient.CoreV1().Namespaces().List(context.TODO(), metaV1.ListOptions{})
	for _, n := range list.Items {
		fmt.Printf("NAMESPACE: %s\n", n.Name)
	}
}
