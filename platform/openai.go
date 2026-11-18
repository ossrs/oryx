// Copyright (c) 2022-2026 Winlin
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
	AIProvider     string `json:"aiProvider"`
	AISecretKey    string `json:"aiSecretKey"`
	AIOrganization string `json:"aiOrganization"`
	AIBaseURL      string `json:"aiBaseURL"`
}

func (v *SrsAssistantProvider) String() string {
	return fmt.Sprintf("provider=%v, secretKey=%vB, baseURL=%v",
		v.AIProvider, len(v.AISecretKey), v.AIBaseURL)
}

type SrsAssistantASR struct {
	AIASREnabled  bool   `json:"aiAsrEnabled"`
	AIASRLanguage string `json:"aiAsrLanguage"`
	AIASRPrompt   string `json:"aiAsrPrompt"`
}

func (v *SrsAssistantASR) String() string {
	return fmt.Sprintf("enabled=%v,language=%v,prompt=%v",
		v.AIASREnabled, v.AIASRLanguage, v.AIASRPrompt)
}

type SrsAssistantChat struct {
	AIChatEnabled   bool   `json:"aiChatEnabled"`
	AIChatModel     string `json:"aiChatModel"`
	AIChatPrompt    string `json:"aiChatPrompt"`
	AIChatMaxWindow int    `json:"aiChatMaxWindow"`
	AIChatMaxWords  int    `json:"aiChatMaxWords"`
}

func (v *SrsAssistantChat) String() string {
	return fmt.Sprintf("enabled=%v,model=%v,prompt=%v,window=%v,words=%v",
		v.AIChatEnabled, v.AIChatModel, v.AIChatPrompt, v.AIChatMaxWindow, v.AIChatMaxWords)
}

type SrsAssistantPost struct {
	AIPostEnabled   bool   `json:"aiPostEnabled"`
	AIPostModel     string `json:"aiPostModel"`
	AIPostPrompt    string `json:"aiPostPrompt"`
	AIPostMaxWindow int    `json:"aiPostMaxWindow"`
	AIPostMaxWords  int    `json:"aiPostMaxWords"`
}

func (v *SrsAssistantPost) String() string {
	return fmt.Sprintf("enabled=%v,model=%v,prompt=%v,window=%v,words=%v",
		v.AIPostEnabled, v.AIPostModel, v.AIPostPrompt, v.AIPostMaxWindow, v.AIPostMaxWords)
}

type SrsAssistantTTS struct {
	AITTSEnabled bool `json:"aiTtsEnabled"`
}

func (v *SrsAssistantTTS) String() string {
	return fmt.Sprintf("enabled=%v", v.AITTSEnabled)
}

type SrsAssistant struct {
	Assistant bool   `json:"assistant"`
	AIName    string `json:"aiName"`
	SrsAssistantProvider
	SrsAssistantASR
	SrsAssistantChat
	SrsAssistantPost
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
