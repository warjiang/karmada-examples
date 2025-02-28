package main

import (
	"context"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/karmada-io/dashboard/pkg/client"
	"github.com/karmada-io/dashboard/pkg/environment"
	"github.com/karmada-io/karmada/pkg/util/names"
	"github.com/spf13/pflag"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

func parseSkippedNamespace(input []string) {
	// --skipped-propagating-namespaces=a
	// --skipped-propagating-namespaces a
	// f := NewFlagSet("test", ContinueOnError)
	var skippedNamespaces []string
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	f.StringSliceVar(&skippedNamespaces, "skipped-propagating-namespaces", []string{"kube-.*"},
		"Comma-separated namespaces that should be skipped from propagating.\n"+
			"Note: 'karmada-system', 'karmada-cluster' and 'karmada-es-.*' are Karmada reserved namespaces that will always be skipped.")
	f.ParseErrorsWhitelist.UnknownFlags = true
	err := f.Parse(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(skippedNamespaces)
}
func main() {
	kubeconfigPath := os.Getenv("KUBECONFIG_PATH")
	kubeContext := os.Getenv("KUBE_CONTEXT")
	skipTLSVerify := os.Getenv("SKIP_TLS_VERIFY") == "true"

	client.InitKubeConfig(
		client.WithUserAgent(environment.UserAgent()),
		client.WithKubeconfig(kubeconfigPath),
		client.WithKubeContext(kubeContext),
		client.WithInsecureTLSSkipVerify(skipTLSVerify),
	)

	kubeClient := client.InClusterClient()
	getResult, err := kubeClient.AppsV1().
		Deployments("karmada-system").
		Get(context.TODO(), "karmada-kube-controller-manager", metaV1.GetOptions{})
	if err != nil {
		panic(err)
	}
	commands := getResult.Spec.Template.Spec.Containers[0].Command
	//skipNsParams := ""
	fmt.Println(commands)
	parseSkippedNamespace(commands)
	fmt.Println(names.IsReservedNamespace("karmada-system"))
	//parseSkippedNamespace([]string{
	//	"--skipped-propagating-namespaces", "a,b",
	//})
	//for _, command := range commands {
	//	if strings.Index(command, "--skipped-propagating-namespaces") == 0 {
	//		skipNsParams = ""
	//	}
	//}
	//fmt.Println(skipNsParams)
	//cmd := app.NewControllerManagerCommand(context.TODO())
	//cmd.Execute()
}
