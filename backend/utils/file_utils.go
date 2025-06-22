package utils

import (
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// DetectLanguage ファイル名から言語を推測
func DetectLanguage(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	languageMap := map[string]string{
		".go":         "go",
		".js":         "javascript",
		".ts":         "typescript",
		".jsx":        "javascript",
		".tsx":        "typescript",
		".py":         "python",
		".java":       "java",
		".c":          "c",
		".cpp":        "cpp",
		".cc":         "cpp",
		".cxx":        "cpp",
		".h":          "c",
		".hpp":        "cpp",
		".cs":         "csharp",
		".php":        "php",
		".rb":         "ruby",
		".rs":         "rust",
		".swift":      "swift",
		".kt":         "kotlin",
		".scala":      "scala",
		".r":          "r",
		".sql":        "sql",
		".sh":         "shell",
		".bash":       "shell",
		".zsh":        "shell",
		".ps1":        "powershell",
		".html":       "html",
		".htm":        "html",
		".css":        "css",
		".scss":       "scss",
		".sass":       "sass",
		".json":       "json",
		".xml":        "xml",
		".yaml":       "yaml",
		".yml":        "yaml",
		".toml":       "toml",
		".md":         "markdown",
		".txt":        "text",
		".dockerfile": "dockerfile",
	}

	if lang, exists := languageMap[ext]; exists {
		return lang
	}

	// Dockerfileなどの特別なファイル名をチェック
	baseName := strings.ToLower(filepath.Base(filename))
	if baseName == "dockerfile" || strings.HasPrefix(baseName, "dockerfile.") {
		return "dockerfile"
	}
	if baseName == "makefile" || strings.HasPrefix(baseName, "makefile.") {
		return "makefile"
	}

	return "unknown"
}

// IsTextFile バイト配列がテキストファイルかどうかを判定
func IsTextFile(content []byte) bool {
	// 空ファイルはテキストとして扱う
	if len(content) == 0 {
		return true
	}

	// UTF-8として有効かチェック
	if !utf8.Valid(content) {
		return false
	}

	// バイナリファイルの特徴的なバイト列をチェック
	binaryIndicators := [][]byte{
		{0x00},                   // NULL byte
		{0xFF, 0xFE},             // UTF-16 LE BOM
		{0xFE, 0xFF},             // UTF-16 BE BOM
		{0x7F, 0x45, 0x4C, 0x46}, // ELF binary
		{0x4D, 0x5A},             // Windows PE/COFF
		{0x50, 0x4B},             // ZIP/JAR
		{0x89, 0x50, 0x4E, 0x47}, // PNG
		{0xFF, 0xD8, 0xFF},       // JPEG
		{0x47, 0x49, 0x46},       // GIF
	}

	for _, indicator := range binaryIndicators {
		if len(content) >= len(indicator) {
			matches := true
			for i, b := range indicator {
				if content[i] != b {
					matches = false
					break
				}
			}
			if matches {
				return false
			}
		}
	}

	// 制御文字が多い場合はバイナリと判定
	controlCount := 0
	sampleSize := 512
	if len(content) < sampleSize {
		sampleSize = len(content)
	}

	for i := 0; i < sampleSize; i++ {
		b := content[i]
		// タブ、改行、復帰以外の制御文字をカウント
		if b < 32 && b != 9 && b != 10 && b != 13 {
			controlCount++
		}
	}

	// 制御文字が30%以上の場合はバイナリと判定
	return float64(controlCount)/float64(sampleSize) < 0.3
}

// GetFileSize ファイルサイズを人間が読みやすい形式で返す
func GetFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return string(rune(bytes)) + " B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"KB", "MB", "GB", "TB", "PB"}
	return string(rune(bytes/div)) + " " + units[exp]
}

// SanitizeFilename ファイル名を安全にする
func SanitizeFilename(filename string) string {
	// 危険な文字を除去
	filename = strings.ReplaceAll(filename, "..", "")
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	filename = strings.ReplaceAll(filename, ":", "_")
	filename = strings.ReplaceAll(filename, "*", "_")
	filename = strings.ReplaceAll(filename, "?", "_")
	filename = strings.ReplaceAll(filename, "\"", "_")
	filename = strings.ReplaceAll(filename, "<", "_")
	filename = strings.ReplaceAll(filename, ">", "_")
	filename = strings.ReplaceAll(filename, "|", "_")

	// 長すぎる場合は切り詰める
	if len(filename) > 255 {
		ext := filepath.Ext(filename)
		base := strings.TrimSuffix(filename, ext)
		if len(base) > 255-len(ext) {
			base = base[:255-len(ext)]
		}
		filename = base + ext
	}

	return filename
}
