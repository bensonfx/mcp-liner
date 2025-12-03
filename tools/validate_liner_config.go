package tools

import (
	"encoding/json"
	"fmt"

	"github.com/bensonfx/mcp-liner/internal/config"
	"github.com/bensonfx/mcp-liner/internal/responses"
	"github.com/bensonfx/mcp-liner/internal/validation"
	"github.com/phuslu/log"
)

// ValidateLinerConfigParams validate_liner_config工具的参数
type ValidateLinerConfigParams struct {
	ConfigContent string `json:"config_content"` // YAML配置内容
}

// ValidateLinerConfig 验证liner配置文件
func ValidateLinerConfig(arguments json.RawMessage) (string, error) {
	var params ValidateLinerConfigParams
	if err := json.Unmarshal(arguments, &params); err != nil {
		log.Error().Err(err).Msg("failed to parse parameters")
		return responses.ErrorResponse(
			fmt.Sprintf("Invalid parameters: %v", err),
			"Please provide 'config_content' parameter with YAML configuration",
		)
	}

	log.Info().Msg("validating liner config")

	// 验证YAML语法
	if err := validation.ValidateYAML(params.ConfigContent); err != nil {
		log.Warn().Err(err).Msg("YAML syntax validation failed")
		return responses.ErrorResponse(
			fmt.Sprintf("YAML syntax error: %v", err),
			"Please check your YAML syntax and ensure it's properly formatted",
		)
	}

	// 解析配置
	cfg, err := config.FromYAML(params.ConfigContent)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse config")
		return responses.ErrorResponse(
			fmt.Sprintf("Config parsing error: %v", err),
			"Please ensure the YAML structure matches liner configuration format",
		)
	}

	// 验证配置逻辑
	result := validation.ValidateConfig(cfg)

	if result.Valid {
		log.Info().Msg("config validation passed")
	} else {
		log.Warn().Int("errors", len(result.Errors)).Msg("config validation failed")
	}

	return responses.ValidationResponse(result)
}
