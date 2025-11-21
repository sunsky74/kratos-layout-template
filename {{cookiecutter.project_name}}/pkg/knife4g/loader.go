package knife4g

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadOpenAPIFromFile 从文件加载 OpenAPI 配置
func LoadOpenAPIFromFile(filepath string) (*OpenAPI3, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read openapi file: %w", err)
	}

	var openapi OpenAPI3
	if err := yaml.Unmarshal(data, &openapi); err != nil {
		return nil, fmt.Errorf("failed to parse openapi yaml: %w", err)
	}

	return &openapi, nil
}

// NewDefaultConfig 创建默认的 Knife4g 配置
// 默认使用 docs/api/openapi.yaml 作为 OpenAPI 配置文件
func NewDefaultConfig(serverName string) (*Config, error) {
	openapiPath := "docs/api/openapi.yaml"
	openapi, err := LoadOpenAPIFromFile(openapiPath)
	if err != nil {
		return nil, err
	}

	return &Config{
		RelativePath: "",
		ServerName:   serverName,
		OpenAPI:      openapi,
	}, nil
}
