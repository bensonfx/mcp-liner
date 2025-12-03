package validation

import (
	"strings"
	"testing"

	"github.com/bensonfx/mcp-liner/internal/config"
)

func TestValidateYAML(t *testing.T) {
	// 测试有效的YAML
	validYAML := `
global:
  log_level: info
  dns_server: https://8.8.8.8/dns-query
`
	err := ValidateYAML(validYAML)
	if err != nil {
		t.Errorf("Valid YAML should not produce error: %v", err)
	}

	// 测试无效的YAML
	invalidYAML := `
global:
  log_level: info
  invalid syntax here
`
	err = ValidateYAML(invalidYAML)
	if err == nil {
		t.Error("Invalid YAML should produce error")
	}
}

func TestValidateConfig(t *testing.T) {
	// 测试有效配置
	cfg := &config.Config{
		Global: config.GlobalConfig{
			LogLevel:  "info",
			DnsServer: "https://8.8.8.8/dns-query",
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
				},
			},
		},
	}

	result := ValidateConfig(cfg)
	if !result.Valid {
		t.Errorf("Valid config should pass validation, errors: %v", result.Errors)
	}

	// 测试缺少server_name的HTTPS配置
	badCfg := &config.Config{
		Global: config.GlobalConfig{},
		Https: []config.HTTPConfig{
			{
				Listen: []string{":443"},
				// Missing ServerName
			},
		},
	}

	result = ValidateConfig(badCfg)
	if result.Valid {
		t.Error("Config with missing server_name should fail validation")
	}
	if len(result.Errors) == 0 {
		t.Error("Should have validation errors")
	}
}

func TestValidateDialerReferences(t *testing.T) {
	// 测试未定义的dialer引用
	cfg := &config.Config{
		Global: config.GlobalConfig{},
		Dialer: map[string]string{
			"local": "local",
		},
		Https: []config.HTTPConfig{
			{
				Listen:     []string{":443"},
				ServerName: []string{"example.org"},
				Forward: config.HTTPForwardConfig{
					Dialer: "undefined_dialer", // 引用未定义的dialer
				},
			},
		},
	}

	result := ValidateConfig(cfg)
	if result.Valid {
		t.Error("Config with undefined dialer should fail validation")
	}

	found := false
	for _, err := range result.Errors {
		if strings.Contains(err.Message, "undefined_dialer") && strings.Contains(err.Message, "not defined") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Should have error about undefined dialer")
	}
}

func TestValidateTunnelConfig(t *testing.T) {
	// 测试缺少必填字段的tunnel配置
	cfg := &config.Config{
		Global: config.GlobalConfig{},
		Tunnel: []config.TunnelConfig{
			{
				// Missing remote_listen, proxy_pass, dialer
			},
		},
	}

	result := ValidateConfig(cfg)
	if result.Valid {
		t.Error("Tunnel config with missing required fields should fail")
	}
	if len(result.Errors) < 3 {
		t.Error("Should have multiple validation errors for missing fields")
	}
}

func TestFormatValidationErrors(t *testing.T) {
	result := &ValidationResult{
		Valid: false,
		Errors: []ValidationError{
			{Field: "test.field", Message: "test error"},
		},
	}

	formatted := FormatValidationErrors(result)
	if !strings.Contains(formatted, "test.field") {
		t.Error("Formatted errors should contain field name")
	}
	if !strings.Contains(formatted, "test error") {
		t.Error("Formatted errors should contain error message")
	}
}
