package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/beevik/etree"
	"regexp"
	"strings"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// ParseToolString 解析大模型工具调用返回的字符串
func ParseToolString(toolString string) (string, string, map[string]interface{}, error) {
	// 使用正则表达式提取 XML 部分
	re := regexp.MustCompile(`(?s)<use_mcp_tool>.*?</use_mcp_tool>`)
	matches := re.FindStringSubmatch(toolString)
	if len(matches) == 0 {
		return "", "", nil, errors.New("no tool XML found in input string")
	}
	xmlStr := matches[0]

	// 解析 XML
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xmlStr); err != nil {
		return "", "", nil, fmt.Errorf("failed to parse XML: %v", err)
	}

	root := doc.SelectElement("use_mcp_tool")
	if root == nil {
		return "", "", nil, errors.New("invalid XML structure: missing use_mcp_tool element")
	}

	// 提取 server_name
	serverNameElem := root.SelectElement("server_name")
	if serverNameElem == nil {
		return "", "", nil, errors.New("missing server_name element")
	}
	serverName := serverNameElem.Text()

	// 提取 tool_name
	toolNameElem := root.SelectElement("tool_name")
	if toolNameElem == nil {
		return "", "", nil, errors.New("missing tool_name element")
	}
	toolName := toolNameElem.Text()

	// 提取并解析 arguments
	argsElem := root.SelectElement("arguments")
	if argsElem == nil {
		return "", "", nil, errors.New("missing arguments element")
	}

	var toolArgs map[string]interface{}
	argsText := strings.TrimSpace(argsElem.Text())
	if argsText != "" {
		if err := json.Unmarshal([]byte(argsText), &toolArgs); err != nil {
			return "", "", nil, fmt.Errorf("invalid tool arguments: %v", err)
		}
	} else {
		toolArgs = make(map[string]interface{})
	}

	return serverName, toolName, toolArgs, nil
}
