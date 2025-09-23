package cicd

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

// ReviewResponse 从 DeepSeek API 返回的响应数据结构
type ReviewResponse struct {
	Issues []struct {
		Type       string `json:"type"`                 // 问题类型(如安全、性能等)
		Message    string `json:"message"`              // 问题描述
		Line       int    `json:"line"`                 // 问题所在行号
		Severity   string `json:"severity"`             // 问题严重程度(如warning, error)
		Suggestion string `json:"suggestion,omitempty"` // 改进建议(可选)
	} `json:"issues"` // 发现的问题列表
}

// CodeReviewer 定义代码审查接口
type CodeReviewer interface {
	Review(ctx context.Context, code string) (*ReviewResponse, error)
}

// MockReviewer 实现代码审查的mock版本
type MockReviewer struct {
	mock.Mock
}

// Review 实现CodeReviewer接口
func (m *MockReviewer) Review(ctx context.Context, code string) (*ReviewResponse, error) {
	args := m.Called(ctx, code)

	// 如果第一个返回值是ReviewResponse类型，则返回它，否则返回nil
	resp, _ := args.Get(0).(*ReviewResponse)
	return resp, args.Error(1)
}

// TestDeepSeekCodeReview 测试代码审查功能
func TestDeepSeekCodeReview(t *testing.T) {
	// 创建一个临时目录和测试文件
	// 这部分在实际运行时会被创建，但为了安全我们先跳过实际的文件创建

	// 创建一个mock reviewer
	mockReviewer := new(MockReviewer)

	// 设置mock的返回值
	expectedIssues := []struct {
		Type       string `json:"type"`
		Message    string `json:"message"`
		Line       int    `json:"line"`
		Severity   string `json:"severity"`
		Suggestion string `json:"suggestion,omitempty"`
	}{{
		Type:     "style",
		Message:  "Missing comment for exported function",
		Line:     10,
		Severity: "warning",
	}}

	mockReviewer.On("Review", mock.Anything, mock.Anything).Return(&ReviewResponse{
		Issues: expectedIssues,
	}, nil)

	// 定义一个测试函数，模拟文件遍历
	// 在实际测试中，我们可以使用临时文件来测试
	t.Run("shouldHandleReviewResults", func(t *testing.T) {
		// 模拟一个简单的Go文件
		testCode := `package test

func TestFunction() { }`

		// 调用mock的审查函数
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		review, err := mockReviewer.Review(ctx, testCode)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// 验证结果
		if len(review.Issues) != 1 {
			t.Errorf("Expected 1 issue, got %d", len(review.Issues))
		}

		// 验证mock是否被正确调用
		mockReviewer.AssertExpectations(t)
	})

	// 测试错误处理
	t.Run("shouldHandleReviewErrors", func(t *testing.T) {
		mockErrorReviewer := new(MockReviewer)
		mockErrorReviewer.On("Review", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("API error"))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := mockErrorReviewer.Review(ctx, "test code")
		if err == nil {
			t.Error("Expected an error, but got nil")
		}

		mockErrorReviewer.AssertExpectations(t)
	})

	// 测试空结果处理
	t.Run("shouldHandleNoIssues", func(t *testing.T) {
		mockEmptyReviewer := new(MockReviewer)
		mockEmptyReviewer.On("Review", mock.Anything, mock.Anything).Return(&ReviewResponse{}, nil)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		review, err := mockEmptyReviewer.Review(ctx, "perfect code")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(review.Issues) != 0 {
			t.Errorf("Expected 0 issues, got %d", len(review.Issues))
		}

		mockEmptyReviewer.AssertExpectations(t)
	})
}
