/*
 * Copyright 2025 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"PaiPai/apps/eino_chat/eino/knowledge_service"
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/redis/go-redis/v9"
)

func init() {
	// 检查环境变量，至少需要配置一种AI服务(.env)
	hasArk := os.Getenv("ARK_API_KEY") != "" && os.Getenv("ARK_EMBEDDING_MODEL") != ""
	hasZhipu := os.Getenv("ZHIPUAI_API_KEY") != ""

	if !hasArk && !hasZhipu {
		log.Fatalf("❌ [ERROR] 至少需要配置一种AI服务:\n" +
			"1. 火山云方舟: ARK_API_KEY 和 ARK_EMBEDDING_MODEL\n" +
			"2. 智谱AI: ZHIPUAI_API_KEY")
	}
}

func main() {
	ctx := context.Background()

	err := indexMarkdownFiles(ctx, "./eino-docs")
	if err != nil {
		// 增强的错误处理
		if strings.Contains(err.Error(), "does not support RediSearch module") {
			log.Printf("错误: %v", err)
			log.Println("解决方案:")
			log.Println("1. 下载并安装Redis Stack (推荐): https://redis.io/download/")
			log.Println("2. 或者单独安装RediSearch模块并在启动Redis时加载")
			log.Println("3. 或者使用云托管的Redis服务，确保启用了RediSearch功能")
		} else {
			log.Printf("索引文档时发生错误: %v", err)
		}
		// 不再使用panic，而是优雅退出
		os.Exit(1)
	}

	log.Println("index success")
}

func indexMarkdownFiles(ctx context.Context, dir string) error {
	log.Printf("开始构建知识库索引，文档目录: %s", dir)

	runner, err := knowledge_service.BuildKnowledgeIndexing(ctx)
	if err != nil {
		log.Printf("构建索引图失败: %v", err)
		return fmt.Errorf("build index graph failed: %w", err)
	}

	// 遍历 dir 下的所有 markdown 文件
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk dir failed: %w", err)
		}
		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".md") {
			fmt.Printf("[skip] not a markdown file: %s\n", path)
			return nil
		}

		fmt.Printf("[start] indexing file: %s\n", path)

		ids, err := runner.Invoke(ctx, document.Source{URI: path})
		if err != nil {
			// 记录详细错误但继续处理其他文件
			log.Printf("[error] 处理文件 %s 时出错: %v", path, err)
			// 不返回错误，继续处理其他文件
			return nil
		}

		fmt.Printf("[done] indexing file: %s, len of parts: %d\n", path, len(ids))

		return nil
	})

	return err
}

type RedisVectorStoreConfig struct {
	RedisKeyPrefix string
	IndexName      string
	Embedding      embedding.Embedder
	Dimension      int
	MinScore       float64
	RedisAddr      string
}

func initVectorIndex(ctx context.Context, config *RedisVectorStoreConfig) (err error) {
	if config.Embedding == nil {
		return fmt.Errorf("embedding cannot be nil")
	}
	if config.Dimension <= 0 {
		return fmt.Errorf("dimension must be positive")
	}

	client := redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
	})

	// 确保在错误时关闭连接
	defer func() {
		if err != nil {
			client.Close()
		}
	}()

	if err = client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	indexName := fmt.Sprintf("%s%s", config.RedisKeyPrefix, config.IndexName)

	// 检查是否支持RediSearch模块
	_, err = client.Do(ctx, "FT.INFO", indexName).Result()
	if err != nil {
		if strings.Contains(err.Error(), "unknown command") {
			log.Println("警告: Redis服务器不支持RediSearch模块，无法创建向量索引")
			log.Println("请使用带有RediSearch模块的Redis版本，或者使用Redis Stack")
			return fmt.Errorf("redis does not support RediSearch module")
		}
		// 如果是索引不存在的错误，继续执行
	}

	// 检查是否存在索引
	exists, err := client.Do(ctx, "FT.INFO", indexName).Result()
	if err != nil {
		if !strings.Contains(err.Error(), "Unknown index name") && !strings.Contains(err.Error(), "unknown command") {
			log.Printf("检查索引是否存在时出错: %v", err)
			return fmt.Errorf("failed to check if index exists: %w", err)
		}
		err = nil
	} else if exists != nil {
		log.Printf("索引 %s 已存在，跳过创建", indexName)
		return nil
	}

	// Create new index
	createIndexArgs := []interface{}{
		"FT.CREATE", indexName,
		"ON", "HASH",
		"PREFIX", "1", config.RedisKeyPrefix,
		"SCHEMA",
		"content", "TEXT",
		"metadata", "TEXT",
		"vector", "VECTOR", "FLAT",
		"6",
		"TYPE", "FLOAT32",
		"DIM", config.Dimension,
		"DISTANCE_METRIC", "COSINE",
	}

	if err = client.Do(ctx, createIndexArgs...).Err(); err != nil {
		log.Printf("创建索引失败: %v", err)
		return fmt.Errorf("failed to create index: %w", err)
	}

	// 验证索引是否创建成功
	if _, err = client.Do(ctx, "FT.INFO", indexName).Result(); err != nil {
		log.Printf("验证索引创建失败: %v", err)
		return fmt.Errorf("failed to verify index creation: %w", err)
	}

	log.Printf("Redis向量索引 %s 创建成功", indexName)
	return nil
}
