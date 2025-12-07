package tools

import (
	"encoding/json"
	"fmt"

	"github.com/bensonfx/mcp-liner/internal/responses"
	"github.com/phuslu/log"
)

// GenerateDialerConfigParams generate_dialer_config工具的参数
type GenerateDialerConfigParams struct {
	Name    string `json:"name"`    // 拨号器名称，如 "cloud"
	Type    string `json:"type"`    // 类型：local, socks5, http2, http3, ssh, wss
	Address string `json:"address"` // 地址，如 "example.com:1080"
}

// GenerateDialerConfig 生成代理拨号器配置
func GenerateDialerConfig(arguments json.RawMessage) (string, error) {
	var params GenerateDialerConfigParams
	if err := json.Unmarshal(arguments, &params); err != nil {
		log.Error().Err(err).Msg("failed to parse parameters")
		return responses.ErrorResponse(
			fmt.Sprintf("Invalid parameters: %v", err),
			"Please provide 'name', 'type', and 'address' parameters",
		)
	}

	log.Info().
		Str("name", params.Name).
		Str("type", params.Type).
		Str("address", params.Address).
		Msg("generating dialer config")

	// 验证参数
	if params.Name == "" {
		return responses.ErrorResponse(
			"Dialer name is required",
			"Please provide a name for the dialer, e.g., 'cloud', 'proxy'",
		)
	}

	validTypes := []string{"local", "socks5", "http2", "http3", "ssh", "wss"}
	isValidType := false
	for _, t := range validTypes {
		if params.Type == t {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return responses.ErrorResponse(
			fmt.Sprintf("Invalid dialer type: %s", params.Type),
			fmt.Sprintf("Supported types: %v", validTypes),
		)
	}

	// 构建拨号器URL
	var dialerURL string
	switch params.Type {
	case "local":
		// 本地 dialer 格式必须为 local:// 或 local://interface（如 local://wg0）
		if params.Address != "" {
			dialerURL = fmt.Sprintf("local://%s", params.Address)
		} else {
			dialerURL = "local://"
		}
	case "socks5":
		dialerURL = fmt.Sprintf("socks5://%s", params.Address)
	case "http2":
		dialerURL = fmt.Sprintf("http2://%s", params.Address)
	case "http3":
		dialerURL = fmt.Sprintf("http3://%s", params.Address)
	case "ssh":
		dialerURL = fmt.Sprintf("ssh://%s", params.Address)
	case "wss":
		dialerURL = fmt.Sprintf("wss://%s", params.Address)
	}

	// 生成YAML配置片段
	yamlContent := fmt.Sprintf(`dialer:
  %s: %s
`, params.Name, dialerURL)

	example := fmt.Sprintf(`Usage example:
1. Add this dialer configuration to your liner config
2. Reference it in http/tunnel configs using: dialer: %s

Example HTTP config:
https:
  - listen: [":443"]
    server_name: ["example.org"]
    forward:
      policy: proxy_pass
      dialer: %s
      log: true
`, params.Name, params.Name)

	log.Info().Str("dialer_name", params.Name).Msg("dialer config generated successfully")
	return responses.ConfigWithExampleResponse(yamlContent, example)
}
