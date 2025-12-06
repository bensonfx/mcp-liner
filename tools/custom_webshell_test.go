package tools

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGenerateWebshellConfig(t *testing.T) {
	// 测试参数
	params := map[string]interface{}{
		"listen":      []string{":8443"},
		"server_name": []string{"shell.test.com"},
		"command":     "bash",
		"home":        "/home/user",
		"auth_table":  "users.csv",
		"location":    "/terminal/",
	}

	jsonParams, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal params: %v", err)
	}

	// 调用工具
	result, err := GenerateWebshellConfig(jsonParams)
	if err != nil {
		t.Fatalf("GenerateWebshellConfig returned error: %v", err)
	}

	// 验证结果包含预期的YAML内容片段
	expectedStrings := []string{
		":8443",
		"shell.test.com",
		"location: /terminal/",
		"enabled: true",
		"command: bash",
		"home: /home/user",
		"auth_table: users.csv",
	}

	for _, s := range expectedStrings {
		if !strings.Contains(result, s) {
			t.Errorf("Expected output to contain %q, but it didn't.\nOutput:\n%s", s, result)
		}
	}
}
