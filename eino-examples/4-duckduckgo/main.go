package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/ddgsearch"
	"log"
	"time"
)

func main() {
	// the duckduckgo sdk cannot work normally
	// https://github.com/deedy5/duckduckgo_search/
	// 创建带配置的新客户端
	cfg := &ddgsearch.Config{
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		Cache:      true,
	}
	client, err := ddgsearch.New(cfg)
	if err != nil {
		log.Fatalf("New of ddgsearch failed, err=%v", err)
	}

	// 配置搜索参数
	params := &ddgsearch.SearchParams{
		Query:      "what is golang",
		Region:     ddgsearch.RegionUS,
		SafeSearch: ddgsearch.SafeSearchModerate,
		TimeRange:  ddgsearch.TimeRangeMonth,
		MaxResults: 10,
	}

	// 执行搜索
	response, err := client.Search(context.Background(), params)
	if err != nil {
		log.Fatalf("Search of ddgsearch failed, err=%v", err)
	}

	// 打印结果
	for i, result := range response.Results {
		fmt.Printf("%d. %s\n   URL: %s\n   Description: %s\n\n",
			i+1, result.Title, result.URL, result.Description)
	}
}
