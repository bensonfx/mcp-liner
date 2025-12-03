// Package responses 提供MCP工具响应格式化功能
package responses

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bensonfx/mcp-liner/internal/validation"
)

// MCPResponse MCP标准响应格式
type MCPResponse struct {
	Content []ContentBlock `json:"content"`
	IsError bool           `json:"isError,omitempty"`
}

// ContentBlock 内容块
type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// SuccessResponse 创建成功响应
// yamlContent: 生成的YAML配置内容
// description: 配置说明（可选）
func SuccessResponse(yamlContent string, description string) (string, error) {
	var textBuilder strings.Builder

	if description != "" {
		textBuilder.WriteString(description)
		textBuilder.WriteString("\n\n")
	}

	textBuilder.WriteString("Generated Liner Configuration:\n\n")
	textBuilder.WriteString("```yaml\n")
	textBuilder.WriteString(yamlContent)
	textBuilder.WriteString("\n```")

	response := MCPResponse{
		Content: []ContentBlock{
			{
				Type: "text",
				Text: textBuilder.String(),
			},
		},
		IsError: false,
	}

	return marshalResponse(response)
}

// ErrorResponse 创建错误响应
// errorMsg: 错误信息
// suggestion: 建议（可选）
func ErrorResponse(errorMsg string, suggestion string) (string, error) {
	var textBuilder strings.Builder

	textBuilder.WriteString("Error: ")
	textBuilder.WriteString(errorMsg)
	textBuilder.WriteString("\n")

	if suggestion != "" {
		textBuilder.WriteString("\nSuggestion: ")
		textBuilder.WriteString(suggestion)
	}

	response := MCPResponse{
		Content: []ContentBlock{
			{
				Type: "text",
				Text: textBuilder.String(),
			},
		},
		IsError: true,
	}

	return marshalResponse(response)
}

// ValidationResponse 创建验证响应
// result: 验证结果
func ValidationResponse(result *validation.ValidationResult) (string, error) {
	if result.Valid {
		response := MCPResponse{
			Content: []ContentBlock{
				{
					Type: "text",
					Text: "✅ Configuration validation passed!\n\nThe liner configuration is valid and ready to use.",
				},
			},
			IsError: false,
		}
		return marshalResponse(response)
	}

	// 构建错误信息
	var textBuilder strings.Builder
	textBuilder.WriteString("❌ Configuration validation failed!\n\n")
	textBuilder.WriteString(fmt.Sprintf("Found %d error(s):\n\n", len(result.Errors)))

	for i, err := range result.Errors {
		textBuilder.WriteString(fmt.Sprintf("%d. **%s**: %s\n", i+1, err.Field, err.Message))
	}

	textBuilder.WriteString("\nPlease fix these errors and try again.")

	response := MCPResponse{
		Content: []ContentBlock{
			{
				Type: "text",
				Text: textBuilder.String(),
			},
		},
		IsError: true,
	}

	return marshalResponse(response)
}

// DocumentationResponse 创建文档响应
// topic: 主题
// content: 文档内容
func DocumentationResponse(topic string, content string) (string, error) {
	var textBuilder strings.Builder

	textBuilder.WriteString(fmt.Sprintf("# Liner Documentation: %s\n\n", topic))
	textBuilder.WriteString(content)

	response := MCPResponse{
		Content: []ContentBlock{
			{
				Type: "text",
				Text: textBuilder.String(),
			},
		},
		IsError: false,
	}

	return marshalResponse(response)
}

// InfoResponse 创建信息响应
// title: 标题
// info: 信息内容
func InfoResponse(title string, info string) (string, error) {
	var textBuilder strings.Builder

	if title != "" {
		textBuilder.WriteString(fmt.Sprintf("## %s\n\n", title))
	}

	textBuilder.WriteString(info)

	response := MCPResponse{
		Content: []ContentBlock{
			{
				Type: "text",
				Text: textBuilder.String(),
			},
		},
		IsError: false,
	}

	return marshalResponse(response)
}

// ConfigWithExampleResponse 创建带示例的配置响应
// yamlContent: 生成的YAML配置
// example: 使用示例
func ConfigWithExampleResponse(yamlContent string, example string) (string, error) {
	var textBuilder strings.Builder

	textBuilder.WriteString("Generated Liner Configuration:\n\n")
	textBuilder.WriteString("```yaml\n")
	textBuilder.WriteString(yamlContent)
	textBuilder.WriteString("\n```\n\n")

	if example != "" {
		textBuilder.WriteString("## Usage Example\n\n")
		textBuilder.WriteString(example)
	}

	response := MCPResponse{
		Content: []ContentBlock{
			{
				Type: "text",
				Text: textBuilder.String(),
			},
		},
		IsError: false,
	}

	return marshalResponse(response)
}

// marshalResponse 序列化响应为JSON
func marshalResponse(response MCPResponse) (string, error) {
	data, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}
	return string(data), nil
}

// SimpleTextResponse 创建简单文本响应（用于调试）
func SimpleTextResponse(text string) (string, error) {
	response := MCPResponse{
		Content: []ContentBlock{
			{
				Type: "text",
				Text: text,
			},
		},
		IsError: false,
	}
	return marshalResponse(response)
}
