# MCP-Liner

MCP-Liner 是一个 MCP (Model Context Protocol) Server，用于辅助生成和管理 [liner](https://github.com/phuslu/liner) 配置文件。通过与Claude Desktop或其他MCP客户端集成，可以快速生成各种场景下的liner配置。

## 功能特性

- 🚀 **快速生成配置** - 支持多种场景模板（HTTP转发、隧道、DNS等）
- ✅ **配置验证** - 自动检查配置语法和逻辑错误
- 📚 **内置文档** - 提供完整的liner使用文档查询
- 🔧 **灵活定制** - 支持自定义拨号器、转发策略等


## 安装与配置

> [!NOTE]
> **MCP服务器通过stdin/stdout通信**
> 可以直接运行测试（会等待JSON-RPC输入），但正常使用应通过Claude Desktop或其他MCP客户端调用。

### 1. 编译
```bash
cd /Users/benson/workspace/liner/mcp-liner
go build -o build/mcp-liner ./cmd/mcp-liner
```

### 2. 配置Claude Desktop

编辑 `claude_desktop_config.json` 文件：

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`

添加以下配置：
```json
{
  "mcpServers": {
    "mcp-liner": {
      "command": "/path/to/your/mcp-liner"
    }
  }
}
```

#### 方式2：使用 uvx (Python)

如果已安装 [uv](https://github.com/astral-sh/uv)，可以使用 `uvx` 运行（无需手动编译）：

```json
{
  "mcpServers": {
    "mcp-liner": {
      "command": "uvx",
      "args": ["mcp-liner"]
    }
  }
}
```

> 提示：如果是本地开发版本，可以使用 `"args": ["--from", "/path/to/mcp-liner", "mcp-liner"]`

重启Claude Desktop即可使用。

## MCP工具列表

### 1. generate_liner_config
生成完整的liner配置文件

**参数**:
```json
{
  "template": "http_forward|tunnel_server|tunnel_client|dns|full",
  "params": {
    "listen": [":443"],
    "server_name": ["example.org"],
    "dialer": "local"
  }
}
```

### 2. validate_liner_config
验证配置文件正确性

**参数**:
```json
{
  "config_content": "yaml配置内容"
}
```

### 3. generate_global_config
生成全局配置

**参数**:
```json
{
  "log_level": "info",
  "dns_server": "https://8.8.8.8/dns-query",
  "disable_http3": false
}
```

### 4. generate_http_config
生成HTTP/HTTPS配置

**参数**:
```json
{
  "listen": [":443"],
  "server_name": ["example.com"],
  "forward_policy": "proxy_pass",
  "dialer": "local",
  "enable_tunnel": false
}
```

### 5. generate_tunnel_config
生成隧道配置

**参数**:
```json
{
  "role": "server|client",
  "listen": [":443"],
  "server_name": ["tunnel.example.org"],
  "auth_table": "auth_user.csv"
}
```

### 6. generate_dns_config
生成DNS配置

**参数**:
```json
{
  "listen": [":53"],
  "proxy_pass": "https://8.8.8.8/dns-query"
}
```

### 7. generate_dialer_config
生成拨号器配置

**参数**:
```json
{
  "name": "cloud",
  "type": "socks5|http2|http3|ssh|wss",
  "address": "example.com:1080"
}
```

### 8. query_liner_docs
查询liner文档

**参数**:
```json
{
  "topic": "global|http|tunnel|dns|dialer|policy"
}
```

### 9. generate_policy_examples
生成Policy模板示例和文档

**参数**:
```json
{
  "config_type": "http_forward|sni_forward|socks_forward|web_doh|dns",
  "policy_type": "geoip|geosite|domain_match|ip_range|file_based|fetch_based|custom"
}
```

### 10. generate_redsocks_config
生成Redsocks透明代理配置（仅限Linux）

**参数**:
```json
{
  "listen": [":12345"],
  "dialer": "proxy",
  "dialer_url": "socks5://127.0.0.1:1080",
  "log": true
}
```

### 11. generate_redsocks_iptables
生成Redsocks iptables规则

**参数**:
```json
{
  "redsocks_port": 12345,
  "lan_interface": "eth0",
  "wan_interface": "eth1",
  "proxy_ports": [80, 443],
  "exclude_cidrs": ["10.0.0.0/8"],
  "format": "iptables-save|shell-script"
}
```

### 12. generate_sni_config
生成SNI路由配置

**参数**:
```json
{
  "enabled": true,
  "policy": "{{ if hasSuffixes \"google.com\" .ServerName }}proxy{{ else }}direct{{ end }}",
  "dialer": "local",
  "log": true
}
```

### 13. generate_stream_config
生成Stream转发配置

**参数**:
```json
{
  "listen": [":3389"],
  "proxy_pass": "192.168.1.100:3389",
  "dialer": "local",
  "proxy_protocol": 0,
  "log": true
}
```

### 14. generate_webshell_config
生成Web Shell配置

**参数**:
```json
{
  "listen": [":443"],
  "server_name": ["shell.example.org"],
  "command": "login",
  "home": "/home/user",
  "auth_table": "auth_user.csv",
  "location": "/shell/"
}
```


## 使用示例

### 示例1：生成HTTP转发配置

在Claude中输入：
```
使用 generate_http_config 工具生成一个HTTPS转发配置，监听443端口，server_name是example.org
```

### 示例2：生成隧道服务端配置

```
使用 generate_tunnel_config 工具生成隧道服务端配置：
- listen: [":443"]
- server_name: ["tunnel.example.org"]
- auth_table: "auth_user.csv"
```

### 示例3：验证配置文件

```
使用 validate_liner_config 工具验证以下配置：
[粘贴你的YAML配置]
```

### 示例4：查询文档

```
使用 query_liner_docs 查询tunnel相关的文档
```

## 开发

### 运行测试
```bash
# 运行所有测试
go test ./... -v

# 仅运行内部模块测试
go test ./internal/... -v

# 仅运行工具测试
go test ./tools/... -v
```

### 项目结构
```
mcp-liner/
├── cmd/mcp-liner/      # 主程序入口
│   └── main.go
├── internal/           # 内部模块
│   ├── config/         # 配置结构定义
│   ├── templates/      # 配置模板
│   ├── validation/     # 配置验证
│   └── responses/      # MCP响应格式化
├── tools/              # MCP工具实现
│   ├── generate_liner_config.go
│   ├── validate_liner_config.go
│   ├── generate_global_config.go
│   ├── generate_http_config.go
│   ├── generate_tunnel_config.go
│   ├── generate_dns_config.go
│   ├── generate_dialer_config.go
│   └── query_liner_docs.go
└── tests/              # 测试代码
    └── integration/
```

## 依赖

- Go 1.23+
- github.com/modelcontextprotocol/go-sdk v0.2.0
- github.com/phuslu/log v1.0.113
- github.com/spf13/cobra v1.8.1
- gopkg.in/yaml.v3 v3.0.1

## CI & Coverage

本项目配置了完整的 GitHub Actions CI 流程：
- **Lint**: GolangCI-Lint (Go) & Ruff (Python)
- **Test**: Go Test & Pytest
- **Build**: 跨平台编译验证

### Coverage 报告
测试覆盖率报告会自动上传至 [Codecov](https://about.codecov.io/)。
要查看覆盖率报告：
1. 确保您的仓库已连接到 Codecov。
2. 配置 `CODECOV_TOKEN` (如果私有仓库)。
3. 在 README 顶部添加 Badge：
   `[![codecov](https://codecov.io/gh/<ORG>/<REPO>/graph/badge.svg?token=<TOKEN>)](https://codecov.io/gh/<ORG>/<REPO>)`

## 版本

当前版本：**v1.0.0**

## License

与liner主项目保持一致

## 相关链接

- [Liner项目](https://github.com/phuslu/liner)
- [MCP协议](https://modelcontextprotocol.io/)
