package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()

	// 创建 ChatModel，默认使用 OpenAI 兼容接口
	// 可通过环境变量配置：OPENAI_API_KEY、OPENAI_MODEL、OPENAI_BASE_URL
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  os.Getenv("OPENAI_API_KEY"),
		Model:   os.Getenv("OPENAI_MODEL"),
		BaseURL: os.Getenv("OPENAI_BASE_URL"),
	})
	if err != nil {
		log.Fatalf("创建 ChatModel 失败: %v", err)
	}

	// 创建 ChatModelAgent
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "ChatBot",
		Description: "一个简单的聊天机器人",
		Instruction: "你是一个乐于助人的助手，请用中文回答用户的问题。",
		Model:       chatModel,
	})
	if err != nil {
		log.Fatalf("创建 Agent 失败: %v", err)
	}

	// 创建 Runner
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: false,
	})

	fmt.Println("=== Eino 聊天机器人 ===")
	fmt.Println("输入你的问题，按回车发送。输入 exit 退出。")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("你: ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if input == "exit" {
			fmt.Println("再见！")
			break
		}
		if input == "" {
			continue
		}

		iter := runner.Query(ctx, input)
		fmt.Print("AI: ")
		for {
			event, ok := iter.Next()
			if !ok {
				break
			}
			if event.Err != nil {
				fmt.Printf("\n错误: %v\n", event.Err)
				break
			}
			if event.Output != nil && event.Output.MessageOutput != nil {
				msg := event.Output.MessageOutput.Message
				if msg != nil && msg.Role == schema.Assistant {
					fmt.Print(msg.Content)
				}
			}
		}
		fmt.Println()
		fmt.Println()
	}
}
