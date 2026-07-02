package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mchenziyi/aisc/orchestration"
	"github.com/mchenziyi/aisc/state"
)

func main() {
	// 先检查是否已冻结，不依赖 API key
	stage, err := state.LoadStage(".", "stage-requirement")
	if err == nil && stage.Status == "frozen" {
		fmt.Printf("✅ 项目已完成冻结，PRD: ./docs/prd-frozen.md\n")
		return
	}

	key := os.Getenv("DEEPSEEK_API_KEY")
	if key == "" {
		fmt.Println("请设置 DEEPSEEK_API_KEY 环境变量")
		os.Exit(1)
	}

	model := os.Getenv("DEEPSEEK_MODEL")
	if model == "" {
		model = "deepseek-v4-flash"
	}

	orch := orchestration.New(key, model)
	runner := orchestration.NewStageRunner(".", orch)

	ctx := context.Background()
	if err := runner.Run(ctx); err != nil {
		fmt.Println("❌", err)
		os.Exit(1)
	}
}
