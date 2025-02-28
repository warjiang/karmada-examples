package main

import (
	"ginext/common"
	"github.com/gin-gonic/gin"
	"github.com/karmada-io/dashboard/pkg/client"
	"k8s.io/client-go/tools/remotecommand"
)

type TerminalResponse struct {
	ID string `json:"id"`
}

func handleExecShell(c *gin.Context) {
	sessionID, err := genTerminalSessionId()
	if err != nil {
		common.Fail(c, err)
		return
	}
	cfg, _, err := client.GetKubeConfig()
	if err != nil {
		common.Fail(c, err)
		return
	}
	terminalSessions.Set(sessionID, TerminalSession{
		id:       sessionID,
		bound:    make(chan error),
		sizeChan: make(chan remotecommand.TerminalSize),
	})
	go WaitForTerminal(client.InClusterClient(), cfg, c, sessionID)
	common.Success(c, TerminalResponse{ID: sessionID})
}
