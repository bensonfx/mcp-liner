// Package validation 提供liner配置验证功能
package validation

import (
	"fmt"
	"strings"

	"github.com/bensonfx/mcp-liner/internal/config"
	"gopkg.in/yaml.v3"
)

// ValidationError 验证错误
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationResult 验证结果
type ValidationResult struct {
	Valid  bool
	Errors []ValidationError
}

// ValidateYAML 验证YAML语法
func ValidateYAML(yamlContent string) error {
	var data interface{}
	err := yaml.Unmarshal([]byte(yamlContent), &data)
	if err != nil {
		return fmt.Errorf("invalid YAML syntax: %w", err)
	}
	return nil
}

// ValidateConfig 验证liner配置
func ValidateConfig(cfg *config.Config) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Errors: []ValidationError{},
	}

	// 验证全局配置
	validateGlobal(&cfg.Global, result)

	// 验证Dialer配置
	validateDialers(cfg.Dialer, result)

	// 验证HTTPS配置
	for i, httpsCfg := range cfg.Https {
		validateHTTPConfig(httpsCfg, fmt.Sprintf("https[%d]", i), result, true)
	}

	// 验证HTTP配置
	for i, httpCfg := range cfg.Http {
		validateHTTPConfig(httpCfg, fmt.Sprintf("http[%d]", i), result, false)
	}

	// 验证隧道配置
	for i, tunnelCfg := range cfg.Tunnel {
		validateTunnelConfig(tunnelCfg, fmt.Sprintf("tunnel[%d]", i), result)
	}

	// 验证DNS配置
	for i, dnsCfg := range cfg.Dns {
		validateDNSConfig(dnsCfg, fmt.Sprintf("dns[%d]", i), result)
	}

	// 验证Socks配置
	for i, socksCfg := range cfg.Socks {
		validateSocksConfig(socksCfg, fmt.Sprintf("socks[%d]", i), result)
	}

	// 验证dialer引用
	validateDialerReferences(cfg, result)

	if len(result.Errors) > 0 {
		result.Valid = false
	}

	return result
}

// validateGlobal 验证全局配置
func validateGlobal(global *config.GlobalConfig, result *ValidationResult) {
	// 验证日志级别
	if global.LogLevel != "" {
		validLevels := []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}
		if !contains(validLevels, global.LogLevel) {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "global.log_level",
				Message: fmt.Sprintf("invalid log level: %s, must be one of: %s", global.LogLevel, strings.Join(validLevels, ", ")),
			})
		}
	}

	// 验证DNS服务器格式
	if global.DnsServer != "" {
		if !strings.HasPrefix(global.DnsServer, "https://") && !strings.HasPrefix(global.DnsServer, "udp://") && !strings.Contains(global.DnsServer, ":") {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "global.dns_server",
				Message: "dns_server should be a valid DNS server address (e.g., 'https://8.8.8.8/dns-query' or '8.8.8.8:53')",
			})
		}
	}
}

// validateDialers 验证拨号器配置
func validateDialers(dialers map[string]string, result *ValidationResult) {
	if len(dialers) == 0 {
		return
	}

	for name, url := range dialers {
		if name == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "dialer",
				Message: "dialer name cannot be empty",
			})
		}
		if url == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("dialer.%s", name),
				Message: "dialer URL cannot be empty",
			})
		}
	}
}

// validateHTTPConfig 验证HTTP/HTTPS配置
func validateHTTPConfig(httpCfg config.HTTPConfig, prefix string, result *ValidationResult, isHTTPS bool) {
	// 验证listen字段
	if len(httpCfg.Listen) == 0 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   fmt.Sprintf("%s.listen", prefix),
			Message: "listen field is required and cannot be empty",
		})
	}

	// HTTPS需要server_name
	if isHTTPS && len(httpCfg.ServerName) == 0 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   fmt.Sprintf("%s.server_name", prefix),
			Message: "server_name field is required for HTTPS configuration",
		})
	}

	// 如果配置了forward，验证forward配置
	if httpCfg.Forward.Policy != "" || httpCfg.Forward.Dialer != "" {
		if httpCfg.Forward.Policy == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("%s.forward.policy", prefix),
				Message: "policy is required when forward is configured",
			})
		}
	}
}

