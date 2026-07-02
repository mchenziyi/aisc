package main

import (
	"context"
	"fmt"
	"os"

	"agentdemo/agent"

	"github.com/mchenziyi/aisc/agents/prompts"
)

func main() {
	key := os.Getenv("DEEPSEEK_API_KEY")
	if key == "" {
		fmt.Println("请设置 DEEPSEEK_API_KEY 环境变量")
		os.Exit(1)
	}

	a, err := agent.New(key, "deepseek-v4-flash", false)
	if err != nil {
		fmt.Println("创建 Agent 失败:", err)
		os.Exit(1)
	}

	pmPrompt, err := prompts.Load("pm", "draft")
	if err != nil {
		fmt.Println("加载 PM 提示词失败:", err)
		os.Exit(1)
	}
	a.SetSystemPrompt(pmPrompt)
	a.MaxMessages = 5000

	ctx := context.Background()
	task := `请根据以下用户需求，输出完整 PRD：

做一个 AI 视频平台。
核心功能：
1. 用户可以上传视频（标题、描述、标签）
2. 用户可以浏览视频列表（支持搜索和分页）
3. 用户可以收藏视频
4. 用户可以看到自己的收藏列表
用户需要登录后才能使用。
技术栈不做限制，先出需求文档。`

	fmt.Println("🚀 正在生成 PRD...")
	result, err := a.Run(ctx, task)
	if err != nil {
		fmt.Println("Agent 运行失败:", err)
		os.Exit(1)
	}
	fmt.Println("--- PRD ---")
	fmt.Println(result)
}
