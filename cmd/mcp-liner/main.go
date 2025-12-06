// mcp-liner - MCP Server for Liner configuration generation
// 用于生成和管理 liner 配置的 MCP Server
package main

import (
	"C"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bensonfx/mcp-liner/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"
)
import "io"

const (
	appName    = "mcp-liner"
	appVersion = "0.0.0"
)

var (
	cancelFunc  context.CancelFunc
	mu          sync.Mutex
	stdinReader *os.File
)

//export Stop
func Stop() {
	mu.Lock()
	defer mu.Unlock()
	if cancelFunc != nil {
		cancelFunc()
	}
	// 关闭 Pipe Reader 以打断阻塞的 Read 操作
	// 这不会关闭真实的系统 Stdin
	if stdinReader != nil {
		_ = stdinReader.Close()
	}
}

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

	// 9. generate_policy_examples - 生成 Policy 模板示例
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_policy_examples",
		Description: "生成 Policy 模板示例和文档，支持 GeoIP、Geosite、域名匹配、IP范围、文件匹配等多种路由策略",
	}, wrapToolHandler(tools.GeneratePolicyExamples))

	// 10. generate_redsocks_config - 生成 Redsocks 透明代理配置
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_redsocks_config",
		Description: "生成 Redsocks 透明代理配置（仅限Linux），支持通过 iptables 重定向 TCP 流量",
	}, wrapToolHandler(tools.GenerateRedsocksConfig))

	// 11. generate_redsocks_iptables - 生成 Redsocks iptables 规则
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_redsocks_iptables",
		Description: "生成 Redsocks 透明代理的 iptables 规则，支持 iptables-save 和 shell 脚本两种格式，包含路由循环防护",
	}, wrapToolHandler(tools.GenerateRedsocksIptables))

	// 12. generate_sni_config - 生成 SNI 路由配置
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_sni_config",
		Description: "生成 SNI（Server Name Indication）路由配置，支持基于 TLS ClientHello 的流量路由",
	}, wrapToolHandler(tools.GenerateSniConfig))

	// 13. generate_stream_config - 生成 Stream 转发配置
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_stream_config",
		Description: "生成 Stream 转发配置，支持 TCP/TLS 端口转发和 PROXY 协议",
	}, wrapToolHandler(tools.GenerateStreamConfig))

	// 14. generate_webshell_config - 生成 Web Shell 配置
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_webshell_config",
		Description: "生成 Web Shell 配置，支持通过浏览器访问终端，可配置命令和认证",
	}, wrapToolHandler(tools.GenerateWebshellConfig))

	// 创建一个可以被信号取消的 context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// 创建 Pipe 模拟 Stdin
	r, w, err := os.Pipe()
	if err != nil {
		log.Error().Err(err).Msg("failed to create pipe")
		os.Exit(1)
	}

	// 启动协程将真实 Stdin 数据复制到 Pipe Writer
	go func() {
		defer func() {
			_ = w.Close() // 显式忽略 Close 的 error
		}()
		_, _ = io.Copy(w, os.Stdin)
	}()

	// 替换 Stdin 为 Pipe Reader
	os.Stdin = r

	mu.Lock()
	cancelFunc = cancel
	stdinReader = r
	mu.Unlock()

	// 运行服务器
	log.Info().Msg("mcp server is running, waiting for connections...")
	if err := server.Run(ctx, mcp.NewStdioTransport()); err != nil {
		if err != context.Canceled {
			log.Error().Err(err).Msg("server error")
			os.Exit(1)
		}
	}
	log.Info().Msg("server stopped gracefully")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

//export RunShared
func RunShared() {
	main()
}
