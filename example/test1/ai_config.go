package main

import "errors"

// AIConfig AI系统配置
type AIConfig struct {
	DefaultModule AiModuleConfig        `yaml:"default_module"`
	DeleteAudio   bool                  `yaml:"delete_audio,omitempty"` // 是否删除临时音频文件
	ASR           map[string]ASRConfig  `yaml:"ASR"`
	TTS           map[string]TTSConfig  `yaml:"TTS"`
	LLM           map[string]LLMConfig  `yaml:"LLM"`
	VLLM          map[string]VLLMConfig `yaml:"VLLLM"`
	CMDExit       []string              `yaml:"CMD_exit"`
}

func (c *AIConfig) GetAsrConfig(module string) (ASRConfig, error) {
	if module == "" {
		module = c.DefaultModule.ASR
	}
	if _, ok := c.ASR[module]; !ok {
		if config, ok := c.ASR[module]; ok {
			return config, nil
		}
		return ASRConfig{}, errors.New("default ASR module is not set")
	}
	return c.ASR[module], nil
}

func (c *AIConfig) GetTTSConfig(module string) (TTSConfig, error) {
	if module == "" {
		module = c.DefaultModule.TTS
	}
	if _, ok := c.TTS[module]; !ok {
		if config, ok := c.TTS[module]; ok {
			return config, nil
		}
		return TTSConfig{}, errors.New("default TTS module is not set")
	}
	return c.TTS[module], nil
}

func (c *AIConfig) GetLLMConfig(module string) (LLMConfig, error) {
	if module == "" {
		module = c.DefaultModule.LLM
	}
	if _, ok := c.LLM[module]; !ok {
		if config, ok := c.LLM[module]; ok {
			return config, nil
		}
		return LLMConfig{}, errors.New("default LLM module is not set")
	}
	return c.LLM[module], nil
}

func (c *AIConfig) GetVLLMConfig(module string) (VLLMConfig, error) {
	if module == "" {
		module = c.DefaultModule.VLLM
	}
	if _, ok := c.VLLM[module]; !ok {
		if config, ok := c.VLLM[module]; ok {
			return config, nil
		}
		return VLLMConfig{}, errors.New("default VLLLM module is not set")
	}
	return c.VLLM[module], nil
}

// AiModuleConfig 默认使用的模块配置
type AiModuleConfig struct {
	ASR  string `yaml:"asr"`
	TTS  string `yaml:"tts"`
	LLM  string `yaml:"llm"`
	VLLM string `yaml:"vllm"`
}

// ASRConfig ASR配置
type ASRConfig struct {
	Name          string `yaml:"name"`
	Type          string `yaml:"type"`
	AppID         string `yaml:"appid,omitempty"`
	AccessToken   string `yaml:"access_token,omitempty"`
	OutputDir     string `yaml:"output_dir,omitempty"`
	EndWindowSize int    `yaml:"end_window_size,omitempty"`
	Addr          string `yaml:"addr,omitempty"`
	APIKey        string `yaml:"api_key,omitempty"`
	Model         string `yaml:"model,omitempty"`
	Voice         string `yaml:"voice,omitempty"`
	Prompt        string `yaml:"prompt,omitempty"`
	Lang          string `yaml:"lang,omitempty"`
}

// TTSConfig TTS配置
type TTSConfig struct {
	Name            string      `yaml:"name"`
	Type            string      `yaml:"type"`
	Voice           string      `yaml:"voice,omitempty"`
	Format          string      `yaml:"format"           json:"format"` // 输出格式
	OutputDir       string      `yaml:"output_dir,omitempty"`
	AppID           string      `yaml:"appid,omitempty"`
	Token           string      `yaml:"token,omitempty"`
	Cluster         string      `yaml:"cluster,omitempty"`
	SampleRate      int         `yaml:"sample_rate,omitempty"`
	SupportedVoices []VoiceInfo `yaml:"supported_voices,omitempty"`
}

// VoiceInfo 语音信息
type VoiceInfo struct {
	Name        string `json:"name" yaml:"name"`
	DisplayName string `json:"display_name" yaml:"display_name"`
	Sex         string `json:"sex" yaml:"sex"`
	Description string `json:"description" yaml:"description"`
	AudioURL    string `json:"audio_url" yaml:"audio_url"`
}

// LLMConfig LLM配置
type LLMConfig struct {
	Name                string                 `yaml:"name"`
	Type                string                 `yaml:"type"`
	ModelName           string                 `yaml:"model_name,omitempty"`
	URL                 string                 `yaml:"url,omitempty"`
	APIKey              string                 `yaml:"api_key,omitempty"`
	Temperature         float64                `yaml:"temperature" json:"temperature"` // 温度参数
	TopP                float64                `yaml:"top_p"       json:"top_p"`       // TopP参数
	BotID               string                 `yaml:"bot_id,omitempty"`
	UserID              string                 `yaml:"user_id,omitempty"`
	ClientID            string                 `yaml:"client_id,omitempty"`
	PublicKey           string                 `yaml:"public_key,omitempty"`
	PrivateKey          string                 `yaml:"private_key,omitempty"`
	PersonalAccessToken string                 `yaml:"personal_access_token,omitempty"`
	Thinking            string                 `yaml:"thinking,omitempty"`
	MaxTokens           int                    `yaml:"max_tokens,omitempty"`
	Extra               map[string]interface{} `yaml:",inline"     json:"extra"` // 额外配置
}

// VLLMConfig VLLLM配置结构（视觉语言大模型）
type VLLMConfig struct {
	Name        string                 `yaml:"name"`
	Type        string                 `yaml:"type"        json:"type"`        // API类型，复用LLM的类型
	ModelName   string                 `yaml:"model_name"  json:"model_name"`  // 模型名称，使用支持视觉的模型
	BaseURL     string                 `yaml:"url"         json:"url"`         // API地址
	APIKey      string                 `yaml:"api_key"     json:"api_key"`     // API密钥
	Temperature float64                `yaml:"temperature" json:"temperature"` // 温度参数
	MaxTokens   int                    `yaml:"max_tokens"  json:"max_tokens"`  // 最大令牌数
	TopP        float64                `yaml:"top_p"       json:"top_p"`       // TopP参数
	Security    SecurityConfig         `yaml:"security"    json:"security"`    // 图片安全配置
	Extra       map[string]interface{} `yaml:",inline"     json:"extra"`       // 额外配置

}

// SecurityConfig 安全配置
type SecurityConfig struct {
	MaxFileSize       int64    `json:"max_file_size" yaml:"max_file_size"`
	MaxPixels         int64    `json:"max_pixels" yaml:"max_pixels"`
	MaxWidth          int      `json:"max_width" yaml:"max_width"`
	MaxHeight         int      `json:"max_height" yaml:"max_height"`
	AllowedFormats    []string `json:"allowed_formats" yaml:"allowed_formats"`
	EnableDeepScan    bool     `json:"enable_deep_scan" yaml:"enable_deep_scan"`
	ValidationTimeout string   `json:"validation_timeout" yaml:"validation_timeout"`
}
