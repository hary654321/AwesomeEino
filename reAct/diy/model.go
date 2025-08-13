package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/model"
)

// newChatModel component initialization function of node 'ChatModel4' in graph 'test'
func newChatModel(ctx context.Context) (cm model.ToolCallingChatModel, err error) {
	arkModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "qwen3:4b",
	})
	if err != nil {
		fmt.Printf("failed to create chat model: %v", err)
		return
	}
	return arkModel, nil
}
