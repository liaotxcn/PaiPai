package cicd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ReviewResponse 从 DeepSeek API 返回的响应数据结构
type ReviewResponse struct {
	Issues []struct {
		Type     string `json:"type"`
		Message  string `json:"message"`
		Line     int    `json:"line"`
		Severity string `json:"severity"`
	} `json:"issues"`
}

// requestDeepSeekReview
func requestDeepSeekReview(code string) (*ReviewResponse, error) {
	// 实际API调用代码
	// 此处为示例返回一个空的ReviewResponse
	return &ReviewResponse{}, nil
}

func TestDeepSeekCodeReview(t *testing.T) {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			content, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("Error reading file %s: %v", path, err)
				return nil
			}

			// 调用 DeepSeek API
			review, err := requestDeepSeekReview(string(content))
			if err != nil {
				t.Errorf("Error reviewing file %s: %v", path, err)
				return nil
			}

			for _, issue := range review.Issues {
				t.Errorf("%s:%d [%s] %s", path, issue.Line, issue.Severity, issue.Message)
			}
		}
		return nil
	})

	if err != nil {
		t.Errorf("Error walking directory: %v", err)
	}
}
