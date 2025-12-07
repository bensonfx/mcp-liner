// Package templates 提供liner配置模板
package templates

import (
	"github.com/bensonfx/mcp-liner/internal/config"
)

// HTTPForwardTemplate 生成HTTP转发配置模板
// dialer: 拨号器名称，如 "local", "proxy"等
// enableLog: 是否启用日志
func HTTPForwardTemplate(listen []string, serverName []string, dialer string, enableLog bool) config.HTTPConfig {
	return config.HTTPConfig{
		Listen:     listen,
		ServerName: serverName,
		Forward: config.HTTPForwardConfig{
			Policy: "proxy_pass",
			Dialer: dialer,
			Log:    enableLog,
		},
	}
}

// HTTPSForwardTemplate 生成HTTPS转发配置模板（带TLS证书）
func HTTPSForwardTemplate(listen []string, serverName []string, certfile, keyfile, dialer string, enableLog bool) config.HTTPConfig {
	return config.HTTPConfig{
		Listen:     listen,
		ServerName: serverName,
		Certfile:   certfile,
		Keyfile:    keyfile,
		Forward: config.HTTPForwardConfig{
			Policy: "proxy_pass",
			Dialer: dialer,
			Log:    enableLog,
		},
	}
}

// HTTPForwardWithPolicyTemplate 生成带自定义policy的HTTP转发配置
// policy: 转发策略，可以是go template字符串
func HTTPForwardWithPolicyTemplate(listen []string, serverName []string, policy, dialer string, enableLog bool) config.HTTPConfig {
	return config.HTTPConfig{
		Listen:     listen,
		ServerName: serverName,
		Forward: config.HTTPForwardConfig{
			Policy: policy,
			Dialer: dialer,
			Log:    enableLog,
		},
	}
}

// TunnelServerTemplate 生成隧道服务端配置模板
// authTable: 认证表文件路径
// allowListens: 允许的监听地址，如 ["127.0.0.1", "240.0.0.0/8"]
func TunnelServerTemplate(listen []string, serverName []string, authTable string, allowListens []string) config.HTTPConfig {
	return config.HTTPConfig{
		Listen:     listen,
		ServerName: serverName,
		Tunnel: config.HTTPTunnelConfig{
			Enabled:         true,
			AuthTable:       authTable,
			AllowListens:    allowListens,
			EnableKeepAlive: true,
			Log:             true,
		},
	}
}

// TunnelClientTemplate 生成隧道客户端配置模板
// remoteListen: 公网服务器上的监听地址，如 ["127.0.0.1:10022"]
// proxyPass: 本地需要转发的地址，如 "192.168.50.1:22"
// dialer: 连接公网服务器的拨号器名称
func TunnelClientTemplate(remoteListen []string, proxyPass, resolver, dialer string) config.TunnelConfig {
	return config.TunnelConfig{
		RemoteListen:    remoteListen,
		ProxyPass:       proxyPass,
		Resolver:        resolver,
		Dialer:          dialer,
		DialTimeout:     5,
		EnableKeepAlive: true,
		Log:             true,
	}
}

// DNSForwardTemplate 生成DNS转发配置模板
// listen: 监听地址，如 [":53"]
// proxyPass: 上游DNS服务器，如 "https://8.8.8.8/dns-query"
func DNSForwardTemplate(listen []string, proxyPass string) config.DnsConfig {
	return config.DnsConfig{
		Listen:    listen,
		Policy:    "forward",
		ProxyPass: proxyPass,
		CacheSize: 4096,
		Log:       true,
	}
}

// DNSWithPolicyTemplate 生成带自定义policy的DNS配置
func DNSWithPolicyTemplate(listen []string, policy, proxyPass string, cacheSize int) config.DnsConfig {
	return config.DnsConfig{
		Listen:    listen,
		Policy:    policy,
		ProxyPass: proxyPass,
		CacheSize: cacheSize,
		Log:       true,
	}
}

// WebProxyTemplate 生成Web代理配置模板
// location: web路径，如 "/"
// pass: 代理目标地址，可以用 go template
func WebProxyTemplate(location, pass string) config.HTTPWebConfig {
	return config.HTTPWebConfig{
		Location: location,
		Proxy: config.HTTPWebProxyConfig{
			Pass: pass,
		},
	}
}

// WebDohTemplate 生成Web DoH配置模板
func WebDohTemplate(location, proxyPass string) config.HTTPWebConfig {
	return config.HTTPWebConfig{
		Location: location,
		Doh: config.HTTPWebDohConfig{
			Enabled:   true,
			ProxyPass: proxyPass,
			CacheSize: 4096,
		},
	}
}

// WebIndexTemplate 生成Web Index配置模板
func WebIndexTemplate(location, root string) config.HTTPWebConfig {
	return config.HTTPWebConfig{
		Location: location,
		Index: config.HTTPWebIndexConfig{
			Root: root,
		},
	}
}

