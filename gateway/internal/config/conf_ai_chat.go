package config

type AiChatConfig struct {
	VectorCollection string `yaml:"vector_collection"`
	DefaultTopK      int    `yaml:"default_top_k"`
	// AssistantEnabled 控制站内 AI 助手是否可用。
	// 未配置时默认开启；生产环境缺少 Chroma/agent-server 时显式设为 false。
	AssistantEnabled *bool `yaml:"assistant_enabled"`
}

// IsAssistantEnabled 返回助手开关，未配置时默认开启。
func (c AiChatConfig) IsAssistantEnabled() bool {
	if c.AssistantEnabled == nil {
		return true
	}
	return *c.AssistantEnabled
}
