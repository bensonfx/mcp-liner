package tools

import (
	"encoding/json"
	"fmt"

	"github.com/bensonfx/mcp-liner/internal/config"
	"github.com/bensonfx/mcp-liner/internal/responses"
	"github.com/phuslu/log"
)

// GenerateSSHConfigParams generate_ssh_config工具的参数
type GenerateSSHConfigParams struct {
	Listen           []string `json:"listen"`            // 监听地址，如 [":2222"]
	ServerVersion    string   `json:"server_version"`    // SSH服务版本字符串
	HostKey          string   `json:"host_key"`          // 主机私钥内容或路径
	AuthTable        string   `json:"auth_table"`        // 认证表路径
	AuthorizedKeys   string   `json:"authorized_keys"`   // 公钥认证文件路径
	Shell            string   `json:"shell"`             // 默认Shell
	Home             string   `json:"home"`              // 用户主目录模板
	DisableKeepalive bool     `json:"disable_keepalive"` // 是否禁用Keepalive
	Log              bool     `json:"log"`               // 是否启用日志
}

// GenerateSSHConfig 生成SSH Server配置
func GenerateSSHConfig(arguments json.RawMessage) (string, error) {
	var params GenerateSSHConfigParams
	if err := json.Unmarshal(arguments, &params); err != nil {
		log.Error().Err(err).Msg("failed to parse parameters")
		return responses.ErrorResponse(
			fmt.Sprintf("Invalid parameters: %v", err),
			"Please provide valid parameters for SSH configuration",
		)
	}

	log.Info().
		Strs("listen", params.Listen).
		Msg("generating SSH config")

	// 设置默认值
	if len(params.Listen) == 0 {
		params.Listen = []string{":2222"}
	}
	if params.HostKey == "" {
		// Default to strict-mode generated key or standard path?
		// liner usually expects a file or PEM content.
		params.HostKey = "ssh_host_key"
	}
	if params.AuthTable == "" {
		params.AuthTable = "auth_user.csv"
	}
	if params.Shell == "" {
		params.Shell = "/bin/bash" // default to bash
	}

	cfg := config.Config{
		Global: config.NewDefaultGlobalConfig(),
		Ssh: []config.SshConfig{
			{
				Listen:           params.Listen,
				ServerVersion:    params.ServerVersion,
				HostKey:          params.HostKey,
				AuthTable:        params.AuthTable,
				AuthorizedKeys:   params.AuthorizedKeys,
				Shell:            params.Shell,
				Home:             params.Home,
				DisableKeepalive: params.DisableKeepalive,
				Log:              params.Log,
			},
		},
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

	description := "Generated SSH Server configuration"
	log.Info().Msg("SSH config generated successfully")
	return responses.SuccessResponse(yamlContent, description)
}
