package code_review

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ReviewRequest 发送给 DeepSeek API 的请求数据结构
type ReviewRequest struct {
	Code     string `json:"code"`     // 需要审查的源代码
	Language string `json:"language"` // 编程语言(该项目为Golang)
	Config   struct {
		Strictness string `json:"strictness,omitempty"` // 审查严格度: low, medium, high
	} `json:"config,omitempty"`
}

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

// requestDeepSeekReview 发送代码到 DeepSeek API 进行审查
// 参数: code - 需要审查的源代码字符串
// 返回值: 审查结果和错误信息
func requestDeepSeekReview(code string) (*ReviewResponse, error) {
	// 准备请求数据
	reqData := ReviewRequest{
		Code:     code,
		Language: "go", // 指定为Golang
	}
	reqData.Config.Strictness = "high" // 设置审查严格度为高

	// 将请求数据编码为JSON
	reqBody, _ := json.Marshal(reqData)

	// 发送HTTP POST请求到DeepSeek API
	resp, err := http.Post(
		"https://api.deepseek.com/v1/code/review", // DeepSeek API端点
		"application/json",                        // 内容类型
		bytes.NewBuffer(reqBody),                  // 请求体
	)
	if err != nil {
		return nil, err // 返回请求错误
	}
	defer resp.Body.Close() // 确保响应体最终关闭

	// 读取响应体
	body, _ := ioutil.ReadAll(resp.Body)

	// 解析JSON响应到结构体
	var reviewResp ReviewResponse
	err = json.Unmarshal(body, &reviewResp)
	if err != nil {
		return nil, err // 返回JSON解析错误
	}

	// 返回审查结果
	return &reviewResp, nil
}

func main() {
	// 示例Go代码，将被发送审查
	sampleCode := `package main

				   import "fmt"

                   func main() {
	               		fmt.Println("PaiPai")
                   }`

	// 调用DeepSeek审查API
	review, err := requestDeepSeekReview(sampleCode)
	if err != nil {
		// 如果审查过程中出错，打印错误并退出
		fmt.Printf("Error: %v\n", err)
		return
	}

	// 打印审查结果
	for _, issue := range review.Issues {
		// 打印每个问题的基本信息: [严重程度] 行号: 问题描述
		fmt.Printf("[%s] Line %d: %s\n", issue.Severity, issue.Line, issue.Message)
		// 如果有改进建议，打印出来
		if issue.Suggestion != "" {
			fmt.Printf("Suggestion: %s\n", issue.Suggestion)
		}
	}
}
