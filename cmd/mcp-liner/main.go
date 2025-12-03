// mcp-liner - MCP Server for Liner configuration generation
// 用于生成和管理 liner 配置的 MCP Server
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/bensonfx/mcp-liner/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)

const (
	appName    = "mcp-liner"
	appVersion = "1.0.0"
)

var rootCmd = &cobra.Command{
	Use:     appName,
	Short:   "MCP Server for Liner configuration generation",
	Long:    `mcp-liner 是一个 MCP Server，用于生成和管理 liner 配置文件。`,
	Version: appVersion,
	Run:     runServer,
}

func init() {
	// 设置日志
	log.DefaultLogger = log.Logger{
		Level:      log.InfoLevel,
		Caller:     1,
		TimeFormat: "15:04:05",
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			EndWithMessage: true,
		},
	}
}

// wrapToolHandler 包装工具处理函数以符合MCP SDK的类型要求
func wrapToolHandler(handler func(json.RawMessage) (string, error)) mcp.ToolHandlerFor[map[string]interface{}, interface{}] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]interface{}]) (*mcp.CallToolResultFor[interface{}], error) {
		// 将参数转换为JSON
		jsonData, err := json.Marshal(params.Arguments)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal parameters: %w", err)
		}

		// 调用原始处理函数
		result, err := handler(jsonData)
		if err != nil {
			return nil, err
		}

		// 返回结果
		return &mcp.CallToolResultFor[interface{}]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: result,
				},
			},
		}, nil
	}
}

func runServer(cmd *cobra.Command, args []string) {
	log.Info().Str("version", appVersion).Msg("starting mcp-liner server")

	// 创建 MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    appName,
		Version: appVersion,
	}, nil)

	// 注册工具
	// 使用包装函数来符合MCP SDK的类型要求

	// 1. generate_liner_config - 生成完整的 liner 配置文件
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_liner_config",
		Description: "生成完整的 liner 配置文件，支持多种场景模板（http_forward, tunnel_server, tunnel_client, dns, full）",
	}, wrapToolHandler(tools.GenerateLinerConfig))

	// 2. validate_liner_config - 验证配置文件
	mcp.AddTool(server, &mcp.Tool{
		Name:        "validate_liner_config",
		Description: "验证 liner 配置文件的正确性，检查语法和逻辑错误",
	}, wrapToolHandler(tools.ValidateLinerConfig))

	// 3. generate_global_config - 生成全局配置
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_global_config",
		Description: "生成 liner 全局配置（日志、DNS、连接池等设置）",
	}, wrapToolHandler(tools.GenerateGlobalConfig))

	// 4. generate_http_config - 生成 HTTP/HTTPS 转发配置
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_http_config",
		Description: "生成 HTTP/HTTPS 转发配置，支持 forward 策略、tunnel 和 web 服务",
	}, wrapToolHandler(tools.GenerateHTTPConfig))

	// 5. generate_tunnel_config - 生成内网穿透配置
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_tunnel_config",
		Description: "生成内网穿透配置，支持 server（服务端）和 client（客户端）两种角色",
	}, wrapToolHandler(tools.GenerateTunnelConfig))

	// 6. generate_dns_config - 生成 DNS 配置
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_dns_config",
		Description: "生成 DNS 服务配置，支持 DNS/DoT/DoH，可配置策略和上游服务器",
	}, wrapToolHandler(tools.GenerateDNSConfig))

	// 7. generate_dialer_config - 生成代理拨号器配置
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_dialer_config",
		Description: "生成代理拨号器配置，支持 local、socks5、http2、http3、ssh、wss 等类型",
	}, wrapToolHandler(tools.GenerateDialerConfig))

	// 8. query_liner_docs - 查询 liner 文档
	mcp.AddTool(server, &mcp.Tool{
		Name:        "query_liner_docs",
		Description: "查询 liner 文档和使用说明，支持按主题查询（global, http, tunnel, dns, dialer, policy）",
	}, wrapToolHandler(tools.QueryLinerDocs))

	// 运行服务器
	log.Info().Msg("mcp server is running, waiting for connections...")
	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		log.Error().Err(err).Msg("server error")
		os.Exit(1)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
