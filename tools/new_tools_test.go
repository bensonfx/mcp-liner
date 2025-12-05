package tools

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGeneratePolicyExamples(t *testing.T) {
	tests := []struct {
		name        string
		params      GeneratePolicyExamplesParams
		wantErr     bool
		wantContain string
	}{
		{
			name: "GeoIP policy example",
			params: GeneratePolicyExamplesParams{
				ConfigType: "http_forward",
				PolicyType: "geoip",
			},
			wantErr:     false,
			wantContain: "GeoIP-based Routing Policy",
		},
		{
			name: "Geosite policy example",
			params: GeneratePolicyExamplesParams{
				ConfigType: "http_forward",
				PolicyType: "geosite",
			},
			wantErr:     false,
			wantContain: "Geosite-based Routing Policy",
		},
		{
			name: "Domain match policy example",
			params: GeneratePolicyExamplesParams{
				ConfigType: "sni_forward",
				PolicyType: "domain_match",
			},
			wantErr:     false,
			wantContain: "Domain Matching Policy",
		},
		{
			name: "IP range policy example",
			params: GeneratePolicyExamplesParams{
				ConfigType: "http_forward",
				PolicyType: "ip_range",
			},
			wantErr:     false,
			wantContain: "IP Range Routing Policy",
		},
		{
			name: "File-based policy example",
			params: GeneratePolicyExamplesParams{
				ConfigType: "socks_forward",
				PolicyType: "file_based",
			},
			wantErr:     false,
			wantContain: "File-based Routing Policy",
		},
		{
			name: "Fetch-based policy example",
			params: GeneratePolicyExamplesParams{
				ConfigType: "http_forward",
				PolicyType: "fetch_based",
			},
			wantErr:     false,
			wantContain: "Fetch-based Routing Policy",
		},
		{
			name: "Custom policy example",
			params: GeneratePolicyExamplesParams{
				ConfigType: "http_forward",
				PolicyType: "custom",
			},
			wantErr:     false,
			wantContain: "Go Template Built-in Functions",
		},
		{
			name: "Custom policy contains Sprig functions",
			params: GeneratePolicyExamplesParams{
				ConfigType: "dns",
				PolicyType: "custom",
			},
			wantErr:     false,
			wantContain: "Sprig Library Functions",
		},
		{
			name: "Invalid policy type",
			params: GeneratePolicyExamplesParams{
				ConfigType: "http_forward",
				PolicyType: "invalid",
			},
			wantErr:     false, // Error responses are returned as formatted content
			wantContain: "Unknown policy_type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.params)
			if err != nil {
				t.Fatalf("failed to marshal params: %v", err)
			}

			result, err := GeneratePolicyExamples(jsonData)
			if err != nil {
				t.Fatalf("GeneratePolicyExamples() unexpected error: %v", err)
			}

			if !strings.Contains(result, tt.wantContain) {
				t.Errorf("GeneratePolicyExamples() result does not contain %q", tt.wantContain)
			}
		})
	}
}

func TestGenerateRedsocksConfig(t *testing.T) {
	tests := []struct {
		name        string
		params      GenerateRedsocksConfigParams
		wantErr     bool
		wantContain string
	}{
		{
			name: "Basic redsocks config",
			params: GenerateRedsocksConfigParams{
				Listen:    []string{":12345"},
				Dialer:    "proxy",
				DialerURL: "socks5://127.0.0.1:1080",
				Log:       true,
			},
			wantErr:     false,
			wantContain: "redsocks:",
		},
		{
			name: "Redsocks with default values",
			params: GenerateRedsocksConfigParams{
				Dialer: "local",
			},
			wantErr:     false,
			wantContain: "listen:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.params)
			if err != nil {
				t.Fatalf("failed to marshal params: %v", err)
			}

			result, err := GenerateRedsocksConfig(jsonData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRedsocksConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !strings.Contains(result, tt.wantContain) {
				t.Errorf("GenerateRedsocksConfig() result does not contain %q", tt.wantContain)
			}
		})
	}
}

