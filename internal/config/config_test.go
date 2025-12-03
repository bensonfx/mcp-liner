package config

import (
	"testing"
)

func TestNewDefaultGlobalConfig(t *testing.T) {
	cfg := NewDefaultGlobalConfig()

	if cfg.LogLevel != "info" {
		t.Errorf("Expected LogLevel to be 'info', got '%s'", cfg.LogLevel)
	}

	if cfg.DnsServer != "https://8.8.8.8/dns-query" {
		t.Errorf("Expected DnsServer to be 'https://8.8.8.8/dns-query', got '%s'", cfg.DnsServer)
	}

	if cfg.DisableHttp3 != false {
		t.Error("Expected DisableHttp3 to be false")
	}
}

func TestConfigToYAML(t *testing.T) {
	cfg := &Config{
		Global: NewDefaultGlobalConfig(),
		Dialer: map[string]string{
			"local": "local",
		},
	}

	yamlStr, err := cfg.ToYAML()
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	if yamlStr == "" {
		t.Error("Expected non-empty YAML string")
	}

	// 验证可以解析回来
	parsedCfg, err := FromYAML(yamlStr)
	if err != nil {
		t.Fatalf("FromYAML failed: %v", err)
	}

	if parsedCfg.Global.LogLevel != cfg.Global.LogLevel {
		t.Error("Parsed config does not match original")
	}
}

func TestNewDefaultHTTPConfig(t *testing.T) {
	listen := []string{":443"}
	serverName := []string{"example.org"}

	cfg := NewDefaultHTTPConfig(listen, serverName)

	if len(cfg.Listen) != 1 || cfg.Listen[0] != ":443" {
		t.Error("Listen not set correctly")
	}

	if len(cfg.ServerName) != 1 || cfg.ServerName[0] != "example.org" {
		t.Error("ServerName not set correctly")
	}

	if cfg.Forward.Policy != "proxy_pass" {
		t.Error("Forward policy not set correctly")
	}
}

func TestNewDefaultTunnelConfig(t *testing.T) {
	remoteListen := []string{"127.0.0.1:10022"}
	proxyPass := "127.0.0.1:22"
	dialer := "cloud"

	cfg := NewDefaultTunnelConfig(remoteListen, proxyPass, dialer)

	if len(cfg.RemoteListen) != 1 || cfg.RemoteListen[0] != "127.0.0.1:10022" {
		t.Error("RemoteListen not set correctly")
	}

	if cfg.ProxyPass != proxyPass {
		t.Error("ProxyPass not set correctly")
	}

	if cfg.Dialer != dialer {
		t.Error("Dialer not set correctly")
	}

	if !cfg.EnableKeepAlive {
		t.Error("Expected EnableKeepAlive to be true")
	}
}

func TestNewDefaultDNSConfig(t *testing.T) {
	listen := []string{":53"}
	proxyPass := "https://8.8.8.8/dns-query"

	cfg := NewDefaultDNSConfig(listen, proxyPass)

	if len(cfg.Listen) != 1 || cfg.Listen[0] != ":53" {
		t.Error("Listen not set correctly")
	}

	if cfg.ProxyPass != proxyPass {
		t.Error("ProxyPass not set correctly")
	}

	if cfg.Policy != "forward" {
		t.Error("Policy should be 'forward'")
	}

	if cfg.CacheSize != 4096 {
		t.Error("CacheSize should be 4096")
	}
}
