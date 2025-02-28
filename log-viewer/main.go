package main

import (
	"gin-ext/common"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/karmada-io/dashboard/pkg/client"
	"github.com/karmada-io/dashboard/pkg/environment"
	"k8s-dashboard/container"
	"k8s-dashboard/logs"
	"net/http"
	"os"
	"strconv"
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
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/log/:namespace/:pod", func(c *gin.Context) {
		namespace := c.Param("namespace")
		podID := c.Param("pod")
		containerID := c.Param("container")
		refTimestamp := c.Query("referenceTimestamp")
		if refTimestamp == "" {
			refTimestamp = logs.NewestTimestamp
		}
		refLineNum, err := strconv.Atoi(c.Query("referenceLineNum"))
		if err != nil {
			refLineNum = 0
		}
		usePreviousLogs := c.Query("previous") == "true"
		offsetFrom, err1 := strconv.Atoi(c.Query("offsetFrom"))
		offsetTo, err2 := strconv.Atoi(c.Query("offsetTo"))
		logFilePosition := c.Query("logFilePosition")
		logSelector := logs.DefaultSelection
		if err1 == nil && err2 == nil {
			logSelector = &logs.Selection{
				ReferencePoint: logs.LogLineId{
					LogTimestamp: logs.LogTimestamp(refTimestamp),
					LineNum:      refLineNum,
				},
				OffsetFrom:      offsetFrom,
				OffsetTo:        offsetTo,
				LogFilePosition: logFilePosition,
			}
		}
		result, err := container.GetLogDetails(client.InClusterClient(), namespace, podID, containerID, logSelector, usePreviousLogs)
		if err != nil {
			common.Fail(c, err)
			return
		}
		common.Success(c, result)
	})
	err := r.Run(serverAddress)
	if err != nil {
		panic(err)
	}
}
