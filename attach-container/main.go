package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/karmada-io/dashboard/pkg/client"
	"github.com/karmada-io/dashboard/pkg/environment"
	"net/http"
	"os"
)

func main() {
	kubeconfigPath := os.Getenv("KUBECONFIG_PATH")
	kubeContext := os.Getenv("KUBE_CONTEXT")
	skipTLSVerify := os.Getenv("SKIP_TLS_VERIFY") == "true"
	serverAddress := os.Getenv("SERVER_ADDRESS")

	client.InitKubeConfig(
		client.WithUserAgent(environment.UserAgent()),
		client.WithKubeconfig(kubeconfigPath),
		client.WithKubeContext(kubeContext),
		client.WithInsecureTLSSkipVerify(skipTLSVerify),
	)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pang")
	})

	apiV1 := r.Group("/api/v1/misc/sockjs")
	apiV1.GET("/pod/:namespace/:pod/shell/:container", handleExecShell)

	http.Handle("/", r)
	http.Handle("/api/sockjs/", CreateAttachHandler("/api/sockjs"))

	err := http.ListenAndServe(serverAddress, nil)
	if err != nil {
		fmt.Printf("Server failed to start. Error: %v\n", err)
	}
}