// FullConfigTemplate 生成完整配置模板
// 包含全局配置、HTTP转发、DNS等常用功能
func FullConfigTemplate() *config.Config {
	return &config.Config{
		Global: config.GlobalConfig{
			LogLevel:        "info",
			DnsServer:       "https://8.8.8.8/dns-query",
			DisableHttp3:    false,
			DialTimeout:     5,
			IdleConnTimeout: 90,
			MaxIdleConns:    100,
		},
		Dialer: map[string]string{
			"local": "local",
		},
		Https: []config.HTTPConfig{
			{
				Listen:     []string{":443"},
				ServerName: []string{"example.org"},
				Forward: config.HTTPForwardConfig{
					Policy: "proxy_pass",
					Dialer: "local",
					Log:    true,
				},
			},
		},
		Dns: []config.DnsConfig{
			{
				Listen:    []string{":53"},
				Policy:    "forward",
				ProxyPass: "https://8.8.8.8/dns-query",
				CacheSize: 4096,
				Log:       true,
			},
		},
	}
}

// SimpleHTTPProxyTemplate 生成简单HTTP代理配置
func SimpleHTTPProxyTemplate(listen []string, dialer string) *config.Config {
	return &config.Config{
		Global: config.NewDefaultGlobalConfig(),
		Dialer: map[string]string{
			"local": "local",
		},
		Http: []config.HTTPConfig{
			{
				Listen: listen,
				Forward: config.HTTPForwardConfig{
					Policy: "proxy_pass",
					Dialer: dialer,
					Log:    true,
				},
			},
		},
	}
}

// TunnelScenarioTemplate 生成隧道场景完整配置
// role: "server" 或 "client"
func TunnelScenarioTemplate(role string, params map[string]interface{}) *config.Config {
	cfg := &config.Config{
		Global: config.NewDefaultGlobalConfig(),
		Dialer: map[string]string{
			"local": "local",
		},
	}

	if role == "server" {
		// 服务端配置：在HTTPS上启用tunnel
		listen, ok := params["listen"].([]string)
		if !ok {
			return cfg
		}
		serverName, ok := params["server_name"].([]string)
		if !ok {
			return cfg
		}
		authTable, ok := params["auth_table"].(string)
		if !ok {
			return cfg
		}
		allowListens, ok := params["allow_listens"].([]string)
		if !ok {
			return cfg
		}

		cfg.Https = []config.HTTPConfig{
			TunnelServerTemplate(listen, serverName, authTable, allowListens),
		}
	} else if role == "client" {
		// 客户端配置：配置tunnel连接到服务端
		remoteListen, ok := params["remote_listen"].([]string)
		if !ok {
			return cfg
		}
		proxyPass, ok := params["proxy_pass"].(string)
		if !ok {
			return cfg
		}
		resolver, ok := params["resolver"].(string)
		if !ok {
			return cfg
		}
		dialer, ok := params["dialer"].(string)
		if !ok {
			return cfg
		}

		// 需要在dialer中配置连接到公网服务器的拨号器
		if dialerConfig, ok := params["dialer_config"].(string); ok {
			cfg.Dialer[dialer] = dialerConfig
		}

		cfg.Tunnel = []config.TunnelConfig{
			TunnelClientTemplate(remoteListen, proxyPass, resolver, dialer),
		}
	}

	return cfg
}

// SimpleHTTPForwardConfig 生成简单的HTTP转发完整配置
func SimpleHTTPForwardConfig(listen []string, serverName []string, dialerName string, dialerURL string) config.Config {
	global := config.NewDefaultGlobalConfig()

	dialers := make(map[string]string)
	dialers["local"] = "local"
	if dialerName != "" && dialerName != "local" && dialerURL != "" {
		dialers[dialerName] = dialerURL
	}

	https := []config.HTTPConfig{
		HTTPForwardTemplate(listen, serverName, dialerName, true),
	}

	return config.Config{
		Global: global,
		Dialer: dialers,
		Https:  https,
	}
}

// SimpleTunnelConfig 生成简单的隧道完整配置
func SimpleTunnelConfig(
	role string,
	listen []string,
	serverName []string,
	authTable string,
	remoteListen []string,
	proxyPass string,
	dialerName string,
	dialerURL string,
) config.Config {
	global := config.NewDefaultGlobalConfig()
	cfg := config.Config{
		Global: global,
	}

	if role == "server" {
		// 服务端配置
		cfg.Https = []config.HTTPConfig{
			TunnelServerTemplate(listen, serverName, authTable, []string{"127.0.0.1", "240.0.0.0/8"}),
		}
	} else if role == "client" {
		// 客户端配置
		dialers := make(map[string]string)
		dialers["local"] = "local"
		if dialerName != "" && dialerURL != "" {
			dialers[dialerName] = dialerURL
		}
		cfg.Dialer = dialers
		cfg.Tunnel = []config.TunnelConfig{
			TunnelClientTemplate(remoteListen, proxyPass, "https://8.8.8.8/dns-query", dialerName),
		}
	}

	return cfg
}

// SimpleDNSConfig 生成简单的DNS完整配置
func SimpleDNSConfig(listen []string, proxyPass string) config.Config {
	global := config.NewDefaultGlobalConfig()

	dns := []config.DnsConfig{
		DNSForwardTemplate(listen, proxyPass),
	}

	return config.Config{
		Global: global,
		Dns:    dns,
	}
}
