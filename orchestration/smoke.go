package orchestration

import (
	"context"
	"fmt"
	"time"

	"github.com/mchenziyi/aisc/logger"
)

// runSmokeLoop 执行冒烟测试，失败时调用 agent 修复。
// 返回修复后的 artifact 内容。
func (sr *StageRunner) runSmokeLoop(
	ctx context.Context,
	artifact string,
	roundNum int,
) (string, error) {
	cfg := sr.cfg
	if cfg.SmokeTester == nil {
		return artifact, nil // 无冒烟测试，跳过
	}

	maxRetries := cfg.MaxSmokeRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}

	for retry := 0; retry <= maxRetries; retry++ {
		t0 := time.Now()
		err := cfg.SmokeTester(sr.Root)
		dur := time.Since(t0).Milliseconds()

		if err == nil {
			sr.log.Log(logger.INFO, "smoke_pass", dur, logger.F{"retries": retry})
			return artifact, nil
		}

		if retry >= maxRetries {
			sr.log.Log(logger.ERROR, "smoke_fail_max_retries", dur, logger.F{
				"retries": retry,
				"error":   err.Error(),
			})
			return artifact, fmt.Errorf("冒烟测试 %d 次后仍失败: %w", retry, err)
		}

		sr.log.Log(logger.INFO, "smoke_fail_retry", dur, logger.F{
			"retry": retry + 1,
			"max":   maxRetries,
			"error": err.Error(),
		})
		fmt.Printf("🔧 冒烟测试失败（第 %d/%d 次），Agent 自动修复...\n", retry+1, maxRetries)

		// 让 Agent 根据错误信息修复代码
		getFixPrompt := cfg.PromptRevise
		if cfg.PromptSmokeFix != nil {
			getFixPrompt = cfg.PromptSmokeFix
		}
		fixPrompt, promptErr := getFixPrompt()
		if promptErr != nil {
			return artifact, fmt.Errorf("load fix prompt: %w", promptErr)
		}
		fixTask := fmt.Sprintf("冒烟测试失败，请根据以下错误修复代码：\n\n%s\n\n当前 artifact:\n%s", err.Error(), artifact)
		t1 := time.Now()
		if len(cfg.Tools) > 0 {
			if tc, ok := sr.Orch.Client.(AgentClientWithTools); ok {
				artifact, err = tc.RunWithTools(ctx, fixPrompt, fixTask, cfg.Tools)
			}
		} else {
			artifact, err = sr.Orch.Client.Run(ctx, fixPrompt, fixTask)
		}
		sr.log.Log(logger.INFO, "smoke_fix", time.Since(t1).Milliseconds(), nil)
		if err != nil {
			return artifact, fmt.Errorf("smoke fix: %w", err)
		}
	}

	return artifact, nil
}
