package cicd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// ReviewRequest 发送给 DeepSeek API 的请求数据结构
type ReviewRequest struct {
	Code     string `json:"code"`     // 需要审查的源代码
	Language string `json:"language"` // 编程语言(该项目为Golang)
	Config   struct {
		Strictness string `json:"strictness,omitempty"` // 审查严格度: low, medium, high
	} `json:"config,omitempty"`
}

// CodeReviewer 定义代码审查接口
type CodeReviewer_dp interface {
	Review(ctx context.Context, code string) (*ReviewResponse_dp, error)
}

// DeepSeekReviewer 实现CodeReviewer接口
type DeepSeekReviewer struct {
	APIKey string
}

// NewDeepSeekReviewer 创建一个新的DeepSeekReviewer实例
func NewDeepSeekReviewer(apiKey string) *DeepSeekReviewer {
	return &DeepSeekReviewer{
		APIKey: apiKey,
	}
}

// Review 实现CodeReviewer接口的Review方法
func (dr *DeepSeekReviewer) Review(ctx context.Context, code string) (*ReviewResponse_dp, error) {
	return requestDeepSeekReview(ctx, code, dr.APIKey)
}

// ReviewResponse 从 DeepSeek API 返回的响应数据结构
type ReviewResponse_dp struct {
	Issues []struct {
		Type       string `json:"type"`                 // 问题类型(如安全、性能等)
		Message    string `json:"message"`              // 问题描述
		Line       int    `json:"line"`                 // 问题所在行号
		Severity   string `json:"severity"`             // 问题严重程度(如warning, error)
		Suggestion string `json:"suggestion,omitempty"` // 改进建议(可选)
	} `json:"issues"` // 发现的问题列表
}

// requestDeepSeekReview 发送代码到 DeepSeek API 进行审查
// 参数: ctx - 上下文，用于超时控制
// 参数: code - 需要审查的源代码字符串
// 参数: apiKey - DeepSeek API密钥
// 返回值: 审查结果和错误信息
func requestDeepSeekReview(ctx context.Context, code string, apiKey string) (*ReviewResponse_dp, error) {
	// 准备请求数据
	reqData := ReviewRequest{
		Code:     code,
		Language: "go", // 指定为Golang
	}
	reqData.Config.Strictness = "high" // 设置审查严格度为高

	// 将请求数据编码为JSON
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.deepseek.com/v1/code/review", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 30 * time.Second, // 设置客户端超时
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close() // 确保响应体最终关闭

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		// 读取错误响应体以获取更多信息
		errBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("API returned non-200 status code: %d, body: %s",
			resp.StatusCode, string(errBody))
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 解析JSON响应到结构体
	var reviewResp ReviewResponse_dp
	err = json.Unmarshal(body, &reviewResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, body: %s", err, string(body))
	}

	// 返回审查结果
	return &reviewResp, nil
}

// RunCodeReview 运行代码审查的公共函数
// 用于被其他程序调用
func RunCodeReview(ctx context.Context, code string) (*ReviewResponse_dp, error) {
	// 从环境变量获取API密钥
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("DEEPSEEK_API_KEY environment variable is not set")
	}

	// 创建DeepSeekReviewer实例
	reviewer := NewDeepSeekReviewer(apiKey)

	// 执行代码审查
	return reviewer.Review(ctx, code)
}

func main() {
	// 从环境变量获取API密钥
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: DEEPSEEK_API_KEY environment variable is not set")
		fmt.Println("Please set the API key before running this program")
		os.Exit(1)
	}

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 示例Go代码，将被发送审查
	sampleCode := `package cicd

import "fmt"

func main() {
	fmt.Println("PaiPai")
}`

	// 使用DeepSeekReviewer执行代码审查
	reviewer := NewDeepSeekReviewer(apiKey)
	review, err := reviewer.Review(ctx, sampleCode)
	if err != nil {
		// 如果审查过程中出错，打印错误并退出
		fmt.Printf("Error: %v\n", err)
		return
	}

	// 打印审查结果
	if len(review.Issues) == 0 {
		fmt.Println("No issues found in the code review.")
	} else {
		fmt.Printf("Found %d issues:\n", len(review.Issues))
		for i, issue := range review.Issues {
			// 打印每个问题的基本信息: [严重程度] 行号: 问题描述
			fmt.Printf("%d. [%s] Line %d: %s\n", i+1, issue.Severity, issue.Line, issue.Message)
			// 如果有改进建议，打印出来
			if issue.Suggestion != "" {
				fmt.Printf("   Suggestion: %s\n", issue.Suggestion)
			}
		}
	}
}
