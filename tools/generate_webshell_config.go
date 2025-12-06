package tools

import (
	"encoding/json"
	"fmt"

	"github.com/bensonfx/mcp-liner/internal/config"
	"github.com/bensonfx/mcp-liner/internal/responses"
	"github.com/phuslu/log"
)

// GenerateWebshellConfigParams generate_webshell_config工具的参数
type GenerateWebshellConfigParams struct {
	Listen     []string `json:"listen"`      // 监听地址
	ServerName []string `json:"server_name"` // 域名
	Command    string   `json:"command"`     // 执行命令
	Home       string   `json:"home"`        // Home目录
	AuthTable  string   `json:"auth_table"`  // 认证表
	Location   string   `json:"location"`    // URL路径
}

// GenerateWebshellConfig 生成Web Shell配置
func GenerateWebshellConfig(arguments json.RawMessage) (string, error) {
	var params GenerateWebshellConfigParams
	if err := json.Unmarshal(arguments, &params); err != nil {
		log.Error().Err(err).Msg("failed to parse parameters")
		return responses.ErrorResponse(
			fmt.Sprintf("Invalid parameters: %v", err),
			"Please provide valid parameters for Webshell configuration",
		)
	}

	log.Info().
		Strs("listen", params.Listen).
		Strs("server_name", params.ServerName).
		Str("command", params.Command).
		Msg("generating Webshell config")

	// 设置默认值
	if len(params.Listen) == 0 {
		params.Listen = []string{":443"}
	}
	if len(params.ServerName) == 0 {
		params.ServerName = []string{"shell.example.org"}
	}
	if params.Command == "" {
		params.Command = "login"
	}
	if params.Location == "" {
		params.Location = "/shell/"
	}

	// 使用模板生成配置
	// 注意：这里我们需要构建一个包含 Web Shell 的 Config
	// 由于 templates 包可能还没有专门针对 Shell 的模板，我们可能需要更新 templates
	// 或者直接在这里构建 Config 对象

	webConfig := config.HTTPWebConfig{
		Location: params.Location,
		Shell: config.HTTPWebShellConfig{
			Enabled:   true,
			Command:   params.Command,
			Home:      params.Home,
			AuthTable: params.AuthTable,
		},
	}
	// 如果 auth_table 没设置，给个默认值? 还是保持为空由用户决定
	if params.AuthTable == "" {
		webConfig.Shell.AuthTable = "auth_user.csv"
	}

	httpConfig := config.HTTPConfig{
		Listen:     params.Listen,
		ServerName: params.ServerName,
		Forward: config.HTTPForwardConfig{
			Policy: "return 404", // Shell 应该是主要功能，Forward 这里可以给个默认
			Log:    true,
		},
		Web: []config.HTTPWebConfig{webConfig},
	}

	cfg := config.Config{
		Global: config.NewDefaultGlobalConfig(),
		Dialer: map[string]string{
			"local": "local",
		},
		Https: []config.HTTPConfig{httpConfig},
	}

	// 转换为YAML
	yamlContent, err := cfg.ToYAML()
	if err != nil {
		log.Error().Err(err).Msg("failed to convert config to YAML")
		return responses.ErrorResponse(
			fmt.Sprintf("Failed to generate YAML: %v", err),
			"",
		)
	}

	description := "Generated Webshell configuration\n"
	description += fmt.Sprintf("Access via https://%s%s", params.ServerName[0], params.Location)

	log.Info().Msg("Webshell config generated successfully")
	return responses.SuccessResponse(yamlContent, description)
}