// validateTunnelConfig 验证隧道配置
func validateTunnelConfig(tunnelCfg config.TunnelConfig, prefix string, result *ValidationResult) {
	// 验证remote_listen字段
	if len(tunnelCfg.RemoteListen) == 0 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   fmt.Sprintf("%s.remote_listen", prefix),
			Message: "remote_listen field is required and cannot be empty",
		})
	}

	// 验证proxy_pass字段
	if tunnelCfg.ProxyPass == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   fmt.Sprintf("%s.proxy_pass", prefix),
			Message: "proxy_pass field is required",
		})
	}

	// 验证dialer字段
	if tunnelCfg.Dialer == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   fmt.Sprintf("%s.dialer", prefix),
			Message: "dialer field is required",
		})
	}
}

// validateDNSConfig 验证DNS配置
func validateDNSConfig(dnsCfg config.DnsConfig, prefix string, result *ValidationResult) {
	// 验证listen字段
	if len(dnsCfg.Listen) == 0 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   fmt.Sprintf("%s.listen", prefix),
			Message: "listen field is required and cannot be empty",
		})
	}

	// 验证proxy_pass字段（如果policy是forward）
	if dnsCfg.Policy == "forward" && dnsCfg.ProxyPass == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   fmt.Sprintf("%s.proxy_pass", prefix),
			Message: "proxy_pass is required when policy is 'forward'",
		})
	}
}

// validateSocksConfig 验证Socks配置
func validateSocksConfig(socksCfg config.SocksConfig, prefix string, result *ValidationResult) {
	// 验证listen字段
	if len(socksCfg.Listen) == 0 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   fmt.Sprintf("%s.listen", prefix),
			Message: "listen field is required and cannot be empty",
		})
	}
}

// validateDialerReferences 验证dialer引用
func validateDialerReferences(cfg *config.Config, result *ValidationResult) {
	// 收集所有定义的dialer
	定义的dialers := make(map[string]bool)
	定义的dialers["local"] = true // local是内置dialer
	for name := range cfg.Dialer {
		定义的dialers[name] = true
	}

	// 检查HTTP配置中的dialer引用
	for i, httpsCfg := range cfg.Https {
		if httpsCfg.Forward.Dialer != "" && !定义的dialers[httpsCfg.Forward.Dialer] {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("https[%d].forward.dialer", i),
				Message: fmt.Sprintf("dialer '%s' is not defined", httpsCfg.Forward.Dialer),
			})
		}
	}

	for i, httpCfg := range cfg.Http {
		if httpCfg.Forward.Dialer != "" && !定义的dialers[httpCfg.Forward.Dialer] {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("http[%d].forward.dialer", i),
				Message: fmt.Sprintf("dialer '%s' is not defined", httpCfg.Forward.Dialer),
			})
		}
	}

	// 检查隧道配置中的dialer引用
	for i, tunnelCfg := range cfg.Tunnel {
		if tunnelCfg.Dialer != "" && !定义的dialers[tunnelCfg.Dialer] {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("tunnel[%d].dialer", i),
				Message: fmt.Sprintf("dialer '%s' is not defined", tunnelCfg.Dialer),
			})
		}
	}

	// 检查Socks配置中的dialer引用
	for i, socksCfg := range cfg.Socks {
		if socksCfg.Forward.Dialer != "" && !定义的dialers[socksCfg.Forward.Dialer] {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("socks[%d].forward.dialer", i),
				Message: fmt.Sprintf("dialer '%s' is not defined", socksCfg.Forward.Dialer),
			})
		}
	}
}

// contains 检查字符串是否在切片中
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// FormatValidationErrors 格式化验证错误
func FormatValidationErrors(result *ValidationResult) string {
	if result.Valid {
		return "Configuration is valid"
	}

	var builder strings.Builder
	builder.WriteString("Configuration validation failed:\n\n")
	for i, err := range result.Errors {
		builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, err.Error()))
	}
	return builder.String()
}
