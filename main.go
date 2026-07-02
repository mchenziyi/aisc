package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mchenziyi/aisc/orchestration"
	"github.com/mchenziyi/aisc/state"
)

func main() {
	// 先检查所有 Stage 是否已全部冻结（不需要 API key）
	allFrozen := true
	for _, id := range []string{"stage-requirement", "stage-api-design", "stage-tech-design"} {
		frozen, err := isStageFrozen(id)
		if err != nil || !frozen {
			allFrozen = false
			break
		}
	}
	if allFrozen {
		fmt.Println("✅ 所有 Stage 已完成冻结，无需运行")
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

	// ─── Stage 1: Requirement ──────────────────────────────
	reqCfg := orchestration.DefaultRequirementConfig()
	if frozen, _ := isStageFrozen("stage-requirement"); frozen {
		fmt.Println("✅ Requirement Stage 已冻结，跳过")
	} else {
		if err := runner.Run(ctx, reqCfg); err != nil {
			fmt.Println("❌ Requirement Stage:", err)
			os.Exit(1)
		}
		fmt.Println("✅ Requirement Stage 完成")
		fmt.Println()
	}

	// ─── Stage 2: API Design ──────────────────────────────
	apiCfg := orchestration.DefaultAPIDesignConfig()
	if frozen, _ := isStageFrozen("stage-api-design"); frozen {
		fmt.Println("✅ API Design Stage 已冻结，跳过")
	} else {
		if err := runner.Run(ctx, apiCfg); err != nil {
			fmt.Println("❌ API Design Stage:", err)
			os.Exit(1)
		}
		fmt.Println("✅ API Design Stage 完成")
		fmt.Println()
	}

	// ─── Stage 3: Tech Design ─────────────────────────────
	techCfg := orchestration.DefaultTechDesignConfig()
	if frozen, _ := isStageFrozen("stage-tech-design"); frozen {
		fmt.Println("✅ Tech Design Stage 已冻结，跳过")
	} else {
		if err := runner.Run(ctx, techCfg); err != nil {
			fmt.Println("❌ Tech Design Stage:", err)
			os.Exit(1)
		}
		fmt.Println("✅ Tech Design Stage 完成")
		fmt.Println()
	}
}

func isStageFrozen(stageID string) (bool, error) {
	s, err := state.LoadStage(".", stageID)
	if err != nil {
		return false, err
	}
	return s.Status == "frozen", nil
}
