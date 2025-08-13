package main

import (
	"context"
	"io"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// newBranch branch initialization method of node 'ChatModel4' in graph 'test'
func newBranch1(ctx context.Context, in *schema.Message) (endNode string, err error) {
	if isToolCall, err := ToolCallChecker(ctx, in); err != nil {
		return "", err
	} else if isToolCall {
		return "ToolsNode1", nil
	}
	return compose.END, nil
}

func ToolCallChecker(_ context.Context, in *schema.Message) (bool, error) {
	if in.ToolCalls == nil {
		return false, nil
	}
	if len(in.ToolCalls) > 0 {
		return true, nil
	}
	return false, nil
}

func newBranch2(ctx context.Context, in []*schema.Message) (endNode string, err error) {
	err = compose.ProcessState[*state](ctx, func(_ context.Context, state *state) error {
		if len(state.ReturnDirectlyToolCallID) > 0 {
			endNode = "Lambda1" // 如果需要直接返回，去Lambda1节点
		} else {
			endNode = "ChatModel4" // 否则返回到模型节点继续循环
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return endNode, nil
}

// 流式分支函数，用于ChatModel4节点
func newBranch1Stream(ctx context.Context, sr *schema.StreamReader[*schema.Message]) (endNode string, err error) {
	defer sr.Close()

	for {
		msg, err := sr.Recv()
		if err == io.EOF {
			return compose.END, nil
		}
		if err != nil {
			return "", err
		}

		if len(msg.ToolCalls) > 0 {
			return "ToolsNode1", nil
		}

		if len(msg.Content) == 0 { // 跳过前面的空块
			continue
		}

		return compose.END, nil
	}
}

// 流式分支函数，用于ToolsNode1节点
func newBranch2Stream(ctx context.Context, sr *schema.StreamReader[[]*schema.Message]) (endNode string, err error) {
	sr.Close()

	err = compose.ProcessState[*state](ctx, func(_ context.Context, state *state) error {
		if len(state.ReturnDirectlyToolCallID) > 0 {
			endNode = "Lambda1" // 如果需要直接返回，去Lambda1节点
		} else {
			endNode = "ChatModel4" // 否则返回到模型节点继续循环
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return endNode, nil
}
