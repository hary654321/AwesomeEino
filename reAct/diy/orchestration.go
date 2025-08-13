package main

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/schema"
)

type state struct {
	Messages                 []*schema.Message
	ReturnDirectlyToolCallID string
}

func Buildtest(ctx context.Context) (r compose.Runnable[[]*schema.Message, *schema.Message], err error) {
	const (
		ChatModel4 = "ChatModel4"
		ToolsNode1 = "ToolsNode1"
		Lambda1    = "Lambda1"
	)
	//创建图
	g := compose.NewGraph[[]*schema.Message, *schema.Message](compose.WithGenLocalState(func(ctx context.Context) *state {
		return &state{Messages: make([]*schema.Message, 0, 12)}
	}))

	// 首先创建工具节点以获取工具信息
	toolsNode1KeyOfToolsNode, err := newToolsNode(ctx)
	if err != nil {
		return nil, err
	}

	// 生成工具信息
	toolInfos, err := genToolInfos(ctx, toolsNode1KeyOfToolsNode)
	if err != nil {
		return nil, err
	}

	// 创建原始模型
	rawModel, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}

	// 使用agent.ChatModelWithTools包装模型，让模型知道可用的工具
	chatModel4KeyOfChatModel, err := agent.ChatModelWithTools(nil, rawModel, toolInfos)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatModelNode(ChatModel4, chatModel4KeyOfChatModel, compose.WithNodeName("ArkModel"),
		compose.WithStatePreHandler(func(ctx context.Context, in []*schema.Message, state *state) ([]*schema.Message, error) {
			state.Messages = append(state.Messages, in...)
			// 复制一份
			copyMessage := make([]*schema.Message, len(state.Messages))
			copy(copyMessage, state.Messages)
			return copyMessage, nil
		}))

	//添加工具节点
	_ = g.AddToolsNode(ToolsNode1, toolsNode1KeyOfToolsNode,
		compose.WithStatePreHandler(func(ctx context.Context, in *schema.Message, state *state) (*schema.Message, error) {
			if in == nil {
				return state.Messages[len(state.Messages)-1], nil
			}
			state.Messages = append(state.Messages, in)
			// 检查是否有需要直接返回的工具调用（这里暂时不设置任何工具直接返回）
			state.ReturnDirectlyToolCallID = ""
			return in, nil
		}))
	//添加lambda节点
	lambda1 := compose.TransformableLambda(func(ctx context.Context, msgs *schema.StreamReader[[]*schema.Message]) (*schema.StreamReader[*schema.Message], error) {
		return schema.StreamReaderWithConvert(msgs, func(msgs []*schema.Message) (*schema.Message, error) {
			var msg *schema.Message
			err := compose.ProcessState[*state](ctx, func(_ context.Context, state *state) error {
				for i := range msgs {
					if msgs[i] != nil && msgs[i].ToolCallID == state.ReturnDirectlyToolCallID {
						msg = msgs[i]
						return nil
					}
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
			if msg == nil {
				return nil, schema.ErrNoValue
			}
			return msg, nil
		}), nil
	})
	_ = g.AddLambdaNode(Lambda1, lambda1)
	//链接节点
	_ = g.AddEdge(compose.START, ChatModel4)
	_ = g.AddEdge(Lambda1, compose.END)
	_ = g.AddBranch(ChatModel4, compose.NewStreamGraphBranch(newBranch1Stream, map[string]bool{compose.END: true, ToolsNode1: true}))
	_ = g.AddBranch(ToolsNode1, compose.NewStreamGraphBranch(newBranch2Stream, map[string]bool{ChatModel4: true, Lambda1: true}))

	r, err = g.Compile(ctx, compose.WithGraphName("test"), compose.WithNodeTriggerMode(compose.AnyPredecessor), compose.WithMaxRunSteps(12))
	if err != nil {
		return nil, err
	}
	return r, err
}

// genToolInfos 生成工具信息，参考react.go的实现
func genToolInfos(ctx context.Context, toolsNode *compose.ToolsNode) ([]*schema.ToolInfo, error) {
	// 创建工具实例
	addTool := GetAddTool()
	subTool := GetSubTool()
	analyzeTool := GetAnalyzeTool()

	tools := []tool.BaseTool{addTool, subTool, analyzeTool}

	toolInfos := make([]*schema.ToolInfo, 0, len(tools))
	for _, t := range tools {
		tl, err := t.Info(ctx)
		if err != nil {
			return nil, err
		}
		toolInfos = append(toolInfos, tl)
	}

	return toolInfos, nil
}
