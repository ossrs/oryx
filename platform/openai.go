// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT
package main

import (
	"fmt"
	"strings"
)

func gptModelSupportSystem(model string) bool {
	if strings.HasPrefix(model, "o1-") {
		return false
	}
	return true
}

func gptModelSupportStream(model string) bool {
	if strings.HasPrefix(model, "o1-") {
		return false
	}
	return true
}

func gptModelSupportMaxTokens(model string, maxTokens int) int {
	if strings.HasPrefix(model, "o1-") {
		return 0
	}
	return maxTokens
}

func gptModelSupportTemperature(model string, temperature float32) float32 {
	if strings.HasPrefix(model, "o1-") {
		return 0.0
	}
	return temperature
}

type SrsAssistantProvider struct {
	// The AI provider.
	AIProvider string `json:"aiProvider"`
	// The AI secret key.
	AISecretKey string `json:"aiSecretKey"`
	// The AI organization.
	AIOrganization string `json:"aiOrganization"`
	// The AI base URL.
	AIBaseURL string `json:"aiBaseURL"`
}

func (v *SrsAssistantProvider) String() string {
	return fmt.Sprintf("provider=%v, secretKey=%vB, baseURL=%v",
		v.AIProvider, len(v.AISecretKey), v.AIBaseURL)
}

type SrsAssistantASR struct {
	// Whether enable the AI ASR.
	AIASREnabled bool `json:"aiAsrEnabled"`
	// The AI asr language.
	AIASRLanguage string `json:"aiAsrLanguage"`
	// The AI asr prompt type. user or user-ai.
	AIASRPrompt string `json:"aiAsrPrompt"`
}

func (v *SrsAssistantASR) String() string {
	return fmt.Sprintf("enabled=%v,language=%v,prompt=%v",
		v.AIASREnabled, v.AIASRLanguage, v.AIASRPrompt)
}

type SrsAssistantChat struct {
	// Whether enable the AI processing.
	AIChatEnabled bool `json:"aiChatEnabled"`
	// The AI model name.
	AIChatModel string `json:"aiChatModel"`
	// The AI chat system prompt.
	AIChatPrompt string `json:"aiChatPrompt"`
	// The AI chat max window.
	AIChatMaxWindow int `json:"aiChatMaxWindow"`
	// The AI chat max words.
	AIChatMaxWords int `json:"aiChatMaxWords"`
}

func (v *SrsAssistantChat) String() string {
	return fmt.Sprintf("enabled=%v,model=%v,prompt=%v,window=%v,words=%v",
		v.AIChatEnabled, v.AIChatModel, v.AIChatPrompt, v.AIChatMaxWindow, v.AIChatMaxWords)
}

type SrsAssistantPost struct {
	// Whether enable the AI post processing.
	AIPostEnabled bool `json:"aiPostEnabled"`
	// The AI model name.
	AIPostModel string `json:"aiPostModel"`
	// The AI chat system prompt.
	AIPostPrompt string `json:"aiPostPrompt"`
	// The AI chat max window.
	AIPostMaxWindow int `json:"aiPostMaxWindow"`
	// The AI chat max words.
	AIPostMaxWords int `json:"aiPostMaxWords"`
}

func (v *SrsAssistantPost) String() string {
	return fmt.Sprintf("enabled=%v,model=%v,prompt=%v,window=%v,words=%v",
		v.AIPostEnabled, v.AIPostModel, v.AIPostPrompt, v.AIPostMaxWindow, v.AIPostMaxWords)
}

type SrsAssistantTTS struct {
	// Whether enable the AI TTS.
	AITTSEnabled bool `json:"aiTtsEnabled"`
}

func (v *SrsAssistantTTS) String() string {
	return fmt.Sprintf("enabled=%v", v.AITTSEnabled)
}

type SrsAssistant struct {
	// Whether enable the AI assistant.
	Assistant bool `json:"assistant"`
	// The AI name.
	AIName string `json:"aiName"`
	// The AI assistant provider.
	SrsAssistantProvider
	// The AI assistant ASR.
	SrsAssistantASR
	// The AI assistant chat.
	SrsAssistantChat
	// The AI assistant post.
	SrsAssistantPost
	// The AI assistant TTS.
	SrsAssistantTTS
}

func NewAssistant(opts ...func(*SrsAssistant)) *SrsAssistant {
	v := &SrsAssistant{}

	v.AIASREnabled = true
	v.AIChatEnabled = true
	v.AIPostEnabled = false
	v.AITTSEnabled = true

	for _, opt := range opts {
		opt(v)
	}

	return v
}

func (v *SrsAssistant) String() string {
	return fmt.Sprintf("assistant=%v, name=%v, provider=<%v>, asr=<%v>, chat=<%v>, post=<%v>, tts=<%v>",
		v.Assistant, v.AIName, v.SrsAssistantProvider.String(), v.SrsAssistantASR.String(), v.SrsAssistantChat.String(),
		v.SrsAssistantPost.String(), v.SrsAssistantTTS.String(),
	)
}
