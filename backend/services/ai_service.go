package services

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

type AIService struct {
	client *openai.Client
}

func NewAIService() *AIService {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		// デモ用のダミークライアント
		return &AIService{client: nil}
	}

	return &AIService{
		client: openai.NewClient(apiKey),
	}
}

// AnalyzeCode コード解析を実行
func (ai *AIService) AnalyzeCode(code, language string) (string, error) {
	if ai.client == nil {
		return ai.mockAnalysis("code_analysis", code, language), nil
	}

	prompt := fmt.Sprintf(`
以下の%sコードを解析して、以下の情報をJSON形式で提供してください：

1. コードの概要と目的
2. 主要な関数・メソッドの一覧
3. 使用されているデザインパターン
4. 潜在的な問題点や改善提案
5. 依存関係の分析

コード：
%s

JSON形式で回答してください。
`, language, code)

	resp, err := ai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 2000,
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

// GenerateDocumentation ドキュメント生成
func (ai *AIService) GenerateDocumentation(code, language string) (string, error) {
	if ai.client == nil {
		return ai.mockAnalysis("documentation", code, language), nil
	}

	prompt := fmt.Sprintf(`
以下の%sコードの技術文書を作成してください。以下の要素を含めてください：

1. API仕様（関数・メソッドの説明）
2. アーキテクチャ概要
3. 使用方法の例
4. 設定方法
5. トラブルシューティング

コード：
%s

Markdown形式で回答してください。
`, language, code)

	resp, err := ai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 3000,
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

// DetectPatterns パターン検出
func (ai *AIService) DetectPatterns(code, language string) (string, error) {
	if ai.client == nil {
		return ai.mockAnalysis("pattern_detection", code, language), nil
	}

	prompt := fmt.Sprintf(`
以下の%sコードを分析して、使用されているデザインパターンやアンチパターンを特定してください：

1. デザインパターン（Singleton, Factory, Observer, etc.）
2. アンチパターン（God Object, Spaghetti Code, etc.）
3. コード品質の評価
4. リファクタリング提案

コード：
%s

JSON形式で回答してください。
`, language, code)

	resp, err := ai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 2000,
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

// AnalyzeDependencies 依存関係分析
func (ai *AIService) AnalyzeDependencies(files []FileInfo) (string, error) {
	if ai.client == nil {
		return ai.mockDependencyAnalysis(files), nil
	}

	// ファイル情報を整理
	fileList := "ファイル一覧:\n"
	for _, file := range files {
		fileList += fmt.Sprintf("- %s (%s)\n", file.Name, file.Language)
	}

	prompt := fmt.Sprintf(`
以下のファイル群の依存関係を分析して、プロジェクト構造を可視化してください：

%s

以下の情報をJSON形式で提供してください：
1. ファイル間の依存関係マップ
2. モジュール構造の分析
3. 循環依存の検出
4. アーキテクチャの改善提案

JSON形式で回答してください。
`, fileList)

	resp, err := ai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 2000,
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

type FileInfo struct {
	Name     string
	Language string
	Content  string
}

// モック関数（デモ用）
func (ai *AIService) mockAnalysis(analysisType, code, language string) string {
	switch analysisType {
	case "code_analysis":
		return `{
  "overview": "コードの概要分析結果（デモ）",
  "functions": ["function1", "function2"],
  "patterns": ["MVC", "Singleton"],
  "issues": ["潜在的な問題1", "潜在的な問題2"],
  "dependencies": ["external_lib1", "external_lib2"]
}`
	case "documentation":
		return "# API Documentation (Demo)\n\n## Overview\nこのは" + language + "で書かれたコードのドキュメントです。\n\n## Functions\n- function1(): 機能1の説明\n- function2(): 機能2の説明\n\n## Usage\n```\n// 使用例\n```"
	case "pattern_detection":
		return `{
  "design_patterns": ["Singleton", "Factory"],
  "anti_patterns": [],
  "code_quality": "Good",
  "refactoring_suggestions": ["提案1", "提案2"]
}`
	default:
		return "分析結果（デモ）"
	}
}

func (ai *AIService) mockDependencyAnalysis(files []FileInfo) string {
	return `{
  "dependency_map": {
    "file1.go": ["file2.go"],
    "file2.go": ["file3.go"]
  },
  "modules": ["module1", "module2"],
  "circular_dependencies": [],
  "architecture_suggestions": ["提案1", "提案2"]
}`
}
