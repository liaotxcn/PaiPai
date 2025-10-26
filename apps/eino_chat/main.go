package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Config 服务配置
type Config struct {
	Port              int    `json:"port"`
	KnowledgeBase     string `json:"knowledge_base"`
	RedisAddress      string `json:"redis_address"`
	RedisPassword     string `json:"redis_password"`
	ARKAPIKey         string `json:"ark_api_key"`
	ARKChatModel      string `json:"ark_chat_model"`
	ARKAPIBaseURL     string `json:"ark_api_base_url"`
	ARKEmbeddingModel string `json:"ark_embedding_model"`
}

// 聊天消息结构
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func main() {
	// 解析命令行参数
	var configPath string
	flag.StringVar(&configPath, "config", ".env", "配置文件路径")
	flag.Parse()

	// 加载环境变量
	if err := godotenv.Load(configPath); err != nil {
		log.Printf("警告: 无法加载配置文件 %s: %v", configPath, err)
		log.Println("使用默认配置继续...")
	}

	// 设置必要的环境变量
	err := os.Setenv("REDIS_ADDR", getEnv("REDIS_ADDRESS", "localhost:6379"))
	if err != nil {
		log.Printf("设置REDIS_ADDR环境变量失败: %v", err)
	}

	err = os.Setenv("ARK_API_KEY", getEnv("ARK_API_KEY", ""))
	if err != nil {
		log.Printf("设置ARK_API_KEY环境变量失败: %v", err)
	}

	err = os.Setenv("ARK_CHAT_MODEL", getEnv("ARK_CHAT_MODEL", "doubao-pro-32k"))
	if err != nil {
		log.Printf("设置ARK_CHAT_MODEL环境变量失败: %v", err)
	}

	err = os.Setenv("ARK_API_BASE_URL", getEnv("ARK_API_BASE_URL", "https://ark.cn-xxxxxx.volces.com/api/v3"))
	if err != nil {
		log.Printf("设置ARK_API_BASE_URL环境变量失败: %v", err)
	}

	err = os.Setenv("ARK_EMBEDDING_MODEL", getEnv("ARK_EMBEDDING_MODEL", "embedding-v2"))
	if err != nil {
		log.Printf("设置ARK_EMBEDDING_MODEL环境变量失败: %v", err)
	}

	// 上下文可以在后续实现真实RAG功能时使用
	// ctx := context.Background()

	// 输出提示信息
	log.Println("注意：由于缺少knowledge_service的具体实现，将使用模拟实现")
	log.Println("请确保ARK_API_KEY环境变量已正确设置")

	// 由于缺少完整的RAG实现依赖，这里提供一个简化版本的交互式聊天
	log.Println("Eino Chat服务启动成功")
	log.Println("输入 'quit' 或 'exit' 退出服务")
	log.Println("注意：当前使用模拟响应，如需真实RAG功能，请确保正确配置环境变量")
	log.Println("输入问题开始对话...")

	// 启动简化版聊天
	startSimpleChat()
}

// 启动简化版交互式聊天
func startSimpleChat() {
	messages := []ChatMessage{}
	for {
		fmt.Print("用户: ")
		var userInput string
		fmt.Scanln(&userInput)

		if userInput == "quit" || userInput == "exit" {
			log.Println("服务已停止")
			break
		}

		// 添加用户消息
		messages = append(messages, ChatMessage{Role: "user", Content: userInput})

		// 生成模拟回答
		response := generateMockResponse(userInput)

		// 添加助手消息
		messages = append(messages, ChatMessage{Role: "assistant", Content: response})

		// 显示回答
		fmt.Printf("Eino: %s\n", response)

		// 保存对话到文件
		saveChatHistory(messages)
	}
}

// 生成模拟回答
func generateMockResponse(query string) string {
	// 简单的模拟回答逻辑
	responses := []string{
		"这是一个基于RAG技术的回答。要使用真实功能，请确保正确配置ARK_API_KEY环境变量。",
		"抱歉，当前使用的是模拟响应模式。请检查.env文件中的配置项是否正确。",
		"感谢您的问题！由于缺少必要的API密钥配置，我暂时只能提供模拟回答。",
		"要启用完整的RAG功能，您需要在.env文件中设置有效的ARK_API_KEY和相关参数。",
		"我是Eino Chat助手。当前使用简化模式运行，如需完整功能，请配置必要的环境变量。",
	}

	// 简单的查询哈希选择不同回答
	hash := 0
	for _, char := range query {
		hash += int(char)
	}

	return responses[hash%len(responses)]
}

// 保存聊天历史
func saveChatHistory(messages []ChatMessage) {
	historyDir := filepath.Join("data", "history")
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		log.Printf("创建历史记录目录失败: %v", err)
		return
	}

	historyFile := filepath.Join(historyDir, "chat_history.json")
	file, err := os.OpenFile(historyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("打开历史记录文件失败: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(messages); err != nil {
		log.Printf("保存历史记录失败: %v", err)
	}
}

// 从环境变量获取字符串值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// 从环境变量获取整数值
func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		var value int
		if _, err := fmt.Sscanf(valueStr, "%d", &value); err == nil {
			return value
		}
	}
	return defaultValue
}
