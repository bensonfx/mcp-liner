package tools

import (
	"encoding/json"
	"fmt"

	"github.com/bensonfx/mcp-liner/internal/config"
	"github.com/bensonfx/mcp-liner/internal/responses"
	"github.com/phuslu/log"
)

// GenerateStreamConfigParams generate_stream_config工具的参数
type GenerateStreamConfigParams struct {
	Listen        []string `json:"listen"`         // 监听地址，如 [":3389"]
	ProxyPass     string   `json:"proxy_pass"`     // 转发目标地址，如 "192.168.1.100:3389"
	Dialer        string   `json:"dialer"`         // 拨号器名称
	DialerURL     string   `json:"dialer_url"`     // 拨号器URL（可选）
	Keyfile       string   `json:"keyfile"`        // TLS密钥文件（可选）
	Certfile      string   `json:"certfile"`       // TLS证书文件（可选）
	ProxyProtocol uint     `json:"proxy_protocol"` // PROXY协议版本（0=禁用，1或2=启用）
	DialTimeout   int      `json:"dial_timeout"`   // 拨号超时（秒）
	SpeedLimit    int64    `json:"speed_limit"`    // 速度限制（字节/秒）
	Log           bool     `json:"log"`            // 是否启用日志
}

// GenerateStreamConfig 生成Stream转发配置
func GenerateStreamConfig(arguments json.RawMessage) (string, error) {
	var params GenerateStreamConfigParams
	if err := json.Unmarshal(arguments, &params); err != nil {
		log.Error().Err(err).Msg("failed to parse parameters")
		return responses.ErrorResponse(
			fmt.Sprintf("Invalid parameters: %v", err),
			"Please provide valid parameters for Stream configuration",
		)
	}

	log.Info().
		Strs("listen", params.Listen).
		Str("proxy_pass", params.ProxyPass).
		Str("dialer", params.Dialer).
		Bool("log", params.Log).
		Msg("generating stream config")

	// 设置默认值
	if len(params.Listen) == 0 {
		params.Listen = []string{":8080"}
	}
	if params.ProxyPass == "" {
		return responses.ErrorResponse(
			"proxy_pass is required",
			"Please specify the target address to forward to",
		)
	}
	if params.Dialer == "" {
		params.Dialer = "local"
	}
	if params.DialTimeout == 0 {
		params.DialTimeout = 5
	}

	// 构建配置
	cfg := config.Config{
		Global: config.NewDefaultGlobalConfig(),
		Dialer: map[string]string{
			"local": "local",
		},
		Stream: []config.StreamConfig{
			{
				Listen:        params.Listen,
				ProxyPass:     params.ProxyPass,
				Dialer:        params.Dialer,
				Keyfile:       params.Keyfile,
				Certfile:      params.Certfile,
				ProxyProtocol: params.ProxyProtocol,
				DialTimeout:   params.DialTimeout,
				SpeedLimit:    params.SpeedLimit,
				Log:           params.Log,
			},
		},
	}

	// 如果提供了dialer URL，添加到配置
	if params.DialerURL != "" && params.Dialer != "local" {
		cfg.Dialer[params.Dialer] = params.DialerURL
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

	description := "Generated Stream forwarding configuration\n\n"
	description += "Stream forwarding provides TCP/TLS port forwarding with these features:\n"
	description += "- Simple TCP port forwarding\n"
	description += "- Optional TLS termination\n"
	description += "- PROXY protocol support (v1 and v2)\n"
	description += "- Speed limiting\n"
	description += "- Custom dialer support\n\n"

	description += "Common use cases:\n"
	description += "- Port forwarding: Forward local port to remote service\n"
	description += "- Load balancing: Distribute traffic across backends\n"
	description += "- TLS termination: Decrypt TLS and forward plain TCP\n"
	description += "- Access control: Add authentication layer to raw TCP services\n\n"

	if params.ProxyProtocol > 0 {
		description += fmt.Sprintf("PROXY protocol v%d enabled - real client IP will be forwarded\n", params.ProxyProtocol)
	}

	if params.Keyfile != "" && params.Certfile != "" {
		description += "TLS enabled - connection will be encrypted\n"
	}

	log.Info().Msg("stream config generated successfully")
	return responses.SuccessResponse(yamlContent, description)
}