func TestGenerateRedsocksIptables(t *testing.T) {
	tests := []struct {
		name        string
		params      GenerateRedsocksIptablesParams
		wantErr     bool
		wantContain string
	}{
		{
			name: "iptables-save format",
			params: GenerateRedsocksIptablesParams{
				RedsocksPort: 12345,
				LANInterface: "eth0",
				WANInterface: "eth1",
				ProxyPorts:   []int{80, 443},
				Format:       "iptables-save",
			},
			wantErr:     false,
			wantContain: "*nat",
		},
		{
			name: "shell-script format",
			params: GenerateRedsocksIptablesParams{
				RedsocksPort: 12345,
				LANInterface: "mlan0",
				WANInterface: "enp1s0",
				Format:       "shell-script",
			},
			wantErr:     false,
			wantContain: "#!/bin/bash",
		},
		{
			name: "Contains loop prevention rules",
			params: GenerateRedsocksIptablesParams{
				RedsocksPort: 12345,
				Format:       "iptables-save",
			},
			wantErr:     false,
			wantContain: "127.0.0.0/8",
		},
		{
			name: "Contains REDSOCKS chain",
			params: GenerateRedsocksIptablesParams{
				RedsocksPort: 12345,
				Format:       "iptables-save",
			},
			wantErr:     false,
			wantContain: ":REDSOCKS",
		},
		{
			name: "Invalid format",
			params: GenerateRedsocksIptablesParams{
				RedsocksPort: 12345,
				Format:       "invalid",
			},
			wantErr:     false, // Error responses are returned as formatted content
			wantContain: "Unknown format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.params)
			if err != nil {
				t.Fatalf("failed to marshal params: %v", err)
			}

			result, err := GenerateRedsocksIptables(jsonData)
			if err != nil {
				t.Fatalf("GenerateRedsocksIptables() unexpected error: %v", err)
			}

			if !strings.Contains(result, tt.wantContain) {
				t.Errorf("GenerateRedsocksIptables() result does not contain %q", tt.wantContain)
			}
		})
	}
}

func TestGenerateSniConfig(t *testing.T) {
	tests := []struct {
		name        string
		params      GenerateSniConfigParams
		wantErr     bool
		wantContain string
	}{
		{
			name: "Basic SNI config",
			params: GenerateSniConfigParams{
				Enabled: true,
				Policy:  "proxy_pass",
				Dialer:  "local",
				Log:     true,
			},
			wantErr:     false,
			wantContain: "sni:",
		},
		{
			name: "SNI with policy template",
			params: GenerateSniConfigParams{
				Enabled:     true,
				Policy:      "{{ if hasSuffixes \"google.com\" .ServerName }}proxy{{ else }}direct{{ end }}",
				Dialer:      "proxy",
				DialerURL:   "socks5://127.0.0.1:1080",
				DisableIpv6: true,
			},
			wantErr:     false,
			wantContain: "policy:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.params)
			if err != nil {
				t.Fatalf("failed to marshal params: %v", err)
			}

			result, err := GenerateSniConfig(jsonData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateSniConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !strings.Contains(result, tt.wantContain) {
				t.Errorf("GenerateSniConfig() result does not contain %q", tt.wantContain)
			}
		})
	}
}

func TestGenerateStreamConfig(t *testing.T) {
	tests := []struct {
		name        string
		params      GenerateStreamConfigParams
		wantErr     bool
		wantContain string
	}{
		{
			name: "Basic stream config",
			params: GenerateStreamConfigParams{
				Listen:    []string{":3389"},
				ProxyPass: "192.168.1.100:3389",
				Dialer:    "local",
				Log:       true,
			},
			wantErr:     false,
			wantContain: "stream:",
		},
		{
			name: "Stream with TLS",
			params: GenerateStreamConfigParams{
				Listen:    []string{":443"},
				ProxyPass: "127.0.0.1:8080",
				Keyfile:   "/path/to/key.pem",
				Certfile:  "/path/to/cert.pem",
			},
			wantErr:     false,
			wantContain: "keyfile:",
		},
		{
			name: "Stream with PROXY protocol",
			params: GenerateStreamConfigParams{
				Listen:        []string{":8080"},
				ProxyPass:     "backend:8080",
				ProxyProtocol: 2,
			},
			wantErr:     false,
			wantContain: "proxy_protocol:",
		},
		{
			name: "Stream without proxy_pass",
			params: GenerateStreamConfigParams{
				Listen: []string{":8080"},
			},
			wantErr:     false, // Error responses are returned as formatted content
			wantContain: "proxy_pass is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.params)
			if err != nil {
				t.Fatalf("failed to marshal params: %v", err)
			}

			result, err := GenerateStreamConfig(jsonData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateStreamConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !strings.Contains(result, tt.wantContain) {
				t.Errorf("GenerateStreamConfig() result does not contain %q", tt.wantContain)
			}
		})
	}
}
